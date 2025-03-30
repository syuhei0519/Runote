import request from 'supertest';
import { PrismaClient } from '@prisma/client';
import app from '../src/app';

const prisma = new PrismaClient();

beforeEach(async () => {
  await prisma.post.deleteMany(); // 初期化
});

afterAll(async () => {
  await prisma.$disconnect();
});

describe('📬 Post API（共通データでのCRUDテスト）', () => {
    let createdPostId: number;
  
    // 各テスト前に投稿を作成
    beforeEach(async () => {
      await prisma.post.deleteMany(); // 前テストの投稿を削除
      const post = await prisma.post.create({
        data: {
          title: '初期タイトル',
          content: '初期コンテンツ',
        },
      });
      createdPostId = post.id;
    });
  
    afterAll(async () => {
      await prisma.$disconnect();
    });
  
    it('POST /posts - 新しい投稿を作成できる', async () => {
      const res = await request(app)
        .post('/posts')
        .send({ title: 'テスト投稿', content: 'これはテストです' });
  
      expect(res.status).toBe(201);
      expect(res.body).toHaveProperty('id');
      expect(res.body.title).toBe('テスト投稿');
    });
  
    it('GET /posts/:id - 投稿を取得できる', async () => {
      const res = await request(app).get(`/posts/${createdPostId}`);
      expect(res.status).toBe(200);
      expect(res.body.id).toBe(createdPostId);
      expect(res.body.title).toBe('初期タイトル');
    });
  
    it('PUT /posts/:id - 投稿を更新できる', async () => {
      const res = await request(app)
        .put(`/posts/${createdPostId}`)
        .send({ title: '更新後タイトル', content: '更新後本文' });
  
      expect(res.status).toBe(200);
      expect(res.body.title).toBe('更新後タイトル');
      expect(res.body.content).toBe('更新後本文');
    });
  
    it('DELETE /posts/:id - 投稿を削除できる', async () => {
      const delRes = await request(app).delete(`/posts/${createdPostId}`);
      expect(delRes.status).toBe(204);
  
      const getRes = await request(app).get(`/posts/${createdPostId}`);
      expect(getRes.status).toBe(404);
    });
  });

describe('🚨 異常系テスト', () => {
    it('GET /posts/:id - 存在しない ID を取得すると 404', async () => {
      const res = await request(app).get('/posts/99999');
      expect(res.status).toBe(404);
      expect(res.body.error).toBe('Post not found');
    });
  
    it('PUT /posts/:id - 存在しない ID を更新すると 404', async () => {
      const res = await request(app).put('/posts/99999').send({
        title: 'タイトル',
        content: '本文',
      });
      expect(res.status).toBe(404);
      expect(res.body.error).toBe('Post not found');
    });
  
    it('PUT /posts/:id - リクエストボディが無効だと 400', async () => {
      const create = await prisma.post.create({
        data: { title: '元', content: '本文' },
      });
  
      const res = await request(app).put(`/posts/${create.id}`).send({});
      expect(res.status).toBe(400);
      expect(res.body.errors).toBeDefined(); // zodのバリデーション結果など
    });
  
    it('DELETE /posts/:id - 存在しない ID を削除すると 404', async () => {
      const res = await request(app).delete('/posts/99999');
      expect(res.status).toBe(404);
      expect(res.body.error).toBe('Post not found');
    });
  
    it('DELETE /posts/:id - 不正な ID を渡すと 400', async () => {
      const res = await request(app).delete('/posts/abc');
      expect(res.status).toBe(400);
      expect(res.body.error).toBe('Invalid ID');
    });
  });  