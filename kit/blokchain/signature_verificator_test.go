package blokchain

import (
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_Blokchain_Verificator_Signature(t *testing.T) {

	// Connect to Ethereum node
	client, err := rpc.Dial("https://eth.public-rpc.com")
	if err != nil {
		log.Fatal(err)
	}
	ethClient := ethclient.NewClient(client)

	// Instantiate the SignatureVerificator
	sv := NewSignatureVerificator(ethClient)
	message := "Hello"
	signatureHash := "0xcb9ae25f658178a16ef6f2c37f984ae2651e967d7a3a168b05d4b75d83ba6ec01be6cd12a2235824b866ce0566048491ea6732aa4f581d785448564a3f67bee31b"
	signer := "0x2bbd10fc8793c35f25de6f25ed1d6cff1402b473"
	// Verify the message signature
	err = sv.VerifyMessage(message, signatureHash, signer)
	assert.Nil(t, err)

}
