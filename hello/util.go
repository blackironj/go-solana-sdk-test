package hello

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

func GetConnection() *client.Client {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	resp, err := c.GetVersion(context.TODO())
	if err != nil {
		log.Fatalf("failed to version info, err: %v", err)
	}

	fmt.Println("Solana version", resp.SolanaCore)

	return c
}

func GetAccountFromFile(fileDir string) types.Account {
	keyBytes, err := ioutil.ReadFile(fileDir)
	if err != nil {
		log.Fatalf("fail to read keypair file, err: %v", err)
	}

	dataBytes := make([]byte, 0)
	if err := json.Unmarshal(keyBytes, &dataBytes); err != nil {
		log.Fatalf("fail to unmarshal json, err: %v", err)
	}

	a, err := types.AccountFromBytes(dataBytes)
	if err != nil {
		log.Fatalf("fail to parse keypair json, err: %v", err)
	}

	return a
}
