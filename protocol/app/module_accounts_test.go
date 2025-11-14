package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/danielvindax/vd-chain/protocol/app"
	bridgemoduletypes "github.com/danielvindax/vd-chain/protocol/x/bridge/types"
	perpetualsmoduletypes "github.com/danielvindax/vd-chain/protocol/x/perpetuals/types"
	rewardsmoduletypes "github.com/danielvindax/vd-chain/protocol/x/rewards/types"
	satypes "github.com/danielvindax/vd-chain/protocol/x/subaccounts/types"
	vaultmoduletypes "github.com/danielvindax/vd-chain/protocol/x/vault/types"
	vestmoduletypes "github.com/danielvindax/vd-chain/protocol/x/vest/types"
	marketmapmoduletypes "github.com/dydxprotocol/slinky/x/marketmap/types"
)

func TestModuleAccountsToAddresses(t *testing.T) {
	expectedModuleAccToAddresses := map[string]string{
		authtypes.FeeCollectorName:                   "vindax17xpfvakm2amg962yls6f84z3kell8c5les5vz4",
		bridgemoduletypes.ModuleName:                 "vindax1zlefkpe3g0vvm9a4h0jf9000lmqutlh9j7tmen",
		distrtypes.ModuleName:                        "vindax1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8wknsqh",
		stakingtypes.BondedPoolName:                  "vindax1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3uj7rsl",
		stakingtypes.NotBondedPoolName:               "vindax1tygms3xhhs3yv487phx3dw4a95jn7t7lgjzjxt",
		govtypes.ModuleName:                          "vindax10d07y265gmmuvt4z0w9aw880jnsr700jntyflm",
		ibctransfertypes.ModuleName:                  "vindax1yl6hdjhmkf37639730gffanpzndzdpmh8kp97t",
		satypes.ModuleName:                           "vindax1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5vcx39",
		perpetualsmoduletypes.InsuranceFundName:      "vindax1c7ptc87hkd54e3r7zjy92q29xkq7t79w69fh2l",
		rewardsmoduletypes.TreasuryAccountName:       "vindax16wrau2x4tsg033xfrrdpae6kxfn9kyuernd6m7",
		rewardsmoduletypes.VesterAccountName:         "vindax1ltyc6y4skclzafvpznpt2qjwmfwgsndp4y7tj7",
		vestmoduletypes.CommunityTreasuryAccountName: "vindax15ztc7xy42tn2ukkc0qjthkucw9ac63pgpwk52v",
		vestmoduletypes.CommunityVesterAccountName:   "vindax1wxje320an3karyc6mjw4zghs300dmrjkwr8wzf",
		icatypes.ModuleName:                          "vindax1vlthgax23ca9syk7xgaz347xmf4nunefwppm9c",
		marketmapmoduletypes.ModuleName:              "vindax16j3d86dww8p2rzdlqsv7wle98cxzjxw6gztvtv",
		vaultmoduletypes.MegavaultAccountName:        "vindax18tkxrnrkqc2t0lr3zxr5g6a4hdvqksylxst3uu",
	}

	require.True(t, len(expectedModuleAccToAddresses) == len(app.GetMaccPerms()),
		"expected %d, got %d", len(expectedModuleAccToAddresses), len(app.GetMaccPerms()))
	for acc, address := range expectedModuleAccToAddresses {
		expectedAddr := authtypes.NewModuleAddress(acc).String()
		require.Equal(t, address, expectedAddr, "module (%v) should have address (%s)", acc, expectedAddr)
	}
}

func TestBlockedAddresses(t *testing.T) {
	expectedBlockedAddresses := map[string]bool{
		"vindax17xpfvakm2amg962yls6f84z3kell8c5les5vz4": true,
		"vindax1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8wknsqh": true,
		"vindax1tygms3xhhs3yv487phx3dw4a95jn7t7lgjzjxt": true,
		"vindax1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3uj7rsl": true,
		"vindax1yl6hdjhmkf37639730gffanpzndzdpmh8kp97t": true,
		"vindax1vlthgax23ca9syk7xgaz347xmf4nunefwppm9c": true,
	}
	require.Equal(t, expectedBlockedAddresses, app.BlockedAddresses())
}

func TestMaccPerms(t *testing.T) {
	maccPerms := app.GetMaccPerms()
	expectedMaccPerms := map[string][]string{
		"bonded_tokens_pool":     {"burner", "staking"},
		"bridge":                 {"minter"},
		"distribution":           nil,
		"fee_collector":          nil,
		"gov":                    {"burner"},
		"insurance_fund":         nil,
		"not_bonded_tokens_pool": {"burner", "staking"},
		"subaccounts":            nil,
		"transfer":               {"minter", "burner"},
		"interchainaccounts":     nil,
		"rewards_treasury":       nil,
		"rewards_vester":         nil,
		"community_treasury":     nil,
		"community_vester":       nil,
		"marketmap":              nil,
		"megavault":              nil,
	}
	require.Equal(t, expectedMaccPerms, maccPerms, "default macc perms list does not match expected")
}

func TestModuleAccountAddrs(t *testing.T) {
	expectedModuleAccAddresses := map[string]bool{
		"vindax17xpfvakm2amg962yls6f84z3kell8c5les5vz4": true, // x/auth.FeeCollector
		"vindax1zlefkpe3g0vvm9a4h0jf9000lmqutlh9j7tmen": true, // x/bridge
		"vindax1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8wknsqh": true, // x/distribution
		"vindax1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3uj7rsl": true, // x/staking.bondedPool
		"vindax1tygms3xhhs3yv487phx3dw4a95jn7t7lgjzjxt": true, // x/staking.notBondedPool
		"vindax10d07y265gmmuvt4z0w9aw880jnsr700jntyflm": true, // x/ gov
		"vindax1yl6hdjhmkf37639730gffanpzndzdpmh8kp97t": true, // ibc transfer
		"vindax1vlthgax23ca9syk7xgaz347xmf4nunefwppm9c": true, // interchainaccounts
		"vindax1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5vcx39": true, // x/subaccount
		"vindax1c7ptc87hkd54e3r7zjy92q29xkq7t79w69fh2l": true, // x/clob.insuranceFund
		"vindax16wrau2x4tsg033xfrrdpae6kxfn9kyuernd6m7": true, // x/rewards.treasury
		"vindax1ltyc6y4skclzafvpznpt2qjwmfwgsndp4y7tj7": true, // x/rewards.vester
		"vindax15ztc7xy42tn2ukkc0qjthkucw9ac63pgpwk52v": true, // x/vest.communityTreasury
		"vindax1wxje320an3karyc6mjw4zghs300dmrjkwr8wzf": true, // x/vest.communityVester
		"vindax16j3d86dww8p2rzdlqsv7wle98cxzjxw6gztvtv": true, // x/marketmap
		"vindax18tkxrnrkqc2t0lr3zxr5g6a4hdvqksylxst3uu": true, // x/vault.megavault
	}

	require.Equal(t, expectedModuleAccAddresses, app.ModuleAccountAddrs())
}
