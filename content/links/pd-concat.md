---
title: "pd.concat"
date: 2026-01-24
draft: false
---

# pd.concat

**Link:** https://pandas.pydata.org/docs/reference/api/pandas.concat.html

## Context

* Combining Pandas Dataframes
    

Let’s say I have a list of dataframes, i want to combine them with certain number but not the contents, just append the next df to the end of the current df. I used the [pd.concat](https://pandas.pydata.org/docs/reference/api/pandas.concat.html) function to combine the slices of pandas dataframe in a single list.

```go
PAGE_LIMIT = 2
PAGES = 5
for i in range(0, PAGES, PAGE_LIMIT):
if i + PAGE_LIMIT <= num_of_pages:
df_batch = pd.concat(df_list[i:i + PAGE_LIMIT], axis=0)
else:
df_batch = pd.concat(df_list[i:], axis=0)
```

**Source:** techstructive-weekly-15
