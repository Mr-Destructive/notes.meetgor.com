package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"blog/internal/db"
	"blog/internal/models"
)

// HandleExportsGet returns published posts as JSON for export
func HandleExportsGet(w http.ResponseWriter, r *http.Request, database *db.DB) {
	ctx := context.Background()

	// Get all published posts
	posts, _, err := database.ListPosts(ctx, &models.ListOptions{
		Status: "published",
		Limit:  1000,
	})
	if err != nil {
		renderHTML(w, "500", fmt.Sprintf("Error fetching posts: %v", err))
		return
	}

	// Return as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

// HandleExportsMarkdown generates markdown files and Hugo configuration
func HandleExportsMarkdown(w http.ResponseWriter, r *http.Request, database *db.DB) {
	if r.Method != "POST" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	ctx := context.Background()

	// Get all published posts
	posts, _, err := database.ListPosts(ctx, &models.ListOptions{
		Status: "published",
		Limit:  1000,
	})
	if err != nil {
		respondExportError(w, 500, fmt.Sprintf("Error fetching posts: %v", err))
		return
	}

	// Create exports directory structure
	exportDir := "./exports"
	contentDir := filepath.Join(exportDir, "content", "posts")
	workflowDir := filepath.Join(exportDir, ".github", "workflows")

	// Create directories
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		respondExportError(w, 500, fmt.Sprintf("Failed to create directories: %v", err))
		return
	}
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		respondExportError(w, 500, fmt.Sprintf("Failed to create workflow directory: %v", err))
		return
	}

	// Write each post as markdown
	for _, post := range posts {
		markdown := postToMarkdown(post)
		filename := filepath.Join(contentDir, post.Slug+".md")
		if err := os.WriteFile(filename, []byte(markdown), 0644); err != nil {
			respondExportError(w, 500, fmt.Sprintf("Failed to write post file: %v", err))
			return
		}
	}

	// Generate hugo.toml
	hugoConfig := generateHugoConfig()
	if err := os.WriteFile(filepath.Join(exportDir, "hugo.toml"), []byte(hugoConfig), 0644); err != nil {
		respondExportError(w, 500, fmt.Sprintf("Failed to write hugo config: %v", err))
		return
	}

	// Generate GitHub Actions workflow
	workflow := generateGitHubWorkflow()
	if err := os.WriteFile(filepath.Join(workflowDir, "deploy.yml"), []byte(workflow), 0644); err != nil {
		respondExportError(w, 500, fmt.Sprintf("Failed to write workflow: %v", err))
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"exported_at": time.Now().Format(time.RFC3339),
		"posts_count": len(posts),
		"files_count": len(posts) + 2, // posts + hugo.toml + workflow
		"output_dir":  exportDir,
		"success":     true,
		"message":     fmt.Sprintf("Successfully exported %d posts to %s", len(posts), exportDir),
	})
}

// postToMarkdown converts a post to markdown with front matter
func postToMarkdown(post *models.Post) string {
	frontMatter := generateFrontMatter(post)
	return fmt.Sprintf("---\n%s---\n\n%s", frontMatter, post.Content)
}

// generateFrontMatter creates YAML front matter for a post
func generateFrontMatter(post *models.Post) string {
	sb := strings.Builder{}

	// Title
	sb.WriteString(fmt.Sprintf("title: \"%s\"\n", escapeYAML(post.Title)))

	// Date
	sb.WriteString(fmt.Sprintf("date: %s\n", post.CreatedAt.Format("2006-01-02")))

	// Slug
	sb.WriteString(fmt.Sprintf("slug: %s\n", post.Slug))

	// Draft status (false for published)
	sb.WriteString("draft: false\n")

	// Type
	sb.WriteString(fmt.Sprintf("type: %s\n", post.TypeID))

	// Description/Excerpt
	if post.Excerpt != "" {
		sb.WriteString(fmt.Sprintf("description: \"%s\"\n", escapeYAML(post.Excerpt)))
	}

	// Tags
	if len(post.Tags) > 0 {
		sb.WriteString("tags:\n")
		for _, tag := range post.Tags {
			sb.WriteString(fmt.Sprintf("  - %s\n", escapeYAML(tag)))
		}
	}

	// Metadata if present
	if len(post.Metadata) > 0 {
		sb.WriteString("metadata:\n")
		for key, value := range post.Metadata {
			// Convert value to string
			valStr := fmt.Sprintf("%v", value)
			sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", key, escapeYAML(valStr)))
		}
	}

	return sb.String()
}

// escapeYAML escapes special characters for YAML
func escapeYAML(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}

// generateHugoConfig creates a hugo.toml configuration file
func generateHugoConfig() string {
	return `baseURL = "https://example.com/"
languageCode = "en-us"
title = "My Blog"
theme = "blog-theme"

[params]
  author = "Your Name"
  description = "A blog about technology and software"

# Content structure
[content]
  [[content.dirs]]
    path = "content/posts"
    singular = false

# Markdown processing
[markup]
  [markup.goldmark]
    [markup.goldmark.renderer]
      hardWraps = true
      xhtml = false

# Output formats
[outputs]
  home = ["HTML", "JSON"]
`
}

// generateGitHubWorkflow creates a GitHub Actions workflow file for deployment
func generateGitHubWorkflow() string {
	return `name: Deploy Hugo Site

on:
  push:
    branches:
      - main
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
        with:
          submodules: recursive

      - name: Setup Hugo
        uses: peaceiris/actions-hugo@v2
        with:
          hugo-version: latest
          extended: true

      - name: Build
        run: hugo --minify

      - name: Deploy
        uses: peaceiris/actions-gh-pages@v3
        if: github.ref == 'refs/heads/main'
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
`
}

// respondExportError sends an error response in JSON format
func respondExportError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
