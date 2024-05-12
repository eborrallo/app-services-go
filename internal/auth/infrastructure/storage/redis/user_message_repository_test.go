package redis

import (
	"app-services-go/configs"
	"app-services-go/kit/cache"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func Test_UserMessageRepository_Save(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	port, _ := strconv.Atoi(mr.Port())
	connection := cache.RedisConnection(configs.RedisConfig{
		RedisHost: mr.Host(),
		RedisPort: port,
	})
	redisCache := cache.NewRedisCache(connection, 10*time.Second)
	repository := NewUserMessageRepository(redisCache)

	repository.SaveMessage("0x123", "Hello World")

	result, err := repository.GetMessage("0x123")
	if err != nil {
		log.Println(err)
	}
	assert.Equal(t, "Hello World", result)
	mr.Close()

}
