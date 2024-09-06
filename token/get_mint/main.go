package main

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
	"log"
)

var mintPubkey = common.PublicKeyFromString("gYqzga5v1RoVWxtfXizHuoyxUpTnzf9WyrXftTkDfpT")

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	getAccountInfoResponse, err := c.GetAccountInfo(context.TODO(), mintPubkey.ToBase58())
	if err != nil {
		log.Fatalf("failed to get account info, err: %v", err)
	}

	mintAccount, err := token.MintAccountFromData(getAccountInfoResponse.Data)
	if err != nil {
		log.Fatalf("failed to parse data to a mint account, err: %v", err)
	}

	fmt.Printf("%+v\n", mintAccount)
	// {MintAuthority:HcNCxoni2Ln5si48s1w8r5TRVH296RQ1MzKeM9FctdPg Supply:0 Decimals:8 IsInitialized:true FreezeAuthority:<nil>}
}
