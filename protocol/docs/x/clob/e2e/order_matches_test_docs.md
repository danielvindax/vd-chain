# Test Documentation: Order Matches E2E Tests

## Overview

This test file verifies **Order Matching Validation** in the CLOB module. The test ensures that order matching operations are correctly validated during DeliverTx, including:
1. IOC orders cannot be matched twice
2. Partially filled conditional IOC orders cannot be matched again
3. IOC orders can match with multiple makers in single operation
4. IOC orders cannot be taker in multiple separate matches

---

## Test Function: TestDeliverTxMatchValidation

### Test Case 1: Failure - Partially Filled IOC Taker Order Cannot Be Matched Twice

### Input
- **Block 1:**
  - Place order: Bob buys 5 at price 40
  - Place order: Alice sells 10 at price 15, IOC
  - Match operation: Alice (IOC taker) matches with Bob (maker), fill 5
- **Block 2:**
  - Place order: Bob buys 5 at price 40
  - Place order: Alice sells 10 at price 15, IOC (same order)
  - Match operation: Attempt to match Alice IOC order again

### Output
- **Block 1 DeliverTx:** SUCCESS
- **Block 2 DeliverTx:** FAIL with error "IOC order is already filled, remaining size is cancelled."

### Why It Runs This Way?

1. **IOC Rule:** IOC orders must be fully filled immediately or cancelled.
2. **Partial Fill:** Order was partially filled in block 1 (5 out of 10).
3. **Remaining Cancelled:** Remaining size (5) was cancelled after partial fill.
4. **Cannot Reuse:** Cannot match the same IOC order again in later block.

---

### Test Case 2: Failure - Cannot Match Partially Filled Conditional IOC Order

### Input
- **Block 1:**
  - Place conditional IOC order: Alice buys 1 BTC at 50,000, TakeProfit trigger at 49,999
  - Place long-term order: Dave sells 0.25 BTC at 50,000
- **Block 2:**
  - Conditional order triggers and partially matches (0.25 BTC filled)
- **Block 3:**
  - Place order: Dave sells 1 BTC at 50,000
  - Match operation: Attempt to match conditional IOC order again

### Output
- **Block 2 DeliverTx:** SUCCESS (partial fill)
- **Block 3 DeliverTx:** FAIL with error `ErrStatefulOrderDoesNotExist`

### Why It Runs This Way?

1. **Conditional IOC:** Conditional IOC order partially filled in block 2.
2. **Order Removed:** After partial fill, conditional IOC order is removed from state.
3. **Cannot Match:** Cannot match order that no longer exists in state.

---

### Test Case 3: Success - IOC Order Matches with Multiple Makers in Single Operation

### Input
- **Orders:**
  - Bob buys 5 at price 40 (maker 1)
  - Bob buys 5 at price 40 (maker 2)
  - Alice sells 10 at price 15, IOC (taker)
- **Match Operation:**
  - Alice IOC order matches with both Bob orders
  - Fill 5 from maker 1, fill 5 from maker 2

### Output
- **DeliverTx:** SUCCESS
- **Result:** Alice order fully filled (10), both Bob orders fully filled (5 each)

### Why It Runs This Way?

1. **Multiple Makers:** IOC order can match with multiple maker orders in single operation.
2. **Full Fill:** Order is fully filled by combining fills from multiple makers.
3. **Single Operation:** All matches happen in one match operation.

---

### Test Case 4: Failure - IOC Order Cannot Be Taker in Multiple Matches

### Input
- **Orders:**
  - Bob buys 5 at price 40 (maker 1)
  - Alice sells 10 at price 15, IOC (taker)
  - Bob buys 5 at price 40 (maker 2)
- **Match Operations:**
  - Match 1: Alice IOC with maker 1 (fill 5)
  - Match 2: Alice IOC with maker 2 (fill 5)

### Output
- **DeliverTx:** FAIL with error "IOC order is already filled, remaining size is cancelled."

### Why It Runs This Way?

1. **Multiple Matches:** IOC order cannot be taker in multiple separate match operations.
2. **First Match:** After first match, order is considered filled/cancelled.
3. **Second Match Fails:** Cannot use same IOC order in second match operation.

---

## Flow Summary

### IOC Order Matching Rules

1. **Single Match Operation:**
   - IOC order can match with multiple makers in one operation
   - All fills happen atomically
   - Order fully filled or cancelled

2. **No Multiple Matches:**
   - IOC order cannot be taker in multiple separate operations
   - After first match, order is filled/cancelled
   - Subsequent matches fail

3. **Partial Fill Handling:**
   - If IOC order partially fills, remaining size is cancelled
   - Order removed from state after partial fill
   - Cannot match again in later blocks

### Key Points

1. **IOC Order Behavior:**
   - Must fill completely immediately or be cancelled
   - Cannot persist across blocks if partially filled
   - Remaining size cancelled after partial fill

2. **Match Operation:**
   - Single match operation can include multiple maker fills
   - All fills happen atomically
   - Order state updated after operation

3. **State Management:**
   - Partially filled IOC orders removed from state
   - Cannot query or match removed orders
   - State consistency maintained

4. **Validation:**
   - DeliverTx validates match operations
   - Checks order existence and state
   - Rejects invalid matches

### Design Rationale

1. **Immediate Execution:** IOC orders ensure immediate execution or cancellation.

2. **State Consistency:** Prevents matching orders that no longer exist.

3. **Atomic Operations:** Single match operation ensures atomic fills.

4. **Safety:** Validation prevents invalid match operations.

