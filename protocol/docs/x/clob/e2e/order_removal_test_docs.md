# Test Documentation: Order Removal E2E Tests

## Overview

This test file verifies **Order Removal** functionality in the CLOB module. Orders can be removed from the order book for various reasons. The test ensures that:
1. Conditional orders are removed when they cross maker orders (PostOnly violation)
2. Conditional IOC orders are removed if not fully filled
3. Self-trading removes maker orders
4. Fully filled orders are removed
5. Under-collateralized orders are removed

---

## Test Function: TestConditionalOrderRemoval

### Test Case 1: Conditional PostOnly Order Crosses Maker - Removed

### Input
- **Subaccounts:**
  - Alice: 10,000 USD
  - Bob: 10,000 USD
- **Orders:**
  - Long-term order: Alice buys 5 at price 10 (maker)
  - Conditional order: Bob sells 10 at price 10, PostOnly, StopLoss trigger at 15
- **Price Update:** Price rises to 14.9 (triggers conditional order)

### Output
- **Alice Order:** Not removed (maker order)
- **Bob Order:** Removed (PostOnly violation - crosses maker)

### Why It Runs This Way?

1. **PostOnly Violation:** When conditional order triggers, it would cross existing maker order.
2. **PostOnly Rule:** PostOnly orders cannot cross existing orders, must be maker.
3. **Removal:** Conditional order is removed instead of crossing.

---

### Test Case 2: Conditional IOC Order Not Fully Filled - Removed

### Input
- **Subaccounts:**
  - Carl: 10,000 USD
  - Dave: 10,000 USD
- **Orders:**
  - Long-term order: Dave sells 0.25 BTC at 50,000
  - Conditional order: Carl buys 0.5 BTC at 50,000, IOC, StopLoss trigger at 50,003
- **Price Update:** Price rises to 50,004 (triggers conditional order)

### Output
- **Dave Order:** Removed (fully filled)
- **Carl Order:** Removed (IOC not fully filled)

### Why It Runs This Way?

1. **IOC Rule:** Immediate-Or-Cancel orders must be fully filled immediately or cancelled.
2. **Partial Fill:** Only 0.25 BTC available, but order wants 0.5 BTC.
3. **Removal:** IOC order is removed because it cannot be fully filled.

---

### Test Case 3: Conditional Self Trade - Removes Maker Order

### Input
- **Subaccount:** Alice with 10,000 USD
- **Orders:**
  - Long-term order: Alice buys 5 at price 10
  - Conditional order: Alice sells 20 at price 10, StopLoss trigger at 15
- **Price Update:** Price rises to 14.9 (triggers conditional order)

### Output
- **Long-term Order:** Removed (self-trade removes maker)
- **Conditional Order:** Not removed (taker in self-trade)

### Why It Runs This Way?

1. **Self-Trade:** Same subaccount has both maker and taker orders.
2. **Maker Removal:** Self-trading removes the maker order to prevent abuse.
3. **Taker Kept:** Taker order (conditional) is kept.

---

### Test Case 4: Fully Filled Maker Orders - Removed

### Input
- **Subaccounts:**
  - Alice: 10,000 USD
  - Bob: 10,000 USD
- **Orders:**
  - Long-term order: Alice buys 5 at price 10
  - Conditional order: Bob sells 50 at price 10, StopLoss trigger at 15
- **Price Update:** Price rises to 14.9 (triggers conditional order)

### Output
- **Alice Order:** Removed (fully filled by conditional order)
- **Bob Order:** Not removed (partially filled, 45 remaining)

### Why It Runs This Way?

1. **Full Fill:** Conditional order fully fills maker order (5 units).
2. **Maker Removal:** Fully filled maker order is removed.
3. **Taker Partial:** Conditional order partially filled, remains on book.

---

### Test Case 5: Under-Collateralized Conditional Taker - Removed

### Input
- **Subaccounts:**
  - Carl: 100,000 USD
  - Dave: 10,000 USD
- **Orders:**
  - Long-term order: Carl buys 1 BTC at 50,000
  - Conditional order: Dave sells 1 BTC at 50,000, StopLoss trigger at 50,003
- **Withdrawal:** Dave withdraws 10,000 USD (becomes under-collateralized)
- **Price Update:** Price rises to 50,002.5 (triggers conditional order)

### Output
- **Carl Order:** Not removed
- **Dave Order:** Removed (fails collateralization check during matching)

### Why It Runs This Way?

1. **Collateralization Check:** When conditional order triggers and tries to match, system checks collateral.
2. **Insufficient Collateral:** Dave doesn't have enough collateral after withdrawal.
3. **Removal:** Order is removed instead of executing trade.

---

## Flow Summary

### Order Removal Reasons

1. **PostOnly Violation:**
   - Order would cross existing maker
   - PostOnly orders must be maker
   - Order removed instead of crossing

2. **IOC Not Fully Filled:**
   - IOC order cannot be fully filled immediately
   - IOC orders must fill completely or be cancelled
   - Order removed

3. **Self-Trade:**
   - Same subaccount has maker and taker orders
   - Maker order removed to prevent abuse
   - Taker order kept

4. **Fully Filled:**
   - Order completely filled by matching
   - Fully filled orders removed from book
   - State updated to reflect fill

5. **Under-Collateralized:**
   - Order fails collateralization check
   - Insufficient margin for position
   - Order removed before execution

### Key Points

1. **Removal Timing:**
   - Orders removed during DeliverTx
   - Removal happens before state update
   - Events emitted for removed orders

2. **Removal Reasons:**
   - Tracked in order removal events
   - Different reasons for different scenarios
   - Used for off-chain tracking

3. **State Consistency:**
   - Removed orders not in state
   - Cannot query removed orders
   - Fill amounts tracked before removal

4. **Event Emission:**
   - Order removal events emitted
   - Include removal reason
   - Used by indexer for off-chain sync

5. **Collateralization:**
   - Checked when order tries to match
   - Must have sufficient margin
   - Under-collateralized orders removed

### Design Rationale

1. **Order Book Integrity:** Removal prevents invalid orders from staying on book.

2. **Risk Management:** Under-collateralized orders removed to prevent bad trades.

3. **Fairness:** Self-trade removal prevents manipulation.

4. **Efficiency:** Removal of invalid orders keeps book clean.

