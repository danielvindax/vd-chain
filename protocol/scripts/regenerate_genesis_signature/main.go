package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmed25519 "github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/danielvindax/vd-chain/protocol/app"
	"github.com/danielvindax/vd-chain/protocol/testutil/constants"
)

// This script regenerates signatures for genesis validator transactions
// after fixing the validator_address prefix from dydxvaloper to vindaxvaloper

const (
	chainID = "dydx-sample-1"
)

// generateValidatorPubkeyFromMnemonic generates an Ed25519 validator public key from a mnemonic
// and returns it as a base64-encoded string
func generateValidatorPubkeyFromMnemonic(mnemonic string) string {
	// Generate Ed25519 validator private key from mnemonic (same as gen-priv-key command)
	validatorPrivKey := tmed25519.GenPrivKeyFromSecret([]byte(mnemonic))
	// Get public key from private key
	validatorPubKey := validatorPrivKey.PubKey()
	// Convert to bytes and encode as base64
	pubkeyBytes := validatorPubKey.Bytes()
	return base64.StdEncoding.EncodeToString(pubkeyBytes)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <validator_name> [value]")
		fmt.Println("Validators: alice, bob, carl, dave")
		fmt.Println("Example: go run main.go dave")
		fmt.Println("Example: go run main.go dave 0")
		fmt.Println("Example: go run main.go bob 500000000000000000000000")
		os.Exit(1)
	}

	validatorName := os.Args[1]
	
	// Parse value (default to 0)
	valueStr := "0"
	if len(os.Args) >= 3 {
		valueStr = os.Args[2]
	}

	// Validator-specific data (pubkey and memo are validator-specific, not account-specific)
	type validatorInfo struct {
		privKey       cryptotypes.PrivKey
		mnemonic      string
		delegatorAddr string
		validatorAddr string
		memo          string
		accountNumber uint64
	}

	validators := map[string]validatorInfo{
		"alice": {
			privKey:       constants.AlicePrivateKey,
			mnemonic:      constants.AliceMnenomic,
			delegatorAddr: constants.AliceAccAddress.String(),
			validatorAddr: constants.AliceValAddress.String(),
			memo:          "17e5e45691f0d01449c84fd4ae87279578cdd7ec@172.17.0.3:26656",
			accountNumber: 0, // Genesis transactions use account number 0
		},
		"bob": {
			privKey:       constants.BobPrivateKey,
			mnemonic:      constants.BobMnenomic,
			delegatorAddr: constants.BobAccAddress.String(),
			validatorAddr: constants.BobValAddress.String(),
			memo:          "b69182310be02559483e42c77b7b104352713166@172.17.0.3:26656",
			accountNumber: 0, // Genesis transactions use account number 0
		},
		"carl": {
			privKey:       constants.CarlPrivateKey,
			mnemonic:      constants.CarlMnenomic,
			delegatorAddr: constants.CarlAccAddress.String(),
			validatorAddr: constants.CarlValAddress.String(),
			memo:          "47539956aaa8e624e0f1d926040e54908ad0eb44@172.17.0.3:26656",
			accountNumber: 0, // Genesis transactions use account number 0
		},
		"dave": {
			privKey:       constants.DavePrivateKey,
			mnemonic:      constants.DaveMnenomic,
			delegatorAddr: constants.DaveAccAddress.String(),
			validatorAddr: constants.DaveValAddress.String(),
			memo:          "5882428984d83b03d0c907c1f0af343534987052@172.17.0.3:26656",
			accountNumber: 0, // Genesis transactions use account number 0
		},
	}

	valInfo, exists := validators[validatorName]
	if !exists {
		fmt.Printf("Unknown validator: %s. Supported: alice, bob, carl, dave\n", validatorName)
		os.Exit(1)
	}

	privKey := valInfo.privKey
	delegatorAddr := valInfo.delegatorAddr
	validatorAddr := valInfo.validatorAddr
	memo := valInfo.memo
	accountNumber := valInfo.accountNumber

	// Generate validator pubkey from mnemonic
	pubkeyBase64 := generateValidatorPubkeyFromMnemonic(valInfo.mnemonic)

	// Decode validator pubkey
	pubkeyBytes, err := base64.StdEncoding.DecodeString(pubkeyBase64)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode pubkey: %v", err))
	}

	// Create Ed25519 pubkey
	ed25519PubKey := &ed25519.PubKey{Key: pubkeyBytes}

	// Create pubkey Any
	pubkeyAny, err := codectypes.NewAnyWithValue(ed25519PubKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to create pubkey Any: %v", err))
	}

	// Create MsgCreateValidator
	msg := &stakingtypes.MsgCreateValidator{
		Description: stakingtypes.Description{
			Moniker:         validatorName,
			Identity:        "",
			Website:         "",
			SecurityContact: "",
			Details:         "",
		},
		Commission: stakingtypes.CommissionRates{
			Rate:          math.LegacyMustNewDecFromStr("1.0"),
			MaxRate:       math.LegacyMustNewDecFromStr("1.0"),
			MaxChangeRate: math.LegacyMustNewDecFromStr("0.01"),
		},
		MinSelfDelegation: math.NewInt(1),
		DelegatorAddress:  delegatorAddr,
		ValidatorAddress:  validatorAddr,
		Pubkey:            pubkeyAny,
		Value: func() sdk.Coin {
			amount, ok := math.NewIntFromString(valueStr)
			if !ok {
				panic(fmt.Sprintf("Failed to parse amount: %s", valueStr))
			}
			return sdk.Coin{
				Denom:  "adv4tnt",
				Amount: amount,
			}
		}(),
	}

	// Setup encoding config
	encodingConfig := app.GetEncodingConfig()
	clientCtx := client.Context{}.WithTxConfig(encodingConfig.TxConfig)
	txBuilder := clientCtx.TxConfig.NewTxBuilder()

	// Set messages
	err = txBuilder.SetMsgs(msg)
	if err != nil {
		panic(fmt.Sprintf("Failed to set messages: %v", err))
	}

	// Set memo
	txBuilder.SetMemo(memo)

	// Set timeout height
	txBuilder.SetTimeoutHeight(0)

	// Set fee
	txBuilder.SetFeeAmount(sdk.Coins{})
	txBuilder.SetGasLimit(200000)

	// First round: set empty signature to get signer info
	sigV2 := signing.SignatureV2{
		PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: 0,
	}

	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		panic(fmt.Sprintf("Failed to set signatures: %v", err))
	}

	// Second round: sign the transaction
	signerData := xauthsigning.SignerData{
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      0,
	}

	signMode := signing.SignMode_SIGN_MODE_DIRECT
	signBytes, err := xauthsigning.GetSignBytesAdapter(
		context.Background(),
		clientCtx.TxConfig.SignModeHandler(),
		signMode,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to get sign bytes: %v", err))
	}

	sig, err := privKey.Sign(signBytes)
	if err != nil {
		panic(fmt.Sprintf("Failed to sign: %v", err))
	}

	sigV2.Data.(*signing.SingleSignatureData).Signature = sig
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		panic(fmt.Sprintf("Failed to set final signatures: %v", err))
	}

	// Get the signature bytes
	sigs, err := txBuilder.GetTx().GetSignaturesV2()
	if err != nil {
		panic(fmt.Sprintf("Failed to get signatures: %v", err))
	}

	if len(sigs) == 0 {
		panic("No signatures found")
	}

	sigData := sigs[0].Data.(*signing.SingleSignatureData)
	signatureBase64 := base64.StdEncoding.EncodeToString(sigData.Signature)

	fmt.Printf("Validator: %s\n", validatorName)
	fmt.Printf("Delegator Address: %s\n", delegatorAddr)
	fmt.Printf("Validator Address: %s\n", validatorAddr)
	fmt.Printf("Value: %s\n", valueStr)
	fmt.Printf("\nValidator Pubkey (base64):\n")
	fmt.Println(pubkeyBase64)
	fmt.Printf("\nNew signature:\n")
	fmt.Println(signatureBase64)
	fmt.Println("\nCopy this signature and pubkey and replace them in genesis.go")
}
