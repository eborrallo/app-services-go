package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Crypy_Encode(t *testing.T) {

	txt := "123456"
	txtEncoded, err := Encode(txt)

	assert.NoError(t, err)

	assert.True(t, Compare(txt, txtEncoded))
}

func Test_Crypy_Encrypt(t *testing.T) {
	t.Setenv("SECRET", "test-secret-0123456789")

	txt := "123456"
	txtEncoded, err := Encrypt(txt)
	assert.NoError(t, err)
	assert.NotEmpty(t, txtEncoded)

	assert.NotEqual(t, txt, txtEncoded)
	txtDecoded, err := Decrypt(txtEncoded)
	assert.NoError(t, err)
	assert.Equal(t, txt, txtDecoded)

}
