package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type Group[T any] struct {
	group singleflight.Group
}

type Cache struct {
	data map[int]int
	mu   sync.Mutex
}

var group Group[int]

func NewCache() *Cache {
	return &Cache{
		data: make(map[int]int),
	}
}

func (c *Cache) Get(key int) int {
	c.mu.Lock()
	v, ok := c.data[key]
	c.mu.Unlock()

	if ok {
		return v
	}

	// genericsを利用したい場合は、以下のように書く
	vv, err, _ := group.Do(fmt.Sprintf("cache_%d", key), func() (int, error) {
		value, err := heavy(key)
		if err != nil {
			return 0, err
		}
		c.Set(key, value)
		return value, nil
	})

	// 2回目以降のキャッシュの更新処理は初回のキャッシュ更新処理を待つ
	// vv, err, _ := group.Do(fmt.Sprintf("cache_%d", key), func() (interface{}, error) {
	// 	value, err := heavy(key)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	c.Set(key, value)
	// 	return value, nil
	// })

	if err != nil {
		panic(err)
	}
	return vv
}

func (g *Group[T]) Do(key string, fn func() (T, error)) (T, error, bool) {
	v, err, shared := g.group.Do(key, func() (any, error) {
		return fn()
	})
	return v.(T), err, shared
}

func (c *Cache) Set(key, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func heavy(key int) (int, error) {
	log.Printf("heavy: %d", key)
	time.Sleep(1 * time.Second)
	return key * 10, nil
}

func main() {
	mCache := NewCache()

	for i := 0; i < 100; i++ {
		go func(i int) {
			mCache.Get(i)
		}(i)
	}

	time.Sleep(2 * time.Second)

	for i := 0; i < 10; i++ {
		log.Println(mCache.Get(i))
	}
}
