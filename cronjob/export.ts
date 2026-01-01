import { createClient } from '@libsql/client';
import { writeFile, mkdir } from 'fs/promises';
import { join } from 'path';

const turso = createClient({
  url: process.env.TURSO_CONNECTION_URL!,
  authToken: process.env.TURSO_AUTH_TOKEN
});

interface Post {
  id: string;
  type_id: string;
  title: string;
  slug: string;
  content: string;
  excerpt?: string;
  tags: string[];
  metadata: Record<string, any>;
  created_at: string;
  updated_at: string;
  published_at?: string;
}

async function exportPosts() {
  try {
    console.log('Fetching published posts from Turso...');

    // Get all published posts
    const result = await turso.execute(`
      SELECT 
        id, type_id, title, slug, content, excerpt, tags, metadata,
        created_at, updated_at, published_at
      FROM posts
      WHERE status = 'published'
      ORDER BY published_at DESC
    `);

    const posts = result.rows.map(row => ({
      id: row.id as string,
      type_id: row.type_id as string,
      title: row.title as string,
      slug: row.slug as string,
      content: row.content as string,
      excerpt: row.excerpt as string,
      tags: JSON.parse((row.tags as string) || '[]'),
      metadata: JSON.parse((row.metadata as string) || '{}'),
      created_at: row.created_at as string,
      updated_at: row.updated_at as string,
      published_at: row.published_at as string
    })) as Post[];

    console.log(`Found ${posts.length} published posts`);

    // Create content directories
    const contentDir = join(process.cwd(), '../hugo/content');
    const typeDir = join(contentDir, posts[0]?.type_id || 'posts');
    
    await mkdir(typeDir, { recursive: true });

    // Export each post as markdown
    for (const post of posts) {
      const frontmatter = generateFrontmatter(post);
      const markdown = `${frontmatter}\n\n${post.content}`;

      const filePath = join(contentDir, post.type_id, `${post.slug}.md`);
      await writeFile(filePath, markdown, 'utf-8');
      console.log(`✓ Exported: ${post.slug}`);
    }

    console.log(`\nSuccessfully exported ${posts.length} posts`);

  } catch (error) {
    console.error('Export failed:', error);
    process.exit(1);
  }
}

function generateFrontmatter(post: Post): string {
  const fm = {
    title: post.title,
    slug: post.slug,
    date: post.published_at || post.created_at,
    lastmod: post.updated_at,
    tags: post.tags,
    type: post.type_id,
    ...(post.excerpt && { summary: post.excerpt }),
    ...post.metadata
  };

  return `---\n${Object.entries(fm)
    .filter(([, v]) => v !== undefined && v !== null)
    .map(([k, v]) => {
      if (Array.isArray(v)) {
        return `${k}: [${v.map(item => `"${item}"`).join(', ')}]`;
      }
      if (typeof v === 'object') {
        return `${k}: ${JSON.stringify(v)}`;
      }
      return `${k}: "${v}"`;
    })
    .join('\n')}\n---`;
}

// Run export
exportPosts().catch(console.error);
