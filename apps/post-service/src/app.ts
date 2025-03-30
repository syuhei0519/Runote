import express from 'express';
import postRoutes from './routes/post';

const app = express();
app.use(express.json());
app.use('/posts', postRoutes);

app.get('/health', (req, res) => {
    res.json({ status: 'ok' });
  });

export default app;