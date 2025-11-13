# Test Documentation: Vest Entry Governance Proposals

## Overview

This test file verifies the **Vest Entries** (entries managing vesting tokens) management functionality through governance proposals. The test includes:
1. **Set Vest Entry:** Create or update vest entry
2. **Delete Vest Entry:** Delete vest entry

Vest entries manage vesting tokens from treasury account to vester account according to a time schedule.

---

## Test Function: TestSetVestEntry_Success

### Test Case 1: Success - Create a New Vest Entry

### Input
- **Genesis State:** No vest entries
- **Proposed Message:**
  - `MsgSetVestEntry`:
    - VesterAccount: "random_vester"
    - TreasuryAccount: "random_treasury"
    - Denom: "avdtn"
    - StartTime: 2023-10-02 00:00:00 UTC
    - EndTime: 2024-10-01 00:00:00 UTC
    - Authority: gov module address

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entry:** Created in state with above information
- **State:** Vest entry can be queried and used

### Why It Runs This Way?

1. **Create Operation:** When vest entry doesn't exist, `MsgSetVestEntry` will create new.
2. **Governance Authority:** Only governance module has permission to create vest entries.
3. **Time Range:** StartTime < EndTime ensures valid vesting period.

---

### Test Case 2: Success - Update an Existing Vest Entry

### Input
- **Genesis State:** Has vest entry with VesterAccount = "random_vester"
- **Proposed Message:**
  - `MsgSetVestEntry` with same VesterAccount but different information:
    - TreasuryAccount: "random_treasury" (new)
    - Denom: "avdtn" (new)
    - StartTime: 2023-10-02 (new)
    - EndTime: 2024-10-01 (new)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entry:** Updated with new information
- **Old Entry:** Overwritten by new entry

### Why It Runs This Way?

1. **Update Operation:** When vest entry already exists, `MsgSetVestEntry` will update instead of creating new.
2. **Idempotency:** Can update multiple times with same VesterAccount.
3. **Flexibility:** Allows changing treasury, denom, or time range.

---

### Test Case 3: Success - Create Two New Vest Entries

### Input
- **Genesis State:** No vest entries
- **Proposed Messages:**
  1. `MsgSetVestEntry` for "random_vester"
  2. `MsgSetVestEntry` for "random_vester_2"

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entries:** Both entries created successfully
- **State:** Both can be queried independently

### Why It Runs This Way?

1. **Batch Creation:** One proposal can create multiple vest entries at once.
2. **Independent Entries:** Each entry is independent, doesn't affect others.
3. **Efficiency:** Allows setting up multiple vesting schedules in one proposal.

---

### Test Case 4: Success - Create and Then Update a Vest Entry

### Input
- **Genesis State:** No vest entries
- **Proposed Messages:**
  1. `MsgSetVestEntry` creates entry for "random_vester"
  2. `MsgSetVestEntry` updates entry for "random_vester" (same proposal)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entry:** Only last entry (from message 2) is saved
- **State:** Entry updated with information from message 2

### Why It Runs This Way?

1. **Sequential Execution:** Messages are executed in order within proposal.
2. **Last Write Wins:** Message 2 overwrites message 1 because same VesterAccount.
3. **Use Case:** Allows adjusting information within the same proposal.

---

### Test Case 5: Success - Update a Vest Entry Twice

### Input
- **Genesis State:** Has vest entry with VesterAccount = "random_vester"
- **Proposed Messages:**
  1. `MsgSetVestEntry` updates entry first time
  2. `MsgSetVestEntry` updates entry second time (same proposal)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entry:** Entry updated twice, final result is from message 2

### Why It Runs This Way?

1. **Multiple Updates:** Can update same entry multiple times in one proposal.
2. **Final State:** Only final state is kept.
3. **Flexibility:** Allows fine-tuning vesting parameters.

---

## Test Function: TestSetVestEntry_Failure

### Test Case 1: Failure - Vester Account is Empty

