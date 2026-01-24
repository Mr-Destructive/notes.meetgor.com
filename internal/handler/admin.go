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
<style>
.stat-card { background: white; padding: 20px; border: 1px solid #e5e5e5; border-radius: 8px; }
.stat-card > div:first-child { font-size: 28px; font-weight: 700; margin-bottom: 8px; color: #0a0a0a; }
.stat-card > div:last-child { color: #666; font-size: 13px; font-weight: 500; }
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 16px; margin-bottom: 24px; }
.card { background: white; border-radius: 8px; margin-bottom: 24px; border: 1px solid #e5e5e5; }
.card-header { padding: 20px; border-bottom: 1px solid #e5e5e5; display: flex; justify-content: space-between; align-items: center; }
.card-header h3 { font-size: 16px; font-weight: 600; margin: 0; }
.card-body { padding: 20px; }
.table { width: 100%; border-collapse: collapse; }
.table th { text-align: left; padding: 12px 8px; background: #fafafa; border-bottom: 1px solid #e5e5e5; font-weight: 600; font-size: 13px; color: #666; }
.table td { padding: 12px 8px; border-bottom: 1px solid #e5e5e5; }
.table tr:last-child td { border-bottom: none; }
.table tr:hover { background: #f9f9f9; }
.btn { padding: 8px 16px; border: 1px solid #e5e5e5; border-radius: 6px; background: white; cursor: pointer; transition: all 0.15s; text-decoration: none; font-size: 13px; font-weight: 500; }
.btn:hover { background: #f9f9f9; border-color: #ccc; }
.btn-primary { background: #0a0a0a; color: white; border: none; }
.btn-primary:hover { background: #1a1a1a; }
.btn-danger { background: #e63946; color: white; border: none; }
.btn-danger:hover { background: #d62828; }
.btn-outline { background: white; color: #0a0a0a; border: 1px solid #e5e5e5; }
.btn-outline:hover { background: #f9f9f9; }
.btn-sm { padding: 6px 12px; font-size: 12px; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 12px; font-size: 12px; font-weight: 600; }
.badge-success { background: #d1e7dd; color: #0f5132; }
.badge-warning { background: #fff3cd; color: #664d03; }
.badge-danger { background: #f8d7da; color: #842029; }
.table-actions { display: flex; gap: 6px; }
</style>
<div style="margin-bottom: 24px;">
	<h2 style="font-size: 18px; font-weight: 600; margin-bottom: 16px;">Overview</h2>
	<div class="stats-grid">
		<div class="stat-card">
			<div>` + strconv.Itoa(len(posts)) + `</div>
			<div>Total Posts</div>
		</div>
		<div class="stat-card">
			<div>` + strconv.Itoa(published) + `</div>
			<div>Published</div>
		</div>
		<div class="stat-card">
			<div>` + strconv.Itoa(draft) + `</div>
			<div>Draft</div>
		</div>
		<div class="stat-card">
			<div>` + strconv.Itoa(len(series)) + `</div>
			<div>Series</div>
		</div>
		<div class="stat-card">
			<div>` + strconv.Itoa(len(types)) + `</div>
			<div>Post Types</div>
		</div>
		<div class="stat-card">
			<div>` + strconv.Itoa(len(tags)) + `</div>
			<div>Tags</div>
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
<style>
.card { background: white; border-radius: 8px; margin-bottom: 24px; border: 1px solid #e5e5e5; }
.card-header { padding: 20px; border-bottom: 1px solid #e5e5e5; display: flex; justify-content: space-between; align-items: center; }
.card-header h3 { font-size: 16px; font-weight: 600; margin: 0; }
.card-body { padding: 20px; }
.table { width: 100%; border-collapse: collapse; }
.table th { text-align: left; padding: 12px 8px; background: #fafafa; border-bottom: 1px solid #e5e5e5; font-weight: 600; font-size: 13px; color: #666; }
.table td { padding: 12px 8px; border-bottom: 1px solid #e5e5e5; }
.table tr:last-child td { border-bottom: none; }
.table tr:hover { background: #f9f9f9; }
.btn { padding: 8px 16px; border: 1px solid #e5e5e5; border-radius: 6px; background: white; cursor: pointer; transition: all 0.15s; text-decoration: none; font-size: 13px; font-weight: 500; }
.btn:hover { background: #f9f9f9; border-color: #ccc; }
.btn-primary { background: #0a0a0a; color: white; border: none; }
.btn-primary:hover { background: #1a1a1a; }
.btn-danger { background: #e63946; color: white; border: none; }
.btn-danger:hover { background: #d62828; }
.btn-outline { background: white; color: #0a0a0a; border: 1px solid #e5e5e5; }
.btn-outline:hover { background: #f9f9f9; }
.btn-sm { padding: 6px 12px; font-size: 12px; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 12px; font-size: 12px; font-weight: 600; }
.badge-success { background: #d1e7dd; color: #0f5132; }
.badge-warning { background: #fff3cd; color: #664d03; }
.badge-danger { background: #f8d7da; color: #842029; }
.search-bar { display: flex; gap: 8px; margin-bottom: 20px; }
.search-bar select { padding: 10px 12px; border: 1px solid #e5e5e5; border-radius: 6px; font-size: 14px; }
.table-actions { display: flex; gap: 6px; }
</style>
<div class="card">
	<div class="card-header">
		<div style="display: flex; justify-content: space-between; align-items: center; width: 100%; gap: 20px;">
			<h3>Posts</h3>
			<button class="btn btn-primary" hx-get="/admin/posts/new" hx-target="#main-content">New Post</button>
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
				<td><strong>` + post.Title + `</strong></td>
				<td>` + post.TypeID + `</td>
				<td>` + statusBadge + `</td>
				<td>` + post.CreatedAt.Format("Jan 2, 2006") + `</td>
				<td>
					<div class="table-actions">
						<button class="btn btn-sm btn-outline" hx-get="/admin/posts/` + post.ID + `/edit" hx-target="#main-content">Edit</button>
						<button class="btn btn-sm btn-danger" hx-delete="/api/posts/` + post.ID + `" hx-confirm="Delete?">Delete</button>
					</div>
				</td>
			</tr>
		`
	}

	html += `
			</tbody>
		</table>
		<div style="margin-top: 20px; display: flex; gap: 10px; align-items: center;">
			<small style="color: #666;">Showing ` + strconv.Itoa(offset+1) + ` to ` + strconv.Itoa(offset+len(posts)) + ` of ` + strconv.Itoa(total) + `</small>
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
		renderHTML(w, "text/html", `<div>Error loading series</div>`)
		return
	}

	html := `
<style>
.card { background: white; border-radius: 8px; margin-bottom: 24px; border: 1px solid #e5e5e5; }
.card-header { padding: 20px; border-bottom: 1px solid #e5e5e5; display: flex; justify-content: space-between; align-items: center; }
.card-header h3 { font-size: 16px; font-weight: 600; margin: 0; }
.card-body { padding: 20px; }
.table { width: 100%; border-collapse: collapse; }
.table th { text-align: left; padding: 12px 8px; background: #fafafa; border-bottom: 1px solid #e5e5e5; font-weight: 600; font-size: 13px; color: #666; }
.table td { padding: 12px 8px; border-bottom: 1px solid #e5e5e5; }
.table tr:last-child td { border-bottom: none; }
.table tr:hover { background: #f9f9f9; }
.btn { padding: 8px 16px; border: 1px solid #e5e5e5; border-radius: 6px; background: white; cursor: pointer; transition: all 0.15s; text-decoration: none; font-size: 13px; font-weight: 500; }
.btn:hover { background: #f9f9f9; border-color: #ccc; }
.btn-primary { background: #0a0a0a; color: white; border: none; }
.btn-primary:hover { background: #1a1a1a; }
.btn-danger { background: #e63946; color: white; border: none; }
.btn-danger:hover { background: #d62828; }
.btn-sm { padding: 6px 12px; font-size: 12px; }
.table-actions { display: flex; gap: 6px; }
</style>
<div class="card">
	<div class="card-header">
		<div style="display: flex; justify-content: space-between; align-items: center; width: 100%; gap: 20px;">
			<h3>Series</h3>
			<button class="btn btn-primary" hx-get="/admin/series/new" hx-target="#main-content">New Series</button>
		</div>
	</div>
	<div class="card-body">
		<table class="table">
			<thead>
				<tr>
					<th>Name</th>
					<th>Slug</th>
					<th>Posts</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
	`

	for _, s := range series {
		html += `
			<tr>
				<td><strong>` + s.Name + `</strong></td>
				<td>` + s.Slug + `</td>
				<td>-</td>
				<td>
					<div class="table-actions">
						<button class="btn btn-sm btn-outline" hx-get="/admin/series/` + s.ID + `/edit" hx-target="#main-content">Edit</button>
						<button class="btn btn-sm btn-danger" hx-delete="/api/series/` + s.ID + `" hx-confirm="Delete?">Delete</button>
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
		renderHTML(w, "text/html", `<div>Error loading post types</div>`)
		return
	}

	html := `
<style>
.card { background: white; border-radius: 8px; margin-bottom: 24px; border: 1px solid #e5e5e5; }
.card-header { padding: 20px; border-bottom: 1px solid #e5e5e5; }
.card-header h3 { font-size: 16px; font-weight: 600; margin: 0; }
.card-body { padding: 20px; }
.table { width: 100%; border-collapse: collapse; }
.table th { text-align: left; padding: 12px 8px; background: #fafafa; border-bottom: 1px solid #e5e5e5; font-weight: 600; font-size: 13px; color: #666; }
.table td { padding: 12px 8px; border-bottom: 1px solid #e5e5e5; }
.table tr:last-child td { border-bottom: none; }
.table tr:hover { background: #f9f9f9; }
</style>
<div class="card">
	<div class="card-header">
		<h3>Post Types (` + strconv.Itoa(len(types)) + `)</h3>
	</div>
	<div class="card-body">
		<table class="table">
			<thead>
				<tr>
					<th>ID</th>
					<th>Name</th>
					<th>Description</th>
				</tr>
			</thead>
			<tbody>
	`

	for _, t := range types {
		html += `
			<tr>
				<td><code>` + t.ID + `</code></td>
				<td>` + t.Name + `</td>
				<td><small>` + t.Description + `</small></td>
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

// HandlePostEditor serves the post editor form
func HandlePostEditor(w http.ResponseWriter, r *http.Request, database *db.DB, postID string) {
	ctx := context.Background()

	var post *models.Post
	var err error

	if postID != "" && postID != "new" {
		post, err = database.GetPost(ctx, postID)
		if err != nil {
			renderHTML(w, "text/html", `<div style="color: red;">Post not found</div>`)
			return
		}
	}

	postTypes, _ := database.GetPostTypes(ctx)

	htmlContent := `
<style>
.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 6px; font-weight: 500; font-size: 14px; }
.form-group input, .form-group textarea, .form-group select { width: 100%; padding: 10px 12px; border: 1px solid #e5e5e5; border-radius: 6px; font-family: inherit; font-size: 14px; }
.form-group input:focus, .form-group textarea:focus, .form-group select:focus { outline: none; border-color: #0a0a0a; box-shadow: 0 0 0 3px rgba(10, 10, 10, 0.05); }
.form-group textarea { min-height: 200px; resize: vertical; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.form-actions { display: flex; gap: 10px; margin-top: 30px; }
.btn { padding: 8px 16px; border: 1px solid #e5e5e5; border-radius: 6px; background: white; cursor: pointer; transition: all 0.15s; text-decoration: none; font-size: 13px; font-weight: 500; }
.btn:hover { background: #f9f9f9; border-color: #ccc; }
.btn-primary { background: #0a0a0a; color: white; border: none; }
.btn-primary:hover { background: #1a1a1a; }
.btn-secondary { background: white; }
.alert { padding: 12px 16px; border-radius: 8px; margin-bottom: 20px; font-size: 14px; }
.alert-success { background: #d1e7dd; color: #0f5132; }
.alert-danger { background: #f8d7da; color: #842029; }
#markdown-preview { padding: 20px; background: #f9f9f9; border-radius: 6px; border: 1px solid #e5e5e5; }
</style>

<div id="editor-message" class="alert" style="display: none;"></div>

<form id="post-form" style="background: white; border-radius: 8px; border: 1px solid #e5e5e5; padding: 20px;">
	<div class="form-group">
		<label for="post-type">Post Type *</label>
		<select id="post-type" name="type_id" onchange="updatePostType()" required>
			<option value="">-- Select a type --</option>
	`

	for _, t := range postTypes {
		selected := ""
		if post != nil && post.TypeID == t.ID {
			selected = "selected"
		}
		htmlContent += `<option value="` + t.ID + `" ` + selected + `>` + t.Name + `</option>`
	}

	htmlContent += `
		</select>
	</div>

	<div class="form-row">
		<div class="form-group">
			<label for="post-title">Title <span id="title-required"></span></label>
			<input type="text" id="post-title" name="title" value="`
	if post != nil {
		htmlContent += html.EscapeString(post.Title)
	}
	htmlContent += `">
		</div>
		<div class="form-group">
			<label for="post-slug">Slug *</label>
			<input type="text" id="post-slug" name="slug" required value="`
	if post != nil {
		htmlContent += html.EscapeString(post.Slug)
	}
	htmlContent += `">
		</div>
	</div>

	<div class="form-group">
		<label for="post-content">Content <span id="content-required"></span></label>
		<textarea id="post-content" name="content">`
	if post != nil {
		htmlContent += html.EscapeString(post.Content)
	}
	htmlContent += `</textarea>
	</div>

	<div class="form-group">
		<label for="post-excerpt">Excerpt</label>
		<textarea id="post-excerpt" name="excerpt" style="min-height: 80px;">`
	if post != nil {
		htmlContent += html.EscapeString(post.Excerpt)
	}
	htmlContent += `</textarea>
	</div>

	<div class="form-group">
		<label for="post-tags">Tags (comma-separated)</label>
		<input type="text" id="post-tags" name="tags" value="`
	if post != nil && len(post.Tags) > 0 {
		htmlContent += html.EscapeString(strings.Join(post.Tags, ", "))
	}
	htmlContent += `">
	</div>

	<div class="form-group">
		<label for="post-status">Status</label>
		<select id="post-status" name="status">
			<option value="draft" `
	if post == nil || post.Status == "draft" {
		htmlContent += "selected"
	}
	htmlContent += `>Draft</option>
			<option value="published" `
	if post != nil && post.Status == "published" {
		htmlContent += "selected"
	}
	htmlContent += `>Published</option>
			<option value="archived" `
	if post != nil && post.Status == "archived" {
		htmlContent += "selected"
	}
	htmlContent += `>Archived</option>
		</select>
	</div>

	<div id="type-specific-fields"></div>

	<div class="form-actions">
		<button type="button" class="btn btn-primary" onclick="publish()">Publish</button>
		<button type="button" class="btn btn-secondary" onclick="saveAsDraft()">Save as Draft</button>
		<button type="button" class="btn btn-secondary" onclick="htmx.ajax('GET', '/admin/posts', {target: '#main-content'})">Cancel</button>
	</div>
</form>

<div style="margin-top: 40px; display: grid; grid-template-columns: 1fr 1fr; gap: 20px;">
	<div>
		<h3 style="font-size: 14px; font-weight: 600; margin-bottom: 12px;">Content Preview</h3>
		<div id="markdown-preview" style="max-height: 500px; overflow-y: auto;"></div>
	</div>
</div>

<script>
const POST_TEMPLATES = {
	'article': { titleRequired: true, contentRequired: true, fields: [] },
	'link': { titleRequired: false, contentRequired: false, fields: [{name: 'url', label: 'URL', type: 'url'}] },
	'page': { titleRequired: true, contentRequired: true, fields: [] },
	'post': { titleRequired: true, contentRequired: true, fields: [] },
	'project': { titleRequired: true, contentRequired: false, fields: [] },
	'newsletter': { titleRequired: true, contentRequired: true, fields: [{name: 'edition', label: 'Edition', type: 'number'}] },
};

function updatePreview() {
	const title = document.getElementById('post-title').value;
	const content = document.getElementById('post-content').value;
	
	if (!content && !title) {
		document.getElementById('markdown-preview').innerHTML = 
			'<p style="color: #999; text-align: center;">No content to preview</p>';
		return;
	}
	
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
			'<p style="color: #e63946;">Error rendering preview: ' + escapeHtml(e.message) + '</p>';
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
	
	const titleRequired = document.getElementById('title-required');
	if (titleRequired) titleRequired.textContent = template.titleRequired ? '*' : '(optional)';
	
	const contentRequired = document.getElementById('content-required');
	if (contentRequired) contentRequired.textContent = template.contentRequired ? '*' : '(optional)';
	
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
	
	const metadata = {};
	const inputs = document.querySelectorAll('[name^="meta_"]');
	for (const input of inputs) {
		const key = input.name.substring(5);
		if (input.value) {
			metadata[key] = input.value;
		}
	}
	
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

window.addEventListener('load', () => {
	` + (func() string {
		if post != nil {
			return fmt.Sprintf(`document.getElementById('post-type').value = '%s';`, escapeSingleQuote(post.TypeID))
		}
		return ""
	})() + `
	updatePostType();
});
</script>
	`

	renderHTML(w, "text/html", htmlContent)
}

// HandleExportPage serves the export interface
func HandleExportPage(w http.ResponseWriter, r *http.Request, database *db.DB) {
	html := `
<style>
.card { background: white; border-radius: 8px; margin-bottom: 24px; border: 1px solid #e5e5e5; }
.card-header { padding: 20px; border-bottom: 1px solid #e5e5e5; }
.card-header h3 { font-size: 16px; font-weight: 600; margin: 0; }
.card-body { padding: 20px; }
.btn { padding: 8px 16px; border: 1px solid #e5e5e5; border-radius: 6px; background: white; cursor: pointer; transition: all 0.15s; text-decoration: none; font-size: 13px; font-weight: 500; }
.btn-primary { background: #0a0a0a; color: white; border: none; }
.btn-primary:hover { background: #1a1a1a; }
</style>
<div class="card">
	<div class="card-header">
		<h3>Export Posts</h3>
	</div>
	<div class="card-body">
		<p style="margin-bottom: 20px; color: #666;">Export all posts in Markdown format for static site generation.</p>
		<a href="/api/exports" class="btn btn-primary" download="posts.zip">Download Export</a>
	</div>
</div>
	`
	renderHTML(w, "text/html", html)
}

func escapeSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

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
