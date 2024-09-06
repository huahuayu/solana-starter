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
	balance, err := c.GetBalance(
		context.TODO(),
		"HcNCxoni2Ln5si48s1w8r5TRVH296RQ1MzKeM9FctdPg",
	)
	if err != nil {
		log.Fatalf("get balance, err: %v", err)
	}
	fmt.Println(balance)
}
