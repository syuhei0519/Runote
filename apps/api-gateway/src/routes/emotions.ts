import express from 'express';
import { proxyRequest } from '../utils/proxy';

const router = express.Router();
const EMOTION_SERVICE_URL = process.env.EMOTION_SERVICE_URL || 'http://emotion-service:8080';

router.all('/*', (req, res) => {
  const targetUrl = `${EMOTION_SERVICE_URL}${req.originalUrl}`;
  proxyRequest(req, res, targetUrl);
});

export default router;