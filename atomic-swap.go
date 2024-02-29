package main

import (
	"sync/atomic"
)

func AtomicSwap(a *int32, b *int32) {
	temp := atomic.SwapInt32(a, atomic.LoadInt32(b))
	atomic.StoreInt32(b, temp)
}
