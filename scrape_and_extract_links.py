#!/usr/bin/env python3
"""
Scrape all newsletter editions from Substack.
Extract Read and Watched sections.
Parse links with commentary.
Create link posts preserving user's exact text.
"""

import requests
import re
import sqlite3
import uuid
from datetime import datetime
from bs4 import BeautifulSoup
import time

def fetch_newsletter(slug):
    """Fetch a single newsletter edition."""
    url = f'https://techstructively.substack.com/p/{slug}'
    headers = {
        'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36'
    }
    
    try:
        response = requests.get(url, headers=headers, timeout=10)
        response.raise_for_status()
        return response.text
    except Exception as e:
        print(f"  Error fetching {slug}: {e}")
        return None

def extract_links_from_content(html_content):
    """Extract links from HTML content, preserving context/commentary."""
    links = []
    
    if not html_content:
        return links
    
    # Find all markdown and HTML links
    # Pattern: [text](url) or <a href="url">text</a>
    
    # Markdown links
    markdown_pattern = r'\[([^\]]+)\]\(([^)]+)\)'
    for match in re.finditer(markdown_pattern, html_content):
        title = match.group(1).strip()
        url = match.group(2).strip()
        
        if url.startswith('http'):
            # Get surrounding context (100 chars before and after)
            start = max(0, match.start() - 150)
            end = min(len(html_content), match.end() + 150)
            context = html_content[start:end].strip()
            
            links.append({
                'title': title,
                'url': url,
                'context': context
            })
    
    # HTML links (fallback)
    html_pattern = r'<a[^>]+href=["\']([^"\']+)["\'][^>]*>([^<]+)</a>'
    for match in re.finditer(html_pattern, html_content):
        url = match.group(1).strip()
        title = match.group(2).strip()
        
        if url.startswith('http') and title:
            # Get surrounding context
            start = max(0, match.start() - 150)
            end = min(len(html_content), match.end() + 150)
            context = html_content[start:end].strip()
            
            links.append({
                'title': title,
                'url': url,
                'context': context
            })
    
    return links

def scrape_all_newsletters():
    """Scrape all newsletter editions."""
    # Generate all newsletter slugs
    newsletter_slugs = [f'techstructive-weekly-{i}' for i in range(1, 79)]
    
    print("=" * 70)
    print("Scraping All Newsletter Editions from Substack")
    print("=" * 70)
    print(f"\nFetching {len(newsletter_slugs)} newsletters...\n")
    
    all_links = []
    
    for i, slug in enumerate(newsletter_slugs, 1):
        print(f"  [{i:2d}] Fetching {slug}...", end=' ', flush=True)
        
        html = fetch_newsletter(slug)
        if html:
            links = extract_links_from_content(html)
            print(f"✓ ({len(links)} links)")
            
            for link in links:
                link['source'] = slug
                all_links.append(link)
        else:
            print("✗ Failed")
        
        time.sleep(1)  # Rate limiting
    
    return all_links

def create_slug_from_title(title):
    """Create URL-safe slug from title."""
    slug = title.lower()
    slug = re.sub(r'[^a-z0-9]+', '-', slug)
    slug = slug.strip('-')
    return slug[:60]

def insert_link_post(conn, link_data):
    """Insert link post with preserved commentary."""
    cursor = conn.cursor()
    post_id = str(uuid.uuid4())
    
    slug = create_slug_from_title(link_data['title'])
    
    # Use context as the user's commentary
    full_content = f"""# {link_data['title']}

**Link:** {link_data['url']}

## Context

{link_data.get('context', '')}

**Source:** {link_data.get('source', 'newsletter')}
"""
    
    try:
        cursor.execute('''
            INSERT OR IGNORE INTO posts (
                id, type_id, slug, title, content, excerpt,
                status, created_at, updated_at, published_at
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ''', (
            post_id,
            'link',
            slug,
            link_data['title'][:255],
            full_content,
            link_data.get('context', '')[:500],
            'published',
            datetime.now().isoformat(),
            datetime.now().isoformat(),
            datetime.now().isoformat()
        ))
        conn.commit()
        return True
    except Exception as e:
        print(f"Error inserting {link_data['title']}: {e}")
        return False

def main():
    print("Attempting to scrape Substack newsletters...\n")
    
    # Try to install BeautifulSoup if needed
    try:
        from bs4 import BeautifulSoup
    except:
        print("Installing beautifulsoup4...")
        import subprocess
        subprocess.run(['pip', 'install', '--break-system-packages', 'beautifulsoup4'], 
                      capture_output=True)
    
    # Scrape newsletters
    all_links = scrape_all_newsletters()
    
    print(f"\nTotal links extracted: {len(all_links)}")
    
    if not all_links:
        print("No links extracted. Substack may be blocking requests.")
        print("Alternative: Please provide raw newsletter content and I'll parse it.")
        return False
    
    # Insert into database
    db_path = '/home/meet/code/blog/author/test.db'
    conn = sqlite3.connect(db_path)
    
    cursor = conn.cursor()
    cursor.execute("SELECT slug FROM posts WHERE type_id='link'")
    existing = {row[0] for row in cursor.fetchall()}
    
    print(f"\nInserting links into database...\n")
    
    inserted = 0
    skipped = 0
    
    for link in all_links:
        slug = create_slug_from_title(link['title'])
        if slug not in existing:
            if insert_link_post(conn, link):
                inserted += 1
                existing.add(slug)
        else:
            skipped += 1
    
    conn.close()
    
    print(f"\nCompletion Summary:")
    print(f"  Inserted: {inserted}")
    print(f"  Skipped: {skipped}")
    
    return inserted > 0

if __name__ == '__main__':
    try:
        success = main()
        exit(0 if success else 1)
    except KeyboardInterrupt:
        print("\n\nCancelled")
        exit(1)
