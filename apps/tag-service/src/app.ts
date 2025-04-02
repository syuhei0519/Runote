import express from 'express';
import tagRoutes from './routes/tag';

const app = express();
app.use(express.json());
app.use('/tags', tagRoutes);

app.get('/health', (req, res) => {
    res.json({ status: 'ok' });
  });

export default app;