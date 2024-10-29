package main

import (
	"log"
	"time"

	"github.com/kuroko-shirai/cache"
)

func main() {
	newCache, err := cache.New[string](&cache.Config{
		Poll: 5 * time.Second,
		TTL:  10 * time.Second,
		CLS:  true,
	})
	if err != nil {
		return
	}

	newCache.Set(1, "one")
	newCache.Set(2, "two")

	log.Println(newCache.Keys())

	for _, key := range newCache.Keys() {
		log.Println(newCache.Has(key))
		v, k := newCache.Get(key)
		log.Println(v, k)
	}
}
