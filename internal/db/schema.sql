-- Post types
CREATE TABLE IF NOT EXISTS post_types (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  description TEXT
);

-- Main posts table
CREATE TABLE IF NOT EXISTS posts (
  id TEXT PRIMARY KEY,
  type_id TEXT NOT NULL,
  title TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  content TEXT NOT NULL,
  excerpt TEXT,
  status TEXT DEFAULT 'draft', -- draft, published, archived
  is_featured BOOLEAN DEFAULT 0,
  tags TEXT, -- JSON array ["tag1", "tag2"]
  metadata TEXT, -- JSON for type-specific data
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  published_at DATETIME,
  FOREIGN KEY(type_id) REFERENCES post_types(id)
);

-- Revisions for version history
CREATE TABLE IF NOT EXISTS revisions (
  id TEXT PRIMARY KEY,
  post_id TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
);

-- Series/Collections
CREATE TABLE IF NOT EXISTS series (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  description TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Post to series mapping
CREATE TABLE IF NOT EXISTS post_series (
  post_id TEXT NOT NULL,
  series_id TEXT NOT NULL,
  order_in_series INT,
  PRIMARY KEY(post_id, series_id),
  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY(series_id) REFERENCES series(id) ON DELETE CASCADE
);

-- Settings
CREATE TABLE IF NOT EXISTS settings (
  key TEXT PRIMARY KEY,
  value TEXT
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_posts_type ON posts(type_id);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at);
CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_series_slug ON series(slug);

-- Insert default post types
INSERT OR IGNORE INTO post_types (id, name, slug, description) VALUES
  ('article', 'Article', 'article', 'Full-length articles'),
  ('review', 'Review', 'review', 'Book, movie, or product reviews'),
  ('thought', 'Thought', 'thought', 'Quick thoughts and reflections'),
  ('link', 'Link', 'link', 'Curated links with commentary'),
  ('til', 'TIL', 'til', 'Today I Learned'),
  ('quote', 'Quote', 'quote', 'Quotations and excerpts'),
  ('list', 'List', 'list', 'Curated lists'),
  ('note', 'Note', 'note', 'Quick notes'),
  ('snippet', 'Snippet', 'snippet', 'Code snippets'),
  ('essay', 'Essay', 'essay', 'Long-form essays'),
  ('tutorial', 'Tutorial', 'tutorial', 'Step-by-step guides'),
  ('interview', 'Interview', 'interview', 'Q&A interviews');
