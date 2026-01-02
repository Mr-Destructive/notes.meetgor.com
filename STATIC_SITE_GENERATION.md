# Static Site Generation & Deployment Guide

## Overview

The blog system includes integrated static site generation (SSG) for exporting content to Hugo, Jekyll, or any Markdown-based static site generator. Export markdown files with proper front matter, Hugo configuration, and GitHub Actions deployment workflow - all generated automatically.

## Features

### Export Capabilities

1. **Markdown Export**
   - All published posts converted to Markdown
   - Hugo-compatible YAML front matter
   - Proper metadata preservation (tags, dates, types)
   - Markdown content exactly as authored

2. **Hugo Configuration**
   - Auto-generated `hugo.toml` with sensible defaults
   - Content structure: `content/posts/`
   - Customizable baseURL and site metadata

3. **Deployment Workflow**
   - GitHub Actions CI/CD workflow auto-generated
   - Builds Hugo site on push to main
   - Deploys to GitHub Pages automatically
   - No manual deployment steps needed

## File Structure

Generated export produces:

```
exports/
├── content/
│   └── posts/
│       ├── post-one.md
│       ├── post-two.md
│       └── ...
├── hugo.toml          (Site configuration)
└── .github/
    └── workflows/
        └── deploy.yml (GitHub Actions workflow)
```

## Generated Front Matter Format

```yaml
---
title: "Post Title"
date: 2024-01-15
slug: post-slug
draft: false
type: article
description: "Post excerpt"
tags: ["tag1", "tag2"]
metadata:
  custom_field: "value"
---

Post content in Markdown...
```

## Export API

### Endpoint: POST /api/exports

Triggers export of all published posts to Markdown format with Hugo configuration.

**Request:**
```bash
curl -X POST http://localhost:8080/api/exports \
  -H "Content-Type: application/json"
```

**Response:**
```json
{
  "exported_at": "2026-01-02T21:22:22Z",
  "posts_count": 5,
  "files_count": 5,
  "output_dir": "./exports",
  "success": true,
  "message": "Successfully exported 5 posts to ./exports"
}
```

### Endpoint: GET /api/exports

Returns JSON export of published posts.

**Request:**
```bash
curl http://localhost:8080/api/exports
```

**Response:**
```json
[
  {
    "id": "post-id",
    "title": "Post Title",
    "slug": "post-slug",
    "content": "...",
    "type_id": "article",
    "status": "published",
    "tags": ["tag1", "tag2"],
    "created_at": "2024-01-15T10:00:00Z",
    ...
  },
  ...
]
```

## Export Interface

### Admin Panel Export Page

Navigate to **Export** in the admin sidebar to:

- View export statistics (total posts, published count)
- View Hugo front matter format example
- Trigger markdown export with confirmation
- See export status and results

### Export Dialog

```
Export & Deploy
  ├── Static Site Generation
  │   └── [Export to Markdown] button
  ├── Export Status
  │   ├── Total Posts: 10
  │   ├── Published Posts: 8
  │   └── Ready to Export: 10 posts
  └── Hugo Integration
      └── Front matter format reference
```

## Deployment Workflow

### GitHub Actions Workflow (Auto-generated)

File: `.github/workflows/deploy.yml`

Triggers on:
- Push to `main` branch
- Manual workflow dispatch

Steps:
1. Checkout code
2. Setup Hugo (latest version)
3. Build site (`hugo --minify`)
4. Deploy to GitHub Pages

No additional configuration needed - just push the exported repository to GitHub.

## Setup Instructions

### 1. Export Posts

From admin panel or API:
```bash
curl -X POST http://localhost:8080/api/exports
```

This creates `./exports/` directory with all content.

### 2. Create GitHub Repository

```bash
cd exports
git init
git add .
git commit -m "Initial blog export"
git remote add origin https://github.com/username/blog.git
git push -u origin main
```

### 3. Enable GitHub Pages

In GitHub repository settings:
- Go to Settings > Pages
- Set source to "GitHub Actions"
- Saves automatically once workflow runs

### 4. Configure Hugo Theme

The generated config uses `theme: blog-theme`. Choose a Hugo theme:

```bash
cd exports
git submodule add https://github.com/user/theme themes/blog-theme
```

Popular minimal themes:
- [Ananke](https://github.com/theNewDynamic/gohugo-theme-ananke)
- [Hugo Book](https://github.com/alex-shpak/hugo-book)
- [Hugo PaperMod](https://github.com/adityatelange/hugo-PaperMod)
- [Minimal](https://github.com/calintat/minimal)

### 5. Customize Configuration

Edit `exports/hugo.toml`:

```toml
baseURL = "https://your-domain.com/"
title = "Your Blog Title"

[params]
  author = "Your Name"
  description = "Blog description"
```

## Hugo Theme Integration

### Creating a Minimal Theme

Basic theme structure for `themes/blog-theme/`:

```
themes/blog-theme/
├── layouts/
│   ├── _default/
│   │   ├── baseof.html
│   │   ├── list.html
│   │   └── single.html
│   └── partials/
│       ├── header.html
│       └── footer.html
├── static/
│   └── css/
│       └── style.css
└── theme.toml
```

See [Hugo Theme Documentation](https://gohugo.io/themes/creating/) for details.

## Example Workflow

### Daily Export to Static Site

Add to your workflow:

```bash
#!/bin/bash

# Export latest posts
curl -X POST http://localhost:8080/api/exports

# Commit and push to GitHub
cd exports
git add -A
git commit -m "Blog update: $(date)"
git push origin main
```

Run via cron:
```bash
0 0 * * * /path/to/export-and-deploy.sh
```

## Limitations & Notes

- Only **published** posts are exported (draft/archived excluded)
- Metadata is preserved as YAML front matter
- Images in posts should use relative paths (manage separately)
- Tags are exported as YAML array
- Custom metadata fields preserved in front matter

## Advanced: Custom Metadata

Posts can include custom metadata:

```json
{
  "title": "Advanced Topic",
  "type_id": "tutorial",
  "metadata": {
    "difficulty": "advanced",
    "duration": "30 minutes",
    "prerequisites": "Go basics"
  }
}
```

Exports as:

```yaml
---
title: "Advanced Topic"
type: tutorial
metadata:
  difficulty: "advanced"
  duration: "30 minutes"
  prerequisites: "Go basics"
---
```

## Troubleshooting

### Posts Not Exporting

- Verify posts have `status: "published"`
- Check database has posts: `GET /api/posts?status=published`
- Ensure markdown export is working: `POST /api/exports`

### GitHub Pages Not Deploying

1. Check workflow ran: GitHub Actions tab in repo
2. Enable Pages in Settings > Pages > source: "GitHub Actions"
3. Verify theme exists in repository
4. Check `hugo.toml` baseURL is correct

### Theme Not Rendering

1. Verify theme folder exists in `themes/`
2. Theme name in config matches folder name
3. Run `hugo serve` locally to debug
4. Check Hugo version compatibility

## Summary

The integrated SSG system makes publishing to static hosting simple:
1. Author in the admin panel
2. Export with one click/API call
3. Push to GitHub
4. Automatic deployment to GitHub Pages

No build tools, no complex setup - just simple Markdown export with Hugo integration.
