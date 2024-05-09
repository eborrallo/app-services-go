package crypt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Payload struct {
	Field_one string `json:"field_one"`
	Field_two string `json:"field_two"`
}

func Test_JWT_Encode(t *testing.T) {
	payload := Payload{
		Field_one: "value_one",
		Field_two: "value_two",
	}
	// Create a new JWT token
	token := CreateToken(payload, 10*time.Minute)

	// Verify that the token is not empty
	assert.NotEmpty(t, token)

	// Verify that the token can be decoded successfully
	decodedPayload := Payload{}
	err := GetPayloadFromToken(token, &decodedPayload)
	assert.NoError(t, err)
	// Verify that the decoded token contains the expected claims
	assert.Equal(t, payload.Field_one, decodedPayload.Field_one)
	assert.Equal(t, payload.Field_two, decodedPayload.Field_two)
}
