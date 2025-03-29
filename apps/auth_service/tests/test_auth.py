import pytest
from httpx import AsyncClient, ASGITransport
from main import app

transport = ASGITransport(app=app)

@pytest.mark.asyncio
async def test_register_login_and_me():
    async with AsyncClient(transport=transport, base_url="http://test") as ac:
        # ✅ 正常登録
        res = await ac.post("/auth/register", json={"username": "alice", "password": "secret"})
        assert res.status_code == 200
        data = res.json()
        assert "user_id" in data

        # ❌ 同じユーザー名で再登録（重複エラー）
        res = await ac.post("/auth/register", json={"username": "alice", "password": "secret"})
        assert res.status_code == 400
        assert res.json()["detail"] == "Username already registered"
        
        # ログイン成功
        res = await ac.post("/auth/login", json={"username": "alice", "password": "secret"})
        assert res.status_code == 200
        token = res.json()["access_token"]
        assert token

        # ログイン失敗
        res = await ac.post("/auth/login", json={"username": "alice", "password": "wrong"})
        assert res.status_code == 401

        # 1. 登録
        res = await ac.post("/auth/register", json={"username": "bob", "password": "secret"})
        assert res.status_code == 200
        user_id = res.json()["user_id"]
        assert user_id

        # 2. ログイン
        res = await ac.post("/auth/login", json={"username": "bob", "password": "secret"})
        assert res.status_code == 200
        token = res.json()["access_token"]
        assert token

        # 3. /me にトークン付きでアクセス
        headers = {"Authorization": f"Bearer {token}"}
        res = await ac.get("/auth/me", headers=headers)
        assert res.status_code == 200
        data = res.json()
        assert data["user_id"] == user_id
        assert data["username"] == "bob"

        # 4. /meにトークンなしでアクセス
        res = await ac.get("/auth/me")
        assert res.status_code == 401
        assert res.json()["detail"] in ["Missing Authorization Header", "Could not validate credentials"]

        # 登録 → ログイン
        await ac.post("/auth/register", json={"username": "logoutuser", "password": "secret"})
        res = await ac.post("/auth/login", json={"username": "logoutuser", "password": "secret"})
        token = res.json()["access_token"]

        # ログアウト
        headers = {"Authorization": f"Bearer {token}"}
        res = await ac.post("/auth/logout", headers=headers)
        assert res.status_code == 200
        assert res.json() == {"message": "Logged out"}