package v_8_2

import (
	store "cosmossdk.io/store/types"
	"github.com/danielvindax/vd-chain/protocol/app/upgrades"
)

const (
	UpgradeName = "v8.2"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:   UpgradeName,
	StoreUpgrades: store.StoreUpgrades{},
}
