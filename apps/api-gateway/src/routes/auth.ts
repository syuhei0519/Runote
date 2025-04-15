import express from 'express';
import { proxyRequest } from '../utils/proxy';

const router = express.Router();
const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://auth-service:8000';

router.all('/*', (req, res) => {
  const targetUrl = `${AUTH_SERVICE_URL}${req.originalUrl}`;
  proxyRequest(req, res, targetUrl);
});

export default router;