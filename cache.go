package cache

import (
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	WithRewritting = iota
	WithCancelling
)

type Cache[K, V comparable] struct {
	cls  bool
	size int
	ttl  time.Duration

	cm map[K]container[V]
	tm slice[K]

	mu sync.Mutex
}

type container[V any] struct {
	value V
	ttl   time.Time
}

func New[K, V comparable](config *Config) (*Cache[K, V], error) {
	cm := make(map[K]container[V], config.Size)
	tm := make(slice[K], 0, config.Size)

	c := &Cache[K, V]{
		cls:  config.CLS,
		size: config.Size,
		ttl:  config.TTL,
		cm:   cm,
		tm:   tm,
		mu:   sync.Mutex{},
	}

	if config.CLS {
		eg := new(errgroup.Group)

		eg.Go(
			func(c *Cache[K, V]) func() error {
				return func() error {
					for {
						c.mu.Lock()
						for key := range c.cm {
							if time.Since(c.cm[key].ttl) > c.ttl {
								delete(c.cm, key)
							}
						}

						c.mu.Unlock()
					}
				}
			}(c),
		)

		go func() {
			eg.Wait()
		}()
	}

	return c, nil
}

func (c *Cache[K, V]) Set(key K, value V) {
	t := time.Now()

	c.cm[key] = container[V]{
		value: value,
		ttl:   t,
	}

	c.tm.append(key)
}

func (c *Cache[K, V]) keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := make([]K, 0, len(c.cm))
	for key := range c.cm {
		keys = append(keys, key)
	}

	return keys
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.cm[key]; ok {
		if c.cls {
			if time.Since(item.ttl) <= c.ttl {
				c.tm.delete(key)

				c.Set(key, item.value)

				return item.value, ok
			}

			return *new(V), false
		}

		return item.value, true
	}

	return *new(V), false
}

// Has - looks up an item under specified key
func (c *Cache[K, V]) Has(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.cm[key]

	return ok
}

func (c *Cache[K, V]) Size() int {
	c.size = len(c.cm)

	return c.size
}
