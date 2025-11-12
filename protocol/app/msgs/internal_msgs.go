package msgs

import (
	upgrade "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensus "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	ibcconn "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	"github.com/danielvindax/vd-chain/protocol/lib"
	accountplus "github.com/danielvindax/vd-chain/protocol/x/accountplus/types"
	affiliates "github.com/danielvindax/vd-chain/protocol/x/affiliates/types"
	blocktime "github.com/danielvindax/vd-chain/protocol/x/blocktime/types"
	bridge "github.com/danielvindax/vd-chain/protocol/x/bridge/types"
	clob "github.com/danielvindax/vd-chain/protocol/x/clob/types"
	delaymsg "github.com/danielvindax/vd-chain/protocol/x/delaymsg/types"
	feetiers "github.com/danielvindax/vd-chain/protocol/x/feetiers/types"
	govplus "github.com/danielvindax/vd-chain/protocol/x/govplus/types"
	listing "github.com/danielvindax/vd-chain/protocol/x/listing/types"
	perpetuals "github.com/danielvindax/vd-chain/protocol/x/perpetuals/types"
	prices "github.com/danielvindax/vd-chain/protocol/x/prices/types"
	ratelimit "github.com/danielvindax/vd-chain/protocol/x/ratelimit/types"
	revshare "github.com/danielvindax/vd-chain/protocol/x/revshare/types"
	rewards "github.com/danielvindax/vd-chain/protocol/x/rewards/types"
	sending "github.com/danielvindax/vd-chain/protocol/x/sending/types"
	stats "github.com/danielvindax/vd-chain/protocol/x/stats/types"
	vault "github.com/danielvindax/vd-chain/protocol/x/vault/types"
	vest "github.com/danielvindax/vd-chain/protocol/x/vest/types"
)

