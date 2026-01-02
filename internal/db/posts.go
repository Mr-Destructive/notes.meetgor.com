package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	gen "blog/internal/db/gen"
	"blog/internal/models"
)

// CreatePost creates a new post
func (d *DB) CreatePost(ctx context.Context, post *models.PostCreate) (*models.Post, error) {
	id := generateID()
	now := time.Now()

	tagsJSON, _ := json.Marshal(post.Tags)
	metadata := string(post.Metadata)
	if metadata == "" {
		metadata = "{}"
	}

	status := post.Status
	if status == "" {
		status = "draft"
	}

	// Convert to sqlc parameters
	params := gen.CreatePostParams{
		ID:          id,
		TypeID:      post.TypeID,
		Title:       post.Title,
		Slug:        post.Slug,
		Content:     post.Content,
		Excerpt:     sql.NullString{String: post.Excerpt, Valid: post.Excerpt != ""},
		Status:      sql.NullString{String: status, Valid: true},
		IsFeatured:  sql.NullBool{Bool: post.IsFeatured, Valid: true},
		Tags:        sql.NullString{String: string(tagsJSON), Valid: true},
		Metadata:    sql.NullString{String: metadata, Valid: true},
		CreatedAt:   sql.NullTime{Time: now, Valid: true},
		UpdatedAt:   sql.NullTime{Time: now, Valid: true},
		PublishedAt: convertTimePtr(post.PublishedAt),
	}

	dbPost, err := d.queries.CreatePost(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return convertPost(dbPost), nil
}

// GetPost retrieves a post by ID or slug
func (d *DB) GetPost(ctx context.Context, idOrSlug string) (*models.Post, error) {
	dbPost, err := d.queries.GetPost(ctx, gen.GetPostParams{
		ID:   idOrSlug,
		Slug: idOrSlug,
	})
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	return convertPost(dbPost), nil
}

// ListPosts lists posts with filters
func (d *DB) ListPosts(ctx context.Context, opts *models.ListOptions) ([]*models.Post, int, error) {
	if opts == nil {
		opts = &models.ListOptions{Limit: 50}
	}
	if opts.Limit == 0 {
		opts.Limit = 50
	}

	// Count posts
	countParams := gen.CountPostsParams{
		Status: nil,
		TypeID: nil,
	}
	if opts.Status != "" {
		countParams.Status = opts.Status
	}
	if opts.Type != "" {
		countParams.TypeID = opts.Type
	}

	total, err := d.queries.CountPosts(ctx, countParams)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// List posts
	listParams := gen.ListPostsParams{
		Status: nil,
		TypeID: nil,
		Offset: int64(opts.Offset),
		Limit:  int64(opts.Limit),
	}
	if opts.Status != "" {
		listParams.Status = opts.Status
	}
	if opts.Type != "" {
		listParams.TypeID = opts.Type
	}

	dbPosts, err := d.queries.ListPosts(ctx, listParams)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list posts: %w", err)
	}

	posts := make([]*models.Post, len(dbPosts))
	for i, p := range dbPosts {
		posts[i] = convertPost(p)
	}

	return posts, int(total), nil
}

