package main

import (
	"context"
	"fmt"
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

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				token.TransferChecked(token.TransferCheckedParam{
					From:     aliceTokenATAPubkey,
					To:       newAccount.PublicKey,
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
2024/09/06 10:55:37 new account: 49mFqeNosQqDk3aj332ayCtmBdVU85utcWAjggZaW8Lw 35Y1JSqrMvcN2kAczXGWxJiwJ2V2MyYVBDsfaVV2BbFhLCmc1SWToaGmwRtbsGxDpf9bK2Uf1fj8HLFFYLzBYDi7
2024/09/06 10:55:37 send raw tx error, err: {"code":-32002,"message":"Transaction simulation failed: Error processing Instruction 0: invalid account data for instruction","data":{"accounts":null,"err":{"InstructionError":[0,"InvalidAccountData"]},"innerInstructions":null,"logs":["Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]","Program log: Instruction: TransferChecked","Program log: Error: InvalidAccountData","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 2985 of 200000 compute units","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA failed: invalid account data for instruction"],"replacementBlockhash":null,"returnData":null,"unitsConsumed":2985}}
*/
