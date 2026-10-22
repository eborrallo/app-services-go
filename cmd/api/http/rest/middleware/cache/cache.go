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

		// Try to serve from cache. Only replay valid JSON bodies; a corrupt or
		// non-JSON cache entry is skipped (logged) rather than crashing.
		if cachedResponse, err := cacheInstance.Get(cacheKey); err == nil {
			var cachedValue interface{}
			if err := json.Unmarshal(cachedResponse, &cachedValue); err != nil {
				log.Printf("cache: skipping unreadable cached response for %q: %v", cacheKey, err)
			} else {
				c.JSON(http.StatusOK, cachedValue)
				c.Abort()
				return
			}
		}

		// Not served from cache: capture the response while processing.
		blw := &CachedResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// Only cache successful responses, and only when the body is JSON. A
		// non-JSON body is logged and skipped instead of crashing the process.
		if c.Writer.Status() != http.StatusOK {
			return
		}
		var cacheValue interface{}
		if err := json.Unmarshal(blw.body.Bytes(), &cacheValue); err != nil {
			log.Printf("cache: not caching non-JSON response for %q: %v", cacheKey, err)
			return
		}
		if err := cacheInstance.Set(cacheKey, cache.Value(cacheValue)); err != nil {
			log.Printf("cache: failed to store response for %q: %v", cacheKey, err)
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
