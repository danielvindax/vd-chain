package v_9_3

import (
	store "cosmossdk.io/store/types"
	"github.com/danielvindax/vd-chain/protocol/app/upgrades"
)

const (
	UpgradeName = "v9.3"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:   UpgradeName,
	StoreUpgrades: store.StoreUpgrades{},
}
