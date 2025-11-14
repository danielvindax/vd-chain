# Test Documentation: TWAP Orders E2E Tests

## Overview

This test file verifies **TWAP (Time-Weighted Average Price) Order** functionality in the CLOB module. TWAP orders split a large order into multiple smaller suborders executed over time at regular intervals. The test ensures that:
1. TWAP orders are split into suborders based on duration and interval
2. Suborders are placed at regular intervals
3. Suborders use oracle price with price tolerance
4. TWAP orders catch up if suborders expire unfilled
5. Duplicate TWAP orders are rejected

---

## Test Function: TestTwapOrderPlacementAndCatchup

### Test Case: Success - TWAP Order Placement and Suborder Execution

### Input
- **TWAP Order:**
  - SubaccountId: Alice_Num0
  - Side: BUY
  - Quantums: 100,000,000,000 (10 BTC)
  - Duration: 300 seconds (5 minutes)
  - Interval: 60 seconds (1 minute)
  - PriceTolerance: 0% (market order)
  - GoodTilBlockTime: 300 seconds from now

### Output
- **TWAP Order Placement:**
  - RemainingLegs: 4 (5 total - 1 triggered)
  - RemainingQuantums: 100,000,000,000
- **First Suborder:**
  - Quantums: 20,000,000,000 (100B / 5 = 20B per leg)
  - Subticks: 200,000,000 ($20,000 oracle price)
  - GoodTilBlockTime: 3 seconds from now
  - Side: BUY (same as parent)
- **After 30 Seconds:**
  - Suborder expired and removed
  - TWAP order still has 4 remaining legs
- **After 60 Seconds Total:**
  - Second suborder placed
  - Quantums: 25,000,000,000 (100B / 4 = 25B, catching up)
  - Subticks: 200,000,000 (same oracle price)

### Why It Runs This Way?

1. **Suborder Calculation:** Total quantums divided by number of legs.
   - 5 legs: 100B / 5 = 20B per leg
   - After 1 leg: 100B / 4 = 25B per leg (catchup)
2. **Interval-Based:** Suborders placed at regular intervals (60 seconds).
3. **Oracle Price:** Suborders use current oracle price.
4. **Catchup Logic:** If suborder expires unfilled, next suborder gets larger size to catch up.

---

## Test Function: TestDuplicateTWAPOrderPlacement

### Test Case: Failure - Duplicate TWAP Order

### Input
- **Block 1:**
  - Place TWAP order: Alice buys 100B quantums over 4 legs
- **Block 2:**
  - Attempt to place same TWAP order (same OrderId)

### Output
- **First Order CheckTx:** SUCCESS
- **Second Order CheckTx:** FAIL
- **Error:** "A stateful order with this OrderId already exists"

### Why It Runs This Way?

1. **Duplicate Detection:** System detects duplicate OrderId.
2. **Stateful Order:** TWAP orders are stateful orders.
3. **Rejection:** Cannot place duplicate stateful order.

---

## Test Function: TestTWAPOrderWithMatchingOrders

### Test Case: Success - TWAP Suborder Matches with Existing Order

### Input
- **TWAP Order:** Alice buys 100B quantums over 4 legs
- **Existing Order:** Bob sells matching quantity at compatible price
- **Suborder:** First suborder placed and matches with Bob's order

### Output
- **Suborder:** Fully filled
- **TWAP Order:** Remaining quantums reduced
- **Next Suborder:** Next suborder placed at next interval

### Why It Runs This Way?

1. **Suborder Matching:** TWAP suborders can match like regular orders.
2. **Fill Tracking:** Fills reduce remaining quantums in TWAP order.
3. **Continued Execution:** TWAP order continues placing suborders until complete.

---

## Flow Summary

### TWAP Order Execution Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. PLACE TWAP ORDER                                         │
│    - Specify total quantums, duration, interval              │
│    - Calculate number of legs = duration / interval          │
│    - Store TWAP order in state                               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. PLACE FIRST SUBORDER                                      │
│    - Calculate suborder size = total / num_legs              │
│    - Get current oracle price                                │
│    - Apply price tolerance                                   │
│    - Place suborder on order book                            │
│    - Schedule next suborder trigger                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. SUBORDER EXECUTION                                        │
│    - Suborder can match with existing orders                  │
│    - If filled: Update remaining quantums                     │
│    - If expired: Remove and catch up                         │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. TRIGGER NEXT SUBORDER                                     │
│    - At next interval time                                   │
│    - Calculate catchup size if previous expired              │
│    - Place next suborder                                     │
│    - Repeat until all legs executed                          │
└─────────────────────────────────────────────────────────────┘
```

### Catchup Logic

```
If previous suborder expired unfilled:
  Remaining quantums = original remaining - 0 (nothing filled)
  Remaining legs = original legs - 1
  Next suborder size = remaining quantums / remaining legs
  
Example:
  Original: 100B quantums, 5 legs
  First suborder: 20B (expires unfilled)
  Catchup: 100B / 4 = 25B per remaining leg
```

### Key Points

1. **TWAP Parameters:**
   - Duration: Total time to execute order
   - Interval: Time between suborders
   - PriceTolerance: Maximum price deviation from oracle
   - Number of legs = duration / interval

2. **Suborder Calculation:**
   - Initial: total_quantums / num_legs
   - Catchup: remaining_quantums / remaining_legs
   - Ensures all quantums executed by end of duration

3. **Oracle Price:**
   - Suborders use current oracle price
   - Price tolerance allows deviation
   - Market orders: tolerance = 0

4. **Suborder Lifecycle:**
   - Placed at interval time
   - Can match with existing orders
   - Expires if not filled by GoodTilBlockTime
   - Removed and next suborder catches up

5. **State Tracking:**
   - TWAP order placement tracked in state
   - Remaining legs and quantums tracked
   - Trigger times scheduled for next suborders

6. **Duplicate Prevention:**
   - Cannot place duplicate TWAP order
   - Same OrderId rejected
   - Prevents accidental duplicate execution

### Design Rationale

1. **Price Impact Reduction:** Splitting large orders reduces market impact.

2. **Time Distribution:** Executes order over time for better average price.

3. **Flexibility:** Configurable duration and interval for different strategies.

4. **Catchup Logic:** Ensures order completes even if some suborders expire.

5. **Oracle Integration:** Uses oracle price for fair execution.

