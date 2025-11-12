package prepare

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	gometrics "github.com/hashicorp/go-metrics"

	"github.com/danielvindax/vd-chain/protocol/lib/ante"
	"github.com/danielvindax/vd-chain/protocol/lib/metrics"
)

// GetGroupMsgOther returns two separate slices of byte txs given a single slice of byte txs and max bytes.
// The first slice contains the first N txs where the total bytes of the N txs is <= max bytes.
// The second slice contains the rest of txs, if any.
func GetGroupMsgOther(availableTxs [][]byte, maxBytes uint64) ([][]byte, [][]byte) {
	var (
		txsToInclude [][]byte
		txsRemainder [][]byte
		byteCount    uint64
	)

	for _, tx := range availableTxs {
		byteCount += uint64(len(tx))
		if byteCount <= maxBytes {
			txsToInclude = append(txsToInclude, tx)
		} else {
			txsRemainder = append(txsRemainder, tx)
		}
	}

	return txsToInclude, txsRemainder
}

// RemoveDisallowMsgs removes any txs that contain a disallowed msg.
func RemoveDisallowMsgs(
	ctx sdk.Context,
	decoder sdk.TxDecoder,
	txs [][]byte,
) [][]byte {
	defer telemetry.ModuleMeasureSince(
		ModuleName,
		time.Now(),
		metrics.RemoveDisallowMsgs,
		metrics.Latency,
	)

	ctx.Logger().Info("RemoveDisallowMsgs: starting to filter txs",
		"total_txs", len(txs),
	)
	var filteredTxs [][]byte
	for i, txBytes := range txs {
		// Decode tx so we can read msgs.
		ctx.Logger().Debug("RemoveDisallowMsgs: decoding tx",
			"tx_index", i,
			"tx_bytes_len", len(txBytes),
		)
		tx, err := decoder(txBytes)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("RemoveDisallowMsgs: failed to decode tx (index %v of %v txs): %v", i, len(txs), err))
			continue // continue to next tx.
		}
		
		ctx.Logger().Debug("RemoveDisallowMsgs: decoded tx",
			"tx_index", i,
			"num_msgs", len(tx.GetMsgs()),
		)

		ctx.Logger().Debug("RemoveDisallowMsgs: decoded tx",
			"tx_index", i,
			"num_msgs", len(tx.GetMsgs()),
		)

		// For each msg in tx, check if it is disallowed.
		containsDisallowMsg := false
		for j, msg := range tx.GetMsgs() {
			msgType := proto.MessageName(msg)
			ctx.Logger().Debug("RemoveDisallowMsgs: checking msg",
				"tx_index", i,
				"msg_index", j,
				"msg_type", msgType,
			)
			if ante.IsDisallowExternalSubmitMsg(msg) {
				ctx.Logger().Info("RemoveDisallowMsgs: found disallowed external submit msg",
					"tx_index", i,
					"msg_index", j,
					"msg_type", msgType,
				)
				telemetry.IncrCounterWithLabels(
					[]string{ModuleName, metrics.RemoveDisallowMsgs, metrics.DisallowMsg, metrics.Count},
					1,
					[]gometrics.Label{metrics.GetLabelForStringValue(metrics.Detail, proto.MessageName(msg))},
				)
				containsDisallowMsg = true
				break // break out of loop over msgs.
			}
			// Check for CLOB order msgs that should be disallowed
			// Note: This check is also done in ProcessProposal, but we log it here for visibility
			ctx.Logger().Debug("RemoveDisallowMsgs: msg passed external submit check",
				"tx_index", i,
				"msg_index", j,
				"msg_type", msgType,
			)
		}

		// If tx contains disallowed msg, skip it.
		if containsDisallowMsg {
			ctx.Logger().Info("RemoveDisallowMsgs: skipping tx with disallowed msg",
				"tx_index", i,
				"tx_bytes_len", len(txBytes),
				"num_msgs", len(tx.GetMsgs()),
			)
			continue // continue to next tx.
		}

		// Otherwise, add tx to filtered txs.
		ctx.Logger().Debug("RemoveDisallowMsgs: tx passed all checks, adding to filtered list",
			"tx_index", i,
			"num_msgs", len(tx.GetMsgs()),
		)
		filteredTxs = append(filteredTxs, txBytes)
	}
	
	ctx.Logger().Info("RemoveDisallowMsgs: finished filtering txs",
		"original_count", len(txs),
		"filtered_count", len(filteredTxs),
		"removed_count", len(txs)-len(filteredTxs),
	)

	ctx.Logger().Info("RemoveDisallowMsgs: finished filtering txs",
		"original_count", len(txs),
		"filtered_count", len(filteredTxs),
		"removed_count", len(txs)-len(filteredTxs),
	)

	return filteredTxs
}
