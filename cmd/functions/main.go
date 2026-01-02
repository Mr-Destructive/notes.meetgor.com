package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"blog/internal/db"
	"blog/internal/handler"
	"blog/internal/models"
	"blog/internal/ssg"

	"github.com/joho/godotenv"
)

var database *db.DB

func init() {
	// Load .env for local development
	godotenv.Load()

	// Initialize database on cold start
	var err error
	ctx := context.Background()
	database, err = db.New(ctx)
	if err != nil {
		log.Printf("warning: database init deferred: %v", err)
	}
}

// Handler is the main Netlify function handler
func Handler(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Initialize database if needed (cold start)
	if database == nil {
		ctx := context.Background()
		var err error
		database, err = db.New(ctx)
		if err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
			return
		}
	}

	// Parse route
	path := strings.TrimPrefix(r.URL.Path, "/.netlify/functions/cms")
	path = strings.Trim(path, "/")
	
	// Root path - serve index.html for admin
	if path == "" {
		serveFile(w, "public/index.html")
		return
	}

	// Admin routes
	if strings.HasPrefix(path, "admin") {
		path = strings.TrimPrefix(path, "admin/")
		parts := strings.Split(path, "/")
		
		resource := ""
		var id, action string
		if len(parts) > 0 && parts[0] != "" {
			resource = parts[0]
		}
		if len(parts) > 1 {
			id = parts[1]
		}
		if len(parts) > 2 {
			action = parts[2]
		}

		switch resource {
		case "dashboard":
			handler.HandleAdminDashboard(w, r, database)
		case "posts":
			if action == "edit" {
				// TODO: Handle post edit view
				respondError(w, http.StatusNotImplemented, "Post editor not yet implemented")
			} else if id == "new" {
				// TODO: Handle new post form
				respondError(w, http.StatusNotImplemented, "New post form not yet implemented")
			} else {
				handler.HandlePostsList(w, r, database)
			}
		case "series":
			if action == "edit" {
				// TODO: Handle series edit view
				respondError(w, http.StatusNotImplemented, "Series editor not yet implemented")
			} else if id == "new" {
				// TODO: Handle new series form
				respondError(w, http.StatusNotImplemented, "New series form not yet implemented")
			} else {
				handler.HandleSeriesList(w, r, database)
			}
		case "types":
			handler.HandlePostTypes(w, r, database)
		case "export":
			handler.HandleExportPage(w, r, database)
		default:
			handler.HandleAdminDashboard(w, r, database)
		}
		return
	}

	// API routes
	path = strings.TrimPrefix(path, "api")
	path = strings.Trim(path, "/")
	
	if path == "" {
		respondJSON(w, http.StatusOK, map[string]string{"status": "ok", "version": "1.0"})
		return
	}

	parts := strings.Split(path, "/")

	resource := parts[0]
	var id, action string
	if len(parts) > 1 {
		id = parts[1]
	}
	if len(parts) > 2 {
		action = parts[2]
	}

	// Route to handlers
	switch resource {
	case "auth":
		handleAuth(w, r, database, id)

	case "posts":
		handlePosts(w, r, database, id, action)

	case "series":
		handleSeries(w, r, database, id, action)

	case "types":
		handleTypes(w, r, database)

	case "tags":
		handleTags(w, r, database)

	case "exports":
		handleExports(w, r, database)

	default:
		respondError(w, http.StatusNotFound, "Resource not found")
	}
}

// handleAuth handles authentication endpoints
func handleAuth(w http.ResponseWriter, r *http.Request, db *db.DB, action string) {
	switch action {
	case "login":
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handler.HandleLogin(w, r, db)

	case "logout":
		handler.HandleLogout(w, r)

	case "verify":
		handler.HandleVerify(w, r)

	default:
		respondError(w, http.StatusNotFound, "Auth endpoint not found")
	}
}

