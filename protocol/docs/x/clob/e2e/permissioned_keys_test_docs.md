# Test Documentation: Permissioned Keys E2E Tests

## Overview

This test file verifies **Permissioned Keys (Smart Account Authenticators)** functionality in the CLOB module. Smart accounts can have multiple authenticators that control which operations can be performed. The test ensures that:
1. Orders can be placed with specific authenticators
2. Authenticators must be enabled (smart account feature)
3. Authenticators must exist and not be removed
4. Authenticators validate message types and signatures
5. Composite authenticators (AllOf, AnyOf) work correctly

---

## Test Function: TestPlaceOrder_PermissionedKeys_Failures

### Test Case 1: Failure - Smart Account Not Enabled

### Input
- **Smart Account:** Not enabled
- **Order:**
  - Bob places order to buy 5 at price 40
  - Authenticators: [0] specified
- **Transaction:** Signed with Bob's private key

### Output
- **CheckTx:** FAIL
- **Error Code:** `ErrSmartAccountNotActive`
- **Error Message:** "Smart account is not active"

### Why It Runs This Way?

1. **Feature Flag:** Smart account feature must be enabled.
2. **Authenticators Invalid:** Cannot use authenticators if feature disabled.
3. **Early Rejection:** CheckTx rejects transaction immediately.

---

### Test Case 2: Failure - Authenticator Not Found

### Input
- **Smart Account:** Enabled
- **Order:**
  - Bob places order to buy 5 at price 40
  - Authenticators: [0] specified
- **State:** No authenticators added to Bob's account

### Output
- **CheckTx:** FAIL
- **Error Code:** `ErrAuthenticatorNotFound`
- **Error Message:** "Authenticator not found"

### Why It Runs This Way?

1. **No Authenticators:** Bob's account has no authenticators added.
2. **Invalid Reference:** Authenticator ID 0 doesn't exist.
3. **Validation:** System checks authenticator exists before use.

---

### Test Case 3: Failure - Authenticator Was Removed

### Input
- **Smart Account:** Enabled
- **Block 2:**
  - Add authenticator: Bob adds AllOf authenticator (ID 0)
- **Block 4:**
  - Remove authenticator: Bob removes authenticator ID 0
- **Block 5:**
  - Place order: Bob places order with authenticator [0]

### Output
- **Add Authenticator:** SUCCESS
- **Remove Authenticator:** SUCCESS
- **Place Order:** FAIL with error `ErrAuthenticatorNotFound`

### Why It Runs This Way?

1. **Authenticator Removed:** Authenticator was removed in block 4.
2. **No Longer Exists:** Authenticator ID 0 no longer exists.
3. **Cannot Use:** Cannot use removed authenticator.

---

### Test Case 4: Success - Authenticator Validates Message Type

### Input
- **Smart Account:** Enabled
- **Authenticator:** AllOf with:
  - SignatureVerification (Bob's key)
  - MessageFilter (allows only `/cosmos.bank.v1beta1.MsgSend`)
- **Order:** Bob places order to buy 5 at price 40
- **Authenticators:** [0] specified

### Output
- **CheckTx:** FAIL
- **Error:** Authenticator doesn't allow CLOB order message type

### Why It Runs This Way?

1. **Message Filter:** Authenticator only allows `MsgSend`.
2. **Order Message:** Order uses `MsgPlaceOrder` (different type).
3. **Filter Rejection:** Message filter rejects non-allowed message types.

---

## Flow Summary

### Permissioned Key Validation Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. CHECK SMART ACCOUNT ENABLED                               │
│    - Verify smart account feature is enabled                  │
│    - Reject if feature disabled                              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. VALIDATE AUTHENTICATORS                                   │
│    - Check authenticator IDs exist                           │
│    - Verify authenticators not removed                       │
│    - Reject if invalid                                        │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. EXECUTE AUTHENTICATORS                                    │
│    - For each authenticator:                                 │
│      * SignatureVerification: Verify signature               │
│      * MessageFilter: Check message type                     │
│      * Composite (AllOf/AnyOf): Evaluate children            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. AUTHENTICATOR RESULT                                      │
│    - AllOf: All children must pass                           │
│    - AnyOf: At least one child must pass                     │
│    - Reject if authentication fails                          │
└─────────────────────────────────────────────────────────────┘
```

### Authenticator Types

1. **SignatureVerification:**
   - Verifies transaction signature
   - Uses specified public key
   - Must match signer

2. **MessageFilter:**
   - Filters allowed message types
   - Only specified message types allowed
   - Rejects other message types

3. **AllOf (Composite):**
   - All child authenticators must pass
   - Logical AND operation
   - All conditions must be met

4. **AnyOf (Composite):**
   - At least one child authenticator must pass
   - Logical OR operation
   - Any condition can be met

### Key Points

1. **Smart Account Feature:**
   - Must be enabled to use authenticators
   - Feature flag controls availability
   - Disabled by default

2. **Authenticator Management:**
   - Authenticators can be added
   - Authenticators can be removed
   - Removed authenticators cannot be used

3. **Message Type Filtering:**
   - Authenticators can restrict message types
   - Only allowed message types pass filter
   - Provides fine-grained access control

4. **Composite Authenticators:**
   - AllOf: All conditions must pass
   - AnyOf: At least one condition must pass
   - Can nest multiple levels

5. **Validation Timing:**
   - Checked at CheckTx
   - Early rejection for invalid authenticators
   - Clear error messages

6. **Security:**
   - Multiple authenticators provide layered security
   - Message filtering prevents unauthorized operations
   - Signature verification ensures authorization

### Design Rationale

1. **Access Control:** Authenticators provide fine-grained access control.

2. **Security:** Multiple authenticators add security layers.

3. **Flexibility:** Composite authenticators allow complex policies.

4. **User Control:** Users can manage their authenticators.

5. **Message Filtering:** Prevents unauthorized message types.

