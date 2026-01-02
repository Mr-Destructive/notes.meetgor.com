package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"blog/internal/db"
	"blog/internal/models"
)

// HandleAdminDashboard serves the main dashboard
func HandleAdminDashboard(w http.ResponseWriter, r *http.Request, database *db.DB) {
	ctx := context.Background()

	// Get stats
	posts, _, err := database.ListPosts(ctx, &models.ListOptions{Limit: 100})
	if err != nil {
		renderHTML(w, "500", "Error loading posts")
		return
	}

	types, err := database.GetPostTypes(ctx)
	if err != nil {
		renderHTML(w, "500", "Error loading post types")
		return
	}

	tags, err := database.GetTags(ctx)
	if err != nil {
		renderHTML(w, "500", "Error loading tags")
		return
	}

	series, err := database.ListSeries(ctx, 100, 0)
	if err != nil {
		renderHTML(w, "500", "Error loading series")
		return
	}

	// Count by status
	draft := 0
	published := 0
	archived := 0
	for _, p := range posts {
		switch p.Status {
		case "draft":
			draft++
		case "published":
			published++
		case "archived":
			archived++
		}
	}

	html := `
<div class="card">
	<div class="card-header">
		<h3>Dashboard</h3>
	</div>
	<div class="card-body">
		<div class="stats-grid" style="display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 20px; margin-bottom: 30px;">
			<div class="stat-card" style="background: white; padding: 20px; border: 1px solid #ddd; border-radius: 4px;">
				<div style="font-size: 32px; font-weight: bold; color: #0066cc;">` + strconv.Itoa(len(posts)) + `</div>
				<div style="color: #666; font-size: 14px;">Total Posts</div>
			</div>
			<div class="stat-card" style="background: white; padding: 20px; border: 1px solid #ddd; border-radius: 4px;">
				<div style="font-size: 32px; font-weight: bold; color: #28a745;">` + strconv.Itoa(published) + `</div>
				<div style="color: #666; font-size: 14px;">Published</div>
			</div>
			<div class="stat-card" style="background: white; padding: 20px; border: 1px solid #ddd; border-radius: 4px;">
				<div style="font-size: 32px; font-weight: bold; color: #ffc107;">` + strconv.Itoa(draft) + `</div>
				<div style="color: #666; font-size: 14px;">Draft</div>
			</div>
			<div class="stat-card" style="background: white; padding: 20px; border: 1px solid #ddd; border-radius: 4px;">
				<div style="font-size: 32px; font-weight: bold; color: #17a2b8;">` + strconv.Itoa(len(series)) + `</div>
				<div style="color: #666; font-size: 14px;">Series</div>
			</div>
			<div class="stat-card" style="background: white; padding: 20px; border: 1px solid #ddd; border-radius: 4px;">
				<div style="font-size: 32px; font-weight: bold; color: #6c757d;">` + strconv.Itoa(len(types)) + `</div>
				<div style="color: #666; font-size: 14px;">Post Types</div>
			</div>
			<div class="stat-card" style="background: white; padding: 20px; border: 1px solid #ddd; border-radius: 4px;">
				<div style="font-size: 32px; font-weight: bold; color: #007bff;">` + strconv.Itoa(len(tags)) + `</div>
				<div style="color: #666; font-size: 14px;">Tags</div>
			</div>
		</div>
	</div>
</div>

<div class="card">
	<div class="card-header">
		<h3>Recent Posts</h3>
	</div>
	<div class="card-body">
		<table class="table">
			<thead>
				<tr>
					<th>Title</th>
					<th>Type</th>
					<th>Status</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
	`

	// Show last 5 posts
	for i, post := range posts {
		if i >= 5 {
			break
		}
		statusBadge := `<span class="badge badge-warning">Draft</span>`
		if post.Status == "published" {
			statusBadge = `<span class="badge badge-success">Published</span>`
		} else if post.Status == "archived" {
			statusBadge = `<span class="badge badge-danger">Archived</span>`
		}

		html += `
				<tr>
					<td><strong>` + post.Title + `</strong></td>
					<td>` + post.TypeID + `</td>
					<td>` + statusBadge + `</td>
					<td>` + post.CreatedAt.Format("Jan 2, 2006") + `</td>
					<td>
						<div class="table-actions">
							<button class="btn btn-sm btn-outline" hx-get="/admin/posts/` + post.ID + `/edit" hx-target="#main-content">Edit</button>
							<button class="btn btn-sm btn-danger" hx-delete="/api/posts/` + post.ID + `" hx-confirm="Delete this post?">Delete</button>
						</div>
					</td>
				</tr>
		`
	}

	html += `
			</tbody>
		</table>
	</div>
</div>
	`

	renderHTML(w, "text/html", html)
}

