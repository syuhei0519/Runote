# routers/test_cleanup.py
from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from app.db import get_db
from app.models.user import User
import os

router = APIRouter()

@router.post("/test/cleanup")
def cleanup_test_data(db: Session = Depends(get_db)):
    print("🧪 cleanup called")
    if os.getenv("ENV") != "test":
        raise HTTPException(status_code=403, detail="Forbidden in non-test env")

    db.query(User).delete()
    db.commit()
    return {"message": "Cleaned up"}