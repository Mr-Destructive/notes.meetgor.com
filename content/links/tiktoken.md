---
title: "TikToken"
date: 2026-01-24
draft: false
---

# TikToken

**Link:** https://github.com/openai/tiktoken

## Context

Using the `next` method in python to get the first non-empty or truth value from a dict or any expression.
    
    ```go
    row = {
        "name": {
            "value": "123",
            "position": [40, 40],
        },
        "city": {
            "value": "mumbai",
            "position": [50, 30],
        }
    }
    row_position = next((attrs["position"] for attrs in row.values() if isinstance(attrs, dict) and attrs.get("value")), None)
    ```
    
* Combining Pandas Dataframes

**Source:** techstructive-weekly-15
