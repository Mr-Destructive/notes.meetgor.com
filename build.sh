#!/bin/bash
set -e

# Ensure index.html exists
if [ ! -f "public/index.html" ] || [ ! -s "public/index.html" ]; then
  cat > public/index.html << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Notes - Technical Posts & Links</title>
    <style>
        body { font-family: system-ui, sans-serif; padding: 3rem 2rem; max-width: 800px; margin: 0 auto; line-height: 1.6; }
        h1 { margin-bottom: 0.5rem; color: #333; }
        .subtitle { color: #666; margin-bottom: 2rem; font-size: 1.1rem; }
        .nav { display: flex; gap: 2rem; flex-wrap: wrap; margin: 2rem 0; }
        .nav-item { flex: 1; min-width: 150px; padding: 1.5rem; border: 1px solid #ddd; border-radius: 8px; text-align: center; transition: all 0.3s; text-decoration: none; color: inherit; }
        .nav-item:hover { box-shadow: 0 4px 12px rgba(0,0,0,0.1); transform: translateY(-2px); border-color: #0066cc; }
        .nav-item strong { display: block; font-size: 1.2rem; margin-bottom: 0.5rem; }
        .nav-item .desc { color: #888; font-size: 0.9rem; }
    </style>
</head>
<body>
    <h1>Notes & Technical Posts</h1>
    <p class="subtitle">A collection of articles, curated links, and technical newsletters</p>
    
    <div class="nav">
        <a href="/type/post/" class="nav-item">
            <strong>📝 Posts</strong>
            <span class="desc">Technical articles</span>
        </a>
        <a href="/type/link/" class="nav-item">
            <strong>🔗 Links</strong>
            <span class="desc">Curated resources</span>
        </a>
        <a href="/type/newsletter/" class="nav-item">
            <strong>📧 Newsletter</strong>
            <span class="desc">Weekly digests</span>
        </a>
    </div>
</body>
</html>
EOF
fi

# Copy public to output
cp -r public/* .
