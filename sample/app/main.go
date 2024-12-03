package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kuroko-shirai/cache"
)

func main() {
	newCache, err := cache.New[int32, string](&cache.Config{
		TTL:  150 * time.Millisecond,
		Size: 10,
		CLS:  true,
	})
	if err != nil {
		return
	}

	newCache.Set(1, "one")
	newCache.Set(2, "two")
	newCache.Set(3, "three")

	{
		v, k := newCache.Get(1)
		log.Println(v, k, newCache.Size())
	}

	time.Sleep(150 * time.Millisecond)

	{
		v, k := newCache.Get(1)
		log.Println(v, k, newCache.Size())
	}

	fmt.Println(">c", newCache.Has(1))

	time.Sleep(150 * time.Millisecond)

	fmt.Println(">c", newCache.Has(1))
}
