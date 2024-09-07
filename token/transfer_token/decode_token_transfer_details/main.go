package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/blocto/solana-go-sdk/rpc"
)

type TokenInfo struct {
	Mint     string
	Decimals uint8
	Symbol   string
	Name     string
}

type TokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       uint8   `json:"decimals"`
	UiAmount       float64 `json:"uiAmount"`
	UiAmountString string  `json:"uiAmountString"`
}

type MintInfo struct {
	Mint   string `json:"mint"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type TransferInfo struct {
	Authority   string      `json:"authority,omitempty"`
	Destination string      `json:"destination"`
	Mint        MintInfo    `json:"mint"`
	Source      string      `json:"source"`
	TokenAmount TokenAmount `json:"tokenAmount"`
}

type TransferData struct {
	Info TransferInfo `json:"info"`
	Type string       `json:"type"`
}

func getTokenInfo(c *client.Client, mintAddress string) (TokenInfo, error) {
	account, err := c.GetAccountInfo(context.Background(), mintAddress)
	if err != nil {
		return TokenInfo{}, err
	}

	// The decimals are stored at byte offset 44 in the mint account data
	if len(account.Data) < 45 {
		return TokenInfo{}, fmt.Errorf("invalid mint account data")
	}

	decimals := account.Data[44]

	// Get metadata
	mintPubKey := common.PublicKeyFromString(mintAddress)
	metadataAddress, err := token_metadata.GetTokenMetaPubkey(mintPubKey)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("failed to find metadata PDA: %v", err)
	}

	metadataAccountInfo, err := c.GetAccountInfo(context.Background(), metadataAddress.ToBase58())
	if err != nil {
		return TokenInfo{}, fmt.Errorf("failed to get metadata account info: %v", err)
	}

	metadata, err := token_metadata.MetadataDeserialize(metadataAccountInfo.Data)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("failed to deserialize metadata: %v", err)
	}

	return TokenInfo{
		Mint:     mintAddress,
		Decimals: decimals,
		Symbol:   metadata.Data.Symbol,
		Name:     metadata.Data.Name,
	}, nil
}

func trimTrailingZeros(s string) string {
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	signature := "4Wo87ndvevGXEprhu4UwBHC3GxpFQiW925uomUFt2yVHk1J61rzZXcSmjTPn73XyhA3TncoXUNZAJxpxkguD1jtF"

	// Query transaction details
	txDetails, err := c.GetTransaction(context.Background(), signature)
	if err != nil {
		log.Fatalf("failed to get transaction details, err: %v\n", err)
	}

	// Cache the account index mapping
	var indexAccountMap = make(map[int]string)
	for i, account := range txDetails.Transaction.Message.Accounts {
		indexAccountMap[i] = account.ToBase58()
	}

	// Extracting transfer details
	for _, instruction := range txDetails.Transaction.Message.Instructions {
		if indexAccountMap[instruction.ProgramIDIndex] == common.TokenProgramID.ToBase58() { // Check if it's a Token Program instruction
			var transferData TransferData
			var amount uint64
			var mint string

			switch instruction.Data[0] {
			case 3: // Transfer
				if len(instruction.Data) < 9 {
					continue
				}
				amount = binary.LittleEndian.Uint64(instruction.Data[1:9])
				transferData.Type = "transfer"
				transferData.Info.Source = indexAccountMap[instruction.Accounts[0]]
				transferData.Info.Destination = indexAccountMap[instruction.Accounts[1]]
				transferData.Info.Authority = indexAccountMap[instruction.Accounts[2]]
				mint = indexAccountMap[instruction.Accounts[0]] // Assume ATA, get mint from source

			case 12: // TransferChecked
				if len(instruction.Data) < 9 {
					continue
				}
				amount = binary.LittleEndian.Uint64(instruction.Data[1:9])
				transferData.Type = "transferChecked"
				transferData.Info.Source = indexAccountMap[instruction.Accounts[0]]
				transferData.Info.Mint.Mint = indexAccountMap[instruction.Accounts[1]]
				transferData.Info.Destination = indexAccountMap[instruction.Accounts[2]]
				transferData.Info.Authority = indexAccountMap[instruction.Accounts[3]]
				mint = indexAccountMap[instruction.Accounts[1]]

			default:
				continue // Not a transfer instruction
			}

			// Fetch token info
			tokenInfo, err := getTokenInfo(c, mint)
			if err != nil {
				fmt.Println("Error fetching token info: ", err)
				continue
			}

			// Calculate UI amount using big.Rat
			amountBig := new(big.Rat).SetUint64(amount)
			divisor := new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(tokenInfo.Decimals)), nil))
			uiAmount := new(big.Rat).Quo(amountBig, divisor)

			// Convert uiAmount to a float64 for JSON marshaling
			uiAmountFloat, _ := uiAmount.Float64()

			transferData.Info.Mint = MintInfo{
				Mint:   tokenInfo.Mint,
				Symbol: tokenInfo.Symbol,
				Name:   tokenInfo.Name,
			}

			transferData.Info.TokenAmount = TokenAmount{
				Amount:         fmt.Sprintf("%d", amount),
				Decimals:       tokenInfo.Decimals,
				UiAmount:       uiAmountFloat,
				UiAmountString: trimTrailingZeros(uiAmount.FloatString(int(tokenInfo.Decimals))),
			}

			// Marshal the struct to JSON
			jsonData, err := json.MarshalIndent(transferData, "", "  ")
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				continue
			}

			// Print the JSON
			fmt.Println(string(jsonData))
		}
	}
}

/*
{
  "info": {
    "authority": "HcNCxoni2Ln5si48s1w8r5TRVH296RQ1MzKeM9FctdPg",
    "destination": "H5uz148KxqRTg6NLeY4qs8faPHVY39rmGjkcZN4dvYt9",
    "mint": {
      "mint": "gYqzga5v1RoVWxtfXizHuoyxUpTnzf9WyrXftTkDfpT",
      "symbol": "COOL",
      "name": "Cool token"
    },
    "source": "BdEcBm46DWCEBFXVHwXhW76RLqzyCpaiJMxgveL8dLEm",
    "tokenAmount": {
      "amount": "10000000",
      "decimals": 8,
      "uiAmount": 0.1,
      "uiAmountString": "0.1"
    }
  },
  "type": "transferChecked"
}
*/
