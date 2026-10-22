package crypt

import (
	"app-services-go/configs"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// keyFunc returns the HMAC secret only after asserting the token was signed
// with an HMAC method. Without this check a token could claim "alg":"none" or a
// different algorithm (algorithm-confusion attack).
func keyFunc(secretKey []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	}
}

// CreateToken signs payload into an HS256 JWT. It returns an error instead of
// panicking. Callers must pass a minimal claims value, never the whole User.
func CreateToken(payload interface{}, expirationTimeInHours time.Duration) (string, error) {
	config, err := configs.GetCryptConfig()
	if err != nil {
		return "", err
	}
	secretKey := []byte(config.Secret)

	p, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"payload": string(p),
			"exp":     time.Now().Add(expirationTimeInHours).Unix(),
		})

	return token.SignedString(secretKey)
}

func VerifyToken(tokenString string) error {
	config, err := configs.GetCryptConfig()
	if err != nil {
		return err
	}
	secretKey := []byte(config.Secret)

	token, err := jwt.Parse(tokenString, keyFunc(secretKey))
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("failed to parse claims")
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("token is expired")
	}

	return nil
}

func GetPayloadFromToken(tokenString string, payload interface{}) error {
	config, err := configs.GetCryptConfig()
	if err != nil {
		return err
	}
	secretKey := []byte(config.Secret)

	token, err := jwt.Parse(tokenString, keyFunc(secretKey))
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("failed to parse claims")
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("token is expired")
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
