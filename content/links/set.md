---
title: "Set"
date: 2026-01-24
draft: false
---

# Set

**Link:** https://pkg.go.dev/net/url#Values.Set

## Context

I learnt this will working with the Meta AI API wrapper in Golang. The API uses payload as a URL encoded body and will append key-value pairs to the request body, the subtle difference can cause nil pointer access if not initialized and used the appropriate method correctly. I think I will write a blog post on this.
    
* Shuffling Two Lists keeping the order of the corresponding index the same:
    
    What I was doing was testing and evaluating some results on data, and that data was coming from a set of files in a folder, I wanted to randomly shuffle those values. I wanted to track the metrics from the data with the filename, so I created this little function that shuffles two or more lists in a random order and maintains a one-on-one index mapping.
    
    ```go
    import random
    List of file names
    file_names = ["file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt"]
    Corresponding sample data for each file
    sample_data = ["Data from file 1", "Data from file 2", "Data from file 3", "Data from file 4", "Data from file 5"]
    Create a list of indices
    indices = list(range(len(file_names)))
    Shuffle the indices
    random.shuffle(indices)
    Reorder both lists using the shuffled indices
    shuffled_file_names = [file_names[i] for i in indices]
    shuffled_sample_data = [sample_data[i] for i in indices]
    print(shuffled_file_names)
    print(shuffled_sample_data)
    """
    ['file3.txt', 'file5.txt', 'file1.txt', 'file2.txt', 'file4.txt']
    ['Data from file 3', 'Data from file 5', 'Data from file 1', 'Data from file 2', 'Data from file 4']
    """
    ```
    

This is the function, nice simple and modular.

```go
import random
def shuffle_lists(*lists):
"""
Shuffles two or more lists while keeping the order of corresponding elements the same.
Args:
*lists: Two or more lists to be shuffled.
Returns:
A tuple of shuffled lists with the same order of corresponding elements.
"""

**Source:** techstructive-weekly-5
