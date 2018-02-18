package xsync

import (
	"sync/atomic"
)

type Counter struct {
	counter uint64
}

func (t *Counter) Get() uint64 {
	return atomic.LoadUint64(&t.counter)
}

func (t *Counter) Add(n uint64) uint64 {
	return atomic.AddUint64(&t.counter, n)
}

func NewCounter() *Counter {
	c := &Counter{}

	return c
}
