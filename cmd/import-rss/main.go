package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Content     string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	Description string `xml:"description"`
}

func main() {
	feedURL := "https://techstructively.substack.com/feed"
	postsDir := "exports/content/posts"

	if err := os.MkdirAll(postsDir, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	resp, err := http.Get(feedURL)
	if err != nil {
		log.Fatalf("Failed to fetch RSS: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read RSS body: %v", err)
	}

	var rss RSS
	if err := xml.Unmarshal(data, &rss); err != nil {
		log.Fatalf("Failed to parse XML: %v", err)
	}

	count := 0
	for _, item := range rss.Channel.Items {
		slug := extractSlug(item.Link)
		if slug == "" {
			continue
		}

		filename := filepath.Join(postsDir, slug+".md")
		if _, err := os.Stat(filename); err == nil {
			// File exists, skip
			// fmt.Printf("Skipping existing: %s\n", slug)
			continue
		}

		// Parse date
		// RFC1123: Fri, 16 Jan 2026 18:45:34 GMT
		t, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			log.Printf("Failed to parse date %q: %v", item.PubDate, err)
			t = time.Now()
		}
		dateStr := t.Format("2006-01-02T15:04:05Z07:00")

		// Prepare frontmatter
		// We use tags=["substack"] to identify source
		frontMatter := fmt.Sprintf("---\ntitle: %q\ndate: %s\nslug: %s\ndraft: false\ntype: post\ndescription: %q\ntags: [\"substack\"]\n---\n\n",
			item.Title, dateStr, slug, escapeQuotes(item.Description))

		content := frontMatter + item.Content

		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			log.Printf("Failed to write %s: %v", filename, err)
			continue
		}

		fmt.Printf("✓ Imported: %s\n", slug)
		count++
	}

	fmt.Printf("Imported %d new posts from Substack\n", count)
}

func extractSlug(url string) string {
	// https://techstructively.substack.com/p/techstructive-weekly-77
	parts := strings.Split(url, "/p/")
	if len(parts) < 2 {
		return ""
	}
	// Remove potential query params
	slug := strings.Split(parts[1], "?")[0]
	slug = strings.TrimSuffix(slug, "/")
	return slug
}

func escapeQuotes(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}
