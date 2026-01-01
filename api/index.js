import app from '../backend/dist/index.js';

export default async (req, res) => {
  try {
    // Convert Vercel request to fetch API Request
    const url = new URL(req.url, `http://${req.headers.host}`);
    const request = new Request(url, {
      method: req.method,
      headers: req.headers,
      body: ['GET', 'HEAD'].includes(req.method) ? undefined : req
    });

    // Call Hono app
    const response = await app.fetch(request);
    
    // Set response status and headers
    res.status(response.status);
    for (const [key, value] of response.headers) {
      res.setHeader(key, value);
    }
    
    // Send body
    res.send(await response.text());
  } catch (error) {
    console.error('Error:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
};
