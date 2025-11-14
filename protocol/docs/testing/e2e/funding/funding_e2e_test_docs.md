# Test Documentation: Funding E2E Tests

## Overview

This test file verifies the **Funding Mechanism** for perpetual markets. Funding is a periodic payment between long and short positions based on the premium (difference between mark price and index price). The test ensures that:
1. Funding premiums are calculated correctly based on order book impact prices
2. Funding index is updated correctly based on premiums
3. Funding settlements are calculated and applied correctly to subaccounts
4. Funding rate clamping works when premiums exceed limits

---

## Test Function: TestFunding

### Test Case 1: Index Price Below Impact Bid, Positive Funding, Longs Pay Shorts

### Input
- **Orders:**
  - Unmatched orders to generate funding premiums:
    - Bob: Sell 2 BTC at 28,005 (impact ask)
    - Alice: Buy 2 BTC at 28,000 (impact bid)
  - Matched orders to set up positions:
    - Bob: Sell 1 BTC at 28,003 (matched)
    - Alice: Buy 0.8 BTC at 28,003 (matched)
    - Carl: Buy 0.2 BTC at 28,003 (matched)
- **Initial Index Price:** 28,002
- **Index Price for Premium:** 27,960 (below impact bid)
- **Oracle Price for Funding Index:** 27,000

### Output
- **Funding Premiums:** ~1,430 ppm (0.143%)
- **Funding Index:** 482
- **Settlements:**
  - Alice (long 0.8 BTC): Pays $3.856
  - Bob (short 1 BTC): Receives $4.82
  - Carl (long 0.2 BTC): Pays $0.964

### Why It Runs This Way?

1. **Premium Calculation:** When index price (27,960) is below impact bid (28,000), premium is positive.
   - Premium = (28,000 / 27,960) - 1 ≈ 0.00143 (0.143%)
2. **Longs Pay Shorts:** Positive premium means longs pay shorts.
3. **Funding Index:** Calculated from premium samples over funding epoch.
4. **Settlement:** Applied when subaccount receives transfer, based on funding index difference.

---

### Test Case 2: Index Price Above Impact Ask, Negative Funding, Final Funding Rate Clamped

### Input
- **Orders:** Same as Test Case 1
- **Initial Index Price:** 28,002
- **Index Price for Premium:** 34,000 (above impact ask)
- **Oracle Price for Funding Index:** 33,500

### Output
- **Funding Premiums:** -176,323 ppm (-17.6%, but clamped)
- **Funding Index:** -50,250 (clamped to -12% based on margin requirements)
- **Settlements:**
  - Alice (long 0.8 BTC): Receives $402 (shorts pay longs)
  - Bob (short 1 BTC): Pays $502.5
  - Carl (long 0.2 BTC): Receives $100.5

### Why It Runs This Way?

1. **Premium Calculation:** When index price (34,000) is above impact ask (28,005), premium is negative.
   - Premium = (28,005 / 34,000) - 1 ≈ -0.176 (-17.6%)
2. **Funding Rate Clamp:** Funding rate is clamped to prevent excessive payments.
   - Clamp = premium_rate_clamp_factor × (initial_margin - maintenance_margin)
   - Clamp = 600% × (5% - 3%) = 12% = 120,000 ppm
3. **Shorts Pay Longs:** Negative premium means shorts pay longs (after clamping).
4. **Settlement:** Applied with clamped funding rate.

---

### Test Case 3: Index Price Between Impact Bid and Ask, Zero Funding

### Input
- **Orders:** Same as Test Case 1
- **Initial Index Price:** 28,002
- **Index Price for Premium:** 28,003 (between impact bid 28,000 and ask 28,005)
- **Oracle Price for Funding Index:** 27,500

### Output
- **Funding Premiums:** None (zero)
- **Funding Index:** 0
- **Settlements:**
  - Alice: $0
  - Bob: $0
  - Carl: $0

### Why It Runs This Way?

1. **Zero Premium:** When index price is between impact bid and ask, premium is zero.
2. **No Funding:** No funding payments when premium is zero.
3. **Funding Index:** Remains unchanged (starts at 0, stays at 0).

---

## Flow Summary

### Funding Calculation Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. PLACE ORDERS                                              │
│    - Place unmatched orders to set impact prices             │
│    - Place matched orders to open positions                  │
│    - Update initial index price                               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. FUNDING SAMPLE EPOCHS                                     │
│    - Advance to funding tick epoch                            │
│    - Update index price for premium calculation               │
│    - Collect premium samples (60 samples per funding tick)   │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. CALCULATE FUNDING PREMIUMS                                │
│    - Calculate premium = (impact_price / index_price) - 1    │
│    - If index < impact_bid: positive premium (longs pay)     │
│    - If index > impact_ask: negative premium (shorts pay)    │
│    - If index between bid/ask: zero premium                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. UPDATE FUNDING INDEX                                      │
│    - Calculate funding index from premium samples            │
│    - Apply funding rate clamp if needed                      │
│    - Update perpetual funding index                           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. SETTLE FUNDING                                            │
│    - When subaccount receives transfer                       │
│    - Calculate settlement = (funding_index - position_index) × size │
│    - Apply settlement to subaccount balance                  │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Premium Calculation:**
   ```
   Index < Impact Bid → Positive Premium → Longs Pay Shorts
   Index > Impact Ask → Negative Premium → Shorts Pay Longs
   Index Between → Zero Premium → No Funding
   ```

2. **Funding Index:**
   ```
   Initial: 0
   After Epoch: Updated based on premium samples
   Clamped: If premium exceeds margin-based limit
   ```

3. **Settlement:**
   ```
   Position Opened: Funding Index = 0
   After Funding: Funding Index Updated
   On Transfer: Settlement = (New Index - Old Index) × Size
   ```

### Key Points

1. **Impact Prices:**
   - Impact bid: Best bid price from order book
   - Impact ask: Best ask price from order book
   - Used to calculate premium when index price is outside range

2. **Premium Calculation:**
   - Premium = (Impact Price / Index Price) - 1
   - Positive: Longs pay shorts
   - Negative: Shorts pay longs
   - Zero: No funding

3. **Funding Rate Clamp:**
   - Prevents excessive funding payments
   - Based on margin requirements
   - Formula: clamp_factor × (initial_margin - maintenance_margin)

4. **Funding Index:**
   - Tracks cumulative funding over time
   - Updated at each funding tick
   - Used to calculate settlement when position is closed or transferred

5. **Settlement:**
   - Calculated when subaccount receives transfer
   - Settlement = (Current Funding Index - Position Funding Index) × Position Size
   - Applied to subaccount balance

6. **Premium Samples:**
   - 60 samples collected per funding tick epoch
   - Samples used to calculate average premium
   - Premium samples reset at start of new funding tick epoch

### Design Rationale

1. **Fairness:** Funding ensures long and short positions are balanced by transferring value based on premium.

2. **Price Discovery:** Premium reflects the difference between mark price (order book) and index price (oracle).

3. **Safety:** Funding rate clamping prevents excessive payments that could cause liquidations.

4. **Efficiency:** Funding index allows efficient settlement calculation without recalculating all historical premiums.

5. **Transparency:** Premium samples are collected over time to ensure fair and accurate funding calculation.

