package main

import (
	"sync/atomic"
	"unsafe"
)

/*
* Node -> good'ol linked list node with data and next ptr
* According to notes from Java's implementation of ConcurrentLinkedQueue:
* > This implementation relies on the fact that in garbage collected systems,
* > there is no possibility of ABA problems due to recycled nodes,
* > so there is no need to use "counted pointers" or related techniques
*
* Need to understand this further but in a nutshell, counters as given in the paper impl are not required here.
 */
type Node[E any] struct {
	item E
	next *Node[E]
}

func NewNode[E any]() *Node[E] {
	return &Node[E]{
		next: nil,
	}
}

func NewNodeWithItem[E any](item E) *Node[E] {
	return &Node[E]{
		item: item,
		next: nil,
	}
}

/*
* NBQueue is a Concurent Non-Blocking Queue based on CAS primitives
 */
type NBQueue[E any] struct {
	Head *Node[E]
	Tail *Node[E]
}

func NewQueue[E any]() Queue[E] {
	node := NewNode[E]()

	return &NBQueue[E]{
		Head: node,
		Tail: node,
	}
}

func (Q *NBQueue[E]) Enqueue(item E) {
	// creating a new node with value to enqueue
	newNode := NewNodeWithItem[E](item)

	var tail *Node[E]
	// since CAS-ing can fail, we keep trying until we succeed to enqueue `node`
	for {
		// save the current tail
		tail = Q.Tail
		// save the next node to tail
		next := tail.next

		// ensure our saved tail is still the Tail
		if tail == Q.Tail {
			// wrt tail -> either we are at the last node or our tail is lagging and we need to advance it further
			if next == nil { // tail is pointing to the last node
				if casNodePtr(&tail.next, nil, newNode) {
					break // successfully enqueued
				}
			} else { // tail not pointing to the last node (some concurrent process enqueued more nodes after we read the tail earlier)
				casNodePtr(&Q.Tail, tail, next) // try to move the tail to newly inserted node
			}
		}
	}
	/*
		We have enqueued the node, we can now try to ensure Q.Tail is pointing to newly inserted node.
		But, its not a necessary thing to do, which allows enquques to be fast as we can allow the Q.Tail to lag.
	*/
	casNodePtr(&Q.Tail, tail, newNode)
}

func (Q *NBQueue[E]) Dequeue() (item E, ok bool) {
	var head *Node[E] // saving current head (which will be dequeued), which allows us to free the node data

	// since CAS-ing can fail (on account that nodes are concurrently dequeued), we must keep trying
	for {
		head = Q.Head
		// read the current tail, required for checking empty Queue case
		tail := Q.Tail
		// next ptr to head (will become the new head later)
		next := head.next

		// ensure that our saved head is still the Head
		if head == Q.Head {
			// either Q is empty, or the tail might be lagging (on account of more items enqueued since we started)
			if head == tail {
				// this means tail is pointing to the last node, so our Q is empty
				if next == nil {
					ok = false
					return
				}
				// advance the tail, in case it is lagging
				casNodePtr(&Q.Tail, tail, next)
			} else {
				// read the value
				item = next.item

				// try to advance the head to the next node, this is a required condition for dequeue to succeed
				if casNodePtr(&Q.Head, head, next) {
					break
				}
			}
		}
	}
	head.next = head // make old head node as pointing to itself so that it can be GC'd
	return item, true
}

// convenience wrapper over the atomic CAS which hides the dirty work
func casNodePtr[E any](addr **Node[E], old *Node[E], new *Node[E]) (swapped bool) {
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(addr)),
		unsafe.Pointer(old),
		unsafe.Pointer(new),
	)
}
