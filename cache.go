package cache

type st[K, V comparable] interface {
	Set(K, V)
	Get(K) (V, bool)
	Process()
}

type Cache[K, V comparable] struct {
	s st[K, V]
}

func New[K, V comparable](config *Config) (*Cache[K, V], error) {
	var s st[K, V]

	s = newCMap[K, V](config.Size, config.CLS, config.TTL)

	c := &Cache[K, V]{
		s: s,
	}

	c.s.Process()

	return c, nil
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.s.Set(key, value)
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	value, ok := c.s.Get(key)

	return value, ok
}
