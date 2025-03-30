import { z } from 'zod';

export const PostSchema = z.object({
  title: z.string().min(1, 'タイトルは必須です'),
  content: z.string().min(1, '本文は必須です'),
});

export type PostInput = z.infer<typeof PostSchema>;