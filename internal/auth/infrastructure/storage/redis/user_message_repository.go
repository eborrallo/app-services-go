package redis

import (
	"app-services-go/kit/cache"
	"encoding/json"
	"strings"
)

type UserMessageRepository struct {
	rd *cache.RedisCache
}

func NewUserMessageRepository(rd *cache.RedisCache) UserMessageRepository {
	return UserMessageRepository{
		rd: rd,
	}
}

type UserMessage struct {
	Message string
}

func (r UserMessageRepository) SaveMessage(address string, message string) {
	key := cache.Key(strings.ToLower(address))
	value := cache.Value(UserMessage{Message: message})
	r.rd.Set(key, value)
}

func (r UserMessageRepository) GetMessage(address string) (string, error) {
	key := cache.Key(strings.ToLower(address))
	result, err := r.rd.Get(key)
	if err != nil {
		return "", err
	}
	var userMessage = UserMessage{}
	err = json.Unmarshal(result, &userMessage)
	return userMessage.Message, err
}
