package main

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/blocto/solana-go-sdk/rpc"
)

type Token struct {
	Address  string // mint address
	Decimals uint8
	Symbol   string
	Name     string
}

func main() {
	c := client.NewClient(rpc.MainnetRPCEndpoint)
	token, err := newToken(c, "So11111111111111111111111111111111111111112")
	fmt.Println(token, err)
}

/*
output: &{So11111111111111111111111111111111111111112 9 SOL Wrapped SOL} <nil>
*/

func newToken(c *client.Client, mintAddress string) (*Token, error) {
	account, err := c.GetAccountInfo(context.Background(), mintAddress)
	if err != nil {
		return nil, err
	}

	// The decimals are stored at byte offset 44 in the mint account data
	if len(account.Data) < 45 {
		return nil, fmt.Errorf("invalid mint account data")
	}

	decimals := account.Data[44]

	// Get metadata
	mintPubKey := common.PublicKeyFromString(mintAddress)
	metadataAddress, err := token_metadata.GetTokenMetaPubkey(mintPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find metadata PDA: %v", err)
	}

	metadataAccountInfo, err := c.GetAccountInfo(context.Background(), metadataAddress.ToBase58())
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata account info: %v", err)
	}

	metadata, err := token_metadata.MetadataDeserialize(metadataAccountInfo.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize metadata: %v", err)
	}

	return &Token{
		Address:  mintAddress,
		Decimals: decimals,
		Symbol:   metadata.Data.Symbol,
		Name:     metadata.Data.Name,
	}, nil
}
