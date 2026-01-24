---
title: "Python’s pass by value and pass by reference"
date: 2026-01-24
draft: false
---

# Python’s pass by value and pass by reference

**Link:** https://www.thepythoncodingstack.com/p/python-pass-by-value-reference-assignment

## Context

1. This is one hell of a reason, Python gets a little more confusing and less friendly.
    2. TLDR of the post is that if you pass a immutable variable/object to a function call in python, you need to return it back from the function (if the function modifies those immutable objects). Because the object is immutable it won’t get updated inside the function, it will be created a new, so we need to assign it to the modified version when the function returns.
    3. But for mutable objects, the function can modify it and we are passing it to the function, so the object will be updated.

**Source:** techstructive-weekly-54