// handlePosts handles post CRUD endpoints
func handlePosts(w http.ResponseWriter, r *http.Request, db *db.DB, id, action string) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			// List posts
			limit := 50
			offset := 0
			if q := r.URL.Query().Get("limit"); q != "" {
				fmt.Sscanf(q, "%d", &limit)
			}
			if q := r.URL.Query().Get("offset"); q != "" {
				fmt.Sscanf(q, "%d", &offset)
			}

			opts := &models.ListOptions{
				Limit:  limit,
				Offset: offset,
				Type:   r.URL.Query().Get("type"),
				Status: r.URL.Query().Get("status"),
				Tag:    r.URL.Query().Get("tag"),
				Series: r.URL.Query().Get("series"),
			}

			posts, total, err := db.ListPosts(ctx, opts)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			respondJSON(w, http.StatusOK, map[string]interface{}{
				"posts": posts,
				"total": total,
			})
		} else {
			// Get single post
			post, err := db.GetPost(ctx, id)
			if err != nil {
				respondError(w, http.StatusNotFound, "Post not found")
				return
			}
			respondJSON(w, http.StatusOK, post)
		}

	case http.MethodPost:
		// Create post
		var req models.PostCreate
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		post, err := db.CreatePost(ctx, &req)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)

	case http.MethodPut:
		// Update post
		if id == "" {
			respondError(w, http.StatusBadRequest, "Post ID required")
			return
		}

		var req models.PostUpdate
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		post, err := db.UpdatePost(ctx, id, &req)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, post)

	case http.MethodDelete:
		// Delete post
		if id == "" {
			respondError(w, http.StatusBadRequest, "Post ID required")
			return
		}

		if err := db.DeletePost(ctx, id); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})

	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleSeries handles series endpoints
func handleSeries(w http.ResponseWriter, r *http.Request, db *db.DB, id, action string) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			// List series
			limit := 50
			offset := 0
			if q := r.URL.Query().Get("limit"); q != "" {
				fmt.Sscanf(q, "%d", &limit)
			}
			if q := r.URL.Query().Get("offset"); q != "" {
				fmt.Sscanf(q, "%d", &offset)
			}

			series, err := db.ListSeries(ctx, limit, offset)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			respondJSON(w, http.StatusOK, series)
		} else if action == "posts" {
			// Get posts in series
			posts, err := db.GetSeriesPosts(ctx, id)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondJSON(w, http.StatusOK, posts)
		} else {
			// Get single series
			series, err := db.GetSeries(ctx, id)
			if err != nil {
				respondError(w, http.StatusNotFound, "Series not found")
				return
			}
			respondJSON(w, http.StatusOK, series)
		}

	case http.MethodPost:
		var req models.SeriesCreate
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		series, err := db.CreateSeries(ctx, &req)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(series)

	case http.MethodDelete:
		if id == "" {
			respondError(w, http.StatusBadRequest, "Series ID required")
			return
		}

		if err := db.DeleteSeries(ctx, id); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})

	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleTypes returns post types
func handleTypes(w http.ResponseWriter, r *http.Request, db *db.DB) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	types, err := db.GetPostTypes(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, types)
}

// handleTags returns all tags
func handleTags(w http.ResponseWriter, r *http.Request, db *db.DB) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	tags, err := db.GetTags(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, tags)
}

// handleExports handles data export and markdown generation
func handleExports(w http.ResponseWriter, r *http.Request, db *db.DB) {
	if r.Method == http.MethodGet {
		handleExportsGet(w, r, db)
	} else if r.Method == http.MethodPost {
		handleExportsPost(w, r, db)
	} else {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleExportsGet returns JSON export of posts
func handleExportsGet(w http.ResponseWriter, r *http.Request, db *db.DB) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		status = "published"
	}

	ctx := r.Context()
	posts, _, err := db.ListPosts(ctx, &models.ListOptions{
		Limit:  1000,
		Status: status,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, posts)
}

// handleExportsPost exports posts as markdown for static site generation
func handleExportsPost(w http.ResponseWriter, r *http.Request, db *db.DB) {
	// Parse export type from URL
	r.ParseForm()
	exportType := r.Form.Get("type")
	if exportType == "" {
		exportType = "markdown"
	}

	// Create temporary output directory
	outputDir := "./exports"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create export directory")
		return
	}

	ctx := r.Context()
	result, err := ssg.ExportAll(ctx, db, outputDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Export failed: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// Helper response functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// serveFile serves a static file
func serveFile(w http.ResponseWriter, filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "File not found")
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func main() {
	// For local development
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting CMS server on :%s", port)
	if err := http.ListenAndServe(":"+port, http.HandlerFunc(Handler)); err != nil {
		log.Fatal(err)
	}
}
