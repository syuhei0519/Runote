from sqlalchemy.orm import Session
from app.db import SessionLocal  # SessionLocal は sessionmaker()
from fastapi import Depends, HTTPException, status, Header
from jose import JWTError, jwt
from app.models.user import User
from app.core.config import SECRET_KEY, ALGORITHM
import os

def get_db() -> Session:
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

def get_authorization_header(authorization: str = Header(None, alias="Authorization")):
    if not authorization:
        raise HTTPException(status_code=401, detail="Missing Authorization Header")
    return authorization    

def extract_token_from_header(authorization: str = Depends(get_authorization_header)) -> str:
    if not authorization.startswith("Bearer "):
        raise HTTPException(status_code=401, detail="Invalid token format")
    return authorization[7:]

def get_current_user(token: str = Depends(extract_token_from_header),
     db: Session = Depends(get_db)) -> User:
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Could not validate credentials",
    )

    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        username: str = payload.get("sub")
        if username is None:
            raise credentials_exception
    except JWTError:
        raise credentials_exception

    user = db.query(User).filter(User.username == username).first()
    if user is None:
        raise credentials_exception
    return user