For our queue implementation we don't have the concept of allocating from the free list as we can utilize the GC to do the recycling.
However, for cases were higher perf is needed, we can use something like `sync.Pool` to alleviate the GC pressure. 