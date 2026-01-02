package util

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// GenerateID creates a unique ID
func GenerateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(s string) string {
	// Convert to lowercase
	slug := strings.ToLower(s)

	// Remove/replace special characters
	slug = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`[\s-]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}

// HashPassword hashes a password with bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword verifies a password against a hash
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// ReadingTime estimates reading time in minutes (roughly 200 words/minute)
func ReadingTime(content string) int {
	words := len(strings.Fields(content))
	minutes := words / 200
	if minutes < 1 {
		minutes = 1
	}
	return minutes
}

// Timestamp returns current time in UTC
func Timestamp() time.Time {
	return time.Now().UTC()
}

// TruncateString truncates a string to max length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
