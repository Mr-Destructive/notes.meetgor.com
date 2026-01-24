---
title: "sum of squares and squares of sum problem in Exercism SQLite track"
date: 2026-01-24
draft: false
---

# sum of squares and squares of sum problem in Exercism SQLite track

**Link:** https://exercism.org/tracks/sqlite/exercises/difference-of-squares

## Context

- SQL Week #3
    - String concatenation in SQLite
          - printf function that exists in SQLite. That is such a great thing to have, C like function, just drop the placeholder for the appropriate type and it will be a formatted string, neat and tidy.
          - || operator for simple concatenation. You can just use “Hello” || “ World!” to get a string “Hello World!”. It’s compact but for large number of strings, or readability, it gets a little tricky.
    - I finally wrapped my head around autoincrement and sqlite_sequence table.
          - Autoincrement is a constraint that will force the newly inserted rows to use primary key id values greater than any existing or deleted primary key id values.
          - If the highest row created is deleted before inserting a new record, the sqlite_sequence comes in handy to fetch the max value of the primary key id across the table, since the value is stored in a separate table called the sqlite_sequence with table_name and the seq columns.
          - The sqlite engine decides to get the max of the seq value and the current max row id (it can effectively get it using B+ trees, as it will be the leftmost node or right-most node, however its stored)
          - I will be creating a separate blog on this, a deep dive on the various cases we can run into if the sqlite_sequence table gets altered.
    - USE common table expressions in sqlite
          - This is a way to create a temporary table (like only valid till the query completes running) and use it in the main query(can have nested queries too)
          - The syntax looks like this :
            ```
            WITH <temp-table-name> AS (SELECT something from somewhere)
            SELECT something, <temp-table-name>.something from elsewhere 
            ```
          - This is something I studied while solving the difference of [sum of squares and squares of sum problem in Exercism SQLite track](https://exercism.org/tracks/sqlite/exercises/difference-of-squares).
    - SQL’s IIF is equivalent for CASE WHEN THEN
          - The IIF is like an handy if else block to use when having nested conditions.
          - This works and looks neat for smaller expressions, like one or two condition max, after that its better to use CASE
          - The syntax looks something like:
            ```
            SELECT IIF(something > 10, "YES", "NO") as answer FROM somewhere;
            ```

### Interesting Links

**Source:** techstructive-weekly-54
