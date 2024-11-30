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

	log.Println(newCache.Keys())

	for _, key := range newCache.Keys() {
		log.Println(newCache.Has(key))
		v, k := newCache.Get(key)
		log.Println(v, k)
	}

	time.Sleep(3 * time.Second)

	log.Println(newCache.Keys())
}
