package cache

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"golang.org/x/sync/errgroup"
)

type ct[V any] struct {
	value V

	ttl time.Time
}

type mt[K comparable, V any] struct {
	cls  bool
	ttl  time.Duration
	size int

	m cmap.ConcurrentMap[K, ct[V]]
}

func newCMap[K, V comparable](
	size int,
	ttl time.Duration,
) *mt[K, V] {
	cls := false
	if ttl != time.Duration(0) {
		cls = true
	}

	return &mt[K, V]{
		cls:  cls,
		ttl:  ttl,
		size: size,
		m: cmap.NewWithCustomShardingFunction[K, ct[V]](func(key K) uint32 {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err := enc.Encode(key); err != nil {
				panic(err)
			}
			h := fnv.New32()
			h.Write(buf.Bytes())
			return h.Sum32()
		}),
	}
}

func (c *mt[K, V]) Keys() []K {
	keys := make([]K, 0, c.m.Count())
	for key := range c.m.Items() {
		keys = append(keys, key)
	}

	return keys
}

func (c *mt[K, V]) Get(key K) (V, bool) {
	if c.cls {
		if item, ok := c.m.Get(key); ok {
			if time.Since(item.ttl) <= c.ttl {
				c.Set(key, item.value)

				return item.value, ok
			}
		}

		return *new(V), false
	}
	item, ok := c.m.Get(key)

	return item.value, ok
}

func (c *mt[K, V]) min() (K, bool) {
	var minKey K
	var minValue ct[V]

	first := true
	found := false
	for key, value := range c.m.Items() {
		if first || value.ttl.Before(minValue.ttl) {
			minKey = key
			minValue = value
			first = false
			found = true
		}
	}

	return minKey, found
}

func (c *mt[K, V]) Set(key K, value V) {
	ttl := time.Now()

	if c.size != 0 {
		if c.m.Count() > c.size {
			minKey, found := c.min()
			if found {
				c.m.Remove(minKey)
			}
		}
	}

	c.m.Set(key, ct[V]{
		value: value,
		ttl:   ttl,
	})
}

func (c *mt[K, V]) Process() {
	if c.cls {
		eg := new(errgroup.Group)

		eg.Go(
			func(c *mt[K, V]) func() error {
				return func() error {
					for {
						for _, key := range c.Keys() {
							item, _ := c.m.Get(key)
							if time.Since(item.ttl) > c.ttl {
								c.m.Remove(key)
							}
						}
					}
				}
			}(c),
		)

		go func() {
			eg.Wait()
		}()
	}
}
