---
title: "How HTTP/2 Works and How to Enable It in Go"
date: 2026-01-24
draft: false
---

# How HTTP/2 Works and How to Enable It in Go

**Link:** https://victoriametrics.com/blog/go-http2/?ref=dailydev

## Context

[**How HTTP/2 Works and How to Enable It in Go**](https://victoriametrics.com/blog/go-http2/?ref=dailydev) This is also another post that I took the time to read and was worth it. I honestly don’t know how HTTP 2 works. To some extent, I know how HTTP 1 works, but if someone went a bit deeper, I would start breaking sweat. I really need to implement HTTP from scratch to understand the network stack—one day or day one. - OH, the article, yes it talked in detail about what is the problem with HTTP 1 and how HTTP somewhat solves it. - It is about breaking down the data into frames and makes sure the client has received the frame even if the previous frame is delayed. The fastest frame is served so, it doesn’t block the latest requests. ## Wrote

**Source:** techstructive-weekly-25
