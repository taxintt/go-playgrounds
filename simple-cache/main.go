package main

import (
	"log"
	"sync"
	"time"
)

const defaultValue = "default"

type Cache struct {
	mu      sync.Mutex
	storage map[string]string
}

func NewCache() *Cache {
	return &Cache{
		storage: make(map[string]string),
	}
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.storage[key] = value
}

func (c *Cache) Get(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.storage[key]

	if ok {
		return value
	}

	go func() {
		value = expensiveCall(key)
		c.Set(key, value)
	}()

	return defaultValue
}

func expensiveCall(key string) string {
	time.Sleep(1 * time.Second)
	return "expensive call"
}

func main() {
	mCache := NewCache()
	log.Println(mCache.Get("key1"))

	time.Sleep(2 * time.Second)
	log.Println(mCache.Get("key1"))
}
