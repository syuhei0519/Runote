import express from 'express';
import { proxyRequest } from '../utils/proxy';

const router = express.Router();
const TAG_SERVICE_URL = process.env.TAG_SERVICE_URL || 'http://tag-service:4000';

router.all('/*', (req, res) => {
  const targetUrl = `${TAG_SERVICE_URL}${req.originalUrl}`;
  proxyRequest(req, res, targetUrl);
});

export default router;