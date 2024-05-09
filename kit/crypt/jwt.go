package crypt

import (
	"app-services-go/configs"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(payload interface{}, expirationTimeInHours time.Duration) string {
	config, err := configs.GetCryptConfig()
	if err != nil {
		panic(err)
	}
	secretKey := []byte(config.Secret)

	p, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"payload": string(p),
			"exp":     time.Now().Add(expirationTimeInHours).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func VerifyToken(tokenString string) error {
	config, err := configs.GetCryptConfig()
	if err != nil {
		panic(err)
	}
	secretKey := []byte(config.Secret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func GetPayloadFromToken(tokenString string, payload interface{}) error {
	config, err := configs.GetCryptConfig()
	if err != nil {
		panic(err)
	}
	secretKey := []byte(config.Secret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		return secretKey, nil
	})

	if err != nil {
		return err
	}
	/*
		if !token.Valid {
			return nil, fmt.Errorf("invalid token")
		}
	*/

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("failed to parse claims")
	}

	payloadJSON, ok := claims["payload"].(string)
	if !ok {
		return fmt.Errorf("payload is not a string")
	}

	if err := json.Unmarshal([]byte(payloadJSON), payload); err != nil {
		return err
	}

	return nil

}
