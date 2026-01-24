---
type: "sqlog"
title: "SQLite dot commands: Output mode separator command"
slug: sqlite-mode-dot-command-separators
date: 2025-09-02
tags: ["frontend", "sql"]
---

## Using the separator for the ouput

If you wanted to use a specific separator for columns and rows while displaying the result set / table, you can use the `.separator` dot command which can take 2 arguments, first as the separator for the column and the second for the row.

So, if we set use `.separator "|" "---"` then it will split the columns with `|` and each row with `---`. 

```
1|The Hobbit|J.R.R. Tolkien|310|1937-09-21|39.99---2|The Fellowship of the Ring|J.R.R. Tolkien|423|1954-07-29|49.99---3|The Two Towers|J.R.R. Tolkien|352|1954-11-11|49.99---4|The Return of the King|J.R.R. Tolkien|416|1955-10-20|49.99---
```

The output looks wired but I was giving a example.

The row separator is by default a `\n` character or `\r\n` on windows, which is for the list or any other mode. However if you want to add those, you need to specify it in the string like below:

```
.separator "|" "\n---"
```

```
>sqlite>.separator "|" "\n---"
sqlite> select * from books;
1|The Hobbit|J.R.R. Tolkien|310|1937-09-21|39.99
