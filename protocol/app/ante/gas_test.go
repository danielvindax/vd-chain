package ante_test

import (
	"reflect"
	"testing"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	
	testapp "github.com/danielvindax/vd-chain/protocol/testutil/app"
	assets "github.com/danielvindax/vd-chain/protocol/x/assets/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/danielvindax/vd-chain/protocol/app/ante"
	testante "github.com/danielvindax/vd-chain/protocol/testutil/ante"
	"github.com/danielvindax/vd-chain/protocol/testutil/constants"
	pricestypes "github.com/danielvindax/vd-chain/protocol/x/prices/types"

	"github.com/stretchr/testify/require"
)

const freeInfiniteGasMeterType = "*types.freeInfiniteGasMeter"

func TestValidateMsgType_FreeInfiniteGasDecorator(t *testing.T) {
	tests := map[string]struct {
		msgOne sdk.Msg
		msgTwo sdk.Msg

		expectFreeInfiniteGasMeter bool
		expectedErr                error
	}{
		"no freeInfiniteGasMeter: no msg": {
			expectFreeInfiniteGasMeter: false,
		},
		"yes freeInfiniteGasMeter: single msg, MsgPlaceOrder": {
			msgOne: constants.Msg_PlaceOrder,

			expectFreeInfiniteGasMeter: true,
		},
		"yes freeInfiniteGasMeter: single msg, Msg_CancelOrder": {
			msgOne: constants.Msg_CancelOrder,

			expectFreeInfiniteGasMeter: true,
		},
		"yes freeInfiniteGasMeter: single msg, MsgUpdateMarketPrices": {
			msgOne: &pricestypes.MsgUpdateMarketPrices{}, // app-injected.

			expectFreeInfiniteGasMeter: true,
		},
		"no freeInfiniteGasMeter: single msg": {
			msgOne: &testdata.TestMsg{Signers: []string{"meh"}},

			expectFreeInfiniteGasMeter: false,
		},
		"no freeInfiniteGasMeter: multi msg, MsgUpdateMarketPrices": {
			msgOne: &pricestypes.MsgUpdateMarketPrices{}, // app-injected.
			msgTwo: &testdata.TestMsg{Signers: []string{"meh"}},

			expectFreeInfiniteGasMeter: false,
		},
		"no freeInfiniteGasMeter: mult msgs, NO off-chain single msg clob tx": {
			msgOne: &testdata.TestMsg{Signers: []string{"meh"}},
			msgTwo: &testdata.TestMsg{Signers: []string{"meh"}},

			expectFreeInfiniteGasMeter: false,
		},
		"no freeInfiniteGasMeter: mult msgs, MsgCancelOrder with Transfer": {
			msgOne: constants.Msg_CancelOrder,
			msgTwo: constants.Msg_Transfer,

			expectFreeInfiniteGasMeter: true,
		},
		"yes freeInfiniteGasMeter: mult msgs, two MsgCancelOrder": {
			msgOne: constants.Msg_CancelOrder,
			msgTwo: constants.Msg_CancelOrder,

			expectFreeInfiniteGasMeter: true,
		},
		"yes freeInfiniteGasMeter: mult msgs, MsgPlaceOrder with Transfer": {
			msgOne: constants.Msg_PlaceOrder,
			msgTwo: constants.Msg_Transfer,

			expectFreeInfiniteGasMeter: true,
		},
		"yes freeInfiniteGasMeter: mult msgs, two MsgPlaceOrder": {
			msgOne: constants.Msg_PlaceOrder,
			msgTwo: constants.Msg_PlaceOrder,

			expectFreeInfiniteGasMeter: true,
		},
		"no freeInfiniteGasMeter: mult msgs, MsgPlaceOrder and MsgCancelOrder": {
			msgOne: constants.Msg_PlaceOrder,
			msgTwo: constants.Msg_CancelOrder,

			expectFreeInfiniteGasMeter: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			suite := testante.SetupTestSuite(t, true)
			suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

			wrappedHandler := ante.NewFreeInfiniteGasDecorator()
			antehandler := sdk.ChainAnteDecorators(wrappedHandler)

			msgs := make([]sdk.Msg, 0)
			if tc.msgOne != nil {
				msgs = append(msgs, tc.msgOne)
			}
			if tc.msgTwo != nil {
				msgs = append(msgs, tc.msgTwo)
			}

			require.NoError(t, suite.TxBuilder.SetMsgs(msgs...))

			// Empty private key, so tx's signature should be empty.
			privs, accNums, accSeqs := []cryptotypes.PrivKey{}, []uint64{}, []uint64{}

			tx, err := suite.CreateTestTx(
				suite.Ctx,
				privs,
				accNums,
				accSeqs,
				suite.Ctx.ChainID(),
				signing.SignMode_SIGN_MODE_DIRECT,
			)
			require.NoError(t, err)

			resultCtx, err := antehandler(suite.Ctx, tx, false)
			require.ErrorIs(t, tc.expectedErr, err)

			meter := resultCtx.GasMeter()

			if !tc.expectFreeInfiniteGasMeter || tc.expectedErr != nil {
				require.NotEqual(t, freeInfiniteGasMeterType, reflect.TypeOf(meter).String())
				require.Equal(t, suite.Ctx, resultCtx)
			} else {
				require.Equal(t, freeInfiniteGasMeterType, reflect.TypeOf(meter).String())
			}
		})
	}
}

