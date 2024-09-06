package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
)

// Fee payer account
var feePayer, _ = types.AccountFromBase58("D8i1DFhgxWkBC52kRUtDRkZL5J5bUDJXssJtsehWYT51txphk8ipWe8goFKJt6638vAmEHVxdovsjmfiHPvKPbS")

// Alice account
var alice, _ = types.AccountFromBase58("5ob5v9uGyJstuENC4pu7ScCdkCPiAZXqFLhu3qUoHB8uePLchCpPmgWyXPZ24NxLBD8dUP6UFNNYXKsFJtFKie74")

// Token Mint Pubkey
var mintPubkey = common.PublicKeyFromString("gYqzga5v1RoVWxtfXizHuoyxUpTnzf9WyrXftTkDfpT")

// Function to set token metadata
func setTokenMetadata(c *client.Client, mintPubkey common.PublicKey, data token_metadata.DataV2) error {
	// Get the metadata account for the mint
	metadataAccount, err := token_metadata.GetTokenMetaPubkey(mintPubkey)
	if err != nil {
		return fmt.Errorf("failed to get metadata account: %v", err)
	}

	// Fetch the latest blockhash for the transaction
	res, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	// Create Metadata Account using CreateMetadataAccountV2
	createMetadataAccount := token_metadata.CreateMetadataAccountV3(token_metadata.CreateMetadataAccountV3Param{
		Metadata:      metadataAccount,
		Mint:          mintPubkey,
		MintAuthority: alice.PublicKey,
		Payer:         feePayer.PublicKey,
		Data:          data,
		IsMutable:     true,
	})

	// Create a transaction with the given instructions
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions:    []types.Instruction{createMetadataAccount},
		}),
		Signers: []types.Account{feePayer, alice}, // Both feePayer and mintAuthority (alice) need to sign
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	// Send the transaction
	sig, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	// Output transaction explorer link
	fmt.Printf("Check the transaction at: https://explorer.solana.com/tx/%s?cluster=devnet\n", sig)

	return nil
}

func main() {
	// Create a new Solana client pointing to the Devnet cluster
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	// Token metadata to be set
	metadataData := token_metadata.DataV2{
		Name:                 "Cool token",            // Token name
		Symbol:               "COOL",                  // Token symbol
		Uri:                  "https://cooltoken.com", // Token URI
		SellerFeeBasisPoints: 0,
		Creators:             nil,
	}

	// Attempt to set the token metadata
	err := setTokenMetadata(c, mintPubkey, metadataData)
	if err != nil {
		log.Fatalf("set token metadata error: %v", err)
	}
}

/*
Check the transaction at: https://explorer.solana.com/tx/En2UPcgjBkn1eP3mV4hGBcVT6GjnMfYeHrEPZcmhwFAkBovBXBPFktTu6mRpAt7uUrCHftSEiqr65hgjdnUJMju?cluster=devnet
*/
