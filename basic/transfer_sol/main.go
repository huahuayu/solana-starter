package main

import (
	"context"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/mr-tron/base58"
	"log"
)

// HcNCxoni2Ln5si48s1w8r5TRVH296RQ1MzKeM9FctdPg
var alice, _ = types.AccountFromBase58("5ob5v9uGyJstuENC4pu7ScCdkCPiAZXqFLhu3qUoHB8uePLchCpPmgWyXPZ24NxLBD8dUP6UFNNYXKsFJtFKie74")

// GLndC8XmRT5o6oBLwn8scDNvFY5MuX78wxJQsW5tXctk
var feePayer, _ = types.AccountFromBase58("D8i1DFhgxWkBC52kRUtDRkZL5J5bUDJXssJtsehWYT51txphk8ipWe8goFKJt6638vAmEHVxdovsjmfiHPvKPbS")

var frank = types.NewAccount()

// Transfer 0.1 SOL from alice to frank, using feePayer to pay for the transaction fee
func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	// log alice account
	log.Printf("frank account: %v, private key: %v\n", frank.PublicKey.ToBase58(), base58.Encode(frank.PrivateKey))

	// to fetch recent blockHash
	recentBlockHashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	// create a transfer tx
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{feePayer, alice},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: recentBlockHashResponse.Blockhash,
			Instructions: []types.Instruction{
				system.Transfer(system.TransferParam{
					From:   alice.PublicKey,
					To:     frank.PublicKey,
					Amount: 1e8, // 0.1 SOL
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a transaction, err: %v", err)
	}

	// send tx
	sig, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	log.Println("signature:", sig)
}

/*
2024/09/04 15:56:38 frank account: 6c7QhVGAvyoGa13gE28jppWGmQ3zEYbvjFxgpVYqE3XL, private key: 4qb6d5onbSDRiR6HzfM4nD8vwuEEG68ZEDQ8626668vUKC33YtYXS9Y9M8K9khDr9ECAH8yrdUsHYYWKY1UUH73t
2024/09/04 15:56:39 signature: 4ZBJsUk3hUKgwE3onS87aJfvNZ38aNEYttzpTbcv6Eew7oUhifdaQG11a7DhBNqkcU1CQgB43pWKXSL5op1LdrGf
*/
