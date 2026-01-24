---
title: "Anatomy of a Request: A deep dive of a http request processing from the  Backend side"
date: 2026-01-24
draft: false
---

# Anatomy of a Request: A deep dive of a http request processing from the  Backend side

**Link:** https://youtu.be/s0r3Aky9I5g

## Context

Woah! That is a ton of computation.
    - On Client: Creating the payload, encryption (write copy), loading in kernel space, sending the data
    - On backend: Received the data, reading to the user space, decryption, decoding (serialization) of the body.
    - So many steps are there, the speaker rightly said, its a fascinating field, the more you go deeper, the more stuff is there to explore and learn.

Double click to interact with video

**Source:** techstructive-weekly-53
