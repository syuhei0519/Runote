// src/routes/test.ts
import { Router } from 'express';
import { prisma } from '../lib/prisma';

const router = Router();

router.post('/test/cleanup', async (_req, res) => {
  if (process.env.NODE_ENV !== 'test') return res.status(403).send('Forbidden');

  await prisma.post.deleteMany(); // post-service
  // await prisma.tag.deleteMany(); // tag-service
  res.status(204).send();
});

export default router;