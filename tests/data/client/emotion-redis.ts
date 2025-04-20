// tests/data/redisClient.ts
import Redis from 'ioredis';
import dotenv from 'dotenv';

dotenv.config();

export function getEmotionRedisClient(): Redis {
  const redisUrl = process.env.EMOTION_REDIS_URL || 'redis://emotion-test-redis:6379/0';

  return new Redis(redisUrl);
}