package main

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
	"log"
)

// GLndC8XmRT5o6oBLwn8scDNvFY5MuX78wxJQsW5tXctk
var feePayer, _ = types.AccountFromBase58("D8i1DFhgxWkBC52kRUtDRkZL5J5bUDJXssJtsehWYT51txphk8ipWe8goFKJt6638vAmEHVxdovsjmfiHPvKPbS")

// HcNCxoni2Ln5si48s1w8r5TRVH296RQ1MzKeM9FctdPg
var alice, _ = types.AccountFromBase58("5ob5v9uGyJstuENC4pu7ScCdkCPiAZXqFLhu3qUoHB8uePLchCpPmgWyXPZ24NxLBD8dUP6UFNNYXKsFJtFKie74")

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	// create a mint account
	mint := types.NewAccount()
	fmt.Println("mint:", mint.PublicKey.ToBase58())

	// get rent
	rentExemptionBalance, err := c.GetMinimumBalanceForRentExemption(
		context.Background(),
		token.MintAccountSize,
	)
	if err != nil {
		log.Fatalf("get min balacne for rent exemption, err: %v", err)
	}

	res, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				system.CreateAccount(system.CreateAccountParam{
					From:     feePayer.PublicKey,
					New:      mint.PublicKey,
					Owner:    common.TokenProgramID,
					Lamports: rentExemptionBalance,
					Space:    token.MintAccountSize,
				}),
				token.InitializeMint(token.InitializeMintParam{
					Decimals:   8,
					Mint:       mint.PublicKey,
					MintAuth:   alice.PublicKey,
					FreezeAuth: nil,
				}),
			},
		}),
		Signers: []types.Account{feePayer, mint},
	})
	if err != nil {
		log.Fatalf("generate tx error, err: %v\n", err)
	}

	sig, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send tx error, err: %v\n", err)
	}

	fmt.Printf("check tx at: https://explorer.solana.com/tx/%s?cluster=devnet\n", sig)
}

/*
mint: gYqzga5v1RoVWxtfXizHuoyxUpTnzf9WyrXftTkDfpT
check tx at: https://explorer.solana.com/tx/2zfGVnNkfEWL91emjAfSb3AfFh7XJK4FvnYAiMgibv4iEddjMLsdQGSvjjiGziXVTKBGVWKAMjnKfzXKYx6ZdXHx?cluster=devnet
*/
