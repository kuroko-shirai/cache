package cache

import "time"

type Config struct {
	CLS  bool // flag for clear cache
	Size int  // size of service's cache
	Mode int

	TTL time.Duration // items life-time in the cache
}
