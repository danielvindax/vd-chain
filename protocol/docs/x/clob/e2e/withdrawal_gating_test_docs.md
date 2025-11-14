# Test Documentation: Withdrawal Gating E2E Tests

## Overview

This test file verifies **Withdrawal Gating** functionality in the CLOB module. When a subaccount has negative Total Net Collateral (TNC) and cannot be deleveraged, withdrawals and transfers from that market are blocked (gated) to protect the system. The test ensures that:
1. Withdrawals are gated when negative TNC subaccounts exist
2. Gating applies to isolated markets separately
3. Gating blocks withdrawals for affected markets
4. Gating unblocks when negative TNC resolved

---

## Test Function: TestWithdrawalGating_NegativeTncSubaccount_BlocksThenUnblocks

### Test Case 1: Withdrawals Gated - Non-Overlapping Bankruptcy Prices

### Input
- **Subaccounts:**
  - Carl: 1 BTC Short, 49,999 USD (negative TNC, undercollateralized)
  - Dave: 1 BTC Long, 50,000 USD (short)
  - Dave_Num1: 10,000 USD
- **Oracle Price:** $50,500 / BTC
- **Liquidation Order:** Dave sells 0.25 BTC at $50,000
- **Liquidation:** Attempted but deleveraging fails (non-overlapping bankruptcy prices)
- **Withdrawal:** Dave_Num1 attempts to withdraw from BTC market

### Output
- **Liquidation:** Fails (deleveraging cannot be performed)
- **Carl State:** Still has negative TNC
- **Withdrawals Gated:** BTC market withdrawals blocked
- **Error:** "WithdrawalsAndTransfersBlocked: failed to apply subaccount updates"
- **Gated Perpetual:** BTC perpetual ID marked as gated
- **Negative TNC Seen At Block:** Block 4

### Why It Runs This Way?

1. **Negative TNC:** Carl has negative TNC (49,999 < 50,000 needed).
2. **Deleveraging Fails:** Cannot deleverage because bankruptcy prices don't overlap.
3. **System Protection:** Withdrawals gated to prevent further capital outflow.
4. **Market Isolation:** Gating applies to specific perpetual/market.

---

### Test Case 2: Withdrawals Gated - Isolated Market

### Input
- **Subaccounts:**
  - Carl: 1 ISO Short, 49 USD (negative TNC)
  - Dave: 1 ISO Long, 50 USD (short)
  - Alice: 1 ISO Long, 10,000 USD (isolated subaccount)
- **Oracle Price:** $50.5 / ISO
- **Liquidation:** Attempted but deleveraging fails
- **Withdrawal:** Alice attempts to withdraw from ISO market

### Output
- **Withdrawals Gated:** ISO market withdrawals blocked for isolated subaccounts
- **Gated Perpetual:** ISO perpetual ID marked as gated
- **Error:** "WithdrawalsAndTransfersBlocked"

### Why It Runs This Way?

1. **Isolated Market:** ISO is isolated market with separate collateral pool.
2. **Isolated Subaccount:** Alice has isolated subaccount for ISO market.
3. **Market-Specific Gating:** Gating applies only to ISO market for isolated subaccounts.
4. **Protection:** Prevents capital outflow from isolated market when negative TNC exists.

---

### Test Case 3: Withdrawals Not Gated - Non-Isolated Subaccount

### Input
- **Subaccounts:**
  - Carl: 1 ISO Short, 49 USD (negative TNC)
  - Dave: 1 ISO Long, 50 USD (short)
  - Alice: 10,000 USD (non-isolated subaccount)
- **Oracle Price:** $50.5 / ISO
- **Liquidation:** Attempted but deleveraging fails
- **Withdrawal:** Alice attempts to withdraw (not from ISO market)

### Output
- **Withdrawals Not Gated:** Alice can withdraw (not from ISO market)
- **Gated Perpetual:** ISO perpetual ID still marked as gated
- **Selective Gating:** Only ISO market withdrawals blocked

### Why It Runs This Way?

1. **Non-Isolated Subaccount:** Alice doesn't have ISO position.
2. **Selective Gating:** Gating only affects withdrawals from gated market.
3. **Other Markets:** Withdrawals from other markets still allowed.

---

## Flow Summary

### Withdrawal Gating Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. DETECT NEGATIVE TNC                                       │
│    - Subaccount has TNC < 0                                  │
│    - Liquidations daemon identifies negative TNC accounts    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. ATTEMPT DELEVERAGING                                      │
│    - Try to deleverage negative TNC account                  │
│    - Check for overlapping bankruptcy prices                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. DELEVERAGING FAILS                                        │
│    - Bankruptcy prices don't overlap                          │
│    - Cannot close position                                   │
│    - Negative TNC persists                                   │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. GATE WITHDRAWALS                                          │
│    - Mark perpetual as gated                                 │
│    - Block withdrawals from gated market                      │
│    - Record block when negative TNC seen                      │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. BLOCK WITHDRAWAL ATTEMPTS                                 │
│    - Reject withdrawal transactions                          │
│    - Return error: "WithdrawalsAndTransfersBlocked"          │
│    - Apply to affected markets only                           │
└─────────────────────────────────────────────────────────────┘
```

### Gating Resolution

```
┌─────────────────────────────────────────────────────────────┐
│ 1. RESOLVE NEGATIVE TNC                                      │
│    - Position closed through deleveraging                    │
│    - Or position closed through matching                     │
│    - Or collateral added                                     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. UNGATE WITHDRAWALS                                        │
│    - No more negative TNC accounts                           │
│    - Remove gating from perpetual                            │
│    - Allow withdrawals again                                 │
└─────────────────────────────────────────────────────────────┘
```

### Key Points

1. **Negative TNC Detection:**
   - Subaccount TNC < 0
   - Detected by liquidations daemon
   - Recorded per perpetual/market

2. **Deleveraging Failure:**
   - Bankruptcy prices don't overlap
   - Cannot find counterparty to deleverage
   - Negative TNC cannot be resolved

3. **Gating Mechanism:**
   - Perpetual marked as gated
   - Withdrawals blocked for gated market
   - Transfers also blocked

4. **Market Isolation:**
   - Isolated markets gated separately
   - Non-isolated subaccounts not affected by isolated market gating
   - Cross-market gating possible

5. **Block Tracking:**
   - Block when negative TNC first seen
   - Used for gating duration tracking
   - Helps identify persistent issues

6. **Error Handling:**
   - Clear error message: "WithdrawalsAndTransfersBlocked"
   - Transaction rejected at CheckTx or DeliverTx
   - State protected from capital outflow

### Design Rationale

1. **System Protection:** Prevents capital outflow when system is at risk.

2. **Risk Containment:** Gating limits exposure to negative TNC accounts.

3. **Market Isolation:** Isolated markets protected separately.

4. **Fairness:** Only affects withdrawals from affected markets.

5. **Transparency:** Clear error messages and block tracking.

