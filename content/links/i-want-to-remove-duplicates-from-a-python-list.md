---
title: "I want to remove duplicates from a python list"
date: 2026-01-24
draft: false
---

# I want to remove duplicates from a python list

**Link:** https://www.thepythoncodingstack.com/p/remove-duplicates-from-python-list

## Context

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
  Reference: Python Coding Stack: [I want to remove duplicates from a python list](https://www.thepythoncodingstack.com/p/remove-duplicates-from-python-list)

[![](https://substackcdn.com/image/fetch/$s_!Dn3k!,w_56,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fab4a59e8-e362-456b-8427-934e87c31a0d_600x600.png)The Python Coding StackI Want to Remove Duplicates from a Python List • How Do I Do It?Another short article today to figure out ways to remove duplicate values from a list. The ideal solution depends on what you really need…Read more16 days ago · 18 likes · 6 comments · Stephen Gruppetta](https://www.thepythoncodingstack.com/p/remove-duplicates-from-python-list?utm_source=substack&utm_campaign=post_embed&utm_medium=web)

## Tech News

**Source:** techstructive-weekly-48