### Input
- **Proposed Message:**
  - `MsgSetVestEntry` with VesterAccount = "" (empty string)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted
- **State:** No change

### Why It Runs This Way?

1. **Validation:** VesterAccount cannot be empty because it's the key to identify entry.
2. **Early Rejection:** Validation occurs at CheckTx, no need to wait for proposal execution.
3. **Data Integrity:** Ensures all entries have valid identifier.

---

### Test Case 2: Failure - Treasury Account is Empty

### Input
- **Proposed Message:**
  - `MsgSetVestEntry` with TreasuryAccount = "" (empty string)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **Required Field:** TreasuryAccount is required because it's the source of tokens to vest.
2. **Validation:** Empty treasury account is invalid for vesting operation.

---

### Test Case 3: Failure - Start Time After End Time

### Input
- **Proposed Message:**
  - `MsgSetVestEntry` with:
    - StartTime: 2024-10-01 (after)
    - EndTime: 2023-10-02 (before)

### Output
- **CheckTx:** FAIL
- **Proposal:** Not submitted

### Why It Runs This Way?

1. **Logical Validation:** StartTime must be < EndTime for valid vesting period.
2. **Time Logic:** Cannot have vesting period with start time after end time.

---

### Test Case 4: Failure - Invalid Authority

### Input
- **Proposed Message:**
  - `MsgSetVestEntry` with Authority = Bob's address (not gov module)

### Output
- **Proposal Submission:** FAIL
- **Proposals:** No proposals created

### Why It Runs This Way?

1. **Authority Check:** Only governance module has permission to set vest entries.
2. **Security:** Ensures only governance can manage vesting schedules.
3. **Early Rejection:** Validation at proposal submission time.

---

### Test Case 5: Failure - One Message Fails Causes Rollback

### Input
- **Proposed Messages:**
  1. `MsgSetVestEntry` valid for "random_vester"
  2. `MsgSetVestEntry` with invalid authority for "random_vester"

### Output
- **Proposal Submission:** FAIL (due to message 2)
- **State:** No entries created (including message 1)

### Why It Runs This Way?

1. **Atomic Execution:** If one message fails, entire proposal fails.
2. **Rollback:** State is rolled back to before proposal execution.
3. **Consistency:** Ensures no partial state.

---

## Test Function: TestDeleteVestEntry_Success

### Test Case 1: Success - Delete One Vest Entry

### Input
- **Genesis State:** Has vest entry with VesterAccount = "random_vester"
- **Proposed Message:**
  - `MsgDeleteVestEntry` with VesterAccount = "random_vester"

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entry:** Removed from state
- **Query:** Cannot query entry (returns error)

### Why It Runs This Way?

1. **Delete Operation:** `MsgDeleteVestEntry` removes entry from state.
2. **Cleanup:** After deletion, entry no longer exists in state.
3. **Governance Control:** Only governance has permission to delete entries.

---

### Test Case 2: Success - Delete Two Vest Entries

### Input
- **Genesis State:** Has 2 vest entries
- **Proposed Messages:**
  1. `MsgDeleteVestEntry` for entry 1
  2. `MsgDeleteVestEntry` for entry 2

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entries:** Both deleted

### Why It Runs This Way?

1. **Batch Deletion:** One proposal can delete multiple entries.
2. **Independent Operations:** Each deletion is independent.

---

## Test Function: TestDeleteVestEntry_Failure

### Test Case 1: Failure - Vest Entry Does Not Exist

### Input
- **Genesis State:** No vest entries
- **Proposed Message:**
  - `MsgDeleteVestEntry` for "random_vester" (does not exist)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** No change

### Why It Runs This Way?

1. **Existence Check:** Cannot delete entry that doesn't exist.
2. **Error Handling:** Proposal fails but doesn't crash system.
3. **State Protection:** State unaffected when deleting non-existent entry.

---

### Test Case 2: Failure - Delete the Same Vest Entry Twice

