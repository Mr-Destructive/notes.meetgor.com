import { Hono } from 'hono';
import { randomUUID } from 'crypto';
import type { Post, PostForm, PostStatus, HonoEnv } from '../types.ts';

export const postsRouter = new Hono<HonoEnv>();

// Middleware: auth check
const requireAuth = async (c: any, next: any) => {
  const token = c.req.header('Authorization')?.replace('Bearer ', '');
  if (!token) {
    return c.json({ error: 'Unauthorized' }, 401);
  }
  // TODO: verify JWT token
  await next();
};

// Create post
postsRouter.post('/', requireAuth, async (c) => {
  try {
    const body: PostForm = await c.req.json();
    const db = (c.env as any).db;

    // Validate required fields
    if (!body.title || !body.content || !body.type_id) {
      return c.json({ error: 'Missing required fields' }, 400);
    }

    // Generate slug if not provided
    const slug = body.slug || generateSlug(body.title);

    // Check slug uniqueness
    const existing = await db.execute(
      'SELECT id FROM posts WHERE slug = ?',
      [slug]
    );
    if (existing.rows.length > 0) {
      return c.json({ error: 'Slug already exists' }, 409);
    }

    const id = randomUUID();
    const now = new Date().toISOString();

    await db.execute(
      `INSERT INTO posts (id, type_id, title, slug, content, excerpt, status, tags, metadata, created_at, updated_at)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      [
        id,
        body.type_id,
        body.title,
        slug,
        body.content,
        body.excerpt || null,
        body.status || 'draft',
        JSON.stringify(body.tags || []),
        JSON.stringify(body.metadata || {}),
        now,
        now
      ]
    );

    return c.json({ id, slug, created_at: now }, 201);
  } catch (error) {
    console.error(error);
    return c.json({ error: 'Failed to create post' }, 500);
  }
});

// Get all posts (with filters)
postsRouter.get('/', async (c) => {
  try {
    const db = (c.env as any).db;
    const type = c.req.query('type');
    const status = c.req.query('status') || 'published';
    const limit = parseInt(c.req.query('limit') || '50');
    const offset = parseInt(c.req.query('offset') || '0');

    let query = 'SELECT * FROM posts WHERE 1=1';
    const params: any[] = [];

    if (type) {
      query += ' AND type_id = ?';
      params.push(type);
    }

    query += ' AND status = ?';
    params.push(status);

    query += ' ORDER BY published_at DESC, created_at DESC LIMIT ? OFFSET ?';
    params.push(limit, offset);

    const result = await db.execute(query, params);
    const posts = result.rows.map(formatPost);

    return c.json(posts);
  } catch (error) {
    console.error(error);
    return c.json({ error: 'Failed to fetch posts' }, 500);
  }
});

// Get single post
postsRouter.get('/:id', async (c) => {
  try {
    const db = (c.env as any).db;
    const result = await db.execute('SELECT * FROM posts WHERE id = ? OR slug = ?', [
      c.req.param('id'),
      c.req.param('id')
    ]);

    if (result.rows.length === 0) {
      return c.json({ error: 'Post not found' }, 404);
    }

    return c.json(formatPost(result.rows[0]));
  } catch (error) {
    return c.json({ error: 'Failed to fetch post' }, 500);
  }
});

// Update post
postsRouter.put('/:id', requireAuth, async (c) => {
  try {
    const db = (c.env as any).db;
    const body: Partial<PostForm> = await c.req.json();
    const id = c.req.param('id');

    // Check post exists
    const existing = await db.execute('SELECT id FROM posts WHERE id = ?', [id]);
    if (existing.rows.length === 0) {
      return c.json({ error: 'Post not found' }, 404);
    }

    const now = new Date().toISOString();
    const updates: string[] = [];
    const values: any[] = [];

    if (body.title) {
      updates.push('title = ?');
      values.push(body.title);
    }
    if (body.content) {
      updates.push('content = ?');
      values.push(body.content);
    }
    if (body.excerpt !== undefined) {
      updates.push('excerpt = ?');
      values.push(body.excerpt);
    }
    if (body.status) {
      updates.push('status = ?');
      values.push(body.status);
    }
    if (body.tags) {
      updates.push('tags = ?');
      values.push(JSON.stringify(body.tags));
    }
    if (body.metadata) {
      updates.push('metadata = ?');
      values.push(JSON.stringify(body.metadata));
    }
    if (body.is_featured !== undefined) {
      updates.push('is_featured = ?');
      values.push(body.is_featured ? 1 : 0);
    }

    if (body.status === 'published' && existing.rows[0].status !== 'published') {
      updates.push('published_at = ?');
      values.push(now);
    }

    updates.push('updated_at = ?');
    values.push(now);
    values.push(id);

    const query = `UPDATE posts SET ${updates.join(', ')} WHERE id = ?`;
    await db.execute(query, values);

    return c.json({ success: true, updated_at: now });
  } catch (error) {
    console.error(error);
    return c.json({ error: 'Failed to update post' }, 500);
  }
});

// Delete post
postsRouter.delete('/:id', requireAuth, async (c) => {
  try {
    const db = (c.env as any).db;
    await db.execute('DELETE FROM posts WHERE id = ?', [c.req.param('id')]);
    return c.json({ success: true });
  } catch (error) {
    return c.json({ error: 'Failed to delete post' }, 500);
  }
});

// Get revisions
postsRouter.get('/:id/revisions', async (c) => {
  try {
    const db = (c.env as any).db;
    const result = await db.execute(
      'SELECT id, content, change_summary, created_at FROM revisions WHERE post_id = ? ORDER BY created_at DESC',
      [c.req.param('id')]
    );
    return c.json(result.rows);
  } catch (error) {
    return c.json({ error: 'Failed to fetch revisions' }, 500);
  }
});

// Helpers
function generateSlug(title: string): string {
  return title
    .toLowerCase()
    .replace(/[^\w\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .trim();
}

function formatPost(row: any): Post {
  return {
    id: row.id,
    type_id: row.type_id,
    title: row.title,
    slug: row.slug,
    content: row.content,
    excerpt: row.excerpt,
    status: row.status,
    is_featured: row.is_featured === 1,
    tags: JSON.parse(row.tags || '[]'),
    metadata: JSON.parse(row.metadata || '{}'),
    created_at: row.created_at,
    updated_at: row.updated_at,
    published_at: row.published_at
  };
}
