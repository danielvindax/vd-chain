# Test Documentation: Wind Down Market Proposal

## Overview

This test file verifies the **Wind Down Market** (market closure) functionality through governance proposals. The test ensures that when a CLOB pair transitions to `STATUS_FINAL_SETTLEMENT`, the system will:
1. Cancel all open stateful orders
2. Perform final settlement deleveraging
3. Block new order placement

---

## Test Case 1: Final Settlement Deleveraging - Non-negative TNC Accounts

### Input
- **Initial Subaccounts:**
  - `Carl_Num0`: 1 BTC Short position with 100,000 USD collateral
  - `Dave_Num0`: 1 BTC Long position with 50,000 USD collateral
- **ClobPair:** BTC-USD with initial status `STATUS_ACTIVE`
- **Proposal:** Transition ClobPair to `STATUS_FINAL_SETTLEMENT`

### Output
- **Subaccounts after final settlement:**
  - `Carl_Num0`: Only 50,000 USD remaining (short position closed at oracle price)
  - `Dave_Num0`: Has 100,000 USD (long position closed at oracle price)
- **ClobPair status:** `STATUS_FINAL_SETTLEMENT`
- **Events:** Indexer events emitted for ClobPair update

### Why It Runs This Way?

1. **Non-negative TNC (Total Net Collateral):** Both subaccounts have positive TNC, meaning sufficient collateral for settlement.
2. **Deleveraging at Oracle Price:** Since both accounts have positive TNC, deleveraging is performed at oracle price (reference price).
3. **Result:** 
   - Carl (short 1 BTC) must pay Dave (long 1 BTC) an amount equal to oracle price
   - If oracle price is 50,000 USD/BTC:
     - Carl initially: 100,000 USD - must pay 50,000 USD = 50,000 USD remaining
     - Dave initially: 50,000 USD + receives 50,000 USD = 100,000 USD

---

## Test Case 2: Final Settlement Deleveraging - Negative TNC Accounts

### Input
- **Initial Subaccounts:**
  - `Carl_Num0`: 1 BTC Short position with 49,999 USD collateral (negative TNC)
  - `Dave_Num0`: 1 BTC Long position with 50,001 USD collateral
- **ClobPair:** BTC-USD with initial status `STATUS_ACTIVE`
- **Proposal:** Transition ClobPair to `STATUS_FINAL_SETTLEMENT`

### Output
- **Subaccounts after final settlement:**
  - `Carl_Num0`: Empty (nothing left due to negative TNC)
  - `Dave_Num0`: Has 100,000 USD
- **ClobPair status:** `STATUS_FINAL_SETTLEMENT`

### Why It Runs This Way?

1. **Negative TNC:** Carl has negative TNC (49,999 USD < 50,000 USD needed to cover short position), meaning this account is undercollateralized.
2. **Deleveraging at Bankruptcy Price:** When an account has negative TNC, deleveraging is performed at "bankruptcy price" - the price at which the account has nothing left after settlement.
3. **Result:**
   - Carl doesn't have enough funds to fully settle, so all of Carl's collateral (49,999 USD) is transferred to Dave
   - Dave receives 49,999 USD from Carl + 50,001 USD initially = 100,000 USD
   - Carl loses everything and the account becomes empty

---

## Test Case 3: Cancel Open Stateful Orders

### Input
- **Subaccounts:**
  - `Alice_Num0`: 10,000 USD
  - `Bob_Num0`: 10,000 USD
- **Preexisting Stateful Orders:**
  - Long-term order from Alice: Buy 5 units at price 5, GoodTilBlockTime = 5
  - Long-term order from Bob: Sell 10 units at price 10, GoodTilBlockTime = 10, PostOnly
  - Conditional order from Alice: Buy 1 BTC at price 50,000, GoodTilBlockTime = 10, StopLoss trigger = 50,001
- **Proposal:** Transition ClobPair to `STATUS_FINAL_SETTLEMENT`

### Output
- **Stateful Orders:** All 3 orders removed from state
- **Indexer Events:** Events emitted for each removed order with reason `ORDER_REMOVAL_REASON_FINAL_SETTLEMENT`
- **ClobPair status:** `STATUS_FINAL_SETTLEMENT`

### Why It Runs This Way?

