# Test Documentation: Rate Limiting E2E Tests

## Overview

This test file verifies **Rate Limiting** functionality in the CLOB module. Rate limits restrict the number of operations (orders, cancellations, leverage updates) a subaccount can perform within a specified number of blocks. The test ensures that:
1. Short-term orders are rate limited
2. Stateful orders are rate limited
3. Order cancellations are rate limited
4. Batch cancellations are rate limited
5. Leverage updates are rate limited
6. Rate limits apply per subaccount

---

## Test Function: TestRateLimitingOrders_RateLimitsAreEnforced

### Test Case 1: Failure - Short-Term Orders with Same Subaccount Exceed Limit

### Input
- **Rate Limit Config:**
  - MaxShortTermOrdersAndCancelsPerNBlocks: 1 order per 2 blocks
- **Block 2:**
  - Place order: Alice buys 5 at price 10, CLOB 0
- **Block 2:**
  - Attempt to place order: Alice buys 5 at price 10, CLOB 1

### Output
- **First Order CheckTx:** SUCCESS
- **Second Order CheckTx:** FAIL with error `ErrBlockRateLimitExceeded`
- **Error Message:** "exceeds configured block rate limit"

### Why It Runs This Way?

1. **Rate Limit:** 1 order per 2 blocks.
2. **First Order:** Consumes the limit for blocks 2-3.
3. **Second Order:** Attempts to place in same block, exceeds limit.
4. **Rejection:** CheckTx rejects second order immediately.

---

### Test Case 2: Failure - Short-Term Orders with Different Subaccounts Exceed Limit

### Input
- **Rate Limit Config:** Same as Test Case 1
- **Block 2:**
  - Place order: Alice_Num0 buys 5 at price 10
- **Block 2:**
  - Attempt to place order: Alice_Num1 buys 5 at price 10

### Output
- **First Order CheckTx:** SUCCESS
- **Second Order CheckTx:** FAIL with error `ErrBlockRateLimitExceeded`

### Why It Runs This Way?

1. **Per Subaccount:** Rate limits apply per subaccount.
2. **Different Subaccounts:** Alice_Num0 and Alice_Num1 are different subaccounts.
3. **Still Limited:** Even different subaccounts of same owner are rate limited.
4. **Owner-Based:** Rate limits may be based on owner address, not just subaccount.

---

### Test Case 3: Failure - Stateful Orders Exceed Limit

### Input
- **Rate Limit Config:**
  - MaxStatefulOrdersPerNBlocks: 1 order per 2 blocks
- **Block 2:**
  - Place long-term order: Alice buys 5 at price 10, CLOB 0
- **Block 2:**
  - Attempt to place long-term order: Alice buys 5 at price 10, CLOB 1

### Output
- **First Order CheckTx:** SUCCESS
- **Second Order CheckTx:** FAIL with error `ErrBlockRateLimitExceeded`

### Why It Runs This Way?

1. **Stateful Limit:** Separate limit for stateful orders.
2. **Same Limit Logic:** Works same as short-term order limits.
3. **Per Subaccount:** Limits apply per subaccount.

---

### Test Case 4: Failure - Order Cancellations Exceed Limit

### Input
- **Rate Limit Config:**
  - MaxShortTermOrdersAndCancelsPerNBlocks: 1 operation per 2 blocks
- **Block 2:**
  - Cancel order: Alice cancels order on CLOB 1
- **Block 2:**
  - Attempt to cancel order: Alice cancels order on CLOB 0

### Output
- **First Cancel CheckTx:** SUCCESS
- **Second Cancel CheckTx:** FAIL with error `ErrBlockRateLimitExceeded`

### Why It Runs This Way?

1. **Cancellation Counts:** Cancellations count toward same limit as orders.
2. **Combined Limit:** Orders and cancellations share the same rate limit.
3. **Same Logic:** Works same as order placement limits.

---

