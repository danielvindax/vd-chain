# Test Documentation: Fee Tiers Module Governance Proposals

## Overview

This test file verifies governance proposal to **update Perpetual Fee Params** in Fee Tiers module. Fee tiers define the fee structure for perpetual trading based on trading volume and activity.

---

## Test Function: TestUpdateFeeTiersModuleParams

### Test Case 1: Success - Update Perpetual Fee Params

### Input
- **Genesis State:**
  - Fee tiers module has different params from proposal
- **Proposed Message:**
  - `MsgUpdatePerpetualFeeParams`:
    - Tiers:
      - Tier 0:
        - Name: "test_tier_0"
        - MakerFeePpm: 11_000
        - TakerFeePpm: 22_000
        - No volume requirements (first tier)
      - Tier 1:
        - Name: "test_tier_1"
        - AbsoluteVolumeRequirement: 200_000
        - TotalVolumeShareRequirementPpm: 100_000
        - MakerVolumeShareRequirementPpm: 50_000
        - MakerFeePpm: 1_000
        - TakerFeePpm: 2_000
    - Authority: gov module

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Fee Params:** Updated with new tier structure

### Why It Runs This Way?

1. **Fee Tier Structure:** Fee tiers define different fee rates based on trading volume and activity.
2. **First Tier:** The first tier (tier 0) must have no volume requirements - it's the default tier for all traders.
3. **Volume Requirements:** Higher tiers require traders to meet volume thresholds to qualify for lower fees.
4. **Maker/Taker Fees:** Different fees for makers (liquidity providers) and takers (liquidity consumers).
5. **Governance Control:** Only governance has permission to update fee tiers.

---

### Test Case 2: Failure - No Tiers

### Input
- **Proposed Message:**
  - `MsgUpdatePerpetualFeeParams` with Tiers = [] (empty array)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **Required Field:** At least one tier is required for fee structure.
2. **Early Validation:** Validation at CheckTx to reject early.
3. **Data Integrity:** Ensures fee tiers module always has valid tier structure.

---

### Test Case 3: Failure - First Tier Has Non-Zero Volume Requirement

### Input
- **Proposed Message:**
  - `MsgUpdatePerpetualFeeParams` with:
    - Tier 0:
      - AbsoluteVolumeRequirement: 1 (non-zero - WRONG!)
      - MakerFeePpm: 1_000
      - TakerFeePpm: 2_000

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **First Tier Rule:** The first tier (tier 0) must have zero volume requirements because it's the default tier for all traders.
2. **Logical Constraint:** If first tier has volume requirements, new traders cannot qualify for any tier.
3. **Early Rejection:** Validation at CheckTx to prevent invalid configuration.

---

### Test Case 4: Failure - Sum of Lowest Make Fee and Taker Fee is Negative

### Input
- **Proposed Message:**
  - `MsgUpdatePerpetualFeeParams` with:
    - Tier 0:
      - MakerFeePpm: -1_000 (negative - lowest maker fee)
      - TakerFeePpm: 2_000
    - Tier 1:
      - MakerFeePpm: -888
      - TakerFeePpm: 500 (lowest taker fee)
    - Sum of lowest fees: -1_000 + 500 = -500 (negative)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **Fee Validation:** The sum of the lowest maker fee and lowest taker fee across all tiers must be non-negative.
2. **Economic Constraint:** Ensures the fee structure is economically viable - at least one combination of maker/taker fees should be non-negative.
3. **Early Rejection:** Validation at CheckTx to prevent invalid fee structure.

---

### Test Case 5: Failure - Invalid Authority

### Input
- **Proposed Message:**
  - `MsgUpdatePerpetualFeeParams` with Authority = fee tiers module address (instead of gov module)

### Output
- **Proposal Submission:** FAIL
- **Proposals:** No proposals created

### Why It Runs This Way?

1. **Authority Check:** Only governance module has permission to update fee tiers params.
2. **Security:** Ensures only governance can change fee structure.
3. **Early Rejection:** Validation at proposal submission time.

---

## Flow Summary

### Update Perpetual Fee Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - Tiers array not empty                                   │
│    - First tier has zero volume requirements                 │
│    - Sum of lowest maker fee + lowest taker fee >= 0        │
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
│    - Update perpetual fee params                              │
│    - Apply new tier structure                                │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Fee Params Update:**
   ```
   Old Tiers → New Tiers (UPDATE)
   ```

2. **Tier Structure:**
   ```
   Tier 0: No volume requirements (default tier)
   Tier 1+: Volume requirements to qualify
   ```

### Key Points

1. **Tier Structure:**
   - First tier (tier 0) must have zero volume requirements
   - Higher tiers require volume thresholds
   - Each tier has maker and taker fees

2. **Volume Requirements:**
   - AbsoluteVolumeRequirement: Absolute trading volume threshold
   - TotalVolumeShareRequirementPpm: Total volume share (PPM)
   - MakerVolumeShareRequirementPpm: Maker volume share (PPM)

3. **Fee Validation:**
   - Sum of lowest maker fee + lowest taker fee must be >= 0
   - Ensures economically viable fee structure
   - Validation at CheckTx

4. **Authority:**
   - Only governance module has permission to update
   - Validation at proposal submission time

5. **Atomic Execution:**
   - If validation fails, entire proposal fails
   - State is rolled back to before execution

### Design Rationale

1. **Governance Control:** Only governance has permission to change fee structure to ensure decentralization and fairness.

2. **Safety:** Validation ensures no invalid states (empty tiers, invalid first tier, negative fee sum).

3. **Flexibility:** Allows adjusting fee tiers when needed to optimize trading incentives.

4. **Economic Viability:** Fee validation ensures fee structure is economically sustainable.

5. **User Experience:** First tier with no requirements ensures all traders can participate.

