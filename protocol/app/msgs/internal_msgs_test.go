package msgs_test

import (
	"sort"
	"testing"

	"github.com/danielvindax/vd-chain/protocol/app/msgs"
	"github.com/danielvindax/vd-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

func TestInternalMsgSamples_All_Key(t *testing.T) {
	expectedAllInternalMsgs := lib.MergeAllMapsMustHaveDistinctKeys(msgs.InternalMsgSamplesGovAuth)
	require.Equal(t, expectedAllInternalMsgs, msgs.InternalMsgSamplesAll)
}

func TestInternalMsgSamples_All_Value(t *testing.T) {
	validateMsgValue(t, msgs.InternalMsgSamplesAll)
}

func TestInternalMsgSamples_Gov_Key(t *testing.T) {
	expectedMsgs := []string{
		// auth
		"/cosmos.auth.v1beta1.MsgUpdateParams",

		// bank
		"/cosmos.bank.v1beta1.MsgSetSendEnabled",
		"/cosmos.bank.v1beta1.MsgSetSendEnabledResponse",
		"/cosmos.bank.v1beta1.MsgUpdateParams",
		"/cosmos.bank.v1beta1.MsgUpdateParamsResponse",

		// consensus
		"/cosmos.consensus.v1.MsgUpdateParams",
		"/cosmos.consensus.v1.MsgUpdateParamsResponse",

		// crisis
		"/cosmos.crisis.v1beta1.MsgUpdateParams",
		"/cosmos.crisis.v1beta1.MsgUpdateParamsResponse",

		// distribution
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpend",
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpendResponse",
		"/cosmos.distribution.v1beta1.MsgUpdateParams",
		"/cosmos.distribution.v1beta1.MsgUpdateParamsResponse",

		// gov
		"/cosmos.gov.v1.MsgExecLegacyContent",
		"/cosmos.gov.v1.MsgExecLegacyContentResponse",
		"/cosmos.gov.v1.MsgUpdateParams",
		"/cosmos.gov.v1.MsgUpdateParamsResponse",

		// slashing
		"/cosmos.slashing.v1beta1.MsgUpdateParams",
		"/cosmos.slashing.v1beta1.MsgUpdateParamsResponse",

		// staking
		"/cosmos.staking.v1beta1.MsgSetProposers",
		"/cosmos.staking.v1beta1.MsgSetProposersResponse",
		"/cosmos.staking.v1beta1.MsgUpdateParams",
		"/cosmos.staking.v1beta1.MsgUpdateParamsResponse",

		// upgrade
		"/cosmos.upgrade.v1beta1.MsgCancelUpgrade",
		"/cosmos.upgrade.v1beta1.MsgCancelUpgradeResponse",
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgradeResponse",

		// ibc
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParams",
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParamsResponse",
		"/ibc.applications.transfer.v1.MsgUpdateParams",
		"/ibc.applications.transfer.v1.MsgUpdateParamsResponse",
		"/ibc.core.client.v1.MsgUpdateParams",
		"/ibc.core.client.v1.MsgUpdateParamsResponse",
		"/ibc.core.connection.v1.MsgUpdateParams",
		"/ibc.core.connection.v1.MsgUpdateParamsResponse",

		// accountplus
		"/vindax.accountplus.MsgSetActiveState",
		"/vindax.accountplus.MsgSetActiveStateResponse",

		// affiliates
		"/vindax.affiliates.MsgUpdateAffiliateOverrides",
		"/vindax.affiliates.MsgUpdateAffiliateOverridesResponse",
		"/vindax.affiliates.MsgUpdateAffiliateParameters",
		"/vindax.affiliates.MsgUpdateAffiliateParametersResponse",
		"/vindax.affiliates.MsgUpdateAffiliateTiers",
		"/vindax.affiliates.MsgUpdateAffiliateTiersResponse",
		"/vindax.affiliates.MsgUpdateAffiliateWhitelist",
		"/vindax.affiliates.MsgUpdateAffiliateWhitelistResponse",

		// blocktime
		"/vindax.blocktime.MsgUpdateDowntimeParams",
		"/vindax.blocktime.MsgUpdateDowntimeParamsResponse",
		"/vindax.blocktime.MsgUpdateSynchronyParams",
		"/vindax.blocktime.MsgUpdateSynchronyParamsResponse",

		// bridge
		"/vindax.bridge.MsgCompleteBridge",
		"/vindax.bridge.MsgCompleteBridgeResponse",
		"/vindax.bridge.MsgUpdateEventParams",
		"/vindax.bridge.MsgUpdateEventParamsResponse",
		"/vindax.bridge.MsgUpdateProposeParams",
		"/vindax.bridge.MsgUpdateProposeParamsResponse",
		"/vindax.bridge.MsgUpdateSafetyParams",
		"/vindax.bridge.MsgUpdateSafetyParamsResponse",

		// clob
		"/vindax.clob.MsgCreateClobPair",
		"/vindax.clob.MsgCreateClobPairResponse",
		"/vindax.clob.MsgUpdateBlockRateLimitConfiguration",
		"/vindax.clob.MsgUpdateBlockRateLimitConfigurationResponse",
		"/vindax.clob.MsgUpdateClobPair",
		"/vindax.clob.MsgUpdateClobPairResponse",
		"/vindax.clob.MsgUpdateEquityTierLimitConfiguration",
		"/vindax.clob.MsgUpdateEquityTierLimitConfigurationResponse",
		"/vindax.clob.MsgUpdateLiquidationsConfig",
		"/vindax.clob.MsgUpdateLiquidationsConfigResponse",

		// delaymsg
		"/vindax.delaymsg.MsgDelayMessage",
		"/vindax.delaymsg.MsgDelayMessageResponse",

		// feetiers
		"/vindax.feetiers.MsgSetMarketFeeDiscountParams",
		"/vindax.feetiers.MsgSetMarketFeeDiscountParamsResponse",
		"/vindax.feetiers.MsgSetStakingTiers",
		"/vindax.feetiers.MsgSetStakingTiersResponse",
		"/vindax.feetiers.MsgUpdatePerpetualFeeParams",
		"/vindax.feetiers.MsgUpdatePerpetualFeeParamsResponse",

		// govplus
		"/vindax.govplus.MsgSlashValidator",
		"/vindax.govplus.MsgSlashValidatorResponse",

		// listing
		"/vindax.listing.MsgSetListingVaultDepositParams",
		"/vindax.listing.MsgSetListingVaultDepositParamsResponse",
		"/vindax.listing.MsgSetMarketsHardCap",
		"/vindax.listing.MsgSetMarketsHardCapResponse",
		"/vindax.listing.MsgUpgradeIsolatedPerpetualToCross",
		"/vindax.listing.MsgUpgradeIsolatedPerpetualToCrossResponse",

		// perpeutals
		"/vindax.perpetuals.MsgCreatePerpetual",
		"/vindax.perpetuals.MsgCreatePerpetualResponse",
		"/vindax.perpetuals.MsgSetLiquidityTier",
		"/vindax.perpetuals.MsgSetLiquidityTierResponse",
		"/vindax.perpetuals.MsgUpdateParams",
		"/vindax.perpetuals.MsgUpdateParamsResponse",
		"/vindax.perpetuals.MsgUpdatePerpetualParams",
		"/vindax.perpetuals.MsgUpdatePerpetualParamsResponse",

		// prices
		"/vindax.prices.MsgCreateOracleMarket",
		"/vindax.prices.MsgCreateOracleMarketResponse",
		"/vindax.prices.MsgUpdateMarketParam",
		"/vindax.prices.MsgUpdateMarketParamResponse",

		// ratelimit
		"/vindax.ratelimit.MsgSetLimitParams",
		"/vindax.ratelimit.MsgSetLimitParamsResponse",

		// revshare
		"/vindax.revshare.MsgSetMarketMapperRevShareDetailsForMarket",
		"/vindax.revshare.MsgSetMarketMapperRevShareDetailsForMarketResponse",
		"/vindax.revshare.MsgSetMarketMapperRevenueShare",
		"/vindax.revshare.MsgSetMarketMapperRevenueShareResponse",
		"/vindax.revshare.MsgSetOrderRouterRevShare",
		"/vindax.revshare.MsgSetOrderRouterRevShareResponse",
		"/vindax.revshare.MsgUpdateUnconditionalRevShareConfig",
		"/vindax.revshare.MsgUpdateUnconditionalRevShareConfigResponse",

		// rewards
		"/vindax.rewards.MsgUpdateParams",
		"/vindax.rewards.MsgUpdateParamsResponse",

		// sending
		"/vindax.sending.MsgSendFromModuleToAccount",
		"/vindax.sending.MsgSendFromModuleToAccountResponse",

		// stats
		"/vindax.stats.MsgUpdateParams",
		"/vindax.stats.MsgUpdateParamsResponse",

		// vault
		"/vindax.vault.MsgUnlockShares",
		"/vindax.vault.MsgUnlockSharesResponse",
		"/vindax.vault.MsgUpdateOperatorParams",
		"/vindax.vault.MsgUpdateOperatorParamsResponse",

		// vest
		"/vindax.vest.MsgDeleteVestEntry",
		"/vindax.vest.MsgDeleteVestEntryResponse",
		"/vindax.vest.MsgSetVestEntry",
		"/vindax.vest.MsgSetVestEntryResponse",
	}

	require.Equal(t, expectedMsgs, lib.GetSortedKeys[sort.StringSlice](msgs.InternalMsgSamplesGovAuth))
}

func TestInternalMsgSamples_Gov_Value(t *testing.T) {
	validateMsgValue(t, msgs.InternalMsgSamplesGovAuth)
}