var (
	// InternalMsgSamplesAll are msgs that are used only used internally.
	InternalMsgSamplesAll = lib.MergeAllMapsMustHaveDistinctKeys(InternalMsgSamplesGovAuth)

	// InternalMsgSamplesGovAuth are msgs that are used only used internally.
	// GovAuth means that these messages must originate from the gov module and
	// signed by gov module account.
	// InternalMsgSamplesAll are msgs that are used only used internally.
	InternalMsgSamplesGovAuth = lib.MergeAllMapsMustHaveDistinctKeys(
		InternalMsgSamplesDefault,
		InternalMsgSamplesDydxCustom,
	)

	// CosmosSDK default modules
	InternalMsgSamplesDefault = map[string]sdk.Msg{
		// auth
		"/cosmos.auth.v1beta1.MsgUpdateParams": &auth.MsgUpdateParams{},

		// bank
		"/cosmos.bank.v1beta1.MsgSetSendEnabled":         &bank.MsgSetSendEnabled{},
		"/cosmos.bank.v1beta1.MsgSetSendEnabledResponse": nil,
		"/cosmos.bank.v1beta1.MsgUpdateParams":           &bank.MsgUpdateParams{},
		"/cosmos.bank.v1beta1.MsgUpdateParamsResponse":   nil,

		// consensus
		"/cosmos.consensus.v1.MsgUpdateParams":         &consensus.MsgUpdateParams{},
		"/cosmos.consensus.v1.MsgUpdateParamsResponse": nil,

		// crisis
		"/cosmos.crisis.v1beta1.MsgUpdateParams":         &crisis.MsgUpdateParams{},
		"/cosmos.crisis.v1beta1.MsgUpdateParamsResponse": nil,

		// distribution
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpend":         &distribution.MsgCommunityPoolSpend{},
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpendResponse": nil,
		"/cosmos.distribution.v1beta1.MsgUpdateParams":               &distribution.MsgUpdateParams{},
		"/cosmos.distribution.v1beta1.MsgUpdateParamsResponse":       nil,

		// gov
		"/cosmos.gov.v1.MsgExecLegacyContent":         &gov.MsgExecLegacyContent{},
		"/cosmos.gov.v1.MsgExecLegacyContentResponse": nil,
		"/cosmos.gov.v1.MsgUpdateParams":              &gov.MsgUpdateParams{},
		"/cosmos.gov.v1.MsgUpdateParamsResponse":      nil,

		// slashing
		"/cosmos.slashing.v1beta1.MsgUpdateParams":         &slashing.MsgUpdateParams{},
		"/cosmos.slashing.v1beta1.MsgUpdateParamsResponse": nil,

		// staking
		"/cosmos.staking.v1beta1.MsgSetProposers":         &staking.MsgSetProposers{},
		"/cosmos.staking.v1beta1.MsgSetProposersResponse": nil,
		"/cosmos.staking.v1beta1.MsgUpdateParams":         &staking.MsgUpdateParams{},
		"/cosmos.staking.v1beta1.MsgUpdateParamsResponse": nil,

		// upgrade
		"/cosmos.upgrade.v1beta1.MsgCancelUpgrade":           &upgrade.MsgCancelUpgrade{},
		"/cosmos.upgrade.v1beta1.MsgCancelUpgradeResponse":   nil,
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade":         &upgrade.MsgSoftwareUpgrade{},
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgradeResponse": nil,

		// ibc
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParams":         &icahosttypes.MsgUpdateParams{},
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParamsResponse": nil,
		"/ibc.applications.transfer.v1.MsgUpdateParams":                         &ibctransfer.MsgUpdateParams{},
		"/ibc.applications.transfer.v1.MsgUpdateParamsResponse":                 nil,
		"/ibc.core.client.v1.MsgUpdateParams":                                   &ibcclient.MsgUpdateParams{},
		"/ibc.core.client.v1.MsgUpdateParamsResponse":                           nil,
		"/ibc.core.connection.v1.MsgUpdateParams":                               &ibcconn.MsgUpdateParams{},
		"/ibc.core.connection.v1.MsgUpdateParamsResponse":                       nil,
	}

	// Custom modules
	InternalMsgSamplesDydxCustom = map[string]sdk.Msg{
		// affiliates
		"/vindax.affiliates.MsgUpdateAffiliateTiers":              &affiliates.MsgUpdateAffiliateTiers{},
		"/vindax.affiliates.MsgUpdateAffiliateTiersResponse":      nil,
		"/vindax.affiliates.MsgUpdateAffiliateWhitelist":          &affiliates.MsgUpdateAffiliateWhitelist{},
		"/vindax.affiliates.MsgUpdateAffiliateWhitelistResponse":  nil,
		"/vindax.affiliates.MsgUpdateAffiliateParameters":         &affiliates.MsgUpdateAffiliateParameters{},
		"/vindax.affiliates.MsgUpdateAffiliateParametersResponse": nil,
		"/vindax.affiliates.MsgUpdateAffiliateOverrides":          &affiliates.MsgUpdateAffiliateOverrides{},
		"/vindax.affiliates.MsgUpdateAffiliateOverridesResponse":  nil,

		// accountplus
		"/vindax.accountplus.MsgSetActiveState":         &accountplus.MsgSetActiveState{},
		"/vindax.accountplus.MsgSetActiveStateResponse": nil,

		// blocktime
		"/vindax.blocktime.MsgUpdateDowntimeParams":          &blocktime.MsgUpdateDowntimeParams{},
		"/vindax.blocktime.MsgUpdateDowntimeParamsResponse":  nil,
		"/vindax.blocktime.MsgUpdateSynchronyParams":         &blocktime.MsgUpdateSynchronyParams{},
		"/vindax.blocktime.MsgUpdateSynchronyParamsResponse": nil,

		// bridge
		"/vindax.bridge.MsgCompleteBridge":              &bridge.MsgCompleteBridge{},
		"/vindax.bridge.MsgCompleteBridgeResponse":      nil,
		"/vindax.bridge.MsgUpdateEventParams":           &bridge.MsgUpdateEventParams{},
		"/vindax.bridge.MsgUpdateEventParamsResponse":   nil,
		"/vindax.bridge.MsgUpdateProposeParams":         &bridge.MsgUpdateProposeParams{},
		"/vindax.bridge.MsgUpdateProposeParamsResponse": nil,
		"/vindax.bridge.MsgUpdateSafetyParams":          &bridge.MsgUpdateSafetyParams{},
		"/vindax.bridge.MsgUpdateSafetyParamsResponse":  nil,

		// clob
		"/vindax.clob.MsgCreateClobPair":                             &clob.MsgCreateClobPair{},
		"/vindax.clob.MsgCreateClobPairResponse":                     nil,
		"/vindax.clob.MsgUpdateBlockRateLimitConfiguration":          &clob.MsgUpdateBlockRateLimitConfiguration{},
		"/vindax.clob.MsgUpdateBlockRateLimitConfigurationResponse":  nil,
		"/vindax.clob.MsgUpdateClobPair":                             &clob.MsgUpdateClobPair{},
		"/vindax.clob.MsgUpdateClobPairResponse":                     nil,
		"/vindax.clob.MsgUpdateEquityTierLimitConfiguration":         &clob.MsgUpdateEquityTierLimitConfiguration{},
		"/vindax.clob.MsgUpdateEquityTierLimitConfigurationResponse": nil,
		"/vindax.clob.MsgUpdateLiquidationsConfig":                   &clob.MsgUpdateLiquidationsConfig{},
		"/vindax.clob.MsgUpdateLiquidationsConfigResponse":           nil,

		// delaymsg
		"/vindax.delaymsg.MsgDelayMessage":         &delaymsg.MsgDelayMessage{},
		"/vindax.delaymsg.MsgDelayMessageResponse": nil,

		// feetiers
		"/vindax.feetiers.MsgUpdatePerpetualFeeParams":           &feetiers.MsgUpdatePerpetualFeeParams{},
		"/vindax.feetiers.MsgUpdatePerpetualFeeParamsResponse":   nil,
		"/vindax.feetiers.MsgSetMarketFeeDiscountParams":         &feetiers.MsgSetMarketFeeDiscountParams{},
		"/vindax.feetiers.MsgSetMarketFeeDiscountParamsResponse": nil,
		"/vindax.feetiers.MsgSetStakingTiers":                    &feetiers.MsgSetStakingTiers{},
		"/vindax.feetiers.MsgSetStakingTiersResponse":            nil,

		// govplus
		"/vindax.govplus.MsgSlashValidator":         &govplus.MsgSlashValidator{},
		"/vindax.govplus.MsgSlashValidatorResponse": nil,

		// listing
		"/vindax.listing.MsgSetMarketsHardCap":                       &listing.MsgSetMarketsHardCap{},
		"/vindax.listing.MsgSetMarketsHardCapResponse":               nil,
		"/vindax.listing.MsgSetListingVaultDepositParams":            &listing.MsgSetListingVaultDepositParams{},
		"/vindax.listing.MsgSetListingVaultDepositParamsResponse":    nil,
		"/vindax.listing.MsgUpgradeIsolatedPerpetualToCross":         &listing.MsgUpgradeIsolatedPerpetualToCross{},
		"/vindax.listing.MsgUpgradeIsolatedPerpetualToCrossResponse": nil,

		// perpetuals
		"/vindax.perpetuals.MsgCreatePerpetual":               &perpetuals.MsgCreatePerpetual{},
		"/vindax.perpetuals.MsgCreatePerpetualResponse":       nil,
		"/vindax.perpetuals.MsgSetLiquidityTier":              &perpetuals.MsgSetLiquidityTier{},
		"/vindax.perpetuals.MsgSetLiquidityTierResponse":      nil,
		"/vindax.perpetuals.MsgUpdateParams":                  &perpetuals.MsgUpdateParams{},
		"/vindax.perpetuals.MsgUpdateParamsResponse":          nil,
		"/vindax.perpetuals.MsgUpdatePerpetualParams":         &perpetuals.MsgUpdatePerpetualParams{},
		"/vindax.perpetuals.MsgUpdatePerpetualParamsResponse": nil,

		// prices
		"/vindax.prices.MsgCreateOracleMarket":         &prices.MsgCreateOracleMarket{},
		"/vindax.prices.MsgCreateOracleMarketResponse": nil,
		"/vindax.prices.MsgUpdateMarketParam":          &prices.MsgUpdateMarketParam{},
		"/vindax.prices.MsgUpdateMarketParamResponse":  nil,

		// ratelimit
		"/vindax.ratelimit.MsgSetLimitParams":         &ratelimit.MsgSetLimitParams{},
		"/vindax.ratelimit.MsgSetLimitParamsResponse": nil,

		// revshare
		"/vindax.revshare.MsgSetMarketMapperRevShareDetailsForMarket":         &revshare.MsgSetMarketMapperRevShareDetailsForMarket{}, //nolint:lll
		"/vindax.revshare.MsgSetMarketMapperRevShareDetailsForMarketResponse": nil,
		"/vindax.revshare.MsgSetMarketMapperRevenueShare":                     &revshare.MsgSetMarketMapperRevenueShare{}, //nolint:lll
		"/vindax.revshare.MsgSetMarketMapperRevenueShareResponse":             nil,
		"/vindax.revshare.MsgSetOrderRouterRevShare":                          &revshare.MsgSetOrderRouterRevShare{}, //nolint:lll
		"/vindax.revshare.MsgSetOrderRouterRevShareResponse":                  nil,
		"/vindax.revshare.MsgUpdateUnconditionalRevShareConfig":               &revshare.MsgUpdateUnconditionalRevShareConfig{}, //nolint:lll
		"/vindax.revshare.MsgUpdateUnconditionalRevShareConfigResponse":       nil,

		// rewards
		"/vindax.rewards.MsgUpdateParams":         &rewards.MsgUpdateParams{},
		"/vindax.rewards.MsgUpdateParamsResponse": nil,

		// sending
		"/vindax.sending.MsgSendFromModuleToAccount":         &sending.MsgSendFromModuleToAccount{},
		"/vindax.sending.MsgSendFromModuleToAccountResponse": nil,

		// stats
		"/vindax.stats.MsgUpdateParams":         &stats.MsgUpdateParams{},
		"/vindax.stats.MsgUpdateParamsResponse": nil,

		// vault
		"/vindax.vault.MsgUnlockShares":                 &vault.MsgUnlockShares{},
		"/vindax.vault.MsgUnlockSharesResponse":         nil,
		"/vindax.vault.MsgUpdateOperatorParams":         &vault.MsgUpdateOperatorParams{},
		"/vindax.vault.MsgUpdateOperatorParamsResponse": nil,

		// vest
		"/vindax.vest.MsgSetVestEntry":            &vest.MsgSetVestEntry{},
		"/vindax.vest.MsgSetVestEntryResponse":    nil,
		"/vindax.vest.MsgDeleteVestEntry":         &vest.MsgDeleteVestEntry{},
		"/vindax.vest.MsgDeleteVestEntryResponse": nil,
	}
)
