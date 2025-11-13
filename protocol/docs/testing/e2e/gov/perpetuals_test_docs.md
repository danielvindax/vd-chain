# Test Documentation: Perpetuals Module Governance Proposals

## Overview

This test file verifies governance proposals related to the **Perpetuals Module**, including:
1. **Update Module Params:** Update parameters of perpetuals module
2. **Update Perpetual Params:** Update parameters of a specific perpetual
3. **Set Liquidity Tier:** Create or update liquidity tier

---

## Test Function: TestUpdatePerpetualsModuleParams

### Test Case 1: Success - Update Module Params

### Input
- **Genesis State:** Perpetuals module has params different from proposal
- **Proposed Message:**
  - `MsgUpdateParams`:
    - FundingRateClampFactorPpm: 123_456
    - PremiumVoteClampFactorPpm: 123_456_789
    - MinNumVotesPerSample: 15
    - Authority: gov module

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Module Params:** Updated with new values

### Why It Runs This Way?

1. **Module-Level Params:** These are parameters that apply to the entire perpetuals module, not a specific perpetual.
2. **Governance Control:** Only governance has permission to update module params.
3. **Validation:** All values must be > 0 to ensure validity.

---

### Test Case 2-4: Failure Cases - Zero Values

### Input
- **Proposed Message:** One of the params = 0:
  - FundingRateClampFactorPpm = 0
  - PremiumVoteClampFactorPpm = 0
  - MinNumVotesPerSample = 0

### Output
- **CheckTx:** FAIL
- **State:** No change

### Why It Runs This Way?

1. **Non-Zero Validation:** All params must be > 0 because they are used in calculations.
2. **Early Rejection:** Validation at CheckTx to reject early, no need to wait for proposal execution.

---

### Test Case 5: Failure - Invalid Authority

### Input
- **Proposed Message:** Authority = perpetuals module (instead of gov module)

### Output
- **Proposal Submission:** FAIL

### Why It Runs This Way?

1. **Authority Check:** Only governance module has permission to update module params.
2. **Security:** Ensures only governance can change module-level settings.

---

## Test Function: TestUpdatePerpetualsParams

### Test Case 1: Success - Update Perpetual Params

### Input
- **Genesis State:** 
  - Has perpetual with ID = 0
  - Has liquidity tier with ID = 123
  - Has market with ID = 4
- **Proposed Message:**
  - `MsgUpdatePerpetualParams`:
    - Id: 0
    - Ticker: "BTC-VDTN" (changed)
    - MarketId: 4
    - DefaultFundingPpm: 500 (changed)
    - LiquidityTier: 123
    - MarketType: unchanged (immutable)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Perpetual Params:** Updated (except MarketType)

### Why It Runs This Way?

1. **Perpetual-Level Update:** Updates params of a specific perpetual.
2. **MarketType Immutable:** MarketType cannot be changed after perpetual is created.
3. **Dependencies:** Must have liquidity tier and market exist first.

---

### Test Case 2: Failure - Empty Ticker

### Input
- **Proposed Message:** Ticker = "" (empty)

### Output
- **CheckTx:** FAIL

### Why It Runs This Way?

1. **Required Field:** Ticker is the identifier of perpetual, cannot be empty.

---

### Test Case 3: Failure - Default Funding PPM Exceeds Maximum

### Input
- **Proposed Message:** DefaultFundingPpm = 1_000_001 (> 1 million)

### Output
- **CheckTx:** FAIL

### Why It Runs This Way?

1. **Boundary Check:** DefaultFundingPpm must be <= 1_000_000 (1 million ppm = 100%).
2. **PPM Format:** PPM (parts per million) has maximum of 1 million.

---

### Test Case 4: Failure - Invalid Authority

### Input
- **Proposed Message:** Authority = perpetuals module (instead of gov)

### Output
- **Proposal Submission:** FAIL

### Why It Runs This Way?

1. **Governance Control:** Only governance has permission to update perpetual params.

---

### Test Case 5: Failure - Liquidity Tier Does Not Exist

