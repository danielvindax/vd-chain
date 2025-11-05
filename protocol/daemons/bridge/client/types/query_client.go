package types

import bridgetypes "github.com/danielvindax/vd-chain/protocol/x/bridge/types"

// BridgeQueryClient is an interface that encapsulates the x/bridge `QueryClient` interface.
type BridgeQueryClient interface {
	bridgetypes.QueryClient
}
