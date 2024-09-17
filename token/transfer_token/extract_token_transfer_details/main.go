package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/shopspring/decimal"
)

type Transfer struct {
	Type                      string `json:"type"`
	TokenAddress              string `json:"tokenAddress"`
	Decimals                  uint8  `json:"decimals"`
	Symbol                    string `json:"symbol"`
	Name                      string `json:"name"`
	Authority                 string `json:"authority"`
	Source                    string `json:"source"`
	Destination               string `json:"destination"`
	Amount                    string `json:"amount"`
	UiAmount                  string `json:"uiAmount"`
	IsInnerInstruction        bool   `json:"isInnerInstruction"`
	OuterInstructionIndex     int    `json:"outerInstructionIndex"`
	OuterInstructionProgramID string `json:"outerInstructionProgramID"`
}

type Token struct {
	Address  string // mint address
	Decimals uint8
	Symbol   string
	Name     string
}

func main() {
	c := client.NewClient(rpc.MainnetRPCEndpoint)
	//c := client.NewClient("https://solana.w3node.com/87989be6c2f6334f58643503881317013360a391a6d0e70b8038ec19d45a1afa/api")
	txHash := "4yoaptWrZcNuyPujYTCT3xtydveKa6MLxJr9v4Ypmr9uMpLRUubj2xupL3F8KRQwKVi2YLvetS34sQWYw9R4YupF"
	transfers, err := decodeTokenTransferInstruction(c, txHash)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, transfer := range transfers {
		fmt.Printf("Transfer: %+v\n", *transfer)
	}
}

/* output:
Transfer: {Type:transfer TokenAddress:So11111111111111111111111111111111111111112 Decimals:9 Symbol:SOL Name:Wrapped SOL Authority:DVnVg4p4uzoQfH48iUfx8EGYE2q34xfDzGwwACYDD9G6 Source:6tFPTzVd4Lg3NVgWgwDb7bVfiUcigLXHoBE3Fernjfqw Destination:DzaqzbktzU4PgpXkxpXLWvGH8BAM6P1Q3JjdjEibsHcB Amount:10000 UiAmount:0.00001 IsInnerInstruction:false OuterInstructionIndex:3 OuterInstructionProgramID:TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA}
Transfer: {Type:transfer TokenAddress:So11111111111111111111111111111111111111112 Decimals:9 Symbol:SOL Name:Wrapped SOL Authority:DVnVg4p4uzoQfH48iUfx8EGYE2q34xfDzGwwACYDD9G6 Source:DzaqzbktzU4PgpXkxpXLWvGH8BAM6P1Q3JjdjEibsHcB Destination:BH99eJBXodXtJRCbE4Z2vpashf19pLW9vGf37PTPWH9D Amount:10000 UiAmount:0.00001 IsInnerInstruction:true OuterInstructionIndex:4 OuterInstructionProgramID:675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8}
Transfer: {Type:transfer TokenAddress:1DZ2M31avcvyXMihcX5Pjtcz4qZeGFuQ2gGSjSwoRms Decimals:6 Symbol:WORMS Name:Worms by Matt Furie Authority:5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1 Source:96RaEiBVEZgpWDCKBhmNMu4E3WiAU1thBGk6NYNqH9eK Destination:FZkQSdvQqWbNh1ASdo9MYUxeUNkoNWdC1Wkv33jGBWLc Amount:10281 UiAmount:0.010281 IsInnerInstruction:true OuterInstructionIndex:4 OuterInstructionProgramID:675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8}
*/

func decodeTokenTransferInstruction(c *client.Client, txHash string) ([]*Transfer, error) {
	// Query transaction details
	tx, err := c.GetTransaction(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("error fetching transaction details: %w", err)
	}

	// Cache the account index mapping
	var indexAccountMap = make(map[int]string)
	for i, account := range tx.AccountKeys {
		indexAccountMap[i] = account.ToBase58()
	}

	// Cache the instruction programIDIndex mapping
	var instructionIndexProgramIDMap = make(map[int]int)
	for i, instruction := range tx.Transaction.Message.Instructions {
		instructionIndexProgramIDMap[i] = instruction.ProgramIDIndex
	}

	// Cache token mint addresses
	var mintAddresses = make(map[string]struct{})
	for _, tokenBalance := range tx.Meta.PreTokenBalances {
		mintAddresses[tokenBalance.Mint] = struct{}{}
	}

	var allTransfers []*Transfer

	// Process outer instructions
	for i, instruction := range tx.Transaction.Message.Instructions {
		programID := indexAccountMap[instruction.ProgramIDIndex]
		if programID == common.TokenProgramID.String() {
			transfer, _ := tryDecodeTransfer(instruction, indexAccountMap)
			if transfer != nil {
				transfer.IsInnerInstruction = false
				transfer.OuterInstructionIndex = i
				transfer.OuterInstructionProgramID = programID
				allTransfers = append(allTransfers, transfer)
			}
		}
	}

	// Process inner instructions
	for _, innerInstructions := range tx.Meta.InnerInstructions {
		outerProgramIDIndex := instructionIndexProgramIDMap[int(innerInstructions.Index)]
		outerProgramID := indexAccountMap[outerProgramIDIndex]
		for _, instruction := range innerInstructions.Instructions {
			programID := indexAccountMap[instruction.ProgramIDIndex]
			if programID == common.TokenProgramID.String() {
				transfer, _ := tryDecodeTransfer(instruction, indexAccountMap)
				if transfer != nil {
					transfer.IsInnerInstruction = true
					transfer.OuterInstructionIndex = int(innerInstructions.Index)
					transfer.OuterInstructionProgramID = outerProgramID
					allTransfers = append(allTransfers, transfer)
				}
			}
		}
	}

	// Get all the authorities from allTransfers
	var authorities = make(map[string]struct{})
	for _, transfer := range allTransfers {
		authorities[transfer.Authority] = struct{}{}
	}

	// Try to derive mint address for transfers without mint info
	for _, transfer := range allTransfers {
		if transfer.TokenAddress == "" {
			derivedMint := tryDeriveMintAddress(transfer.Source, transfer.Destination, authorities, mintAddresses)
			if derivedMint != "" {
				transfer.TokenAddress = derivedMint
			}
		}
	}

	// Get and populate token info for all transfers(which mint is not empty)
	for _, transfer := range allTransfers {
		if transfer.TokenAddress != "" {
			token, err := newToken(c, transfer.TokenAddress)
			if err == nil && token != nil {
				transfer.Symbol = token.Symbol
				transfer.Name = token.Name
				transfer.Decimals = token.Decimals
				if token.Decimals <= 0 {
					transfer.UiAmount = transfer.Amount
				} else {
					amount, _ := decimal.NewFromString(transfer.Amount)
					divisor := decimal.New(1, int32(token.Decimals))
					uiAmount := amount.Div(divisor)
					transfer.UiAmount = uiAmount.String()
				}
			}
		}

	}

	return allTransfers, nil
}

