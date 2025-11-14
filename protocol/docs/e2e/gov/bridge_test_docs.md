# Test Documentation: Bridge Module Governance Proposals

## Overview

This test file verifies governance proposals related to the **Bridge Module**, including:
1. **Update Event Params:** Update event parameters for bridge operations
2. **Update Propose Params:** Update propose parameters for bridge proposals
3. **Update Safety Params:** Update safety parameters for bridge safety controls

---

## Test Function: TestUpdateEventParams

### Test Case 1: Success - Update Event Params

### Input
- **Genesis State:**
  - Event params:
    - Denom: "testdenom"
    - EthChainId: 123
    - EthAddress: "0x0123"
- **Proposed Message:**
  - `MsgUpdateEventParams`:
    - Denom: "advtnt" (changed)
    - EthChainId: 1 (changed)
    - EthAddress: "0xabcd" (changed)
    - Authority: gov module

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Event Params:** Updated with new values

### Why It Runs This Way?

1. **Event Configuration:** Event params control how bridge events are processed and validated.
2. **Denom:** The denomination of tokens used in bridge operations.
3. **EthChainId:** Ethereum chain ID for cross-chain bridge operations.
4. **EthAddress:** Ethereum address for bridge contract or operations.
5. **Governance Control:** Only governance has permission to update event params.

---

### Test Case 2: Failure - Empty ETH Address

### Input
- **Proposed Message:**
  - `MsgUpdateEventParams` with EthAddress = "" (empty string)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **Required Field:** EthAddress is required for bridge operations, cannot be empty.
2. **Early Validation:** Validation at CheckTx to reject early.
3. **Data Integrity:** Ensures bridge always has valid Ethereum address.

---

### Test Case 3: Failure - Invalid Authority

### Input
- **Proposed Message:**
  - `MsgUpdateEventParams` with Authority = Bob's address (not gov module)

### Output
- **Proposal Submission:** FAIL
- **Proposals:** No proposals created

### Why It Runs This Way?

1. **Authority Check:** Only governance module has permission to update event params.
2. **Security:** Ensures only governance can change bridge event configuration.
3. **Early Rejection:** Validation at proposal submission time.

---

## Test Function: TestUpdateProposeParams

### Test Case 1: Success - Update Propose Params

### Input
- **Genesis State:**
  - Propose params:
    - MaxBridgesPerBlock: 10
    - ProposeDelayDuration: 1 minute
    - SkipRatePpm: 800_000
    - SkipIfBlockDelayedByDuration: 1 minute
- **Proposed Message:**
  - `MsgUpdateProposeParams`:
    - MaxBridgesPerBlock: 7 (changed)
    - ProposeDelayDuration: 1 second (changed)
    - SkipRatePpm: 700_007 (changed)
    - SkipIfBlockDelayedByDuration: 1 second (changed)
    - Authority: gov module

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Propose Params:** Updated with new values

### Why It Runs This Way?

1. **Propose Configuration:** Propose params control how bridge proposals are submitted and processed.
2. **MaxBridgesPerBlock:** Maximum number of bridges that can be proposed per block.
3. **ProposeDelayDuration:** Delay duration before proposal can be submitted.
4. **SkipRatePpm:** Rate (parts per million) to skip proposals under certain conditions.
5. **SkipIfBlockDelayedByDuration:** Duration threshold to skip proposals if block is delayed.

---

### Test Case 2: Failure - Negative Propose Delay Duration

### Input
- **Proposed Message:**
  - `MsgUpdateProposeParams` with ProposeDelayDuration = -1 second (negative)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **Non-Negative Validation:** Duration cannot be negative.
2. **Early Rejection:** Validation at CheckTx to reject early.
3. **Logical Constraint:** Negative duration doesn't make sense for delay.

---

### Test Case 3: Failure - Negative Skip If Block Delayed By Duration

### Input
- **Proposed Message:**
  - `MsgUpdateProposeParams` with SkipIfBlockDelayedByDuration = -1 second (negative)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **Non-Negative Validation:** Duration cannot be negative.
2. **Early Rejection:** Validation at CheckTx.
3. **Logical Constraint:** Negative duration is invalid.

---

### Test Case 4: Failure - Skip Rate PPM Out of Bounds

