package cache

import (
	"encoding/json"
	"errors"
	"github.com/patrickmn/go-cache"
	"log"
	"time"
)

const (
	defaultExpiration = 5 * time.Minute
	purgeTime         = 10 * time.Minute
)

type InMemoryCache struct {
	client *cache.Cache
}

func NewInMemoryCache(expirationTime time.Duration) *InMemoryCache {
	if expirationTime == 0 {
		expirationTime = defaultExpiration
	}
	obj := cache.New(expirationTime, purgeTime)
	return &InMemoryCache{
		client: obj,
	}
}

func (c *InMemoryCache) Get(key Key) (item []byte, err error) {
	value, ok := c.client.Get(string(key))

	if ok {
		log.Println("from cache")
		res, err := json.Marshal(value)
		if err != nil {
			log.Fatal("Error")
		}
		return res, nil
	}
	return nil, errors.New("not found")

}

func (c *InMemoryCache) Set(key Key, value Value) {
	c.client.Set(string(key), value, cache.DefaultExpiration)
}
