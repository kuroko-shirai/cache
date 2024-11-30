package cache

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/kuroko-shirai/task"
)

type Cache[K, V comparable] struct {
	cm map[K]container[V]

	ttl  time.Duration
	poll time.Duration

	mu sync.Mutex
}

type container[V any] struct {
	value V
	ttl   time.Time
}

func New[K, V comparable](config *Config) (*Cache[K, V], error) {
	cm := make(map[K]container[V])

	if config.CLS && config.TTL <= config.Poll {
		return nil, errors.New("the polling frequency should not exceed the lifetime of the cache entry")
	}

	c := &Cache[K, V]{
		cm:   cm,
		ttl:  config.TTL,
		poll: config.Poll,
		mu:   sync.Mutex{},
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
						time.Sleep(c.poll)
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

func (c *Cache[K, V]) Keys() []K {
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

	if time.Since(c.cm[key].ttl) <= c.ttl {
		c.cm[key] = container[V]{
			value: c.cm[key].value,
			ttl:   time.Now(),
		}

		item, ok := c.cm[key]

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
