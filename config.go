package cache

import "time"

type Config struct {
	Size int

	Poll time.Duration // frequency of polling items in cache
	TTL  time.Duration // items life-time in the cache

	CLS bool // flag for clear cache
}
