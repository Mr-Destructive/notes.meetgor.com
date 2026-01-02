package db

import (
	"context"
	"testing"

	"blog/internal/models"
)

func TestCreateSeries(t *testing.T) {
	ctx := context.Background()

	series := &models.SeriesCreate{
		Name:        "Test Series",
		Slug:        "test-series",
		Description: "A test series",
	}

	result, err := testDB.CreateSeries(ctx, series)
	if err != nil {
		t.Fatalf("CreateSeries failed: %v", err)
	}

	if result.ID == "" {
		t.Error("Series ID should not be empty")
	}
	if result.Name != series.Name {
		t.Errorf("Name mismatch: got %q, want %q", result.Name, series.Name)
	}
	if result.Slug != series.Slug {
		t.Errorf("Slug mismatch: got %q, want %q", result.Slug, series.Slug)
	}
}

func TestGetSeries(t *testing.T) {
	ctx := context.Background()

	// Create series
	series := &models.SeriesCreate{
		Name:        "Get Series Test",
		Slug:        "get-series-test",
		Description: "Test",
	}

	created, err := testDB.CreateSeries(ctx, series)
	if err != nil {
		t.Fatalf("CreateSeries failed: %v", err)
	}

	// Get by ID
	retrieved, err := testDB.GetSeries(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetSeries by ID failed: %v", err)
	}

	if retrieved.Name != created.Name {
		t.Errorf("Name mismatch: got %q, want %q", retrieved.Name, created.Name)
	}

	// Get by slug
	bySlug, err := testDB.GetSeries(ctx, "get-series-test")
	if err != nil {
		t.Fatalf("GetSeries by slug failed: %v", err)
	}

	if bySlug.ID != created.ID {
		t.Errorf("ID mismatch: got %q, want %q", bySlug.ID, created.ID)
	}
}

func TestListSeries(t *testing.T) {
	ctx := context.Background()

	// Create multiple series
	for i := 0; i < 3; i++ {
		series := &models.SeriesCreate{
			Name:        "List Test Series " + string(rune(i)),
			Slug:        "list-series-" + string(rune(i)),
			Description: "Test",
		}
		_, err := testDB.CreateSeries(ctx, series)
		if err != nil {
			t.Fatalf("CreateSeries failed: %v", err)
		}
	}

	// List series
	series, err := testDB.ListSeries(ctx, 50, 0)
	if err != nil {
		t.Fatalf("ListSeries failed: %v", err)
	}

	if len(series) < 3 {
		t.Errorf("Series count: got %d, want at least 3", len(series))
	}
}

func TestDeleteSeries(t *testing.T) {
	ctx := context.Background()

	// Create series
	series := &models.SeriesCreate{
		Name:        "Delete Series Test",
		Slug:        "delete-series-test",
		Description: "Test",
	}

	created, err := testDB.CreateSeries(ctx, series)
	if err != nil {
		t.Fatalf("CreateSeries failed: %v", err)
	}

	// Delete series
	err = testDB.DeleteSeries(ctx, created.ID)
	if err != nil {
		t.Fatalf("DeleteSeries failed: %v", err)
	}

	// Verify deletion
	_, err = testDB.GetSeries(ctx, created.ID)
	if err == nil {
		t.Error("Series should not exist after deletion")
	}
}

