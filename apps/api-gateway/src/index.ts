import express from 'express';
import dotenv from 'dotenv';
import postsRoutes from './routes/posts';
import authRoutes from './routes/auth';
import emotionsRoutes from './routes/emotions';
import tagsRoutes from './routes/tags';

// .env を読み込む
let envFile = 'dev.env';

if (process.env.NODE_ENV === 'test') {
  envFile = 'test.env';
} else if (process.env.NODE_ENV === 'production') {
  envFile = 'prd.env';
}

dotenv.config({ path: envFile });
console.log(`✅ Loaded env file: ${envFile}`);

const app = express();
const PORT = process.env.PORT || 3000;

app.use(express.json());

// 各ルートに振り分け（/api/* は Nginx が除去済み）
app.use('/posts', postsRoutes);
app.use('/auth', authRoutes);
app.use('/emotions', emotionsRoutes);
app.use('/tags', tagsRoutes);

// 健康チェック（任意）
app.get('/', (_req, res) => {
  res.send('Runote API Gateway is running!');
});

app.listen(PORT, () => {
  console.log(`API Gateway running at http://localhost:${PORT}`);
});