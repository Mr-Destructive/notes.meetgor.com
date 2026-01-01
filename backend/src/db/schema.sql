-- Post types table
CREATE TABLE IF NOT EXISTS post_types (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  slug TEXT UNIQUE NOT NULL,
  description TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Main posts table
CREATE TABLE IF NOT EXISTS posts (
  id TEXT PRIMARY KEY,
  type_id TEXT NOT NULL,
  title TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  content TEXT NOT NULL,
  excerpt TEXT,
  status TEXT DEFAULT 'draft',
  is_featured BOOLEAN DEFAULT 0,
  tags TEXT,
  metadata TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  published_at DATETIME,
  FOREIGN KEY(type_id) REFERENCES post_types(id) ON DELETE RESTRICT
);

-- Revision history
CREATE TABLE IF NOT EXISTS revisions (
  id TEXT PRIMARY KEY,
  post_id TEXT NOT NULL,
  content TEXT NOT NULL,
  change_summary TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
);

-- App settings
CREATE TABLE IF NOT EXISTS settings (
  key TEXT PRIMARY KEY,
  value TEXT,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_posts_type ON posts(type_id);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_published ON posts(published_at) WHERE status = 'published';
CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
CREATE INDEX IF NOT EXISTS idx_revisions_post ON revisions(post_id);

-- Insert default post types
INSERT OR IGNORE INTO post_types (id, name, slug, description) VALUES
('article', 'Article', 'article', 'Long-form technical or editorial content'),
('review', 'Review', 'review', 'Book, movie, product, or media reviews'),
('thought', 'Thought', 'thought', 'Quick thoughts and reflections'),
('link', 'Link', 'link', 'Curated links with commentary'),
('til', 'TIL', 'til', 'Today I Learned - quick learnings'),
('quote', 'Quote', 'quote', 'Interesting quotes and citations'),
('list', 'List', 'list', 'Curated lists and collections'),
('note', 'Note', 'note', 'Quick notes and snippets'),
('snippet', 'Snippet', 'snippet', 'Code snippets with syntax highlighting'),
('essay', 'Essay', 'essay', 'Personal essays and long-form writing'),
('tutorial', 'Tutorial', 'tutorial', 'Step-by-step guides and tutorials'),
('interview', 'Interview', 'interview', 'Q&A format interviews'),
('experiment', 'Experiment', 'experiment', 'Technical experiments and learnings');
