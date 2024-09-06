package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
)

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	// token account address
	getAccountInfoResponse, err := c.GetAccountInfo(context.TODO(), "BdEcBm46DWCEBFXVHwXhW76RLqzyCpaiJMxgveL8dLEm")
	if err != nil {
		log.Fatalf("failed to get account info, err: %v", err)
	}

	tokenAccount, err := token.TokenAccountFromData(getAccountInfoResponse.Data)
	if err != nil {
		log.Fatalf("failed to parse data to a token account, err: %v", err)
	}

	fmt.Printf("%+v\n", tokenAccount)
	// {Mint:gYqzga5v1RoVWxtfXizHuoyxUpTnzf9WyrXftTkDfpT Owner:HcNCxoni2Ln5si48s1w8r5TRVH296RQ1MzKeM9FctdPg Amount:0 Delegate:<nil> State:1 IsNative:<nil> DelegatedAmount:0 CloseAuthority:<nil>}
	// after mint to the account {Mint:gYqzga5v1RoVWxtfXizHuoyxUpTnzf9WyrXftTkDfpT Owner:HcNCxoni2Ln5si48s1w8r5TRVH296RQ1MzKeM9FctdPg Amount:100000000 Delegate:<nil> State:1 IsNative:<nil> DelegatedAmount:0 CloseAuthority:<nil>}
}
