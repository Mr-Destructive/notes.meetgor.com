---
title: "JSON PATCH"
date: 2026-01-24
draft: false
---

# JSON PATCH

**Link:** https://datatracker.ietf.org/doc/html/rfc6902

## Context

- I learned a lot of things about the JSON PATCH method, I will be writing a detailed post about using JSON Patch in the article **100 days of Golang: HTTP PATCH Method, **which should be live after this newsletter.
    - The HTTP PATCH method is like a PUT request but for updating only specific fields and not the entire resource in that sense, you only send the fields to be updated compared to the PUT request where you have to send the entire resource (including the ones that you don’t want to update).
    - [JSON PATCH](https://datatracker.ietf.org/doc/html/rfc6902) is a specific type of PATCH method, where the payload is a JSON patch document, which includes
          - The operation to be performed
          - The field (path) to be updated
          - The value of the field/path to update to
    - [JSON Merge PATCH](https://datatracker.ietf.org/doc/html/rfc7386) is another type of PATCH method where the payload is a JSON body with the fields and values you want to update (like put but only the fields to update)

**Source:** techstructive-weekly-24
