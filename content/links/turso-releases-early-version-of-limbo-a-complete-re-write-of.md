---
title: "Turso releases early version of Limbo: A complete re-write of SQLite in Rust"
date: 2026-01-24
draft: false
---

# Turso releases early version of Limbo: A complete re-write of SQLite in Rust

**Link:** https://turso.tech/blog/introducing-limbo-a-complete-rewrite-of-sqlite-in-rust

## Context

[Introducing Limbo: A complete rewrite of SQLite in Rust](https://turso.tech/blog/introducing-limbo-a-complete-rewrite-of-sqlite-in-rust) This is pretty cool, I was thinking it might be already [there](https://github.com/rusqlite/rusqlite), yes definitely, but they are just not re-writing it, they are forking and adding features on top of it which is absolutely wild. Surely there won’t be time gains, but they now have a lot more control over what needs to be changed and included while adding more or even upstreaming from SQLite-core. People might call re-writing a waste of time (especially for such a well-developed and stable tool), but people forget they are making something from scratch gives you a whole different depth of understanding than just forking it. It will pay dividends slowly in the long run, pulling from upstream might be challenging and tedious but since SQLite is rock-solid, there won’t be any breaking changes that might get added to it, so a win-win for Turso in my opinion.

**Source:** techstructive-weekly-20
