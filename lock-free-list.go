package main

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type LNode struct {
	value int
	next  unsafe.Pointer
}

type List struct {
	head unsafe.Pointer
}

func (l *List) Add(value int) {
	node := &LNode{value: value}

	for {
		oldHead := atomic.LoadPointer(&l.head)
		node.next = oldHead

		if atomic.CompareAndSwapPointer(&l.head, oldHead, unsafe.Pointer(node)) {
			break
		}
	}
}

func (l *List) Print() {
	curr := atomic.LoadPointer(&l.head)

	for curr != nil {
		node := (*LNode)(curr)
		fmt.Println(node.value)
		curr = atomic.LoadPointer(&node.next)
	}
}
