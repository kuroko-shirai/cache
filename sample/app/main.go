package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kuroko-shirai/cache"
)

func main() {
	newCache, err := cache.New[int, string](&cache.Config{
		TTL:  150 * time.Millisecond,
		Size: 5,
		CLS:  true,
	})
	if err != nil {
		return
	}

	for i := range 10 {
		newCache.Set(i, fmt.Sprintf("worker-%d", i))
	}

	{
		v, k := newCache.Get(1)
		log.Println(v, k)
		v, k = newCache.Get(7)
		log.Println(v, k)
	}

	time.Sleep(150 * time.Millisecond)

	{
		v, k := newCache.Get(7)
		log.Println(v, k)
	}

	time.Sleep(150 * time.Millisecond)
}
