# Test Documentation: Equity Tier Limit E2E Tests

## Overview

This test file verifies **Equity Tier Limit** functionality in the CLOB module. Equity tier limits restrict the number of stateful orders (long-term and conditional orders) a subaccount can have open based on their Total Net Collateral (TNC). The test ensures that:
1. Subaccounts with lower TNC have fewer allowed stateful orders
2. Subaccounts with higher TNC have more allowed stateful orders
3. Order cancellation can free up slots for new orders
4. Both long-term and conditional orders count toward limits

---

## Test Function: TestPlaceOrder_EquityTierLimit

### Test Case 1: Failure - Long-Term Order Exceeds Max Open Stateful Orders

### Input
- **Subaccount:** Alice with TNC < $5,000
- **Existing Orders:**
  - 1 conditional order (StopLoss)
- **New Order:** Long-term order
- **Equity Tier Config:**
  - Tier 0: $0 TNC → 0 orders
  - Tier 1: $5,000 TNC → 1 order
  - Tier 2: $70,000 TNC → 100 orders

### Output
- **CheckTx:** FAIL (after advancing block)
- **Error:** Would exceed max open stateful orders

### Why It Runs This Way?

1. **Tier Limit:** Alice is in tier 1 (TNC < $5,000), limit = 1 order.
2. **Already Has 1:** Already has 1 conditional order.
3. **Exceeds Limit:** New long-term order would exceed limit of 1.
4. **Rejection:** Order rejected to prevent exceeding limit.

---

### Test Case 2: Failure - Conditional Order Exceeds Max Open Stateful Orders

### Input
- **Subaccount:** Alice with TNC < $5,000
- **Existing Orders:**
  - 1 long-term order
- **New Order:** Conditional order (StopLoss)
- **Equity Tier Config:** Same as Test Case 1

### Output
- **CheckTx:** FAIL (after advancing block)
- **Error:** Would exceed max open stateful orders

### Why It Runs This Way?

1. **Same Limit:** Conditional orders count toward same limit as long-term orders.
2. **Already Has 1:** Already has 1 long-term order.
3. **Exceeds Limit:** New conditional order would exceed limit of 1.

---

### Test Case 3: Success - Order Cancellation Frees Up Slot

### Input
- **Subaccount:** Alice with TNC < $5,000
- **Existing Orders:**
  - 1 conditional order (StopLoss)
- **Cancellation:** Cancel the conditional order
- **New Order:** Long-term order (same block)
- **Equity Tier Config:** Same as Test Case 1

### Output
- **Cancellation:** SUCCESS
- **New Order:** SUCCESS
- **Final State:** 1 long-term order (conditional cancelled)

### Why It Runs This Way?

1. **Cancellation First:** Conditional order cancelled, freeing up slot.
2. **Slot Available:** After cancellation, slot available for new order.
3. **Same Block:** Cancellation and placement in same block works.
4. **Limit Respected:** Final state has 1 order, within limit.

---

### Test Case 4: Failure - Conditional Order Would Exceed Limit (Untriggered)

### Input
- **Subaccount:** Alice with TNC < $5,000
- **Existing Orders:**
  - 1 long-term order
- **New Order:** Conditional order (TakeProfit, untriggered)
- **Equity Tier Config:** Same as Test Case 1

### Output
- **CheckTx:** FAIL (after advancing block)
- **Error:** Would exceed max open stateful orders

### Why It Runs This Way?

1. **Untriggered Counts:** Untriggered conditional orders count toward limit.
2. **Same Limit:** Both triggered and untriggered conditional orders count.
3. **Exceeds Limit:** New conditional order would exceed limit.

---

## Flow Summary

### Equity Tier Limit Check

```
┌─────────────────────────────────────────────────────────────┐
│ 1. CALCULATE SUBACCOUNT TNC                                 │
│    - Get subaccount's Total Net Collateral                  │
│    - Include all positions and assets                        │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. DETERMINE EQUITY TIER                                     │
│    - Find tier based on TNC amount                           │
│    - Tier 0: $0 TNC → 0 orders                              │
│    - Tier 1: $5,000 TNC → 1 order                           │
│    - Tier 2: $70,000 TNC → 100 orders                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. COUNT EXISTING STATEFUL ORDERS                            │
│    - Count long-term orders                                  │
│    - Count conditional orders (triggered and untriggered)   │
│    - Count orders in same block (uncommitted)                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. VALIDATE LIMIT                                            │
│    - Check if new order would exceed limit                   │
│    - Consider cancellations in same block                    │
│    - Reject if would exceed limit                            │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Equity Tiers:**
   ```
   Tier 0: $0 TNC → 0 orders
   Tier 1: $5,000 TNC → 1 order
   Tier 2: $70,000 TNC → 100 orders
   ```

2. **Order Counting:**
   ```
   Long-term orders: Count toward limit
   Conditional orders (triggered): Count toward limit
   Conditional orders (untriggered): Count toward limit
   ```

### Key Points

1. **TNC-Based Limits:**
   - Limits based on Total Net Collateral
   - Higher TNC = more allowed orders
   - Protects system from order book spam

2. **Stateful Orders:**
   - Long-term orders count toward limit
   - Conditional orders count toward limit
   - Short-term orders don't count (expire same block)

3. **Same Block Logic:**
   - Cancellations free up slots immediately
   - Can cancel and place in same block
   - Uncommitted orders count toward limit

4. **Untriggered Conditionals:**
   - Untriggered conditional orders count toward limit
   - Must have slot available when placing
   - Triggering doesn't change count (same order)

5. **Validation Timing:**
   - Checked when placing order
   - After block advancement (for committed orders)
   - Considers same-block cancellations

### Design Rationale

1. **Resource Management:** Limits prevent order book spam from low-collateral accounts.

2. **Fairness:** Higher collateral accounts get more order slots.

3. **Flexibility:** Cancellations allow users to manage their order slots.

4. **Safety:** Prevents system overload from excessive stateful orders.

