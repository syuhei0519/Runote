import yaml
from main import app
from fastapi.openapi.utils import get_openapi

schema = get_openapi(
    title=app.title,
    version=app.version,
    routes=app.routes,
)

with open("openapi.yaml", "w") as f:
    yaml.dump(schema, f, sort_keys=False, allow_unicode=True)