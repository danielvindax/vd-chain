# Tài liệu Test: Permissioned Keys E2E Tests

## Tổng quan

File test này xác minh chức năng **Permissioned Keys (Smart Account Authenticators)** trong CLOB module. Smart accounts có thể có nhiều authenticators kiểm soát các operations nào có thể được thực hiện. Test đảm bảo rằng:
1. Orders có thể được đặt với specific authenticators
2. Authenticators phải được enable (smart account feature)
3. Authenticators phải tồn tại và không bị xóa
4. Authenticators validate message types và signatures
5. Composite authenticators (AllOf, AnyOf) hoạt động đúng

---

## Test Function: TestPlaceOrder_PermissionedKeys_Failures

### Test Case 1: Thất bại - Smart Account Không được Enable

### Đầu vào
- **Smart Account:** Không được enable
- **Order:**
  - Bob đặt order để mua 5 ở giá 40
  - Authenticators: [0] được chỉ định
- **Transaction:** Được ký với private key của Bob

### Đầu ra
- **CheckTx:** FAIL
- **Error Code:** `ErrSmartAccountNotActive`
- **Error Message:** "Smart account is not active"

### Tại sao chạy theo cách này?

1. **Feature Flag:** Smart account feature phải được enable.
2. **Authenticators Invalid:** Không thể sử dụng authenticators nếu feature disabled.
3. **Early Rejection:** CheckTx từ chối transaction ngay lập tức.

---

### Test Case 2: Thất bại - Authenticator Không Tìm thấy

### Đầu vào
- **Smart Account:** Được enable
- **Order:**
  - Bob đặt order để mua 5 ở giá 40
  - Authenticators: [0] được chỉ định
- **State:** Không có authenticators được thêm vào tài khoản của Bob

### Đầu ra
- **CheckTx:** FAIL
- **Error Code:** `ErrAuthenticatorNotFound`
- **Error Message:** "Authenticator not found"

### Tại sao chạy theo cách này?

1. **No Authenticators:** Tài khoản của Bob không có authenticators được thêm.
2. **Invalid Reference:** Authenticator ID 0 không tồn tại.
3. **Validation:** Hệ thống kiểm tra authenticator tồn tại trước khi sử dụng.

---

### Test Case 3: Thất bại - Authenticator Đã bị Xóa

### Đầu vào
- **Smart Account:** Được enable
- **Block 2:**
  - Add authenticator: Bob thêm AllOf authenticator (ID 0)
- **Block 4:**
  - Remove authenticator: Bob xóa authenticator ID 0
- **Block 5:**
  - Place order: Bob đặt order với authenticator [0]

### Đầu ra
- **Add Authenticator:** SUCCESS
- **Remove Authenticator:** SUCCESS
- **Place Order:** FAIL với lỗi `ErrAuthenticatorNotFound`

### Tại sao chạy theo cách này?

1. **Authenticator Removed:** Authenticator được xóa trong block 4.
2. **No Longer Exists:** Authenticator ID 0 không còn tồn tại.
3. **Cannot Use:** Không thể sử dụng removed authenticator.

---

### Test Case 4: Thành công - Authenticator Validate Message Type

### Đầu vào
- **Smart Account:** Được enable
- **Authenticator:** AllOf với:
  - SignatureVerification (key của Bob)
  - MessageFilter (chỉ cho phép `/cosmos.bank.v1beta1.MsgSend`)
- **Order:** Bob đặt order để mua 5 ở giá 40
- **Authenticators:** [0] được chỉ định

### Đầu ra
- **CheckTx:** FAIL
- **Error:** Authenticator không cho phép CLOB order message type

### Tại sao chạy theo cách này?

1. **Message Filter:** Authenticator chỉ cho phép `MsgSend`.
2. **Order Message:** Order sử dụng `MsgPlaceOrder` (loại khác).
3. **Filter Rejection:** Message filter từ chối non-allowed message types.

---

## Tóm tắt Flow

### Permissioned Key Validation Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. KIỂM TRA SMART ACCOUNT ENABLED                           │
│    - Xác minh smart account feature được enable            │
│    - Từ chối nếu feature disabled                           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. VALIDATE AUTHENTICATORS                                  │
│    - Kiểm tra authenticator IDs tồn tại                     │
│    - Xác minh authenticators không bị xóa                   │
│    - Từ chối nếu invalid                                    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. THỰC THI AUTHENTICATORS                                  │
│    - Cho mỗi authenticator:                                 │
│      * SignatureVerification: Xác minh signature           │
│      * MessageFilter: Kiểm tra message type                 │
│      * Composite (AllOf/AnyOf): Đánh giá children           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. AUTHENTICATOR RESULT                                     │
│    - AllOf: Tất cả children phải pass                      │
│    - AnyOf: Ít nhất một child phải pass                    │
│    - Từ chối nếu authentication thất bại                    │
└─────────────────────────────────────────────────────────────┘
```

### Authenticator Types

1. **SignatureVerification:**
   - Xác minh transaction signature
   - Sử dụng specified public key
   - Phải khớp với signer

2. **MessageFilter:**
   - Lọc allowed message types
   - Chỉ specified message types được phép
   - Từ chối các message types khác

3. **AllOf (Composite):**
   - Tất cả child authenticators phải pass
   - Logical AND operation
   - Tất cả conditions phải được đáp ứng

4. **AnyOf (Composite):**
   - Ít nhất một child authenticator phải pass
   - Logical OR operation
   - Bất kỳ condition nào có thể được đáp ứng

### Điểm quan trọng

1. **Smart Account Feature:**
   - Phải được enable để sử dụng authenticators
   - Feature flag kiểm soát availability
   - Disabled theo mặc định

2. **Authenticator Management:**
   - Authenticators có thể được thêm
   - Authenticators có thể được xóa
   - Removed authenticators không thể được sử dụng

3. **Message Type Filtering:**
   - Authenticators có thể giới hạn message types
   - Chỉ allowed message types pass filter
   - Cung cấp fine-grained access control

4. **Composite Authenticators:**
   - AllOf: Tất cả conditions phải pass
   - AnyOf: Ít nhất một condition phải pass
   - Có thể nest nhiều levels

5. **Validation Timing:**
   - Được kiểm tra tại CheckTx
   - Early rejection cho invalid authenticators
   - Error messages rõ ràng

6. **Bảo mật:**
   - Nhiều authenticators cung cấp layered security
   - Message filtering ngăn chặn unauthorized operations
   - Signature verification đảm bảo authorization

### Lý do thiết kế

1. **Access Control:** Authenticators cung cấp fine-grained access control.

2. **Bảo mật:** Nhiều authenticators thêm security layers.

3. **Flexibility:** Composite authenticators cho phép complex policies.

4. **User Control:** Users có thể quản lý authenticators của họ.

5. **Message Filtering:** Ngăn chặn unauthorized message types.
