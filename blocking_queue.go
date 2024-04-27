package main

import "sync"

/*
* BlockingQueue is a Concurrent Blocking Queue with 2-lock based approach.
 */
type BlockingQueue[E any] struct {
	Head  *Node[E]
	Tail  *Node[E]
	hlock sync.Mutex
	tlock sync.Mutex
}

func NewBlockingQueue[E any]() Queue[E] {
	node := NewNode[E]()

	return &BlockingQueue[E]{
		Head: node,
		Tail: node,
	}
}

func (Q *BlockingQueue[E]) Enqueue(item E) {
	// create new node wit value to enqueue
	newNode := NewNodeWithItem[E](item)

	// acquire lock on tail before accessing it
	Q.tlock.Lock()
	defer Q.tlock.Unlock() // lock is released automatically upon exit

	// link the new node at the end of the list
	Q.Tail.next = newNode
	// point tail to the new node
	// note that we need not check for concurrent modification since we are modifying the tail in a lock,
	// so we are guaranteed that no modifications have taken place during this time.
	Q.Tail = newNode
}

func (Q *BlockingQueue[E]) Dequeue() (item E, ok bool) {
	// acquire lock on head before accessing it
	Q.hlock.Lock()
	defer Q.hlock.Unlock() // lock released upon exit from method

	// read head
	head := Q.Head
	// read next ptr to head (will become the new head later)
	next := head.next

	// this means that the queue is empty, so nothing to dequeue
	if next == nil {
		return item, false
	}

	// read the value
	item = next.item
	// point head to the next node
	Q.Head = next

	// make old head as pointing to itself so it is garbage collected
	head.next = head
	return item, true
}
