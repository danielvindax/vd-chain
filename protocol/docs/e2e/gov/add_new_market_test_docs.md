# Test Documentation: Add New Market Proposal

## Overview

This test file verifies the **add new market** functionality through governance proposals. The test ensures that when creating a new market, the system will:
1. Create Oracle Market (price feed)
2. Create Perpetual contract
3. Create CLOB Pair with INITIALIZING status
4. Use DelayMessage to transition CLOB Pair to ACTIVE after a number of blocks
5. Enable market in market map

---

## Test Case 1: Success with 4 Standard Messages (Delay Blocks = 10)

### Input
- **Proposed Messages:**
  1. `MsgCreateOracleMarket`: Create market param with ID = 1001
  2. `MsgCreatePerpetual`: Create perpetual with ID = 1001
  3. `MsgCreateClobPair`: Create CLOB pair with ID = 1001, status = INITIALIZING
  4. `MsgDelayMessage`: Delay message to update CLOB pair to ACTIVE after 10 blocks
- **Market Map:** Market initially disabled in market map

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Market Param:** Created with ID = 1001
- **Market Price:** Initialized with price = 0
- **Perpetual:** Created with ID = 1001
- **ClobPair:** 
  - Initially: status = INITIALIZING
  - After 10 blocks: status = ACTIVE
- **Market Map:** Market enabled after CLOB pair transitions to ACTIVE

### Why It Runs This Way?

1. **Message Order is Critical:** Messages must be executed in order:
   - Oracle Market first (needed for price feed)
   - Perpetual next (needed for CLOB pair)
   - CLOB Pair after (depends on perpetual)
   - DelayMessage last (to activate CLOB pair)

2. **Delay Blocks = 10:** CLOB pair is not activated immediately but must wait 10 blocks to ensure:
   - All dependencies are fully set up
   - Oracle has time to update price
   - System has time to validate state

3. **Market Map Integration:** Market must be enabled in market map to allow trading, which only happens after CLOB pair transitions to ACTIVE.

---

## Test Case 2: Success with Delay Blocks = 1

### Input
- **Proposed Messages:** Same as Test Case 1
- **Delay Blocks:** 1 (instead of 10)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **ClobPair:** Transitions to ACTIVE after 1 block

### Why It Runs This Way?

1. **Minimum Delay:** This test ensures delay blocks can be 1 (minimum), not necessarily 10.
2. **Fast Activation:** Some cases may need to activate market faster, delay = 1 allows this.

---

## Test Case 3: Success with Delay Blocks = 0

### Input
- **Proposed Messages:** Same as Test Case 1
- **Delay Blocks:** 0 (no delay)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **ClobPair:** No delay, but still needs to be activated through delay message mechanism

### Why It Runs This Way?

1. **Zero Delay:** This test ensures delay = 0 is also supported, meaning message can be executed immediately in the next block.
2. **Edge Case:** This is an edge case to ensure system handles delay = 0 correctly.

---

## Test Case 4: Success with Delayed UpdateClobPair Message Failure

### Input
- **Proposed Messages:**
  1. `MsgCreateOracleMarket`: ID = 1001
  2. `MsgCreatePerpetual`: ID = 1001
  3. `MsgCreateClobPair`: ID = 1001
  4. `MsgDelayMessage`: Contains `MsgUpdateClobPair` with ClobPairId = 9999 (does not exist)
- **Delay Blocks:** 10

### Output
- **Proposal Status:** `PROPOSAL_STATUS_PASSED` (proposal still passes because other messages succeeded)
- **ClobPair:** Still at INITIALIZING status (not updated because delayed message failed)
- **Market Map:** Market still disabled

### Why It Runs This Way?

1. **Delayed Message Failure:** When delayed message fails, it doesn't fail the entire proposal because proposal has already been executed successfully.
2. **Partial Success:** Other messages (create market, perpetual, clob pair) still succeed, only delayed update message fails.
3. **State Consistency:** ClobPair remains at INITIALIZING, unaffected by failed delayed message.

---

## Test Case 5: Fail - Incorrectly Ordered Messages

### Input
- **Proposed Messages (Wrong Order):**
  1. `MsgCreateOracleMarket`: ID = 1001
  2. `MsgCreateClobPair`: ID = 1001 (before creating perpetual - WRONG!)
  3. `MsgCreatePerpetual`: ID = 1001
  4. `MsgDelayMessage`: Update CLOB pair

### Output
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Nothing created (full rollback)

### Why It Runs This Way?

1. **Dependency Order:** CLOB Pair depends on Perpetual, so Perpetual must be created first.
2. **Atomic Execution:** Proposal execution is atomic - if one message fails, entire proposal fails and state is rolled back.
3. **Validation:** System validates dependencies and rejects if order is wrong.

---

## Test Case 6: Fail - Existing Objects

### Input
- **Proposed Messages:**
  1. `MsgCreateOracleMarket`: ID = 5 (already exists in genesis)
  2. `MsgCreatePerpetual`: ID = 5 (already exists)
  3. `MsgCreateClobPair`: ID = 5 (already exists)
  4. `MsgDelayMessage`: Update CLOB pair

