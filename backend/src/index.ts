import 'dotenv/config';
import { serve } from '@hono/node-server';
import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { logger } from 'hono/logger';
import { serveStatic } from 'hono/node-server/serve-static';
import { createClient } from '@libsql/client';
import { authRouter } from './api/auth.ts';
import { postsRouter } from './api/posts.ts';

const app = new Hono();

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
  c.env.db = db;
  await next();
});

// Health check
app.get('/health', (c) => {
  return c.json({ status: 'ok' });
});

// API routes
app.route('/api/auth', authRouter);
app.route('/api/posts', postsRouter);

// Serve frontend static files
app.use('/js/*', serveStatic({ root: '../frontend' }));
app.use('/html/*', serveStatic({ root: '../frontend' }));
app.use('/', serveStatic({ path: '../frontend/login.html', root: '.' }));

// Serve login.html as root
app.get('/', (c) => {
  return c.html(`
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="utf-8">
      <meta name="viewport" content="width=device-width, initial-scale=1">
      <title>Blog Admin</title>
    </head>
    <body>
      <script>
        window.location.href = '/login.html';
      </script>
    </body>
    </html>
  `);
});

// Serve login page
app.get('/login.html', serveStatic({ path: '../frontend/login.html' }));

// Serve editor page
app.get('/editor.html', serveStatic({ path: '../frontend/editor.html' }));

// Serve editor JS
app.get('/editor.js', serveStatic({ path: '../frontend/editor.js' }));

// 404 handler
app.notFound((c) => {
  return c.json({ error: 'Not found' }, 404);
});

// Error handler
app.onError((err, c) => {
  console.error('Error:', err);
  return c.json({ error: 'Internal server error' }, 500);
});

// Start server
const port = parseInt(process.env.API_PORT || '3000');

serve({
  fetch: app.fetch,
  port
});

console.log(`🚀 Blog API running on http://localhost:${port}`);

// For serverless (Vercel, etc)
export default app;
