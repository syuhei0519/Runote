import express from 'express';
import postRoutes from './routes/post';
import testRoutes from './routes/test';

const app = express();
app.use(express.json());
app.use('/posts', postRoutes);

app.get('/health', (req, res) => {
    res.json({ status: 'ok' });
  });

if (process.env.NODE_ENV === 'test') {
    app.use('/test', testRoutes);
}  
export default app;