package cache

import (
	"app-services-go/kit/cache"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

// Middleware creates a Gin middleware for caching responses.
func Middleware(redis redis.UniversalClient, cacheDuration time.Duration) gin.HandlerFunc {
	cacheInstance := cache.NewRedisCache(redis, cacheDuration)
	return func(c *gin.Context) {
		// Generate a cache key based on the request URL
		cacheKey := cache.Key(c.Request.URL.String())
		// Try to retrieve the response from the cache
		cachedResponse, err := cacheInstance.Get(cacheKey)
		// If not found in the cache, proceed with the request and capture the response
		blw := &CachedResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		if err == nil {
			// If the response is found in the cache, serve it and return
			var cachedValue interface{}

			err = json.Unmarshal(cachedResponse, &cachedValue)
			if err != nil {
				log.Fatal(err)
			}
			c.JSON(http.StatusOK, cachedValue)
			c.Abort()

			return
		}
		// Continue processing the request
		c.Next()

		// Store the response in the cache if it was successful
		if c.Writer.Status() == http.StatusOK {
			var cacheValue interface{}
			err = json.Unmarshal(blw.body.Bytes(), &cacheValue)
			if err != nil {
				log.Fatal(err)
			}
			cacheInstance.Set(cacheKey, cache.Value(cacheValue))
		}
	}
}

type CachedResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CachedResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
