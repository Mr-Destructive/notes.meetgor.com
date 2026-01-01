const app = require('../backend/dist/index.js').default;

module.exports = async (req, res) => {
  try {
    // Build full URL
    const proto = req.headers['x-forwarded-proto'] || 'http';
    const host = req.headers['x-forwarded-host'] || req.headers.host;
    const url = new URL(req.url, `${proto}://${host}`);

    // Create fetch-compatible request
    const request = new Request(url.toString(), {
      method: req.method,
      headers: new Map(Object.entries(req.headers)),
      body: ['GET', 'HEAD'].includes(req.method) ? undefined : req.body
    });

    // Call Hono app
    const response = await app.fetch(request);
    
    // Send response
    res.status(response.status);
    response.headers.forEach((value, key) => {
      res.setHeader(key, value);
    });
    
    res.send(await response.text());
  } catch (error) {
    console.error('Error:', error);
    res.status(500).json({ error: error.message });
  }
};
