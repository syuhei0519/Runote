// utils/proxy.ts
import { Request, Response } from 'express';
import axios from 'axios';
import jwt from 'jsonwebtoken';
import dotenv from 'dotenv';

dotenv.config();

const JWT_SECRET = process.env.JWT_SECRET || 'supersecretkey';

// 認証不要ルート（method + path をセットで管理）
const noAuthRoutes = [
  { method: 'POST', path: '/auth/login' },
  { method: 'POST', path: '/auth/register' },
  { method: 'POST', path: '/auth/test/cleanup' },
  { method: 'POST', path: '/emotions/test/cleanup' },
  { method: 'POST', path: '/posts/test/cleanup' },
  { method: 'POST', path: '/tags/test/cleanup' },
  { method: 'OPTIONS', path: '*' } // CORSプリフライトをスキップ対象に
];

function shouldSkipAuth(req: Request): boolean {
  return noAuthRoutes.some(route => {
    const methodMatches = req.method === route.method;
    const pathMatches =
      route.path === '*' || req.originalUrl.startsWith(route.path);
    return methodMatches && pathMatches;
  });
}

export async function proxyRequest(req: Request, res: Response, targetUrl: string) {
  if (shouldSkipAuth(req)) {
    try {
      const result = await axios({
        method: req.method as any,
        url: targetUrl,
        headers: {
          ...req.headers,
          'Content-Type': 'application/json',
        },
        data: req.body || {},
      });
      return res.status(result.status).json(result.data);
    } catch (err: any) {
      return res.status(err.response?.status || 500).json({
        error: 'Upstream service error',
        detail: err.response?.data || err.message,
      });
    }
  }

  // 🔐 認証チェック（スキップ対象以外）
  const authHeader = req.headers.authorization;
  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return res.status(401).json({ error: 'Authorization token missing or malformed' });
  }

  const token = authHeader.split(' ')[1];
  let decoded: string | jwt.JwtPayload;

  try {
    decoded = jwt.verify(token, JWT_SECRET);
  } catch (err: any) {
    if (err.name === 'TokenExpiredError') {
      return res.status(401).json({ error: 'Token expired' });
    }
    return res.status(403).json({ error: 'Invalid token', reason: err.message });
  }

  // 🔁 サービスに中継
  try {
    const result = await axios({
      method: req.method as any,
      url: targetUrl,
      headers: {
        ...req.headers,
        'Authorization': req.headers.authorization,
        'X-User-Id': (decoded as any).sub,
      },
      data: req.body
    });

    res.status(result.status).json(result.data);
  } catch (err: any) {
    console.error(`[Proxy Error] ${err.message}`);
    res.status(err.response?.status || 500).json({
      error: 'Upstream service error',
      detail: err.message,
    });
  }
}
