package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

type Post struct {
	ID        string
	Title     string
	Slug      string
	Content   string
	Excerpt   sql.NullString
	TypeID    string
	Status    string
	CreatedAt string
	Tags      sql.NullString
}

func main() {
	// Get Turso connection URL and token from env
	connURL := os.Getenv("TURSO_CONNECTION_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	if connURL == "" || authToken == "" {
		log.Fatal("TURSO_CONNECTION_URL and TURSO_AUTH_TOKEN required")
	}

	// Connect to Turso
	dbURL := connURL + "?authToken=" + authToken
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Ping failed: %v", err)
	}
	cancel()

	// Query posts
	rows, err := db.Query("SELECT id, title, slug, content, excerpt, type_id, status, created_at, tags FROM posts WHERE status = 'published'")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	// Create exports directory
	postsDir := "exports/content/posts"
	if err := os.MkdirAll(postsDir, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	count := 0
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Content, &post.Excerpt, &post.TypeID, &post.Status, &post.CreatedAt, &post.Tags); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		// Default type to "posts"
		typeID := post.TypeID
		if typeID == "" {
			typeID = "posts"
		}

		// Handle nullable fields
		excerpt := ""
		if post.Excerpt.Valid {
			excerpt = post.Excerpt.String
		}
		
		tags := "[]"
		if post.Tags.Valid {
			tags = post.Tags.String
		}

		// Build front matter
		frontMatter := fmt.Sprintf(`---
		title: "%s"
		date: %s
		slug: %s
		draft: false
		type: %s
		description: "%s"
		tags: %s
		---
		
		`, escapeYAML(post.Title), post.CreatedAt[:10], post.Slug, typeID, escapeYAML(excerpt), tags)

		// Write file
		filePath := filepath.Join(postsDir, post.Slug+".md")
		content := frontMatter + post.Content
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			log.Printf("Write error for %s: %v", filePath, err)
			continue
		}

		fmt.Printf("✓ Exported: %s\n", post.Slug)
		count++
	}

	fmt.Printf("\n✓ Total: %d posts exported\n", count)
}

func escapeYAML(s string) string {
	// Simple YAML escaping
	if s == "" {
		return ""
	}
	// Replace newlines and quotes
	result := ""
	for _, c := range s {
		if c == '"' {
			result += "\\\""
		} else if c == '\n' {
			result += " "
		} else {
			result += string(c)
		}
	}
	return result
}
