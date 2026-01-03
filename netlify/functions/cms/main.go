package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"blog/internal/db"
	"blog/internal/handler"
	gen "blog/internal/db/gen"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

//go:embed ui/*
var uiFiles embed.FS

func init() {
	log.Println("CMS function initializing...")
}

// lambdaHandler is the AWS Lambda handler for Netlify Functions
func lambdaHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request: %s %s", req.HTTPMethod, req.Path)
	ctx := context.Background()

	// Handle CORS preflight
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
		}, nil
	}

	// Parse path early to check if it's API or UI
	fullPath := strings.TrimPrefix(req.Path, "/.netlify/functions/cms")
	fullPath = strings.Trim(fullPath, "/")

	// Serve root path with admin dashboard shell
	if fullPath == "" {
		return serveAdminShell()
	}
	
	// Serve UI files for legacy routes
	if fullPath == "login" || fullPath == "dashboard" || fullPath == "editor" || 
	   fullPath == "editor.js" || fullPath == "dashboard.js" {
		return serveUI(fullPath)
	}
	
	// Serve CSS files
	if strings.HasPrefix(fullPath, "css/") {
		if fullPath == "css/admin.css" {
			return serveAdminCSS()
		}
		return respondError(404, "CSS file not found"), nil
	}

	// Health check for API root
	if fullPath == "api" {
		return respondJSON(200, map[string]string{"status": "ok", "version": "1.0"}), nil
	}

	// Connect to Turso database
	dbName := os.Getenv("TURSO_CONNECTION_URL")
	dbToken := os.Getenv("TURSO_AUTH_TOKEN")

	if dbName == "" || dbToken == "" {
		return respondError(500, "Database credentials not configured"), nil
	}

	dbString := fmt.Sprintf("%s?authToken=%s", dbName, dbToken)
	sqldb, err := sql.Open("libsql", dbString)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return respondError(500, "Database connection failed"), nil
	}
	defer sqldb.Close()

	// Verify database connection
	if err := sqldb.PingContext(ctx); err != nil {
		log.Printf("Database ping error: %v", err)
		return respondError(500, "Database connection failed"), nil
	}

	// Initialize schema if needed (soft fail - continues even if tables already exist)
	if err := initSchemaIfNotExists(ctx, sqldb); err != nil {
		log.Printf("Schema initialization warning (non-fatal): %v", err)
		// Don't return error - tables may already exist
	}

	// Create sqlc queries
	queries := gen.New(sqldb)
	
	// Handle /admin/* routes with database access
	if strings.HasPrefix(fullPath, "admin/") {
		return handleAdminRoute(ctx, req, fullPath, sqldb)
	}

	// Ensure specific query tables are initialized
	if err := queries.InitPostTables(ctx); err != nil {
		log.Printf("Post tables init warning (non-fatal): %v", err)
	}
	if err := queries.InitSeriesTables(ctx); err != nil {
		log.Printf("Series tables init warning (non-fatal): %v", err)
	}

	// API routes
	path := strings.TrimPrefix(fullPath, "api")
	path = strings.Trim(path, "/")

	parts := strings.Split(path, "/")
	resource := parts[0]
	var id, action string
	if len(parts) > 1 {
		id = parts[1]
	}
	if len(parts) > 2 {
		action = parts[2]
	}

	// Route to handlers
	switch resource {
	case "auth":
		return handleAuth(req, ctx, id)
	case "posts":
		return handlePosts(req, ctx, queries, id, action)
	case "types":
		return handleTypes(req, ctx, queries)
	case "tags":
		return handleTags(req, ctx, queries)
	case "exports":
		return handleExports(req, ctx, queries)
	default:
		return respondError(404, "Resource not found"), nil
	}
}

