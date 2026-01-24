---
title: "SQLite: How it works: Richard Hipp"
date: 2026-01-24
draft: false
---

# SQLite: How it works: Richard Hipp

**Link:** https://youtu.be/ZSKLA81tBis

## Context

What a banger of a presentation and talk. Explained so much, in depth, in such a short time. It helped me understand what SQLite actually is, it’s a parser + virtual machine to run the core part and basically the fopen function in C to actually perform the operation.
    - One unique insight here is
          - Reading 10 files content from disk is slower than reading those file contents from SQLite
          - Why? Because the database file is opened once and the reading happens in that instance only, data is stored in pages (fragments of memory), so it’s just a matter of reading bytes at a specific order.
          - But reading 10 different files on disk will make you use fopen 10 times, and that is slow!
          - 200 IQ move from SQLite team, have never seen such a beautiful solution to almost all the problems in the data world.

Double click to interact with video

**Source:** techstructive-weekly-53
