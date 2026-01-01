const app = require('../backend/dist/index.js');

module.exports = async (req, res) => {
  try {
    const response = await app.fetch(req);
    res.status(response.status);
    for (const [key, value] of response.headers.entries()) {
      res.setHeader(key, value);
    }
    res.send(await response.text());
  } catch (error) {
    console.error(error);
    res.status(500).json({ error: 'Internal server error' });
  }
};
