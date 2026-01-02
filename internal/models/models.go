package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Post represents a blog post
type Post struct {
	ID          string          `json:"id"`
	TypeID      string          `json:"type_id"`
	Title       string          `json:"title"`
	Slug        string          `json:"slug"`
	Content     string          `json:"content"`
	Excerpt     string          `json:"excerpt"`
	Status      string          `json:"status"` // draft, published, archived
	IsFeatured  bool            `json:"is_featured"`
	Tags        []string        `json:"tags"`
	Metadata    json.RawMessage `json:"metadata"` // JSON metadata
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	PublishedAt *time.Time      `json:"published_at"`
}

// PostCreate is the request body for creating a post
type PostCreate struct {
	TypeID      string          `json:"type_id"`
	Title       string          `json:"title"`
	Slug        string          `json:"slug"`
	Content     string          `json:"content"`
	Excerpt     string          `json:"excerpt"`
	Tags        []string        `json:"tags"`
	Metadata    json.RawMessage `json:"metadata"`
	IsFeatured  bool            `json:"is_featured"`
	Status      string          `json:"status"`
	PublishedAt *time.Time      `json:"published_at"`
}

// PostUpdate is the request body for updating a post
type PostUpdate struct {
	TypeID      *string         `json:"type_id"`
	Title       *string         `json:"title"`
	Slug        *string         `json:"slug"`
	Content     *string         `json:"content"`
	Excerpt     *string         `json:"excerpt"`
	Tags        []string        `json:"tags"`
	Metadata    json.RawMessage `json:"metadata"`
	IsFeatured  *bool           `json:"is_featured"`
	Status      *string         `json:"status"`
	PublishedAt *time.Time      `json:"published_at"`
}

// PostType represents a post type template
type PostType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// Revision represents a post revision
type Revision struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Series represents a collection of posts
type Series struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// SeriesCreate is the request body for creating a series
type SeriesCreate struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// PostSeries represents a post in a series
type PostSeries struct {
	PostID       string `json:"post_id"`
	SeriesID     string `json:"series_id"`
	OrderInSeries int    `json:"order_in_series"`
}

// ListOptions for filtering and pagination
type ListOptions struct {
	Limit  int
	Offset int
	Type   string // post type filter
	Status string // draft, published, archived
	Tag    string // filter by tag
	Series string // filter by series
}

// TagCount represents a tag with post count
type TagCount struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// Scan helpers for database scanning
func (p *Post) ScanRow(row *sql.Row) error {
	var tags sql.NullString
	var metadata sql.NullString
	var publishedAt sql.NullTime

	err := row.Scan(
		&p.ID, &p.TypeID, &p.Title, &p.Slug, &p.Content,
		&p.Excerpt, &p.Status, &p.IsFeatured, &tags, &metadata,
		&p.CreatedAt, &p.UpdatedAt, &publishedAt,
	)
	if err != nil {
		return err
	}

	if tags.Valid {
		json.Unmarshal([]byte(tags.String), &p.Tags)
	}
	if metadata.Valid {
		p.Metadata = json.RawMessage(metadata.String)
	}
	if publishedAt.Valid {
		p.PublishedAt = &publishedAt.Time
	}

	return nil
}

func (p *Post) ScanRows(rows *sql.Rows) error {
	var tags sql.NullString
	var metadata sql.NullString
	var publishedAt sql.NullTime

	err := rows.Scan(
		&p.ID, &p.TypeID, &p.Title, &p.Slug, &p.Content,
		&p.Excerpt, &p.Status, &p.IsFeatured, &tags, &metadata,
		&p.CreatedAt, &p.UpdatedAt, &publishedAt,
	)
	if err != nil {
		return err
	}

	if tags.Valid {
		json.Unmarshal([]byte(tags.String), &p.Tags)
	}
	if metadata.Valid {
		p.Metadata = json.RawMessage(metadata.String)
	}
	if publishedAt.Valid {
		p.PublishedAt = &publishedAt.Time
	}

	return nil
}
