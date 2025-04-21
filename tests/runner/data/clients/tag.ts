// tag-serviceテストデータ登録用DB接続用クライアントを提供(PostgreSQL)
import { Client } from 'pg';
import dotenv from 'dotenv';

dotenv.config();

export async function getTagClient(): Promise<Client> {
  const databaseUrl = process.env.TAG_DATABASE_URL || 'postgresql://taguser:tagpass@tag-test-db:5432/tag_test';

  if (!databaseUrl) {
    throw new Error('POST_DATABASE_URL is not defined in environment variables');
  }

  const client = new Client({
    connectionString: databaseUrl,
  });

  await client.connect();
  return client;
}