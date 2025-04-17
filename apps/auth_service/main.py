import os
import time
from sqlalchemy.exc import OperationalError
from fastapi import FastAPI
from app.api import router  # app/api/__init__.py に router が定義されている想定
from app.db import Base, engine
from app.api.endpoints import test_cleanup
from sqlalchemy import text

app = FastAPI(
    title="Runote Auth Service",
    version="0.1.0"
)

def wait_for_mysql(max_retries=10, delay=3):
    for i in range(max_retries):
        try:
            with engine.connect() as conn:
                conn.execute(text("SELECT 1"))
                print("✅ MySQL ready")
                return
        except OperationalError as e:
            print(f"❌ Retry {i+1}/{max_retries} - MySQL not ready yet: {e}")
            time.sleep(delay)
    raise RuntimeError("❌ MySQL did not become available in time")

app.include_router(router)

if os.getenv("ENV") == "test":
    app.include_router(test_cleanup.router, prefix="/auth")

wait_for_mysql()

Base.metadata.create_all(bind=engine)

@app.get("/health")
def health_check():
    return {"status": "ok"}