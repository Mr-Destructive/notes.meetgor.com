package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

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

// fetchPageTitle fetches the title from a URL by reading its HTML meta tags
func fetchPageTitle(url string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible)")
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return ""
	}
	defer resp.Body.Close()

	// Read limited amount of HTML to find og:title or <title>
	limitReader := io.LimitReader(resp.Body, 50000)
	body, _ := io.ReadAll(limitReader)
	html := string(body)

	// Try og:title first
	ogPattern := regexp.MustCompile(`<meta\s+property=["']og:title["']\s+content=["']([^"']+)["']`)
	if matches := ogPattern.FindStringSubmatch(html); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// Try title tag
	titlePattern := regexp.MustCompile(`<title>([^<]+)</title>`)
	if matches := titlePattern.FindStringSubmatch(html); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}

// initDB initializes the database on first request
func initDB() error {
	if database != nil {
		return nil
	}

	d, err := db.New(context.Background())
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize schema with soft fail
	if err := d.InitSchema(context.Background()); err != nil {
		log.Printf("Schema initialization warning (non-fatal): %v", err)
	}

	database = d
	return nil
}

// Handler is the main Netlify function handler
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s", r.Method, r.URL.Path)

	// Initialize database on first request
	if err := initDB(); err != nil {
		log.Printf("Database initialization error: %v", err)
		respondError(w, http.StatusInternalServerError, "Database connection failed")
		return
	}
	
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
	
	// Root path - serve admin dashboard
	if path == "" {
		serveAdminIndex(w, r)
		return
	}

	// API health endpoint
	if path == "api" {
		respondJSON(w, http.StatusOK, map[string]string{"status": "ok", "version": "1.0"})
		return
	}

	// CSS file serving
	if strings.HasPrefix(path, "css/") {
		serveCSSFile(w, r, path)
		return
	}

	// Database initialization happens on demand in individual handlers

	// Admin routes
	if strings.HasPrefix(path, "admin/") {
		path = strings.TrimPrefix(path, "admin/")
		path = strings.Trim(path, "/")
		
		parts := strings.Split(path, "/")
		resource := parts[0]
		var id, action string
		if len(parts) > 1 {
			id = parts[1]
		}
		if len(parts) > 2 {
			action = parts[2]
		}

		switch resource {
		case "dashboard":
			h.HandleAdminDashboard(w, r, database)
		case "posts":
			if action == "new" || (id != "" && action == "edit") {
				h.HandlePostEditor(w, r, database, id)
			} else {
				h.HandlePostsList(w, r, database)
			}
		case "series":
			if action == "new" || (id != "" && action == "edit") {
				// TODO: HandleSeriesEditor
				respondError(w, http.StatusNotImplemented, "Series editor not yet implemented")
			} else {
				h.HandleSeriesList(w, r, database)
			}
		case "types":
			h.HandlePostTypes(w, r, database)
		case "exports":
			h.HandleExportPage(w, r, database)
		default:
			respondError(w, http.StatusNotFound, "Admin page not found")
		}
		return
	}

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

		// Fetch title from link if empty and type is "link"
		if (req.Title == "" || req.Title == nil) && req.TypeID == "link" {
			if metadata, ok := req.Metadata.(map[string]interface{}); ok {
				if sourceURL, exists := metadata["source_url"].(string); exists && sourceURL != "" {
					if title := fetchPageTitle(sourceURL); title != "" {
						req.Title = title
					}
				}
			}
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

// serveCSSFile serves CSS files from public/css/
func serveCSSFile(w http.ResponseWriter, r *http.Request, path string) {
	// Simple CSS serving - in production this would be embedded or served by Netlify
	filename := strings.TrimPrefix(path, "css/")
	
	// Only allow admin.css for now
	if filename == "admin.css" {
		css := `* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

:root {
    --primary: #1a1a1a;
    --secondary: #666;
    --border: #e0e0e0;
    --bg-light: #f9f9f9;
    --bg-white: #ffffff;
    --text: #333;
    --text-light: #666;
    --success: #27ae60;
    --warning: #f39c12;
    --danger: #e74c3c;
}

html, body {
    height: 100%;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    font-size: 14px;
    line-height: 1.6;
    color: var(--text);
    background: var(--bg-white);
}

.card {
    background: var(--bg-white);
    border: 1px solid var(--border);
    border-radius: 4px;
    margin-bottom: 24px;
    overflow: hidden;
}

.card-header {
    padding: 20px;
    border-bottom: 1px solid var(--border);
    background: var(--bg-light);
}

.card-header h3 {
    font-size: 16px;
    font-weight: 600;
    margin: 0;
}

.card-body {
    padding: 20px;
}

.form-group {
    margin-bottom: 20px;
}

.form-group label {
    display: block;
    margin-bottom: 6px;
    font-weight: 500;
    font-size: 13px;
}

.form-group input,
.form-group textarea,
.form-group select {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid var(--border);
    border-radius: 3px;
    font-family: inherit;
    font-size: 13px;
    transition: border-color 0.2s;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
    outline: none;
    border-color: var(--primary);
    box-shadow: inset 0 0 0 1px var(--primary);
}

.form-group textarea {
    resize: vertical;
    min-height: 200px;
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 12px;
}

.btn {
    padding: 10px 16px;
    border: 1px solid transparent;
    border-radius: 3px;
    cursor: pointer;
    font-size: 13px;
    font-weight: 500;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    transition: all 0.2s;
    white-space: nowrap;
}

.btn-primary {
    background: var(--primary);
    color: white;
}

.btn-success {
    background: var(--success);
    color: white;
}

.btn-danger {
    background: var(--danger);
    color: white;
}

.btn-outline {
    background: transparent;
    border-color: var(--border);
    color: var(--text);
}

.btn-outline:hover {
    background: var(--bg-light);
}

.btn-sm {
    padding: 6px 12px;
    font-size: 12px;
}

.table {
    width: 100%;
    border-collapse: collapse;
}

.table th {
    padding: 12px;
    text-align: left;
    font-weight: 600;
    border-bottom: 2px solid var(--border);
    background: var(--bg-light);
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.table td {
    padding: 12px;
    border-bottom: 1px solid var(--border);
    font-size: 13px;
}

.table tbody tr:hover {
    background: var(--bg-light);
}

.table-actions {
    display: flex;
    gap: 8px;
}

.search-bar {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
}

.search-bar select {
    padding: 10px 12px;
    border: 1px solid var(--border);
    border-radius: 3px;
    font-size: 13px;
}

.badge {
    display: inline-flex;
    align-items: center;
    padding: 4px 10px;
    border-radius: 3px;
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.badge-success {
    background: #d5f4e6;
    color: #27ae60;
}

.badge-warning {
    background: #fdebd0;
    color: #f39c12;
}

.badge-danger {
    background: #fadbd8;
    color: #e74c3c;
}

.alert {
    padding: 12px 16px;
    border-radius: 3px;
    margin-bottom: 20px;
    border-left: 4px solid;
    font-size: 13px;
}

.alert-success {
    background: #d5f4e6;
    color: #27ae60;
    border-color: #27ae60;
}

.alert-danger {
    background: #fadbd8;
    color: #e74c3c;
    border-color: #e74c3c;
}`
		w.Header().Set("Content-Type", "text/css")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, css)
		return
	}
	
	respondError(w, http.StatusNotFound, "CSS file not found")
}

