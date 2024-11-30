package main

import (
	"log"
	"time"

	"github.com/kuroko-shirai/cache"
)

func main() {
	newCache, err := cache.New[int32, string](&cache.Config{
		TTL:  150 * time.Millisecond,
		Size: 2,
		CLS:  false,
	})
	if err != nil {
		return
	}

	newCache.Set(1, "one")
	newCache.Set(2, "two")
	newCache.Set(3, "three")

	{
		v, k := newCache.Get(1)
		log.Println(v, k)
	}

	time.Sleep(150 * time.Millisecond)

	{
		v, k := newCache.Get(1)
		log.Println(v, k)
	}
}
