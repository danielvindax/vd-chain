package keeper_test

import (
	"context"
	"testing"

	testapp "github.com/danielvindax/vd-chain/protocol/testutil/app"
	"github.com/danielvindax/vd-chain/protocol/x/bridge/keeper"
	"github.com/danielvindax/vd-chain/protocol/x/bridge/types"
)

func setupMsgServer(t *testing.T) (keeper.Keeper, types.MsgServer, context.Context) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	return k, keeper.NewMsgServerImpl(k), ctx
}
