# Test Documentation: Perpetuals Authorities E2E Tests

## Overview

This test file verifies **Authority Management** functionality in the Perpetuals module. Authorities are addresses that have permission to perform privileged operations in the perpetuals module. The test ensures that:
1. Governance module is recognized as an authority
2. DelayMsg module is recognized as an authority
3. Invalid addresses are not recognized as authorities
4. Authority checks work correctly

---

## Test Function: TestHasAuthority

### Test Case 1: Success - Governance Module is Authority

### Input
- **Authority Address:** Governance module address
  - Address: `authtypes.NewModuleAddress(govtypes.ModuleName)`
- **Check:** `HasAuthority(authorityAddress)`

### Output
- **Result:** `true`
- **Authority:** Governance module is recognized as an authority

### Why It Runs This Way?

1. **Governance Authority:** Governance module needs authority to update perpetual parameters through proposals.
2. **Module Address:** Module addresses are derived from module names.
3. **Permission Check:** System checks if address has authority before allowing privileged operations.

---

### Test Case 2: Success - DelayMsg Module is Authority

### Input
- **Authority Address:** DelayMsg module address
  - Address: `authtypes.NewModuleAddress(delaymsgtypes.ModuleName)`
- **Check:** `HasAuthority(authorityAddress)`

### Output
- **Result:** `true`
- **Authority:** DelayMsg module is recognized as an authority

### Why It Runs This Way?

1. **DelayMsg Authority:** DelayMsg module needs authority to execute delayed messages that update perpetual parameters.
2. **Delayed Updates:** Allows scheduled parameter updates through delayed messages.
3. **Module Integration:** DelayMsg module integrates with perpetuals for parameter updates.

---

### Test Case 3: Failure - Random Invalid Address is Not Authority

### Input
- **Authority Address:** Random invalid address
  - Address: `"random"`
- **Check:** `HasAuthority(authorityAddress)`

### Output
- **Result:** `false`
- **Authority:** Random address is not recognized as an authority

### Why It Runs This Way?

1. **Security:** Only authorized addresses can perform privileged operations.
2. **Validation:** System validates authority before allowing operations.
3. **Access Control:** Prevents unauthorized access to perpetual parameters.

---

## Flow Summary

### Authority Check Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. AUTHORITY REQUEST                                         │
│    - Address provided for authority check                    │
│    - System checks if address is authorized                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. AUTHORITY VALIDATION                                     │
│    - Check if address matches governance module             │
│    - Check if address matches delaymsg module                │
│    - Check if address is in authorized list                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. RESULT                                                    │
│    - Return true if authorized                              │
│    - Return false if not authorized                          │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Authority States:**
   ```
   Address → Authority Check → Authorized / Not Authorized
   ```

2. **Authorized Addresses:**
   ```
   Governance Module Address → Authorized
   DelayMsg Module Address → Authorized
   Other Addresses → Not Authorized
   ```

### Key Points

1. **Module Authorities:**
   - Governance module: Can update parameters through proposals
   - DelayMsg module: Can execute delayed parameter updates
   - Other modules: Not authorized by default

2. **Authority Check:**
   - `HasAuthority(address)` checks if address is authorized
   - Returns boolean result
   - Used before allowing privileged operations

3. **Security:**
   - Only authorized addresses can perform privileged operations
   - Prevents unauthorized parameter updates
   - Ensures system integrity

### Design Rationale

1. **Access Control:** Authority system provides fine-grained access control for perpetual parameters.

2. **Module Integration:** Allows other modules (governance, delaymsg) to interact with perpetuals.

3. **Security:** Prevents unauthorized access to critical system parameters.

4. **Flexibility:** Can be extended to add more authorized addresses if needed.

