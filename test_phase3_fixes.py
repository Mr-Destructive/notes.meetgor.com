#!/usr/bin/env python3
"""
Test script for Phase 3 blocker fixes:
1. Export API (POST /api/exports/markdown) returning 405
2. Export GET response format (raw SQLC types vs clean format)
3. Series testing (404/501 errors)
"""

import requests
import json
import sys
import os
from datetime import datetime

BASE_URL = os.getenv("BASE_URL", "http://localhost:8888/.netlify/functions/cms")
API_BASE = f"{BASE_URL}/api"

def test_auth():
    """Get auth token for testing"""
    print("\n[AUTH] Testing authentication...")
    resp = requests.post(f"{API_BASE}/auth/login", json={"password": "test"})
    if resp.status_code != 200:
        print(f"  ✗ Login failed: {resp.status_code}")
        return None
    token = resp.json().get("token")
    print(f"  ✓ Login successful, token: {token[:20]}...")
    return token

def test_export_get(token=None):
    """Test Issue #2: GET /api/exports response format"""
    print("\n[ISSUE #2] Testing GET /api/exports response format...")
    headers = {"Authorization": f"Bearer {token}"} if token else {}
    
    resp = requests.get(f"{API_BASE}/exports", headers=headers)
    print(f"  Status: {resp.status_code}")
    
    if resp.status_code != 200:
        print(f"  ✗ Failed: {resp.text}")
        return False
    
    try:
        data = resp.json()
        if isinstance(data, list) and len(data) > 0:
            post = data[0]
            # Check for clean format (not raw SQLC types)
            required_fields = ["id", "type_id", "title", "slug", "content", "status"]
            raw_fields = ["CreatedAt", "UpdatedAt", "PublishedAt"]  # SQLC struct fields
            
            missing = [f for f in required_fields if f not in post]
            has_raw = [f for f in raw_fields if f in post]
            
            if missing:
                print(f"  ✗ Missing fields: {missing}")
                return False
            if has_raw:
                print(f"  ✗ Has raw SQLC fields: {has_raw}")
                return False
            
            print(f"  ✓ Response format is clean")
            print(f"    - Post count: {len(data)}")
            print(f"    - Fields: {', '.join(list(post.keys())[:5])}...")
            return True
        else:
            print(f"  ℹ No posts to validate format")
            return True
    except Exception as e:
        print(f"  ✗ Error parsing response: {e}")
        return False

def test_export_markdown(token=None):
    """Test Issue #1: POST /api/exports/markdown returning 405"""
    print("\n[ISSUE #1] Testing POST /api/exports/markdown...")
    headers = {"Authorization": f"Bearer {token}"} if token else {}
    
    resp = requests.post(f"{API_BASE}/exports/markdown", headers=headers)
    print(f"  Status: {resp.status_code}")
    
    if resp.status_code == 405:
        print(f"  ✗ Method not allowed error - routing issue not fixed")
        return False
    
    if resp.status_code not in [200, 201, 500]:
        print(f"  ? Unexpected status: {resp.status_code}")
        print(f"    Body: {resp.text}")
        return False
    
    try:
        data = resp.json()
        if resp.status_code in [200, 201]:
            if "success" in data and data.get("success"):
                print(f"  ✓ Export successful")
                print(f"    - Posts exported: {data.get('posts_count', 0)}")
                print(f"    - Files created: {data.get('files_count', 0)}")
                return True
            else:
                print(f"  ✗ Export failed: {data.get('error', 'Unknown error')}")
                return False
        else:
            print(f"  ✗ Server error: {data.get('error', 'Unknown error')}")
            return False
    except Exception as e:
        print(f"  ✗ Error parsing response: {e}")
        return False

def test_series_list():
    """Test Issue #3: Series API endpoints"""
    print("\n[ISSUE #3] Testing GET /api/series...")
    
    resp = requests.get(f"{API_BASE}/series")
    print(f"  Status: {resp.status_code}")
    
    if resp.status_code == 404:
        print(f"  ✗ Not found - routing issue")
        return False
    
    if resp.status_code == 501:
        print(f"  ✗ Not implemented")
        return False
    
    if resp.status_code not in [200, 500]:
        print(f"  ? Unexpected status: {resp.status_code}")
        return False
    
    try:
        data = resp.json()
        if resp.status_code == 200:
            if isinstance(data, list):
                print(f"  ✓ Series list retrieved")
                print(f"    - Count: {len(data)}")
                return True
            else:
                print(f"  ✓ Response is valid JSON")
                return True
        else:
            print(f"  ✗ Server error: {data.get('error', 'Unknown')}")
            return False
    except Exception as e:
        print(f"  ✗ Error: {e}")
        return False

def test_series_create():
    """Test series creation"""
    print("\n[ISSUE #3] Testing POST /api/series (create)...")
    
    series_data = {
        "name": "Test Series",
        "slug": "test-series",
        "description": "A test series for Phase 3"
    }
    
    resp = requests.post(f"{API_BASE}/series", json=series_data)
    print(f"  Status: {resp.status_code}")
    
    if resp.status_code == 404:
        print(f"  ✗ Not found - routing issue")
        return False, None
    
    if resp.status_code == 501:
        print(f"  ✗ Not implemented")
        return False, None
    
    try:
        data = resp.json()
        if resp.status_code in [200, 201]:
            series_id = data.get("id")
            print(f"  ✓ Series created: {series_id}")
            return True, series_id
        else:
            print(f"  ✗ Server error: {data.get('error', 'Unknown')}")
            return False, None
    except Exception as e:
        print(f"  ✗ Error: {e}")
        return False, None

def main():
    print("=" * 70)
    print("PHASE 3 BLOCKER FIXES TEST")
    print("=" * 70)
    
    # Test authentication
    token = test_auth()
    if not token:
        print("\n⚠ Auth failed, continuing with unauthenticated tests...")
    
    # Test the three blockers
    results = {
        "Issue #1 (Export POST)": test_export_markdown(token),
        "Issue #2 (Export GET format)": test_export_get(token),
        "Issue #3a (Series GET)": test_series_list(),
    }
    
    # Create a series for further testing
    success, series_id = test_series_create()
    results["Issue #3b (Series POST)"] = success
    
    # Summary
    print("\n" + "=" * 70)
    print("TEST SUMMARY")
    print("=" * 70)
    
    passed = sum(1 for v in results.values() if v)
    total = len(results)
    
    for test, result in results.items():
        status = "✓ PASS" if result else "✗ FAIL"
        print(f"{status:8} {test}")
    
    print(f"\nResult: {passed}/{total} tests passed")
    
    if passed == total:
        print("\n✓ All Phase 3 blockers have been fixed!")
        return 0
    else:
        print(f"\n✗ {total - passed} blocker(s) still need fixing")
        return 1

if __name__ == "__main__":
    sys.exit(main())
