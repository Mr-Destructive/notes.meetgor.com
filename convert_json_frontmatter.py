#!/usr/bin/env python3
import os
import json
import glob
from datetime import datetime

posts_dir = "/home/meet/code/blog/author/exports/content/posts"

for filepath in glob.glob(os.path.join(posts_dir, "*.md")):
    with open(filepath, 'r') as f:
        content = f.read()
    
    # Check if it starts with JSON
    if not content.startswith('{'):
        continue
    
    # Find the end of JSON object
    brace_count = 0
    json_end = 0
    for i, char in enumerate(content):
        if char == '{':
            brace_count += 1
        elif char == '}':
            brace_count -= 1
            if brace_count == 0:
                json_end = i + 1
                break
    
    if json_end == 0:
        print(f"Could not find JSON end in {filepath}")
        continue
    
    try:
        json_str = content[:json_end]
        rest_content = content[json_end:].strip()
        data = json.loads(json_str)
        
        # Convert to YAML frontmatter
        yaml_lines = ["---"]
        
        # Add type
        post_type = data.get('type', 'post')
        yaml_lines.append(f"type: {post_type}")
        
        # Add title
        if 'title' in data:
            yaml_lines.append(f'title: "{data["title"]}"')
        
        # Add subtitle if exists
        if 'subtitle' in data:
            yaml_lines.append(f'subtitle: "{data["subtitle"]}"')
        
        # Add date
        if 'date' in data:
            date_str = data['date']
            try:
                # Try to parse and reformat
                dt = datetime.fromisoformat(date_str.replace('+0530', '+05:30').replace(' +0530', '+05:30'))
                yaml_lines.append(f"date: {dt.isoformat()}")
            except:
                yaml_lines.append(f'date: "{date_str}"')
        
        # Add slug if exists
        if 'slug' in data:
            yaml_lines.append(f'slug: "{data["slug"]}"')
        
        # Add series if exists
        if 'series' in data:
            series = data['series']
            if isinstance(series, list) and len(series) > 0:
                yaml_lines.append(f'series: ["{series[0]}"]')
        
        # Add tags if exists
        if 'tags' in data:
            tags = data['tags']
            if isinstance(tags, list) and len(tags) > 0:
                yaml_lines.append("tags:")
                for tag in tags:
                    yaml_lines.append(f"  - {tag}")
        
        # Add image_url if exists
        if 'image_url' in data:
            yaml_lines.append(f'image_url: {data["image_url"]}')
        
        # Add cover if exists
        if 'cover' in data:
            yaml_lines.append(f'cover: {data["cover"]}')
        
        yaml_lines.append("---")
        yaml_frontmatter = "\n".join(yaml_lines)
        new_content = yaml_frontmatter + "\n\n" + rest_content
        
        with open(filepath, 'w') as f:
            f.write(new_content)
        
        print(f"Converted: {os.path.basename(filepath)}")
    
    except Exception as e:
        print(f"Error processing {filepath}: {e}")
