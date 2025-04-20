// post-serviceテストデータ登録用DB接続用クライアントを提供(PostgreSQL)
import { Client } from 'pg';
import dotenv from 'dotenv';

dotenv.config();

export async function getPostClient(): Promise<Client> {
  const databaseUrl = process.env.POST_DATABASE_URL || 'postgresql://test:test@post-test-db:5432/post_test';

  if (!databaseUrl) {
    throw new Error('POST_DATABASE_URL is not defined in environment variables');
  }

  const client = new Client({
    connectionString: databaseUrl,
  });

  await client.connect();
  return client;
}