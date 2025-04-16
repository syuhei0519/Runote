import request from 'supertest';
import jwt from 'jsonwebtoken';

const api = request('http://nginx:8080'); // Nginx経由（Dockerネットワーク内で通信）

const testUser = {
    password: 'secret123',
    username: 'testuser'
};

let jwtToken = '';

/**
 * テスト用DBを初期化（auth-service にある /auth/cleanup を使用）
 */
async function resetDatabase() {
    const res = await api
      .post('/api/auth/test/cleanup')
      .set('Content-Type', 'application/json')
      .send({}); // ✅ 空のJSONボディを送る
    expect(res.status).toBe(200);
}

/**
 * ユーザーを登録
 */
async function registerUser() {
  const res = await api
    .post('/api/auth/register')
    .send(testUser);
  expect([200, 201]).toContain(res.status);
}

/**
 * JWTをログインAPI経由で取得
 */
async function fetchJwt(): Promise<string> {
    const testUser = { username: 'testuser', password: 'password' };
  
    // 毎回初期化
    await api.post('/api/auth/test/cleanup').send({});
    console.log('🧹 Cleanup done');
  
    const registerRes = await api.post('/api/auth/register').send(testUser);
    expect([200, 201]).toContain(registerRes.status);
  
    const loginRes = await api.post('/api/auth/login').send(testUser);
    expect(loginRes.status).toBe(200);
    expect(loginRes.body.access_token).toBeDefined();
  
    const token = `Bearer ${loginRes.body.access_token}`;
    return token;
}  

describe('🧪 Runote API Gateway E2E via Nginx', () => {
  beforeAll(async () => {
    // ステップ順：DB初期化 → ユーザー作成 → トークン取得
    await resetDatabase();
    await registerUser();
    jwtToken = await fetchJwt();
  });

  it('✅ GET /api/posts → 200 + data（要JWT）', async () => {
    const res = await api
      .get('/api/posts')
      .set('Authorization', jwtToken);

    expect(res.statusCode).toBe(200);
    expect(res.body).toBeDefined();
  });

  it('🔒 GET /api/posts → 401 when no token', async () => {
    const res = await api.get('/api/posts');
    expect(res.status).toBe(401);
  });

  it('🛑 GET /api/posts → 403 with invalid token', async () => {
    const res = await api
      .get('/api/posts')
      .set('Authorization', 'Bearer invalid.token.here');

    expect(res.status).toBe(403);
  });

  it('🌐 OPTIONS /api/posts → 204 (CORS preflight)', async () => {
    const res = await api
      .options('/api/posts')
      .set('Origin', 'http://localhost:3000')
      .set('Access-Control-Request-Method', 'POST');

    expect(res.status).toBe(204);
    expect(res.headers['access-control-allow-origin']).toBe('*');
  });

  it('📝 POST /api/posts → 201 + created（要JWT）', async () => {
    const jwtToken = await fetchJwt();
    const res = await api
      .post('/api/posts')
      .set('Authorization', jwtToken)
      .send({
        title: 'New Post',
        content: 'Test content',
        tagIds: [],
        emotion: {
          label: 'happy',
          type: 'preset',
          intensity: 3,
        },
      });
  
    expect(res.statusCode).toBe(201);
    expect(res.body).toHaveProperty('id');
    expect(res.body.title).toBe('New Post');
  });
  
  it('⏳ POST /api/posts → 401 with expired JWT', async () => {
    // 手動で期限切れトークンを作成（有効期限を過去にする）
    const expiredToken = jwt.sign(
      { sub: 'testuser', exp: Math.floor(Date.now() / 1000) - 60 }, // 60秒前
      process.env.JWT_SECRET || 'runote-dev-secret'
    );
  
    const res = await api
      .post('/api/posts')
      .set('Authorization', `Bearer ${expiredToken}`)
      .send({
        title: 'Expired Token Post',
        content: 'Should fail',
      });
  
    expect(res.status).toBe(401);
  });  
});