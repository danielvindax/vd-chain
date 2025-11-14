# Tài liệu Test: Authorization Module

## Tổng quan

File test này xác minh chức năng **Authorization (Authz) Module**. Module authz cho phép một tài khoản (granter) cấp quyền cho tài khoản khác (grantee) để thực thi một số message thay mặt họ. Test này đảm bảo rằng:
1. External messages (như `MsgSend`) có thể được cấp và thực thi
2. Internal messages không thể được cấp hoặc thực thi qua authz
3. App-injected messages bị chặn
4. Nested authz messages bị chặn
5. Unsupported messages bị chặn
6. Custom dYdX messages bị chặn

---

## Test Function: TestAuthz

### Test Case 1: Thành công - Alice cấp quyền cho Bob để gửi từ tài khoản của cô ấy

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:**
  - Granter: Alice
  - Grantee: Bob
  - Authorization: Generic authorization cho `MsgSend`
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgSend` từ Alice đến Bob, số tiền: 1 USDC

### Đầu ra
- **CheckTx:** SUCCESS
- **DeliverTx:** SUCCESS
- **Số dư Alice:** Giảm 1 USDC + phí (5 cents)
- **Số dư Bob:** Tăng 1 USDC - phí (5 cents)

### Tại sao chạy theo cách này?

1. **External Messages:** `MsgSend` là external message có thể được cấp qua authz.
2. **Cấp quyền:** Alice cấp quyền cho Bob để gửi token từ tài khoản của cô ấy.
3. **Thực thi:** Bob thực thi thành công thao tác gửi thay mặt Alice.
4. **Thanh toán phí:** Mỗi giao dịch (grant và exec) trả phí riêng biệt.

---

### Test Case 2: Thất bại - Bob cố gắng vote thay mặt Alice mà không có quyền

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** Không có
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgVote` thay mặt Alice

### Đầu ra
- **CheckTx:** SUCCESS
- **DeliverTx:** FAIL với lỗi `ErrNoAuthorizationFound`

### Tại sao chạy theo cách này?

1. **Không có quyền:** Bob không có quyền vote thay mặt Alice.
2. **CheckTx Pass:** CheckTx không validate authz permissions, chỉ validate format message.
3. **DeliverTx Fail:** Authz keeper validate permissions trong DeliverTx và từ chối.

---

### Test Case 3: Thất bại - Cấp quyền cho Internal Messages không cho phép thực thi

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:**
  - Granter: Alice
  - Grantee: Bob
  - Authorization: Generic authorization cho `MsgUpdateParams` (internal message)
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgUpdateParams` với authority = gov module

### Đầu ra
- **CheckTx:** SUCCESS
- **DeliverTx:** FAIL với lỗi `ErrNoAuthorizationFound`

### Tại sao chạy theo cách này?

1. **Internal Messages:** Internal messages (như `MsgUpdateParams`) không thể được thực thi qua authz.
2. **Bảo mật:** Điều này ngăn chặn thực thi trái phép các thao tác đặc quyền.
3. **Grant thành công:** Grant được chấp nhận, nhưng thực thi bị chặn.

---

### Test Case 4: Thất bại - Bob cố gắng cập nhật Gov Params (Authority = Gov)

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** Không có
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgUpdateParams` với authority = gov module

### Đầu ra
- **CheckTx:** SUCCESS
- **DeliverTx:** FAIL với lỗi `ErrNoAuthorizationFound`

### Tại sao chạy theo cách này?

1. **Không có quyền:** Bob không có quyền thực thi messages thay mặt gov module.
2. **Internal Message:** Ngay cả khi được cấp, internal messages không thể được thực thi qua authz.

---

### Test Case 5: Thất bại - Bob cố gắng cập nhật Gov Params (Authority = Bob)

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** Không có
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgUpdateParams` với authority = Bob

### Đầu ra
- **CheckTx:** SUCCESS
- **DeliverTx:** FAIL với lỗi `ErrInvalidSigner`

### Tại sao chạy theo cách này?

1. **Invalid Authority:** Bob không có trong danh sách authorized signers để tạo CLOB pairs.
2. **Authority Check:** Message tự thất bại vì Bob không có authority cần thiết.
3. **Lỗi khác:** Thất bại với `ErrInvalidSigner` thay vì `ErrNoAuthorizationFound`.

---

### Test Case 6: Thất bại - Bob cố gắng đề xuất Operations (App Injected)

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** Không có
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgProposedOperations`

### Đầu ra
- **CheckTx:** FAIL với lỗi `ErrInvalidRequest`
- **DeliverTx:** Không đạt được

### Tại sao chạy theo cách này?

1. **App-Injected Messages:** `MsgProposedOperations` được inject bởi app, không được submit bởi users.
2. **Ante Handler:** Ante handler từ chối các messages này tại CheckTx.
3. **Bảo mật:** Ngăn chặn users submit app-internal messages.

---

### Test Case 7: Thất bại - Double Nested Authz Message

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** Không có
- **MsgExec:**
  - Grantee: Bob
  - Message: Một `MsgExec` khác (nested)

### Đầu ra
- **CheckTx:** FAIL với lỗi `ErrInvalidRequest`
- **DeliverTx:** Không đạt được

### Tại sao chạy theo cách này?

