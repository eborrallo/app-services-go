package blokchain

import (
	"app-services-go/configs"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func RpcConection() *ethclient.Client {
	conf, err := configs.GetBlokchainConfig()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	client, err := rpc.Dial(conf.RpcUrl)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return ethclient.NewClient(client)
}
