# Test Documentation: Reduce-Only Orders E2E Tests

## Overview

This test file verifies **Reduce-Only Order** functionality in the CLOB module. Reduce-only orders are orders that can only reduce an existing position, not open a new position or increase an existing position. The test ensures that:
1. Reduce-only orders can partially match to reduce position
2. Reduce-only orders cannot increase position size
3. Reduce-only orders work with IOC orders
4. Reduce-only orders work across multiple blocks

---

## Test Function: TestReduceOnlyOrders

### Test Case 1: Success - IOC Reduce-Only Order Partially Matches, Maker Fully Filled (Same Block)

### Input
- **Subaccounts:**
  - Carl: 100,000 USD
  - Alice: 1 BTC Long, 500,000 USD
- **Orders (Block 1):**
  - Carl: Buy 10 at price 500,000 (maker)
  - Alice: Sell 15 at price 500,000, IOC, Reduce-Only (taker)
- **Match:** Alice order matches with Carl order

### Output
- **Carl Order:** Fully filled (10)
- **Alice Order:** Partially filled (10), remaining cancelled
- **Carl Position:** 10 (new position opened)
- **Alice Position:** 0.999999 (reduced from 1 BTC)

### Why It Runs This Way?

1. **Reduce-Only:** Alice's order can only reduce her 1 BTC long position.
2. **Partial Match:** Order matches 10 units, reducing position from 1 BTC to 0.999999 BTC.
3. **IOC Behavior:** Remaining 5 units cancelled (IOC rule).
4. **Maker Filled:** Carl's maker order fully filled.

---

### Test Case 2: Success - IOC Reduce-Only Order Partially Matches (Second Block)

### Input
- **Subaccounts:**
  - Carl: 100,000 USD
  - Alice: 1 BTC Long, 500,000 USD
- **Orders:**
  - Block 1: Carl buys 10 at price 500,000
  - Block 2: Alice sells 15 at price 500,000, IOC, Reduce-Only

### Output
- **Carl Order:** Fully filled (10)
- **Alice Order:** Partially filled (10), remaining cancelled
- **Alice Position:** Reduced from 1 BTC to 0.999999 BTC

### Why It Runs This Way?

1. **Cross-Block Matching:** Reduce-only order can match with orders from previous blocks.
2. **Same Behavior:** Reduce-only logic works same across blocks.
3. **Position Reduction:** Position reduced by matched amount.

---

### Test Case 3: Success - IOC Reduce-Only Order Partially Matches, Maker Partially Filled

### Input
- **Subaccounts:**
  - Carl: 100,000 USD
  - Alice: 1 BTC Long, 500,000 USD
- **Orders:**
  - Block 1: Carl buys 80 at price 500,000
  - Block 2: Alice sells 15 at price 500,000, IOC, Reduce-Only

### Output
- **Carl Order:** Partially filled (15), 65 remaining on book
- **Alice Order:** Partially filled (15), remaining cancelled
- **Alice Position:** Reduced from 1 BTC to 0.9999985 BTC

### Why It Runs This Way?

1. **Partial Fill Both:** Both orders partially filled.
2. **Maker Remains:** Carl's order remains on book with remaining size.
3. **Position Reduced:** Alice's position reduced by filled amount.

---

## Flow Summary

### Reduce-Only Order Logic

```
┌─────────────────────────────────────────────────────────────┐
│ 1. CHECK EXISTING POSITION                                   │
│    - Query subaccount's current position                     │
│    - Determine position side (long/short)                    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. VALIDATE ORDER DIRECTION                                  │
│    - Reduce-only buy: Only valid if short position exists    │
│    - Reduce-only sell: Only valid if long position exists    │
│    - Reject if no position or wrong direction               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. CALCULATE MAX FILL AMOUNT                                 │
│    - Max fill = min(order size, position size)               │
│    - Cannot fill more than existing position                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. EXECUTE MATCH                                            │
│    - Fill up to max fill amount                              │
│    - Reduce position by fill amount                          │
│    - Cancel remaining if IOC                                 │
└─────────────────────────────────────────────────────────────┘
```

### Key Points

1. **Reduce-Only Constraint:**
   - Can only reduce existing position
   - Cannot open new position
   - Cannot increase position size

2. **Position Direction:**
   - Reduce-only buy: Only works with short position
   - Reduce-only sell: Only works with long position
   - Must match position direction

3. **Fill Amount:**
   - Limited by existing position size
   - Cannot fill more than position
   - Partial fills allowed

4. **IOC Compatibility:**
   - Reduce-only works with IOC orders
   - Remaining size cancelled if not fully filled
   - Immediate execution or cancellation

5. **Cross-Block Matching:**
   - Can match with orders from previous blocks
   - Position checked at match time
   - Works same as same-block matching

### Design Rationale

1. **Risk Management:** Reduce-only orders help users close positions without accidentally opening new ones.

2. **Position Control:** Ensures users can only reduce risk, not increase it.

3. **Flexibility:** Works with various order types (IOC, regular, etc.).

4. **Safety:** Prevents accidental position increases.

