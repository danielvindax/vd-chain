package main

import (
	"flag"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danielvindax/vd-chain/protocol/app/config"
	"github.com/danielvindax/vd-chain/protocol/x/subaccounts/types"
)

/*
Generate a valid SubaccountId with "vindax" prefix

Usage:
	go run scripts/gen_subaccount/gen_subaccount.go
	go run scripts/gen_subaccount/gen_subaccount.go -number 5
*/
func main() {
	// Set address prefixes to "vindax"
	config.SetAddressPrefixes()

	// Parse flags
	var number uint
	flag.UintVar(&number, "number", 0, "subaccount number (0-128000)")
	flag.Parse()

	// Generate a random private key
	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey()
	addr := pubKey.Address()

	// Convert to bech32 with vindax prefix
	bech32Addr := sdk.AccAddress(addr).String()

	// Create SubaccountId
	subaccountId := types.SubaccountId{
		Owner:  bech32Addr,
		Number: uint32(number),
	}

	// Validate to ensure it's correct
	if err := subaccountId.Validate(); err != nil {
		fmt.Printf("ERROR: Generated invalid SubaccountId: %v\n", err)
		return
	}

	fmt.Println("Generated SubaccountId:")
	fmt.Printf("Owner:  %s\n", bech32Addr)
	fmt.Printf("Number: %d\n", number)
	fmt.Println()
	fmt.Println("You can use this in your test:")
	fmt.Printf("owner:  %q,\n", bech32Addr)
	fmt.Printf("number: %d,\n", number)
	fmt.Println()
	fmt.Println("Or as SubaccountId struct:")
	fmt.Printf("&types.SubaccountId{\n")
	fmt.Printf("    Owner:  %q,\n", bech32Addr)
	fmt.Printf("    Number: %d,\n", number)
	fmt.Printf("}\n")
}

