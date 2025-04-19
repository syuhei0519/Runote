// src/routes/test.ts
import { Router } from 'express';
import { PrismaClient } from '@prisma/client';

const prisma = new PrismaClient();

const router = Router();

router.post('/cleanup', async (_req, res) => {
  if (process.env.NODE_ENV !== 'test') {
    return res.status(403).json({ message: 'Forbidden' });
  }

  try {
    // 他のリレーションがある場合は先に削除
    await prisma.postTag?.deleteMany(); // 中間テーブルがある場合（存在しなければ削除OK）
    await prisma.tag.deleteMany();

    res.status(200).send();
  } catch (error) {
    console.error('Cleanup Error:', error);
    res.status(500).json({ message: 'Cleanup failed' });
  }
});

export default router;