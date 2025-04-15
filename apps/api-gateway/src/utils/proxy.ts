import { Request, Response } from 'express';
import axios from 'axios';
import jwt from 'jsonwebtoken';
import dotenv from 'dotenv';

dotenv.config();

const JWT_SECRET = process.env.JWT_SECRET || 'default-secret';

export async function proxyRequest(req: Request, res: Response, targetUrl: string) {
  const authHeader = req.headers.authorization;

  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return res.status(401).json({ error: 'Authorization token missing or malformed' });
  }

  const token = authHeader.split(' ')[1];
  let decoded: string | jwt.JwtPayload;

  try {
    decoded = jwt.verify(token, JWT_SECRET);
    console.log('JWT Verified:', decoded);
  } catch (err) {
    return res.status(403).json({ error: 'Invalid token' });
  }

  try {
    const result = await axios({
      method: req.method as any,
      url: targetUrl,
      headers: {
        ...req.headers,
        'X-User-Id': (decoded as any).sub,
      },
      data: req.body,
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