package main

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
	"log"
)

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	// should pass a token account address
	// in Solana, each token account is associated with a specific mint. This means that when you create a token account, you specify the mint that the token account is associated with. Once this association is made, it cannot be changed.  Therefore, when you query the balance of a token account, you don't need to specify the mint address because the token account already has that information. The Solana protocol knows which mint the token account is associated with, and it uses this information to correctly interpret the balance of the token account.  In other words, the balance of a token account is inherently tied to the mint that it's associated with, so there's no need to specify the mint when querying the balance. The mint information is already encapsulated within the token account itself.
	tokenAmount, err := c.GetTokenAccountBalance(
		context.Background(),
		"HeCBh32JJ8DxcjTyc6q46tirHR8hd2xj3mGoAcQ7eduL",
	)
	if err != nil {
		log.Fatalln("get balance error", err)
	}
	// the smallest unit like lamports
	fmt.Println("balance", tokenAmount.Amount)
	// the decimals of mint which token account holds
	fmt.Println("decimals", tokenAmount.Decimals)
}
