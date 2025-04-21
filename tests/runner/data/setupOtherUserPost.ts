// デフォルト以外のユーザーと投稿を作成
import { postClient } from './clients/post';
import { userClient } from './clients/auth';
import bcrypt from 'bcrypt';

export default async function setupOtherUserPost() {
  // ハッシュ済みパスワード（プレーン: "password123"）
  const passwordHash = await bcrypt.hash('password123', 10);

  // ユーザーID=999 を作成（存在しなければ）
  await userClient.query(`
    INSERT INTO "user" (id, name, password)
    VALUES (999, 'otheruser', $1)
    ON CONFLICT (id) DO NOTHING;
  `, [passwordHash]);

  // 投稿ID=2 を作成（存在しなければ）
  await postClient.query(`
    INSERT INTO post (id, title, content, "userId")
    VALUES (2, '他人の投稿', 'これは他人の投稿です', 999)
    ON CONFLICT (id) DO NOTHING;
  `);
}