### Output
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Nothing new created

### Why It Runs This Way?

1. **Idempotency:** Cannot create objects with IDs that already exist.
2. **Error Handling:** System detects conflict and rejects proposal.
3. **State Protection:** Ensures no duplicate IDs in the system.

---

## Test Case 7: Fail - Invalid Signer (Proposal Submission)

### Input
- **Proposed Messages:**
  1. `MsgCreateOracleMarket`: Authority = CLOB module address (WRONG! Must be gov module)
  2. `MsgCreatePerpetual`: Authority = gov module (correct)
  3. `MsgCreateClobPair`: Authority = gov module (correct)
  4. `MsgDelayMessage`: Authority = gov module (correct)

### Output
- **Proposal Submission:** FAIL (cannot submit)
- **Proposals:** No proposals created

### Why It Runs This Way?

1. **Authority Validation:** Each message must have correct authority:
   - `MsgCreateOracleMarket` must have authority = gov module
   - Validation occurs at proposal submission time
2. **Early Rejection:** Proposal is rejected immediately when submitted, no need to wait for execution.
3. **Security:** Ensures only correct authority can create objects.

---

## Test Case 8: Fail - Invalid Signer on MsgDelayMessage

### Input
- **Proposed Messages:**
  1. `MsgCreateOracleMarket`: Authority = gov module (correct)
  2. `MsgCreatePerpetual`: Authority = gov module (correct)
  3. `MsgCreateClobPair`: Authority = gov module (correct)
  4. `MsgDelayMessage`: 
     - Authority = gov module (correct)
     - But wrapped message (`MsgUpdateClobPair`) has authority = gov module (WRONG! Must be delaymsg module)

### Output
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Market, Perpetual, CLOB Pair created but proposal fails when executing delayed message

### Why It Runs This Way?

1. **Nested Authority:** `MsgDelayMessage` contains another message (`MsgUpdateClobPair`), and that message also has its own authority.
2. **Delayed Validation:** Authority of wrapped message is only validated when delayed message is executed, not when proposal is submitted.
3. **Partial State:** Previous messages executed successfully, but delayed message failure causes proposal to fail.

---

## Flow Summary

### Add New Market Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. INITIALIZE GENESIS STATE                                  │
│    - Market map with market disabled                         │
│    - No market/perpetual/clob pair with new ID               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                                │
│    - Validate authority of all messages                      │
│    - Validate message order                                  │
│    - Validate no duplicate IDs                                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTION (If Submit Successful)                 │
│    a. Execute MsgCreateOracleMarket                          │
│       → Create MarketParam and MarketPrice (price = 0)       │
│    b. Execute MsgCreatePerpetual                              │
│       → Create Perpetual contract                            │
│    c. Execute MsgCreateClobPair                               │
│       → Create CLOB Pair with INITIALIZING status            │
│    d. Execute MsgDelayMessage                                │
│       → Schedule MsgUpdateClobPair to execute after N blocks  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. DELAYED MESSAGE EXECUTION                                 │
│    - After N blocks (0, 1, or 10), delayed message executes │
│    - Validate authority of wrapped message                   │
│    - Update CLOB Pair status: INITIALIZING → ACTIVE          │
│    - Enable market in market map                             │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. TRADING ENABLED                                           │
│    - After CLOB Pair = ACTIVE and market enabled             │
│    - Users can place orders (but need oracle price > 0)       │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **CLOB Pair Status Transition:**
   ```
   Does not exist → INITIALIZING → ACTIVE
   ```

2. **Market Map State:**
   ```
   Disabled → Enabled (when CLOB Pair = ACTIVE)
   ```

3. **Oracle Price:**
   ```
   Does not exist → 0 (initialized) → Actual price (from oracle)
   ```

### Key Points

1. **Message Order:** Message order is CRITICAL:
   - Oracle Market → Perpetual → CLOB Pair → DelayMessage
   - Wrong order will cause proposal to fail

2. **Authority:**
   - Proposal messages: Authority = gov module
   - Delayed UpdateClobPair: Authority = delaymsg module
   - Validation occurs at both submission and execution time

3. **Delay Blocks:**
   - Can be 0, 1, or any number
   - Allows system time to setup before activation

4. **Atomic Execution:**
   - If one message fails, entire proposal fails
   - State is rolled back to before proposal execution

5. **Market Map Integration:**
   - Market must be enabled in market map for trading
   - Only enabled after CLOB Pair = ACTIVE

### Design Rationale

1. **Safety:** Delay mechanism ensures market is not activated immediately, allowing validation and setup.

2. **Dependency Management:** Message order ensures dependencies are created correctly.

3. **Flexibility:** Delay blocks can be adjusted as needed.

4. **Error Handling:** Atomic execution ensures state consistency - no partial state.

5. **Integration:** Market map integration ensures market can be discovered and traded.
