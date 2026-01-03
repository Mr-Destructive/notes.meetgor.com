# Hybrid Post Management Strategy

## Architecture

The blog now uses a hybrid approach for managing posts:

### 1. **Git-Committed Posts** (Historical/Imported)
- **Source**: `/exports/content/posts/` (55 posts from JSON migration)
- **Management**: Manually committed to Git
- **Type**: Mix of posts and link posts
- **Purpose**: Archive and searchable history

### 2. **Database Posts** (Dynamic/Current)
- **Source**: Turso database (via CMS)
- **Management**: Automatic sync workflow (`sync-posts-turso.yml`)
- **Type**: Posts created via admin interface
- **Purpose**: New, live content

## Workflow

### Current Flow

```
┌─────────────────────────────────────────────────┐
│         Git-Committed Posts                     │
│    (Imported from JSON archives)                │
│    - 55 historical posts                        │
│    - Persistent in version control              │
│    - Searchable and indexed                     │
└──────────────┬──────────────────────────────────┘
               │
               ├─→ exports/content/posts/
               │   ├── post1.md
               │   ├── post2.md
               │   └── ... (55 files)
               │
               └─→ Hugo Build
                   └─→ 69 Pages Generated
                       ├── Homepage with all posts
                       ├── Posts section (regular posts)
                       ├── Links section (link posts)
                       └── Tag pages

┌─────────────────────────────────────────────────┐
│     Database Posts (Optional)                   │
│   (When sync workflow runs)                     │
│  ┌───────────────────────────────────────────┐  │
│  │ Turso DB → Export Tool → MD files         │  │
│  │ (Every 6 hours or manual trigger)         │  │
│  └──────────┬──────────────────────────────────┘ │
│             │                                    │
│             └─→ Merges with Git posts            │
│                                                  │
└──────────────────────────────────────────────────┘
```

## Current State

### Git-Committed Posts
- **55 posts** from `drafts.json`, `temp-blog.json`, `links-blog.json`
- Always available, never removed
- Version controlled and backed up

### Database Posts (Via Sync Workflow)
- When `sync-posts-turso.yml` runs:
  1. Exports all published posts from Turso
  2. Converts to markdown
  3. Writes to `exports/content/posts/`
  4. Hugo builds the site
  5. Commits and pushes to trigger Vercel deployment

## Testing the Hybrid Approach

```bash
# Current state: 55 git-committed posts
ls exports/content/posts/ | wc -l
# Output: 55

# When sync-posts-turso.yml runs, it will:
# 1. Add any new posts from Turso
# 2. Update existing posts if they changed
# 3. Keep git-committed posts intact

# After sync, the exports/content/posts/ directory will have:
# - All git-committed posts (55)
# - Any posts from Turso database (if any are published)
```

## Merge Behavior

If a post exists in both sources (unlikely but possible):
- **Database export wins** (newer data)
- Slug matching ensures deduplication
- Git history preserves original import

## Advantages

✅ **Persistent History** - Imported posts never lost  
✅ **Version Control** - Full audit trail in Git  
✅ **Dynamic Updates** - New posts via CMS sync  
✅ **Flexibility** - Can add/edit posts two ways  
✅ **Searchability** - All posts indexed and discoverable  

## Implementation Notes

1. **Export Tool** handles both:
   - Creating new markdown from DB
   - Overwriting existing files if needed

2. **Hugo** treats both sources equally:
   - Doesn't distinguish between git-committed and DB posts
   - All posts are published if in exports/content/posts/

3. **Vercel Deployment** triggered by:
   - Push to `public/**` (from deploy-vercel.yml)
   - Automatic when sync workflow commits

## Future Improvements

- [ ] Add CLI to directly create posts without CMS
- [ ] Tag posts by source (imported/database/manual)
- [ ] Create import log for audit trail
- [ ] Implement post merge strategy (prefer git/db)
- [ ] Add post deduplication based on title hash
