package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Crypy_Encode(t *testing.T) {

	txt := "123456"
	txtEncoded, error := Encode(txt)

	assert.NoError(t, error)

	assert.True(t, Compare(txt, txtEncoded))
}

func Test_Crypy_Encrypt(t *testing.T) {

	txt := "123456"
	txtEncoded := Encrypt(txt)
	assert.NotEmpty(t, txtEncoded)

	assert.NotEqual(t, txt, txtEncoded)
	txtDecoded := Decrypt(txtEncoded)
	assert.Equal(t, txt, txtDecoded)

}
