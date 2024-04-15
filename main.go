package main

import (
	"fmt"
	"sync"
)

func main() {
	Q := NewQueue()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		Q.Enqueue(1)
		Q.Enqueue(2)
		Q.Enqueue(3)
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
