package util

import (
	"strings"
	"testing"
)

func TestGenerateID(t *testing.T) {
	id1 := GenerateID()
	id2 := GenerateID()

	if id1 == "" {
		t.Error("ID should not be empty")
	}

	if id1 == id2 {
		t.Error("IDs should be unique")
	}

	// Should be hex encoded
	if !strings.ContainsAny(id1, "0123456789abcdef") {
		t.Error("ID should be hex encoded")
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"My Amazing Post!", "my-amazing-post"},
		{"Test   Multiple   Spaces", "test-multiple-spaces"},
		{"CamelCaseTitle", "camelcasetitle"},
		{"hello-world", "hello-world"},
		{"   Leading and trailing   ", "leading-and-trailing"},
	}

	for _, test := range tests {
		result := GenerateSlug(test.input)
		if result != test.expected {
			t.Errorf("GenerateSlug(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	if hash == password {
		t.Error("Hash should not equal plaintext password")
	}

	// Should look like bcrypt hash
	if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") {
		t.Error("Hash should start with bcrypt prefix")
	}
}

func TestCheckPassword_Valid(t *testing.T) {
	password := "testPassword123"

	hash, _ := HashPassword(password)

	if !CheckPassword(hash, password) {
		t.Error("CheckPassword should return true for correct password")
	}
}

func TestCheckPassword_Invalid(t *testing.T) {
	password := "testPassword123"
	wrongPassword := "wrongPassword456"

	hash, _ := HashPassword(password)

	if CheckPassword(hash, wrongPassword) {
		t.Error("CheckPassword should return false for incorrect password")
	}
}

func TestReadingTime(t *testing.T) {
	tests := []struct {
		content string
		minTime int
		maxTime int
	}{
		{
			content: strings.Repeat("word ", 100),
			minTime: 0,
			maxTime: 1,
		},
		{
			content: strings.Repeat("word ", 400),
			minTime: 1,
			maxTime: 3,
		},
		{
			content: strings.Repeat("word ", 2000),
			minTime: 8,
			maxTime: 12,
		},
	}

	for _, test := range tests {
		result := ReadingTime(test.content)
		if result < test.minTime || result > test.maxTime {
			t.Errorf("ReadingTime(%d words) = %d, want between %d-%d", 
				len(strings.Fields(test.content)), result, test.minTime, test.maxTime)
		}
	}
}

func TestReadingTime_Minimum(t *testing.T) {
	// Very short content should still return 1 minute
	content := "This is short"
	result := ReadingTime(content)

	if result < 1 {
		t.Errorf("ReadingTime should return at least 1, got %d", result)
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input   string
		maxLen  int
		expected string
	}{
		{"Hello World", 5, "Hello..."},
		{"Hello World", 11, "Hello World"},
		{"Hello World", 20, "Hello World"},
		{"Test", 4, "Test"},
		{"Test", 2, "Te..."},
		{"", 10, ""},
	}

	for _, test := range tests {
		result := TruncateString(test.input, test.maxLen)
		if result != test.expected {
			t.Errorf("TruncateString(%q, %d) = %q, want %q", 
				test.input, test.maxLen, result, test.expected)
		}
	}
}

func TestTimestamp(t *testing.T) {
	ts := Timestamp()

	if ts.IsZero() {
		t.Error("Timestamp should not be zero")
	}

	// Should be in UTC
	if ts.Location().String() != "UTC" {
		t.Errorf("Timestamp should be UTC, got %s", ts.Location().String())
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "testPassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashPassword(password)
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "testPassword123"
	hash, _ := HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPassword(hash, password)
	}
}

func BenchmarkGenerateSlug(b *testing.B) {
	text := "This is a test title for slug generation"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateSlug(text)
	}
}
