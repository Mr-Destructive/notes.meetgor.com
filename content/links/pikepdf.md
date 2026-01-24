---
title: "pikepdf"
date: 2026-01-24
draft: false
---

# pikepdf

**Link:** https://pikepdf.readthedocs.io/en/latest/

## Context

```go
    import pikepdf
    from pathlib import Path
    def split_pdf(input_pdf_path, output_dir):
    pdf = pikepdf.Pdf.open(input_pdf_path)
    output_dir = Path(output_dir)
    output_dir.mkdir(exist_ok=True)
    for i, page in enumerate(pdf.pages):
    output_pdf_path = output_dir / f'page_{i+1}.pdf'
    with pikepdf.Pdf.new() as output_pdf:
    output_pdf.pages.append(page)
    output_pdf.save(output_pdf_path)
    ```
    
* LLMs are taking over. I mean not quality-wise, but they seem to be everywhere, almost every company is trying to use AI Agents to make their business look smart (but is not actually). This is a harsh or soft truth that I have to accept and move ahead to leverage in things I would do and build.

**Source:** techstructive-weekly-2
