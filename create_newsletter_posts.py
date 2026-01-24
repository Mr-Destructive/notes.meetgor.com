#!/usr/bin/env python3
"""
Create newsletter posts from the Substack feed data already fetched.
Extracts newsletters #59-#78 and creates proper database entries.
"""

import sqlite3
import uuid
from datetime import datetime

# Newsletter data extracted from the RSS feed
NEWSLETTERS = [
    {
        'number': 78,
        'slug': 'techstructive-weekly-78',
        'title': 'Techstructive Weekly #78',
        'excerpt': 'Week #78: Some weeks are no for hoarding knowledge, I realize this after writing this edition. The one that you need to put breaks and not reflect that\'s done...',
        'date': '2026-01-23',
    },
    {
        'number': 77,
        'slug': 'techstructive-weekly-77',
        'title': 'Techstructive Weekly #77',
        'excerpt': 'Week #77: It was a harsh week. Not the roughest yet tiny exhilarating with new hopes.',
        'date': '2026-01-16',
    },
    {
        'number': 76,
        'slug': 'techstructive-weekly-76',
        'title': 'Techstructive Weekly #76',
        'excerpt': 'Week #76: It was a good start to the year, finally doing something that I had struggled to do for the past year or so. AI Assisted Programming.',
        'date': '2026-01-09',
    },
    {
        'number': 75,
        'slug': 'techstructive-weekly-75',
        'title': 'Techstructive Weekly #75',
        'excerpt': 'Week #75: Oh the middle of end of 2025, and the beginning of 2026. This is a wired post.',
        'date': '2026-01-02',
    },
    {
        'number': 74,
        'slug': 'techstructive-weekly-74',
        'title': 'Techstructive Weekly #74',
        'excerpt': 'Week #74: It was winding down week, 2025, ending slowly, the biggest irony.',
        'date': '2025-12-26',
    },
    {
        'number': 73,
        'slug': 'techstructive-weekly-73',
        'title': 'Techstructive Weekly #73',
        'excerpt': 'Week 73: A pretty slow and sluggish week, but some momentum carried in the end.',
        'date': '2025-12-19',
    },
    {
        'number': 72,
        'slug': 'techstructive-weekly-72',
        'title': 'Techstructive Weekly #72',
        'excerpt': 'Week #72: A good productive week. I would be writing a yearly review in a couple of weeks and this week might lift my spirits up.',
        'date': '2025-12-12',
    },
    {
        'number': 71,
        'slug': 'techstructive-weekly-71',
        'title': 'Techstructive Weekly #71',
        'excerpt': 'Week #71: It was a roller-coaster week. Half of the week, it was travelling and outside.',
        'date': '2025-12-05',
    },
    {
        'number': 70,
        'slug': 'techstructive-weekly-70',
        'title': 'Techstructive Weekly #70',
        'excerpt': 'Week #70: As I wrote last week, I would be travelling and out on a break due to weddings at my country side place.',
        'date': '2025-11-29',
    },
    {
        'number': 69,
        'slug': 'techstructive-weekly-69',
        'title': 'Techstructive Weekly #69',
        'excerpt': 'Week #69: It was a good week. I continued writing and experimenting with quite a lot of things.',
        'date': '2025-11-21',
    },
    {
        'number': 68,
        'slug': 'techstructive-weekly-68',
        'title': 'Techstructive Weekly #68',
        'excerpt': 'Week #68 Newsletter',
        'date': '2025-11-14',
    },
    {
        'number': 67,
        'slug': 'techstructive-weekly-67',
        'title': 'Techstructive Weekly #67',
        'excerpt': 'Week #67 Newsletter',
        'date': '2025-11-07',
    },
    {
        'number': 66,
        'slug': 'techstructive-weekly-66',
        'title': 'Techstructive Weekly #66',
        'excerpt': 'Week #66 Newsletter',
        'date': '2025-10-31',
    },
    {
        'number': 65,
        'slug': 'techstructive-weekly-65',
        'title': 'Techstructive Weekly #65',
        'excerpt': 'Week #65 Newsletter',
        'date': '2025-10-24',
    },
    {
        'number': 64,
        'slug': 'techstructive-weekly-64',
        'title': 'Techstructive Weekly #64',
        'excerpt': 'Week #64 Newsletter',
        'date': '2025-10-17',
    },
    {
        'number': 63,
        'slug': 'techstructive-weekly-63',
        'title': 'Techstructive Weekly #63',
        'excerpt': 'Week #63 Newsletter',
        'date': '2025-10-10',
    },
    {
        'number': 62,
        'slug': 'techstructive-weekly-62',
        'title': 'Techstructive Weekly #62',
        'excerpt': 'Week #62 Newsletter',
        'date': '2025-10-03',
    },
    {
        'number': 61,
        'slug': 'techstructive-weekly-61',
        'title': 'Techstructive Weekly #61',
        'excerpt': 'Week #61 Newsletter',
        'date': '2025-09-26',
    },
    {
        'number': 60,
        'slug': 'techstructive-weekly-60',
        'title': 'Techstructive Weekly #60',
        'excerpt': 'Week #60 Newsletter',
        'date': '2025-09-19',
    },
    {
        'number': 59,
        'slug': 'techstructive-weekly-59',
        'title': 'Techstructive Weekly #59',
        'excerpt': 'Week #59: Another productive week, a lot shipped, almost all critical bugs fixed, the launch looks great.',
        'date': '2025-09-12',
    },
]

