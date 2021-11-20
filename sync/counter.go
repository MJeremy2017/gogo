package sync

import "sync"

type Counter struct {
	mu sync.Mutex
	value int
}

func (c *Counter) Incr() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += 1
}

func (c *Counter) Value() int {
	return c.value
}