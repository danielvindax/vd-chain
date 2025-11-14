# Test Documentation: CLOB App E2E Tests

## Overview

This test file verifies core **CLOB application** functionality including order hydration, concurrent operations, transaction validation, and statistics tracking. The tests ensure that:
1. Long-term orders are properly hydrated from state to memclob on startup
2. Orders can be matched during hydration
3. Concurrent matches and cancels work correctly
4. Transaction signature validation works correctly
5. Statistics are tracked correctly for trading activity

---

## Test Function: TestHydrationInPreBlocker

### Test Case: Success - Hydrate Long-Term Order from State

### Input
- **Genesis State:**
  - Long-term order exists in state (not in memclob)
  - Order: Carl_Num0, Buy 1 BTC at 50,000 USDC, GoodTilTime = 10
  - Order placed at block 1
  - Order expiration: time.Unix(50, 0)
- **Block:** Advance to block 2

### Output
- **Order in State:** Order exists in state storage
- **Order in MemClob:** Order is hydrated and exists in memclob
- **Order on Orderbook:** Order is visible on the orderbook

### Why It Runs This Way?

1. **State Hydration:** Long-term orders are stored in state but need to be loaded into memclob on startup.
2. **PreBlocker:** PreBlocker is called before each block to hydrate orders from state.
3. **Order Visibility:** Orders must be in memclob to be visible on the orderbook.
4. **Persistence:** Orders persist across restarts, so hydration is critical.

---

## Test Function: TestHydrationWithMatchPreBlocker

### Test Case: Success - Hydrate Orders That Match During Hydration

### Input
- **Genesis State:**
  - Carl: Long-term buy order (1 BTC at 50,000 USDC)
  - Dave: Long-term sell order (1 BTC at 50,000 USDC)
  - Both orders placed at block 1
  - Both orders expire at time.Unix(10, 0)
- **Block:** Advance to block 2

### Output
- **PreBlocker:** Orders are hydrated and matched during PreBlocker
- **State Changes Discarded:** State changes from PreBlocker are discarded (IsCheckTx = true)
- **Operations Queue:** Match operation is added to operations queue
- **Block 2:** Orders are fully filled and removed from state
- **Final State:**
  - Carl: Long 1 BTC, USDC balance decreased by 50,000 USDC
  - Dave: Short 1 BTC, USDC balance increased by 50,000 USDC

### Why It Runs This Way?

1. **Hydration Matching:** Orders that match during hydration should generate match operations.
2. **State Isolation:** PreBlocker runs in CheckTx context, so state changes are discarded.
3. **Operations Queue:** Matches are queued and processed in the next block.
4. **Full Fill:** Matching orders are fully filled and removed from state.

---

## Test Function: TestConcurrentMatchesAndCancels

### Test Case: Success - Concurrent Matches and Cancels

### Input
- **Accounts:** 1000 random accounts
- **Orders:**
  - 300 orders that match (150 buys, 150 sells)
    - 50 orders of size 5, 10, 15 each for both sides
    - Total matched volume: 1,500 quantums
  - 700 orders that are cancelled
    - Orders placed then immediately cancelled
- **Execution:** All CheckTx calls executed concurrently
- **Block:** Advance to block 3

### Output
- **Matched Orders:** All 300 orders are fully filled
- **Cancelled Orders:** All 700 orders are cancelled (no fills)
- **No Data Races:** Test passes with `-race` flag enabled

### Why It Runs This Way?

1. **Concurrency Testing:** Tests that the system handles concurrent operations correctly.
2. **Race Detection:** Uses Go's race detector to find data races.
3. **Mixed Operations:** Tests both matches and cancels happening concurrently.
4. **Stress Test:** 1000 accounts with concurrent operations stress tests the system.

---

## Test Function: TestFailsDeliverTxWithIncorrectlySignedPlaceOrderTx

### Test Case: Failure - Incorrectly Signed Order Placement

### Input
- **Order:** Alice's order (from Alice_Num0)
- **Signer:** Bob's private key (incorrect signer)
- **Transaction:** Order placement transaction signed by Bob

### Output
- **DeliverTx:** FAIL
- **Error:** "invalid pubkey: MsgProposedOperations is invalid"
- **Transaction:** Rejected

### Why It Runs This Way?

1. **Signature Validation:** Transactions must be signed by the correct account.
2. **Security:** Prevents unauthorized order placement.
3. **DeliverTx Validation:** Validation happens in DeliverTx, not just CheckTx.

