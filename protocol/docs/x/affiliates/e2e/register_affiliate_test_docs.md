# Test Documentation: Affiliates Register E2E Tests

## Overview

This test file verifies **Affiliate Registration** functionality in the Affiliates module. Affiliates allow users to register referral relationships where a referee (user) can be associated with an affiliate (referrer). The test ensures that:
1. Only the referee can register the affiliate relationship
2. The affiliate cannot register themselves
3. Unrelated addresses cannot register the relationship
4. Signature validation works correctly

---

## Test Function: TestRegisterAffiliateInvalidSigner

### Test Case 1: Success - Valid Signer (Referee)

### Input
- **Referee:** Bob's address
- **Affiliate:** Alice's address
- **Signer:** Bob's private key (referee)
- **Message:** `MsgRegisterAffiliate`
  - Referee: Bob's address
  - Affiliate: Alice's address

### Output
- **CheckTx:** SUCCESS
- **Transaction:** Accepted
- **Relationship:** Affiliate relationship registered

### Why It Runs This Way?

1. **Referee Authorization:** Only the referee can register the affiliate relationship.
2. **Self-Registration:** Referee must sign the transaction themselves.
3. **Relationship Creation:** Creates a referral relationship between referee and affiliate.

---

### Test Case 2: Failure - Invalid Signer (Affiliate)

### Input
- **Referee:** Bob's address
- **Affiliate:** Alice's address
- **Signer:** Alice's private key (affiliate, incorrect signer)
- **Message:** `MsgRegisterAffiliate`
  - Referee: Bob's address
  - Affiliate: Alice's address

### Output
- **CheckTx:** FAIL
- **Error:** "pubKey does not match signer address"
- **Transaction:** Rejected

### Why It Runs This Way?

1. **Authorization:** Only the referee can register the relationship.
2. **Security:** Prevents affiliates from registering themselves.
3. **Signature Validation:** System validates that signer matches referee address.

---

### Test Case 3: Failure - Invalid Signer (Non-Related Address)

### Input
- **Referee:** Bob's address
- **Affiliate:** Alice's address
- **Signer:** Carl's private key (unrelated address, incorrect signer)
- **Message:** `MsgRegisterAffiliate`
  - Referee: Bob's address
  - Affiliate: Alice's address

### Output
- **CheckTx:** FAIL
- **Error:** "pubKey does not match signer address"
- **Transaction:** Rejected

### Why It Runs This Way?

1. **Authorization:** Only the referee can register the relationship.
2. **Security:** Prevents third parties from registering relationships.
3. **Signature Validation:** System validates that signer matches referee address.

---

## Flow Summary

### Affiliate Registration Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. CREATE MESSAGE                                            │
│    - Referee address                                         │
│    - Affiliate address                                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SIGN TRANSACTION                                          │
│    - Referee signs with their private key                   │
│    - Signature included in transaction                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. CHECKTX VALIDATION                                        │
│    - Validate signature matches referee address              │
│    - Check if signer is referee                              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. REGISTRATION                                              │
│    - If valid: Register affiliate relationship              │
│    - If invalid: Reject transaction                          │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Registration States:**
   ```
   Not Registered → Registration Request → Validated → Registered / Rejected
   ```

2. **Signer Validation:**
   ```
   Transaction → Extract Signer → Compare with Referee → Authorized / Unauthorized
   ```

### Key Points

1. **Referee Authorization:**
   - Only the referee can register the affiliate relationship
   - Referee must sign the transaction
   - Signature must match referee address

2. **Security:**
   - Affiliate cannot register themselves
   - Third parties cannot register relationships
   - Signature validation prevents unauthorized registrations

3. **Relationship Creation:**
   - Creates referral relationship between referee and affiliate
   - Relationship can be used for rewards and tracking
   - One-to-one relationship (one referee, one affiliate)

### Design Rationale

1. **User Control:** Referee controls their own affiliate registration.

2. **Security:** Signature validation prevents unauthorized registrations.

3. **Simplicity:** Simple one-to-one relationship model.

4. **Flexibility:** Allows users to choose their affiliate.