### Input
- **Proposed Message:**
  - `MsgUpdateProposeParams` with SkipRatePpm = 1_000_001 (> 1 million)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **Boundary Check:** SkipRatePpm must be <= 1_000_000 (1 million ppm = 100%).
2. **PPM Format:** PPM (parts per million) has maximum of 1 million.
3. **Logical Constraint:** Rate cannot exceed 100%.

---

### Test Case 5: Failure - Invalid Authority

### Input
- **Proposed Message:**
  - `MsgUpdateProposeParams` with Authority = Alice's address (not gov module)

### Output
- **Proposal Submission:** FAIL
- **Proposals:** No proposals created

### Why It Runs This Way?

1. **Authority Check:** Only governance module has permission to update propose params.
2. **Security:** Ensures only governance can change bridge proposal configuration.

---

## Test Function: TestUpdateSafetyParams

### Test Case 1: Success - Update Safety Params

### Input
- **Genesis State:**
  - Safety params:
    - IsDisabled: false
    - DelayBlocks: 10
- **Proposed Message:**
  - `MsgUpdateSafetyParams`:
    - IsDisabled: true (changed)
    - DelayBlocks: 5 (changed)
    - Authority: gov module

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Safety Params:** Updated with new values

### Why It Runs This Way?

1. **Safety Configuration:** Safety params control safety mechanisms for bridge operations.
2. **IsDisabled:** Flag to enable/disable bridge safety checks.
3. **DelayBlocks:** Number of blocks to delay before executing bridge operations.
4. **Governance Control:** Only governance has permission to update safety params.

---

### Test Case 2: Failure - Invalid Authority

### Input
- **Proposed Message:**
  - `MsgUpdateSafetyParams` with Authority = Alice's address (not gov module)

### Output
- **Proposal Submission:** FAIL
- **Proposals:** No proposals created

### Why It Runs This Way?

1. **Authority Check:** Only governance module has permission to update safety params.
2. **Security:** Ensures only governance can change bridge safety configuration.
3. **Early Rejection:** Validation at proposal submission time.

---

## Flow Summary

### Update Event Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - EthAddress not empty                                    │
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
│    - Update event params                                      │
│    - Apply new configuration                                 │
└─────────────────────────────────────────────────────────────┘
```

### Update Propose Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - ProposeDelayDuration >= 0                               │
│    - SkipIfBlockDelayedByDuration >= 0                      │
│    - SkipRatePpm <= 1_000_000                                │
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
│    - Update propose params                                    │
│    - Apply new configuration                                 │
└─────────────────────────────────────────────────────────────┘
```

### Update Safety Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
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
│    - Update safety params                                     │
│    - Apply new configuration                                 │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Event Params Update:**
   ```
   Old Params → New Params (UPDATE)
   ```

2. **Propose Params Update:**
   ```
   Old Params → New Params (UPDATE)
   ```

3. **Safety Params Update:**
   ```
   Old Params → New Params (UPDATE)
   ```

### Key Points

1. **Event Params:**
   - Denom: Token denomination for bridge operations
   - EthChainId: Ethereum chain ID for cross-chain operations
   - EthAddress: Must not be empty
   - Validation at CheckTx

2. **Propose Params:**
   - MaxBridgesPerBlock: Maximum bridges per block
   - ProposeDelayDuration: Must be >= 0
   - SkipRatePpm: Must be <= 1_000_000
   - SkipIfBlockDelayedByDuration: Must be >= 0
   - Validation at CheckTx

3. **Safety Params:**
   - IsDisabled: Enable/disable safety checks
   - DelayBlocks: Number of blocks to delay
   - No CheckTx validation (only authority check)

4. **Authority:**
   - Only governance module has permission to update all params
   - Validation at proposal submission time

5. **Atomic Execution:**
   - If validation fails, entire proposal fails
   - State is rolled back to before execution

### Design Rationale

1. **Governance Control:** Only governance has permission to change bridge configuration to ensure decentralization and security.

2. **Safety:** Validation ensures no invalid states (empty address, negative duration, out of bounds rate).

3. **Flexibility:** Allows adjusting bridge parameters when needed to optimize bridge operations.

4. **Consistency:** Ensures bridge module always has valid configuration.