// handleAuth handles authentication endpoints
func handleAuth(req events.APIGatewayProxyRequest, ctx context.Context, action string) (events.APIGatewayProxyResponse, error) {
	switch action {
	case "login":
		if req.HTTPMethod != "POST" {
			return respondError(405, "Method not allowed"), nil
		}
		var loginReq struct {
			Password string `json:"password"`
		}
		if err := json.NewDecoder(strings.NewReader(req.Body)).Decode(&loginReq); err != nil {
			return respondError(400, "Invalid request"), nil
		}

		// Accept any non-empty password for now (set ADMIN_PASSWORD env var in Netlify UI for security)
		if loginReq.Password == "" {
			return respondError(401, "Password required"), nil
		}
		
		// For development/demo: accept "test" or any password with ADMIN_PASSWORD env var
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if adminPassword != "" && loginReq.Password != adminPassword {
			return respondError(401, "Invalid credentials"), nil
		}
		
		log.Printf("Login successful for user (password: %q)", loginReq.Password)

		// Generate JWT token
		token, err := generateToken()
		if err != nil {
			log.Printf("Token generation error: %v", err)
			return respondError(500, "Failed to generate token"), nil
		}

		return respondJSON(200, map[string]interface{}{
			"token":      token,
			"expires_at": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
		}), nil

	case "verify":
		if req.HTTPMethod != "GET" {
			return respondError(405, "Method not allowed"), nil
		}
		token := getToken(req)
		if token == "" {
			return respondError(401, "No token"), nil
		}
		if !verifyToken(token) {
			return respondError(401, "Invalid token"), nil
		}
		return respondJSON(200, map[string]bool{"valid": true}), nil

	case "logout":
		return respondJSON(200, map[string]string{"status": "logged out"}), nil

	default:
		return respondError(404, "Auth endpoint not found"), nil
	}
}

// generateToken creates a simple JWT-like token (basic implementation for Netlify)
func generateToken() (string, error) {
	// For production, use a proper JWT library
	// This is a simple base64-encoded timestamp token for demo
	payload := fmt.Sprintf("admin:%d", time.Now().Unix())
	return base64.StdEncoding.EncodeToString([]byte(payload)), nil
}

// getToken extracts token from Authorization header or cookie
func getToken(req events.APIGatewayProxyRequest) string {
	if auth := req.Headers["Authorization"]; auth != "" {
		if len(auth) > 7 && auth[:7] == "Bearer " {
			return auth[7:]
		}
	}
	if cookies := req.Headers["Cookie"]; cookies != "" {
		for _, cookie := range strings.Split(cookies, ";") {
			if strings.Contains(cookie, "auth_token=") {
				parts := strings.Split(strings.TrimSpace(cookie), "=")
				if len(parts) == 2 {
					return parts[1]
				}
			}
		}
	}
	return ""
}

// verifyToken validates the token (basic implementation)
func verifyToken(token string) bool {
	if _, err := base64.StdEncoding.DecodeString(token); err != nil {
		return false
	}
	return true
}

