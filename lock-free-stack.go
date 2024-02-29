package main

import (
	"sync/atomic"
	"unsafe"
)

type SNode struct {
	value int
	next  unsafe.Pointer
}

type Stack struct {
	top unsafe.Pointer
}

func (s *Stack) Push(value int) {
	node := &SNode{value: value}

	for {
		oldTop := atomic.LoadPointer(&s.top)
		node.next = oldTop

		if atomic.CompareAndSwapPointer(&s.top, oldTop, unsafe.Pointer(node)) {
			break
		}
	}
}

func (s *Stack) Pop() (int, bool) {
	for {
		oldTop := atomic.LoadPointer(&s.top)
		if oldTop == nil {
			return 0, false
		}

		newTop := (*SNode)(oldTop).next

		if atomic.CompareAndSwapPointer(&s.top, oldTop, unsafe.Pointer(newTop)) {
			return (*SNode)(oldTop).value, true
		}
	}
}
