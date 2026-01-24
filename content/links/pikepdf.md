---
title: "pikepdf"
date: 2026-01-24
draft: false
---

# pikepdf

**Link:** https://pikepdf.readthedocs.io/en/latest/

## Context

* Reading and Writing a PDF file with [pikepdf](https://pikepdf.readthedocs.io/en/latest/) (useful for splitting and manipulating files)
    
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

**Source:** techstructive-weekly-2
