import { z } from 'zod';

const EmotionSchema = z.object({
  type: z.enum(['preset', 'custom']),
  label: z.string().min(1, '感情ラベルは必須です'),
  presetKey: z.string().optional(),
  intensity: z.number().int().min(1).max(5)
});

export const PostSchema = z.object({
  title: z.string().min(1, 'タイトルは必須です'),
  content: z.string().min(1, '本文は必須です'),
  tagIds: z.array(z.number().int().positive()).optional(),
  emotion: EmotionSchema.optional()
});

export type PostInput = z.infer<typeof PostSchema>;