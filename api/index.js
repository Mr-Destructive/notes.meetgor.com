import '../backend/dist/index.js';
import app from '../backend/dist/index.js';

export default (req, res) => {
  return app.fetch(req);
};
