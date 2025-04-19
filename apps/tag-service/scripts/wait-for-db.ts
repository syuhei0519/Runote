import { Client } from 'pg';

const MAX_RETRIES = 20;
const DELAY_MS = 3000;

const dbUrl = process.env.DATABASE_URL;

if (!dbUrl) {
  console.error('❌ DATABASE_URL が未設定です');
  process.exit(1);
}

const waitForPostgres = async () => {
  for (let i = 1; i <= MAX_RETRIES; i++) {
    try {
      console.log(`⏳ PostgreSQL 接続試行中... (${i}/${MAX_RETRIES})`);
      const client = new Client({ connectionString: dbUrl });
      await client.connect();
      await client.query('SELECT 1'); // 簡易クエリで readiness を確認
      await client.end();

      console.log('✅ PostgreSQL 接続成功');
      return;
    } catch (err) {
      console.log(`❌ 接続失敗 (${i}/${MAX_RETRIES}): ${(err as Error).message}`);
      await new Promise((res) => setTimeout(res, DELAY_MS));
    }
  }

  console.error('💥 PostgreSQL が接続可能になる前にタイムアウトしました');
  process.exit(1);
};

waitForPostgres();