func newToken(c *client.Client, mintAddress string) (*Token, error) {
	account, err := c.GetAccountInfo(context.Background(), mintAddress)
	if err != nil {
		return nil, err
	}

	// The decimals are stored at byte offset 44 in the mint account data
	if len(account.Data) < 45 {
		return nil, fmt.Errorf("invalid mint account data")
	}

	decimals := account.Data[44]

	// Get metadata
	mintPubKey := common.PublicKeyFromString(mintAddress)
	metadataAddress, err := token_metadata.GetTokenMetaPubkey(mintPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find metadata PDA: %v", err)
	}

	metadataAccountInfo, err := c.GetAccountInfo(context.Background(), metadataAddress.ToBase58())
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata account info: %v", err)
	}

	metadata, err := token_metadata.MetadataDeserialize(metadataAccountInfo.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize metadata: %v", err)
	}

	return &Token{
		Address:  mintAddress,
		Decimals: decimals,
		Symbol:   metadata.Data.Symbol,
		Name:     metadata.Data.Name,
	}, nil
}

func tryDecodeTransfer(instruction types.CompiledInstruction, indexAccountMap map[int]string) (*Transfer, error) {
	if len(instruction.Data) == 0 {
		return nil, nil
	}

	instructionType := instruction.Data[0]

	var amount uint64
	var sourceAccount, destinationAccount, authorityAccount, mintAddress string

	switch instructionType {
	case 3: // Transfer
		if len(instruction.Data) < 9 || len(instruction.Accounts) < 3 {
			return nil, fmt.Errorf("invalid transfer instruction data: data length %d, accounts length %d", len(instruction.Data), len(instruction.Accounts))
		}
		amount = binary.LittleEndian.Uint64(instruction.Data[1:9])
		sourceAccount = indexAccountMap[instruction.Accounts[0]]
		destinationAccount = indexAccountMap[instruction.Accounts[1]]
		authorityAccount = indexAccountMap[instruction.Accounts[2]]

		// We don't have the mint address for regular transfers, so we'll leave it empty
		mintAddress = ""

	case 12: // TransferChecked
		if len(instruction.Data) < 10 || len(instruction.Accounts) < 4 {
			return nil, fmt.Errorf("invalid transfer checked instruction data: data length %d, accounts length %d", len(instruction.Data), len(instruction.Accounts))
		}
		amount = binary.LittleEndian.Uint64(instruction.Data[1:9])
		sourceAccount = indexAccountMap[instruction.Accounts[0]]
		mintAddress = indexAccountMap[instruction.Accounts[1]]
		destinationAccount = indexAccountMap[instruction.Accounts[2]]
		authorityAccount = indexAccountMap[instruction.Accounts[3]]

	default:
		return nil, nil // Skip unsupported instructions
	}

	transfer := &Transfer{
		Type:         instructionTypeToString(instructionType),
		Source:       sourceAccount,
		Destination:  destinationAccount,
		Authority:    authorityAccount,
		TokenAddress: mintAddress,
		Amount:       fmt.Sprintf("%d", amount),
	}

	return transfer, nil
}

func instructionTypeToString(instructionType byte) string {
	switch instructionType {
	case 3:
		return "transfer"
	case 12:
		return "transferChecked"
	default:
		return "unknown"
	}
}

// tryDeriveMintAddress assume source or destination address is ATA address, compare it to derivedATA, if it's match, then return mint address
func tryDeriveMintAddress(source, destination string, authorities map[string]struct{}, mintAddresses map[string]struct{}) string {
	for authority := range authorities {
		for mintAddress, _ := range mintAddresses {
			derivedATA, _ := deriveAssociatedTokenAddress(
				common.PublicKeyFromString(authority),
				common.PublicKeyFromString(mintAddress),
			)
			if derivedATA.ToBase58() == source || derivedATA.ToBase58() == destination {
				return mintAddress
			}
		}
	}
	return ""
}

// deriveAssociatedTokenAddress derives the associated token address for a given owner and mint
func deriveAssociatedTokenAddress(owner, mint common.PublicKey) (common.PublicKey, error) {
	seeds := [][]byte{
		owner.Bytes(),
		common.TokenProgramID.Bytes(),
		mint.Bytes(),
	}

	programDerivedAddress, _, err := common.FindProgramAddress(seeds, common.SPLAssociatedTokenAccountProgramID)
	if err != nil {
		return common.PublicKey{}, fmt.Errorf("failed to find program address: %w", err)
	}

	return programDerivedAddress, nil
}