1. **Stateful Orders:** These are orders that persist across multiple blocks (long-term and conditional orders), different from short-term orders that only exist for 1 block.
2. **Must Cancel When Wind Down:** When market is closed, all pending orders must be canceled as they can no longer be executed.
3. **Events:** Indexer needs to be notified about order cancellations to update off-chain state.
4. **Both Sides:** This test ensures orders from both sides (buy and sell) are canceled, including conditional orders.

---

## Test Case 4: Block New Order Placement

### Input
- **Subaccounts:**
  - `Alice_Num0`: 10,000 USD
  - `Bob_Num0`: 10,000 USD
  - `Carl_Num0`: 10,000 USD
- **Orders attempting to place (after wind down):**
  - Short-term orders: Buy 10, Sell 15, IOC order
  - Long-term orders: Buy 5, Sell 10
  - Conditional orders: Buy 1 BTC with stop loss
- **Proposal:** ClobPair has been transitioned to `STATUS_FINAL_SETTLEMENT`

### Output
- **All CheckTx responses:** FAIL with log containing "trading is disabled for clob pair"
- **No orders placed:** All orders are rejected

### Why It Runs This Way?

1. **Protect Final Settlement State:** Once market is in final settlement state, new orders are not allowed because:
   - Market is in the process of closing
   - Only final settlement deleveraging is allowed
   - No new trading activity can occur

2. **Validation at CheckTx:** Validation is performed at `CheckTx` to reject orders early, before they enter the mempool.

3. **All Order Types:** This test ensures short-term, long-term, and conditional orders are all blocked, regardless of order type.

---

## Flow Summary

### Wind Down Market Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. INITIALIZE GENESIS STATE                                  │
│    - Create ClobPair with ACTIVE status                     │
│    - Create subaccounts with positions                      │
│    - Place stateful orders (if any)                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                                │
│    - Create MsgUpdateClobPair with FINAL_SETTLEMENT status  │
│    - Submit proposal through governance module               │
│    - Validators vote (in test: all vote YES)                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTED                                         │
│    - Proposal status: PROPOSAL_STATUS_PASSED                │
│    - ClobPair status transitions to FINAL_SETTLEMENT        │
│    - Indexer events emitted                                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. CANCEL STATEFUL ORDERS                                    │
│    - All long-term orders removed                            │
│    - All conditional orders removed                          │
│    - Indexer events emitted for each removed order           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. FINAL SETTLEMENT DELEVERAGING                             │
│    - Liquidations daemon provides SubaccountOpenPositionInfo │
│    - System identifies accounts needing deleverage           │
│    - Perform deleveraging:                                   │
│      * Non-negative TNC → deleverage at oracle price         │
│      * Negative TNC → deleverage at bankruptcy price        │
│    - Update subaccount balances                              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 6. BLOCK NEW ORDERS                                          │
│    - CheckTx validation rejects all new orders               │
│    - Log: "trading is disabled for clob pair"                │
│    - Applies to all order types                              │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **ClobPair Status Transition:**
   ```
   STATUS_ACTIVE → STATUS_FINAL_SETTLEMENT
   ```

2. **Subaccount State Changes:**
   - Positions are closed (deleveraged)
   - Balances updated based on oracle/bankruptcy price
   - Negative TNC accounts may be wiped clean

3. **Order State:**
   - Stateful orders: Removed from state
   - New orders: Rejected at CheckTx

### Key Points

1. **Timing:** Final settlement deleveraging only occurs after proposal execution and liquidations daemon has provided position information.

2. **Price Determination:**
   - **Oracle Price:** Used for accounts with positive TNC
   - **Bankruptcy Price:** Used for accounts with negative TNC (price at which account has nothing left)

3. **Event Emission:** Indexer needs to be notified about:
   - ClobPair status update
   - Stateful order removals
   - To sync state off-chain

4. **Validation:** CheckTx validation ensures no new orders can be placed after market has been wound down.

### Design Rationale

1. **Safety:** Wind down market is a critical process that must ensure:
   - All positions are settled correctly
   - No new trading activity
   - State is updated consistently

2. **Fairness:** 
   - Non-negative TNC accounts are settled at oracle price (fair value)
   - Negative TNC accounts are settled at bankruptcy price (lose everything)

3. **Transparency:** Indexer events ensure off-chain systems can track all changes.

