## Concurrent Data Structures and Algorithms

Repository hosting sample code and examples from implementations of various blocking and/or non-blocking concurrent data-structures and the algorithms powering them.

### Concurrent Queue Algorithms

Implements blocking and non-blocking concurrent queues as described in the paper by **Michael & Scott**: https://www.cs.rochester.edu/~scott/papers/1996_PODC_queues.pdf

The paper describes the following two Concurrent Queue Algorithms:
- **Non-Blocking Concurrent Queue**: A simple, non-blocking, practical and fast concurrent queue algorithm which utilizes the universal atomic primitives like *compare_and_swap* (*load_linked/store_conditional* in ARM like processors) which is commonly provided by all hardware today.

- **Blocking Queue**: A queue with separate head and tail pointer locks which allows only one enqueue and one dequeue to proceed at a time. It is recommended to use with hardware which only provides simple atomic primitives like *test_and_set*. However, the paper recommends using a single-lock version in case there are only 1-2 contending processes.

#### Notes on Implementation
- **ABA Problem**: In the paper, the authors have described the [ABA problem](https://www.baeldung.com/cs/aba-concurrency) which is prevalent in most non-blocking algorithms. There are some solutions mentioned in the paper like using modification counter along with pointers, but those require the use of double-word atomic compare_and_swap primitives which are not commonly found. 
    - However, in our case, the language of choice, Go, is a garbage collected language. The benefit of that is, we do not need to explicitly recycle dequeued nodes and maintain a free list. Thus, the ABA problem cannot occcur as each newly created node will have an unique reference (since the old node is still alive and its reference is valid). When same node is enqueued again, it will have a new reference and thus fail the CAS check. Atleast, this is how it is justified in the [ConcurrentLinkedQueue](https://github.com/openjdk-mirror/jdk7u-jdk/blob/f4d80957e89a19a29bb9f9807d2a28351ed7f7df/src/share/classes/java/util/concurrent/ConcurrentLinkedQueue.java#L113) implementation of the same algorithm in Java.

    - In non-GC languages, like C++ we have the concept of [Hazard Pointers](https://wiki.sei.cmu.edu/confluence/display/c/CON09-C.+Avoid+the+ABA+problem+when+using+lock-free+algorithms) which allows the thread to keep track of the pointers it is using thus keeping track when it is being modified or used by the other thread.

- For our Queue implementation we don't have the concept of allocating from the free list as we can utilize the GC to collect and reclaim the node memory. However, for cases were higher perf is needed, we can use something like `sync.Pool` to alleviate the GC pressure.

- **Go Generics**: I have tried to define the contracts of this algorithm which allow it to leverage the Golang generics, thus the implementations can work with any data type in a type-safe manner.

#### Benchmarks
 *coming soon*