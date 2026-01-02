package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	gen "blog/internal/db/gen"
	"blog/internal/models"
)

// CreateSeries creates a new series
func (d *DB) CreateSeries(ctx context.Context, series *models.SeriesCreate) (*models.Series, error) {
	id := generateID()
	now := time.Now()

	params := gen.CreateSeriesParams{
		ID:          id,
		Name:        series.Name,
		Slug:        series.Slug,
		Description: sql.NullString{String: series.Description, Valid: series.Description != ""},
		CreatedAt:   sql.NullTime{Time: now, Valid: true},
	}

	dbSeries, err := d.queries.CreateSeries(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create series: %w", err)
	}

	return convertSeries(dbSeries), nil
}

// GetSeries retrieves a series by ID or slug
func (d *DB) GetSeries(ctx context.Context, idOrSlug string) (*models.Series, error) {
	dbSeries, err := d.queries.GetSeries(ctx, gen.GetSeriesParams{
		ID:   idOrSlug,
		Slug: idOrSlug,
	})
	if err != nil {
		return nil, fmt.Errorf("series not found: %w", err)
	}

	return convertSeries(dbSeries), nil
}

// ListSeries lists all series
func (d *DB) ListSeries(ctx context.Context, limit, offset int) ([]*models.Series, error) {
	if limit == 0 {
		limit = 50
	}

	params := gen.ListSeriesParams{
		Offset: int64(offset),
		Limit:  int64(limit),
	}

	dbSeries, err := d.queries.ListSeries(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list series: %w", err)
	}

	series := make([]*models.Series, len(dbSeries))
	for i, s := range dbSeries {
		series[i] = convertSeries(s)
	}

	return series, nil
}

// UpdateSeries updates a series
func (d *DB) UpdateSeries(ctx context.Context, id string, name, slug, description string) (*models.Series, error) {
	params := gen.UpdateSeriesParams{
		ID:          id,
		Name:        name,
		Slug:        slug,
		Description: sql.NullString{String: description, Valid: description != ""},
	}

	dbSeries, err := d.queries.UpdateSeries(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update series: %w", err)
	}

	return convertSeries(dbSeries), nil
}

// DeleteSeries deletes a series
func (d *DB) DeleteSeries(ctx context.Context, id string) error {
	// Check if series exists
	_, err := d.GetSeries(ctx, id)
	if err != nil {
		return fmt.Errorf("series not found")
	}

	if err := d.queries.DeleteSeries(ctx, id); err != nil {
		return fmt.Errorf("failed to delete series: %w", err)
	}

	return nil
}

// AddPostToSeries adds a post to a series
func (d *DB) AddPostToSeries(ctx context.Context, postID, seriesID string, order int) error {
	// Manual execution because sqlc doesn't handle the repeated parameter in ON CONFLICT
	query := `
		INSERT INTO post_series (post_id, series_id, order_in_series)
		VALUES (?, ?, ?)
		ON CONFLICT(post_id, series_id) DO UPDATE SET order_in_series = ?
	`
	_, err := d.conn.ExecContext(ctx, query, postID, seriesID, order, order)
	if err != nil {
		return fmt.Errorf("failed to add post to series: %w", err)
	}

	return nil
}

// RemovePostFromSeries removes a post from a series
func (d *DB) RemovePostFromSeries(ctx context.Context, postID, seriesID string) error {
	params := gen.RemovePostFromSeriesParams{
		PostID:   postID,
		SeriesID: seriesID,
	}

	if err := d.queries.RemovePostFromSeries(ctx, params); err != nil {
		return fmt.Errorf("failed to remove post from series: %w", err)
	}

	return nil
}

// GetSeriesPosts retrieves all posts in a series
func (d *DB) GetSeriesPosts(ctx context.Context, seriesID string) ([]*models.Post, error) {
	dbPosts, err := d.queries.GetSeriesPosts(ctx, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to get series posts: %w", err)
	}

	posts := make([]*models.Post, len(dbPosts))
	for i, p := range dbPosts {
		posts[i] = convertPost(p)
	}

	return posts, nil
}

// GetPostSeries retrieves all series a post belongs to
func (d *DB) GetPostSeries(ctx context.Context, postID string) ([]*models.Series, error) {
	dbSeries, err := d.queries.GetPostSeries(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post series: %w", err)
	}

	series := make([]*models.Series, len(dbSeries))
	for i, s := range dbSeries {
		series[i] = convertSeries(s)
	}

	return series, nil
}

// Helper function to convert sqlc Series to business Series
func convertSeries(s gen.Series) *models.Series {
	series := &models.Series{
		ID:   s.ID,
		Name: s.Name,
		Slug: s.Slug,
	}

	if s.Description.Valid {
		series.Description = s.Description.String
	}
	if s.CreatedAt.Valid {
		series.CreatedAt = s.CreatedAt.Time
	}

	return series
}
