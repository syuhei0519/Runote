// auth-serviceテストデータ登録用DB接続用クライアントを提供(MySQL)
import mysql from 'mysql2/promise';
import dotenv from 'dotenv';

dotenv.config();

export async function getAuthClient(): Promise<mysql.Connection> {
  const databaseUrl = process.env.AUTH_DATABASE_URL || 'mysql+pymysql://test:test@auth-test-db:3306/auth_test';
  if (!databaseUrl) {
    throw new Error('AUTH_DATABASE_URL is not defined in environment variables');
  }

  const connection = await mysql.createConnection(databaseUrl);
  return connection;
}