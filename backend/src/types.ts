export type PostType = 
  | 'article' 
  | 'review' 
  | 'thought' 
  | 'link' 
  | 'til' 
  | 'quote' 
  | 'list' 
  | 'note' 
  | 'snippet' 
  | 'essay' 
  | 'tutorial' 
  | 'interview' 
  | 'experiment';

export type PostStatus = 'draft' | 'published' | 'archived';

export interface Post {
  id: string;
  type_id: PostType;
  title: string;
  slug: string;
  content: string;
  excerpt?: string;
  status: PostStatus;
  is_featured: boolean;
  tags: string[];
  metadata: Record<string, any>;
  created_at: string;
  updated_at: string;
  published_at?: string;
}

export interface PostForm {
  type_id: PostType;
  title: string;
  slug?: string;
  content: string;
  excerpt?: string;
  status?: PostStatus;
  is_featured?: boolean;
  tags?: string[];
  metadata?: Record<string, any>;
}

// Type-specific metadata shapes
export interface ReviewMetadata {
  type: 'book' | 'movie' | 'product' | 'article';
  rating?: number; // 1-5
  author?: string;
  subject_link?: string;
}

export interface LinkMetadata {
  source_url: string;
  domain?: string;
}

export interface TilMetadata {
  category?: string;
  difficulty?: 'beginner' | 'intermediate' | 'advanced';
}

export interface QuoteMetadata {
  author: string;
  source?: string;
  context?: string;
}

export interface ListMetadata {
  items: string[];
  list_type: 'ordered' | 'unordered';
}

export interface SnippetMetadata {
  language: string;
  syntax_highlight?: string;
}

export interface EssayMetadata {
  reading_time?: number;
  word_count?: number;
}

export interface TutorialMetadata {
  difficulty?: 'beginner' | 'intermediate' | 'advanced';
  estimated_time?: number; // in minutes
  tools?: string[];
}

export interface InterviewMetadata {
  interviewee_name: string;
  role?: string;
  company?: string;
}

// Template definitions for each post type
export const POST_TEMPLATES = {
  article: {
    fields: ['title', 'content', 'excerpt', 'tags'],
    defaultMetadata: { category: '' }
  },
  review: {
    fields: ['title', 'subject', 'rating', 'content', 'tags'],
    defaultMetadata: { type: 'book', rating: 5 } as ReviewMetadata
  },
  thought: {
    fields: ['title', 'content', 'tags'],
    defaultMetadata: {}
  },
  link: {
    fields: ['title', 'url', 'excerpt', 'tags'],
    defaultMetadata: { source_url: '' } as LinkMetadata
  },
  til: {
    fields: ['title', 'content', 'tags', 'category'],
    defaultMetadata: { difficulty: 'intermediate' } as TilMetadata
  },
  quote: {
    fields: ['text', 'author', 'source', 'tags'],
    defaultMetadata: { author: '', source: '' } as QuoteMetadata
  },
  list: {
    fields: ['title', 'items', 'list_type', 'tags'],
    defaultMetadata: { items: [], list_type: 'unordered' } as ListMetadata
  },
  note: {
    fields: ['title', 'content', 'tags'],
    defaultMetadata: { is_private: false }
  },
  snippet: {
    fields: ['title', 'code', 'language', 'tags'],
    defaultMetadata: { language: 'javascript' } as SnippetMetadata
  },
  essay: {
    fields: ['title', 'content', 'excerpt', 'tags'],
    defaultMetadata: {} as EssayMetadata
  },
  tutorial: {
    fields: ['title', 'content', 'difficulty', 'estimated_time', 'tags'],
    defaultMetadata: { difficulty: 'beginner' } as TutorialMetadata
  },
  interview: {
    fields: ['title', 'content', 'interviewee_name', 'role', 'tags'],
    defaultMetadata: { interviewee_name: '' } as InterviewMetadata
  },
  experiment: {
    fields: ['title', 'content', 'tags', 'tools'],
    defaultMetadata: {}
  }
} as const;
