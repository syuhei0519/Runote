import { Router, Request, Response } from 'express';
import { PrismaClient } from '@prisma/client';
import { PostSchema } from '../validators/post';
import type { z } from 'zod';
import { ZodError } from 'zod';
import axios from 'axios';
import type { PostInput } from '../validators/post';

const router = Router();
const prisma = new PrismaClient();

// GET /posts
router.get('/', async (_req: Request, res: Response) => {
    const posts = await prisma.post.findMany({ orderBy: { createdAt: 'desc' } });
    res.json(posts);
  });
  
// GET /posts/:id
router.get('/:id', async (req: Request, res: Response) => {
    const id = Number(req.params.id);
    if (isNaN(id)) {
      return res.status(400).json({ error: 'Invalid ID' });
    }
  
    const post = await prisma.post.findUnique({ where: { id } });
    if (!post) return res.status(404).json({ error: 'Post not found' });
  
    res.json(post);
});
  
// POST /posts
router.post('/', async (req: Request, res: Response) => {
  const parsed = PostSchema.safeParse(req.body);
  if (!parsed.success) {
    return res.status(400).json({ error: parsed.error.format() });
  }

  const { title, content, tagIds = [], emotion } = parsed.data as PostInput;

  try {
    const newPost = await prisma.post.create({
      data: {
        title,
        content,
        tags: {
          create: tagIds.map(tagId => ({ tagId }))
        }
      }
    });

    // emotion がある場合だけ emotion-service を呼ぶ、失敗しても無視
    if (emotion?.label && emotion?.type && typeof emotion.intensity === 'number') {
      try {
        await axios.post('http://emotion-service:8080/emotions', {
          postId: newPost.id,
          ...emotion
        });
      } catch (e) {
        console.error('emotion-service unreachable');
      }
    }    

    res.status(201).json(newPost);
  } catch (error) {
    console.error('投稿作成失敗:', error);
    res.status(500).json({ error: 'Internal Server Error' });
  }
});

// PUT /posts/:id
router.put('/:id', async (req: Request, res: Response) => {
    const id = Number(req.params.id);
    if (isNaN(id)) return res.status(400).json({ error: 'Invalid ID' });

    try {
      const parsed = PostSchema.parse(req.body);
      const existing = await prisma.post.findUnique({ where: { id } });
      if (!existing) return res.status(404).json({ error: 'Post not found' });
  
      const updated = await prisma.post.update({
        where: { id },
        data: parsed,
      });
  
      return res.json(updated);
    } catch (err) {
      if (err instanceof ZodError) {
        return res.status(400).json({ errors: err.errors }); // ✅ Zodのエラー配列をそのまま返す
      }
      return res.status(500).json({ error: 'Internal Server Error' });
    }
});
  
// DELETE /posts/:id
router.delete('/:id', async (req: Request, res: Response) => {
    const id = Number(req.params.id);
    if (isNaN(id)) return res.status(400).json({ error: 'Invalid ID' });
  
    const existing = await prisma.post.findUnique({ where: { id } });
    if (!existing) return res.status(404).json({ error: 'Post not found' });
  
    await prisma.post.delete({ where: { id } });
    res.status(204).send(); // ✅ No Content を返す
});
  
export default router;