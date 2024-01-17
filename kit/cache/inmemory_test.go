package cache

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

func Test_KeyNotExist(t *testing.T) {
	cache := NewInMemoryCache(0)

	keyNotExist := Key("nonExistentKey")
	_, err := cache.Get(keyNotExist)
	if err == nil || err.Error() != "not found" {
		t.Errorf("Expected 'not found' error for a non-existent key, but got: %v", err)
	}
}

func Test_StoreString(t *testing.T) {
	cache := NewInMemoryCache(0)

	keyString := Key("stringKey")
	valueString := Value("stringValue")
	cache.Set(keyString, valueString)
	resultString, err := cache.Get(keyString)
	if err != nil {
		t.Errorf("Error retrieving string value from cache: %v", err)
	}
	resultStringStr, _ := strconv.Unquote(string(resultString))
	expectedResultString := "stringValue"
	if resultStringStr != expectedResultString {
		t.Errorf("Unexpected result for string. Expected: %s, Got: %s", expectedResultString, resultString)
	}
}

func Test_StoreStruct(t *testing.T) {
	cache := NewInMemoryCache(0)

	type TestStruct struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	keyStruct := Key("structKey")
	valueStruct := TestStruct{Field1: "value1", Field2: 42}
	cache.Set(keyStruct, valueStruct)
	resultStruct, err := cache.Get(keyStruct)
	if err != nil {
		t.Errorf("Error retrieving struct value from cache: %v", err)
	}
	var expectedResultStruct TestStruct
	if err := json.Unmarshal(resultStruct, &expectedResultStruct); err != nil {
		t.Errorf("Error unmarshaling struct result: %v", err)
	}
	if expectedResultStruct != valueStruct {
		t.Errorf("Unexpected result for struct. Expected: %+v, Got: %+v", valueStruct, expectedResultStruct)
	}
}

func Test_ExpiredTime(t *testing.T) {
	cache := NewInMemoryCache(1 * time.Millisecond)

	keyExpired := Key("expiredKey")
	valueExpired := Value("expiredValue")
	cache.Set(keyExpired, valueExpired)

	// Sleep for a while to allow cache to expire
	time.Sleep(1 * time.Second)

	_, err := cache.Get(keyExpired)
	if err == nil || err.Error() != "not found" {
		t.Errorf("Expected 'not found' error after expiration, but got: %v", err)
	}
}
