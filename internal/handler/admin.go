package handler

import (
	"context"
	"html"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

// HandleSeriesEditor serves the series editor form
func HandleSeriesEditor(w http.ResponseWriter, r *http.Request, database *db.DB, seriesID string) {
	ctx := context.Background()

	var series *models.Series
	if seriesID != "" {
		var err error
		series, err = database.GetSeries(ctx, seriesID)
		if err != nil {
			renderHTML(w, "text/html", `<div class="alert alert-danger">Series not found</div>`)
			return
		}
	}

	title := "New Series"
	buttonText := "Create Series"
	method := "POST"
	endpoint := "/api/series"
	nameValue := ""
	slugValue := ""
	descValue := ""

	if series != nil {
		title = "Edit Series"
		buttonText = "Update Series"
		method = "PUT"
		endpoint = "/api/series/" + series.ID
		nameValue = series.Name
		slugValue = series.Slug
		descValue = series.Description
	}

	html := `
<div class="card">
	<div class="card-header">
		<div style="display: flex; justify-content: space-between; align-items: center;">
			<h3 id="editor-title">` + title + `</h3>
			<button class="btn btn-outline" hx-get="/admin/series" hx-target="#main-content">← Back to Series</button>
		</div>
	</div>
	<div class="card-body">
		<form id="series-form" onsubmit="handleSeriesSubmit(event)">
			<div class="form-group">
				<label for="series-name">Series Name *</label>
				<input type="text" id="series-name" name="name" placeholder="Enter series name" required value="` + nameValue + `" />
			</div>

			<div class="form-group">
				<label for="series-slug">Slug *</label>
				<input type="text" id="series-slug" name="slug" placeholder="series-slug" required value="` + slugValue + `" />
			</div>

			<div class="form-group">
				<label for="series-description">Description</label>
				<textarea id="series-description" name="description" placeholder="Optional description" rows="4">` + descValue + `</textarea>
			</div>

			<div style="display: flex; gap: 10px; margin-top: 24px;">
				<button type="button" class="btn btn-primary" onclick="saveSeries()">💾 ` + buttonText + `</button>
				<button type="button" class="btn btn-outline" hx-get="/admin/series" hx-target="#main-content">Cancel</button>
			</div>
		</form>

		<div id="series-message" style="margin-top: 16px; display: none;"></div>
	</div>
</div>

<script>
const SERIES_METHOD = '` + method + `';
const SERIES_ENDPOINT = '` + endpoint + `';

function saveSeries() {
	const name = document.getElementById('series-name').value;
	const slug = document.getElementById('series-slug').value;
	const description = document.getElementById('series-description').value;

	if (!name.trim()) {
		showMessage('Series name is required', 'danger');
		return;
	}

	if (!slug.trim()) {
		showMessage('Slug is required', 'danger');
		return;
	}

	const seriesData = {
		name: name,
		slug: slug,
		description: description
	};

	fetch(SERIES_ENDPOINT, {
		method: SERIES_METHOD,
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(seriesData)
	})
	.then(r => {
		if (!r.ok) return r.json().then(e => { throw new Error(e.error || 'Request failed'); });
		return r.json();
	})
	.then(data => {
		showMessage((SERIES_METHOD === 'POST' ? 'Created' : 'Updated') + ' series successfully!', 'success');
		setTimeout(() => {
			htmx.ajax('GET', '/admin/series', {target: '#main-content'});
		}, 1000);
	})
	.catch(err => {
		showMessage('Error: ' + err.message, 'danger');
	});
}

function showMessage(message, type) {
	const msg = document.getElementById('series-message');
	msg.className = 'alert alert-' + type;
	msg.textContent = message;
	msg.style.display = 'block';
	setTimeout(() => msg.style.display = 'none', 5000);
}
</script>
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

// HandlePostEditor serves the post editor form for new or existing posts
func HandlePostEditor(w http.ResponseWriter, r *http.Request, database *db.DB, postID string) {
	ctx := context.Background()

	types, err := database.GetPostTypes(ctx)
	if err != nil {
		renderHTML(w, "text/html", `<div class="alert alert-danger">Error loading post types</div>`)
		return
	}

	// Convert types to JavaScript array
	typesJSON := "["
	for i, t := range types {
		if i > 0 {
			typesJSON += ","
		}
		typesJSON += fmt.Sprintf(`{"id":"%s","name":"%s","description":"%s"}`, t.ID, t.Name, t.Description)
	}
	typesJSON += "]"

	// Fetch existing post if editing
	var post *models.Post
	if postID != "" {
		p, err := database.GetPost(ctx, postID)
		if err != nil {
			renderHTML(w, "text/html", `<div class="alert alert-danger">Post not found</div>`)
			return
		}
		post = p
	}

	htmlContent := `
<div class="card">
	<div class="card-header">
		<div style="display: flex; justify-content: space-between; align-items: center;">
			<h3 id="editor-title">New Post</h3>
			<button class="btn btn-outline" hx-get="/admin/posts" hx-target="#main-content">← Back to Posts</button>
		</div>
	</div>
	<div class="card-body">
		<div class="editor-tabs">
			<button type="button" class="editor-tab active" onclick="switchTab('edit')">✏️ Edit</button>
			<button type="button" class="editor-tab" onclick="switchTab('preview')">👁️ Preview</button>
		</div>
		
		<div id="edit-content" class="editor-content active">
			<form id="post-form" onsubmit="handlePostSubmit(event)">
				<div class="form-group">
				<label for="post-type">Post Type *</label>
				<select id="post-type" name="type_id" required onchange="updatePostType()">
					<option value="">Select a post type...</option>
	`

	for _, t := range types {
		selected := ""
		if post != nil && post.TypeID == t.ID {
			selected = "selected"
		}
		htmlContent += fmt.Sprintf(`<option value="%s" %s>%s</option>`, t.ID, selected, t.Name)
	}

	htmlContent += `
				</select>
			</div>

			<div class="form-group">
				<label for="post-title">Title <span id="title-required">*</span></label>
				<input type="text" id="post-title" name="title" placeholder="Enter post title"`
			if post != nil {
			htmlContent += fmt.Sprintf(` value="%s"`, html.EscapeString(post.Title))
			}
			htmlContent += ` />
			</div>

			<div class="form-group">
				<label for="post-slug">Slug *</label>
				<input type="text" id="post-slug" name="slug" placeholder="post-slug" required`
			if post != nil {
			htmlContent += fmt.Sprintf(` value="%s"`, html.EscapeString(post.Slug))
			}
			htmlContent += ` />
			</div>

			<div class="form-group">
				<label for="post-content">Content <span id="content-required">*</span></label>
				<textarea id="post-content" name="content" placeholder="Enter post content" rows="12">`
			if post != nil {
			htmlContent += html.EscapeString(post.Content)
			}
			htmlContent += `</textarea>
			</div>

			<div class="form-group">
				<label for="post-excerpt">Excerpt</label>
				<textarea id="post-excerpt" name="excerpt" placeholder="Optional excerpt" rows="3">`
			if post != nil {
			htmlContent += html.EscapeString(post.Excerpt)
			}
			htmlContent += `</textarea>
			</div>

			<div class="form-group">
				<label for="post-tags">Tags (comma-separated)</label>
				<input type="text" id="post-tags" name="tags" placeholder="tag1, tag2, tag3"`
			if post != nil && len(post.Tags) > 0 {
			htmlContent += ` value="`
			for i, tag := range post.Tags {
				if i > 0 {
					htmlContent += ", "
				}
				htmlContent += tag
			}
			htmlContent += `"`
			}
			htmlContent += ` />
			</div>

			<div id="type-specific-fields"></div>

			<div class="form-group">
				<label for="post-status">Status *</label>
				<select id="post-status" name="status" required>
					<option value="draft">Draft</option>
					<option value="published">Published</option>
					<option value="archived">Archived</option>
				</select>`
	if post != nil {
		htmlContent += fmt.Sprintf(`<script>document.getElementById('post-status').value = '%s';</script>`, escapeSingleQuote(post.Status))
	}
	htmlContent += `
			</div>

			<div style="display: flex; gap: 10px; margin-top: 24px;">
				<button type="button" class="btn btn-primary" onclick="saveAsDraft()">💾 Save as Draft</button>
				<button type="button" class="btn btn-success" onclick="publish()">🚀 Publish</button>
				<button type="button" class="btn btn-outline" hx-get="/admin/posts" hx-target="#main-content">Cancel</button>
			</div>
			</form>

			<div id="editor-message" style="margin-top: 16px; display: none;"></div>
		</div>
		
		<div id="preview-content" class="editor-content">
			<div class="preview" id="markdown-preview">
				<p style="color: #999; text-align: center;">Select content above to see preview</p>
			</div>
		</div>
	</div>
</div>

<style>
.editor-tabs {
	display: flex;
	gap: 0;
	border-bottom: 1px solid #ddd;
	margin-bottom: 20px;
}
.editor-tab {
	padding: 12px 20px;
	border: none;
	background: #f5f5f5;
	cursor: pointer;
	font-size: 14px;
	border-bottom: 3px solid transparent;
	transition: all 0.2s;
}
.editor-tab:hover {
	background: #e8e8e8;
}
.editor-tab.active {
	background: white;
	border-bottom-color: #0066cc;
	color: #0066cc;
	font-weight: 600;
}
.editor-content {
	display: none;
}
.editor-content.active {
	display: block;
}
.preview {
	padding: 20px;
	background: white;
	border: 1px solid #ddd;
	border-radius: 4px;
	min-height: 400px;
	font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
.preview h1, .preview h2, .preview h3, .preview h4, .preview h5, .preview h6 {
	margin: 20px 0 10px 0;
	font-weight: 600;
}
.preview p {
	line-height: 1.6;
	margin: 10px 0;
	color: #333;
}
.preview code {
	background: #f4f4f4;
	padding: 2px 6px;
	border-radius: 3px;
	font-family: 'Monaco', 'Courier New', monospace;
}
.preview pre {
	background: #f4f4f4;
	padding: 12px;
	border-radius: 4px;
	overflow-x: auto;
}
.preview blockquote {
	border-left: 4px solid #ddd;
	padding-left: 15px;
	margin-left: 0;
	color: #666;
}
.preview ul, .preview ol {
	margin: 10px 0;
	padding-left: 30px;
}
</style>

<script>
const POST_TYPES = ` + typesJSON + `;
const POST_TEMPLATES = {
	article: { titleRequired: true, contentRequired: true, fields: [] },
	link: { titleRequired: false, contentRequired: false, fields: [{ name: "source_url", label: "Source URL", type: "url" }] },
	quote: { titleRequired: false, contentRequired: true, fields: [{ name: "author", label: "Author", type: "text" }, { name: "source", label: "Source", type: "text" }] },
	tutorial: { titleRequired: true, contentRequired: true, fields: [{ name: "difficulty", label: "Difficulty", type: "select", options: ["beginner", "intermediate", "advanced"] }, { name: "estimated_time", label: "Estimated Time", type: "text" }] },
	list: { titleRequired: true, contentRequired: false, fields: [{ name: "list_type", label: "List Type", type: "select", options: ["ordered", "unordered"] }] },
	thought: { titleRequired: false, contentRequired: true, fields: [] },
	snippet: { titleRequired: true, contentRequired: true, fields: [{ name: "language", label: "Language", type: "text" }] },
	series: { titleRequired: true, contentRequired: true, fields: [] },
	review: { titleRequired: true, contentRequired: true, fields: [{ name: "rating", label: "Rating", type: "number" }, { name: "subject", label: "Subject", type: "text" }] },
	announcement: { titleRequired: true, contentRequired: true, fields: [] },
	photo: { titleRequired: false, contentRequired: false, fields: [{ name: "image_url", label: "Image URL", type: "url" }] },
	video: { titleRequired: false, contentRequired: false, fields: [{ name: "video_url", label: "Video URL", type: "url" }] }
};

function switchTab(tab) {
	// Update tabs
	document.querySelectorAll('.editor-tab').forEach(el => el.classList.remove('active'));
	event.target.classList.add('active');
	
	// Update content
	document.querySelectorAll('.editor-content').forEach(el => el.classList.remove('active'));
	document.getElementById(tab + '-content').classList.add('active');
	
	// Update preview if switching to it
	if (tab === 'preview') {
		updatePreview();
	}
}

function updatePreview() {
	const title = document.getElementById('post-title').value;
	const content = document.getElementById('post-content').value;
	
	if (!content) {
		document.getElementById('markdown-preview').innerHTML = 
			'<p style="color: #999; text-align: center;">No content to preview</p>';
		return;
	}
	
	// Check if marked is available
	if (typeof marked === 'undefined') {
		document.getElementById('markdown-preview').innerHTML = 
			'<p style="color: #999; text-align: center;">Loading markdown parser...</p>';
		return;
	}
	
	try {
		const html = marked.parse ? marked.parse(content) : marked(content);
		let preview = '';
		
		if (title) {
			preview += '<h1>' + escapeHtml(title) + '</h1>';
		}
		
		preview += html;
		document.getElementById('markdown-preview').innerHTML = preview;
	} catch (e) {
		console.error('Preview error:', e);
		document.getElementById('markdown-preview').innerHTML = 
			'<p style="color: #e74c3c;">Error rendering preview: ' + escapeHtml(e.message) + '</p>';
	}
}

function escapeHtml(text) {
	const map = {
		'&': '&amp;',
		'<': '&lt;',
		'>': '&gt;',
		'"': '&quot;',
		"'": '&#039;'
	};
	return text.replace(/[&<>"']/g, c => map[c]);
}

// Auto-update preview when content changes
document.addEventListener('DOMContentLoaded', () => {
	const contentArea = document.getElementById('post-content');
	if (contentArea) {
		contentArea.addEventListener('input', updatePreview);
		contentArea.addEventListener('change', updatePreview);
	}
	const titleArea = document.getElementById('post-title');
	if (titleArea) {
		titleArea.addEventListener('input', updatePreview);
	}
});

function updatePostType() {
	const typeId = document.getElementById('post-type').value;
	const template = POST_TEMPLATES[typeId] || {};
	
	// Update title requirement
	const titleRequired = document.getElementById('title-required');
	if (titleRequired) titleRequired.textContent = template.titleRequired ? '*' : '(optional)';
	
	// Update content requirement
	const contentRequired = document.getElementById('content-required');
	if (contentRequired) contentRequired.textContent = template.contentRequired ? '*' : '(optional)';
	
	// Generate type-specific fields
	const fieldsContainer = document.getElementById('type-specific-fields');
	fieldsContainer.innerHTML = '';
	
	if (template.fields && template.fields.length > 0) {
		const fieldset = document.createElement('fieldset');
		fieldset.style.marginTop = '20px';
		fieldset.style.paddingTop = '20px';
		fieldset.style.borderTop = '1px solid #ddd';
		
		const legend = document.createElement('legend');
		legend.textContent = 'Type-Specific Fields';
		legend.style.fontSize = '14px';
		legend.style.fontWeight = '600';
		legend.style.marginBottom = '12px';
		fieldset.appendChild(legend);
		
		for (const field of template.fields) {
			const group = document.createElement('div');
			group.className = 'form-group';
			
			const label = document.createElement('label');
			label.setAttribute('for', 'field-' + field.name);
			label.textContent = field.label;
			group.appendChild(label);
			
			let input;
			if (field.type === 'select') {
				input = document.createElement('select');
				input.id = 'field-' + field.name;
				input.name = 'meta_' + field.name;
				for (const opt of (field.options || [])) {
					const option = document.createElement('option');
					option.value = opt;
					option.textContent = opt;
					input.appendChild(option);
				}
			} else if (field.type === 'textarea') {
				input = document.createElement('textarea');
				input.id = 'field-' + field.name;
				input.name = 'meta_' + field.name;
				input.rows = 3;
			} else {
				input = document.createElement('input');
				input.type = field.type;
				input.id = 'field-' + field.name;
				input.name = 'meta_' + field.name;
				input.placeholder = 'Enter ' + field.label.toLowerCase();
			}
			
			group.appendChild(input);
			fieldset.appendChild(group);
		}
		
		fieldsContainer.appendChild(fieldset);
	}
}

function saveAsDraft() {
	document.getElementById('post-status').value = 'draft';
	handlePostSubmit();
}

function publish() {
	document.getElementById('post-status').value = 'published';
	handlePostSubmit();
}

function showMessage(message, type) {
	const msg = document.getElementById('editor-message');
	msg.className = 'alert alert-' + type;
	msg.textContent = message;
	msg.style.display = 'block';
	setTimeout(() => msg.style.display = 'none', 5000);
}

function handlePostSubmit(event) {
	if (event) event.preventDefault();
	
	const typeId = document.getElementById('post-type').value;
	if (!typeId) {
		showMessage('Please select a post type', 'danger');
		return;
	}
	
	const template = POST_TEMPLATES[typeId] || {};
	
	const title = document.getElementById('post-title').value;
	if (template.titleRequired && !title.trim()) {
		showMessage('Title is required for this post type', 'danger');
		return;
	}
	
	const content = document.getElementById('post-content').value;
	if (template.contentRequired && !content.trim()) {
		showMessage('Content is required for this post type', 'danger');
		return;
	}
	
	const slug = document.getElementById('post-slug').value;
	if (!slug.trim()) {
		showMessage('Slug is required', 'danger');
		return;
	}
	
	// Build metadata from type-specific fields
	const metadata = {};
	const inputs = document.querySelectorAll('[name^="meta_"]');
	for (const input of inputs) {
		const key = input.name.substring(5);
		if (input.value) {
			metadata[key] = input.value;
		}
	}
	
	// Parse tags
	const tagsStr = document.getElementById('post-tags').value;
	const tags = tagsStr ? tagsStr.split(',').map(t => t.trim()).filter(t => t) : [];
	
	const postData = {
		type_id: typeId,
		title: title || null,
		slug: slug,
		content: content || null,
		excerpt: document.getElementById('post-excerpt').value || '',
		tags: tags,
		metadata: Object.keys(metadata).length > 0 ? metadata : {},
		status: document.getElementById('post-status').value
	};
	
	const postId = '` + postID + `';
	const method = postId ? 'PUT' : 'POST';
	const url = postId ? '/api/posts/' + postId : '/api/posts';
	
	fetch(url, {
		method: method,
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(postData)
	})
	.then(r => {
		if (!r.ok) return r.json().then(e => { throw new Error(e.error || 'Request failed'); });
		return r.json();
	})
	.then(data => {
		showMessage((postId ? 'Updated' : 'Created') + ' post successfully!', 'success');
		setTimeout(() => {
			htmx.ajax('GET', '/admin/posts', {target: '#main-content'});
		}, 1000);
	})
	.catch(err => {
		showMessage('Error: ' + err.message, 'danger');
	});
}

// Initialize on load
window.addEventListener('load', () => {
	`
	if post != nil {
		htmlContent += fmt.Sprintf(`document.getElementById('post-type').value = '%s';`, escapeSingleQuote(post.TypeID))
	}
	htmlContent += `
	updatePostType();
});
</script>
	`

	renderHTML(w, "text/html", htmlContent)
}

// escapeSingleQuote escapes single quotes for JavaScript strings
func escapeSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
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
