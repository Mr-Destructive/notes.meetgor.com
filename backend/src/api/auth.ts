import { Hono } from 'hono';
import { SignJWT } from 'jose';
import type { HonoEnv } from '../types';

export const authRouter = new Hono();

const JWT_SECRET = new TextEncoder().encode(process.env.JWT_SECRET || 'dev-secret-key');
const ADMIN_PASSWORD = process.env.ADMIN_PASSWORD || '';

// Login endpoint
authRouter.post('/login', async (c) => {
  try {
    const { password } = await c.req.json() as { password: string };

    if (!password) {
      return c.json({ error: 'Password required' }, 400);
    }

    // Direct password comparison
    const isValid = password === ADMIN_PASSWORD;

    if (!isValid) {
      return c.json({ error: 'Invalid credentials' }, 401);
    }

    // Create JWT token
    const token = await new SignJWT({ admin: true })
      .setProtectedHeader({ alg: 'HS256' })
      .setExpirationTime('7d')
      .sign(JWT_SECRET);

    return c.json({
      token,
      expires_in: 604800 // 7 days in seconds
    });
  } catch (error) {
    console.error(error);
    return c.json({ error: 'Authentication failed' }, 500);
  }
});

// Verify token endpoint
authRouter.post('/verify', async (c) => {
  try {
    const token = c.req.header('Authorization')?.replace('Bearer ', '');

    if (!token) {
      return c.json({ valid: false }, 401);
    }

    // TODO: verify JWT
    return c.json({ valid: true });
  } catch (error) {
    return c.json({ valid: false }, 401);
  }
});

// Refresh token endpoint
authRouter.post('/refresh', async (c) => {
  try {
    const token = c.req.header('Authorization')?.replace('Bearer ', '');

    if (!token) {
      return c.json({ error: 'No token provided' }, 401);
    }

    // TODO: verify old token and issue new one
    const newToken = await new SignJWT({ admin: true })
      .setProtectedHeader({ alg: 'HS256' })
      .setExpirationTime('7d')
      .sign(JWT_SECRET);

    return c.json({
      token: newToken,
      expires_in: 604800
    });
  } catch (error) {
    return c.json({ error: 'Token refresh failed' }, 500);
  }
});

// Logout (client-side, just remove token)
authRouter.post('/logout', (c) => {
  return c.json({ success: true });
});
