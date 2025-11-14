# Test Documentation: Isolated Subaccount Transfers E2E Tests

## Overview

This test file verifies **Isolated Subaccount Transfer** functionality in the Sending module. Isolated subaccounts are subaccounts that are isolated to specific perpetual markets. The test ensures that:
1. Transfers between isolated and non-isolated subaccounts work correctly
2. Collateral pools are updated correctly when transferring between different market types
3. Transfers between isolated subaccounts in different markets work correctly
4. Transfers fail when collateral pools have insufficient funds
5. Transfers within the same isolated market don't move collateral

---

## Test Function: TestTransfer_Isolated_Non_Isolated_Subaccounts

### Test Case 1: Success - Transfer from Isolated to Non-Isolated Subaccount

### Input
- **Subaccounts:**
  - Alice_Num0: Isolated subaccount with 1 ISO long position, 10,000 USDC
  - Bob_Num0: Non-isolated subaccount with 10,000 USDC
- **Collateral Pools:**
  - Cross collateral pool: 10,000 USDC
  - Isolated market collateral pool: 10,000 USDC
- **Transfer:**
  - From: Alice_Num0 (isolated)
  - To: Bob_Num0 (non-isolated)
  - Amount: 100 USDC

### Output
- **Subaccounts:**
  - Alice_Num0: Balance decreased by 100 USDC
  - Bob_Num0: Balance increased by 100 USDC
- **Collateral Pools:**
  - Cross collateral pool: Increased by 100 USDC (10,100 USDC)
  - Isolated market collateral pool: Decreased by 100 USDC (9,900 USDC)

### Why It Runs This Way?

1. **Collateral Pool Movement:** When transferring from isolated to non-isolated, collateral moves from isolated pool to cross pool.
2. **Isolation:** Isolated subaccounts have separate collateral pools for each market.
3. **Pool Updates:** Collateral pools must be updated to maintain balance.

---

### Test Case 2: Success - Transfer from Non-Isolated to Isolated Subaccount

### Input
- **Subaccounts:**
  - Alice_Num0: Isolated subaccount with 1 ISO long position, 10,000 USDC
  - Bob_Num0: Non-isolated subaccount with 10,000 USDC
- **Collateral Pools:**
  - Cross collateral pool: 10,000 USDC
  - Isolated market collateral pool: 10,000 USDC
- **Transfer:**
  - From: Bob_Num0 (non-isolated)
  - To: Alice_Num0 (isolated)
  - Amount: 100 USDC

### Output
- **Subaccounts:**
  - Alice_Num0: Balance increased by 100 USDC
  - Bob_Num0: Balance decreased by 100 USDC
- **Collateral Pools:**
  - Cross collateral pool: Decreased by 100 USDC (9,900 USDC)
  - Isolated market collateral pool: Increased by 100 USDC (10,100 USDC)

### Why It Runs This Way?

1. **Collateral Pool Movement:** When transferring from non-isolated to isolated, collateral moves from cross pool to isolated pool.
2. **Reverse Flow:** Opposite direction of Test Case 1.
3. **Pool Updates:** Collateral pools must be updated to maintain balance.

---

### Test Case 3: Success - Transfer Between Isolated Subaccounts in Different Markets

### Input
- **Subaccounts:**
  - Alice_Num0: Isolated subaccount in ISO market, 1 ISO long, 10,000 USDC
  - Bob_Num0: Isolated subaccount in ISO2 market, 1 ISO2 long, 10,000 USDC
- **Collateral Pools:**
  - ISO market collateral pool: 10,000 USDC
  - ISO2 market collateral pool: 10,000 USDC
- **Transfer:**
  - From: Alice_Num0 (ISO market)
  - To: Bob_Num0 (ISO2 market)
  - Amount: 100 USDC

### Output
- **Subaccounts:**
  - Alice_Num0: Balance decreased by 100 USDC
  - Bob_Num0: Balance increased by 100 USDC
- **Collateral Pools:**
  - ISO market collateral pool: Decreased by 100 USDC (9,900 USDC)
  - ISO2 market collateral pool: Increased by 100 USDC (10,100 USDC)

### Why It Runs This Way?

1. **Different Markets:** Isolated subaccounts in different markets have separate collateral pools.
2. **Pool Movement:** Collateral moves from one isolated pool to another.
3. **Isolation:** Each isolated market maintains its own collateral pool.

---

### Test Case 4: Failure - Insufficient Funds in Isolated Collateral Pool

### Input
- **Subaccounts:**
  - Alice_Num0: Isolated subaccount with 1 ISO long position, 10,000 USDC
  - Bob_Num0: Non-isolated subaccount with 10,000 USDC
- **Collateral Pools:**
  - Cross collateral pool: 10,000 USDC
  - Isolated market collateral pool: 0 USDC (empty)