func TestSubmitTxnWithGas(t *testing.T) {
	tests := map[string]struct {
		gasFee       sdk.Coins
		responseCode uint32
		logMessage   string
	}{
		"Success - 5 cents usdc gas fee": {
			gasFee:       constants.TestFeeCoins_5Cents,
			responseCode: errors.SuccessABCICode,
		},
		"Success - 5 cents native token gas fee": {
			gasFee:       constants.TestFeeCoins_5Cents_NativeToken,
			responseCode: errors.SuccessABCICode,
		},
		"Failure: 0 gas fee": {
			gasFee:       sdk.Coins{},
			responseCode: sdkerrors.ErrInsufficientFee.ABCICode(),
			logMessage: "insufficient fees; got:  required: 25000000000000000adv4tnt," +
				"25000ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5: insufficient fee",
		},
		"Failure: unsupported gas fee denom": {
			gasFee: sdk.Coins{
				// 1BTC, which is not supported as a gas fee denom, and should be plenty to cover gas.
				sdk.NewCoin(constants.BtcUsd.Denom, sdkmath.NewInt(100_000_000)),
			},
			responseCode: sdkerrors.ErrInsufficientFee.ABCICode(),
			logMessage: "insufficient fees; got: 100000000btc-denom required: 25000000000000000adv4tnt," +
				"25000ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5: insufficient fee",
		},
	}
		for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := &banktypes.MsgSend{
				FromAddress: constants.BobAccAddress.String(),
				ToAddress:   constants.AliceAccAddress.String(),
				Amount: []sdk.Coin{
					sdk.NewCoin(assets.AssetUsdc.Denom, sdkmath.NewInt(1)),
				},
			}

			// Add Bob and Alice accounts to genesis state
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *authtypes.GenesisState) {
						// Add Bob account
						bobAccount := authtypes.NewBaseAccount(
							constants.BobAccAddress,
							constants.BobPrivateKey.PubKey(),
							0,
							0,
						)
						genesisState.Accounts = append(genesisState.Accounts, codectypes.UnsafePackAny(bobAccount))

						// Add Alice account
						aliceAccount := authtypes.NewBaseAccount(
							constants.AliceAccAddress,
							constants.AlicePrivateKey.PubKey(),
							0,
							0,
						)
						genesisState.Accounts = append(genesisState.Accounts, codectypes.UnsafePackAny(aliceAccount))
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						// Add bank balances for Bob
						// Bob needs: 1 USDC to send + gas fees (50000 USDC for USDC gas fee test, or native tokens for native token gas fee test)
						// Give Bob plenty of both USDC and native tokens
						bobFound := false
						for i, bal := range genesisState.Balances {
							if bal.Address == constants.BobAccAddress.String() {
								// Update existing balance
								genesisState.Balances[i] = banktypes.Balance{
									Address: constants.BobAccAddress.String(),
									Coins: sdk.NewCoins(
										// 1 million USDC (plenty for sending 1 USDC + gas fees)
										sdk.NewCoin(assets.AssetUsdc.Denom, sdkmath.NewInt(1_000_000_000_000)),
										// 1 million native tokens (plenty for gas fees)
										sdk.NewCoin(constants.TestNativeTokenDenom, sdkmath.NewInt(1_000_000_000_000_000_000_000_000)),
									),
								}
								bobFound = true
								break
							}
						}
						if !bobFound {
							// Add new balance for Bob
							genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
								Address: constants.BobAccAddress.String(),
								Coins: sdk.NewCoins(
									// 1 million USDC (plenty for sending 1 USDC + gas fees)
									sdk.NewCoin(assets.AssetUsdc.Denom, sdkmath.NewInt(1_000_000_000_000)),
									// 1 million native tokens (plenty for gas fees)
									sdk.NewCoin(constants.TestNativeTokenDenom, sdkmath.NewInt(1_000_000_000_000_000_000_000_000)),
								),
							})
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			msgSendCheckTx := testapp.MustMakeCheckTxWithPrivKeySupplier(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: constants.BobAccAddress.String(),
					Gas:                  constants.TestGasLimit,
					FeeAmt:               tc.gasFee,
				},
				constants.GetPrivateKeyFromAddress,
				msg,
			)

			checkTx := tApp.CheckTx(msgSendCheckTx)
			// Sanity check that gas was used.
			require.Greater(t, checkTx.GasUsed, int64(0))
			require.Equal(t, tc.responseCode, checkTx.Code)
			if tc.responseCode != errors.SuccessABCICode {
				require.Equal(t, tc.logMessage, checkTx.Log)
			}
		})
	}
}
