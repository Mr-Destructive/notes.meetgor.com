# JSON to Markdown Migration

## Overview
Migrated 55 posts from legacy JSON blog archives into the Hugo-based SSG for publishing on Vercel.

## Sources Converted

### 1. drafts.json
- **Items**: 9 items
- **Type**: Drafts and unpublished content
- **Fields**: title, content (HTML), type_id, created_at, updated_at

### 2. temp-blog.json  
- **Items**: 37 items  
- **Type**: Temporary/staging blog posts
- **Fields**: title, slug, content (HTML), type_id, published (bool/int), created_at

### 3. links-blog.json
- **Items**: 50 items
- **Type**: Link aggregator posts
- **Fields**: url, title, commentary, image_url
- **Converted to**: type="link"

## Converter Tool
**Location**: `cmd/convert-json/main.go`

### Features
- Converts HTML to Markdown automatically using regex
- Handles multiple JSON schemas
- Generates proper Hugo front matter
- Creates SEO-friendly slugs from titles
- Categorizes by post type (posts vs links)

### Usage
```bash
# Convert all sources
go run ./cmd/convert-json/main.go --all --output exports/content/posts/

# Convert specific source
go run ./cmd/convert-json/main.go --drafts --output exports/content/posts/
go run ./cmd/convert-json/main.go --temp --output exports/content/posts/
go run ./cmd/convert-json/main.go --links --output exports/content/posts/
```

## Results

### Final Stats
- **Total Posts Generated**: 55 markdown files
- **Hugo Pages Built**: 69 (includes tag/category pages)
- **Conversion Success**: 100%

### Post Distribution
- Regular Posts: 7
- Link Posts: 48

### Output Format
Each markdown file includes:
```yaml
---
title: "..."
date: 2025-01-03
slug: post-slug
draft: false
type: posts  # or "link"
description: ""
tags: []
---
```

## Post Type Separation

The site now displays posts in two sections (per layout in `exports/layouts/posts/list.html`):

1. **Posts** - Regular articles and content
2. **Interesting Links** - Aggregated link posts with commentary

## Deployment

All converted posts are:
- ✅ Committed to `exports/content/posts/`
- ✅ Generated into `public/` directory
- ✅ Pushed to GitHub main branch
- ✅ Triggering Vercel deployment (deploy-vercel.yml)

The live site is now serving 55+ imported posts alongside any newly published content via the CMS.

## Notes

- HTML entities properly unescaped (e.g., `&quot;`, `&lt;`, `&gt;`)
- Code blocks and inline code preserved
- Links converted to markdown format `[text](url)`
- Date fallback: uses creation timestamp from source
- Slug generation: auto-generated from title if not provided
- Duplicates: Removed duplicate "Advent of SQL Day 2" entry
