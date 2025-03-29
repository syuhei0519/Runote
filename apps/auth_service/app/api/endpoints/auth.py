from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from app.core.security import verify_token, create_access_token, verify_password, hash_password
from datetime import timedelta
from app.deps import get_db, get_current_user
from sqlalchemy.orm import Session
from app.models.user import User

router = APIRouter()

class RegisterRequest(BaseModel):
    username: str
    password: str

class LoginRequest(BaseModel):
    username: str
    password: str

# 仮のインメモリユーザーDB（後でDBに置き換え）
# fake_users = {}

@router.get("/me")
def read_current_user(current_user: User = Depends(get_current_user)):
    return {
        "user_id": current_user.id,
        "username": current_user.username,
    }

@router.post("/register")
def register(req: RegisterRequest, db: Session = Depends(get_db)):
    # 既存ユーザー確認
    existing_user = db.query(User).filter(User.username == req.username).first()
    if existing_user:
        raise HTTPException(status_code=400, detail="Username already registered")

    # パスワードハッシュ化
    hashed_pw = hash_password(req.password)

    # ユーザー作成
    user = User(username=req.username, hashed_password=hashed_pw)
    db.add(user)
    db.commit()
    db.refresh(user)

    return {"message": "registered", "user_id": user.id}

@router.post("/login")
def login(req: LoginRequest, db: Session = Depends(get_db)):
    user = db.query(User).filter(User.username == req.username).first()
    if not user or not verify_password(req.password, user.hashed_password):
        raise HTTPException(status_code=401, detail="Invalid credentials")

    access_token = create_access_token(data={"sub": user.username})
    return {"access_token": access_token, "token_type": "bearer"}

@router.post("/logout")
def logout(current_user: User = Depends(get_current_user)):
    # 実際には何もせず「ログアウトした」と返す
    return {"message": "Logged out"}