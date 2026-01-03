package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type DraftItem struct {
	ID        string `json:"id"`
	TypeID    string `json:"type_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Metadata  string `json:"metadata"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type TempItem struct {
	ID        string      `json:"id"`
	TypeID    string      `json:"type_id"`
	Slug      string      `json:"slug"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	CreatedAt int64       `json:"created_at"`
	UpdatedAt int64       `json:"updated_at"`
	Published interface{} `json:"published"` // Can be bool or int
}

type LinkItem struct {
	ID         interface{} `json:"id"` // Can be string or number
	URL        string      `json:"url"`
	Title      string      `json:"title"`
	Commentary string      `json:"commentary"`
	ImageURL   string      `json:"image_url"`
}

func htmlToMarkdown(htmlContent string) string {
	// Unescape HTML entities
	content := html.UnescapeString(htmlContent)

	// Replace common HTML tags with markdown equivalents
	content = regexp.MustCompile(`<h1[^>]*>(.*?)</h1>`).ReplaceAllString(content, "# $1")
	content = regexp.MustCompile(`<h2[^>]*>(.*?)</h2>`).ReplaceAllString(content, "## $1")
	content = regexp.MustCompile(`<h3[^>]*>(.*?)</h3>`).ReplaceAllString(content, "### $1")
	content = regexp.MustCompile(`<h4[^>]*>(.*?)</h4>`).ReplaceAllString(content, "#### $1")
	content = regexp.MustCompile(`<strong[^>]*>(.*?)</strong>`).ReplaceAllString(content, "**$1**")
	content = regexp.MustCompile(`<b[^>]*>(.*?)</b>`).ReplaceAllString(content, "**$1**")
	content = regexp.MustCompile(`<em[^>]*>(.*?)</em>`).ReplaceAllString(content, "*$1*")
	content = regexp.MustCompile(`<i[^>]*>(.*?)</i>`).ReplaceAllString(content, "*$1*")
	content = regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*>(.*?)</a>`).ReplaceAllString(content, "[$2]($1)")
	content = regexp.MustCompile(`<code[^>]*>(.*?)</code>`).ReplaceAllString(content, "`$1`")
	content = regexp.MustCompile(`<pre[^>]*>(.*?)</pre>`).ReplaceAllString(content, "```\n$1\n```")
	content = regexp.MustCompile(`<p[^>]*>(.*?)</p>`).ReplaceAllString(content, "$1\n")
	content = regexp.MustCompile(`<li[^>]*>(.*?)</li>`).ReplaceAllString(content, "- $1")
	content = regexp.MustCompile(`<ul[^>]*>(.*?)</ul>`).ReplaceAllString(content, "$1")
	content = regexp.MustCompile(`<ol[^>]*>(.*?)</ol>`).ReplaceAllString(content, "$1")
	content = regexp.MustCompile(`<br[^>]*>`).ReplaceAllString(content, "\n")
	content = regexp.MustCompile(`<hr[^>]*>`).ReplaceAllString(content, "---")

	// Clean up extra whitespace
	content = regexp.MustCompile(`\n\n+`).ReplaceAllString(content, "\n\n")
	content = strings.TrimSpace(content)

	return content
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func msToDate(ms int64) string {
	return time.UnixMilli(ms).Format("2006-01-02")
}

func convertDrafts(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	var items []DraftItem
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	count := 0
	for _, item := range items {
		if item.Title == "" || item.Title == "Untitled" {
			continue
		}

		slug := generateSlug(item.Title)
		date := msToDate(item.CreatedAt)

		// Convert HTML to Markdown
		content := htmlToMarkdown(item.Content)

		// Determine type (default to "posts")
		typeID := "posts"
		if item.TypeID != "" {
			typeID = item.TypeID
		}

		// Build front matter
		frontMatter := fmt.Sprintf("---\ntitle: %q\ndate: %s\nslug: %s\ndraft: false\ntype: %s\ndescription: \"\"\ntags: []\n---\n\n",
			item.Title, date, slug, typeID)

		// Write file
		filePath := filepath.Join(outputPath, slug+".md")
		content = frontMatter + content
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			log.Printf("Error writing %s: %v", filePath, err)
			continue
		}

		fmt.Printf("✓ Converted draft: %s\n", slug)
		count++
	}

	fmt.Printf("✓ Total drafts converted: %d\n", count)
	return nil
}

