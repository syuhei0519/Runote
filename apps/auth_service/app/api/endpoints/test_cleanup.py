# routers/test_cleanup.py
from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from db import get_db
from models import User  # 必要なモデルをインポート

router = APIRouter()

@router.post("/test/cleanup")
def cleanup_test_data(db: Session = Depends(get_db)):
    if not is_test_env():
        return {"error": "Forbidden in non-test env"}

    db.query(User).delete()
    db.commit()
    return {"message": "Cleaned up"}