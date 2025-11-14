# Test Documentation: Authorization Module

## Overview

This test file verifies the **Authorization (Authz) Module** functionality. The authz module allows one account (granter) to grant permissions to another account (grantee) to execute certain messages on their behalf. This test ensures that:
1. External messages (like `MsgSend`) can be granted and executed
2. Internal messages cannot be granted or executed via authz
3. App-injected messages are blocked
4. Nested authz messages are blocked
5. Unsupported messages are blocked
6. Custom dYdX messages are blocked

---

## Test Function: TestAuthz

### Test Case 1: Success - Alice Grants Permission to Bob to Send from Her Account

### Input
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:**
  - Granter: Alice
  - Grantee: Bob
  - Authorization: Generic authorization for `MsgSend`
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgSend` from Alice to Bob, amount: 1 USDC

### Output
- **CheckTx:** SUCCESS
- **DeliverTx:** SUCCESS
- **Alice Balance:** Decreased by 1 USDC + fees (5 cents)
- **Bob Balance:** Increased by 1 USDC - fees (5 cents)

### Why It Runs This Way?

1. **External Messages:** `MsgSend` is an external message that can be granted via authz.
2. **Permission Grant:** Alice grants Bob permission to send tokens from her account.
3. **Execution:** Bob successfully executes the send operation on behalf of Alice.
4. **Fee Payment:** Each transaction (grant and exec) pays fees separately.

---

### Test Case 2: Failure - Bob Tries to Vote on Behalf of Alice Without Permission

### Input
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** None
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgVote` on behalf of Alice

### Output
- **CheckTx:** SUCCESS
- **DeliverTx:** FAIL with error `ErrNoAuthorizationFound`

### Why It Runs This Way?

1. **No Permission:** Bob doesn't have permission to vote on behalf of Alice.
2. **CheckTx Passes:** CheckTx doesn't validate authz permissions, only message format.
3. **DeliverTx Fails:** Authz keeper validates permissions during DeliverTx and rejects.

---

### Test Case 3: Failure - Granting Permissions for Internal Messages Doesn't Allow Execution

### Input
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:**
  - Granter: Alice
  - Grantee: Bob
  - Authorization: Generic authorization for `MsgUpdateParams` (internal message)
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgUpdateParams` with authority = gov module

### Output
- **CheckTx:** SUCCESS
- **DeliverTx:** FAIL with error `ErrNoAuthorizationFound`

### Why It Runs This Way?

1. **Internal Messages:** Internal messages (like `MsgUpdateParams`) cannot be executed via authz.
2. **Security:** This prevents unauthorized execution of privileged operations.
3. **Grant Succeeds:** Grant is accepted, but execution is blocked.

---

### Test Case 4: Failure - Bob Tries to Update Gov Params (Authority = Gov)

### Input
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** None
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgUpdateParams` with authority = gov module

### Output
- **CheckTx:** SUCCESS
- **DeliverTx:** FAIL with error `ErrNoAuthorizationFound`

### Why It Runs This Way?

1. **No Permission:** Bob doesn't have permission to execute messages on behalf of gov module.
2. **Internal Message:** Even if granted, internal messages cannot be executed via authz.

---

### Test Case 5: Failure - Bob Tries to Update Gov Params (Authority = Bob)

### Input
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** None
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgUpdateParams` with authority = Bob

### Output
- **CheckTx:** SUCCESS
- **DeliverTx:** FAIL with error `ErrInvalidSigner`

### Why It Runs This Way?

1. **Invalid Authority:** Bob is not in the list of authorized signers for creating CLOB pairs.
2. **Authority Check:** The message itself fails because Bob doesn't have the required authority.
3. **Different Error:** This fails with `ErrInvalidSigner` instead of `ErrNoAuthorizationFound`.

---

### Test Case 6: Failure - Bob Tries to Propose Operations (App Injected)

### Input
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** None
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgProposedOperations`

### Output
- **CheckTx:** FAIL with error `ErrInvalidRequest`
- **DeliverTx:** Not reached

### Why It Runs This Way?

1. **App-Injected Messages:** `MsgProposedOperations` is injected by the app, not submitted by users.
2. **Ante Handler:** The ante handler rejects these messages at CheckTx.
3. **Security:** Prevents users from submitting app-internal messages.

---

### Test Case 7: Failure - Double Nested Authz Message

### Input
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** None
- **MsgExec:**
  - Grantee: Bob
  - Message: Another `MsgExec` (nested)