// generateID creates a unique ID for posts
func generateID() string {
	// Simple ID generation - in production use UUID
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

// handlePosts handles /posts CRUD endpoints
func handlePosts(req events.APIGatewayProxyRequest, ctx context.Context, queries *gen.Queries, id, action string) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		if id != "" {
			// Get single post by ID or slug
			post, err := queries.GetPost(ctx, gen.GetPostParams{
				ID:   id,
				Slug: id,
			})
			if err != nil {
				log.Printf("GetPost error: %v", err)
				return respondError(404, "Post not found"), nil
			}
			return respondJSON(200, post), nil
		}

		// List posts with optional filters
		status := req.QueryStringParameters["status"]
		// Note: empty status means get all posts (for authenticated users viewing dashboard)
		// For public API, callers can specify status=published to get only published posts

		limit := int64(50)
		offset := int64(0)
		
		// Use nil for status if empty to get all posts
		var statusParam interface{} = nil
		if status != "" {
			statusParam = status
		}

		posts, err := queries.ListPosts(ctx, gen.ListPostsParams{
			Status: statusParam,
			TypeID: nil,
			Offset: offset,
			Limit:  limit,
		})
		if err != nil {
			log.Printf("ListPosts error: %v", err)
			return respondError(500, "Failed to fetch posts"), nil
		}

		return respondJSON(200, posts), nil

	case "POST":
		// Create new post
		var postReq struct {
			TypeID   string                 `json:"type_id"`
			Title    *string                `json:"title"`
			Slug     string                 `json:"slug"`
			Content  *string                `json:"content"`
			Excerpt  *string                `json:"excerpt"`
			Status   string                 `json:"status"`
			Tags     []string               `json:"tags"`
			Metadata map[string]interface{} `json:"metadata"`
		}
		if err := json.NewDecoder(strings.NewReader(req.Body)).Decode(&postReq); err != nil {
			log.Printf("POST decode error: %v", err)
			return respondError(400, "Invalid request body"), nil
		}

		// Generate ID
		postID := generateID()
		now := time.Now()

		// Convert tags and metadata to JSON
		tagsJSON, _ := json.Marshal(postReq.Tags)
		metaJSON, _ := json.Marshal(postReq.Metadata)

		// Helper to convert *string to sql.NullString
		stringToNull := func(s *string) sql.NullString {
			if s == nil || *s == "" {
				return sql.NullString{Valid: false}
			}
			return sql.NullString{String: *s, Valid: true}
		}

		title := ""
		if postReq.Title != nil {
			title = *postReq.Title
		}
		content := ""
		if postReq.Content != nil {
			content = *postReq.Content
		}

		post, err := queries.CreatePost(ctx, gen.CreatePostParams{
			ID:      postID,
			TypeID:  postReq.TypeID,
			Title:   title,
			Slug:    postReq.Slug,
			Content: content,
			Excerpt: stringToNull(postReq.Excerpt),
			Status:  sql.NullString{String: postReq.Status, Valid: postReq.Status != ""},
			Tags:    sql.NullString{String: string(tagsJSON), Valid: len(tagsJSON) > 0 && string(tagsJSON) != "[]"},
			Metadata: sql.NullString{String: string(metaJSON), Valid: len(metaJSON) > 0 && string(metaJSON) != "{}"},
			CreatedAt: sql.NullTime{Time: now, Valid: true},
			UpdatedAt: sql.NullTime{Time: now, Valid: true},
		})
		if err != nil {
			log.Printf("CreatePost error: %v", err)
			return respondError(500, "Failed to create post"), nil
		}

		return respondJSON(201, post), nil

	case "PUT":
		if id == "" {
			return respondError(400, "Post ID required"), nil
		}

		var updateReq struct {
			Title      *string                `json:"title"`
			Slug       *string                `json:"slug"`
			Content    *string                `json:"content"`
			Excerpt    *string                `json:"excerpt"`
			Status     *string                `json:"status"`
			Tags       []string               `json:"tags"`
			Metadata   map[string]interface{} `json:"metadata"`
			IsFeatured *bool                  `json:"is_featured"`
		}
		if err := json.NewDecoder(strings.NewReader(req.Body)).Decode(&updateReq); err != nil {
			return respondError(400, "Invalid request body"), nil
		}

		// Convert tags and metadata to JSON if provided
		tagsJSON := sql.NullString{Valid: false}
		if len(updateReq.Tags) > 0 {
			b, _ := json.Marshal(updateReq.Tags)
			jsonStr := string(b)
			if jsonStr != "[]" && jsonStr != "" {
				tagsJSON = sql.NullString{String: jsonStr, Valid: true}
			}
		}

		metaJSON := sql.NullString{Valid: false}
		if updateReq.Metadata != nil && len(updateReq.Metadata) > 0 {
			b, _ := json.Marshal(updateReq.Metadata)
			jsonStr := string(b)
			if jsonStr != "{}" && jsonStr != "" {
				metaJSON = sql.NullString{String: jsonStr, Valid: true}
			}
		}

		// Helper function to convert *string to sql.NullString for nullable fields
		stringToNullable := func(s *string) sql.NullString {
			if s == nil || *s == "" {
				return sql.NullString{Valid: false}
			}
			return sql.NullString{String: *s, Valid: true}
		}

		// Helper function to convert *bool to sql.NullBool
		boolToNull := func(b *bool) sql.NullBool {
			if b == nil {
				return sql.NullBool{Valid: false}
			}
			return sql.NullBool{Bool: *b, Valid: true}
		}

		// Title and Slug are NOT NULL in DB, so must provide values (don't change if not provided)
		title := updateReq.Title
		if title == nil || *title == "" {
			// Need to fetch current value from DB
			currentPost, err := queries.GetPost(ctx, gen.GetPostParams{ID: id, Slug: id})
			if err != nil {
				return respondError(404, "Post not found"), nil
			}
			if title == nil {
				title = &currentPost.Title
			} else if *title == "" {
				// Empty string - keep current value
				title = &currentPost.Title
			}
		}

		slug := updateReq.Slug
		if slug == nil || *slug == "" {
			currentPost, err := queries.GetPost(ctx, gen.GetPostParams{ID: id, Slug: id})
			if err != nil {
				return respondError(404, "Post not found"), nil
			}
			if slug == nil {
				slug = &currentPost.Slug
			} else if *slug == "" {
				slug = &currentPost.Slug
			}
		}

		content := updateReq.Content
		if content == nil || *content == "" {
			currentPost, err := queries.GetPost(ctx, gen.GetPostParams{ID: id, Slug: id})
			if err != nil {
				return respondError(404, "Post not found"), nil
			}
			if content == nil {
				content = &currentPost.Content
			} else if *content == "" {
				content = &currentPost.Content
			}
		}

		post, err := queries.UpdatePost(ctx, gen.UpdatePostParams{
			ID:         id,
			Title:      *title,
			Slug:       *slug,
			Content:    *content,
			Excerpt:    stringToNullable(updateReq.Excerpt),
			Status:     stringToNullable(updateReq.Status),
			Tags:       tagsJSON,
			Metadata:   metaJSON,
			PublishedAt: sql.NullTime{Valid: false}, // Don't change on update
			IsFeatured: boolToNull(updateReq.IsFeatured),
			UpdatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			log.Printf("UpdatePost error: %v", err)
			return respondError(500, "Failed to update post"), nil
		}

		return respondJSON(200, post), nil

	case "DELETE":
		if id == "" {
			return respondError(400, "Post ID required"), nil
		}

		if err := queries.DeletePost(ctx, id); err != nil {
			log.Printf("DeletePost error: %v", err)
			return respondError(500, "Failed to delete post"), nil
		}

		return respondJSON(200, map[string]string{"status": "deleted"}), nil

	default:
		return respondError(405, "Method not allowed"), nil
	}
}

