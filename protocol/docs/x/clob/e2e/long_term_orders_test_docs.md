# Test Documentation: Long-Term Orders E2E Tests

## Overview

This test file verifies **Long-Term Order** functionality in the CLOB module. Long-term orders are orders that persist across multiple blocks until they are filled, cancelled, or expire. The test ensures that:
1. Long-term orders can be placed and persist across blocks
2. Long-term orders can be cancelled
3. Order cancellation and placement in same block is handled correctly
4. Fully filled orders cannot be cancelled

---

## Test Function: TestPlaceOrder_StatefulCancelFollowedByPlaceInSameBlockErrorsInCheckTx

### Test Case: Failure - Cancel and Place Same Order in Same Block

### Input
- **Block 2:**
  - Place long-term order: Alice buys 5 at price 10
- **Block 3:**
  - Cancel the long-term order
  - Attempt to place the same order again

### Output
- **Cancel CheckTx:** SUCCESS
- **Place CheckTx:** FAIL with error "An uncommitted stateful order cancellation with this OrderId already exists"
- **Final State:** Order cancelled, new order not placed

### Why It Runs This Way?

1. **Uncommitted Cancellation:** When an order is cancelled in the same block, the cancellation is uncommitted.
2. **Conflict Detection:** System detects conflict between cancellation and placement of same order.
3. **Early Rejection:** CheckTx rejects the placement to prevent invalid state.

---

## Test Function: TestCancelFullyFilledStatefulOrderInSameBlockItIsFilled

### Test Case: Failure - Cancel Fully Filled Order

### Input
- **Block 2:**
  - Place long-term order: Alice buys 5 at price 10
- **Block 3:**
  - Place matching order: Bob sells 5 at price 10 (fully fills Alice's order)
  - Attempt to cancel Alice's order

### Output
- **Match CheckTx:** SUCCESS
- **Cancel CheckTx:** SUCCESS
- **DeliverTx:** Cancel transaction FAILS with error `ErrStatefulOrderCancellationFailedForAlreadyRemovedOrder`
- **Final State:** Order fully filled, cancellation fails

### Why It Runs This Way?

1. **Order Filled First:** Matching order fills the long-term order before cancellation executes.
2. **Cancellation Fails:** Cannot cancel an order that has already been removed (filled).
3. **Transaction Ordering:** DeliverTx processes transactions in order, so fill happens before cancellation.

---

## Test Function: TestCancelStatefulOrder

### Test Case 1: Success - Cancel Order in Same Block

### Input
- **Block 2:**
  - Place long-term order
  - Cancel the same order

### Output
- **Both CheckTx:** SUCCESS
- **Final State:** Order does not exist in state

### Why It Runs This Way?

1. **Same Block Cancellation:** Order can be cancelled in the same block it's placed.
2. **State Cleanup:** Order is removed from state immediately.

---

### Test Case 2: Success - Cancel Order in Future Block

### Input
- **Block 2:**
  - Place long-term order
- **Block 3:**
  - Cancel the order

### Output
- **Place CheckTx:** SUCCESS
- **Cancel CheckTx:** SUCCESS
- **Final State:** Order removed from state

### Why It Runs This Way?

1. **Persistent Orders:** Long-term orders persist across blocks.
2. **Future Cancellation:** Orders can be cancelled in any future block before expiration.

---

## Flow Summary

### Long-Term Order Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│ 1. PLACE ORDER                                              │
│    - Order placed on order book                              │
│    - Order persists in state                                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. ORDER PERSISTS                                            │
│    - Order remains on book across blocks                     │
│    - Can be matched at any time                              │
│    - Can be cancelled at any time                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. ORDER TERMINATION                                         │
│    - Option A: Fully filled → Removed                        │
│    - Option B: Cancelled → Removed                           │
│    - Option C: Expired → Removed                             │
└─────────────────────────────────────────────────────────────┘
```

### Cancel Order Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. SUBMIT CANCELLATION                                       │
│    - Create MsgCancelOrderStateful                            │
│    - Specify order ID and good til block time                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CHECKTX VALIDATION                                        │
│    - Verify order exists in state                            │
│    - Check cancellation parameters                           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. DELIVERTX EXECUTION                                       │
│    - Verify order still exists (not filled)                  │
│    - Remove order from state                                 │
│    - Emit order removal events                               │
└─────────────────────────────────────────────────────────────┘
```

### Key Points

1. **Long-Term Orders:**
   - Persist across multiple blocks
   - GoodTilBlockTime specifies expiration time
   - Can be matched, cancelled, or expire

2. **Cancellation:**
   - Can cancel in same block or future blocks
   - Cannot cancel if order already filled
   - Cannot cancel if order doesn't exist

3. **Conflict Detection:**
   - System detects conflicts between cancellation and placement
   - CheckTx rejects conflicting operations early
   - Prevents invalid state

4. **State Management:**
   - Orders tracked in keeper state
   - Cancellation removes order from state
   - Fill removes order from state

5. **Transaction Ordering:**
   - DeliverTx processes transactions in order
   - First transaction to modify order wins
   - Later conflicting transactions fail

### Design Rationale

1. **Flexibility:** Long-term orders allow users to set orders that persist until filled or cancelled.

2. **Safety:** Conflict detection prevents invalid state from cancellation/placement conflicts.

3. **Efficiency:** Early rejection at CheckTx prevents wasted computation.

4. **Consistency:** Transaction ordering ensures deterministic state updates.

