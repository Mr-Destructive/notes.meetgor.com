---
title: "pd.concat"
date: 2026-01-24
draft: false
---

# pd.concat

**Link:** https://pandas.pydata.org/docs/reference/api/pandas.concat.html

## Context

```go
PAGE_LIMIT = 2
PAGES = 5
for i in range(0, PAGES, PAGE_LIMIT):
if i + PAGE_LIMIT <= num_of_pages:
df_batch = pd.concat(df_list[i:i + PAGE_LIMIT], axis=0)
else:
df_batch = pd.concat(df_list[i:], axis=0)
```

* Deleting elements from the list given the indices
    
    I know this could be easily gotten from GPT but feels good to do it yourself sometimes.
    

```go
def delete_elements_by_indices(lst, indices):
indices_set = set(indices)
return [item for idx, item in enumerate(lst) if idx not in indices_set]
lst = [10, 20, 30, 40, 50]
indices_to_delete = [1, 3]
result = delete_elements_by_indices(lst, indices_to_delete)
print(result)
Output: [10, 30, 50]
```

Fun tidbit,

Chat GPT search is just awesome, it recognises me? really?

**Source:** techstructive-weekly-15
