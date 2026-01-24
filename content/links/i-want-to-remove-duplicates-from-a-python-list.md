---
title: "I want to remove duplicates from a python list"
date: 2026-01-24
draft: false
---

# I want to remove duplicates from a python list

**Link:** https://www.thepythoncodingstack.com/p/remove-duplicates-from-python-list

## Context

Javascript has quirky behaviour with null, undefined and what equality operator
    - We don’t know if the value is null or undefined, the object might be still undefined but it will be a truthy value
    - The equality operator is very wired, this teaches us that too much flexibility is also bad
          - For instance, the ‘5’ == 5 will be true but ‘5’ === 5 won’t be
          - The triple equal is a type check whereas the double equal is a value check after type casting, which might be a bit unpredictable as the data gets complex.
- Getting the unique elements from a list in Python without changing the order
    - After python 3.7 changes the dictionary, the order of the keys inserted is maintained, so we can use that to create a map of the elements in the list as a key in the dictionary and return the unique elements.
    - Neat little trick, could be well often be used widely in many cases.

  ```
  # this might change the order
  list(set(queue))

  # this will preserve the order
  # works for python > 3.7 
  list(dict.fromkeys(queue))
  ```

**Source:** techstructive-weekly-48
