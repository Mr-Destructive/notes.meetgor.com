package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"blog/internal/db"
	h "blog/internal/handler"
	"blog/internal/models"
	"blog/internal/ssg"

	"github.com/joho/godotenv"
)

var database *db.DB

func init() {
	godotenv.Load()
	log.Println("Function initializing...")
}

// Handler is the main Netlify function handler
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s", r.Method, r.URL.Path)
	
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse route
	path := strings.TrimPrefix(r.URL.Path, "/.netlify/functions/cms")
	path = strings.Trim(path, "/")
	
	// Root path - health check
	if path == "" {
		respondJSON(w, http.StatusOK, map[string]string{"status": "ok", "message": "CMS is running"})
		return
	}

	// API health endpoint
	if path == "api" {
		respondJSON(w, http.StatusOK, map[string]string{"status": "ok", "version": "1.0"})
		return
	}

	// Database initialization happens on demand in individual handlers

	// API routes
	if strings.HasPrefix(path, "api/") {
		path = strings.TrimPrefix(path, "api/")
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
		return
	}

	respondError(w, http.StatusNotFound, "Not found")
}



// handleAuth handles authentication endpoints
func handleAuth(w http.ResponseWriter, r *http.Request, db *db.DB, action string) {
	switch action {
	case "login":
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		h.HandleLogin(w, r, db)
	case "logout":
		h.HandleLogout(w, r)
	case "verify":
		h.HandleVerify(w, r)
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
			post, err := db.GetPost(ctx, id)
			if err != nil {
				respondError(w, http.StatusNotFound, "Post not found")
				return
			}
			respondJSON(w, http.StatusOK, post)
		}

	case http.MethodPost:
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
			posts, err := db.GetSeriesPosts(ctx, id)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondJSON(w, http.StatusOK, posts)
		} else {
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

// main() is only for local development
// Netlify Functions automatically exports Handler
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting CMS server on :%s (local dev only)", port)
	if err := http.ListenAndServe(":"+port, http.HandlerFunc(Handler)); err != nil {
		log.Fatal(err)
	}
}
