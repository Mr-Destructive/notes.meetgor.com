---
title: "Meta AI API wrapper"
date: 2026-01-24
draft: false
---

# Meta AI API wrapper

**Link:** https://github.com/Mr-Destructive/meta-ai-golang

## Context

* Difference between cookies Add and Set in URL Values:
    
    [Add](https://pkg.go.dev/net/url#Values.Add): Appends the value to the key without replacing existing values (useful for handling multiple values for a single key).
    
    [Set](https://pkg.go.dev/net/url#Values.Set): Replaces the existing value for the key (ensures that only one value is associated with the key).
    
    I learnt this will working with the Meta AI API wrapper in Golang. The API uses payload as a URL encoded body and will append key-value pairs to the request body, the subtle difference can cause nil pointer access if not initialized and used the appropriate method correctly. I think I will write a blog post on this.

**Source:** techstructive-weekly-5
