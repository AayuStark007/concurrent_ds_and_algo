package main

type Queue interface {
	Enqueue(value int)
	Dequeue() (int, bool)
}
