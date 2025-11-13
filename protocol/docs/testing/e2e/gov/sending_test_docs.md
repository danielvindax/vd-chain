# Test Documentation: Sending Module Governance Proposals

## Overview

This test file verifies governance proposal to **send tokens from module account to user account or another module account**. This is the mechanism for governance to distribute tokens from treasury or module accounts.

---

## Test Function: TestSendFromModuleToAccount

### Test Case 1: Success - Send from Module to User Account

### Input
- **Genesis State:**
  - Community Treasury module has balance: 200 avdtn
- **Proposed Message:**
  - `MsgSendFromModuleToAccount`:
    - SenderModuleName: "community_treasury"
    - Recipient: Alice's address
    - Coin: 123 avdtn
    - Authority: gov module

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Module Balance:** 200 - 123 = 77 avdtn
- **Alice Balance:** Increased by 123 avdtn

### Why It Runs This Way?

1. **Module-to-Account Transfer:** Allows governance to distribute tokens from module accounts (like treasury) to user accounts.
2. **Use Cases:** 
   - Airdrops
   - Rewards distribution
   - Treasury disbursements
3. **Governance Control:** Only governance has permission to perform transfers from module accounts.

---

### Test Case 2: Success - Send from Module to Module Account

### Input
- **Genesis State:**
  - Community Treasury module has balance: 123 avdtn
- **Proposed Message:**
  - `MsgSendFromModuleToAccount`:
    - SenderModuleName: "community_treasury"
    - Recipient: Community Vester module address
    - Coin: 123 avdtn
    - Authority: gov module

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Treasury Balance:** 0 avdtn (all transferred)
- **Vester Balance:** Increased by 123 avdtn

### Why It Runs This Way?

1. **Module-to-Module Transfer:** Allows transferring tokens between module accounts.
2. **Use Cases:**
   - Funding vesting contracts
   - Rebalancing module accounts
   - Treasury management
3. **Flexibility:** Supports both user accounts and module accounts as recipients.

---

### Test Case 3: Failure - Insufficient Balance

### Input
- **Genesis State:**
  - Community Treasury module has balance: 123 avdtn
- **Proposed Message:**
  - `MsgSendFromModuleToAccount`:
    - Coin: 124 avdtn (more than balance)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** No change (balances remain unchanged)

### Why It Runs This Way?

1. **Balance Check:** Module account must have sufficient balance to transfer.
2. **Execution-Time Validation:** Validation occurs when proposal executes.
3. **State Protection:** Ensures no negative balances.

---

### Test Case 4: Failure - Invalid Authority

### Input
- **Proposed Message:**
  - `MsgSendFromModuleToAccount` with Authority = sending module address (instead of gov module)

### Output
- **Proposal Submission:** FAIL
- **Proposals:** No proposals created

### Why It Runs This Way?

1. **Authority Check:** Only governance module has permission to send from module accounts.
2. **Security:** Ensures only governance can control module account transfers.
3. **Early Rejection:** Validation at proposal submission time.

---

## Flow Summary

### Send From Module To Account Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - Authority = gov module                                  │
│    - SenderModuleName valid                                 │
│    - Recipient address valid                                │
│    - Coin amount > 0                                         │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                                │
│    - Validate authority                                      │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTION                                        │
│    - Check sender module balance >= coin amount             │
│    - Transfer coins from sender module                       │
│    - Transfer coins to recipient                             │
│    - Update balances                                         │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Balance Changes:**
   ```
   Sender Module: Balance -= Amount
   Recipient: Balance += Amount
   ```

2. **Validation Points:**
   ```
   Submission: Authority check
   Execution: Balance check
   ```

### Key Points

1. **Balance Validation:**
   - Sender module must have sufficient balance
   - Validation occurs at execution time
   - If insufficient, proposal fails and state rolls back

2. **Authority:**
   - Only governance module has permission
   - Validation at proposal submission time

3. **Recipient Types:**
   - Supports both user accounts and module accounts
   - Recipient address must be valid

4. **Atomic Execution:**
   - If balance insufficient, entire proposal fails
   - State is rolled back to before execution

5. **Use Cases:**
   - Treasury disbursements
   - Rewards distribution
   - Module account rebalancing
   - Funding vesting contracts

### Design Rationale

1. **Governance Control:** Only governance has permission to transfer from module accounts to ensure decentralization and security.

2. **Safety:** Balance validation ensures no negative balances or insufficient funds.

3. **Flexibility:** Supports both user and module accounts as recipients to support many use cases.

4. **Transparency:** All transfers from module accounts go through governance proposals, ensuring transparency.
