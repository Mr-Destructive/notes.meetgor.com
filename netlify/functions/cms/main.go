package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"blog/internal/db"
	h "blog/internal/handler"
	"blog/internal/models"
	"blog/internal/ssg"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

var database *db.DB

func init() {
	// Skip .env loading in Lambda
	log.Println("CMS function initializing...")
}

// handler is the AWS Lambda handler for Netlify Functions
func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request: %s %s", req.HTTPMethod, req.Path)

	// Note: Database initialization disabled for Lambda
	// Lambda is ephemeral - use Turso remote DB instead
	// if database == nil {
	// 	var err error
	// 	database, err = db.New(ctx)
	// 	...
	// }

	// Parse path
	path := strings.TrimPrefix(req.Path, "/.netlify/functions/cms")
	path = strings.TrimPrefix(path, "/api")
	path = strings.Trim(path, "/")

	// Health check
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
	case "auth":
		return handleAuth(req, database, id)
	case "posts":
		return handlePosts(req, database, id, action)
	case "series":
		return handleSeries(req, database, id, action)
	case "types":
		return handleTypes(req, database)
	case "tags":
		return handleTags(req, database)
	case "exports":
		return handleExports(req, database)
	default:
		return respondError(404, "Resource not found"), nil
	}
}

// Stub handlers that convert Lambda events to http.Request format
func handleAuth(req events.APIGatewayProxyRequest, db *db.DB, action string) (events.APIGatewayProxyResponse, error) {
	// TODO: Implement proper handlers
	return respondError(501, "Not implemented"), nil
}

func handlePosts(req events.APIGatewayProxyRequest, db *db.DB, id, action string) (events.APIGatewayProxyResponse, error) {
	return respondError(501, "Not implemented"), nil
}

func handleSeries(req events.APIGatewayProxyRequest, db *db.DB, id, action string) (events.APIGatewayProxyResponse, error) {
	return respondError(501, "Not implemented"), nil
}

func handleTypes(req events.APIGatewayProxyRequest, db *db.DB) (events.APIGatewayProxyResponse, error) {
	return respondError(501, "Not implemented"), nil
}

func handleTags(req events.APIGatewayProxyRequest, db *db.DB) (events.APIGatewayProxyResponse, error) {
	return respondError(501, "Not implemented"), nil
}

func handleExports(req events.APIGatewayProxyRequest, db *db.DB) (events.APIGatewayProxyResponse, error) {
	return respondError(501, "Not implemented"), nil
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
