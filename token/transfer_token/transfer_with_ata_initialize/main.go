package main

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/program/associated_token_account"
	"github.com/mr-tron/base58"
	"log"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
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

	newAccount := types.NewAccount()
	log.Println("new account:", newAccount.PublicKey.ToBase58(), base58.Encode(newAccount.PrivateKey))

	ata, _, err := common.FindAssociatedTokenAddress(newAccount.PublicKey, mintPubkey)
	if err != nil {
		log.Fatalf("find ata error, err: %v", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				associated_token_account.Create(associated_token_account.CreateParam{
					Funder:                 feePayer.PublicKey,
					Owner:                  newAccount.PublicKey,
					Mint:                   mintPubkey,
					AssociatedTokenAccount: ata,
				}),
				token.TransferChecked(token.TransferCheckedParam{
					From:     aliceTokenATAPubkey,
					To:       ata,
					Mint:     mintPubkey,
					Auth:     alice.PublicKey,
					Signers:  []common.PublicKey{},
					Amount:   1e7,
					Decimals: 8,
				}),
			},
		}),
		Signers: []types.Account{feePayer, alice},
	})
	if err != nil {
		log.Fatalf("failed to new tx, err: %v", err)
	}

	txhash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send raw tx error, err: %v\n", err)
	}

	fmt.Printf("check tx at: https://explorer.solana.com/tx/%s?cluster=devnet\n", txhash)
}

/*
2024/09/06 11:37:28 new account: DfHYap9MpUyNjEkVKUC8jwwy7TrsPWrWnaAMzdZUroAg 59UvmLfepKgEoUdGWjjachF3neFuqRZkrmo6BgBn1EMWg7squTRZxuCB3L9zRT8bBzdU9QaLHhR5byEroVFx4fC4
check tx at: https://explorer.solana.com/tx/4Wo87ndvevGXEprhu4UwBHC3GxpFQiW925uomUFt2yVHk1J61rzZXcSmjTPn73XyhA3TncoXUNZAJxpxkguD1jtF?cluster=devnet
*/
