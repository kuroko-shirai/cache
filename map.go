package cache

import (
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type ct[V any] struct {
	value V

	ttl time.Time
}

type cmap[K, V comparable] struct {
	cls  bool
	ttl  time.Duration
	size int

	m map[K]ct[V]

	mu  sync.Mutex
	smu sync.Mutex
}

func newCMap[K, V comparable](
	size int,
	cls bool,
	ttl time.Duration,
) *cmap[K, V] {
	return &cmap[K, V]{
		cls:  cls,
		ttl:  ttl,
		size: size,
		m:    make(map[K]ct[V], size),
		mu:   sync.Mutex{},
		smu:  sync.Mutex{},
	}
}

func (c *cmap[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := make([]K, 0, len(c.m))
	for key := range c.m {
		keys = append(keys, key)
	}

	return keys
}

func (c *cmap[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.put(key, value)
}

func (c *cmap[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cls {
		if item, ok := (c.m)[key]; ok {
			if time.Since(item.ttl) <= c.ttl {
				c.put(key, item.value)

				return item.value, ok
			}
		}

		return *new(V), false
	}
	item, ok := (c.m)[key]

	return item.value, ok
}

func (c *cmap[K, V]) min() (K, bool) {
	var minKey K
	var minValue ct[V]

	first := true
	found := false
	for key, value := range c.m {
		if first || value.ttl.Before(minValue.ttl) {
			minKey = key
			minValue = value
			first = false
			found = true
		}
	}

	return minKey, found
}

func (c *cmap[K, V]) put(key K, value V) {
	ttl := time.Now()

	if c.size != 0 {
		if len(c.m) > c.size {
			minKey, found := c.min()
			if found {
				delete(c.m, minKey)
			}
		}
	}

	c.m[key] = ct[V]{
		value: value,
		ttl:   ttl,
	}
}

func (c *cmap[K, V]) Process() {
	if c.cls {
		eg := new(errgroup.Group)

		eg.Go(
			func(c *cmap[K, V]) func() error {
				return func() error {
					for {
						c.smu.Lock()

						for _, key := range c.Keys() {
							c.mu.Lock()

							if time.Since((c).m[key].ttl) > c.ttl {
								delete((c).m, key)
							}

							c.mu.Unlock()
						}

						c.smu.Unlock()
					}
				}
			}(c),
		)

		go func() {
			eg.Wait()
		}()
	}
}
