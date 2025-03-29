from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from app.core.security import verify_token, create_access_token, verify_password, hash_password
from datetime import timedelta

router = APIRouter()

class RegisterRequest(BaseModel):
    username: str
    password: str

class LoginRequest(BaseModel):
    username: str
    password: str

# 仮のインメモリユーザーDB（後でDBに置き換え）
fake_users = {}

@router.get("/me")
def read_me(current_user: dict = Depends(verify_token)):
    return {"user": current_user}

@router.post("/register")
def register(req: RegisterRequest):
    if req.username in fake_users:
        raise HTTPException(status_code=400, detail="User already exists")
    
    hashed = hash_password(req.password)
    fake_users[req.username] = hashed
    return {"message": "registered"}

@router.post("/login")
def login(req: LoginRequest):
    user_pw = fake_users.get(req.username)
    if not user_pw or not verify_password(req.password, user_pw):
        raise HTTPException(status_code=401, detail="Invalid credentials")
    
    access_token = create_access_token({"sub": req.username}, expires_delta=timedelta(minutes=30))
    return {"access_token": access_token, "token_type": "bearer"}