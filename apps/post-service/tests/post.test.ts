import request from 'supertest';
import { PrismaClient } from '@prisma/client';
import app from '../src/app';
import axios from 'axios';
jest.mock('axios'); 

const prisma = new PrismaClient();

beforeAll(async () => {
  await prisma.tag.createMany({
    data: [
      { id: 1, name: 'happy' },
      { id: 2, name: 'tired' }
    ],
    skipDuplicates: true
  });
});

// 各テスト前に投稿を作成
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
      await prisma.postTag.deleteMany();
      await prisma.post.deleteMany();
      jest.clearAllMocks();
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
  
    it('POST /posts - 新しい投稿を emotion 付きで作成できる', async () => {
      (axios.post as jest.Mock).mockResolvedValue({ status: 201 });
  
      const res = await request(app)
        .post('/posts')
        .send({
          title: '感情付き投稿',
          content: 'これは感情テストです',
          tagIds: [1, 2],
          emotion: {
            type: 'custom',
            label: '超エモい',
            intensity: 5,
          }
        });
  
      expect(res.status).toBe(201);
      expect(res.body).toHaveProperty('id');
      expect(res.body.title).toBe('感情付き投稿');
      expect(axios.post).toHaveBeenCalledWith( // ★ emotion-service 呼び出し確認
        expect.stringContaining('/emotions'),
        expect.objectContaining({
          label: '超エモい',
          intensity: 5,
        })
      );
    });
  
    it('POST /posts - emotionなしでも投稿可能', async () => {
      const res = await request(app)
        .post('/posts')
        .send({
          title: 'シンプル投稿',
          content: '感情なし',
          tagIds: [1]
        });
  
      expect(res.status).toBe(201);
      expect(res.body.title).toBe('シンプル投稿');
      expect(axios.post).not.toHaveBeenCalled(); // ★ emotion-serviceは呼ばれない
    });

    it('POST /posts - emotion-service が失敗しても投稿は成功する', async () => {
      (axios.post as jest.Mock).mockRejectedValue(new Error('emotion-service unreachable'));
    
      const res = await request(app)
        .post('/posts')
        .send({
          title: 'エラー時もOK',
          content: 'emotionは失敗',
          emotion: {
            type: 'preset',
            label: 'tired',
            presetKey: 'tired',
            intensity: 2
          }
        });
    
      expect(res.status).toBe(201); // ✅ 投稿自体は成功
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
  }
);

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
  }
);

describe('中間テーブルへの反映テスト', () => {
    it('POST /posts - tagIds に応じて PostTag が作成される', async () => {
      const res = await request(app)
        .post('/posts')
        .send({
          title: 'タグ付き投稿',
          content: 'これはタグ付きです',
          tagIds: [1, 2]
        });
    
      expect(res.status).toBe(201);
      const postId = res.body.id;
    
      // PostTag に正しく登録されているか確認
      const tagLinks = await prisma.postTag.findMany({ where: { postId } });
      const tagIds = tagLinks.map(t => t.tagId).sort();
      expect(tagIds).toEqual([1, 2]);
    });
  }
);