// handleTypes returns post types
func handleTypes(req events.APIGatewayProxyRequest, ctx context.Context, queries *gen.Queries) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != "GET" {
		return respondError(405, "Method not allowed"), nil
	}

	types, err := queries.GetPostTypes(ctx)
	if err != nil {
		log.Printf("GetPostTypes error: %v", err)
		return respondError(500, "Failed to fetch types"), nil
	}

	return respondJSON(200, types), nil
}

// handleTags returns all posts (for now, tags would require aggregation)
func handleTags(req events.APIGatewayProxyRequest, ctx context.Context, queries *gen.Queries) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != "GET" {
		return respondError(405, "Method not allowed"), nil
	}

	// GetTags not implemented yet - return empty array
	return respondJSON(200, []map[string]interface{}{}), nil
}

// handleExports returns published posts for export
func handleExports(req events.APIGatewayProxyRequest, ctx context.Context, queries *gen.Queries) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != "GET" {
		return respondError(405, "Method not allowed"), nil
	}

	// Get published posts only
	posts, err := queries.ListPosts(ctx, gen.ListPostsParams{
		Status: "published",
		TypeID: nil,
		Offset: 0,
		Limit:  1000,
	})
	if err != nil {
		log.Printf("ListPosts error: %v", err)
		return respondError(500, "Failed to fetch posts"), nil
	}

	return respondJSON(200, posts), nil
}

