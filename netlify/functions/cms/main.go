package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	gen "blog/internal/db/gen"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func init() {
	log.Println("CMS function initializing...")
}

// handler is the AWS Lambda handler for Netlify Functions
func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request: %s %s", req.HTTPMethod, req.Path)
	ctx := context.Background()

	// Health check
	if req.Path == "/.netlify/functions/cms" || req.Path == "/.netlify/functions/cms/" {
		return respondJSON(200, map[string]string{"status": "ok", "message": "CMS function running"}), nil
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

	// Create sqlc queries
	queries := gen.New(db)

	// Parse path
	path := strings.TrimPrefix(req.Path, "/.netlify/functions/cms")
	path = strings.TrimPrefix(path, "/api")
	path = strings.Trim(path, "/")

	if path == "" {
		return respondJSON(200, map[string]string{"status": "ok", "version": "1.0"}), nil
	}

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

// handlePosts handles GET /posts
func handlePosts(req events.APIGatewayProxyRequest, ctx context.Context, queries *gen.Queries, id, action string) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != "GET" {
		return respondError(405, "Method not allowed"), nil
	}

	if id != "" {
		// Get single post
		post, err := queries.GetPost(ctx, id)
		if err != nil {
			log.Printf("GetPost error: %v", err)
			return respondError(404, "Post not found"), nil
		}
		return respondJSON(200, post), nil
	}

	// List posts with optional filters
	limit := 50
	offset := 0
	status := req.QueryStringParameters["status"]
	if status == "" {
		status = "published"
	}

	posts, err := queries.ListPosts(ctx)
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

// handleTags returns all tags
func handleTags(req events.APIGatewayProxyRequest, ctx context.Context, queries *gen.Queries) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != "GET" {
		return respondError(405, "Method not allowed"), nil
	}

	tags, err := queries.GetTags(ctx)
	if err != nil {
		log.Printf("GetTags error: %v", err)
		return respondError(500, "Failed to fetch tags"), nil
	}

	return respondJSON(200, tags), nil
}

// handleExports returns published posts for export
func handleExports(req events.APIGatewayProxyRequest, ctx context.Context, queries *gen.Queries) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != "GET" {
		return respondError(405, "Method not allowed"), nil
	}

	// Get published posts only
	posts, err := queries.ListPosts(ctx)
	if err != nil {
		log.Printf("ListPosts error: %v", err)
		return respondError(500, "Failed to fetch posts"), nil
	}

	return respondJSON(200, posts), nil
}
// Helper functions
func respondJSON(status int, data interface{}) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
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
