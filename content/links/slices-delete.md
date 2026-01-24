---
title: "Slices.Delete"
date: 2026-01-24
draft: false
---

# Slices.Delete

**Link:** https://pkg.go.dev/slices#Delete

## Context

* [Slices.Delete](https://pkg.go.dev/slices#Delete) in Golang is Wired
    
    I mean just look at this
    

```go
s := []int{1,2,3}
slices.Delete(s, 1, 2)
fmt.Println(s)
// [1, 3, 3]
s = []int{1,2,3}
s = slices.Delete(s, 1, 2)
fmt.Println(s)
// [1,3]
```

I mean, If the function is returning the modified slice, why are we mutating the original one? Do one or the other, not both.

The safer route will be like this then:

```go
s = []int{1,2,3}
newS := slices.Delete(slices.Clone(s), 1, 2)
fmt.Println(s)
// [1, 2, 3]
fmt.Println(newS)
// [1,3]
```

I need to dive deep into why this is the way it is. Looks pretty confusing to mean.  
AFTER A WHILE:

[This article](https://medium.com/google-cloud/go-slices-deleting-items-and-memory-usage-81419317db3d) explains a bit clearly why the original slice is useless after the operation and we need to pass a copy of that in order to avoid it’s mutation.

## Tech News

**Source:** techstructive-weekly-19