func TestAddPostToSeries(t *testing.T) {
	ctx := context.Background()

	// Create series and post
	series := &models.SeriesCreate{
		Name:        "Post Series Test",
		Slug:        "post-series-test",
		Description: "Test",
	}
	createdSeries, _ := testDB.CreateSeries(ctx, series)

	post := &models.PostCreate{
		TypeID:  "article",
		Title:   "Series Post Test",
		Slug:    "series-post-test",
		Content: "Content",
		Status:  "draft",
	}
	createdPost, _ := testDB.CreatePost(ctx, post)

	// Add post to series
	err := testDB.AddPostToSeries(ctx, createdPost.ID, createdSeries.ID, 1)
	if err != nil {
		t.Fatalf("AddPostToSeries failed: %v", err)
	}

	// Verify post is in series
	posts, err := testDB.GetSeriesPosts(ctx, createdSeries.ID)
	if err != nil {
		t.Fatalf("GetSeriesPosts failed: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("Posts in series: got %d, want 1", len(posts))
	}

	if posts[0].ID != createdPost.ID {
		t.Errorf("Post ID mismatch: got %q, want %q", posts[0].ID, createdPost.ID)
	}
}

func TestGetSeriesPosts(t *testing.T) {
	ctx := context.Background()

	// Create series and multiple posts
	series := &models.SeriesCreate{
		Name:        "Multiple Posts Series",
		Slug:        "multi-posts-series",
		Description: "Test",
	}
	createdSeries, _ := testDB.CreateSeries(ctx, series)

	postIDs := make([]string, 3)
	for i := 0; i < 3; i++ {
		post := &models.PostCreate{
			TypeID:  "article",
			Title:   "Multi Series Post " + string(rune(i)),
			Slug:    "multi-post-" + string(rune(i)),
			Content: "Content",
			Status:  "draft",
		}
		created, _ := testDB.CreatePost(ctx, post)
		postIDs[i] = created.ID
		testDB.AddPostToSeries(ctx, created.ID, createdSeries.ID, i+1)
	}

	// Get posts
	posts, err := testDB.GetSeriesPosts(ctx, createdSeries.ID)
	if err != nil {
		t.Fatalf("GetSeriesPosts failed: %v", err)
	}

	if len(posts) != 3 {
		t.Errorf("Posts count: got %d, want 3", len(posts))
	}

	// Verify ordering
	for i, post := range posts {
		if post.ID != postIDs[i] {
			t.Errorf("Post order mismatch at %d", i)
		}
	}
}

func TestGetPostSeries(t *testing.T) {
	ctx := context.Background()

	// Create post
	post := &models.PostCreate{
		TypeID:  "article",
		Title:   "Multi Series Post",
		Slug:    "multi-series-post",
		Content: "Content",
		Status:  "draft",
	}
	createdPost, _ := testDB.CreatePost(ctx, post)

	// Create multiple series and add post
	seriesIDs := make([]string, 2)
	for i := 0; i < 2; i++ {
		series := &models.SeriesCreate{
			Name:        "Post Multi Series " + string(rune(i)),
			Slug:        "post-multi-series-" + string(rune(i)),
			Description: "Test",
		}
		created, _ := testDB.CreateSeries(ctx, series)
		seriesIDs[i] = created.ID
		testDB.AddPostToSeries(ctx, createdPost.ID, created.ID, i+1)
	}

	// Get post's series
	series, err := testDB.GetPostSeries(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("GetPostSeries failed: %v", err)
	}

	if len(series) != 2 {
		t.Errorf("Series count: got %d, want 2", len(series))
	}
}

func TestRemovePostFromSeries(t *testing.T) {
	ctx := context.Background()

	// Create series and post
	series := &models.SeriesCreate{
		Name:        "Remove Post Series",
		Slug:        "remove-post-series",
		Description: "Test",
	}
	createdSeries, _ := testDB.CreateSeries(ctx, series)

	post := &models.PostCreate{
		TypeID:  "article",
		Title:   "Remove Post Test",
		Slug:    "remove-post-test",
		Content: "Content",
		Status:  "draft",
	}
	createdPost, _ := testDB.CreatePost(ctx, post)

	// Add post to series
	testDB.AddPostToSeries(ctx, createdPost.ID, createdSeries.ID, 1)

	// Verify it's in series
	posts, _ := testDB.GetSeriesPosts(ctx, createdSeries.ID)
	if len(posts) != 1 {
		t.Fatalf("Post should be in series before removal")
	}

	// Remove post from series
	err := testDB.RemovePostFromSeries(ctx, createdPost.ID, createdSeries.ID)
	if err != nil {
		t.Fatalf("RemovePostFromSeries failed: %v", err)
	}

	// Verify removal
	posts, _ = testDB.GetSeriesPosts(ctx, createdSeries.ID)
	if len(posts) != 0 {
		t.Errorf("Posts in series after removal: got %d, want 0", len(posts))
	}
}
