package main

import (
	"context"
	"log"

	"blog/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	ctx := context.Background()

	// Initialize database
	database, err := db.New(ctx)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize schema
	if err := database.InitSchema(ctx); err != nil {
		log.Fatalf("failed to initialize schema: %v", err)
	}

	log.Println("✓ Database initialized successfully")
	log.Println("✓ Schema created")

	// List post types to verify
	types, err := database.GetPostTypes(ctx)
	if err != nil {
		log.Fatalf("failed to get post types: %v", err)
	}

	log.Printf("✓ Found %d post types:", len(types))
	for _, t := range types {
		log.Printf("  - %s (%s)", t.Name, t.ID)
	}
}
