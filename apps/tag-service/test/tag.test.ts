import request from 'supertest';
import app from '../src/app';
import { PrismaClient } from '@prisma/client';
import { execSync } from 'child_process';

const prisma = new PrismaClient();

// テスト用 DB のマイグレーションとクライアント生成
beforeAll(() => {
    execSync('npx prisma migrate deploy --schema=./prisma/schema.prisma', {
      stdio: 'inherit',
    });
});

beforeEach(async () => {
  await prisma.tag.deleteMany(); // クリーンな状態でテスト
});

afterAll(async () => {
  await prisma.$disconnect();
});

describe('タグAPI', () => {
  it('POST /tags - タグを作成できる', async () => {
    const res = await request(app)
      .post('/tags')
      .send({ name: '朝ラン' });

    expect(res.statusCode).toBe(201);
    expect(res.body).toHaveProperty('id');
    expect(res.body.name).toBe('朝ラン');
  });

  it('GET /tags - タグ一覧が取得できる', async () => {
    await prisma.tag.create({ data: { name: '夜ラン' } });

    const res = await request(app).get('/tags');

    expect(res.statusCode).toBe(200);
    expect(Array.isArray(res.body)).toBe(true);
    expect(res.body[0].name).toBe('夜ラン');
  });

  it('PUT /tags/:id - タグを更新できる', async () => {
    // まずタグを1つ作成
    const tag = await prisma.tag.create({ data: { name: '旧タグ名' } });
  
    // 更新リクエスト送信
    const res = await request(app)
      .put(`/tags/${tag.id}`)
      .send({ name: '新タグ名' });
  
    expect(res.statusCode).toBe(200);
    expect(res.body).toHaveProperty('id', tag.id);
    expect(res.body.name).toBe('新タグ名');
  });
  
  it('PUT /tags/:id - 無効な名前はバリデーションエラーになる', async () => {
    const tag = await prisma.tag.create({ data: { name: 'バリデート対象' } });
  
    const res = await request(app)
      .put(`/tags/${tag.id}`)
      .send({ name: '' }); // 無効なデータ
  
    expect(res.statusCode).toBe(400);
    expect(res.body).toHaveProperty('errors');
    expect(Array.isArray(res.body.errors)).toBe(true);
  });  

  it('DELETE /tags/:id - タグを削除できる', async () => {
    const tag = await prisma.tag.create({ data: { name: '削除タグ' } });

    const res = await request(app).delete(`/tags/${tag.id}`);

    expect(res.statusCode).toBe(204);
  });
});