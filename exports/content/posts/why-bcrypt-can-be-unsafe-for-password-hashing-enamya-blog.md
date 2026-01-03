---
title: "Why bcrypt Can Be Unsafe for Password Hashing ? | enamya | blog"
date: 2025-12-20
slug: why-bcrypt-can-be-unsafe-for-password-hashing-enamya-blog
draft: false
type: link
description: ""
tags: ["link"]
---

Nice to know that bcrypt is not safe for passwords greater than 72 characters, who even would store such a long password.
But that is the thing, subtle decisions, like this is not a password, so we can use bcrypt, and bam you would be wrong

use Argon for this