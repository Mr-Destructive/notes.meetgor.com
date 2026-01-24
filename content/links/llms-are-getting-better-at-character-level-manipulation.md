---
date: 2025-10-17
draft: false
link: https://blog.burkert.me/posts/llm_evolution_character_manipulation/
preview_description: 'Recently, I have been testing how well the newest generations
  of large language models (such as GPT-5 or Claude 4.5) handle natural language,
  specifically counting characters, manipulating characters in a sentences, or solving
  encoding and ciphers. Surprisingly, the newest models were able to solve these kinds
  of tasks, unlike previous generations of LLMs.

  Character manipulation LLMs handle individual characters poorly. This is due to
  all text being encoded as tokens via the LLM tokenizer and its vocabulary. Individual
  tokens typically represent clusters of characters, sometimes even full words (especially
  in English and other common languages in the training dataset). This makes any considerations
  on a more granular level than tokens fairly difficult, although LLMs have been capable
  of certain simple tasks (such as spelling out individual characters in a word) for
  a while.'
preview_image: https://blog.burkert.me/llm-character-text-manipulation.jpg
tags:
- testing
title: LLMs are getting better at character-level manipulation
---

# LLMs are getting better at character-level manipulation

**Link:** https://blog.burkert.me/posts/llm_evolution_character_manipulation/

## Context

Its evident from the test that newer and larger models are better at generalizing Base64 encoding and decoding. So that implies they will get better at character-level manipulation and analysis. Sadly the how many r’s in strawberry problem will be solvable by LLMs Thinking is out of the equation, the crux here is the tokenisation, the better sense of the word you have, the better it understands, but the fine balance between less and more context is critical, and I think it is still being fine tuned to get a sweet spot.

**Source:** techstructive-weekly-64
