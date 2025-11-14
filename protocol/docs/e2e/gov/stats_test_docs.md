# Test Documentation: Stats Module Governance Proposals

## Overview

This test file verifies governance proposal to **update Stats Module Params**. Stats module manages statistics tracking and window duration for various metrics.

---

## Test Function: TestUpdateParams

### Test Case 1: Success - Update Stats Module Params

### Input
- **Genesis State:**
  - Stats module has params:
    - WindowDuration: different from proposal
- **Proposed Message:**
  - `MsgUpdateParams`:
    - WindowDuration: 1 hour (changed)
    - Authority: gov module

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Module Params:** Updated with new WindowDuration

### Why It Runs This Way?

1. **Stats Configuration:** WindowDuration controls the time window for statistics calculations.
2. **Time Window:** This parameter defines how long statistics are tracked and aggregated.
3. **Governance Control:** Only governance has permission to update stats params.
4. **Flexibility:** Allows adjusting the statistics window when needed.

---

### Test Case 2: Failure - Invalid Authority

### Input
- **Proposed Message:**
  - `MsgUpdateParams` with Authority = stats module address (instead of gov module)

### Output
- **Proposal Submission:** FAIL
- **Proposals:** No proposals created
- **State:** No change

### Why It Runs This Way?

1. **Authority Check:** Only governance module has permission to update stats params.
2. **Security:** Ensures only governance can change statistics configuration.
3. **Early Rejection:** Validation at proposal submission time.

---

## Flow Summary

### Update Stats Module Params Process

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
│    - Update stats module params                               │
│    - Apply new WindowDuration                                 │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Params Update:**
   ```
   Old Params → New Params (UPDATE)
   ```

2. **WindowDuration:**
   ```
   Old Duration → New Duration (UPDATE)
   ```

### Key Points

1. **WindowDuration:**
   - Defines the time window for statistics tracking
   - Used for aggregating and calculating statistics
   - Can be adjusted to optimize statistics collection

2. **Authority:**
   - Only governance module has permission to update
   - Validation at proposal submission time

3. **Atomic Execution:**
   - If validation fails, entire proposal fails
   - State is rolled back to before execution

4. **Simplicity:**
   - Stats module has minimal params (only WindowDuration)
   - Simple update process

### Design Rationale

1. **Governance Control:** Only governance has permission to change statistics configuration to ensure decentralization.

2. **Flexibility:** Allows adjusting the statistics window when needed to optimize metrics collection.

3. **Simplicity:** Minimal parameters keep the module simple and focused.

4. **Consistency:** Ensures stats module always has valid configuration.

