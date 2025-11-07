package main

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	appconfig "github.com/dydxprotocol/v4-chain/protocol/app/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	// Cosmos SDK modules
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// IBC
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"

	// dYdX v4 modules
	bridgemoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	marketmapmoduletypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	perpetualsmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	rewardsmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vestmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	vaultmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

type modItem struct {
	Label string
	Name  string
}

func main() {
	// Set HRP "vindax" (required before deriving address)
	appconfig.SetAddressPrefixes()

	fmt.Println("🔹 Module Accounts for vindax chain:")

	// Full list of account names to derive module address
	// Note: With Cosmos, most use ModuleName; some pools/IF use separate constants.
	items := []modItem{
		// Cosmos SDK & Core
		{Label: "fee_collector", Name: authtypes.FeeCollectorName},
		{Label: "distribution", Name: distrtypes.ModuleName},
		{Label: "gov", Name: govtypes.ModuleName},
		{Label: "staking_bonded_pool", Name: stakingtypes.BondedPoolName},
		{Label: "staking_not_bonded_pool", Name: stakingtypes.NotBondedPoolName},

		// IBC
		{Label: "ibc_transfer", Name: ibctransfertypes.ModuleName},
		{Label: "ica_host", Name: icatypes.ModuleName},

		// dYdX v4 specific modules
		{Label: "bridge", Name: bridgemoduletypes.ModuleName},
		{Label: "marketmap", Name: marketmapmoduletypes.ModuleName},
		{Label: "perpetuals_insurance_fund", Name: perpetualsmoduletypes.InsuranceFundName},
		{Label: "rewards_treasury", Name: rewardsmoduletypes.TreasuryAccountName},
		{Label: "rewards_vester", Name: rewardsmoduletypes.VesterAccountName},
		{Label: "vest_community_treasury", Name: vestmoduletypes.CommunityTreasuryAccountName},
		{Label: "vest_community_vester", Name: vestmoduletypes.CommunityVesterAccountName},
		{Label: "satypes", Name: satypes.ModuleName},
		{Label: "vault_megavault", Name: vaultmoduletypes.MegavaultAccountName},
		{Label: "delaymsg", Name: delaymsgtypes.ModuleName},

		// (project-specific) If you use these modules with module account, add:
		// {Label: "clob", Name: "clob"},
		// {Label: "perpetuals", Name: "perpetuals"},
		// {Label: "vault", Name: "vault"},
		// {Label: "rewards", Name: "rewards"},
		// {Label: "stats", Name: "stats"},
	}

	// Print all module addresses
	for _, it := range items {
		addr := authtypes.NewModuleAddress(it.Name)
		fmt.Printf("%-28s → %s\n", it.Label, addr.String())
	}

	// Example: print additional "module vault" and "vault subaccount" addresses
	vaultModuleAddr := authtypes.NewModuleAddress("vault")
	fmt.Printf("%-28s → %s\n", "vault(module_name='vault')", vaultModuleAddr.String())

	// Subaccount: actually also AccAddress from module address
	vaultAcc := sdk.AccAddress(vaultModuleAddr)
	fmt.Printf("%-28s → %s\n", "vault Acc", vaultAcc.String())
}
