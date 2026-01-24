---
title: "Using set.update to append a list of keys of a dictionary"
date: 2026-01-24
draft: false
---

# Using set.update to append a list of keys of a dictionary

**Link:** https://docs.python.org/3/library/stdtypes.html#frozenset.update

## Context

* [Using set.update to append a list of keys of a dictionary](https://docs.python.org/3/library/stdtypes.html#frozenset.update): Let’s say I have a list of dictionaries of some sort, I want to keep track of all the unique keys in those dictionaries, the dirty and the long way would be this:
    
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

**Source:** techstructive-weekly-6
