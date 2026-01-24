---
title: "GGUF, the file format for storing LLM model weights"
date: 2026-01-24
draft: false
---

# GGUF, the file format for storing LLM model weights

**Link:** https://vickiboykis.com/2024/02/28/gguf-the-long-way-around/

## Context

I wanted to evaluate an idea for a project. Running models from a file, and this file format is what I needed. Using this format and a binding with llama.cpp and other libraries, this can be used for inference later to actually run the model
        
* Using llama.cpp python bindings to run a model with a gguf file
    
    * We can use the llama.cpp or other library binding to load the file in memory, and the binding library will use the inference to get the tokens out from the given prompt
        
* Creating Python Lambda functions in Vercel
    
    * The snippet is what you need to get up and running with Python serverless functions in Vercel
        
        ```bash
        import json

**Source:** techstructive-weekly-40
