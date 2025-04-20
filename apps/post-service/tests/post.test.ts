import request from 'supertest';
import { PrismaClient } from '@prisma/client';
import app from '../src/app';
import axios from 'axios';
jest.mock('axios');

const prisma = new PrismaClient();
let testUserId: number;

beforeAll(async () => {
  // ユーザー作成（ID 1）
  const user = await prisma.user.upsert({
    where: { id: 1 },
    update: {},
    create: { id: 1, name: 'testuser' },
  });
  testUserId = user.id;

  // タグ作成
  await prisma.tag.createMany({
    data: [
      { id: 1, name: 'happy' },
      { id: 2, name: 'tired' },
    ],
    skipDuplicates: true,
  });
});

beforeEach(async () => {
  await prisma.postTag.deleteMany();
  await prisma.post.deleteMany();
  jest.clearAllMocks();
});

afterAll(async () => {
  await prisma.$disconnect();
});

describe('📬 Post API（共通データでのCRUDテスト）', () => {
  let createdPostId: number;

  beforeEach(async () => {
    const post = await prisma.post.create({
      data: {
        title: '初期タイトル',
        content: '初期コンテンツ',
        userId: testUserId,
      },
    });
    createdPostId = post.id;
  });

  it('POST /posts - 新しい投稿を emotion 付きで作成できる', async () => {
    (axios.post as jest.Mock).mockResolvedValue({ status: 201 });

    const res = await request(app)
      .post('/posts')
      .set('X-User-Id', `${testUserId}`)
      .send({
        title: '感情付き投稿',
        content: 'これは感情テストです',
        tagIds: [1, 2],
        emotion: {
          type: 'custom',
          label: '超エモい',
          intensity: 5,
        },
      });

    expect(res.status).toBe(201);
    expect(res.body).toHaveProperty('id');
    expect(res.body.title).toBe('感情付き投稿');
    expect(axios.post).toHaveBeenCalledWith(
      expect.stringContaining('/emotions'),
      expect.objectContaining({ label: '超エモい', intensity: 5 })
    );
  });

  it('POST /posts - emotionなしでも投稿可能', async () => {
    const res = await request(app)
      .post('/posts')
      .set('X-User-Id', `${testUserId}`)
      .send({
        title: 'シンプル投稿',
        content: '感情なし',
        tagIds: [1],
      });

    expect(res.status).toBe(201);
    expect(res.body.title).toBe('シンプル投稿');
    expect(axios.post).not.toHaveBeenCalled();
  });

  it('POST /posts - emotion-service が失敗しても投稿は成功する', async () => {
    (axios.post as jest.Mock).mockRejectedValue(new Error('emotion-service unreachable'));

    const res = await request(app)
      .post('/posts')
      .set('X-User-Id', `${testUserId}`)
      .send({
        title: 'エラー時もOK',
        content: 'emotionは失敗',
        emotion: {
          type: 'preset',
          label: 'tired',
          presetKey: 'tired',
          intensity: 2,
        },
      });

    expect(res.status).toBe(201);
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
      .set('X-User-Id', `${testUserId}`)
      .send({ title: '更新後タイトル', content: '更新後本文' });

    expect(res.status).toBe(200);
    expect(res.body.title).toBe('更新後タイトル');
    expect(res.body.content).toBe('更新後本文');
  });

  it('DELETE /posts/:id - 投稿を削除できる', async () => {
    const delRes = await request(app)
      .delete(`/posts/${createdPostId}`)
      .set('X-User-Id', `${testUserId}`);
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
    const res = await request(app)
      .put('/posts/99999')
      .set('X-User-Id', `${testUserId}`)
      .send({
        title: 'タイトル',
        content: '本文',
      });
    expect(res.status).toBe(404);
  });

  it('PUT /posts/:id - リクエストボディが無効だと 400', async () => {
    const create = await prisma.post.create({
      data: { title: '元', content: '本文', userId: testUserId },
    });

    const res = await request(app)
      .put(`/posts/${create.id}`)
      .set('X-User-Id', `${testUserId}`)
      .send({});
    expect(res.status).toBe(400);
    expect(res.body.errors).toBeDefined();
  });

  it('DELETE /posts/:id - 存在しない ID を削除すると 404', async () => {
    const res = await request(app)
      .delete('/posts/99999')
      .set('X-User-Id', `${testUserId}`);
    expect(res.status).toBe(404);
  });

  it('DELETE /posts/:id - 不正な ID を渡すと 400', async () => {
    const res = await request(app)
      .delete('/posts/abc')
      .set('X-User-Id', `${testUserId}`);
    expect(res.status).toBe(400);
  });

  it('PUT /posts/:id - 他人の投稿は更新できない', async () => {
    const anotherUser = await prisma.user.upsert({
      where: { id: 999 },
      update: {},
      create: { id: 999, name: 'anotheruser' },
    });

    const post = await prisma.post.create({
      data: {
        title: '他人の投稿',
        content: '他人の本文',
        userId: anotherUser.id,
      },
    });

    const res = await request(app)
      .put(`/posts/${post.id}`)
      .set('X-User-Id', `${testUserId}`)
      .send({ title: '不正更新', content: 'これは許されない' });

    expect(res.status).toBe(403);
    expect(res.body.error).toBe('forbidden');
  });

  it('DELETE /posts/:id - 他人の投稿は削除できない', async () => {
    const post = await prisma.post.create({
      data: {
        title: '他人の投稿',
        content: '削除できない投稿',
        userId: 999,
      },
    });

    const res = await request(app)
      .delete(`/posts/${post.id}`)
      .set('X-User-Id', `${testUserId}`);
    expect(res.status).toBe(403);
    expect(res.body.error).toBe('forbidden');
  });
});

describe('中間テーブルへの反映テスト', () => {
  it('POST /posts - tagIds に応じて PostTag が作成される', async () => {
    const res = await request(app)
      .post('/posts')
      .set('X-User-Id', `${testUserId}`)
      .send({
        title: 'タグ付き投稿',
        content: 'これはタグ付きです',
        tagIds: [1, 2],
      });

    expect(res.status).toBe(201);
    const postId = res.body.id;

    const tagLinks = await prisma.postTag.findMany({ where: { postId } });
    const tagIds = tagLinks.map(t => t.tagId).sort();
    expect(tagIds).toEqual([1, 2]);
  });
});
