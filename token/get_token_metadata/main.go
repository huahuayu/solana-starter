package main

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/blocto/solana-go-sdk/rpc"
	"log"
)

const (
	USDCMintAddress = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
)

func GetTokenMetadata(c *client.Client, mintAddress string) (*token_metadata.Metadata, error) {
	mintPubKey := common.PublicKeyFromString(mintAddress)
	metadataAddress, err := token_metadata.GetTokenMetaPubkey(mintPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find metadata PDA: %v", err)
	}

	accountInfo, err := c.GetAccountInfo(context.Background(), metadataAddress.ToBase58())
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %v", err)
	}

	metadata, err := token_metadata.MetadataDeserialize(accountInfo.Data)
	if err != nil {
		log.Fatalf("failed to parse metaAccount, err: %v", err)
	}
	//spew.Dump(metadata)

	return &metadata, nil
}

func main() {
	c := client.NewClient(rpc.MainnetRPCEndpoint)

	tokenMetadata, err := GetTokenMetadata(c, USDCMintAddress)
	if err != nil {
		log.Fatalf("failed to retrieve token metadata: %v", err)
	}
	fmt.Printf("Token Symbol: %s, Token Name: %s\n", tokenMetadata.Data.Symbol, tokenMetadata.Data.Name)
}

/*
Token Symbol: USDC, Token Name: USD Coin
*/
