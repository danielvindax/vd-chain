# Test Documentation: Isolated Subaccount Orders E2E Tests

## Overview

This test file verifies **Isolated Subaccount Order** functionality in the CLOB module. Isolated subaccounts are subaccounts that can only trade in specific isolated markets. The test ensures that:
1. Isolated subaccounts cannot place orders for cross-market (non-isolated) perpetuals
2. Isolated subaccounts can place orders for their isolated market
3. Isolated subaccounts can match with other subaccounts in isolated market
4. Isolated subaccounts use isolated collateral pool

---

## Test Function: TestIsolatedSubaccountOrders

### Test Case 1: Failure - Isolated Subaccount Cannot Place Cross-Market Order

### Input
- **Subaccounts:**
  - Alice: 1 ISO Long, 10,000 USD (isolated to ISO market)
  - Bob: 10,000 USD
- **Perpetuals:**
  - BTC-USD (market 0)
  - ETH-USD (market 1)
  - ISO-USD (market 3, isolated)
- **Orders:**
  - Alice: Attempts to buy 5 BTC at price 10 (CLOB 0 - cross-market)
  - Bob: Sells 5 BTC at price 10 (CLOB 0)

### Output
- **Alice Order:** Rejected (invalid for isolated subaccount)
- **Bob Order:** Accepted
- **Orders Filled:** None (Alice's order invalid)
- **Subaccounts:** Unchanged

### Why It Runs This Way?

1. **Isolation Constraint:** Alice's subaccount is isolated to ISO market only.
2. **Cross-Market Rejection:** Cannot place orders for BTC market (market 0).
3. **Validation:** System validates subaccount can trade in requested market.
4. **Protection:** Prevents isolated subaccounts from trading cross-markets.

---

### Test Case 2: Success - Isolated Subaccount Places Order in Isolated Market

### Input
- **Subaccounts:**
  - Alice: 1 ISO Long, 10,000 USD (isolated to ISO market)
  - Bob: 10,000 USD
- **Perpetuals:**
  - ISO-USD (market 3, isolated)
- **Orders:**
  - Alice: Buys 1 ISO at price 10 (CLOB 3 - isolated market)
  - Bob: Sells 1 ISO at price 10 (CLOB 3)

### Output
- **Both Orders:** Accepted
- **Orders Matched:** Both orders fully matched
- **Alice Position:** 2 ISO Long (1 existing + 1 from match)
- **Bob Position:** 1 ISO Short

### Why It Runs This Way?

1. **Isolated Market:** Alice can trade in ISO market (her isolated market).
2. **Order Matching:** Orders match normally in isolated market.
3. **Position Updates:** Positions updated correctly after match.

---

### Test Case 3: Success - Isolated Subaccount Uses Isolated Collateral Pool

### Input
- **Subaccounts:**
  - Alice: 1 ISO Long, 10,000 USD (isolated)
  - Bob: 10,000 USD
- **Collateral Pools:**
  - Main pool: 10,000 USD
  - ISO isolated pool: 10,000 USD
- **Orders:**
  - Alice: Buys 1 ISO at price 10
  - Bob: Sells 1 ISO at price 10

### Output
- **Orders Matched:** Successfully
- **Collateral Pool:** ISO isolated pool used for Alice
- **Main Pool:** Not affected

### Why It Runs This Way?

1. **Isolated Pool:** Isolated subaccounts use separate collateral pool.
2. **Pool Isolation:** Main pool and isolated pools are separate.
3. **Risk Isolation:** Isolated markets have isolated risk.

---

## Flow Summary

### Isolated Subaccount Order Validation

```
┌─────────────────────────────────────────────────────────────┐
│ 1. RECEIVE ORDER                                             │
│    - Order specifies CLOB pair / perpetual                   │
│    - Order specifies subaccount ID                           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CHECK SUBACCOUNT ISOLATION                                │
│    - Query subaccount's isolated market (if any)              │
│    - Check if subaccount is isolated                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. VALIDATE MARKET                                           │
│    - If isolated: Check if order market = isolated market    │
│    - If not isolated: Allow any market                        │
│    - Reject if isolated subaccount tries cross-market         │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. PROCESS ORDER                                             │
│    - If valid: Process normally                              │
│    - Use isolated collateral pool if applicable               │
│    - Update positions in isolated market                      │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Subaccount Isolation:**
   ```
   Non-Isolated: Can trade any market
   Isolated: Can only trade isolated market
   ```

2. **Collateral Pools:**
   ```
   Main Pool: For non-isolated markets
   Isolated Pool: For isolated markets (per market)
   ```

### Key Points

1. **Isolation Constraint:**
   - Isolated subaccounts can only trade their isolated market
   - Cannot place orders for other markets
   - Validation at CheckTx

2. **Market Matching:**
   - Isolated subaccounts can match with any subaccount in isolated market
   - Matching works normally within isolated market
   - Cross-market matching prevented

3. **Collateral Pools:**
   - Isolated markets have separate collateral pools
   - Isolated subaccounts use isolated pool
   - Risk isolated per market

4. **Position Management:**
   - Positions tracked per market
   - Isolated positions separate from cross-market positions
   - Collateral requirements per pool

5. **Validation:**
   - CheckTx validates market access
   - Early rejection for invalid orders
   - Clear error messages

### Design Rationale

1. **Risk Isolation:** Isolated markets prevent risk spillover to main system.

2. **Capital Efficiency:** Isolated markets can have different risk parameters.

3. **Market Segregation:** Prevents isolated subaccounts from affecting main markets.

4. **Flexibility:** Allows new markets with different risk profiles.

