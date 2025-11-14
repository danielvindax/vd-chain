# Test Documentation: Conditional Orders E2E Tests

## Overview

This test file verifies **Conditional Order** functionality in the CLOB module. Conditional orders are orders that are triggered when certain price conditions are met. The test ensures that:
1. Conditional orders are placed but not triggered when conditions aren't met
2. Conditional orders are triggered when price conditions are met
3. Different trigger types (TakeProfit, StopLoss) work correctly
4. Conditional orders can match with existing orders when triggered

---

## Test Function: TestConditionalOrder

### Test Case 1: TakeProfit/Buy - Not Triggered (No Price Update)

### Input
- **Subaccount:** Alice with 100,000 USD
- **Order:** Conditional Buy 1 BTC at price 50,000, TakeProfit trigger at 49,999
- **Price Updates:** None

### Output
- **Order State:** Exists in state, not triggered
- **Triggered State:** false after all blocks

### Why It Runs This Way?

1. **TakeProfit Logic:** TakeProfit/Buy triggers when price goes below trigger price.
2. **No Price Update:** Without price update, trigger condition never met.
3. **Order Persists:** Order remains in state waiting for trigger condition.

---

### Test Case 2: StopLoss/Buy - Not Triggered (No Price Update)

### Input
- **Subaccount:** Alice with 100,000 USD
- **Order:** Conditional Buy 1 BTC at price 50,000, StopLoss trigger at 50,001
- **Price Updates:** None

### Output
- **Order State:** Exists in state, not triggered
- **Triggered State:** false after all blocks

### Why It Runs This Way?

1. **StopLoss Logic:** StopLoss/Buy triggers when price goes above trigger price.
2. **No Price Update:** Without price update, trigger condition never met.
3. **Order Persists:** Order remains in state.

---

### Test Case 3: TakeProfit/Buy - Triggered by Price Update

### Input
- **Subaccount:** Alice with 100,000 USD
- **Order:** Conditional Buy 1 BTC at price 50,000, TakeProfit trigger at 49,999
- **Price Update:** Price drops to 49,997 (below trigger)

### Output
- **Order State:** Triggered and placed on order book
- **Triggered State:** true after price update
- **Order:** Can now match with existing orders

### Why It Runs This Way?

1. **Price Condition Met:** Price (49,997) < trigger (49,999), condition met.
2. **Order Triggered:** Conditional order becomes active order.
3. **Matching:** Triggered order can match with existing orders.

---

### Test Case 4: StopLoss/Buy - Triggered by Price Update

### Input
- **Subaccount:** Alice with 100,000 USD
- **Order:** Conditional Buy 1 BTC at price 50,000, StopLoss trigger at 50,001
- **Price Update:** Price rises to 50,003 (above trigger)

### Output
- **Order State:** Triggered and placed on order book
- **Triggered State:** true after price update

### Why It Runs This Way?

1. **Price Condition Met:** Price (50,003) > trigger (50,001), condition met.
2. **Order Triggered:** Conditional order becomes active order.

---

### Test Case 5: TakeProfit/Sell - Triggered and Matched

### Input
- **Subaccounts:**
  - Bob: 100,000 USD, 1 BTC long
  - Alice: 100,000 USD
- **Orders:**
  - Long-term order: Bob sells 1 BTC at 50,000
  - Conditional order: Alice buys 1 BTC at 50,000, TakeProfit trigger at 49,999
- **Price Update:** Price drops to 49,997

### Output
- **Conditional Order:** Triggered
- **Orders Matched:** Both orders fully matched
- **Positions:** Alice has 1 BTC long, Bob has 0 BTC

### Why It Runs This Way?

1. **Trigger Condition:** Price drops below trigger, order triggered.
2. **Immediate Matching:** Triggered order matches with existing order on book.
3. **Trade Execution:** Both orders filled, positions updated.

---

## Flow Summary

### Conditional Order Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│ 1. PLACE CONDITIONAL ORDER                                  │
│    - Order placed in conditional state                       │
│    - Trigger condition specified                             │
│    - Order not active on order book                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. WAIT FOR TRIGGER                                          │
│    - Monitor price updates                                   │
│    - Check trigger condition each block                      │
│    - Order remains in conditional state                      │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. TRIGGER CONDITION MET                                     │
│    - Price update meets trigger condition                    │
│    - Order transitions to active state                       │
│    - Order placed on order book                              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. ORDER ACTIVE                                             │
│    - Order can match with existing orders                    │
│    - Order behaves like regular order                        │
│    - Can be filled, cancelled, or expire                     │
└─────────────────────────────────────────────────────────────┘
```

### Trigger Types

1. **TakeProfit/Buy:**
   - Triggers when price < trigger price
   - Used to buy at lower price (profit taking)

2. **StopLoss/Buy:**
   - Triggers when price > trigger price
   - Used to buy at higher price (stop loss protection)

3. **TakeProfit/Sell:**
   - Triggers when price > trigger price
   - Used to sell at higher price (profit taking)

4. **StopLoss/Sell:**
   - Triggers when price < trigger price
   - Used to sell at lower price (stop loss protection)

### Key Points

1. **Conditional State:**
   - Orders start in conditional state
   - Not active on order book until triggered
   - Cannot match until triggered

2. **Trigger Conditions:**
   - Checked each block after price updates
   - Price must cross trigger price to activate
   - Once triggered, order becomes active

3. **Price Updates:**
   - Oracle price updates trigger condition checks
   - Price updates come from price feed
   - Multiple price updates can occur per block

4. **Order Matching:**
   - Once triggered, order behaves like regular order
   - Can match immediately if compatible order exists
   - Can remain on book if no match

5. **State Tracking:**
   - Triggered state tracked per order
   - Order state transitions: conditional → triggered → filled/cancelled

### Design Rationale

1. **Risk Management:** Conditional orders allow users to set automatic orders based on price movements.

2. **Flexibility:** Different trigger types support various trading strategies.

3. **Efficiency:** Orders only become active when conditions are met, reducing order book clutter.

4. **Safety:** Trigger conditions prevent accidental order execution.

