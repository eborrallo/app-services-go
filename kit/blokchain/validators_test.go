package blokchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Blokchain_Validators_ValidAddress_Error(t *testing.T) {

	address := "0x1234567890"
	result := IsValidEthereumAddress(address)
	assert.False(t, result)

	address = "123"
	result = IsValidEthereumAddress(address)
	assert.False(t, result)
}

func Test_Blokchain_Validators_ValidAddress_Success(t *testing.T) {

	address := "0x1234567890123456789012345678901234567890"
	result := IsValidEthereumAddress(address)
	assert.True(t, result)

}
