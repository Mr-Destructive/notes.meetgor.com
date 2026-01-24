---
date: 2025-10-17
draft: false
link: https://marcobambini.substack.com/p/the-secret-life-of-a-local-first
preview_description: A deep dive into how CRDT-powered local-first apps track, merge,
  and sync every INSERT, UPDATE, and DELETE inside SQLite.
preview_image: https://substackcdn.com/image/fetch/$s_!LXtf!,w_1200,h_675,c_fill,f_jpg,q_auto:good,fl_progressive:steep,g_auto/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F1cdd80d7-058a-4349-97ab-dc24a1555fe1_1920x1080.heic
title: Secret Life of Local First Value
tags:
- newsletter
- sql
---
# Secret Life of Local First Value

**Link:** https://marcobambini.substack.com/p/the-secret-life-of-a-local-first

## Context

This is lovely. So well explained what CRDTs are. It’s like a log of what happened in a row of a table. Like column-level details of updation/insertion and deletion. It makes sense now. The metadata table is the crux of this structure. What would happen if the database itself crashes? That is unlikely, I think. SQLite cannot crash at least locally. Nice thinking here.

**Source:** techstructive-weekly-64
