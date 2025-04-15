import express from 'express';
import { proxyRequest } from '../utils/proxy';

const router = express.Router();
const POST_SERVICE_URL = process.env.POST_SERVICE_URL || 'http://post-service:3000';

router.all('/*', (req, res) => {
  const targetUrl = `${POST_SERVICE_URL}${req.originalUrl}`;
  proxyRequest(req, res, targetUrl);
});

export default router;