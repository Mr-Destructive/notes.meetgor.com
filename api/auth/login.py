import json
from datetime import datetime, timedelta
import jwt
import os

ADMIN_PASSWORD = os.getenv("ADMIN_PASSWORD", "password")
JWT_SECRET = os.getenv("JWT_SECRET", "dev-secret-key")

def handler(request):
    if request.method != "POST":
        return {
            "statusCode": 405,
            "body": json.dumps({"error": "Method not allowed"})
        }
    
    try:
        body = json.loads(request.body)
        password = body.get("password")
        
        if password != ADMIN_PASSWORD:
            return {
                "statusCode": 401,
                "body": json.dumps({"error": "Invalid credentials"})
            }
        
        # Create JWT token
        payload = {
            "admin": True,
            "exp": datetime.utcnow() + timedelta(days=7)
        }
        token = jwt.encode(payload, JWT_SECRET, algorithm="HS256")
        
        return {
            "statusCode": 200,
            "body": json.dumps({
                "token": token,
                "expires_in": 604800
            })
        }
    except Exception as e:
        return {
            "statusCode": 400,
            "body": json.dumps({"error": str(e)})
        }
