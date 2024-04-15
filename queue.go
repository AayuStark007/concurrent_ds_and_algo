package main

type Queue[T any] interface {
	Enqueue(value T)
	Dequeue() (T, bool)
}
