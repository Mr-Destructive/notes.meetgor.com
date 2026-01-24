---
title: "Using set.update to append a list of keys of a dictionary"
date: 2026-01-24
draft: false
---

# Using set.update to append a list of keys of a dictionary

**Link:** https://docs.python.org/3/library/stdtypes.html#frozenset.update

## Context

```go
    tables = [
        {"id": 1, "name": "Alice", "age": 25, "email": "alice@example.com"},
        {"id": 2, "name": "Bob", "city": "New York", "phone": "555-1234"},
        {"id": 3, "country": "USA", "zip": "12345", "email": "charlie@example.com"},
        {"id": 4, "name": "Dana", "state": "California", "city": "San Francisco"}
    ]
    unique_keys = set()
    for table in tables:
    for key in table:
    unique_keys.add(key)
    print(unique_keys)
    {'id', 'name', 'age', 'email', 'city', 'phone', 'country', 'zip', 'state'}
    ```
    
    A more cleaner way would be this:
    
    ```go
    for table in tables:
    unique_keys.update(table.keys())
    ```
    
* OpenAI Function Call is not good compared to non-functional prompts. By functional prompt I mean the structure of the response is provided as an object and the LLM has to respond adhering to that structure, this looks good, but not sure why it goofs up the actual text provided to it. Whereas with the normal(non-functional) prompt the response is much better, as we have more control over the things that can be added, validations, and specific structure. I am surprised that normal prompts can also give responses in a structured way that is too consistent.

**Source:** techstructive-weekly-6
