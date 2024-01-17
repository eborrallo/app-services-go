package cache_test

import (
	"app-services-go/kit/cache"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
	redis.UniversalClient
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func Test_KeyNotExist(t *testing.T) {
	mockClient := new(MockRedisClient)
	sut := cache.NewRedisCache(mockClient, 1*time.Millisecond)

	key := cache.Key("nonexistentKey")

	// Set up the mock to return redis.Nil for Get
	mockClient.On("Get", mock.Anything, string(key)).Return(redis.NewStringResult("", redis.Nil))

	result, err := sut.Get(key)

	assert.Error(t, err, "not found")
	assert.Nil(t, result)

	mockClient.AssertExpectations(t)
}

func Test_StoreString(t *testing.T) {
	mockClient := new(MockRedisClient)
	sut := cache.NewRedisCache(mockClient, 1*time.Millisecond)

	key := cache.Key("stringKey")
	value := cache.Value("stringValue")

	// Set up the mock for successful Set
	mockClient.On("Set", mock.Anything, string(key), mock.Anything, mock.Anything).Return(redis.NewStatusResult("", nil))

	sut.Set(key, value)

	mockClient.AssertExpectations(t)
}

func Test_StoreStruct(t *testing.T) {
	mockClient := new(MockRedisClient)
	sut := cache.NewRedisCache(mockClient, 1*time.Millisecond)

	key := cache.Key("structKey")
	value := struct {
		Field1 string
		Field2 int
	}{
		Field1: "Hello",
		Field2: 42,
	}

	// Set up the mock for successful Set
	mockClient.On("Set", mock.Anything, string(key), mock.Anything, mock.Anything).Return(redis.NewStatusResult("", nil))

	sut.Set(key, value)

	mockClient.AssertExpectations(t)
}

func Test_ExpiredTime(t *testing.T) {
	mockClient := new(MockRedisClient)
	sut := cache.NewRedisCache(mockClient, 1*time.Millisecond)

	key := cache.Key("expiredKey")
	value := cache.Value("expiredValue")

	// Set up the mock for successful Set
	mockClient.On("Set", mock.Anything, string(key), mock.Anything, mock.Anything).Return(redis.NewStatusResult("", nil))

	sut.Set(key, value)

	// Set up the mock for Get to return redis.Nil after expiration
	mockClient.On("Get", mock.Anything, string(key)).Return(redis.NewStringResult("", redis.Nil))

	// Advance time by more than the expiration duration
	mockClient.On("Get", mock.Anything, string(key)).After(6 * time.Minute).Return(redis.NewStringResult("", nil))

	result, err := sut.Get(key)

	assert.Error(t, err, "not found")
	assert.Nil(t, result)

	mockClient.AssertExpectations(t)
}

func Test_FetchString(t *testing.T) {
	mockClient := new(MockRedisClient)
	sut := cache.NewRedisCache(mockClient, 1*time.Millisecond)

	key := cache.Key("stringKey")
	expectedValue := cache.Value("stringValue")

	// Set up the mock for successful Get
	mockClient.On("Get", mock.Anything, string(key)).Return(redis.NewStringResult(expectedValue.(string), nil))

	result, err := sut.Get(key)

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, string(result))

	mockClient.AssertExpectations(t)
}

func Test_FetchStruct(t *testing.T) {
	mockClient := new(MockRedisClient)
	sut := cache.NewRedisCache(mockClient, 1*time.Millisecond)

	key := cache.Key("structKey")
	expectedValue := struct {
		Field1 string
		Field2 int
	}{
		Field1: "Hello",
		Field2: 42,
	}

	// Set up the mock for successful Get
	jsonValue, _ := json.Marshal(expectedValue)
	mockClient.On("Get", mock.Anything, string(key)).Return(redis.NewStringResult(string(jsonValue), nil))

	result, err := sut.Get(key)
	assert.NoError(t, err)

	var structResult struct {
		Field1 string
		Field2 int
	}
	err = json.Unmarshal(result, &structResult)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, structResult)

	mockClient.AssertExpectations(t)
}