func convertTemp(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	var items []TempItem
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	count := 0
	for _, item := range items {
		if item.Title == "" {
			continue
		}
		
		// Check if published (handle bool or int)
		published := false
		switch v := item.Published.(type) {
		case bool:
			published = v
		case float64:
			published = v > 0
		}
		if !published {
			continue
		}

		slug := item.Slug
		if slug == "" {
			slug = generateSlug(item.Title)
		}
		date := msToDate(item.CreatedAt)

		// Convert HTML to Markdown
		content := htmlToMarkdown(item.Content)

		// Determine type
		typeID := "posts"
		if item.TypeID != "" {
			typeID = item.TypeID
		}

		// Build front matter
		frontMatter := fmt.Sprintf("---\ntitle: %q\ndate: %s\nslug: %s\ndraft: false\ntype: %s\ndescription: \"\"\ntags: []\n---\n\n",
			item.Title, date, slug, typeID)

		// Write file
		filePath := filepath.Join(outputPath, slug+".md")
		content = frontMatter + content
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			log.Printf("Error writing %s: %v", filePath, err)
			continue
		}

		fmt.Printf("✓ Converted temp: %s\n", slug)
		count++
	}

	fmt.Printf("✓ Total temp items converted: %d\n", count)
	return nil
}

func convertLinks(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	var items []LinkItem
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	count := 0
	for _, item := range items {
		if item.Title == "" {
			continue
		}

		slug := generateSlug(item.Title)
		date := time.Now().AddDate(0, 0, -count).Format("2006-01-02")

		// Use commentary as content, or URL if no commentary
		content := item.Commentary
		if content == "" {
			content = fmt.Sprintf("Source: [%s](%s)", item.URL, item.URL)
		}

		// Build front matter for link type
		frontMatter := fmt.Sprintf("---\ntitle: %q\ndate: %s\nslug: %s\ndraft: false\ntype: link\ndescription: \"\"\ntags: [\"link\"]\n---\n\n",
			item.Title, date, slug)

		// Write file
		filePath := filepath.Join(outputPath, slug+".md")
		content = frontMatter + content
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			log.Printf("Error writing %s: %v", filePath, err)
			continue
		}

		fmt.Printf("✓ Converted link: %s\n", slug)
		count++
	}

	fmt.Printf("✓ Total links converted: %d\n", count)
	return nil
}

func main() {
	draftsFlag := flag.Bool("drafts", false, "Convert drafts.json")
	tempFlag := flag.Bool("temp", false, "Convert temp-blog.json")
	linksFlag := flag.Bool("links", false, "Convert links-blog.json")
	allFlag := flag.Bool("all", false, "Convert all files")
	outputFlag := flag.String("output", "", "Output directory (required)")

	flag.Parse()

	if *outputFlag == "" {
		log.Fatal("--output flag is required")
	}

	if err := os.MkdirAll(*outputFlag, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	if *allFlag || *draftsFlag {
		if err := convertDrafts("/home/meet/code/blog/drafts.json", *outputFlag); err != nil {
			log.Fatalf("Error converting drafts: %v", err)
		}
	}

	if *allFlag || *tempFlag {
		if err := convertTemp("/home/meet/code/blog/temp-blog.json", *outputFlag); err != nil {
			log.Fatalf("Error converting temp: %v", err)
		}
	}

	if *allFlag || *linksFlag {
		if err := convertLinks("/home/meet/code/blog/links-blog.json", *outputFlag); err != nil {
			log.Fatalf("Error converting links: %v", err)
		}
	}

	fmt.Println("\n✓ Conversion complete!")
}
