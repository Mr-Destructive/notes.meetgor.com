---
title: "Go escape analysis and why my function return worked.<!-- --> - Bonnie Simon"
date: 2025-11-25
slug: go-escape-analysis-and-why-my-function-return-worked-bonnie-simon
draft: false
type: link
description: ""
tags: ["link"]
---

Nice, actually it remembers the pointer if it returns on the function, it allocates on the heap and not the stack so it can stay till no other references hold in the program