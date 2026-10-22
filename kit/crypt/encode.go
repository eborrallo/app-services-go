package crypt

import (
	"app-services-go/configs"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// Md5 is kept only for backwards compatibility with old data. Do NOT use it for
// passwords: it is unsalted and fast. Use Encode/Compare (bcrypt) instead.
func Md5(txt string) string {
	hash := md5.Sum([]byte(txt))
	return hex.EncodeToString(hash[:])
}

// Encode hashes a secret (e.g. a password) with bcrypt at the default cost.
func Encode(txt string) (string, error) {
	txtEncoded, err := bcrypt.GenerateFromPassword([]byte(txt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(txtEncoded), nil
}

// Compare reports whether txt matches the bcrypt-hashed txtEncoded.
func Compare(txt string, txtEncoded string) bool {
	return bcrypt.CompareHashAndPassword([]byte(txtEncoded), []byte(txt)) == nil
}

// aesKey derives a fixed 32-byte AES-256 key from the configured secret. This
// keeps the key length always valid regardless of the secret's length, and
// means the AES key is not the raw JWT signing secret.
func aesKey(secret string) []byte {
	sum := sha256.Sum256([]byte("aes:" + secret))
	return sum[:]
}

// Encrypt returns a base64url AES-CFB ciphertext of text. It returns an error
// instead of panicking on misconfiguration.
func Encrypt(text string) (string, error) {
	config, err := configs.GetCryptConfig()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey(config.Secret))
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)
	// The IV needs to be unique, but not secure. Prepend it to the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt reverses Encrypt. It returns an error (never panics) on malformed
// input, so untrusted tokens from the URL cannot crash the request.
func Decrypt(cryptoText string) (string, error) {
	config, err := configs.GetCryptConfig()
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", fmt.Errorf("invalid token encoding: %w", err)
	}

	block, err := aes.NewCipher(aesKey(config.Secret))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("invalid token: ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
