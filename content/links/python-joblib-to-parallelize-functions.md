---
title: "Python joblib to parallelize functions"
date: 2026-01-24
draft: false
---

# Python joblib to parallelize functions

**Link:** https://joblib.readthedocs.io/en/stable/parallel.html

## Context

[Python joblib to parallelize functions](https://joblib.readthedocs.io/en/stable/parallel.html): This is one of the libraries and ideas that I learned in my daily work to get things done. I used it basically to call a function with a list of arguments (different) multiple calls at the same time to save time. So, if a single function call takes 10 seconds and I have 3 function calls to make, sequentially it would take 30 seconds but with parallel calls, those could be done with just over 10 seconds.

**Source:** techstructive-weekly-6
