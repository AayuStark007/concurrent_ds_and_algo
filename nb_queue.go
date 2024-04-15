package main

import (
	"sync/atomic"
	"unsafe"
)

type Pointer struct {
	ptr   *Node
	count uint32
}

type Node struct {
	value int
	next  Pointer
}

func NewNode() *Node {
	return &Node{
		next: Pointer{
			ptr:   nil,
			count: 0,
		},
	}
}

/*
* NBQueue is a Concurent Non-Blocking Queue based on CAS primitives
 */
type NBQueue struct {
	Head Pointer
	Tail Pointer
}

func NewQueue() Queue {
	var node *Node = NewNode()

	return &NBQueue{
		Head: Pointer{
			ptr:   node,
			count: 0,
		},
		Tail: Pointer{
			ptr:   node,
			count: 0,
		},
	}
}

func (Q *NBQueue) Enqueue(value int) {
	// creating a new node with value to enqueue
	var node *Node = NewNode()
	node.value = value
	node.next.ptr = nil

	var tail Pointer
	// since CAS-ing can fail, we keep trying until we succeed to enqueue `node`
	for {
		// save the current tail
		tail = Q.Tail
		// save the next node to tail
		next := tail.ptr.next

		// ensure our saved tail is still the Tail
		if tail == Q.Tail {
			// wrt tail -> either we are at the last node or our tail is lagging and we need to advance it further
			if next.ptr == nil { // tail is pointing to the last node
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(
					unsafe.Pointer(&tail.ptr.next.ptr)),
					unsafe.Pointer(next.ptr),
					unsafe.Pointer(node)) {
					break // successfully enqueued
				}
			} else { // tail not pointing to the last node (some concurrent process enqueued more nodes after we read the tail earlier)
				atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&Q.Tail.ptr)),
					unsafe.Pointer(tail.ptr),
					unsafe.Pointer(next.ptr)) // try to move the tail to newly inserted node
			}
		}
	}
	/*
		We have enqueued the node, we can now try to ensure Q.Tail is pointing to newly inserted node.
		But, its not a necessary thing to do, which allows enquques to be fast as we can allow the Q.Tail to lag.
	*/
	atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&Q.Tail.ptr)),
		unsafe.Pointer(tail.ptr),
		unsafe.Pointer(node))
}

func (Q *NBQueue) Dequeue() (value int, status bool) {
	var head Pointer // saving current head (which will be dequeued), which allows us to free the node data

	// since CAS-ing can fail (on account that nodes are concurrently dequeued), we must keep trying
	for {
		head = Q.Head
		// read the current tail, required for checking empty Queue case
		tail := Q.Tail
		// next ptr to head (will become the new head later)
		next := head.ptr.next

		// ensure that our saved head is still the Head
		if head == Q.Head {
			// either Q is empty, or the tail might be lagging (on account of more items enqueued since we started)
			if head.ptr == tail.ptr {
				// this means tail is pointing to the last node, so our Q is empty
				if next.ptr == nil {
					status = false
					return
				}
				// advance the tail, in case it is lagging
				atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&Q.Tail.ptr)),
					unsafe.Pointer(tail.ptr),
					unsafe.Pointer(next.ptr),
				)
			} else {
				// read the value
				value = next.ptr.value

				// try to advance the head to the next node, this is a required condition for dequeue to succeed
				if atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&Q.Head.ptr)),
					unsafe.Pointer(head.ptr),
					unsafe.Pointer(next.ptr),
				) {
					break
				}
			}
		}
	}
	head.ptr = head.ptr // make old head node as pointing to itself so that it can be GC'd
	return value, true
}
