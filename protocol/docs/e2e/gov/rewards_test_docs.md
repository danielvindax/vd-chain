# Test Documentation: Rewards Module Governance Proposals

## Overview

This test file verifies governance proposal to **update Rewards Module Params**. Rewards module manages reward distribution to users based on trading activity.

---

## Test Function: TestUpdateRewardsModuleParams

### Test Case 1: Success - Update Rewards Module Params

### Input
- **Genesis State:**
  - Rewards module has params:
    - TreasuryAccount: "test_treasury"
    - Denom: "avdtn"
    - DenomExponent: -18
    - MarketId: 1234
    - FeeMultiplierPpm: 700_000
- **Proposed Message:**
  - `MsgUpdateParams`:
    - TreasuryAccount: "test_treasury" (unchanged)
    - Denom: "avdtn" (unchanged)
    - DenomExponent: -5 (changed from -18)
    - MarketId: 0 (changed from 1234)
    - FeeMultiplierPpm: 700_001 (changed from 700_000)
    - Authority: gov module

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Module Params:** Updated with new values

### Why It Runs This Way?

1. **Rewards Configuration:** These params control how rewards are calculated and distributed.
2. **FeeMultiplierPpm:** This is the multiplier (parts per million) to calculate rewards based on trading fees.
3. **DenomExponent:** Exponent of denom to convert between different units.
4. **MarketId:** Market ID to track rewards for specific market.

---

### Test Case 2: Failure - Treasury Account is Empty

### Input
- **Proposed Message:**
  - `MsgUpdateParams` with TreasuryAccount = "" (empty string)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **Required Field:** TreasuryAccount is the source of funds to distribute rewards, cannot be empty.
2. **Early Validation:** Validation at CheckTx to reject early.
3. **Data Integrity:** Ensures rewards module always has valid treasury account.

---

### Test Case 3: Failure - Denom is Invalid

### Input
- **Proposed Message:**
  - `MsgUpdateParams` with Denom = "7avdtn" (invalid - starts with number)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **Denom Format:** Denom must follow standard format (cannot start with number).
2. **Validation:** Cosmos SDK has rules about denom format.
3. **Early Rejection:** Validation at CheckTx to prevent invalid proposals.

---

### Test Case 4: Failure - Fee Multiplier PPM Greater Than 1 Million

### Input
- **Proposed Message:**
  - `MsgUpdateParams` with FeeMultiplierPpm = 1_000_001 (> 1 million)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **Boundary Check:** FeeMultiplierPpm must be <= 1_000_000 (1 million ppm = 100%).
2. **PPM Format:** PPM (parts per million) has maximum of 1 million.
3. **Logical Constraint:** Multiplier cannot be > 100% in this context.

---

### Test Case 5: Failure - Invalid Authority

### Input
- **Proposed Message:**
  - `MsgUpdateParams` with Authority = rewards module address (instead of gov module)

### Output
- **Proposal Submission:** FAIL
- **Proposals:** No proposals created

### Why It Runs This Way?

1. **Authority Check:** Only governance module has permission to update rewards params.
2. **Security:** Ensures only governance can change rewards configuration.
3. **Early Rejection:** Validation at proposal submission time.

---

## Flow Summary

### Update Rewards Module Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - TreasuryAccount not empty                              │
│    - Denom valid (standard format)                          │
│    - FeeMultiplierPpm <= 1_000_000                          │
│    - Authority = gov module                                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                                │
│    - Validate authority                                      │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTION                                        │
│    - Update rewards module params                            │
│    - Apply new configuration                                 │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Params Update:**
   ```
   Old Params → New Params (UPDATE)
   ```

2. **Validation Points:**
   ```
   CheckTx: TreasuryAccount, Denom, FeeMultiplierPpm
   Submission: Authority
   ```

### Key Points

1. **TreasuryAccount:**
   - Must not be empty
   - This is the source of funds to distribute rewards
   - Validation at CheckTx

2. **Denom Format:**
   - Must follow Cosmos SDK denom format
   - Cannot start with number
   - Validation at CheckTx

3. **FeeMultiplierPpm:**
   - Must be <= 1_000_000 (100%)
   - This is the multiplier to calculate rewards
   - Validation at CheckTx

4. **Authority:**
   - Only governance module has permission
   - Validation at proposal submission time

5. **Atomic Execution:**
   - If validation fails, entire proposal fails
   - State is rolled back to before execution

### Design Rationale

1. **Governance Control:** Only governance has permission to change rewards configuration to ensure decentralization.

2. **Safety:** Validation ensures no invalid states (empty treasury, invalid denom, out of bounds multiplier).

3. **Flexibility:** Allows adjusting rewards parameters when needed to optimize rewards distribution.

4. **Consistency:** Ensures rewards module always has valid configuration.