### Input
- **Genesis State:** Has 1 vest entry
- **Proposed Messages:**
  1. `MsgDeleteVestEntry` for entry (succeeds)
  2. `MsgDeleteVestEntry` for same entry (fails because already deleted)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Entry still deleted (from message 1), but proposal fails due to message 2

### Why It Runs This Way?

1. **Sequential Execution:** Message 1 succeeds, message 2 fails.
2. **Atomic Failure:** When message 2 fails, entire proposal fails and state rolls back.
3. **State Consistency:** Entry is not deleted because proposal failed.

---

### Test Case 3: Failure - Second Entry to Delete Does Not Exist

### Input
- **Genesis State:** Has 1 vest entry ("random_vester")
- **Proposed Messages:**
  1. `MsgDeleteVestEntry` for "random_vester" (exists)
  2. `MsgDeleteVestEntry` for "random_vester_2" (does not exist)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** No change (rollback)

### Why It Runs This Way?

1. **Partial Failure:** Message 1 succeeds but message 2 fails.
2. **Rollback:** Entire proposal fails and state rolls back.
3. **Consistency:** Ensures no partial deletion.

---

### Test Case 4: Failure - Invalid Authority

### Input
- **Genesis State:** Has vest entry
- **Proposed Message:**
  - `MsgDeleteVestEntry` with Authority = Bob's address (not gov)

### Output
- **Proposal Submission:** FAIL
- **State:** No change

### Why It Runs This Way?

1. **Authority Validation:** Only governance module has permission to delete entries.
2. **Security:** Ensures only governance can manage vesting.

---

## Flow Summary

### Set Vest Entry Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - VesterAccount not empty                                 │
│    - TreasuryAccount not empty                               │
│    - StartTime < EndTime                                     │
│    - Authority = gov module                                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                                │
│    - Validate authority                                      │
│    - Validate message format                                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTION                                        │
│    - Check if vest entry exists                               │
│    - If exists: UPDATE                                       │
│    - If not exists: CREATE                                   │
│    - Store in state                                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. VERIFY STATE                                              │
│    - Query vest entry                                        │
│    - Verify information is correct                           │
└─────────────────────────────────────────────────────────────┘
```

### Delete Vest Entry Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - VesterAccount not empty                                 │
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
│    - Check if vest entry exists                               │
│    - If exists: DELETE                                       │
│    - If not exists: FAIL                                     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. VERIFY STATE                                              │
│    - Query vest entry → Should return error                  │
│    - Verify entry no longer in state                          │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Vest Entry State:**
   ```
   Does not exist → Exists (CREATE)
   Exists → Updated (UPDATE)
   Exists → Does not exist (DELETE)
   ```

2. **Proposal Status:**
   ```
   SUBMITTED → PASSED (if all messages succeed)
   SUBMITTED → FAILED (if any message fails)
   ```

### Key Points

1. **Idempotency:** 
   - `SetVestEntry` can be called multiple times with same VesterAccount
   - `DeleteVestEntry` only succeeds if entry exists

2. **Atomic Execution:**
   - If one message fails, entire proposal fails
   - State is rolled back to before proposal execution

3. **Authority:**
   - Only governance module has permission to set/delete vest entries
   - Validation occurs at both submission and execution time

4. **Validation:**
   - VesterAccount and TreasuryAccount cannot be empty
   - StartTime must be < EndTime
   - Entry must exist to delete

5. **Batch Operations:**
   - One proposal can contain multiple messages
   - All messages must succeed for proposal to pass

### Design Rationale

1. **Governance Control:** Only governance has permission to manage vesting schedules to ensure decentralization and security.

2. **Flexibility:** Allows updating vest entries to adjust vesting schedules when needed.

3. **Safety:** Validation ensures data integrity and prevents invalid states.

4. **Atomic Operations:** Ensures state consistency - no partial updates.

5. **Error Handling:** Clear error messages and proper rollback when errors occur.