// serveAdminIndex serves the admin dashboard HTML
func serveAdminIndex(w http.ResponseWriter, r *http.Request) {
	// Read and serve the admin index
	indexHTML := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Blog Admin</title>
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<link rel="stylesheet" href="/css/admin.css">
	<style>
		body {
			margin: 0;
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
			background: #f5f5f5;
		}
		
		.admin-shell {
			display: flex;
			height: 100vh;
			overflow: hidden;
		}
		
		.sidebar {
			width: 220px;
			background: #1a1a1a;
			color: white;
			padding: 20px;
			overflow-y: auto;
			border-right: 1px solid #333;
		}
		
		.sidebar h2 {
			margin: 0 0 30px 0;
			font-size: 18px;
			font-weight: 600;
			border-bottom: 1px solid #333;
			padding-bottom: 15px;
		}
		
		.nav-section {
			margin-bottom: 30px;
		}
		
		.nav-section-title {
			font-size: 11px;
			font-weight: 700;
			text-transform: uppercase;
			color: #888;
			margin-bottom: 12px;
			letter-spacing: 1px;
		}
		
		.nav-link {
			display: block;
			padding: 10px 12px;
			color: #ccc;
			text-decoration: none;
			border-radius: 4px;
			margin-bottom: 4px;
			font-size: 14px;
			cursor: pointer;
			border: none;
			background: none;
			text-align: left;
			width: 100%;
		}
		
		.nav-link:hover {
			background: #333;
			color: white;
		}
		
		.nav-link.active {
			background: #0066cc;
			color: white;
			font-weight: 600;
		}
		
		.main-content {
			flex: 1;
			overflow: hidden;
			display: flex;
			flex-direction: column;
		}
		
		.topbar {
			background: white;
			border-bottom: 1px solid #e0e0e0;
			padding: 15px 30px;
			display: flex;
			justify-content: space-between;
			align-items: center;
		}
		
		.topbar h1 {
			margin: 0;
			font-size: 24px;
			font-weight: 600;
		}
		
		.content-area {
			flex: 1;
			overflow-y: auto;
			padding: 30px;
		}
		
		.htmx-request.htmx-settling #main-content {
			opacity: 0.6;
		}
		
		.htmx-request.htmx-settling .spinner {
			display: inline-block;
		}
		
		.spinner {
			display: none;
			width: 16px;
			height: 16px;
			border: 2px solid #f3f3f3;
			border-top: 2px solid #0066cc;
			border-radius: 50%;
			animation: spin 1s linear infinite;
		}
		
		@keyframes spin {
			0% { transform: rotate(0deg); }
			100% { transform: rotate(360deg); }
		}
		
		@media (max-width: 768px) {
			.admin-shell {
				flex-direction: column;
			}
			
			.sidebar {
				width: 100%;
				max-height: 60px;
				padding: 10px 20px;
				display: flex;
				align-items: center;
				justify-content: space-between;
				overflow-x: auto;
				overflow-y: hidden;
			}
			
			.sidebar h2 {
				margin: 0;
				font-size: 14px;
				border: none;
				padding: 0;
			}
			
			.nav-section {
				margin: 0;
				display: flex;
				gap: 10px;
			}
			
			.nav-section-title {
				display: none;
			}
			
			.content-area {
				padding: 20px;
			}
		}
	</style>
