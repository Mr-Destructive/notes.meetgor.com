---
title: "Notes"
type: "page"
---

A collection of technical notes and posts

## Browse by Type

<style>
.post-types {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 1rem;
  margin-top: 2rem;
}

.post-type-link {
  padding: 1.25rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.375rem;
  text-decoration: none;
  color: inherit;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  text-align: center;
}

.post-type-link:hover {
  border-color: #3b82f6;
  background-color: #f8fafc;
}

.post-type-title {
  font-weight: 600;
  color: #1f2937;
  font-size: 0.95rem;
}

.post-type-count {
  font-size: 0.75rem;
  color: #9ca3af;
}
</style>

<div class="post-types">
<a href="/posts/" class="post-type-link">
<span class="post-type-title">All Posts</span>
<span class="post-type-count">Browse all</span>
</a>
</div>

Or view by category at `/type/`
