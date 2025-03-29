from fastapi import FastAPI
from app.api import router  # app/api/__init__.py に router が定義されている想定

app = FastAPI(
    title="Runote Auth Service",
    version="0.1.0"
)

app.include_router(router)

@app.get("/health")
def health_check():
    return {"status": "ok"}