---

## Test Function: TestFailsDeliverTxWithUnsignedTransactions

### Test Case: Failure - Unsigned Order Placement

### Input
- **Order:** Alice's order (from Alice_Num0)
- **Transaction:** Order placement transaction with no signatures

### Output
- **DeliverTx:** FAIL
- **Error:** "Error: no signatures supplied: MsgProposedOperations is invalid"
- **Transaction:** Rejected

### Why It Runs This Way?

1. **Signature Requirement:** All transactions must be signed.
2. **Security:** Prevents unsigned transactions from being processed.
3. **Validation:** Signature validation happens in DeliverTx.

---

## Test Function: TestStats

### Test Case: Success - Statistics Tracking

### Input
- **Epochs:** Multiple epochs with trading activity
- **Orders:**
  - Block 2-5: Alice (maker) and Bob (taker) trade 10,000 notional
  - Block 6: Alice and Bob trade 5,000 notional (same epoch)
  - Block 8: Alice and Bob trade 5,000 notional (new epoch)
- **Time:** Epochs advance based on StatsEpochDuration

### Output
- **User Stats:**
  - Alice: MakerNotional = 20,000, TakerNotional = 0
  - Bob: MakerNotional = 0, TakerNotional = 20,000
- **Global Stats:** NotionalTraded = 20,000
- **Epoch Stats:**
  - Epoch 0: 15,000 notional
  - Epoch 2: 5,000 notional
- **Window Expiration:** Stats expire after window duration

### Why It Runs This Way?

1. **Statistics Tracking:** Tracks trading activity for rewards and analytics.
2. **Maker/Taker:** Distinguishes between maker and taker roles.
3. **Epochs:** Statistics are tracked per epoch.
4. **Window:** Statistics expire after a window duration.
5. **Aggregation:** User stats, global stats, and epoch stats are all tracked.

---

## Flow Summary

### Order Hydration Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. GENESIS STATE                                             │
│    - Long-term orders stored in state                        │
│    - Orders not in memclob                                   │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. PREBLOCKER                                                │
│    - Load orders from state                                  │
│    - Hydrate orders into memclob                             │
│    - Check for matches                                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. ORDERBOOK                                                 │
│    - Orders visible on orderbook                             │
│    - Orders can be matched                                   │
└─────────────────────────────────────────────────────────────┘
```

### Concurrent Operations Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. CONCURRENT CHECKTX                                         │
│    - Multiple goroutines execute CheckTx                     │
│    - Orders placed and cancelled concurrently                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. BLOCK ADVANCEMENT                                         │
│    - All transactions included in block                      │
│    - Matches and cancels processed                           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. VERIFICATION                                              │
│    - Matched orders fully filled                             │
│    - Cancelled orders not filled                             │
│    - No data races detected                                  │
└─────────────────────────────────────────────────────────────┘
```

### Statistics Tracking Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. ORDER MATCHING                                            │
│    - Orders matched in block                                 │
│    - Maker and taker identified                              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. STATISTICS UPDATE                                         │
│    - User stats updated                                      │
│    - Global stats updated                                    │
│    - Epoch stats updated                                     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. EPOCH ADVANCEMENT                                          │
│    - New epoch starts                                        │
│    - Previous epoch stats preserved                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. WINDOW EXPIRATION                                         │
│    - Old epoch stats expire                                 │
│    - Stats removed from window                               │
└─────────────────────────────────────────────────────────────┘
```

### Key Points

1. **Order Hydration:**
   - Long-term orders must be hydrated from state on startup
   - Hydration happens in PreBlocker
   - Orders become visible on orderbook after hydration

2. **Concurrent Operations:**
   - System must handle concurrent CheckTx calls
   - Race detector helps find data races
   - Stress testing with many accounts

3. **Transaction Validation:**
   - Signatures must be valid
   - Signer must match transaction sender
   - Validation happens in DeliverTx

4. **Statistics:**
   - Maker/Taker distinction is important
   - Statistics tracked per epoch
   - Statistics expire after window duration

### Design Rationale

1. **State Hydration:** Ensures orders persist across restarts and are available for matching.

2. **Concurrency:** System must handle high concurrency in production.

3. **Security:** Signature validation prevents unauthorized transactions.

4. **Analytics:** Statistics enable rewards and analytics features.

