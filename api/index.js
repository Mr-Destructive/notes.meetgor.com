const { Hono } = require('hono');
const { serve } = require('@hono/node-server');
const { cors } = require('hono/cors');
const { logger } = require('hono/logger');
const { createClient } = require('@libsql/client');

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
  c.env = { db };
  await next();
});

// Health check
app.get('/health', (c) => {
  return c.json({ status: 'ok' });
});

// Placeholder API routes (full implementation in backend)
app.post('/api/auth/login', async (c) => {
  try {
    const { password } = await c.req.json();
    
    // Direct password comparison
    const isValid = password === process.env.ADMIN_PASSWORD;
    
    if (!isValid) {
      return c.json({ error: 'Invalid credentials' }, 401);
    }
    
    // For now, return a simple token
    return c.json({
      token: 'temp-token-' + Date.now(),
      expires_in: 604800
    });
  } catch (error) {
    return c.json({ error: 'Authentication failed' }, 500);
  }
});

// 404 handler
app.notFound((c) => {
  return c.json({ error: 'Not found' }, 404);
});

// Error handler
app.onError((err, c) => {
  console.error('Error:', err);
  return c.json({ error: 'Internal server error' }, 500);
});

// Export for Vercel
module.exports = app;