// HandlePostsList serves the posts list view
func HandlePostsList(w http.ResponseWriter, r *http.Request, database *db.DB) {
	ctx := context.Background()

	// Get filters from query params
	limit := 20
	offset := 0
	status := r.URL.Query().Get("status")
	postType := r.URL.Query().Get("type")

	// Parse pagination
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	opts := &models.ListOptions{
		Limit:  limit,
		Offset: offset,
		Status: status,
		Type:   postType,
	}

	posts, total, err := database.ListPosts(ctx, opts)
	if err != nil {
		renderHTML(w, "text/html", `<div class="alert alert-danger">Error loading posts</div>`)
		return
	}

	types, _ := database.GetPostTypes(ctx)

	html := `
<div class="card">
	<div class="card-header">
		<div style="display: flex; justify-content: space-between; align-items: center;">
			<h3>Posts</h3>
			<button class="btn btn-primary" hx-get="/admin/posts/new" hx-target="#main-content">+ New Post</button>
		</div>
	</div>
	<div class="card-body">
		<div class="search-bar">
			<select hx-get="/admin/posts" hx-target="#main-content" name="type" hx-trigger="change">
				<option value="">All Types</option>
	`

	for _, t := range types {
		selected := ""
		if t.ID == postType {
			selected = "selected"
		}
		html += `<option value="` + t.ID + `" ` + selected + `>` + t.Name + `</option>`
	}

	html += `
			</select>
			<select hx-get="/admin/posts" hx-target="#main-content" name="status" hx-trigger="change">
				<option value="">All Statuses</option>
				<option value="draft"`

	if status == "draft" {
		html += ` selected`
	}

	html += `>Draft</option>
				<option value="published"`

	if status == "published" {
		html += ` selected`
	}

	html += `>Published</option>
				<option value="archived"`

	if status == "archived" {
		html += ` selected`
	}

	html += `>Archived</option>
			</select>
		</div>

		<table class="table">
			<thead>
				<tr>
					<th>Title</th>
					<th>Type</th>
					<th>Status</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
	`

	for _, post := range posts {
		statusBadge := `<span class="badge badge-warning">Draft</span>`
		if post.Status == "published" {
			statusBadge = `<span class="badge badge-success">Published</span>`
		} else if post.Status == "archived" {
			statusBadge = `<span class="badge badge-danger">Archived</span>`
		}

		html += `
				<tr>
					<td><strong>` + post.Title + `</strong><br><small style="color: #666;">` + post.Slug + `</small></td>
					<td>` + post.TypeID + `</td>
					<td>` + statusBadge + `</td>
					<td>` + post.CreatedAt.Format("Jan 2, 2006 15:04") + `</td>
					<td>
						<div class="table-actions">
							<button class="btn btn-sm btn-outline" hx-get="/admin/posts/` + post.ID + `/edit" hx-target="#main-content">Edit</button>
							<button class="btn btn-sm btn-danger" hx-delete="/api/posts/` + post.ID + `" hx-confirm="Delete this post?" hx-target="#main-content">Delete</button>
						</div>
					</td>
				</tr>
		`
	}

	html += `
			</tbody>
		</table>

		<div style="text-align: center; margin-top: 20px; color: #666;">
			Showing ` + strconv.Itoa(len(posts)) + ` of ` + strconv.Itoa(total) + ` posts
		</div>
	</div>
</div>
	`

	renderHTML(w, "text/html", html)
}

// HandleSeriesList serves the series list view
func HandleSeriesList(w http.ResponseWriter, r *http.Request, database *db.DB) {
	ctx := context.Background()

	series, err := database.ListSeries(ctx, 100, 0)
	if err != nil {
		renderHTML(w, "text/html", `<div class="alert alert-danger">Error loading series</div>`)
		return
	}

	html := `
<div class="card">
	<div class="card-header">
		<div style="display: flex; justify-content: space-between; align-items: center;">
			<h3>Series</h3>
			<button class="btn btn-primary" hx-get="/admin/series/new" hx-target="#main-content">+ New Series</button>
		</div>
	</div>
	<div class="card-body">
		<table class="table">
			<thead>
				<tr>
					<th>Name</th>
					<th>Slug</th>
					<th>Description</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
	`

	for _, s := range series {
		html += `
				<tr>
					<td><strong>` + s.Name + `</strong></td>
					<td><code>` + s.Slug + `</code></td>
					<td>` + s.Description + `</td>
					<td>` + s.CreatedAt.Format("Jan 2, 2006") + `</td>
					<td>
						<div class="table-actions">
							<button class="btn btn-sm btn-outline" hx-get="/admin/series/` + s.ID + `/edit" hx-target="#main-content">Edit</button>
							<button class="btn btn-sm btn-danger" hx-delete="/api/series/` + s.ID + `" hx-confirm="Delete this series?" hx-target="#main-content">Delete</button>
						</div>
					</td>
				</tr>
		`
	}

	html += `
			</tbody>
		</table>
	</div>
</div>
	`

	renderHTML(w, "text/html", html)
}

