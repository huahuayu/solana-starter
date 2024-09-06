package main

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/mr-tron/base58"
	"log"
)

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)
	account := types.NewAccount()
	fmt.Printf("created account: %v, private key: %v\n", account.PublicKey.ToBase58(), base58.Encode(account.PrivateKey))
	sig, err := c.RequestAirdrop(context.TODO(), account.PublicKey.ToBase58(), 1e9)
	if err != nil {
		log.Fatalf("failed to request airdrop, err: %v", err)
	}
	fmt.Printf("requested airdrop, signature: %v\n", sig)
	fmt.Printf("check tx at: https://explorer.solana.com/tx/%s?cluster=devnet\n", sig)
}

/* alice
created account: HcNCxoni2Ln5si48s1w8r5TRVH296RQ1MzKeM9FctdPg, private key: 5ob5v9uGyJstuENC4pu7ScCdkCPiAZXqFLhu3qUoHB8uePLchCpPmgWyXPZ24NxLBD8dUP6UFNNYXKsFJtFKie74
requested airdrop, signature: 53scqsHpne23owvYdAaZq5DWtzpFtvYW5w64VybrZQ3cuqRojZNzzhQ7EEUL6K26m1FwzNAgLd5yWzP4cPqnuEw2
check tx at: https://explorer.solana.com/tx/53scqsHpne23owvYdAaZq5DWtzpFtvYW5w64VybrZQ3cuqRojZNzzhQ7EEUL6K26m1FwzNAgLd5yWzP4cPqnuEw2?cluster=devnet
*/

/* bob
created account: GLndC8XmRT5o6oBLwn8scDNvFY5MuX78wxJQsW5tXctk, private key: D8i1DFhgxWkBC52kRUtDRkZL5J5bUDJXssJtsehWYT51txphk8ipWe8goFKJt6638vAmEHVxdovsjmfiHPvKPbS
requested airdrop, signature: 2r8VvroDMH9qpGbswZuKt7DaKx12fyzTkc7xJ9Th2KBrT14sYViEd3sHnw67EE6HppXKHZLWKXTsqMRP6kYZ6Lc4
check tx at: https://explorer.solana.com/tx/2r8VvroDMH9qpGbswZuKt7DaKx12fyzTkc7xJ9Th2KBrT14sYViEd3sHnw67EE6HppXKHZLWKXTsqMRP6kYZ6Lc4?cluster=devnet
*/
