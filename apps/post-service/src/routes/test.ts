// src/routes/test.ts
import { Router } from 'express';
import { PrismaClient } from '@prisma/client';

const router = Router();
const prisma = new PrismaClient();

router.post('/test/cleanup', async (_req, res) => {
  if (process.env.NODE_ENV !== 'test') return res.status(403).send('Forbidden');

  await prisma.post.deleteMany(); // post-service
  // await prisma.tag.deleteMany(); // tag-service
  res.status(204).send();
});

export default router;