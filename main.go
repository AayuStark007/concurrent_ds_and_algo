package main

import (
	"fmt"
	"sync"

	"github.com/aayustark007/concurrent_ds_and_algo/queue"
)

func main() {
	// Q := NewNonBlockingQueue[interface{}]()
	Q := queue.NewBlockingQueue[interface{}]()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		Q.Enqueue(queue.NewNode[int]())
		Q.Enqueue(2)
		Q.Enqueue("ef")
		Q.Enqueue(4)
		Q.Enqueue(5)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		Q.Enqueue(6)
		Q.Enqueue(7)
		Q.Enqueue(8)
		Q.Enqueue(9)

		if val, ok := Q.Dequeue(); ok {
			fmt.Printf("goro1: %d\n", val)
		}

		if val, ok := Q.Dequeue(); ok {
			fmt.Printf("goro2: %d\n", val)
		}

		if val, ok := Q.Dequeue(); ok {
			fmt.Printf("goro3: %d\n", val)
		}
	}()

	wg.Wait()

	for {
		value, res := Q.Dequeue()
		if !res {
			break
		}
		fmt.Println(value)
	}

	fmt.Println("Done!")
}