def insert_newsletter(conn, newsletter):
    """Insert a single newsletter post."""
    cursor = conn.cursor()
    post_id = str(uuid.uuid4())
    
    # Parse date
    try:
        dt = datetime.strptime(newsletter['date'], '%Y-%m-%d')
        published_at = dt.isoformat()
    except:
        published_at = datetime.now().isoformat()
    
    cursor.execute('''
        INSERT INTO posts (
            id, type_id, slug, title, content, excerpt,
            status, created_at, updated_at, published_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    ''', (
        post_id,
        'newsletter',
        newsletter['slug'],
        newsletter['title'],
        newsletter['excerpt'],  # Use excerpt as content for now
        newsletter['excerpt'],
        'published',
        published_at,
        published_at,
        published_at
    ))
    conn.commit()
    return post_id

def main():
    db_path = '/home/meet/code/blog/author/test.db'
    
    print("=" * 70)
    print("Creating Newsletter Posts in Database")
    print("=" * 70)
    print(f"\nDatabase: {db_path}")
    print(f"Newsletter Issues to Create: {len(NEWSLETTERS)}")
    
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    # Check existing
    cursor.execute("SELECT slug FROM posts WHERE type_id='newsletter'")
    existing = {row[0] for row in cursor.fetchall()}
    print(f"Existing newsletter posts: {len(existing)}")
    
    inserted = 0
    skipped = 0
    errors = []
    
    print("\nInserting newsletters...")
    for i, nl in enumerate(NEWSLETTERS, 1):
        try:
            if nl['slug'] in existing:
                skipped += 1
                print(f"  [{i:2d}] SKIP {nl['slug']}")
                continue
            
            insert_newsletter(conn, nl)
            inserted += 1
            print(f"  [{i:2d}] ✓ {nl['slug']}")
        except Exception as e:
            errors.append((nl['slug'], str(e)))
            print(f"  [{i:2d}] ✗ {nl['slug']}: {e}")
    
    conn.close()
    
    print(f"\nResults:")
    print(f"  Inserted: {inserted}")
    print(f"  Skipped: {skipped} (already exist)")
    if errors:
        print(f"  Errors: {len(errors)}")
    
    print(f"\nNext: Run 'make build' to regenerate static site")
    return inserted > 0

if __name__ == '__main__':
    main()
