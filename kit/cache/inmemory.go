package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
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
			// Never kill the process over one bad cached value: report it.
			return nil, fmt.Errorf("error marshalling cached value: %w", err)
		}
		return res, nil
	}
	return nil, errors.New("not found")

}

func (c *InMemoryCache) Set(key Key, value Value) error {
	c.client.Set(string(key), value, cache.DefaultExpiration)
	return nil
}
