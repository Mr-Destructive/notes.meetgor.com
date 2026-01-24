---
type: sqlog
title: 'SQLite dot commands: header command'
slug: sqlite-dot-command-header
date: 2025-09-03
tags:
- sql
---
The `.headers` command is crucial as might effect the way the result set is displayed. The first row would be set to the name of the columns in the output of the relevant commands like `.output`, `.once`, or even your simple standard out queries if it is set on.

```
.headers on

OR

.header on
```

Some how either of them work. You need to set it to `on` to enable the headers in the output of the result set. And you can turn off with `.headers off`.

Which modes are effected with this command if set on or off?

Well we need to think about how effected means
1. Only added in the first row
2. Added in each row
3. No effect

