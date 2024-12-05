package cache

import "time"

type Config struct {
	Size int // size of service's cache

	TTL time.Duration // items life-time in the cache
}
