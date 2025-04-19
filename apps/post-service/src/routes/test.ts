import { Router } from 'express';
import { PrismaClient } from '@prisma/client';

const router = Router();
const prisma = new PrismaClient();

router.post('/cleanup', async (_req, res) => {
  if (process.env.NODE_ENV !== 'test') return res.status(403).send('Forbidden');

  await prisma.post.deleteMany(); // post-service
  res.sendStatus(204);
});

export default router;