---
title: "SQLite’s WAL Mode is faster than DELETE Mode"
date: 2026-01-24
draft: false
---

# SQLite’s WAL Mode is faster than DELETE Mode

**Link:** https://youtu.be/qf0GqRz-c74?si=HZ_1yav_DFOzyiOn

## Context

This is so well explained, first showed everything what each one is and then the benchmark just makes everything clear.
    - The WAL mode basically writes the changes in a separate file and merges to the original db file whenever required, hence there is no overhead when reading or writing multiple writers or readers.
    - The delete mode is like a backup, a journal, it keeps pages of the data that are to be changed and after it is committed it deletes the file, that clearly looks slow.

Double click to interact with video

**Source:** techstructive-weekly-52
