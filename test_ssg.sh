#!/bin/bash

# Test script for Static Site Generation
# Verifies all features and links work correctly

echo "=========================================="
echo "SSG Build and Link Test"
echo "=========================================="
echo ""

cd exports || exit 1

# Clean previous build
rm -rf public

# Build
echo "1. Building Hugo site..."
hugo --minify 2>&1 | grep -E "(Total|Pages|Error)"
if [ $? -ne 0 ]; then
    echo "  ✗ Build failed"
    exit 1
fi
echo "  ✓ Build successful"
echo ""

# Check structure
echo "2. Checking generated files..."
files_ok=true

check_file() {
    if [ -f "$1" ]; then
        echo "  ✓ $1"
    else
        echo "  ✗ Missing: $1"
        files_ok=false
    fi
}

check_file "public/index.html"
check_file "public/posts/index.html"
check_file "public/posts/getting-started-go/index.html"
check_file "public/tags/index.html"
check_file "public/tags/go/index.html"
check_file "public/sitemap.xml"

if [ "$files_ok" = false ]; then
    echo ""
    echo "  ✗ Some files missing"
    exit 1
fi
echo ""

# Check home page
echo "3. Testing home page..."
home=$(cat public/index.html)

contains() {
    if echo "$home" | grep -q "$1"; then
        echo "  ✓ Contains: $2"
        return 0
    else
        echo "  ✗ Missing: $2"
        return 1
    fi
}

contains '<h1><a href=/>Notes</a></h1>' "Site title"
contains 'A collection of technical notes and posts' "Site description"
contains 'href=/posts/' "Posts link"
contains 'Getting Started with Go' "Post title"
contains 'href=https://notes.meetgor.com/posts/getting-started-go/' "Post link"
contains 'Learn Go basics' "Post description"
echo ""

# Check posts list
echo "4. Testing posts list..."
posts_list=$(cat public/posts/index.html)

if echo "$posts_list" | grep -q "Getting Started with Go"; then
    echo "  ✓ Post appears in list"
else
    echo "  ✗ Post missing from list"
    exit 1
fi

if echo "$posts_list" | grep -q "href=https://notes.meetgor.com/posts/getting-started-go/"; then
    echo "  ✓ Post link is correct"
else
    echo "  ✗ Post link incorrect"
    exit 1
fi
echo ""

# Check individual post
echo "5. Testing individual post page..."
post=$(cat public/posts/getting-started-go/index.html)

if echo "$post" | grep -q "<h1>Getting Started with Go</h1>"; then
    echo "  ✓ Post title rendered"
else
    echo "  ✗ Post title missing"
    exit 1
fi

if echo "$post" | grep -q "Go is a powerful language"; then
    echo "  ✓ Post content rendered"
else
    echo "  ✗ Post content missing"
    exit 1
fi

if echo "$post" | grep -q "<h2.*>Installation</h2>"; then
    echo "  ✓ Markdown headers parsed"
else
    echo "  ✗ Markdown not processed"
    exit 1
fi

if echo "$post" | grep -q "href=/"; then
    echo "  ✓ Home breadcrumb link"
else
    echo "  ✗ Breadcrumb link missing"
    exit 1
fi
echo ""

# Check tags
echo "6. Testing tag pages..."
tag_go=$(cat public/tags/go/index.html)

if echo "$tag_go" | grep -q "Getting Started with Go"; then
    echo "  ✓ Post shows in tag page"
else
    echo "  ✗ Post missing from tag"
    exit 1
fi
echo ""

# Check navigation
echo "7. Testing navigation links..."
if echo "$post" | grep -q '<a href=/posts/>Posts</a>'; then
    echo "  ✓ Posts nav link"
else
    echo "  ✗ Posts nav link missing"
    exit 1
fi

if echo "$post" | grep -q '<a href=/>Home</a>'; then
    echo "  ✓ Home nav link"
else
    echo "  ✗ Home nav link missing"
    exit 1
fi
echo ""

# Summary
echo "=========================================="
echo "✓ All tests passed!"
echo "=========================================="
echo ""
echo "Site structure:"
find public -name "*.html" | sort | sed 's/^/  /'
echo ""
echo "Ready for deployment!"
