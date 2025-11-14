# Test Documentation: Batch Cancel E2E Tests

## Overview

This test file verifies **Batch Cancel Order** functionality in the CLOB module. Batch cancel allows users to cancel multiple orders in a single transaction, grouped by CLOB pair. The test ensures that:
1. Multiple orders can be cancelled in a single batch cancel transaction
2. Batch cancel works for unfilled, partially filled, and fully filled orders
3. Batch cancel respects GoodTilBlock constraints
4. Batch cancel works across multiple blocks

---

## Test Function: TestBatchCancelSingleCancelFunctionality

### Test Case 1: Success - Cancel Unfilled Short-Term Order

### Input
- **Block 1:**
  - Place order: Alice buys 5 at price 10, GTB 5
- **Block 1:**
  - Batch cancel: Cancel Alice's order with client ID 0 on CLOB pair 0, GTB 5

### Output
- **Order:** Removed from memclob
- **Cancel Expiration:** Set to block 5
- **Fill Amount:** 0 (unfilled)

### Why It Runs This Way?

1. **Batch Cancel:** Single order cancelled via batch cancel message.
2. **Unfilled Order:** Order was never matched, so fill amount is 0.
3. **Removal:** Order removed from memclob immediately.

---

### Test Case 2: Success - Batch Cancel Partially Filled Order (Same Block)

### Input
- **Block 1:**
  - Place order: Alice buys 5 at price 10, GTB 5
  - Place order: Bob sells 4 at price 10, GTB 20 (matches with Alice)
  - Batch cancel: Cancel Alice's order

### Output
- **Order:** Removed from memclob
- **Fill Amount:** 4 (40% filled)
- **Cancel Expiration:** Set to block 5

### Why It Runs This Way?

1. **Partial Fill:** Order was partially filled (4 out of 5) before cancellation.
2. **Same Block:** Cancellation happens in same block as partial fill.
3. **Remaining Cancelled:** Remaining unfilled portion (1) is cancelled.

---

### Test Case 3: Success - Cancel Partially Filled Order (Next Block)

### Input
- **Block 1:**
  - Place order: Alice buys 5 at price 10, GTB 5
  - Place order: Bob sells 4 at price 10, GTB 20 (matches)
- **Block 2:**
  - Batch cancel: Cancel Alice's order

### Output
- **Order:** Removed from memclob
- **Fill Amount:** 4 (40% filled)
- **Cancel Expiration:** Set to block 5

### Why It Runs This Way?

1. **Cross-Block:** Cancellation happens in block after partial fill.
2. **Same Behavior:** Works same as same-block cancellation.
3. **Fill Preserved:** Fill amount from previous block is preserved.

---

### Test Case 4: Success - Cancel Fully Filled Order

### Input
- **Block 1:**
  - Place order: Alice buys 5 at price 10, GTB 5
  - Place order: Bob sells 5 at price 10, GTB 20 (fully matches)
- **Block 2:**
  - Batch cancel: Cancel Alice's order

### Output
- **Order:** Removed from memclob (already removed after fill)
- **Fill Amount:** 5 (100% filled)
- **Cancel Expiration:** Set to block 5

### Why It Runs This Way?

1. **Fully Filled:** Order was fully filled in block 1.
2. **Already Removed:** Order already removed from memclob after fill.
3. **Cancel Succeeds:** Cancellation succeeds even though order already filled.

---

### Test Case 5: Failure - Cancel with GTB < Order GTB Does Not Remove Order

### Input
- **Block 1:**
  - Place order: Alice buys 5 at price 10, GTB 20
- **Block 2:**
  - Batch cancel: Cancel with GTB 5 (less than order's GTB 20)

### Output
- **Order:** Still in memclob (not removed)
- **Cancel Expiration:** Set to block 5
- **Order Expiration:** Still block 20

### Why It Runs This Way?

1. **GTB Constraint:** Cancel GTB (5) < Order GTB (20).
2. **Order Not Removed:** Order remains because cancel expires before order.
3. **Cancel Expires First:** Cancel expires at block 5, order at block 20.

---

## Flow Summary

### Batch Cancel Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. CREATE BATCH CANCEL                                       │
│    - Specify subaccount ID                                   │
│    - Group orders by CLOB pair                                │
│    - List client IDs to cancel                                │
│    - Set GoodTilBlock                                         │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CHECKTX VALIDATION                                        │
│    - Validate message format                                 │
│    - Check subaccount exists                                  │
│    - Verify CLOB pairs exist                                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. DELIVERTX EXECUTION                                       │
│    - For each CLOB pair:                                      │
│      * For each client ID:                                    │
│        - Find order by subaccount + client ID + CLOB pair    │
│        - Cancel order (if exists and not expired)            │
│    - Set cancel expiration                                    │
│    - Remove orders from memclob                               │
└─────────────────────────────────────────────────────────────┘
```

### Key Points

1. **Batch Structure:**
   - Orders grouped by CLOB pair
   - Multiple client IDs per CLOB pair
   - Single transaction cancels multiple orders

2. **Order Matching:**
   - Orders matched by: SubaccountId + ClientId + ClobPairId
   - Must match exactly to cancel
   - Non-existent orders ignored

3. **GTB Constraints:**
   - Cancel GTB must be >= order GTB to remove order
   - If cancel GTB < order GTB, order remains
   - Cancel expires at its GTB

4. **Fill Handling:**
   - Unfilled orders: Cancelled completely
   - Partially filled: Remaining portion cancelled
   - Fully filled: Cancel succeeds but order already removed

5. **Cross-Block:**
   - Can cancel orders from previous blocks
   - Fill amounts preserved
   - Works same as same-block cancellation

### Design Rationale

1. **Efficiency:** Batch cancel allows cancelling multiple orders in one transaction.

2. **Flexibility:** Grouping by CLOB pair allows selective cancellation.

3. **Safety:** GTB constraints prevent premature order removal.

4. **Consistency:** Works consistently across blocks.

