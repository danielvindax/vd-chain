package grpc

import pricetypes "github.com/danielvindax/vd-chain/protocol/x/prices/types"

type QueryServer interface {
	pricetypes.QueryServer
}
