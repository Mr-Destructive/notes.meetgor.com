package main

import (
	"compress/gzip"
	"encoding/json"
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

// JSON API response from rss2json
type RSS2JSONResponse struct {
	Items []RSS2JSONItem `json:"items"`
}

type RSS2JSONItem struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	PubDate     string `json:"pubDate"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

func main() {
	// Use RSS proxy service to bypass Substack blocking
	feedURL := "https://api.rss2json.com/v1/api.json?rss_url=https://techstructively.substack.com/feed"
	postsDir := "exports/content/posts"

	if err := os.MkdirAll(postsDir, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	var data []byte
	var lastErr error
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("GET", feedURL, nil)
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		// Set realistic headers to avoid blocking
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml, */*")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("DNT", "1")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("Referer", "https://techstructively.substack.com/")

		client := &http.Client{
			Timeout: 15 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil // Allow redirects
			},
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			log.Printf("Attempt %d/%d: Failed to fetch RSS: %v", attempt, maxRetries, err)
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second) // Exponential backoff
				continue
			}
			log.Fatalf("Failed to fetch RSS after %d attempts: %v", maxRetries, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
			log.Printf("Attempt %d/%d: Failed to fetch RSS: HTTP %d", attempt, maxRetries, resp.StatusCode)
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			log.Fatalf("Failed to fetch RSS after %d attempts: HTTP %d", maxRetries, resp.StatusCode)
		}

		var readErr error
		var body io.Reader = resp.Body

		// Handle gzip encoding
		if resp.Header.Get("Content-Encoding") == "gzip" {
			gr, err := gzip.NewReader(resp.Body)
			if err != nil {
				lastErr = err
				log.Printf("Attempt %d/%d: Failed to create gzip reader: %v", attempt, maxRetries, err)
				resp.Body.Close()
				if attempt < maxRetries {
					time.Sleep(time.Duration(attempt*2) * time.Second)
					continue
				}
				log.Fatalf("Failed to create gzip reader after %d attempts: %v", maxRetries, err)
			}
			body = gr
			defer gr.Close()
		}

		data, readErr = io.ReadAll(body)
		if readErr != nil {
			lastErr = readErr
			log.Printf("Attempt %d/%d: Failed to read RSS body: %v", attempt, maxRetries, readErr)
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			log.Fatalf("Failed to read RSS body after %d attempts: %v", maxRetries, readErr)
		}

		// Successfully fetched, break out of retry loop
		lastErr = nil
		break
	}

	if lastErr != nil && len(data) == 0 {
		log.Fatalf("Failed to fetch RSS: %v", lastErr)
	}

	// Debug: log first 500 chars of response
	dataStr := string(data)
	if len(dataStr) > 500 {
		log.Printf("DEBUG: Response preview: %s...", dataStr[:500])
	} else {
		log.Printf("DEBUG: Response preview: %s", dataStr)
	}

	var jsonResp RSS2JSONResponse
	if err := json.Unmarshal(data, &jsonResp); err != nil {
		log.Printf("Failed to parse JSON: %v", err)
		log.Printf("Response length: %d bytes", len(data))
		log.Fatalf("Response preview: %s", dataStr)
	}

	count := 0
	for _, item := range jsonResp.Items {
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

		// Parse date - rss2json returns: 2026-01-30 18:45:23
		dateFormats := []string{
			"2006-01-02 15:04:05",           // rss2json format
			time.RFC1123,                     // Original Substack format
			"2006-01-02T15:04:05Z07:00",    // ISO format
		}
		var t time.Time
		var err error
		for _, format := range dateFormats {
			t, err = time.Parse(format, item.PubDate)
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Printf("Failed to parse date %q: %v (using current time)", item.PubDate, err)
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
