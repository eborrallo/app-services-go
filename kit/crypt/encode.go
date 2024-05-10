package crypt

import (
	"app-services-go/configs"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

func Md5(txt string) string {
	hash := md5.Sum([]byte(txt))
	return hex.EncodeToString(hash[:])
}
func Encode(txt string) (string, error) {

	txtEncoded, error := bcrypt.GenerateFromPassword([]byte(txt), bcrypt.MinCost)

	if error != nil {
		return "", error
	}
	return string(txtEncoded[:]), nil

}
func Compare(txt string, txtEncoded string) bool {

	error := bcrypt.CompareHashAndPassword([]byte(txtEncoded), []byte(txt))

	return error == nil

}

func Encrypt(text string) string {
	config, err := configs.GetCryptConfig()
	if err != nil {
		panic(err)
	}
	key := []byte(config.Secret)

	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

func Decrypt(cryptoText string) string {
	config, err := configs.GetCryptConfig()
	if err != nil {
		panic(err)
	}
	key := []byte(config.Secret)

	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
