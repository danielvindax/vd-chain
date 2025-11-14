# Test Documentation: Prices Module Governance Proposals

## Overview

This test file verifies governance proposal to **update Market Param** in Prices module. Market param contains information about price feed configuration for a market.

---

## Test Function: TestUpdateMarketParam

### Test Case 1: Success - Update Market Param

### Input
- **Genesis State:**
  - Has market param with ID = 0, Pair = "btc-avdtn", MinPriceChangePpm = 1_000
  - Market exists in market map
- **Proposed Message:**
  - `MsgUpdateMarketParam`:
    - Id: 0
    - Pair: "btc-avdtn" (unchanged)
    - MinPriceChangePpm: 2_002 (changed from 1_000)
    - Authority: gov module

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Market Param:** MinPriceChangePpm updated to 2_002

### Why It Runs This Way?

1. **Market Param Update:** Allows adjusting minimum price change threshold for a market.
2. **MinPriceChangePpm:** This is the threshold (parts per million) to determine when price change is large enough to be considered significant change.
3. **Governance Control:** Only governance has permission to update market params to ensure consistency.

---

### Test Case 2: Failure - Market Param Does Not Exist

### Input
- **Genesis State:** Only has market param with ID = 0
- **Proposed Message:** 
  - `MsgUpdateMarketParam` with ID = 1 (does not exist)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** No change

### Why It Runs This Way?

1. **Existence Check:** Cannot update market param that doesn't exist.
2. **Execution-Time Validation:** Validation occurs when proposal executes, not at submission.
3. **State Protection:** Ensures no partial updates.

---

### Test Case 3: Failure - New Pair Name Does Not Exist in Market Map

### Input
- **Genesis State:** 
  - Has market param with ID = 0, Pair = "btc-avdtn"
  - Market map only has "btc-avdtn"
- **Proposed Message:**
  - `MsgUpdateMarketParam` with Pair = "nonexistent-pair"

### Output
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** No change

### Why It Runs This Way?

1. **Market Map Integration:** Market param must reference a pair that exists in market map.
2. **Consistency:** Ensures market param and market map always sync with each other.
3. **Dependency:** Market param depends on market map configuration.

---

### Test Case 4: Failure - Empty Pair

### Input
- **Proposed Message:**
  - `MsgUpdateMarketParam` with Pair = "" (empty string)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **Required Field:** Pair is the identifier of market, cannot be empty.
2. **Early Validation:** Validation at CheckTx to reject early.
3. **Data Integrity:** Ensures all market params have valid pair name.

---

### Test Case 5: Failure - Invalid Authority

### Input
- **Proposed Message:**
  - `MsgUpdateMarketParam` with Authority = Alice's address (not gov module)

### Output
- **Proposal Submission:** FAIL
- **Proposals:** No proposals created

### Why It Runs This Way?

1. **Authority Check:** Only governance module has permission to update market params.
2. **Security:** Ensures only governance can change price feed configuration.
3. **Early Rejection:** Validation at proposal submission time.

---

## Flow Summary

### Update Market Param Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - Pair not empty                                          │
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
│    - Check market param exists                               │
│    - Check pair exists in market map                         │
│    - Update market param                                     │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Market Param State:**
   ```
   Exists → Updated (UPDATE)
   Does not exist → FAIL
   ```

2. **Market Map Integration:**
   ```
   Pair must exist in market map
   Cannot update with non-existent pair
   ```

### Key Points

1. **Existence Validation:**
   - Market param must exist to update
   - Pair must exist in market map

2. **Authority:**
   - Only governance module has permission to update
   - Validation at both submission and execution time

3. **MinPriceChangePpm:**
   - This is the threshold to determine significant price changes
   - Can be adjusted to fine-tune price update frequency

4. **Market Map Dependency:**
   - Market param depends on market map
   - Ensures consistency between two systems

5. **Atomic Execution:**
   - If validation fails, entire proposal fails
   - State is rolled back to before execution

### Design Rationale

1. **Governance Control:** Only governance has permission to change price feed configuration to ensure decentralization and consistency.

2. **Consistency:** Market param and market map must always sync with each other to ensure price feeds work correctly.

3. **Safety:** Validation ensures no invalid states (empty pair, non-existent market).

4. **Flexibility:** Allows adjusting MinPriceChangePpm to optimize price update frequency.