### Output
- **CheckTx:** FAIL with error `ErrInvalidRequest`
- **DeliverTx:** Not reached

### Why It Runs This Way?

1. **Nested Authz:** Authz messages cannot be nested (wrapping another `MsgExec`).
2. **Ante Handler:** The ante handler rejects nested authz messages at CheckTx.
3. **Security:** Prevents complex nested authorization chains.

---

### Test Case 8: Failure - Unsupported Transaction Type

### Input
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** None
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgUpdateParams` from ICA controller module (unsupported)

### Output
- **CheckTx:** FAIL with error `ErrInvalidRequest`
- **DeliverTx:** Not reached

### Why It Runs This Way?

1. **Unsupported Messages:** Some message types are not supported in authz.
2. **Ante Handler:** The ante handler maintains a whitelist of supported messages.
3. **Security:** Prevents execution of potentially dangerous or unsupported operations.

---

### Test Case 9: Failure - Bob Wraps dYdX Custom Messages

### Input
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** None
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgPlaceOrder` (dYdX custom message)

### Output
- **CheckTx:** FAIL with error `ErrInvalidRequest`
- **DeliverTx:** Not reached

### Why It Runs This Way?

1. **Custom Messages:** dYdX custom messages (like `MsgPlaceOrder`) are not supported in authz.
2. **Ante Handler:** The ante handler blocks custom dYdX messages.
3. **Security:** Prevents unauthorized trading operations via authz.

---

## Flow Summary

### Successful Authz Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. GRANTER GRANTS PERMISSION                                │
│    - Alice grants Bob permission to execute MsgSend         │
│    - CheckTx: SUCCESS                                        │
│    - DeliverTx: SUCCESS                                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. GRANTEE EXECUTES MESSAGE                                  │
│    - Bob executes MsgSend on behalf of Alice                 │
│    - CheckTx: SUCCESS                                        │
│    - DeliverTx: SUCCESS                                       │
│    - Transfer executed                                        │
└─────────────────────────────────────────────────────────────┘
```

### Failed Authz Flow (No Permission)

```
┌─────────────────────────────────────────────────────────────┐
│ 1. GRANTEE ATTEMPTS EXECUTION                                │
│    - Bob tries to execute message without permission        │
│    - CheckTx: SUCCESS (format valid)                        │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. DELIVERTX VALIDATION                                      │
│    - Authz keeper checks for authorization                  │
│    - No authorization found                                  │
│    - DeliverTx: FAIL with ErrNoAuthorizationFound          │
└─────────────────────────────────────────────────────────────┘
```

### Failed Authz Flow (Blocked at CheckTx)

```
┌─────────────────────────────────────────────────────────────┐
│ 1. GRANTEE ATTEMPTS EXECUTION                                │
│    - Bob tries to execute blocked message type              │
│    - Ante handler validates message type                     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CHECKTX REJECTION                                         │
│    - Message type not allowed in authz                       │
│    - CheckTx: FAIL with ErrInvalidRequest                   │
│    - DeliverTx: Not reached                                  │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Authorization State:**
   ```
   No Authorization → Authorization Granted → Authorization Used
   ```

2. **Message Execution:**
   ```
   CheckTx: Validates message format
   DeliverTx: Validates authorization and executes
   ```

### Key Points

1. **External vs Internal Messages:**
   - External messages (like `MsgSend`) can be granted and executed
   - Internal messages (like `MsgUpdateParams`) cannot be executed via authz

2. **Validation Points:**
   - CheckTx: Validates message format and type
   - DeliverTx: Validates authorization permissions

3. **Blocked Message Types:**
   - App-injected messages (`MsgProposedOperations`)
   - Nested authz messages (double `MsgExec`)
   - Unsupported messages (ICA controller messages)
   - Custom dYdX messages (`MsgPlaceOrder`)

4. **Fee Payment:**
   - Each transaction (grant and exec) pays fees separately
   - Granter pays fees for grant transaction
   - Grantee pays fees for exec transaction

5. **Security:**
   - Only external, whitelisted messages can be executed via authz
   - Internal and privileged operations are blocked
   - Prevents unauthorized access to sensitive operations

### Design Rationale

1. **Security:** Authz is limited to safe, external operations to prevent unauthorized access to privileged functions.

2. **Flexibility:** Allows users to delegate certain operations (like sending tokens) to other accounts.

3. **Validation:** Multiple validation layers (ante handler, CheckTx, DeliverTx) ensure only allowed operations can be executed.

4. **Blocking:** App-injected and custom messages are blocked to prevent abuse and maintain system integrity.

