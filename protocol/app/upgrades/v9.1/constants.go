package v_9_1

import (
	store "cosmossdk.io/store/types"
	"github.com/danielvindax/vd-chain/protocol/app/upgrades"
)

const (
	UpgradeName = "v9.1"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:   UpgradeName,
	StoreUpgrades: store.StoreUpgrades{},
}
