from dotenv import load_dotenv
load_dotenv()

import os

SECRET_KEY = os.getenv("JWT_SECRET", "defaultsecret")
ALGORITHM = os.getenv("JWT_ALGORITHM", "HS256")
ACCESS_TOKEN_EXPIRE_MINUTES = int(os.getenv("JWT_EXPIRE_MINUTES", "30"))