// initSchemaIfNotExists creates tables if they don't exist
func initSchemaIfNotExists(ctx context.Context, db *sql.DB) error {
	// Create tables if not exists
	schema := `
	CREATE TABLE IF NOT EXISTS post_types (
	  id TEXT PRIMARY KEY,
	  name TEXT NOT NULL,
	  slug TEXT UNIQUE NOT NULL,
	  description TEXT
	);

	CREATE TABLE IF NOT EXISTS posts (
	  id TEXT PRIMARY KEY,
	  type_id TEXT NOT NULL,
	  title TEXT NOT NULL,
	  slug TEXT UNIQUE NOT NULL,
	  content TEXT NOT NULL,
	  excerpt TEXT,
	  status TEXT DEFAULT 'draft',
	  is_featured BOOLEAN DEFAULT 0,
	  tags TEXT,
	  metadata TEXT,
	  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	  published_at DATETIME,
	  FOREIGN KEY(type_id) REFERENCES post_types(id)
	);

	CREATE TABLE IF NOT EXISTS revisions (
	  id TEXT PRIMARY KEY,
	  post_id TEXT NOT NULL,
	  content TEXT NOT NULL,
	  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS series (
	  id TEXT PRIMARY KEY,
	  name TEXT NOT NULL,
	  slug TEXT UNIQUE NOT NULL,
	  description TEXT,
	  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS post_series (
	  post_id TEXT NOT NULL,
	  series_id TEXT NOT NULL,
	  order_in_series INT,
	  PRIMARY KEY(post_id, series_id),
	  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
	  FOREIGN KEY(series_id) REFERENCES series(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS settings (
	  key TEXT PRIMARY KEY,
	  value TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_posts_type ON posts(type_id);
	CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
	CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at);
	CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
	CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_series_slug ON series(slug);

	INSERT OR IGNORE INTO post_types (id, name, slug, description) VALUES
	  ('article', 'Article', 'article', 'Full-length articles'),
	  ('review', 'Review', 'review', 'Book, movie, or product reviews'),
	  ('thought', 'Thought', 'thought', 'Quick thoughts and reflections'),
	  ('link', 'Link', 'link', 'Curated links with commentary'),
	  ('til', 'TIL', 'til', 'Today I Learned'),
	  ('quote', 'Quote', 'quote', 'Quotations and excerpts'),
	  ('list', 'List', 'list', 'Curated lists'),
	  ('note', 'Note', 'note', 'Quick notes'),
	  ('snippet', 'Snippet', 'snippet', 'Code snippets'),
	  ('essay', 'Essay', 'essay', 'Long-form essays'),
	  ('tutorial', 'Tutorial', 'tutorial', 'Step-by-step guides'),
	  ('interview', 'Interview', 'interview', 'Q&A interviews');
	`

	if _, err := db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("Schema initialized")
	return nil
}

// serveAdminShell returns the admin dashboard HTML shell
func serveAdminShell() (events.APIGatewayProxyResponse, error) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Blog Admin</title>
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<link rel="stylesheet" href="/css/admin.css">
	<style>
		body {
			margin: 0;
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
			background: #f5f5f5;
		}
		.admin-shell {
			display: flex;
			height: 100vh;
			overflow: hidden;
		}
		.sidebar {
			width: 220px;
			background: #1a1a1a;
			color: white;
			padding: 20px;
			overflow-y: auto;
			border-right: 1px solid #333;
		}
		.sidebar h2 {
			margin: 0 0 30px 0;
			font-size: 18px;
			font-weight: 600;
			border-bottom: 1px solid #333;
			padding-bottom: 15px;
		}
		.nav-section {
			margin-bottom: 30px;
		}
		.nav-section-title {
			font-size: 11px;
			font-weight: 700;
			text-transform: uppercase;
			color: #888;
			margin-bottom: 12px;
			letter-spacing: 1px;
		}
		.nav-link {
			display: block;
			padding: 10px 12px;
			color: #ccc;
			text-decoration: none;
			border-radius: 4px;
			margin-bottom: 4px;
			font-size: 14px;
			cursor: pointer;
			border: none;
			background: none;
			text-align: left;
			width: 100%;
		}
		.nav-link:hover {
			background: #333;
			color: white;
		}
		.nav-link.active {
			background: #0066cc;
			color: white;
			font-weight: 600;
		}
		.main-content {
			flex: 1;
			overflow: hidden;
			display: flex;
			flex-direction: column;
		}
		.topbar {
			background: white;
			border-bottom: 1px solid #e0e0e0;
			padding: 15px 30px;
			display: flex;
			justify-content: space-between;
			align-items: center;
		}
		.topbar h1 {
			margin: 0;
			font-size: 24px;
			font-weight: 600;
		}
		.content-area {
			flex: 1;
			overflow-y: auto;
			padding: 30px;
		}
		.spinner {
			display: none;
			width: 16px;
			height: 16px;
			border: 2px solid #f3f3f3;
			border-top: 2px solid #0066cc;
			border-radius: 50%;
			animation: spin 1s linear infinite;
		}
		@keyframes spin {
			0% { transform: rotate(0deg); }
			100% { transform: rotate(360deg); }
		}
	</style>
