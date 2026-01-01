import 'dotenv/config';
import { serve } from '@hono/node-server';
import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { logger } from 'hono/logger';
import { createClient } from '@libsql/client';
import { authRouter } from './api/auth';
import { postsRouter } from './api/posts';
import type { HonoEnv } from './types';

const app = new Hono<HonoEnv>();

// Middleware
app.use(logger());
app.use(cors({
  origin: '*',
  allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowHeaders: ['Content-Type', 'Authorization']
}));

// Initialize database client
const db = createClient({
  url: process.env.TURSO_CONNECTION_URL || 'file:./dev.db',
  authToken: process.env.TURSO_AUTH_TOKEN
});

// Share db instance
app.use(async (c, next) => {
  (c as any).env = { db };
  await next();
});

// Health check
app.get('/health', (c) => {
  return c.json({ status: 'ok' });
});

// API routes
app.route('/api/auth', authRouter);
app.route('/api/posts', postsRouter);

// 404 handler
app.notFound((c) => {
  return c.json({ error: 'Not found' }, 404);
});

// Error handler
app.onError((err, c) => {
  console.error('Error:', err);
  return c.json({ error: 'Internal server error' }, 500);
});

// For serverless (Vercel, etc)
module.exports = app;
module.exports.default = app;

// Only start server in development (local, not serverless)
if (require.main === module) {
  const port = parseInt(process.env.API_PORT || '3000');
  serve({
    fetch: app.fetch,
    port
  });
  console.log(`🚀 Blog API running on http://localhost:${port}`);
}
