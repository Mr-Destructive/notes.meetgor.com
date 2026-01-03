package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

// handler is the AWS Lambda handler for Netlify Functions
func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	// Serve UI files (login is default at root, dashboard, editor, assets)
	if fullPath == "" || fullPath == "login" || fullPath == "dashboard" || fullPath == "editor" || 
	   fullPath == "editor.js" || fullPath == "dashboard.js" {
		page := fullPath
		if page == "" {
			page = "login"
		}
		return serveUI(page)
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
	db, err := sql.Open("libsql", dbString)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return respondError(500, "Database connection failed"), nil
	}
	defer db.Close()

	// Verify database connection
	if err := db.PingContext(ctx); err != nil {
		log.Printf("Database ping error: %v", err)
		return respondError(500, "Database connection failed"), nil
	}

	// Initialize schema if needed (soft fail - continues even if tables already exist)
	if err := initSchemaIfNotExists(ctx, db); err != nil {
		log.Printf("Schema initialization warning (non-fatal): %v", err)
		// Don't return error - tables may already exist
	}

	// Create sqlc queries
	queries := gen.New(db)

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
			Title    string                 `json:"title"`
			Slug     string                 `json:"slug"`
			Content  string                 `json:"content"`
			Excerpt  string                 `json:"excerpt"`
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

		post, err := queries.CreatePost(ctx, gen.CreatePostParams{
			ID:      postID,
			TypeID:  postReq.TypeID,
			Title:   postReq.Title,
			Slug:    postReq.Slug,
			Content: postReq.Content,
			Excerpt: sql.NullString{String: postReq.Excerpt, Valid: postReq.Excerpt != ""},
			Status:  sql.NullString{String: postReq.Status, Valid: postReq.Status != ""},
			Tags:    sql.NullString{String: string(tagsJSON), Valid: len(tagsJSON) > 0},
			Metadata: sql.NullString{String: string(metaJSON), Valid: len(metaJSON) > 0},
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
		tagsJSON := sql.NullString{}
		if len(updateReq.Tags) > 0 {
			b, _ := json.Marshal(updateReq.Tags)
			tagsJSON = sql.NullString{String: string(b), Valid: true}
		}

		metaJSON := sql.NullString{}
		if updateReq.Metadata != nil {
			b, _ := json.Marshal(updateReq.Metadata)
			metaJSON = sql.NullString{String: string(b), Valid: true}
		}

		// Helper function to convert *string to sql.NullString
		stringToNull := func(s *string) sql.NullString {
			if s == nil {
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

		post, err := queries.UpdatePost(ctx, gen.UpdatePostParams{
			ID:         id,
			Title:      stringToNull(updateReq.Title),
			Slug:       stringToNull(updateReq.Slug),
			Content:    stringToNull(updateReq.Content),
			Excerpt:    stringToNull(updateReq.Excerpt),
			Status:     stringToNull(updateReq.Status),
			Tags:       tagsJSON,
			Metadata:   metaJSON,
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

// serveUI serves the embedded UI files
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
	lambda.Start(handler)
}