// UpdatePost updates a post
func (d *DB) UpdatePost(ctx context.Context, id string, update *models.PostUpdate) (*models.Post, error) {
	// Get current post
	post, err := d.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create revision of old content
	if err := d.CreateRevision(ctx, post.ID, post.Content); err != nil {
		log.Printf("warning: failed to create revision: %v", err)
	}

	// Apply updates
	if update.Title != nil {
		post.Title = *update.Title
	}
	if update.Slug != nil {
		post.Slug = *update.Slug
	}
	if update.Content != nil {
		post.Content = *update.Content
	}
	if update.Excerpt != nil {
		post.Excerpt = *update.Excerpt
	}
	if update.Status != nil {
		post.Status = *update.Status
	}
	if update.IsFeatured != nil {
		post.IsFeatured = *update.IsFeatured
	}
	if len(update.Tags) > 0 {
		post.Tags = update.Tags
	}
	if len(update.Metadata) > 0 {
		post.Metadata = update.Metadata
	}
	if update.PublishedAt != nil {
		post.PublishedAt = update.PublishedAt
	}

	tagsJSON, _ := json.Marshal(post.Tags)
	metadata := string(post.Metadata)
	if metadata == "" {
		metadata = "{}"
	}

	// Manual execution to handle NULL values in COALESCE properly
	query := `
		UPDATE posts SET
		  title = COALESCE(?, title),
		  slug = COALESCE(?, slug),
		  content = COALESCE(?, content),
		  excerpt = COALESCE(?, excerpt),
		  status = COALESCE(?, status),
		  is_featured = COALESCE(?, is_featured),
		  tags = COALESCE(?, tags),
		  metadata = COALESCE(?, metadata),
		  published_at = COALESCE(?, published_at),
		  updated_at = ?
		WHERE id = ?
		RETURNING id, type_id, title, slug, content, excerpt, status, is_featured, tags, metadata, created_at, updated_at, published_at
	`

	row := d.conn.QueryRowContext(ctx, query,
		sql.NullString{String: post.Title, Valid: true},
		sql.NullString{String: post.Slug, Valid: true},
		sql.NullString{String: post.Content, Valid: true},
		sql.NullString{String: post.Excerpt, Valid: post.Excerpt != ""},
		sql.NullString{String: post.Status, Valid: true},
		sql.NullBool{Bool: post.IsFeatured, Valid: true},
		sql.NullString{String: string(tagsJSON), Valid: true},
		sql.NullString{String: metadata, Valid: true},
		convertTimePtr(post.PublishedAt),
		time.Now(),
		id,
	)

	var dbPost gen.Post
	err = row.Scan(
		&dbPost.ID,
		&dbPost.TypeID,
		&dbPost.Title,
		&dbPost.Slug,
		&dbPost.Content,
		&dbPost.Excerpt,
		&dbPost.Status,
		&dbPost.IsFeatured,
		&dbPost.Tags,
		&dbPost.Metadata,
		&dbPost.CreatedAt,
		&dbPost.UpdatedAt,
		&dbPost.PublishedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return convertPost(dbPost), nil
}

// DeletePost deletes a post
func (d *DB) DeletePost(ctx context.Context, id string) error {
	// Check if post exists
	_, err := d.GetPost(ctx, id)
	if err != nil {
		return fmt.Errorf("post not found")
	}

	if err := d.queries.DeletePost(ctx, id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// GetPostTypes retrieves all post types
func (d *DB) GetPostTypes(ctx context.Context) ([]*models.PostType, error) {
	types, err := d.queries.GetPostTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get post types: %w", err)
	}

	result := make([]*models.PostType, len(types))
	for i, t := range types {
		result[i] = &models.PostType{
			ID:          t.ID,
			Name:        t.Name,
			Slug:        t.Slug,
			Description: t.Description.String,
		}
	}

	return result, nil
}

// GetTags retrieves all unique tags with counts
func (d *DB) GetTags(ctx context.Context) ([]*models.TagCount, error) {
	// Get all published posts
	listParams := gen.ListPostsParams{
		Status: "published",
		TypeID: nil,
		Offset: 0,
		Limit:  10000, // Get all
	}

	dbPosts, err := d.queries.ListPosts(ctx, listParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	// Extract tags from posts
	tagMap := make(map[string]int)
	for _, p := range dbPosts {
		if p.Tags.Valid {
			var tags []string
			if err := json.Unmarshal([]byte(p.Tags.String), &tags); err != nil {
				continue
			}
			for _, tag := range tags {
				if tag != "" {
					tagMap[tag]++
				}
			}
		}
	}

	var tags []*models.TagCount
	for tag, count := range tagMap {
		tags = append(tags, &models.TagCount{Tag: tag, Count: count})
	}

	return tags, nil
}

// CreateRevision creates a revision backup
func (d *DB) CreateRevision(ctx context.Context, postID, content string) error {
	id := generateID()
	params := gen.CreateRevisionParams{
		ID:        id,
		PostID:    postID,
		Content:   content,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	return d.queries.CreateRevision(ctx, params)
}

// GetRevisions retrieves post revisions
func (d *DB) GetRevisions(ctx context.Context, postID string) ([]*models.Revision, error) {
	dbRevisions, err := d.queries.GetRevisions(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get revisions: %w", err)
	}

	revisions := make([]*models.Revision, len(dbRevisions))
	for i, r := range dbRevisions {
		revisions[i] = &models.Revision{
			ID:        r.ID,
			PostID:    r.PostID,
			Content:   r.Content,
			CreatedAt: r.CreatedAt.Time,
		}
	}

	return revisions, nil
}

// Slug generation and ID generation helpers
func generateID() string {
	return strings.TrimRight(strings.NewReplacer("+", "-", "/", "_").Replace(
		fmt.Sprintf("%x", time.Now().UnixNano()),
	), "=")
}

// Helper functions to convert between sqlc models and business models
func convertPost(p gen.Post) *models.Post {
	post := &models.Post{
		ID:         p.ID,
		TypeID:     p.TypeID,
		Title:      p.Title,
		Slug:       p.Slug,
		Content:    p.Content,
		IsFeatured: p.IsFeatured.Bool,
	}

	if p.Excerpt.Valid {
		post.Excerpt = p.Excerpt.String
	}
	if p.Status.Valid {
		post.Status = p.Status.String
	}
	if p.Tags.Valid {
		json.Unmarshal([]byte(p.Tags.String), &post.Tags)
	}
	if p.Metadata.Valid {
		post.Metadata = json.RawMessage(p.Metadata.String)
	}
	if p.CreatedAt.Valid {
		post.CreatedAt = p.CreatedAt.Time
	}
	if p.UpdatedAt.Valid {
		post.UpdatedAt = p.UpdatedAt.Time
	}
	if p.PublishedAt.Valid {
		post.PublishedAt = &p.PublishedAt.Time
	}

	return post
}

func convertTimePtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}
