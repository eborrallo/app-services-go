package blokchain

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type SignatureVerificator struct {
	eip1271ABI       string
	eip1271MagicLink string
	EMPTY_BYTE       string
	ethersProvider   *ethclient.Client
}

func NewSignatureVerificator(ethersProvider *ethclient.Client) SignatureVerificator {
	return SignatureVerificator{
		eip1271ABI:       "function isValidSignature(bytes32 _message, bytes _signature) public view returns (bytes4 magicValue)",
		eip1271MagicLink: "isValidSignature(bytes32,bytes)",
		EMPTY_BYTE:       "0x",
		ethersProvider:   ethersProvider,
	}
}

func (sv *SignatureVerificator) VerifyMessage(message string, signature string, ethereumAddress string) error {
	code, err := sv.ethersProvider.CodeAt(context.Background(), common.HexToAddress(ethereumAddress), nil)
	if err != nil {
		fmt.Println("Error getting code", err)
		return err
	}

	var valid bool
	if bytes.Equal(code, common.FromHex(sv.EMPTY_BYTE)) {
		// Normal wallet
		sig := hexutil.MustDecode(signature)
		msg := []byte(message)
		msg = accounts.TextHash(msg)
		if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
			sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
		}

		recovered, err := crypto.SigToPub(msg, sig)
		if err != nil {
			return err
		}

		recoveredAddr := crypto.PubkeyToAddress(*recovered)
		valid = strings.EqualFold(ethereumAddress, recoveredAddr.Hex())

	} else {
		// Smart contract wallet
		/*
			contractAddr := common.HexToAddress(ethereumAddress)
			contract, err := NewEIP1271(contractAddr, sv.ethersProvider)
			if err != nil {
				return err
			}

			hash := crypto.Keccak256([]byte(message))
			returnValue, err := contract.IsValidSignature(&bind.CallOpts{}, hash, signature)
			if err != nil {
				// Signature is not valid
				valid = false
			} else {
				valid = bytes.Equal(returnValue[:], common.FromHex(sv.eip1271MagicLink))
			}
		*/
	}

	if !valid {
		return fmt.Errorf("signature not valid")
	}
	return nil
}
