from fastapi import FastAPI, HTTPException, Header
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel
import os
from datetime import datetime, timedelta
import jwt
from libsql_client import create_client
import pathlib

app = FastAPI()

# Mount static files from public directory
public_dir = pathlib.Path(__file__).parent.parent / "public"
if public_dir.exists():
    app.mount("/", StaticFiles(directory=str(public_dir), html=True), name="static")

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Database client
db = create_client(
    url=os.getenv("TURSO_CONNECTION_URL", "file:./dev.db"),
    auth_token=os.getenv("TURSO_AUTH_TOKEN")
)

# Config
ADMIN_PASSWORD = os.getenv("ADMIN_PASSWORD", "password")
JWT_SECRET = os.getenv("JWT_SECRET", "dev-secret-key")

# Models
class LoginRequest(BaseModel):
    password: str

class PostCreate(BaseModel):
    title: str
    content: str
    type_id: str
    slug: str = None
    excerpt: str = None
    tags: list = []
    metadata: dict = {}
    status: str = "draft"

# Auth endpoints
@app.post("/api/auth/login")
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

@app.post("/api/auth/verify")
async def verify(authorization: str = Header(None)):
    """Verify JWT token"""
    if not authorization:
        raise HTTPException(status_code=401, detail="No token")
    
    try:
        token = authorization.replace("Bearer ", "")
        jwt.decode(token, JWT_SECRET, algorithms=["HS256"])
        return {"valid": True}
    except:
        raise HTTPException(status_code=401, detail="Invalid token")

# Post endpoints
@app.get("/api/posts")
async def get_posts(status: str = None, type_id: str = None):
    """Get posts"""
    query = "SELECT * FROM posts WHERE 1=1"
    if status:
        query += f" AND status = '{status}'"
    if type_id:
        query += f" AND type_id = '{type_id}'"
    query += " ORDER BY published_at DESC"
    
    try:
        result = db.execute(query)
        return [dict(row) for row in result.rows]
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/posts/{post_id}")
async def get_post(post_id: str):
    """Get single post"""
    try:
        result = db.execute(f"SELECT * FROM posts WHERE id = '{post_id}'")
        if result.rows:
            return dict(result.rows[0])
        raise HTTPException(status_code=404, detail="Post not found")
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/posts")
async def create_post(post: PostCreate, authorization: str = Header(None)):
    """Create post"""
    if not authorization:
        raise HTTPException(status_code=401, detail="Unauthorized")
    
    try:
        import uuid
        post_id = str(uuid.uuid4())
        
        query = f"""
        INSERT INTO posts (id, type_id, title, slug, content, excerpt, status, tags, metadata, created_at, updated_at)
        VALUES ('{post_id}', '{post.type_id}', '{post.title}', '{post.slug or post.title.lower().replace(" ", "-")}', 
                '{post.content}', '{post.excerpt or ""}', '{post.status}', '{post.tags}', '{post.metadata}', 
                datetime('now'), datetime('now'))
        """
        db.execute(query)
        return {"id": post_id, "status": "created"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

# Health check
@app.get("/health")
async def health():
    return {"status": "ok"}
