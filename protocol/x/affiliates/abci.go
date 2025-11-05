package affiliates

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danielvindax/vd-chain/protocol/lib/log"
	"github.com/danielvindax/vd-chain/protocol/x/affiliates/keeper"
)

func EndBlocker(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) {
	if err := keeper.AggregateAffiliateReferredVolumeForFills(ctx); err != nil {
		log.ErrorLogWithError(ctx, "error aggregating affiliate volume for fills", err)
	}
}
