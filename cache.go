package cache

import (
	"log"
	"time"

	"github.com/kuroko-shirai/task"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type Config struct {
	Poll time.Duration // frequency of polling items in cache
	TTL  time.Duration // items life-time in the cache
	CLS  bool          // flag for clear cache
}

type Cache[T any] struct {
	cm   cmap.ConcurrentMap[int32, container[T]]
	ttl  time.Duration
	poll time.Duration
}

type container[T any] struct {
	value T
	ttl   time.Time
}

func New[T any](config *Config) (*Cache[T], error) {
	cm := cmap.NewWithCustomShardingFunction[int32, container[T]](func(key int32) uint32 {
		return uint32(key)
	})

	newCache := &Cache[T]{
		cm:   cm,
		ttl:  config.TTL,
		poll: config.Poll,
	}

	if config.CLS {
		newTask := task.New(
			func(recovery any) {
				log.Printf("got panic: %!w", recovery)
			},
			func(cache *Cache[T]) func() {
				return func() {
					cache.process()
				}
			}(newCache),
		)
		newTask.Do()
	}

	return newCache, nil
}

func (cache *Cache[T]) Set(key int32, value T) {
	cache.cm.Set(key, container[T]{
		value: value,
		ttl:   time.Now(),
	})
}

func (cache *Cache[T]) Keys() []int32 {
	return cache.cm.Keys()
}

func (cache *Cache[T]) Get(key int32) (T, bool) {
	if time.Since(cache.cm.Items()[key].ttl) <= cache.ttl {
		cache.cm.Set(key, container[T]{
			value: cache.cm.Items()[key].value,
			ttl:   time.Now(),
		})
		item, ok := cache.cm.Get(key)

		return item.value, ok
	}

	return *new(T), false
}

// flush - removes an old element from the map by ttl.
func (cache *Cache[T]) flush(key int32) {
	if time.Since(cache.cm.Items()[key].ttl) > cache.ttl {
		cache.cm.Remove(key)
	}
}

// Has - looks up an item under specified key
func (cache *Cache[T]) Has(key int32) bool {
	return cache.cm.Has(key)
}

// process -
func (cache *Cache[T]) process() {
	for {
		for _, key := range cache.cm.Keys() {
			cache.flush(key)
		}

		time.Sleep(cache.poll)
	}
}