// HandlePostTypes serves the post types view
func HandlePostTypes(w http.ResponseWriter, r *http.Request, database *db.DB) {
	ctx := context.Background()

	types, err := database.GetPostTypes(ctx)
	if err != nil {
		renderHTML(w, "text/html", `<div class="alert alert-danger">Error loading post types</div>`)
		return
	}

	html := `
<div class="card">
	<div class="card-header">
		<h3>Post Types</h3>
	</div>
	<div class="card-body">
		<table class="table">
			<thead>
				<tr>
					<th>Name</th>
					<th>ID</th>
					<th>Description</th>
				</tr>
			</thead>
			<tbody>
	`

	for _, t := range types {
		html += `
				<tr>
					<td><strong>` + t.Name + `</strong></td>
					<td><code>` + t.ID + `</code></td>
					<td>` + t.Description + `</td>
				</tr>
		`
	}

	html += `
			</tbody>
		</table>
		<p style="margin-top: 16px; color: #666; font-size: 12px;">
			Post types are pre-defined and cannot be modified. These define the structure of different post categories.
		</p>
	</div>
</div>
	`

	renderHTML(w, "text/html", html)
}

// HandleExportPage serves the export configuration page
func HandleExportPage(w http.ResponseWriter, r *http.Request, database *db.DB) {
	ctx := context.Background()

	posts, total, err := database.ListPosts(ctx, &models.ListOptions{Limit: 1000})
	if err != nil {
		renderHTML(w, "text/html", `<div class="alert alert-danger">Error loading posts</div>`)
		return
	}

	html := `
<div class="card">
	<div class="card-header">
		<h3>Export & Deploy</h3>
	</div>
	<div class="card-body">
		<div style="margin-bottom: 24px;">
			<h4 style="font-size: 14px; font-weight: 600; margin-bottom: 12px;">Static Site Generation</h4>
			<p style="color: #666; font-size: 13px; margin-bottom: 16px;">
				Export your blog posts as Markdown files compatible with Hugo, Jekyll, or any static site generator.
			</p>
			<button class="btn btn-primary" hx-post="/api/exports/markdown" hx-confirm="Export all posts to Markdown? This will create markdown files for each post.">
				Export to Markdown
			</button>
		</div>

		<div style="padding-top: 16px; border-top: 1px solid #e0e0e0;">
			<h4 style="font-size: 14px; font-weight: 600; margin-bottom: 12px;">Export Status</h4>
			<table class="table" style="margin: 0;">
				<tr>
					<td><strong>Total Posts</strong></td>
					<td style="text-align: right;">` + strconv.Itoa(total) + `</td>
				</tr>
				<tr>
					<td><strong>Published Posts</strong></td>
					<td style="text-align: right;">` + func() string {
		count := 0
		for _, p := range posts {
			if p.Status == "published" {
				count++
			}
		}
		return strconv.Itoa(count)
	}() + `</td>
				</tr>
				<tr>
					<td><strong>Ready to Export</strong></td>
					<td style="text-align: right;"><span class="badge badge-success">` + strconv.Itoa(total) + ` posts</span></td>
				</tr>
			</table>
		</div>

		<div style="padding-top: 16px; border-top: 1px solid #e0e0e0; margin-top: 16px;">
			<h4 style="font-size: 14px; font-weight: 600; margin-bottom: 12px;">Hugo Integration</h4>
			<p style="color: #666; font-size: 13px; margin-bottom: 16px;">
				Exported markdown files are formatted for Hugo with proper front matter. Deploy to GitHub and set up automatic syncing.
			</p>
			<code style="background: #f9f9f9; padding: 12px; display: block; border-radius: 3px; font-size: 12px; color: #333; overflow-x: auto;">
---<br>
title: "Post Title"<br>
date: 2024-01-01<br>
draft: false<br>
slug: post-slug<br>
type: article<br>
---<br>
Post content in Markdown
			</code>
		</div>
	</div>
</div>
	`

	renderHTML(w, "text/html", html)
}

// renderHTML is a helper to render HTML content
func renderHTML(w http.ResponseWriter, contentType string, html string) {
	if contentType == "500" {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
	}
	fmt.Fprint(w, html)
}
