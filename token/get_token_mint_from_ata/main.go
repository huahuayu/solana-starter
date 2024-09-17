package main

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
)

func main() {
	c := client.NewClient(rpc.MainnetRPCEndpoint)
	mint, err := getTokenMintFromATA(c, common.PublicKeyFromString("3mHBG2nm6Y9inWayRE7qgfeYMocaoZScfAxizWf19zrS"))
	if err != nil {
		return
	}
	fmt.Println(mint)
}

// output: gXduukdwXJbVw1AjpPcnzmiPFxFHTPSE8yL74LUDfgC

func getTokenMintFromATA(c *client.Client, ataAddress common.PublicKey) (common.PublicKey, error) {
	accountInfo, err := c.GetAccountInfo(context.TODO(), ataAddress.ToBase58())
	if err != nil {
		return common.PublicKey{}, err
	}

	ataInfo, err := token.TokenAccountFromData(accountInfo.Data)
	if err != nil {
		return common.PublicKey{}, fmt.Errorf("failed to parse token account data: %w", err)
	}

	return ataInfo.Mint, nil
}