</head>
<body>
	<div class="admin-shell">
		<div class="sidebar">
			<h2>Blog Admin</h2>
			<div class="nav-section">
				<div class="nav-section-title">Main</div>
				<a class="nav-link" hx-get="/admin/dashboard" hx-target="#main-content" onclick="updateActiveNav(this)">📊 Dashboard</a>
				<a class="nav-link" hx-get="/admin/posts" hx-target="#main-content" onclick="updateActiveNav(this)">📝 Posts</a>
				<a class="nav-link" hx-get="/admin/series" hx-target="#main-content" onclick="updateActiveNav(this)">📚 Series</a>
			</div>
			<div class="nav-section">
				<div class="nav-section-title">Config</div>
				<a class="nav-link" hx-get="/admin/types" hx-target="#main-content" onclick="updateActiveNav(this)">🏷️ Post Types</a>
				<a class="nav-link" hx-get="/admin/exports" hx-target="#main-content" onclick="updateActiveNav(this)">📤 Export</a>
			</div>
		</div>
		<div class="main-content">
			<div class="topbar">
				<h1 id="page-title">Dashboard</h1>
				<div class="spinner"></div>
			</div>
			<div class="content-area" id="main-content">
				<div class="card">
					<div class="card-body" style="text-align: center; padding: 60px 20px;">
						<p style="color: #999; font-size: 16px;">Loading...</p>
					</div>
				</div>
			</div>
		</div>
	</div>
	<script src="https://cdn.jsdelivr.net/npm/marked@11.1.1/marked.min.js"></script>
	<script>
		function updateActiveNav(element) {
			document.querySelectorAll('.nav-link').forEach(el => el.classList.remove('active'));
			element.classList.add('active');
			const text = element.textContent.trim();
			document.getElementById('page-title').textContent = text;
		}
		document.addEventListener('DOMContentLoaded', () => {
			htmx.ajax('GET', '/admin/dashboard', { target: '#main-content' });
			document.querySelector('.nav-link').classList.add('active');
		});
	</script>
</body>
</html>`
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       html,
		Headers: map[string]string{
			"Content-Type":                "text/html; charset=utf-8",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

// handleAdminRoute handles /admin/* routes with database access
func handleAdminRoute(ctx context.Context, req events.APIGatewayProxyRequest, fullPath string, sqldb *sql.DB) (events.APIGatewayProxyResponse, error) {
	path := strings.TrimPrefix(fullPath, "admin/")
	path = strings.Trim(path, "/")
	
	// Create DB wrapper
	appDB := db.FromSQL(sqldb)
	
	// Parse the admin path
	parts := strings.Split(path, "/")
	resource := parts[0]
	var id, action string
	if len(parts) > 1 {
		id = parts[1]
	}
	if len(parts) > 2 {
		action = parts[2]
	}
	
	// Create a wrapper HTTP request from the Lambda request
	httpReq, _ := http.NewRequest(req.HTTPMethod, req.Path, strings.NewReader(req.Body))
	httpReq = httpReq.WithContext(ctx)
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	for k, v := range req.QueryStringParameters {
		httpReq.URL.RawQuery += fmt.Sprintf("%s=%s&", k, v)
	}
	
	// Create response wrapper
	respWriter := &lambdaResponseWriter{
		headers: make(map[string]string),
	}
	
	// Route to handlers
	switch resource {
	case "dashboard":
		handler.HandleAdminDashboard(respWriter, httpReq, appDB)
	case "posts":
		// Check for editor routes: /posts/new or /posts/{id}/edit
		if id == "new" || action == "edit" {
			handler.HandlePostEditor(respWriter, httpReq, appDB, id)
		} else {
			handler.HandlePostsList(respWriter, httpReq, appDB)
		}
	case "series":
		// Check for editor routes: /series/new or /series/{id}/edit
		if id == "new" || action == "edit" {
			respWriter.WriteHeader(http.StatusNotImplemented)
			fmt.Fprint(respWriter, `<div class="alert alert-danger">Series editor not yet implemented</div>`)
		} else {
			handler.HandleSeriesList(respWriter, httpReq, appDB)
		}
	case "types":
		handler.HandlePostTypes(respWriter, httpReq, appDB)
	case "exports":
		handler.HandleExportPage(respWriter, httpReq, appDB)
	default:
		return respondError(404, "Admin page not found"), nil
	}
	
	// Convert response
	return events.APIGatewayProxyResponse{
		StatusCode: respWriter.statusCode,
		Body:       respWriter.body.String(),
		Headers:    respWriter.headers,
	}, nil
}

// lambdaResponseWriter implements http.ResponseWriter for Lambda
type lambdaResponseWriter struct {
	headers    map[string]string
	body       strings.Builder
	statusCode int
}

func (w *lambdaResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *lambdaResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	w.headers["Content-Type"] = "text/html; charset=utf-8"
	w.headers["Access-Control-Allow-Origin"] = "*"
	return w.body.Write(b)
}

func (w *lambdaResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// serveAdminCSS serves the admin CSS
func serveAdminCSS() (events.APIGatewayProxyResponse, error) {
	css := `* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

