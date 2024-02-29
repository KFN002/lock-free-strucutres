package main

import (
	"sync/atomic"
	"unsafe"
)

type Node struct {
	value int
	next  *Node
}

type LockFreeQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func NewLockFreeQueue() *LockFreeQueue {
	dummy := &Node{}
	return &LockFreeQueue{
		head: unsafe.Pointer(dummy),
		tail: unsafe.Pointer(dummy),
	}
}

func (q *LockFreeQueue) Enqueue(value int) {
	newNode := &Node{value: value}

	for {
		tail := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&((*Node)(tail)).next)))

		if tail == atomic.LoadPointer(&q.tail) {
			if next == nil {
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&((*Node)(tail)).next)), nil, unsafe.Pointer(newNode)) {
					atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(newNode))
					return
				}
			} else {
				atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(next))
			}
		}
	}
}

func (q *LockFreeQueue) Dequeue() (int, bool) {
	for {
		head := atomic.LoadPointer(&q.head)
		next := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&((*Node)(head)).next)))

		if head == atomic.LoadPointer(&q.head) {
			if next == nil {
				return 0, false
			}
			if atomic.CompareAndSwapPointer(&q.head, head, unsafe.Pointer(next)) {
				return (*Node)(next).value, true
			}
		}
	}
}
