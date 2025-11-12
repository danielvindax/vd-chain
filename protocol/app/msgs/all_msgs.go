package msgs

import (
	"github.com/danielvindax/vd-chain/protocol/lib"
)

var (
	// AllTypeMessages is a list of all messages and types that are used in the app.
	// This list comes from the app's `InterfaceRegistry`.
	AllTypeMessages = map[string]struct{}{
		// auth
		"/cosmos.auth.v1beta1.BaseAccount":      {},
		"/cosmos.auth.v1beta1.ModuleAccount":    {},
		"/cosmos.auth.v1beta1.ModuleCredential": {},
		"/cosmos.auth.v1beta1.MsgUpdateParams":  {},

		// authz
		"/cosmos.authz.v1beta1.GenericAuthorization": {},
		"/cosmos.authz.v1beta1.MsgExec":              {},
		"/cosmos.authz.v1beta1.MsgExecResponse":      {},
		"/cosmos.authz.v1beta1.MsgGrant":             {},
		"/cosmos.authz.v1beta1.MsgGrantResponse":     {},
		"/cosmos.authz.v1beta1.MsgRevoke":            {},
		"/cosmos.authz.v1beta1.MsgRevokeResponse":    {},

		// bank
		"/cosmos.bank.v1beta1.MsgMultiSend":              {},
		"/cosmos.bank.v1beta1.MsgMultiSendResponse":      {},
		"/cosmos.bank.v1beta1.MsgSend":                   {},
		"/cosmos.bank.v1beta1.MsgSendResponse":           {},
		"/cosmos.bank.v1beta1.MsgSetSendEnabled":         {},
		"/cosmos.bank.v1beta1.MsgSetSendEnabledResponse": {},
		"/cosmos.bank.v1beta1.MsgUpdateParams":           {},
		"/cosmos.bank.v1beta1.MsgUpdateParamsResponse":   {},
		"/cosmos.bank.v1beta1.SendAuthorization":         {},
		"/cosmos.bank.v1beta1.Supply":                    {},

		// consensus
		"/cosmos.consensus.v1.MsgUpdateParams":         {},
		"/cosmos.consensus.v1.MsgUpdateParamsResponse": {},

		// crisis
		"/cosmos.crisis.v1beta1.MsgUpdateParams":            {},
		"/cosmos.crisis.v1beta1.MsgUpdateParamsResponse":    {},
		"/cosmos.crisis.v1beta1.MsgVerifyInvariant":         {},
		"/cosmos.crisis.v1beta1.MsgVerifyInvariantResponse": {},

		// crypto
		"/cosmos.crypto.ed25519.PrivKey":            {},
		"/cosmos.crypto.ed25519.PubKey":             {},
		"/cosmos.crypto.multisig.LegacyAminoPubKey": {},
		"/cosmos.crypto.secp256k1.PrivKey":          {},
		"/cosmos.crypto.secp256k1.PubKey":           {},
		"/cosmos.crypto.secp256r1.PubKey":           {},

		// distribution
		"/cosmos.distribution.v1beta1.CommunityPoolSpendProposal":             {},
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpend":                  {},
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpendResponse":          {},
		"/cosmos.distribution.v1beta1.MsgDepositValidatorRewardsPool":         {},
		"/cosmos.distribution.v1beta1.MsgDepositValidatorRewardsPoolResponse": {},
		"/cosmos.distribution.v1beta1.MsgFundCommunityPool":                   {},
		"/cosmos.distribution.v1beta1.MsgFundCommunityPoolResponse":           {},
		"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress":                  {},
		"/cosmos.distribution.v1beta1.MsgSetWithdrawAddressResponse":          {},
		"/cosmos.distribution.v1beta1.MsgUpdateParams":                        {},
		"/cosmos.distribution.v1beta1.MsgUpdateParamsResponse":                {},
		"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward":             {},
		"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorRewardResponse":     {},
		"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission":         {},
		"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommissionResponse": {},

		// evidence
		"/cosmos.evidence.v1beta1.Equivocation":              {},
		"/cosmos.evidence.v1beta1.MsgSubmitEvidence":         {},
		"/cosmos.evidence.v1beta1.MsgSubmitEvidenceResponse": {},

		// feegrant
		"/cosmos.feegrant.v1beta1.AllowedMsgAllowance":        {},
		"/cosmos.feegrant.v1beta1.BasicAllowance":             {},
		"/cosmos.feegrant.v1beta1.MsgGrantAllowance":          {},
		"/cosmos.feegrant.v1beta1.MsgGrantAllowanceResponse":  {},
		"/cosmos.feegrant.v1beta1.MsgPruneAllowances":         {},
		"/cosmos.feegrant.v1beta1.MsgPruneAllowancesResponse": {},
		"/cosmos.feegrant.v1beta1.MsgRevokeAllowance":         {},
		"/cosmos.feegrant.v1beta1.MsgRevokeAllowanceResponse": {},
		"/cosmos.feegrant.v1beta1.PeriodicAllowance":          {},

		// gov
		"/cosmos.gov.v1.MsgCancelProposal":              {},
		"/cosmos.gov.v1.MsgCancelProposalResponse":      {},
		"/cosmos.gov.v1.MsgDeposit":                     {},
		"/cosmos.gov.v1.MsgDepositResponse":             {},
		"/cosmos.gov.v1.MsgExecLegacyContent":           {},
		"/cosmos.gov.v1.MsgExecLegacyContentResponse":   {},
		"/cosmos.gov.v1.MsgSubmitProposal":              {},
		"/cosmos.gov.v1.MsgSubmitProposalResponse":      {},
		"/cosmos.gov.v1.MsgUpdateParams":                {},
		"/cosmos.gov.v1.MsgUpdateParamsResponse":        {},
		"/cosmos.gov.v1.MsgVote":                        {},
		"/cosmos.gov.v1.MsgVoteResponse":                {},
		"/cosmos.gov.v1.MsgVoteWeighted":                {},
		"/cosmos.gov.v1.MsgVoteWeightedResponse":        {},
		"/cosmos.gov.v1beta1.MsgDeposit":                {},
		"/cosmos.gov.v1beta1.MsgDepositResponse":        {},
		"/cosmos.gov.v1beta1.MsgSubmitProposal":         {},
		"/cosmos.gov.v1beta1.MsgSubmitProposalResponse": {},
		"/cosmos.gov.v1beta1.MsgVote":                   {},
		"/cosmos.gov.v1beta1.MsgVoteResponse":           {},
		"/cosmos.gov.v1beta1.MsgVoteWeighted":           {},
		"/cosmos.gov.v1beta1.MsgVoteWeightedResponse":   {},
		"/cosmos.gov.v1beta1.TextProposal":              {},

		// params
		"/cosmos.params.v1beta1.ParameterChangeProposal": {},

		// slashing
		"/cosmos.slashing.v1beta1.MsgUnjail":               {},
		"/cosmos.slashing.v1beta1.MsgUnjailResponse":       {},
		"/cosmos.slashing.v1beta1.MsgUpdateParams":         {},
		"/cosmos.slashing.v1beta1.MsgUpdateParamsResponse": {},

		// staking
		"/cosmos.staking.v1beta1.MsgBeginRedelegate":                   {},
		"/cosmos.staking.v1beta1.MsgBeginRedelegateResponse":           {},
		"/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation":         {},
		"/cosmos.staking.v1beta1.MsgCancelUnbondingDelegationResponse": {},
		"/cosmos.staking.v1beta1.MsgCreateValidator":                   {},
		"/cosmos.staking.v1beta1.MsgCreateValidatorResponse":           {},
		"/cosmos.staking.v1beta1.MsgDelegate":                          {},
		"/cosmos.staking.v1beta1.MsgDelegateResponse":                  {},
		"/cosmos.staking.v1beta1.MsgEditValidator":                     {},
		"/cosmos.staking.v1beta1.MsgEditValidatorResponse":             {},
		"/cosmos.staking.v1beta1.MsgSetProposers":                      {},
		"/cosmos.staking.v1beta1.MsgSetProposersResponse":              {},
		"/cosmos.staking.v1beta1.MsgUndelegate":                        {},
		"/cosmos.staking.v1beta1.MsgUndelegateResponse":                {},
		"/cosmos.staking.v1beta1.MsgUpdateParams":                      {},
		"/cosmos.staking.v1beta1.MsgUpdateParamsResponse":              {},
		"/cosmos.staking.v1beta1.StakeAuthorization":                   {},

		// tx
		"/cosmos.tx.v1beta1.Tx": {},

		// upgrade
		"/cosmos.upgrade.v1beta1.CancelSoftwareUpgradeProposal": {},
		"/cosmos.upgrade.v1beta1.MsgCancelUpgrade":              {},
		"/cosmos.upgrade.v1beta1.MsgCancelUpgradeResponse":      {},
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade":            {},
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgradeResponse":    {},
		"/cosmos.upgrade.v1beta1.SoftwareUpgradeProposal":       {},

		// affiliates
		"/vindax.affiliates.MsgRegisterAffiliate":                 {},
		"/vindax.affiliates.MsgRegisterAffiliateResponse":         {},
		"/vindax.affiliates.MsgUpdateAffiliateTiers":              {},
		"/vindax.affiliates.MsgUpdateAffiliateTiersResponse":      {},
		"/vindax.affiliates.MsgUpdateAffiliateWhitelist":          {},
		"/vindax.affiliates.MsgUpdateAffiliateWhitelistResponse":  {},
		"/vindax.affiliates.MsgUpdateAffiliateParameters":         {},
		"/vindax.affiliates.MsgUpdateAffiliateParametersResponse": {},
		"/vindax.affiliates.MsgUpdateAffiliateOverrides":          {},
		"/vindax.affiliates.MsgUpdateAffiliateOverridesResponse":  {},

		// accountplus
		"/vindax.accountplus.MsgAddAuthenticator":            {},
		"/vindax.accountplus.MsgAddAuthenticatorResponse":    {},
		"/vindax.accountplus.MsgRemoveAuthenticator":         {},
		"/vindax.accountplus.MsgRemoveAuthenticatorResponse": {},
		"/vindax.accountplus.MsgSetActiveState":              {},
		"/vindax.accountplus.MsgSetActiveStateResponse":      {},
		"/vindax.accountplus.TxExtension":                    {},

		// blocktime
		"/vindax.blocktime.MsgUpdateDowntimeParams":          {},
		"/vindax.blocktime.MsgUpdateDowntimeParamsResponse":  {},
		"/vindax.blocktime.MsgUpdateSynchronyParams":         {},
		"/vindax.blocktime.MsgUpdateSynchronyParamsResponse": {},

		// bridge
		"/vindax.bridge.MsgAcknowledgeBridges":          {},
		"/vindax.bridge.MsgAcknowledgeBridgesResponse":  {},
		"/vindax.bridge.MsgCompleteBridge":              {},
		"/vindax.bridge.MsgCompleteBridgeResponse":      {},
		"/vindax.bridge.MsgUpdateEventParams":           {},
		"/vindax.bridge.MsgUpdateEventParamsResponse":   {},
		"/vindax.bridge.MsgUpdateProposeParams":         {},
		"/vindax.bridge.MsgUpdateProposeParamsResponse": {},
		"/vindax.bridge.MsgUpdateSafetyParams":          {},
		"/vindax.bridge.MsgUpdateSafetyParamsResponse":  {},

		// clob
		"/vindax.clob.MsgBatchCancel":                                {},
		"/vindax.clob.MsgBatchCancelResponse":                        {},
		"/vindax.clob.MsgCancelOrder":                                {},
		"/vindax.clob.MsgCancelOrderResponse":                        {},
		"/vindax.clob.MsgCreateClobPair":                             {},
		"/vindax.clob.MsgCreateClobPairResponse":                     {},
		"/vindax.clob.MsgPlaceOrder":                                 {},
		"/vindax.clob.MsgPlaceOrderResponse":                         {},
		"/vindax.clob.MsgProposedOperations":                         {},
		"/vindax.clob.MsgProposedOperationsResponse":                 {},
		"/vindax.clob.MsgUpdateBlockRateLimitConfiguration":          {},
		"/vindax.clob.MsgUpdateBlockRateLimitConfigurationResponse":  {},
		"/vindax.clob.MsgUpdateClobPair":                             {},
		"/vindax.clob.MsgUpdateClobPairResponse":                     {},
		"/vindax.clob.MsgUpdateEquityTierLimitConfiguration":         {},
		"/vindax.clob.MsgUpdateEquityTierLimitConfigurationResponse": {},
		"/vindax.clob.MsgUpdateLiquidationsConfig":                   {},
		"/vindax.clob.MsgUpdateLiquidationsConfigResponse":           {},
		"/vindax.clob.MsgUpdateLeverage":                             {},
		"/vindax.clob.MsgUpdateLeverageResponse":                     {},

		// delaymsg
		"/vindax.delaymsg.MsgDelayMessage":         {},
		"/vindax.delaymsg.MsgDelayMessageResponse": {},

		// feetiers
		"/vindax.feetiers.MsgUpdatePerpetualFeeParams":           {},
		"/vindax.feetiers.MsgUpdatePerpetualFeeParamsResponse":   {},
		"/vindax.feetiers.MsgSetMarketFeeDiscountParams":         {},
		"/vindax.feetiers.MsgSetMarketFeeDiscountParamsResponse": {},
		"/vindax.feetiers.MsgSetStakingTiers":                    {},
		"/vindax.feetiers.MsgSetStakingTiersResponse":            {},

		// govplus
		"/vindax.govplus.MsgSlashValidator":         {},
		"/vindax.govplus.MsgSlashValidatorResponse": {},

		// listing
		"/vindax.listing.MsgSetMarketsHardCap":                       {},
		"/vindax.listing.MsgSetMarketsHardCapResponse":               {},
		"/vindax.listing.MsgCreateMarketPermissionless":              {},
		"/vindax.listing.MsgCreateMarketPermissionlessResponse":      {},
		"/vindax.listing.MsgSetListingVaultDepositParams":            {},
		"/vindax.listing.MsgSetListingVaultDepositParamsResponse":    {},
		"/vindax.listing.MsgUpgradeIsolatedPerpetualToCross":         {},
		"/vindax.listing.MsgUpgradeIsolatedPerpetualToCrossResponse": {},

		// perpetuals
		"/vindax.perpetuals.MsgAddPremiumVotes":               {},
		"/vindax.perpetuals.MsgAddPremiumVotesResponse":       {},
		"/vindax.perpetuals.MsgCreatePerpetual":               {},
		"/vindax.perpetuals.MsgCreatePerpetualResponse":       {},
		"/vindax.perpetuals.MsgSetLiquidityTier":              {},
		"/vindax.perpetuals.MsgSetLiquidityTierResponse":      {},
		"/vindax.perpetuals.MsgUpdateParams":                  {},
		"/vindax.perpetuals.MsgUpdateParamsResponse":          {},
		"/vindax.perpetuals.MsgUpdatePerpetualParams":         {},
		"/vindax.perpetuals.MsgUpdatePerpetualParamsResponse": {},

		// prices
		"/vindax.prices.MsgCreateOracleMarket":         {},
		"/vindax.prices.MsgCreateOracleMarketResponse": {},
		"/vindax.prices.MsgUpdateMarketPrices":         {},
		"/vindax.prices.MsgUpdateMarketPricesResponse": {},
		"/vindax.prices.MsgUpdateMarketParam":          {},
		"/vindax.prices.MsgUpdateMarketParamResponse":  {},

		// ratelimit
		"/vindax.ratelimit.MsgSetLimitParams":         {},
		"/vindax.ratelimit.MsgSetLimitParamsResponse": {},

		// sending
		"/vindax.sending.MsgCreateTransfer":                  {},
		"/vindax.sending.MsgCreateTransferResponse":          {},
		"/vindax.sending.MsgDepositToSubaccount":             {},
		"/vindax.sending.MsgDepositToSubaccountResponse":     {},
		"/vindax.sending.MsgWithdrawFromSubaccount":          {},
		"/vindax.sending.MsgWithdrawFromSubaccountResponse":  {},
		"/vindax.sending.MsgSendFromModuleToAccount":         {},
		"/vindax.sending.MsgSendFromModuleToAccountResponse": {},

		// stats
		"/vindax.stats.MsgUpdateParams":         {},
		"/vindax.stats.MsgUpdateParamsResponse": {},

		// vault
		"/vindax.vault.MsgAllocateToVault":                    {},
		"/vindax.vault.MsgAllocateToVaultResponse":            {},
		"/vindax.vault.MsgDepositToMegavault":                 {},
		"/vindax.vault.MsgDepositToMegavaultResponse":         {},
		"/vindax.vault.MsgRetrieveFromVault":                  {},
		"/vindax.vault.MsgRetrieveFromVaultResponse":          {},
		"/vindax.vault.MsgSetVaultParams":                     {},
		"/vindax.vault.MsgSetVaultParamsResponse":             {},
		"/vindax.vault.MsgSetVaultQuotingParams":              {}, // deprecated
		"/vindax.vault.MsgUnlockShares":                       {},
		"/vindax.vault.MsgUnlockSharesResponse":               {},
		"/vindax.vault.MsgUpdateDefaultQuotingParams":         {},
		"/vindax.vault.MsgUpdateDefaultQuotingParamsResponse": {},
		"/vindax.vault.MsgUpdateOperatorParams":               {},
		"/vindax.vault.MsgUpdateOperatorParamsResponse":       {},
		"/vindax.vault.MsgUpdateParams":                       {}, // deprecated
		"/vindax.vault.MsgWithdrawFromMegavault":              {},
		"/vindax.vault.MsgWithdrawFromMegavaultResponse":      {},

		// vest
		"/vindax.vest.MsgSetVestEntry":            {},
		"/vindax.vest.MsgSetVestEntryResponse":    {},
		"/vindax.vest.MsgDeleteVestEntry":         {},
		"/vindax.vest.MsgDeleteVestEntryResponse": {},

		// revshare
		"/vindax.revshare.MsgSetMarketMapperRevenueShare":                     {},
		"/vindax.revshare.MsgSetMarketMapperRevenueShareResponse":             {},
		"/vindax.revshare.MsgSetOrderRouterRevShare":                          {},
		"/vindax.revshare.MsgSetOrderRouterRevShareResponse":                  {},
		"/vindax.revshare.MsgSetMarketMapperRevShareDetailsForMarket":         {},
		"/vindax.revshare.MsgSetMarketMapperRevShareDetailsForMarketResponse": {},
		"/vindax.revshare.MsgUpdateUnconditionalRevShareConfig":               {},
		"/vindax.revshare.MsgUpdateUnconditionalRevShareConfigResponse":       {},

		// rewards
		"/vindax.rewards.MsgUpdateParams":         {},
		"/vindax.rewards.MsgUpdateParamsResponse": {},

		// ibc.applications
		"/ibc.applications.transfer.v1.MsgTransfer":             {},
		"/ibc.applications.transfer.v1.MsgTransferResponse":     {},
		"/ibc.applications.transfer.v1.MsgUpdateParams":         {},
		"/ibc.applications.transfer.v1.MsgUpdateParamsResponse": {},
		"/ibc.applications.transfer.v1.TransferAuthorization":   {},

		// ibc.core.channel
		"/ibc.core.channel.v1.Channel":                        {},
		"/ibc.core.channel.v1.Counterparty":                   {},
		"/ibc.core.channel.v1.MsgAcknowledgement":             {},
		"/ibc.core.channel.v1.MsgAcknowledgementResponse":     {},
		"/ibc.core.channel.v1.MsgChannelCloseConfirm":         {},
		"/ibc.core.channel.v1.MsgChannelCloseConfirmResponse": {},
		"/ibc.core.channel.v1.MsgChannelCloseInit":            {},
		"/ibc.core.channel.v1.MsgChannelCloseInitResponse":    {},
		"/ibc.core.channel.v1.MsgChannelOpenAck":              {},
		"/ibc.core.channel.v1.MsgChannelOpenAckResponse":      {},
		"/ibc.core.channel.v1.MsgChannelOpenConfirm":          {},
		"/ibc.core.channel.v1.MsgChannelOpenConfirmResponse":  {},
		"/ibc.core.channel.v1.MsgChannelOpenInit":             {},
		"/ibc.core.channel.v1.MsgChannelOpenInitResponse":     {},
		"/ibc.core.channel.v1.MsgChannelOpenTry":              {},
		"/ibc.core.channel.v1.MsgChannelOpenTryResponse":      {},
		"/ibc.core.channel.v1.MsgRecvPacket":                  {},
		"/ibc.core.channel.v1.MsgRecvPacketResponse":          {},
		"/ibc.core.channel.v1.MsgTimeout":                     {},
		"/ibc.core.channel.v1.MsgTimeoutOnClose":              {},
		"/ibc.core.channel.v1.MsgTimeoutOnCloseResponse":      {},
		"/ibc.core.channel.v1.MsgTimeoutResponse":             {},
		"/ibc.core.channel.v1.Packet":                         {},

		// ibc.core.client
		"/ibc.core.client.v1.ClientUpdateProposal":          {},
		"/ibc.core.client.v1.Height":                        {},
		"/ibc.core.client.v1.MsgCreateClient":               {},
		"/ibc.core.client.v1.MsgCreateClientResponse":       {},
		"/ibc.core.client.v1.MsgIBCSoftwareUpgrade":         {},
		"/ibc.core.client.v1.MsgIBCSoftwareUpgradeResponse": {},
		"/ibc.core.client.v1.MsgRecoverClient":              {},
		"/ibc.core.client.v1.MsgRecoverClientResponse":      {},
		"/ibc.core.client.v1.MsgSubmitMisbehaviour":         {},
		"/ibc.core.client.v1.MsgSubmitMisbehaviourResponse": {},
		"/ibc.core.client.v1.MsgUpdateClient":               {},
		"/ibc.core.client.v1.MsgUpdateClientResponse":       {},
		"/ibc.core.client.v1.MsgUpgradeClient":              {},
		"/ibc.core.client.v1.MsgUpgradeClientResponse":      {},
		"/ibc.core.client.v1.MsgUpdateParams":               {},
		"/ibc.core.client.v1.MsgUpdateParamsResponse":       {},
		"/ibc.core.client.v1.UpgradeProposal":               {},

		// ibc.core.commitment
		"/ibc.core.commitment.v1.MerklePath":   {},
		"/ibc.core.commitment.v1.MerklePrefix": {},
		"/ibc.core.commitment.v1.MerkleProof":  {},
		"/ibc.core.commitment.v1.MerkleRoot":   {},

		// ibc.core.connection
		"/ibc.core.connection.v1.ConnectionEnd":                    {},
		"/ibc.core.connection.v1.Counterparty":                     {},
		"/ibc.core.connection.v1.MsgConnectionOpenAck":             {},
		"/ibc.core.connection.v1.MsgConnectionOpenAckResponse":     {},
		"/ibc.core.connection.v1.MsgConnectionOpenConfirm":         {},
		"/ibc.core.connection.v1.MsgConnectionOpenConfirmResponse": {},
		"/ibc.core.connection.v1.MsgConnectionOpenInit":            {},
		"/ibc.core.connection.v1.MsgConnectionOpenInitResponse":    {},
		"/ibc.core.connection.v1.MsgConnectionOpenTry":             {},
		"/ibc.core.connection.v1.MsgConnectionOpenTryResponse":     {},
		"/ibc.core.connection.v1.MsgUpdateParams":                  {},
		"/ibc.core.connection.v1.MsgUpdateParamsResponse":          {},

		// ibc.lightclients
		"/ibc.lightclients.localhost.v2.ClientState":     {},
		"/ibc.lightclients.tendermint.v1.ClientState":    {},
		"/ibc.lightclients.tendermint.v1.ConsensusState": {},
		"/ibc.lightclients.tendermint.v1.Header":         {},
		"/ibc.lightclients.tendermint.v1.Misbehaviour":   {},

		// ica messages
		// Note: the `interchain_accounts.controller` messages are not actually used by the app,
		// since ICA Controller Keeper is initialized as nil.
		// However, since the ica.AppModuleBasic{} needs to be passed to basic_mananger as a whole, these messages
		// registered in the interface registry.
		"/ibc.applications.interchain_accounts.v1.InterchainAccount":                               {},
		"/ibc.applications.interchain_accounts.controller.v1.MsgSendTx":                            {},
		"/ibc.applications.interchain_accounts.controller.v1.MsgSendTxResponse":                    {},
		"/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccount":         {},
		"/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccountResponse": {},
		"/ibc.applications.interchain_accounts.controller.v1.MsgUpdateParams":                      {},
		"/ibc.applications.interchain_accounts.controller.v1.MsgUpdateParamsResponse":              {},
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParams":                            {},
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParamsResponse":                    {},

		// slinky marketmap messages
		"/slinky.marketmap.v1.MsgCreateMarkets":                   {},
		"/slinky.marketmap.v1.MsgCreateMarketsResponse":           {},
		"/slinky.marketmap.v1.MsgParams":                          {},
		"/slinky.marketmap.v1.MsgParamsResponse":                  {},
		"/slinky.marketmap.v1.MsgRemoveMarkets":                   {},
		"/slinky.marketmap.v1.MsgRemoveMarketsResponse":           {},
		"/slinky.marketmap.v1.MsgRemoveMarketAuthorities":         {},
		"/slinky.marketmap.v1.MsgRemoveMarketAuthoritiesResponse": {},
		"/slinky.marketmap.v1.MsgUpdateMarkets":                   {},
		"/slinky.marketmap.v1.MsgUpdateMarketsResponse":           {},
		"/slinky.marketmap.v1.MsgUpsertMarkets":                   {},
		"/slinky.marketmap.v1.MsgUpsertMarketsResponse":           {},
	}

	// DisallowMsgs are messages that cannot be externally submitted.
	DisallowMsgs = lib.MergeAllMapsMustHaveDistinctKeys(
		AppInjectedMsgSamples,
		InternalMsgSamplesAll,
		NestedMsgSamples,
		UnsupportedMsgSamples,
	)

	// AllowMsgs are messages that can be externally submitted.
	AllowMsgs = NormalMsgs
)