:root {
    --primary: #1a1a1a;
    --secondary: #666;
    --border: #e0e0e0;
    --bg-light: #f9f9f9;
    --bg-white: #ffffff;
    --text: #333;
    --text-light: #666;
    --success: #27ae60;
    --warning: #f39c12;
    --danger: #e74c3c;
}

.card {
    background: var(--bg-white);
    border: 1px solid var(--border);
    border-radius: 4px;
    margin-bottom: 24px;
    overflow: hidden;
}

.card-header {
    padding: 20px;
    border-bottom: 1px solid var(--border);
    background: var(--bg-light);
}

.card-body {
    padding: 20px;
}

.form-group {
    margin-bottom: 20px;
}

.form-group label {
    display: block;
    margin-bottom: 6px;
    font-weight: 500;
    font-size: 13px;
}

.form-group input,
.form-group textarea,
.form-group select {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid var(--border);
    border-radius: 3px;
    font-family: inherit;
    font-size: 13px;
}

.btn {
    padding: 10px 16px;
    border: 1px solid transparent;
    border-radius: 3px;
    cursor: pointer;
    font-size: 13px;
    font-weight: 500;
}

.btn-primary { background: var(--primary); color: white; }
.btn-success { background: var(--success); color: white; }
.btn-danger { background: var(--danger); color: white; }
.btn-outline { background: transparent; border-color: var(--border); color: var(--text); }

.table {
    width: 100%;
    border-collapse: collapse;
}

.table th {
    padding: 12px;
    text-align: left;
    font-weight: 600;
    border-bottom: 2px solid var(--border);
    background: var(--bg-light);
}

.table td {
    padding: 12px;
    border-bottom: 1px solid var(--border);
}

.badge {
    display: inline-flex;
    padding: 4px 10px;
    border-radius: 3px;
    font-size: 11px;
    font-weight: 600;
}

.badge-success { background: #d5f4e6; color: #27ae60; }
.badge-warning { background: #fdebd0; color: #f39c12; }
.badge-danger { background: #fadbd8; color: #e74c3c; }

.alert {
    padding: 12px 16px;
    border-radius: 3px;
    margin-bottom: 20px;
    border-left: 4px solid;
}

.alert-success { background: #d5f4e6; color: #27ae60; border-color: #27ae60; }
.alert-danger { background: #fadbd8; color: #e74c3c; border-color: #e74c3c; }`
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       css,
		Headers: map[string]string{
			"Content-Type":                "text/css; charset=utf-8",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func serveUI(page string) (events.APIGatewayProxyResponse, error) {
	var filename, contentType string
	
	switch page {
	case "dashboard":
		filename = "ui/dashboard.html"
		contentType = "text/html; charset=utf-8"
	case "dashboard.js":
		filename = "ui/dashboard.js"
		contentType = "application/javascript; charset=utf-8"
	case "editor":
		filename = "ui/editor.html"
		contentType = "text/html; charset=utf-8"
	case "editor.js":
		filename = "ui/editor.js"
		contentType = "application/javascript; charset=utf-8"
	case "login", "":
		filename = "ui/login.html"
		contentType = "text/html; charset=utf-8"
	default:
		return respondError(404, "Page not found"), nil
	}

	content, err := uiFiles.ReadFile(filename)
	if err != nil {
		log.Printf("Failed to read UI file %s: %v", filename, err)
		return respondError(500, "Failed to load page"), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(content),
		Headers: map[string]string{
			"Content-Type":                contentType,
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

// Helper functions
func respondJSON(status int, data interface{}) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}

func respondError(status int, message string) events.APIGatewayProxyResponse {
	return respondJSON(status, map[string]string{"error": message})
}

// main starts the Lambda handler
func main() {
	lambda.Start(lambdaHandler)
}
