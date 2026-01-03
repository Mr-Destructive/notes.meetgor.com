# Editor Testing Checklist

## Issues Fixed

### 1. Link Post - Source URL Not Saving
**Problem**: Link posts weren't saving the `source_url` metadata.
**Fixed**: 
- Updated editor.js to properly handle type-specific fields
- `getMetadata()` now correctly collects fields and stores them in metadata
- Go backend accepts metadata as map[string]interface{} and marshals to JSON

### 2. [object Object] in Excerpt
**Problem**: Metadata was being stringified incorrectly, appearing as "[object Object]" in excerpt.
**Fixed**:
- Changed excerpt to use `trim()` to handle empty strings
- Fixed getMetadata() to only return objects, not stringified versions
- Go backend properly marshals metadata to JSON string

### 3. Title Not Mandatory for Some Post Types
**Problem**: Title was always required, but link/quote posts should have optional titles.
**Fixed**:
- Added `titleRequired` and `contentRequired` flags to POST_TEMPLATES
- `changePostType()` dynamically updates label from "Title *" to "Title (optional)"
- `savePost()` and `publishPost()` validate based on post type requirements
- Go backend accepts `*string` for Title/Content to handle nulls

### 4. Post Saved Instead of Draft
**Problem**: Status field might have been incorrect.
**Fixed**:
- Explicitly sets `status: 'draft'` in savePost()
- Explicitly sets `status: 'published'` in publishPost()
- Verified Go backend saves correct status value

### 5. Auto-Title for Link Posts
**Added**:
- `updateTitleFromUrl()` function extracts hostname/path from URL
- Called on source_url change event
- Only auto-fills if title is empty

## Test Cases

### Test 1: Create Article Post
1. Go to `/editor`
2. Select "article" type
3. Enter Title: "My First Article"
4. Enter Content: "This is my first article..."
5. Click "Save as Draft"
6. **Expected**: Post saved as draft, visible in dashboard

### Test 2: Create Link Post
1. Go to `/editor`  
2. Select "link" type
3. **Note**: Title field should show "(optional)"
4. **Note**: Content field should show "(optional)"
5. Enter Source URL: "https://example.com/interesting-article"
6. **Expected**: Title auto-filled with "interesting-article"
7. Add Tags: "reading", "learning"
8. Click "Save as Draft"
9. **Expected**: Post saved with metadata.source_url populated
10. Go to dashboard
11. Click draft post to edit
12. **Expected**: source_url field populated correctly, title and tags preserved

### Test 3: Create Quote Post
1. Go to `/editor`
2. Select "quote" type
3. **Note**: Title should be optional
4. **Note**: Content should be required
5. Enter Content: "Great quote text"
6. Enter Author (metadata): "John Doe"
7. Enter Source (metadata): "Article Title"
8. Click "Save as Draft"
9. **Expected**: Post saved, metadata contains author and source

### Test 4: Publish Article
1. Go to `/editor`
2. Select "article" type
3. Enter Title, Content
4. Click "Save as Draft" first
5. Modify content slightly
6. Click "Publish"
7. **Expected**: Post status changed to "published"
8. Go to dashboard
9. **Expected**: Post appears in "Published" filter

### Test 5: Edit Existing Post
1. Create and save a draft post
2. Go to dashboard
3. Click post to edit
4. **Expected**: All fields populated correctly
5. **Expected**: Metadata fields populated correctly
6. **Expected**: Tags restored
7. Modify content
8. Click "Save as Draft"
9. **Expected**: Draft updated without changing status

### Test 6: Empty Metadata Handling
1. Create a thought post (has no type-specific fields)
2. Save as draft
3. **Expected**: metadata is {} in database (not null or [object Object])
4. Edit the post
5. **Expected**: No errors, form loads correctly

## Database Verification

After each test, verify database:
```sql
SELECT id, type_id, title, status, metadata, tags FROM posts;
```

Look for:
- ✓ metadata is valid JSON (not "{}" when empty, not "[object Object]")
- ✓ tags is valid JSON array
- ✓ status is "draft" or "published"
- ✓ NULL fields for optional fields (title, content, excerpt)

## Metadata Examples

### Link Post Metadata
```json
{"source_url": "https://example.com"}
```

### Quote Post Metadata
```json
{"author": "John Doe", "source": "Some Article"}
```

### Tutorial Post Metadata
```json
{"difficulty": "intermediate", "estimated_time": "30 min", "tools": ["tool1", "tool2"]}
```

### List Post Metadata
```json
{"items": ["item1", "item2", "item3"], "list_type": "unordered"}
```

## UI Testing Notes

- Type-specific fields only appear after selecting post type
- "Save as Draft" button has 💾 emoji
- "Publish" button has 🚀 emoji
- "← Dashboard" button returns to dashboard
- Alert messages disappear after 5 seconds
- Form submission prevented (type="button" on buttons)
