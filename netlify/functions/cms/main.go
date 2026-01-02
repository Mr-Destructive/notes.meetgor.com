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

	// Serve UI files (login is default at root, editor, assets)
	if fullPath == "" || fullPath == "login" || fullPath == "editor" || fullPath == "editor.js" {
		return serveUI(strings.TrimPrefix(fullPath, ""))
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

		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if adminPassword == "" {
			// Default to "test" if not set
			adminPassword = "test"
		}

		log.Printf("Login attempt with password, admin password set: %v", adminPassword != "test")
		if loginReq.Password != adminPassword {
			log.Printf("Password mismatch: got %q, expected %q", loginReq.Password, adminPassword)
			return respondError(401, "Invalid credentials"), nil
		}

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

// handlePosts handles GET /posts
func handlePosts(req events.APIGatewayProxyRequest, ctx context.Context, queries *gen.Queries, id, action string) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != "GET" {
		return respondError(405, "Method not allowed"), nil
	}

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
	if status == "" {
		status = "published"
	}

	limit := int64(50)
	offset := int64(0)

	posts, err := queries.ListPosts(ctx, gen.ListPostsParams{
		Status: status,
		TypeID: nil,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		log.Printf("ListPosts error: %v", err)
		return respondError(500, "Failed to fetch posts"), nil
	}

	return respondJSON(200, posts), nil
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
