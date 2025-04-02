import { Router, Request, Response } from 'express';
import { PrismaClient } from '@prisma/client';
import { TagSchema } from '../validators/tag';
import { ZodError } from 'zod';

const router = Router();
const prisma = new PrismaClient();

// GET /tags
router.get('/', async (_req: Request, res: Response) => {
  const tags = await prisma.tag.findMany({ orderBy: { createdAt: 'desc' } });
  res.json(tags);
});

// GET /tags/:id
router.get('/:id', async (req: Request, res: Response) => {
  const id = Number(req.params.id);
  if (isNaN(id)) return res.status(400).json({ error: 'Invalid ID' });

  const tag = await prisma.tag.findUnique({ where: { id } });
  if (!tag) return res.status(404).json({ error: 'Tag not found' });

  res.json(tag);
});

// POST /tags
router.post('/', async (req: Request, res: Response) => {
  const parsed = TagSchema.safeParse(req.body);
  if (!parsed.success) {
    return res.status(400).json({ error: parsed.error.format() });
  }

  const newTag = await prisma.tag.create({
    data: parsed.data,
  });

  res.status(201).json(newTag);
});

// PUT /tags/:id
router.put('/:id', async (req: Request, res: Response) => {
  const id = Number(req.params.id);
  if (isNaN(id)) return res.status(400).json({ error: 'Invalid ID' });

  try {
    const parsed = TagSchema.parse(req.body);
    const existing = await prisma.tag.findUnique({ where: { id } });
    if (!existing) return res.status(404).json({ error: 'Tag not found' });

    const updated = await prisma.tag.update({
      where: { id },
      data: parsed,
    });

    return res.json(updated);
  } catch (err) {
    if (err instanceof ZodError) {
      return res.status(400).json({ errors: err.errors });
    }
    return res.status(500).json({ error: 'Internal Server Error' });
  }
});

// DELETE /tags/:id
router.delete('/:id', async (req: Request, res: Response) => {
  const id = Number(req.params.id);
  if (isNaN(id)) return res.status(400).json({ error: 'Invalid ID' });

  const existing = await prisma.tag.findUnique({ where: { id } });
  if (!existing) return res.status(404).json({ error: 'Tag not found' });

  await prisma.tag.delete({ where: { id } });
  res.status(204).send();
});

export default router;
