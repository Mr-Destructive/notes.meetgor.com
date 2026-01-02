package db

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"blog/internal/models"
)

var testDB *DB

func TestMain(m *testing.M) {
	// Setup test database
	var err error
	testDB, err = NewLocal(context.Background(), ":memory:")
	if err != nil {
		panic(err)
	}

	// Initialize schema
	if err := testDB.InitSchema(context.Background()); err != nil {
		panic(err)
	}

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

func TestCreatePost(t *testing.T) {
	ctx := context.Background()

	post := &models.PostCreate{
		TypeID:   "article",
		Title:    "Test Article",
		Slug:     "test-article",
		Content:  "# Test\n\nContent here",
		Excerpt:  "Test excerpt",
		Tags:     []string{"test", "go"},
		Status:   "draft",
		Metadata: json.RawMessage(`{"key":"value"}`),
	}

	result, err := testDB.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	if result.ID == "" {
		t.Error("Post ID should not be empty")
	}
	if result.Title != post.Title {
		t.Errorf("Title mismatch: got %q, want %q", result.Title, post.Title)
	}
	if result.Status != "draft" {
		t.Errorf("Status mismatch: got %q, want %q", result.Status, "draft")
	}
	if len(result.Tags) != 2 {
		t.Errorf("Tags count mismatch: got %d, want 2", len(result.Tags))
	}
}

func TestGetPost(t *testing.T) {
	ctx := context.Background()

	// Create a post first
	post := &models.PostCreate{
		TypeID:  "article",
		Title:   "Get Test",
		Slug:    "get-test",
		Content: "Content",
		Excerpt: "Excerpt",
		Status:  "draft",
	}

	created, err := testDB.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	// Test GetPost by ID
	retrieved, err := testDB.GetPost(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetPost by ID failed: %v", err)
	}

	if retrieved.Title != created.Title {
		t.Errorf("Title mismatch: got %q, want %q", retrieved.Title, created.Title)
	}

	// Test GetPost by slug
	bySlug, err := testDB.GetPost(ctx, "get-test")
	if err != nil {
		t.Fatalf("GetPost by slug failed: %v", err)
	}

	if bySlug.ID != created.ID {
		t.Errorf("ID mismatch: got %q, want %q", bySlug.ID, created.ID)
	}
}

func TestListPosts(t *testing.T) {
	ctx := context.Background()

	// Create multiple posts
	for i := 0; i < 5; i++ {
		post := &models.PostCreate{
			TypeID:  "article",
			Title:   "List Test " + string(rune(i)),
			Slug:    "list-test-" + string(rune(i)),
			Content: "Content",
			Status:  "draft",
		}
		_, err := testDB.CreatePost(ctx, post)
		if err != nil {
			t.Fatalf("CreatePost failed: %v", err)
		}
	}

	// List all posts
	opts := &models.ListOptions{
		Limit:  50,
		Offset: 0,
	}

	posts, total, err := testDB.ListPosts(ctx, opts)
	if err != nil {
		t.Fatalf("ListPosts failed: %v", err)
	}

	if total < 5 {
		t.Errorf("Total posts count: got %d, want at least 5", total)
	}

	if len(posts) < 5 {
		t.Errorf("Posts count: got %d, want at least 5", len(posts))
	}
}

func TestUpdatePost(t *testing.T) {
	ctx := context.Background()

	// Create post
	post := &models.PostCreate{
		TypeID:  "article",
		Title:   "Update Test",
		Slug:    "update-test",
		Content: "Original content",
		Status:  "draft",
	}

	created, err := testDB.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	// Update post
	newTitle := "Updated Title"
	newStatus := "published"
	update := &models.PostUpdate{
		Title:  &newTitle,
		Status: &newStatus,
	}

	updated, err := testDB.UpdatePost(ctx, created.ID, update)
	if err != nil {
		t.Fatalf("UpdatePost failed: %v", err)
	}

	if updated.Title != newTitle {
		t.Errorf("Title not updated: got %q, want %q", updated.Title, newTitle)
	}

	if updated.Status != newStatus {
		t.Errorf("Status not updated: got %q, want %q", updated.Status, newStatus)
	}
}

func TestDeletePost(t *testing.T) {
	ctx := context.Background()

	// Create post
	post := &models.PostCreate{
		TypeID:  "article",
		Title:   "Delete Test",
		Slug:    "delete-test",
		Content: "Content",
		Status:  "draft",
	}

	created, err := testDB.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	// Delete post
	err = testDB.DeletePost(ctx, created.ID)
	if err != nil {
		t.Fatalf("DeletePost failed: %v", err)
	}

	// Verify deletion
	_, err = testDB.GetPost(ctx, created.ID)
	if err == nil {
		t.Error("Post should not exist after deletion")
	}
}

func TestGetPostTypes(t *testing.T) {
	ctx := context.Background()

	types, err := testDB.GetPostTypes(ctx)
	if err != nil {
		t.Fatalf("GetPostTypes failed: %v", err)
	}

	if len(types) != 12 {
		t.Errorf("Expected 12 post types, got %d", len(types))
	}

	// Check for known types
	typeNames := make(map[string]bool)
	for _, pt := range types {
		typeNames[pt.ID] = true
	}

	expectedTypes := []string{"article", "review", "thought", "link", "til", "quote", "list", "note", "snippet", "essay", "tutorial", "interview"}
	for _, expected := range expectedTypes {
		if !typeNames[expected] {
			t.Errorf("Missing post type: %s", expected)
		}
	}
}

func TestRevisions(t *testing.T) {
	ctx := context.Background()

	// Create post
	post := &models.PostCreate{
		TypeID:  "article",
		Title:   "Revision Test",
		Slug:    "revision-test",
		Content: "Original content",
		Status:  "draft",
	}

	created, err := testDB.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	// Create revision
	originalContent := "Original content"
	err = testDB.CreateRevision(ctx, created.ID, originalContent)
	if err != nil {
		t.Fatalf("CreateRevision failed: %v", err)
	}

	// Update post (should create another revision)
	newContent := "Updated content"
	update := &models.PostUpdate{
		Content: &newContent,
	}
	testDB.UpdatePost(ctx, created.ID, update)

	// Get revisions
	revisions, err := testDB.GetRevisions(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetRevisions failed: %v", err)
	}

	if len(revisions) < 1 {
		t.Errorf("Expected at least 1 revision, got %d", len(revisions))
	}
}

func TestFilterByStatus(t *testing.T) {
	ctx := context.Background()

	// Create posts with different statuses
	statuses := []string{"draft", "published", "archived"}
	for _, status := range statuses {
		post := &models.PostCreate{
			TypeID:  "article",
			Title:   "Status Test " + status,
			Slug:    "status-test-" + status,
			Content: "Content",
			Status:  status,
		}
		_, err := testDB.CreatePost(ctx, post)
		if err != nil {
			t.Fatalf("CreatePost failed: %v", err)
		}
	}

	// Filter by status
	opts := &models.ListOptions{
		Limit:  50,
		Status: "published",
	}

	posts, _, err := testDB.ListPosts(ctx, opts)
	if err != nil {
		t.Fatalf("ListPosts failed: %v", err)
	}

	for _, p := range posts {
		if p.Status != "published" {
			t.Errorf("Filter not working: got status %q", p.Status)
		}
	}
}

func TestFilterByType(t *testing.T) {
	ctx := context.Background()

	// Create posts with different types
	types := []string{"article", "review"}
	for _, typeID := range types {
		post := &models.PostCreate{
			TypeID:  typeID,
			Title:   "Type Test " + typeID,
			Slug:    "type-test-" + typeID,
			Content: "Content",
			Status:  "draft",
		}
		_, err := testDB.CreatePost(ctx, post)
		if err != nil {
			t.Fatalf("CreatePost failed: %v", err)
		}
	}

	// Filter by type
	opts := &models.ListOptions{
		Limit: 50,
		Type:  "review",
	}

	posts, _, err := testDB.ListPosts(ctx, opts)
	if err != nil {
		t.Fatalf("ListPosts failed: %v", err)
	}

	for _, p := range posts {
		if p.TypeID != "review" {
			t.Errorf("Filter not working: got type %q", p.TypeID)
		}
	}
}

func TestPostMetadata(t *testing.T) {
	ctx := context.Background()

	metadata := json.RawMessage(`{"category":"tech","reading_time":5}`)
	post := &models.PostCreate{
		TypeID:   "article",
		Title:    "Metadata Test",
		Slug:     "metadata-test",
		Content:  "Content",
		Status:   "draft",
		Metadata: metadata,
	}

	created, err := testDB.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	if string(created.Metadata) != string(metadata) {
		t.Errorf("Metadata not stored correctly")
	}
}

func TestPostTags(t *testing.T) {
	ctx := context.Background()

	tags := []string{"golang", "testing", "database"}
	post := &models.PostCreate{
		TypeID:  "article",
		Title:   "Tags Test",
		Slug:    "tags-test",
		Content: "Content",
		Status:  "draft",
		Tags:    tags,
	}

	created, err := testDB.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	if len(created.Tags) != len(tags) {
		t.Errorf("Tags count mismatch: got %d, want %d", len(created.Tags), len(tags))
	}

	for i, tag := range created.Tags {
		if tag != tags[i] {
			t.Errorf("Tag mismatch at %d: got %q, want %q", i, tag, tags[i])
		}
	}
}

func TestPublishedAt(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	post := &models.PostCreate{
		TypeID:      "article",
		Title:       "Published Test",
		Slug:        "published-test",
		Content:     "Content",
		Status:      "published",
		PublishedAt: &now,
	}

	created, err := testDB.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	if created.PublishedAt == nil {
		t.Error("PublishedAt should not be nil")
	}

	if created.PublishedAt.Unix() != now.Unix() {
		t.Errorf("PublishedAt mismatch")
	}
}
