---
title: "Memoization"
date: 2026-01-24
draft: false
---

# Memoization

**Link:** https://dev.to/mwong068/what-is-memoization-4359

## Context

[**Memoization**](https://dev.to/mwong068/what-is-memoization-4359) **for caching problems where things are too repetitive** While solving the advent of code problem day 11, this concept could be handy as it helps to store the already computed values in a map or a cache or data store, whatever you wanna call it at your scale. However, the essence is to store the computed value for the later stage, where there are clear-cut rules or defined behavior of the computation required to do certain things. In the context of this problem, we have to check how many stones will be there at the end of 25 blinks, basically at each blink, there are certain rules like if there is a stone numbered 0, it changes to 1, there is one more rule and you can check the whole question [here](https://adventofcode.com/2024/day/11). It can be solved using a map where we know a stone number can be split into n number of stones, so we use that when we encounter that number in the future and don’t re-compute it again.

**Source:** techstructive-weekly-20
