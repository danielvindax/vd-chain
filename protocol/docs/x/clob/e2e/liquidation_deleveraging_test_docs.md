# Test Documentation: Liquidation and Deleveraging E2E Tests

## Overview

This test file verifies **Liquidation and Deleveraging** functionality in the CLOB module. When a subaccount becomes undercollateralized (below maintenance margin), it can be liquidated. The test ensures that:
1. Liquidations respect position block limits (MinPositionNotionalLiquidated, MaxPositionPortionLiquidatedPpm)
2. Liquidations respect subaccount block limits (MaxNotionalLiquidated)
3. Liquidations work for both long and short positions
4. Insurance fund covers losses when needed

---

## Test Function: TestLiquidationConfig

### Test Case 1: Liquidating Short - Respects MinPositionNotionalLiquidated

### Input
- **Subaccounts:**
  - Carl: 1 BTC Short, 50,499 USD collateral (undercollateralized)
  - Dave: 1 BTC Long, 50,000 USD collateral
- **Order:** Dave sells 1 BTC at 50,000
- **Liquidation Config:**
  - MinPositionNotionalLiquidated: $100,000
  - MaxPositionPortionLiquidatedPpm: 1% (10,000 ppm)
  - Oracle Price: 50,000

### Output
- **Liquidation:** Entire position liquidated (1 BTC)
- **Carl Balance:** 50,499 - 50,000 - 250 (fees) = 249 USD
- **Dave Balance:** 50,000 + 50,000 = 100,000 USD

### Why It Runs This Way?

1. **Minimum Notional:** 1% of $50,000 = $500, but minimum is $100,000.
2. **Entire Position:** Since $500 < $100,000, entire position is liquidated.
3. **Full Liquidation:** All 1 BTC is liquidated to meet minimum requirement.

---

### Test Case 2: Liquidating Long - Respects MaxPositionPortionLiquidatedPpm

### Input
- **Subaccounts:**
  - Carl: 1 BTC Short, 100,000 USD
  - Dave: 1 BTC Long, 49,501 USD (undercollateralized)
- **Order:** Carl buys 1 BTC at 50,000
- **Liquidation Config:**
  - MinPositionNotionalLiquidated: $1,000
  - MaxPositionPortionLiquidatedPpm: 10% (100,000 ppm)
  - Oracle Price: 50,000

### Output
- **Liquidation:** 10% of position liquidated (0.1 BTC)
- **Dave Balance:** -49,501 + 5,000 - 25 (fees) = -44,526 USD
- **Dave Position:** 0.9 BTC long remaining

### Why It Runs This Way?

1. **Portion Limit:** 10% of $50,000 = $5,000 worth of BTC.
2. **Partial Liquidation:** Only 0.1 BTC (10%) is liquidated.
3. **Remaining Position:** 0.9 BTC position remains.

---

### Test Case 3: Liquidating Short - Respects MaxNotionalLiquidated

### Input
- **Subaccounts:**
  - Carl: 1 BTC Short, 50,499 USD (undercollateralized)
  - Dave: 1 BTC Long, 50,000 USD
- **Order:** Dave sells 1 BTC at 49,500
- **Liquidation Config:**
  - MaxNotionalLiquidated: $5,000 per block
  - Oracle Price: 50,000

### Output
- **Liquidation:** Only $5,000 worth liquidated (0.1 BTC)
- **Carl Balance:** 50,499 - 5,000 - 25 (fees) = 45,474 USD
- **Carl Position:** 0.9 BTC short remaining

### Why It Runs This Way?

1. **Subaccount Limit:** Maximum $5,000 can be liquidated per block.
2. **Partial Liquidation:** Only 0.1 BTC ($5,000 worth) is liquidated.
3. **Remaining Position:** 0.9 BTC position remains for future liquidation.

---

## Flow Summary

### Liquidation Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. DETECT UNDERCOLLATERALIZATION                            │
│    - Subaccount TNC < maintenance margin                    │
│    - Liquidations daemon identifies liquidatable accounts    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CALCULATE LIQUIDATION AMOUNT                              │
│    - Apply position block limits                             │
│      * MinPositionNotionalLiquidated                         │
│      * MaxPositionPortionLiquidatedPpm                       │
│    - Apply subaccount block limits                           │
│      * MaxNotionalLiquidated                                 │
│    - Take minimum of all limits                              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. FIND MATCHING ORDERS                                      │
│    - Search order book for matching orders                   │
│    - Use fillable price config                               │
│    - Match at best available price                           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. EXECUTE LIQUIDATION                                       │
│    - Close portion of position                              │
│    - Transfer funds to counterparty                          │
│    - Charge liquidation fee                                 │
│    - Update subaccount state                                 │
└─────────────────────────────────────────────────────────────┘
```

### Liquidation Limits

1. **Position Block Limits:**
   - MinPositionNotionalLiquidated: Minimum $ amount to liquidate
   - MaxPositionPortionLiquidatedPpm: Maximum % of position to liquidate
   - Applied per position per block

2. **Subaccount Block Limits:**
   - MaxNotionalLiquidated: Maximum $ amount per subaccount per block
   - MaxQuantumsInsuranceLost: Maximum insurance fund loss per block
   - Applied per subaccount per block

### Key Points

1. **Liquidation Triggers:**
   - Subaccount TNC < maintenance margin
   - Detected by liquidations daemon
   - Liquidatable accounts identified each block

2. **Liquidation Amount:**
   - Calculated based on multiple limits
   - Minimum of position limits and subaccount limits
   - Ensures controlled liquidation rate

3. **Price Discovery:**
   - Uses fillable price config
   - Matches at best available price on order book
   - May use insurance fund if no matching orders

4. **Liquidation Fees:**
   - Charged to liquidated account
   - MaxLiquidationFeePpm sets maximum fee
   - Fees compensate liquidators

5. **Partial Liquidation:**
   - Can liquidate portion of position
   - Remaining position stays open
   - Can be liquidated again in future blocks

6. **Insurance Fund:**
   - Covers losses when liquidation price is unfavorable
   - MaxQuantumsInsuranceLost limits fund exposure
   - Protects protocol from excessive losses

### Design Rationale

1. **Risk Management:** Liquidation limits prevent excessive liquidation in single block.

2. **Market Stability:** Controlled liquidation rate prevents market disruption.

3. **Fairness:** Limits ensure all liquidatable accounts get fair treatment.

4. **Safety:** Insurance fund protects protocol from extreme market conditions.

