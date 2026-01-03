package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	gen "blog/internal/db/gen"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema embed.FS

type DB struct {
	conn    *sql.DB
	queries *gen.Queries
}

// New creates a new database connection from DATABASE_URL env var
func New(ctx context.Context) (*DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	// Parse URL - handle file:// and libsql:// formats
	if dbURL[:5] == "file:" {
		return NewLocal(ctx, dbURL[5:])
	}

	// For Turso/libsql URLs, we'd need a different driver
	// For now, treat as file path
	conn, err := sql.Open("sqlite3", dbURL[5:])
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		conn:    conn,
		queries: NewQueries(conn),
	}, nil
}

// NewLocal creates a local SQLite connection
func NewLocal(ctx context.Context, path string) (*DB, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		conn:    conn,
		queries: NewQueries(conn),
	}, nil
}

// NewQueries is a wrapper around the sqlc New function to avoid naming conflicts
func NewQueries(db gen.DBTX) *gen.Queries {
	return gen.New(db)
}

// FromSQL creates a DB instance from an existing sql.DB connection
func FromSQL(sqldb *sql.DB) *DB {
	return &DB{
		conn:    sqldb,
		queries: NewQueries(sqldb),
	}
}

// InitSchema initializes the database schema from embedded SQL
func (d *DB) InitSchema(ctx context.Context) error {
	schemaSQL, err := schema.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema: %w", err)
	}

	if _, err := d.conn.ExecContext(ctx, string(schemaSQL)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Schema initialized successfully")
	return nil
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.conn.Close()
}

// Conn returns the underlying sql.DB connection
func (d *DB) Conn() *sql.DB {
	return d.conn
}

// Queries returns the sqlc-generated queries interface
func (d *DB) Queries() *gen.Queries {
	return d.queries
}
