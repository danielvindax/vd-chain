# Test Documentation: Builder Code Orders E2E Tests

## Overview

This test file verifies **Builder Code** functionality in the CLOB module. Builder codes allow order builders to receive fees when their orders are matched. The test ensures that:
1. Orders with builder codes can be placed successfully
2. Builder fees are correctly calculated and paid when orders match
3. Orders without builder codes work normally
4. Fee calculation is based on the fill amount and fee percentage

---

## Test Function: TestBuilderCodeOrders

### Test Case: Order with Builder Code Fills and Fees are Paid

### Input
- **Orders:**
  - Alice: Buy order with builder code
    - SubaccountId: Alice_Num0
    - ClobPairId: 0
    - Side: BUY
    - Quantums: 10,000,000,000 (1 BTC)
    - Subticks: 500,000,000 (50,000 USDC/BTC)
    - GoodTilBlock: 20
    - BuilderCodeParameters:
      - BuilderAddress: Carl_Num0 account address
      - FeePpm: 1000 (0.1%)
  - Bob: Sell order without builder code
    - SubaccountId: Bob_Num0
    - ClobPairId: 0
    - Side: SELL
    - Quantums: 10,000,000,000 (1 BTC)
    - Subticks: 500,000,000 (50,000 USDC/BTC)
    - GoodTilBlock: 20
- **Match Operation:** Orders match at block 2

### Output
- **CheckTx:** Both orders pass CheckTx validation
- **Order Fill:** Both orders are fully filled (10,000,000,000 quantums)
- **Builder Fee:** Carl receives 50,000,000 quantums (0.1% of 50,000 USDC)
- **Builder Balance:** Carl's balance increases by the builder fee amount

### Why It Runs This Way?

1. **Builder Code Mechanism:** Tests the fee-sharing mechanism where order builders receive fees.
2. **Fee Calculation:** Fee = (Fill Amount × Price × FeePpm) / 1,000,000
   - Fill Amount: 10,000,000,000 quantums (1 BTC)
   - Price: 500,000,000 subticks (50,000 USDC/BTC)
   - FeePpm: 1000 (0.1%)
   - Fee = (10,000,000,000 × 500,000,000 × 1000) / (1,000,000 × 10^8) = 50,000,000 quantums
3. **Fee Payment:** Builder fee is paid from the matched order proceeds to the builder address.
4. **Order Matching:** Both orders match completely, triggering fee payment.
5. **Balance Verification:** Builder's balance is checked before and after match to verify fee payment.

---

## Flow Summary

### Builder Code Order Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. PLACE ORDERS                                              │
│    - Alice places buy order with builder code                │
│    - Bob places sell order without builder code              │
│    - Both orders pass CheckTx                                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. ADVANCE BLOCK                                             │
│    - Orders are matched in block 2                           │
│    - Match operation is processed                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. CALCULATE BUILDER FEE                                     │
│    - Fee = (Fill Amount × Price × FeePpm) / 1,000,000      │
│    - Fee deducted from order proceeds                        │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. PAY BUILDER FEE                                           │
│    - Fee transferred to builder address                      │
│    - Builder balance increases                               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. VERIFY RESULTS                                            │
│    - Orders are fully filled                                 │
│    - Builder fee is paid correctly                           │
│    - Builder balance matches expected                        │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Order States:**
   ```
   Placed → CheckTx Passed → Matched → Filled → Fee Paid
   ```

2. **Builder Fee Calculation:**
   ```
   Fill Amount × Price × FeePpm / 1,000,000 = Fee
   ```

3. **Balance Updates:**
   ```
   Pre-Match Balance → Match → Fee Payment → Post-Match Balance
   ```

### Key Points

1. **Builder Code Parameters:**
   - BuilderAddress: Address that receives the fee
   - FeePpm: Fee percentage in parts per million (1000 = 0.1%)

2. **Fee Payment:**
   - Fee is paid from the matched order proceeds
   - Only orders with builder codes generate fees
   - Fee is calculated based on fill amount and price

3. **Order Matching:**
   - Orders must match for fees to be paid
   - Partial fills result in proportional fees
   - Full fills result in full fee calculation

4. **Balance Verification:**
   - Builder balance is checked before match
   - Balance delta is calculated after match
   - Delta should equal the builder fee

### Design Rationale

1. **Fee Sharing:** Builder codes incentivize order builders to provide liquidity.

2. **Fee Calculation:** Fee is proportional to trade size and price, ensuring fair compensation.

3. **Flexibility:** Orders can have builder codes or not, allowing for different fee structures.

4. **Verification:** Balance checks ensure fees are correctly calculated and paid.