### Input
- **Genesis State:** Only has liquidity tier ID = 123
- **Proposed Message:** LiquidityTier = 124 (does not exist)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`

### Why It Runs This Way?

1. **Dependency Check:** Perpetual must reference an existing liquidity tier.
2. **Execution-Time Validation:** Validation occurs when proposal executes, not at submission.

---

### Test Case 6: Failure - Market ID Does Not Exist

### Input
- **Genesis State:** Only has market ID = 4
- **Proposed Message:** MarketId = 5 (does not exist)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`

### Why It Runs This Way?

1. **Market Dependency:** Perpetual must reference an existing market.
2. **Execution Validation:** Check occurs when executing proposal.

---

## Test Function: TestSetLiquidityTier

### Test Case 1: Success - Create New Liquidity Tier

### Input
- **Genesis State:** No liquidity tier ID = 5678
- **Proposed Message:**
  - `MsgSetLiquidityTier`:
    - Id: 5678
    - Name: "Test Tier"
    - InitialMarginPpm: 765_432
    - MaintenanceFractionPpm: 345_678
    - ImpactNotional: 654_321

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Liquidity Tier:** Created new

### Why It Runs This Way?

1. **Create Operation:** When liquidity tier doesn't exist, will create new.
2. **Idempotency:** Can update after creation.

---

### Test Case 2: Success - Update Existing Liquidity Tier

### Input
- **Genesis State:** Has liquidity tier ID = 5678
- **Proposed Message:** `MsgSetLiquidityTier` with same ID

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Liquidity Tier:** Updated

### Why It Runs This Way?

1. **Update Operation:** When liquidity tier already exists, will update instead of creating new.

---

### Test Case 3-5: Failure Cases - Invalid Values

### Input
- **Proposed Message:** One of the cases:
  - InitialMarginPpm = 1_000_001 (> maximum)
  - MaintenanceFractionPpm = 1_000_001 (> maximum)
  - ImpactNotional = 0

### Output
- **CheckTx:** FAIL

### Why It Runs This Way?

1. **Boundary Validation:** 
   - PPM values must be <= 1_000_000
   - ImpactNotional must be > 0 (cannot be = 0)

---

### Test Case 6: Failure - Invalid Authority

### Input
- **Proposed Message:** Authority = perpetuals module (instead of gov)

### Output
- **Proposal Submission:** FAIL

### Why It Runs This Way?

1. **Governance Control:** Only governance has permission to set liquidity tiers.

---

## Flow Summary

### Update Module Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - All params > 0                                          │
│    - Authority = gov module                                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                                │
│    - Validate authority                                    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTION                                        │
│    - Update module params                                    │
│    - Apply to all perpetuals                                │
└─────────────────────────────────────────────────────────────┘
```

### Update Perpetual Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - Ticker not empty                                        │
│    - DefaultFundingPpm <= 1_000_000                         │
│    - Authority = gov module                                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTION                                        │
│    - Check perpetual exists                                  │
│    - Check liquidity tier exists                              │
│    - Check market exists                                     │
│    - Update params (except MarketType)                        │
└─────────────────────────────────────────────────────────────┘
```

### Set Liquidity Tier Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - InitialMarginPpm <= 1_000_000                          │
│    - MaintenanceFractionPpm <= 1_000_000                     │
│    - ImpactNotional > 0                                      │
│    - Authority = gov module                                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTION                                        │
│    - Check if liquidity tier exists                           │
│    - If exists: UPDATE                                       │
│    - If not exists: CREATE                                   │
└─────────────────────────────────────────────────────────────┘
```

### Key Points

1. **Module vs Perpetual Params:**
   - Module params: Apply to entire module
   - Perpetual params: Apply to specific perpetual

2. **Immutable Fields:**
   - MarketType of perpetual cannot be changed after creation

3. **Dependencies:**
   - Perpetual must reference existing liquidity tier and market
   - Validation occurs at execution time

4. **PPM Format:**
   - PPM (parts per million) has maximum = 1_000_000 (100%)
   - All PPM values must be <= 1_000_000

5. **Authority:**
   - All updates must have authority = gov module
   - Validation at both submission and execution time

### Design Rationale

1. **Governance Control:** Only governance has permission to change perpetuals configuration to ensure decentralization.

2. **Safety:** Validation ensures no invalid states (zero values, out of bounds).

3. **Dependency Management:** Ensures perpetuals only reference existing objects.

4. **Flexibility:** Allows updating params when needed to adjust market conditions.