### Test Case 5: Failure - Batch Cancellations Exceed Limit

### Input
- **Rate Limit Config:**
  - MaxShortTermOrdersAndCancelsPerNBlocks: 2 operations per 2 blocks
- **Block 2:**
  - Batch cancel: Alice cancels 3 orders (counts as 1 operation)
- **Block 2:**
  - Attempt batch cancel: Alice cancels 3 orders

### Output
- **First Batch Cancel CheckTx:** SUCCESS
- **Second Batch Cancel CheckTx:** FAIL with error `ErrBlockRateLimitExceeded`

### Why It Runs This Way?

1. **Batch Counts as One:** Batch cancel counts as 1 operation, not per order.
2. **Limit Exceeded:** Second batch cancel exceeds limit of 2 per 2 blocks.
3. **Efficiency:** Batch operations are more efficient for rate limits.

---

### Test Case 6: Failure - Leverage Updates Exceed Limit

### Input
- **Rate Limit Config:**
  - MaxLeverageUpdatesPerNBlocks: 1 update per 2 blocks
- **Block 2:**
  - Update leverage: Alice updates leverage for perpetual 0 to 5x
- **Block 2:**
  - Attempt to update leverage: Alice updates leverage for perpetual 1 to 10x

### Output
- **First Update CheckTx:** SUCCESS
- **Second Update CheckTx:** FAIL with error `ErrBlockRateLimitExceeded`

### Why It Runs This Way?

1. **Separate Limit:** Leverage updates have separate rate limit.
2. **Per Subaccount:** Limits apply per subaccount.
3. **Prevents Spam:** Prevents excessive leverage update operations.

---

## Flow Summary

### Rate Limit Check Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. IDENTIFY OPERATION TYPE                                   │
│    - Short-term order/cancel                                 │
│    - Stateful order                                          │
│    - Leverage update                                         │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. GET RATE LIMIT CONFIG                                     │
│    - Find limit for operation type                           │
│    - Get NumBlocks and Limit                                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. COUNT RECENT OPERATIONS                                   │
│    - Count operations in last N blocks                        │
│    - Include current block                                   │
│    - Count per subaccount                                    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. VALIDATE LIMIT                                            │
│    - Check if count >= limit                                 │
│    - Reject if exceeds limit                                 │
│    - Allow if within limit                                   │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Rate Limit Config:**
   ```
   MaxShortTermOrdersAndCancelsPerNBlocks: Limit per N blocks
   MaxStatefulOrdersPerNBlocks: Limit per N blocks
   MaxLeverageUpdatesPerNBlocks: Limit per N blocks
   ```

2. **Operation Counting:**
   ```
   Count operations in sliding window of N blocks
   Include operations from current block
   Count per subaccount
   ```

### Key Points

1. **Per Subaccount Limits:**
   - Limits apply per subaccount
   - Different subaccounts have separate limits
   - Owner address may also be considered

2. **Sliding Window:**
   - Count operations in last N blocks
   - Window slides as blocks advance
   - Operations expire after N blocks

3. **Operation Types:**
   - Short-term orders and cancellations share limit
   - Stateful orders have separate limit
   - Leverage updates have separate limit

4. **Batch Operations:**
   - Batch cancel counts as 1 operation
   - More efficient than individual cancels
   - Still subject to rate limits

5. **CheckTx Validation:**
   - Rate limits checked at CheckTx
   - Early rejection prevents wasted computation
   - Error code: `ErrBlockRateLimitExceeded`

6. **Block Advancement:**
   - Limits persist across blocks
   - Operations count toward limit for N blocks
   - After N blocks, operations no longer count

### Design Rationale

1. **Spam Prevention:** Rate limits prevent order book spam from single subaccount.

2. **Fairness:** Ensures all users have fair access to order book.

3. **System Stability:** Prevents system overload from excessive operations.

4. **Flexibility:** Different limits for different operation types.

