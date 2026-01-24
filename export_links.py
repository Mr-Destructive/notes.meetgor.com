#!/usr/bin/env python3
import sqlite3
import os
from datetime import datetime

db_path = "/tmp/test.db"
output_dir = "/home/meet/code/blog/author/content/links"

os.makedirs(output_dir, exist_ok=True)

conn = sqlite3.connect(db_path)
conn.row_factory = sqlite3.Row
cursor = conn.cursor()

cursor.execute("""
    SELECT p.title, p.slug, p.content, p.created_at
    FROM posts p
    JOIN post_types pt ON p.type_id = pt.id
    WHERE pt.name = 'Link'
    ORDER BY p.created_at DESC
""")

count = 0
for row in cursor.fetchall():
    title = row['title']
    slug = row['slug']
    content = row['content']
    created_at = row['created_at']
    
    # Parse date
    try:
        date_obj = datetime.fromisoformat(created_at.replace('Z', '+00:00'))
        date_str = date_obj.strftime('%Y-%m-%d')
    except:
        date_str = "2024-01-01"
    
    # Create markdown frontmatter
    md_content = f"""---
title: "{title.replace('"', '\\"')}"
date: {date_str}
draft: false
---

{content}
"""
    
    # Write file
    filename = f"{output_dir}/{slug}.md"
    with open(filename, 'w') as f:
        f.write(md_content)
    
    count += 1
    if count % 100 == 0:
        print(f"Exported {count} links...")

print(f"Total links exported: {count}")
conn.close()
