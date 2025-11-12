package process

import (
	"slices"

	errorsmod "cosmossdk.io/errors"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danielvindax/vd-chain/protocol/lib"
)

const (
	minTxsCount                   = 4
	proposedOperationsTxIndex     = 0
	updateMarketPricesTxLenOffset = -1
	addPremiumVotesTxLenOffset    = -2
	acknowledgeBridgesTxLenOffset = -3
	lastOtherTxLenOffset          = acknowledgeBridgesTxLenOffset
	firstOtherTxIndex             = proposedOperationsTxIndex + 1
)

func init() {
	txIndicesAndOffsets := []int{
		proposedOperationsTxIndex,
		acknowledgeBridgesTxLenOffset,
		addPremiumVotesTxLenOffset,
		updateMarketPricesTxLenOffset,
	}
	if minTxsCount != len(txIndicesAndOffsets) {
		panic("minTxsCount does not match expected count of Txs.")
	}
	if lib.ContainsDuplicates(txIndicesAndOffsets) {
		panic("Duplicate indices/offsets defined for Txs.")
	}
	if slices.Min[[]int](txIndicesAndOffsets) != lastOtherTxLenOffset {
		panic("lastTxLenOffset is not the lowest offset")
	}
	if slices.Max[[]int](txIndicesAndOffsets)+1 != firstOtherTxIndex {
		panic("firstOtherTxIndex is <= the maximum offset")
	}
	txIndicesForMinTxsCount := []int{
		proposedOperationsTxIndex,
		acknowledgeBridgesTxLenOffset + minTxsCount,
		addPremiumVotesTxLenOffset + minTxsCount,
		updateMarketPricesTxLenOffset + minTxsCount,
	}
	if minTxsCount != len(txIndicesForMinTxsCount) {
		panic("minTxsCount does not match expected count of Txs.")
	}
	if lib.ContainsDuplicates(txIndicesForMinTxsCount) {
		panic("Overlapping indices and offsets defined for Txs.")
	}
	if minTxsCount != firstOtherTxIndex-lastOtherTxLenOffset {
		panic("Unexpected gap between firstOtherTxIndex and lastOtherTxLenOffset which is greater than minTxsCount")
	}
}

// ProcessProposalTxs is used as an intermediary struct to validate a proposed list of txs
// for `ProcessProposal`.
type ProcessProposalTxs struct {
	// Single msg txs.
	ProposedOperationsTx *ProposedOperationsTx
	AcknowledgeBridgesTx *AcknowledgeBridgesTx
	AddPremiumVotesTx    *AddPremiumVotesTx
	UpdateMarketPricesTx *UpdateMarketPricesTx // abstract over MarketPriceUpdates from VEs or default.

	// Multi msgs txs.
	OtherTxs []*OtherMsgsTx
}

