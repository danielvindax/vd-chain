package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danielvindax/vd-chain/protocol/app/config"
)

/*
Generate a valid bech32 address with "vindax" prefix

Usage:
	go run scripts/gen_address/gen_address.go
*/
func main() {
	// Set address prefixes to "vindax"
	config.SetAddressPrefixes()

	// Generate a random private key
	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey()
	addr := pubKey.Address()

	// Convert to bech32 with vindax prefix
	bech32Addr := sdk.AccAddress(addr).String()

	fmt.Println("Generated bech32 address with 'vindax' prefix:")
	fmt.Println(bech32Addr)
	fmt.Println()
	fmt.Println("You can use this address in your test:")
	fmt.Printf("owner:  %q,\n", bech32Addr)
}

