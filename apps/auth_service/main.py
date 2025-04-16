import os
from fastapi import FastAPI
from app.api import router  # app/api/__init__.py に router が定義されている想定
from app.db import Base, engine
from app.api.endpoints import test_cleanup

app = FastAPI(
    title="Runote Auth Service",
    version="0.1.0"
)

app.include_router(router)

if os.getenv("ENV") == "test":
    app.include_router(test_cleanup.router, prefix="/auth")

Base.metadata.create_all(bind=engine)

@app.get("/health")
def health_check():
    return {"status": "ok"}