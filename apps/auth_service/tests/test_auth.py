import pytest
from httpx import AsyncClient, ASGITransport
from main import app

transport = ASGITransport(app=app)

@pytest.mark.asyncio
async def test_register_and_login():
    async with AsyncClient(transport=transport, base_url="http://test") as ac:
        # 新規登録
        res = await ac.post("/auth/register", json={"username": "alice", "password": "secret"})
        assert res.status_code == 200
        assert res.json() == {"message": "registered"}

        # ログイン成功
        res = await ac.post("/auth/login", json={"username": "alice", "password": "secret"})
        assert res.status_code == 200
        token = res.json()["access_token"]
        assert token

        # ログイン失敗
        res = await ac.post("/auth/login", json={"username": "alice", "password": "wrong"})
        assert res.status_code == 401