package xsync

import (
	"testing"

	"context"

	"golang.org/x/sync/errgroup"
)

var n uint64 = 500

func TestAtomicCounter(t *testing.T) {
	c := NewCounter()
	t.Log(c.Get())

	g, _ := errgroup.WithContext(context.Background())

	var i uint64
	for i = 0; i < n; i++ {
		g.Go(func() error {
			c.Add(1)

			return nil
		})
	}

	g.Wait()

	if c.Get() != n {
		t.Fatalf("expected %d received %d", n, c.Get())
	}

	t.Log(c.Get())
}

func TestCounter(t *testing.T) {
	c := 0
	t.Log(c)

	g, _ := errgroup.WithContext(context.Background())

	var i uint64
	for i = 0; i < n; i++ {
		g.Go(func() error {
			c += 1

			return nil
		})
	}

	g.Wait()

	t.Log(c)
}