- **Transfer:**
  - From: Alice_Num0 (isolated)
  - To: Bob_Num0 (non-isolated)
  - Amount: 100 USDC

### Output
- **DeliverTx:** FAIL
- **Error:** "insufficient funds"
- **Error Code:** `ErrInsufficientFunds`
- **Subaccounts:** No changes (transfer failed)
- **Collateral Pools:** No changes (transfer failed)

### Why It Runs This Way?

1. **Pool Balance:** Collateral pool must have sufficient funds for transfer.
2. **Validation:** System validates pool balance before allowing transfer.
3. **Failure Handling:** Transfer fails if pool has insufficient funds.

---

### Test Case 5: Failure - Insufficient Funds Between Isolated Markets

### Input
- **Subaccounts:**
  - Alice_Num0: Isolated subaccount in ISO market, 1 ISO long, 10,000 USDC
  - Bob_Num0: Isolated subaccount in ISO2 market, 1 ISO2 long, 10,000 USDC
- **Collateral Pools:**
  - ISO market collateral pool: 0 USDC (empty)
  - ISO2 market collateral pool: 10,000 USDC
- **Transfer:**
  - From: Alice_Num0 (ISO market)
  - To: Bob_Num0 (ISO2 market)
  - Amount: 100 USDC

### Output
- **DeliverTx:** FAIL
- **Error:** "insufficient funds"
- **Error Code:** `ErrInsufficientFunds`
- **Subaccounts:** No changes (transfer failed)
- **Collateral Pools:** No changes (transfer failed)

### Why It Runs This Way?

1. **Pool Balance:** Source collateral pool must have sufficient funds.
2. **Validation:** System validates pool balance before allowing transfer.
3. **Failure Handling:** Transfer fails if source pool has insufficient funds.

---

### Test Case 6: Success - Transfer Within Same Isolated Market

### Input
- **Subaccounts:**
  - Alice_Num0: Isolated subaccount in ISO market, 1 ISO long, 10,000 USDC
  - Bob_Num0: Isolated subaccount in ISO market, 1 ISO long, 10,000 USDC
- **Collateral Pools:**
  - ISO market collateral pool: 10,000 USDC
- **Transfer:**
  - From: Alice_Num0 (ISO market)
  - To: Bob_Num0 (ISO market)
  - Amount: 100 USDC

### Output
- **Subaccounts:**
  - Alice_Num0: Balance decreased by 100 USDC
  - Bob_Num0: Balance increased by 100 USDC
- **Collateral Pools:**
  - ISO market collateral pool: No change (10,000 USDC)

### Why It Runs This Way?

1. **Same Market:** Both subaccounts are in the same isolated market.
2. **No Pool Movement:** Collateral stays within the same pool.
3. **Efficiency:** No need to move collateral between pools.

---

## Flow Summary

### Isolated Subaccount Transfer Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. CREATE TRANSFER MESSAGE                                  │
│    - Sender subaccount ID                                   │
│    - Receiver subaccount ID                                 │
│    - Asset ID and amount                                     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. DETERMINE MARKET TYPES                                   │
│    - Check if sender is isolated                            │
│    - Check if receiver is isolated                           │
│    - Identify market types                                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. VALIDATE COLLATERAL POOLS                                 │
│    - Check source pool balance                               │
│    - Verify sufficient funds                                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. UPDATE COLLATERAL POOLS                                  │
│    - If different markets: Move collateral between pools    │
│    - If same market: No pool movement                        │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. UPDATE SUBACCOUNTS                                       │
│    - Decrease sender balance                                 │
│    - Increase receiver balance                               │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Transfer States:**
   ```
   Create Transfer → Validate Pools → Update Pools → Update Subaccounts → Complete
   ```

2. **Collateral Pool Updates:**
   ```
   Isolated → Non-Isolated: Isolated Pool ↓, Cross Pool ↑
   Non-Isolated → Isolated: Cross Pool ↓, Isolated Pool ↑
   Isolated → Isolated (Different): Source Pool ↓, Dest Pool ↑
   Isolated → Isolated (Same): No Change
   ```

### Key Points

1. **Collateral Pools:**
   - Cross collateral pool: For non-isolated subaccounts
   - Isolated market pools: One pool per isolated market
   - Pools must maintain balance

2. **Transfer Rules:**
   - Different markets: Collateral moves between pools
   - Same market: No pool movement
   - Insufficient funds: Transfer fails

3. **Validation:**
   - Pool balance must be sufficient
   - Transfer amount must be valid
   - Subaccounts must exist

### Design Rationale

1. **Isolation:** Isolated markets maintain separate collateral pools for risk management.

2. **Pool Management:** Collateral pools ensure sufficient funds for positions.

3. **Efficiency:** Same-market transfers don't move collateral for efficiency.

4. **Safety:** Validation prevents transfers when pools have insufficient funds.