</head>
<body>
	<div class="admin-shell">
		<div class="sidebar">
			<h2>Blog Admin</h2>
			
			<div class="nav-section">
				<div class="nav-section-title">Main</div>
				<a class="nav-link" hx-get="/admin/dashboard" hx-target="#main-content" onclick="updateActiveNav(this)">📊 Dashboard</a>
				<a class="nav-link" hx-get="/admin/posts" hx-target="#main-content" onclick="updateActiveNav(this)">📝 Posts</a>
				<a class="nav-link" hx-get="/admin/series" hx-target="#main-content" onclick="updateActiveNav(this)">📚 Series</a>
			</div>
			
			<div class="nav-section">
				<div class="nav-section-title">Config</div>
				<a class="nav-link" hx-get="/admin/types" hx-target="#main-content" onclick="updateActiveNav(this)">🏷️ Post Types</a>
				<a class="nav-link" hx-get="/admin/exports" hx-target="#main-content" onclick="updateActiveNav(this)">📤 Export</a>
			</div>
		</div>
		
		<div class="main-content">
			<div class="topbar">
				<h1 id="page-title">Dashboard</h1>
				<div class="spinner"></div>
			</div>
			
			<div class="content-area" id="main-content">
				<div class="card">
					<div class="card-body" style="text-align: center; padding: 60px 20px;">
						<p style="color: #999; font-size: 16px;">Loading...</p>
					</div>
				</div>
			</div>
		</div>
	</div>

	<script src="https://cdn.jsdelivr.net/npm/marked@11.1.1/marked.min.js"></script>
	<script>
		function updateActiveNav(element) {
			document.querySelectorAll('.nav-link').forEach(el => el.classList.remove('active'));
			element.classList.add('active');
			
			const text = element.textContent.trim();
			document.getElementById('page-title').textContent = text;
		}

		// Load dashboard on startup
		document.addEventListener('DOMContentLoaded', () => {
			htmx.ajax('GET', '/admin/dashboard', { target: '#main-content' });
			document.querySelector('.nav-link').classList.add('active');
		});
	</script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, indexHTML)
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
