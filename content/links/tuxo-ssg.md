---
title: "Tuxo SSG"
date: 2026-01-24
draft: false
---

# Tuxo SSG

**Link:** https://github.com/Mr-Destructive/tuxo/blob/main/.github/workflows/cronjob.yml

## Context

* **Setting a Github Pages on GitHub Actions Workflow**:  
    I used this snippet in my [Tuxo SSG](https://github.com/Mr-Destructive/tuxo/blob/main/.github/workflows/cronjob.yml), which would serve as the output directory of the generated SSG files to GitHub Pages.
    
    ```go
        - name: GitHub Pages
          uses: crazy-max/ghaction-github-pages@v3
          with:
            target_branch: output-branch
            build_dir: my_app/
            jekyll: false
          env:
            GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    ```
    

## Tech News

**Source:** techstructive-weekly-1
