package main

import (
	"log"
	"time"

	"github.com/kuroko-shirai/cache"
)

func main() {
	newCache, err := cache.New[int32, string](&cache.Config{
		Poll: 50 * time.Millisecond,
		TTL:  100 * time.Millisecond,
		CLS:  true,
	})
	if err != nil {
		return
	}

	newCache.Set(1, "one")
	newCache.Set(2, "two")
	newCache.Set(3, "three")

	v, k := newCache.Get(1)
	log.Println(v, k)

	time.Sleep(3 * time.Second)
}
