package ssg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"blog/internal/db"
	"blog/internal/models"
)

// ExportMarkdown exports all posts as Markdown files with Hugo front matter
func ExportMarkdown(ctx context.Context, database *db.DB, outputDir string) (int, error) {
	// Get all posts
	posts, _, err := database.ListPosts(ctx, &models.ListOptions{Limit: 10000})
	if err != nil {
		return 0, fmt.Errorf("failed to list posts: %w", err)
	}

	// Create output directory
	postsDir := filepath.Join(outputDir, "content", "posts")
	if err := os.MkdirAll(postsDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create directory: %w", err)
	}

	exported := 0
	for _, post := range posts {
		// Only export published posts
		if post.Status != "published" {
			continue
		}

		// Build front matter
		frontMatter := buildFrontMatter(post)

		// Combine front matter + content
		fileContent := frontMatter + "\n\n" + post.Content

		// Write file
		fileName := fmt.Sprintf("%s.md", post.Slug)
		filePath := filepath.Join(postsDir, fileName)

		if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
			return exported, fmt.Errorf("failed to write file %s: %w", filePath, err)
		}

		exported++
	}

	return exported, nil
}

// buildFrontMatter creates Hugo-compatible YAML front matter
func buildFrontMatter(post *models.Post) string {
	front := "---\n"
	front += fmt.Sprintf("title: \"%s\"\n", escapeYAML(post.Title))
	front += fmt.Sprintf("date: %s\n", post.CreatedAt.Format("2006-01-02"))
	front += fmt.Sprintf("slug: %s\n", post.Slug)
	front += fmt.Sprintf("draft: %v\n", post.Status != "published")
	front += fmt.Sprintf("type: %s\n", post.TypeID)

	if post.Excerpt != "" {
		front += fmt.Sprintf("description: \"%s\"\n", escapeYAML(post.Excerpt))
	}

	// Add tags if present
	if post.Tags != nil && len(post.Tags) > 0 {
		tagsStr := "["
		for i, tag := range post.Tags {
			if i > 0 {
				tagsStr += ", "
			}
			tagsStr += fmt.Sprintf("\"%s\"", escapeYAML(tag))
		}
		tagsStr += "]"
		front += fmt.Sprintf("tags: %s\n", tagsStr)
	}

	// Add metadata if present
	if post.Metadata != nil && len(post.Metadata) > 0 {
		var metadata map[string]interface{}
		if err := json.Unmarshal(post.Metadata, &metadata); err == nil && len(metadata) > 0 {
			front += "metadata:\n"
			for key, value := range metadata {
				if str, ok := value.(string); ok {
					front += fmt.Sprintf("  %s: \"%s\"\n", key, escapeYAML(str))
				} else {
					front += fmt.Sprintf("  %s: %v\n", key, value)
				}
			}
		}
	}

	front += "---"
	return front
}

// escapeYAML escapes special characters in YAML strings
func escapeYAML(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

// GenerateHugoConfig generates Hugo configuration file
func GenerateHugoConfig(outputDir string) error {
	hugoConfig := `baseURL = "https://example.com/"
languageCode = "en-us"
title = "Blog"
theme = "blog-theme"

[params]
  author = "Author"
  description = "A minimalist blog"

[[menu.main]]
  name = "Home"
  url = "/"
  weight = 1

[[menu.main]]
  name = "About"
  url = "/about/"
  weight = 2
`

	configPath := filepath.Join(outputDir, "hugo.toml")
	return os.WriteFile(configPath, []byte(hugoConfig), 0644)
}

// CreateDeployConfig generates GitHub Actions workflow for deployment
func CreateDeployConfig(outputDir string) error {
	workflowDir := filepath.Join(outputDir, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return err
	}

	workflow := `name: Deploy Blog

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
      
      - name: Setup Hugo
        uses: peaceiris/actions-hugo@v2
        with:
          hugo-version: 'latest'
      
      - name: Build site
        run: hugo --minify
      
      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
`

	workflowPath := filepath.Join(workflowDir, "deploy.yml")
	return os.WriteFile(workflowPath, []byte(workflow), 0644)
}

// ExportResult contains the export operation results
type ExportResult struct {
	ExportedAt  time.Time `json:"exported_at"`
	PostsCount  int       `json:"posts_count"`
	OutputDir   string    `json:"output_dir"`
	FilesCount  int       `json:"files_count"`
	Success     bool      `json:"success"`
	Message     string    `json:"message"`
}

// ExportAll runs complete export: markdown files + Hugo config + deploy workflow
func ExportAll(ctx context.Context, database *db.DB, outputDir string) (*ExportResult, error) {
	result := &ExportResult{
		ExportedAt: time.Now(),
		OutputDir:  outputDir,
	}

	// Export markdown files
	count, err := ExportMarkdown(ctx, database, outputDir)
	if err != nil {
		result.Message = fmt.Sprintf("Export failed: %v", err)
		return result, err
	}
	result.FilesCount = count

	// Generate Hugo config
	if err := GenerateHugoConfig(outputDir); err != nil {
		result.Message = fmt.Sprintf("Hugo config generation failed: %v", err)
		return result, err
	}

	// Create GitHub Actions workflow
	if err := CreateDeployConfig(outputDir); err != nil {
		result.Message = fmt.Sprintf("Deploy config generation failed: %v", err)
		return result, err
	}

	result.PostsCount = count
	result.Success = true
	result.Message = fmt.Sprintf("Successfully exported %d posts to %s", count, outputDir)

	return result, nil
}
