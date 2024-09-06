package main

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
	"log"
)

// GLndC8XmRT5o6oBLwn8scDNvFY5MuX78wxJQsW5tXctk
var feePayer, _ = types.AccountFromBase58("D8i1DFhgxWkBC52kRUtDRkZL5J5bUDJXssJtsehWYT51txphk8ipWe8goFKJt6638vAmEHVxdovsjmfiHPvKPbS")

// HcNCxoni2Ln5si48s1w8r5TRVH296RQ1MzKeM9FctdPg
var alice, _ = types.AccountFromBase58("5ob5v9uGyJstuENC4pu7ScCdkCPiAZXqFLhu3qUoHB8uePLchCpPmgWyXPZ24NxLBD8dUP6UFNNYXKsFJtFKie74")

var mintPubkey = common.PublicKeyFromString("gYqzga5v1RoVWxtfXizHuoyxUpTnzf9WyrXftTkDfpT")

var aliceTokenATAPubkey = common.PublicKeyFromString("BdEcBm46DWCEBFXVHwXhW76RLqzyCpaiJMxgveL8dLEm")

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	res, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				token.MintToChecked(token.MintToCheckedParam{
					Mint:     mintPubkey,
					Auth:     alice.PublicKey,
					Signers:  []common.PublicKey{},
					To:       aliceTokenATAPubkey,
					Amount:   1e8,
					Decimals: 8,
				}),
			},
		}),
		Signers: []types.Account{feePayer, alice},
	})
	if err != nil {
		log.Fatalf("generate tx error, err: %v\n", err)
	}

	txhash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send raw tx error, err: %v\n", err)
	}

	fmt.Printf("check tx at: https://explorer.solana.com/tx/%s?cluster=devnet\n", txhash)
}

/*
check tx at: https://explorer.solana.com/tx/3TP1UDkbWWKurLmBSCyP6nii6J2vghxsWRtSmANMjPN1NsdNzSoZQhyPjdxQd4TeAkeGyFX1qMu8wVKRchA2QQP5?cluster=devnet
*/
