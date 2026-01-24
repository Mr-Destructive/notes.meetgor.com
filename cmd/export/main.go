package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

type Post struct {
	ID          string
	Title       string
	Slug        string
	Content     string
	Excerpt     sql.NullString
	TypeID      string
	Status      string
	CreatedAt   string
	PublishedAt sql.NullString
	Tags        sql.NullString
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

	// Query posts - include published_at and order by latest
	rows, err := db.Query("SELECT id, title, slug, content, excerpt, type_id, status, created_at, published_at, tags FROM posts WHERE status = 'published' ORDER BY COALESCE(published_at, created_at) DESC")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Content, &post.Excerpt, &post.TypeID, &post.Status, &post.CreatedAt, &post.PublishedAt, &post.Tags); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		// Clean slug - remove trailing slashes which cause file write errors
		post.Slug = strings.Trim(post.Slug, "/")
		if post.Slug == "" {
			log.Printf("Empty slug for post ID: %s", post.ID)
			continue
		}

		// Determine target directory based on type_id
		// Map type_id to folder name
		typeDir := "posts"
		if post.TypeID != "" {
			// Some common mappings
			switch strings.ToLower(post.TypeID) {
			case "newsletter":
				typeDir = "newsletter"
			case "link", "links":
				typeDir = "links"
			case "thought", "thoughts":
				typeDir = "thoughts"
			case "quote", "quotes":
				typeDir = "quotes"
			default:
				typeDir = post.TypeID
			}
		}

		postsDir := filepath.Join("exports/content", typeDir)
		if err := os.MkdirAll(postsDir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", postsDir, err)
		}

		// Handle nullable fields
		excerpt := ""
		if post.Excerpt.Valid {
			excerpt = post.Excerpt.String
		}
		
		tags := "[]"
		if post.Tags.Valid && post.Tags.String != "" {
			tags = post.Tags.String
		}

		// Parse and format date properly for Hugo
		dateStr := post.CreatedAt
		if post.PublishedAt.Valid && post.PublishedAt.String != "" {
			dateStr = post.PublishedAt.String
		}

		parsedDate, dateErr := time.Parse("2006-01-02 15:04:05", dateStr)
		if dateErr != nil {
			parsedDate, dateErr = time.Parse("2006-01-02T15:04:05Z", dateStr)
			if dateErr != nil {
				if len(dateStr) >= 10 {
					dateStr = dateStr[:10]
				}
			} else {
				dateStr = parsedDate.Format("2006-01-02T15:04:05Z07:00")
			}
		} else {
			dateStr = parsedDate.Format("2006-01-02T15:04:05Z07:00")
		}

		// Build front matter
		frontMatter := fmt.Sprintf("---\ntitle: %q\ndate: %s\nslug: %s\ndraft: false\ntype: %s\ndescription: %q\ntags: %s\n---\n\n", 
			post.Title, dateStr, post.Slug, typeDir, excerpt, tags)

		// Write file
		filePath := filepath.Join(postsDir, post.Slug+".md")
		content := frontMatter + strings.TrimSpace(post.Content)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			log.Printf("Write error for %s: %v", filePath, err)
			continue
		}

		fmt.Printf("✓ Exported [%s]: %s\n", typeDir, post.Slug)
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
