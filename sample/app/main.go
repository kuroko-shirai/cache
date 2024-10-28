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
	})
	if err != nil {
		return
	}

	// go newCache.Process()

	newCache.Set(1, "one")
	newCache.Set(2, "two")

	log.Println(newCache.Keys())

	for _, key := range newCache.Keys() {
		log.Println(newCache.Has(key))
		log.Println(newCache.Get(key))
	}
}
