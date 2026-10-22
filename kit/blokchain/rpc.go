package blokchain

import (
	"app-services-go/configs"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func RpcConection() (*ethclient.Client, error) {
	conf, err := configs.GetBlokchainConfig()
	if err != nil {
		return nil, err
	}

	client, err := rpc.Dial(conf.RpcUrl)
	if err != nil {
		return nil, err
	}
	return ethclient.NewClient(client), nil
}