// DecodeProcessProposalTxs returns a new `processProposalTxs`.
func DecodeProcessProposalTxs(
	ctx sdk.Context,
	decoder sdk.TxDecoder,
	req *abci.RequestProcessProposal,
	bridgeKeeper ProcessBridgeKeeper,
	pricesTxDecoder UpdateMarketPriceTxDecoder,
) (*ProcessProposalTxs, error) {
	ctx.Logger().Info("DecodeProcessProposalTxs: starting",
		"total_txs", len(req.Txs),
		"height", ctx.BlockHeight(),
	)
	
	// Check len (accounting for offset from injected vote-extensions if applicable)
	offset := pricesTxDecoder.GetTxOffset(ctx)
	injectedTxCount := minTxsCount + offset
	numTxs := len(req.Txs)
	ctx.Logger().Debug("DecodeProcessProposalTxs: tx count check",
		"offset", offset,
		"injected_tx_count", injectedTxCount,
		"num_txs", numTxs,
		"min_required", injectedTxCount,
	)
	if numTxs < injectedTxCount {
		ctx.Logger().Error("DecodeProcessProposalTxs: insufficient txs",
			"expected_min", injectedTxCount,
			"actual", numTxs,
		)
		return nil, errorsmod.Wrapf(
			ErrUnexpectedNumMsgs,
			"Expected the proposal to contain at least %d txs, but got %d",
			injectedTxCount,
			numTxs,
		)
	}

	// Price updates.
	ctx.Logger().Debug("DecodeProcessProposalTxs: decoding UpdateMarketPricesTx")
	updatePricesTx, err := pricesTxDecoder.DecodeUpdateMarketPricesTx(
		ctx,
		req.Txs,
	)
	if err != nil {
		ctx.Logger().Error("DecodeProcessProposalTxs: failed to decode UpdateMarketPricesTx", "error", err)
		return nil, err
	}
	ctx.Logger().Debug("DecodeProcessProposalTxs: successfully decoded UpdateMarketPricesTx")

	// Operations.
	// if vote-extensions were injected, offset will be incremented.
	operationsTxIndex := proposedOperationsTxIndex + offset
	ctx.Logger().Debug("DecodeProcessProposalTxs: decoding ProposedOperationsTx",
		"tx_index", operationsTxIndex,
	)
	operationsTx, err := DecodeProposedOperationsTx(decoder, req.Txs[operationsTxIndex])
	if err != nil {
		ctx.Logger().Error("DecodeProcessProposalTxs: failed to decode ProposedOperationsTx",
			"tx_index", operationsTxIndex,
			"error", err,
		)
		return nil, err
	}
	ctx.Logger().Debug("DecodeProcessProposalTxs: successfully decoded ProposedOperationsTx")

	// Acknowledge bridges.
	acknowledgeBridgesTxIndex := numTxs + acknowledgeBridgesTxLenOffset
	ctx.Logger().Debug("DecodeProcessProposalTxs: decoding AcknowledgeBridgesTx",
		"tx_index", acknowledgeBridgesTxIndex,
	)
	acknowledgeBridgesTx, err := DecodeAcknowledgeBridgesTx(
		ctx,
		bridgeKeeper,
		decoder,
		req.Txs[acknowledgeBridgesTxIndex],
	)
	if err != nil {
		ctx.Logger().Error("DecodeProcessProposalTxs: failed to decode AcknowledgeBridgesTx",
			"tx_index", acknowledgeBridgesTxIndex,
			"error", err,
		)
		return nil, err
	}
	ctx.Logger().Debug("DecodeProcessProposalTxs: successfully decoded AcknowledgeBridgesTx")

	// Funding samples.
	addPremiumVotesTxIndex := numTxs + addPremiumVotesTxLenOffset
	ctx.Logger().Debug("DecodeProcessProposalTxs: decoding AddPremiumVotesTx",
		"tx_index", addPremiumVotesTxIndex,
	)
	addPremiumVotesTx, err := DecodeAddPremiumVotesTx(decoder, req.Txs[addPremiumVotesTxIndex])
	if err != nil {
		ctx.Logger().Error("DecodeProcessProposalTxs: failed to decode AddPremiumVotesTx",
			"tx_index", addPremiumVotesTxIndex,
			"error", err,
		)
		return nil, err
	}
	ctx.Logger().Debug("DecodeProcessProposalTxs: successfully decoded AddPremiumVotesTx")

	// Other txs.
	// if vote-extensions were injected, offset will be incremented.
	allOtherTxs := make([]*OtherMsgsTx, numTxs-injectedTxCount)
	ctx.Logger().Info("DecodeProcessProposalTxs: starting to decode OtherTxs",
		"total_other_txs", len(allOtherTxs),
		"first_other_tx_index", firstOtherTxIndex+offset,
		"last_other_tx_index", numTxs+lastOtherTxLenOffset,
	)
	for i, txBytes := range req.Txs[firstOtherTxIndex+offset : numTxs+lastOtherTxLenOffset] {
		ctx.Logger().Debug("DecodeProcessProposalTxs: decoding OtherTx",
			"tx_index", i,
			"tx_bytes_len", len(txBytes),
			"absolute_index", firstOtherTxIndex+offset+i,
		)
		otherTx, err := DecodeOtherMsgsTx(decoder, txBytes)
		if err != nil {
			ctx.Logger().Error("DecodeProcessProposalTxs: failed to decode OtherTx",
				"tx_index", i,
				"absolute_index", firstOtherTxIndex+offset+i,
				"tx_bytes_len", len(txBytes),
				"error", err,
			)
			return nil, err
		}

		ctx.Logger().Debug("DecodeProcessProposalTxs: successfully decoded OtherTx",
			"tx_index", i,
			"num_msgs", len(otherTx.GetMsgs()),
		)
		allOtherTxs[i] = otherTx
	}
	ctx.Logger().Info("DecodeProcessProposalTxs: finished decoding OtherTxs",
		"total_decoded", len(allOtherTxs),
	)

	ctx.Logger().Info("DecodeProcessProposalTxs: successfully decoded all txs",
		"total_other_txs", len(allOtherTxs),
	)
	return &ProcessProposalTxs{
		ProposedOperationsTx: operationsTx,
		AcknowledgeBridgesTx: acknowledgeBridgesTx,
		AddPremiumVotesTx:    addPremiumVotesTx,
		UpdateMarketPricesTx: updatePricesTx,
		OtherTxs:             allOtherTxs,
	}, nil
}

// Validate performs `ValidateBasic` on the underlying msgs that are part of the txs.
// Returns nil if all are valid. Otherwise, returns error.
//
// Exception: for UpdateMarketPricesTx, we perform "in-memory stateful" validation
// to ensure that the new proposed prices are "valid" in comparison to index prices.
func (ppt *ProcessProposalTxs) Validate() error {
	// Validate single msg txs.
	singleTxs := []SingleMsgTx{
		ppt.ProposedOperationsTx,
		ppt.AddPremiumVotesTx,
		ppt.AcknowledgeBridgesTx,
		ppt.UpdateMarketPricesTx,
	}
	for _, smt := range singleTxs {
		if err := smt.Validate(); err != nil {
			return err
		}
	}

	// Validate multi msgs txs.
	for _, mmt := range ppt.OtherTxs {
		if err := mmt.Validate(); err != nil {
			return err
		}
	}

	return nil
}
