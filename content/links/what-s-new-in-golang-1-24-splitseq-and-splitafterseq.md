---
title: "What’s new in Golang 1.24, SplitSeq and SplitAfterSeq"
date: 2026-01-24
draft: false
---

# What’s new in Golang 1.24, SplitSeq and SplitAfterSeq

**Link:** https://youtu.be/hee9KUhvQsY?si=gwAjbrxtEQfL1NHC

## Context

Last weekend I dived into the golang 1.24, and set up the system to install the latest version of Go that later became a TIL by experimenting and exploring the new change, the video about Split and the new SplitSeq was created.
    
    The basic difference is about how we store and iterate the splitter strings, the prior stores it as a slice, the other creates a specific type as an iterator to iterate it on the go.
    

* Wrote: TIL: Build Golang from Source for v1.23+  
    My first post on the substack is an awkward TIL, I don’t care if it provides value or not, it is just a thing I learned myself while doing something. So, thought of sharing it here.  
    This was learned while creating the above-mentioned video, I wanted to explore the features coming to Go in 1.24 which is supposed to be released in February 2025. So, while installing I found a couple of ways to build golang from the source, with a bit of trial and error was able to get the simplest way I can think of here.

**Source:** techstructive-weekly-17