1. **Nested Authz:** Authz messages không thể được nested (wrap một `MsgExec` khác).
2. **Ante Handler:** Ante handler từ chối nested authz messages tại CheckTx.
3. **Bảo mật:** Ngăn chặn chuỗi authorization phức tạp.

---

### Test Case 8: Thất bại - Unsupported Transaction Type

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** Không có
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgUpdateParams` từ ICA controller module (unsupported)

### Đầu ra
- **CheckTx:** FAIL với lỗi `ErrInvalidRequest`
- **DeliverTx:** Không đạt được

### Tại sao chạy theo cách này?

1. **Unsupported Messages:** Một số loại message không được hỗ trợ trong authz.
2. **Ante Handler:** Ante handler duy trì whitelist của các messages được hỗ trợ.
3. **Bảo mật:** Ngăn chặn thực thi các thao tác nguy hiểm hoặc không được hỗ trợ.

---

### Test Case 9: Thất bại - Bob wrap dYdX Custom Messages

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: 100,000 USD
  - Bob_Num0: 100,000 USD
- **MsgGrant:** Không có
- **MsgExec:**
  - Grantee: Bob
  - Message: `MsgPlaceOrder` (dYdX custom message)

### Đầu ra
- **CheckTx:** FAIL với lỗi `ErrInvalidRequest`
- **DeliverTx:** Không đạt được

### Tại sao chạy theo cách này?

1. **Custom Messages:** dYdX custom messages (như `MsgPlaceOrder`) không được hỗ trợ trong authz.
2. **Ante Handler:** Ante handler chặn custom dYdX messages.
3. **Bảo mật:** Ngăn chặn thao tác trading trái phép qua authz.

---

## Tóm tắt Flow

### Successful Authz Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. GRANTER CẤP QUYỀN                                       │
│    - Alice cấp quyền cho Bob để thực thi MsgSend            │
│    - CheckTx: SUCCESS                                        │
│    - DeliverTx: SUCCESS                                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. GRANTEE THỰC THI MESSAGE                                 │
│    - Bob thực thi MsgSend thay mặt Alice                    │
│    - CheckTx: SUCCESS                                        │
│    - DeliverTx: SUCCESS                                       │
│    - Transfer được thực thi                                  │
└─────────────────────────────────────────────────────────────┘
```

### Failed Authz Flow (No Permission)

```
┌─────────────────────────────────────────────────────────────┐
│ 1. GRANTEE CỐ GẮNG THỰC THI                                 │
│    - Bob cố gắng thực thi message mà không có quyền         │
│    - CheckTx: SUCCESS (format hợp lệ)                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. DELIVERTX VALIDATION                                     │
│    - Authz keeper kiểm tra authorization                    │
│    - Không tìm thấy authorization                            │
│    - DeliverTx: FAIL với ErrNoAuthorizationFound          │
└─────────────────────────────────────────────────────────────┘
```

### Failed Authz Flow (Blocked at CheckTx)

```
┌─────────────────────────────────────────────────────────────┐
│ 1. GRANTEE CỐ GẮNG THỰC THI                                 │
│    - Bob cố gắng thực thi loại message bị chặn             │
│    - Ante handler validate loại message                     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CHECKTX REJECTION                                        │
│    - Loại message không được phép trong authz               │
│    - CheckTx: FAIL với ErrInvalidRequest                   │
│    - DeliverTx: Không đạt được                                │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Authorization State:**
   ```
   Không có Authorization → Authorization được cấp → Authorization được sử dụng
   ```

2. **Message Execution:**
   ```
   CheckTx: Validate format message
   DeliverTx: Validate authorization và thực thi
   ```

### Điểm quan trọng

1. **External vs Internal Messages:**
   - External messages (như `MsgSend`) có thể được cấp và thực thi
   - Internal messages (như `MsgUpdateParams`) không thể được thực thi qua authz

2. **Validation Points:**
   - CheckTx: Validate format và loại message
   - DeliverTx: Validate authorization permissions

3. **Blocked Message Types:**
   - App-injected messages (`MsgProposedOperations`)
   - Nested authz messages (double `MsgExec`)
   - Unsupported messages (ICA controller messages)
   - Custom dYdX messages (`MsgPlaceOrder`)

4. **Fee Payment:**
   - Mỗi giao dịch (grant và exec) trả phí riêng biệt
   - Granter trả phí cho giao dịch grant
   - Grantee trả phí cho giao dịch exec

5. **Bảo mật:**
   - Chỉ external, whitelisted messages có thể được thực thi qua authz
   - Internal và privileged operations bị chặn
   - Ngăn chặn truy cập trái phép vào các thao tác nhạy cảm

### Lý do thiết kế

1. **Bảo mật:** Authz bị giới hạn ở các thao tác external an toàn để ngăn chặn truy cập trái phép vào các chức năng đặc quyền.

2. **Linh hoạt:** Cho phép users ủy quyền một số thao tác (như gửi token) cho các tài khoản khác.

3. **Validation:** Nhiều lớp validation (ante handler, CheckTx, DeliverTx) đảm bảo chỉ các thao tác được phép mới có thể được thực thi.

4. **Blocking:** App-injected và custom messages bị chặn để ngăn chặn lạm dụng và duy trì tính toàn vẹn hệ thống.
