import dotenv from 'dotenv';
import app from './app';

dotenv.config(); // .env の読み込み

const port = process.env.PORT || 4000;

app.listen(port, () => {
  console.log(`✅ Server listening on http://localhost:${port}`);
})