---
title: "The fastest way to detect vowel in a string (Python)"
date: 2026-01-24
draft: false
---

# The fastest way to detect vowel in a string (Python)

**Link:** https://austinhenley.com/blog/vowels.html

## Context

1. Wow, this dude just found 11 legit (almost 13) ways to detect vowels in a string in python.Such a great depth, the benchmarks feels so intuitive as why each way performs the way it does.
    2. Here are all the ways it did it
          1. For loop: Simple, readable. Fastest for small strings
          2. C-Styled for loop: Uses or comparisons, but surprisingly much slower
          3. Nested for loop: Totally exhaustive, but slow
          4. Set intersection: Clever and clean. Great when strings are long or vowels are sparse
          5. Generator expression: Pythonic one-liner. Reasonably fast, readable
          6. Recursion: Functional but inefficient. Crashes on long strings
          7. Regex search: Shockingly fast. Calls C-level code internally
          8. Regex replace: Works but inefficient. Doesn’t short-circuit
          9. Filter: Readable but wasteful because it processes the whole string
          10. Map: Similar to filter but slightly better
          11. Prime Numbers: Extremely creative. Maps characters to primes, uses GCD. Way too slow to be practical
    3. Would like to do something in Golang, it sounds so fun that I can’t stop thinking about so many ways to do so trivial things.

**Source:** techstructive-weekly-54
