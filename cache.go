package cache

import (
	"log"
	"sync"
	"time"

	"github.com/kuroko-shirai/task"
)

type Cache[K, V comparable] struct {
	cls bool

	cm map[K]container[V]

	ttl time.Duration

	mu sync.Mutex
}

type container[V any] struct {
	value V
	ttl   time.Time
}

func New[K, V comparable](config *Config) (*Cache[K, V], error) {
	cm := make(map[K]container[V])

	c := &Cache[K, V]{
		cls: config.CLS,
		cm:  cm,
		ttl: config.TTL,
		mu:  sync.Mutex{},
	}

	if config.CLS {
		g := task.WithRecover(
			func(r any, args ...any) {
				log.Println("panic:", r)
			},
		)

		g.Do(
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
			if err := g.Wait(); err != nil {
				log.Printf("errors: %s", err.Error())
			}
		}()
	}

	return c, nil
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.cm[key] = container[V]{
		value: value,
		ttl:   time.Now(),
	}
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
				c.Set(key, item.value)

				return item.value, ok
			}
			return *new(V), false
		}
		return item.value, ok
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
