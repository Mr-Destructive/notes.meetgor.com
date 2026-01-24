package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

// YouTubeMetadata holds information about a YouTube video
type YouTubeMetadata struct {
	Title        string `json:"title"`
	AuthorName   string `json:"author_name"`
	ThumbnailURL string `json:"thumbnail_url"`
	VideoID      string `json:"video_id"`
}

// ExtractYouTubeID parses a YouTube URL and returns the video ID
func ExtractYouTubeID(url string) string {
	patterns := []string{
		`v=([a-zA-Z0-9_-]{11})`,
		`youtu\.be/([a-zA-Z0-9_-]{11})`,
		`embed/([a-zA-Z0-9_-]{11})`,
		`live/([a-zA-Z0-9_-]{11})`,
		`v/([a-zA-Z0-9_-]{11})`,
	}
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

// FetchYouTubeMetadata retrieves video info via oEmbed
func FetchYouTubeMetadata(videoID string) (*YouTubeMetadata, error) {
	if videoID == "" {
		return nil, fmt.Errorf("empty video ID")
	}

	url := fmt.Sprintf("https://www.youtube.com/oembed?url=https://www.youtube.com/watch?v=%s&format=json", videoID)
	
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch metadata: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var meta YouTubeMetadata
	if err := json.Unmarshal(body, &meta); err != nil {
		return nil, err
	}

	meta.VideoID = videoID
	return &meta, nil
}
