import express from 'express';
import tagRoutes from './routes/tag';
import testRoutes from './routes/test';

const app = express();
app.use(express.json());
app.use('/tags', tagRoutes);

app.get('/health', (req, res) => {
    res.json({ status: 'ok' });
  });

if (process.env.NODE_ENV === 'test') {
  app.use('/test', testRoutes); // ← ここで有効化
}

export default app;