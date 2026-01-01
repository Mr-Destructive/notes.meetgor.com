from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from datetime import datetime, timedelta
import jwt
import os

app = FastAPI()

class LoginRequest(BaseModel):
    password: str

ADMIN_PASSWORD = os.getenv("ADMIN_PASSWORD", "password")
JWT_SECRET = os.getenv("JWT_SECRET", "dev-secret-key")

@app.post("")
async def login(request: LoginRequest):
    """Login endpoint"""
    if request.password != ADMIN_PASSWORD:
        raise HTTPException(status_code=401, detail="Invalid credentials")
    
    # Create JWT token
    payload = {
        "admin": True,
        "exp": datetime.utcnow() + timedelta(days=7)
    }
    token = jwt.encode(payload, JWT_SECRET, algorithm="HS256")
    
    return {
        "token": token,
        "expires_in": 604800
    }
