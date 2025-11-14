# Tài liệu Test: Vest Entry Governance Proposals

## Tổng quan

File test này xác minh chức năng quản lý **Vest Entries** (entries quản lý vesting tokens) thông qua governance proposals. Test bao gồm:
1. **Set Vest Entry:** Tạo hoặc cập nhật vest entry
2. **Delete Vest Entry:** Xóa vest entry

Vest entries quản lý vesting tokens từ treasury account đến vester account theo lịch trình thời gian.

---

## Test Function: TestSetVestEntry_Success

### Test Case 1: Thành công - Tạo Vest Entry Mới

### Đầu vào
- **Genesis State:** Không có vest entries
- **Proposed Message:**
  - `MsgSetVestEntry`:
    - VesterAccount: "random_vester"
    - TreasuryAccount: "random_treasury"
    - Denom: "avdtn"
    - StartTime: 2023-10-02 00:00:00 UTC
    - EndTime: 2024-10-01 00:00:00 UTC
    - Authority: gov module address

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entry:** Được tạo trong state với thông tin trên
- **State:** Vest entry có thể được query và sử dụng

### Tại sao chạy theo cách này?

1. **Create Operation:** Khi vest entry không tồn tại, `MsgSetVestEntry` sẽ tạo mới.
2. **Governance Authority:** Chỉ governance module có quyền tạo vest entries.
3. **Time Range:** StartTime < EndTime đảm bảo vesting period hợp lệ.

---

### Test Case 2: Thành công - Cập nhật Vest Entry Hiện có

### Đầu vào
- **Genesis State:** Có vest entry với VesterAccount = "random_vester"
- **Proposed Message:**
  - `MsgSetVestEntry` với cùng VesterAccount nhưng thông tin khác:
    - TreasuryAccount: "random_treasury" (mới)
    - Denom: "avdtn" (mới)
    - StartTime: 2023-10-02 (mới)
    - EndTime: 2024-10-01 (mới)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entry:** Được cập nhật với thông tin mới
- **Old Entry:** Bị ghi đè bởi entry mới

### Tại sao chạy theo cách này?

1. **Update Operation:** Khi vest entry đã tồn tại, `MsgSetVestEntry` sẽ cập nhật thay vì tạo mới.
2. **Idempotency:** Có thể cập nhật nhiều lần với cùng VesterAccount.
3. **Flexibility:** Cho phép thay đổi treasury, denom, hoặc time range.

---

### Test Case 3: Thành công - Tạo Hai Vest Entries Mới

### Đầu vào
- **Genesis State:** Không có vest entries
- **Proposed Messages:**
  1. `MsgSetVestEntry` cho "random_vester"
  2. `MsgSetVestEntry` cho "random_vester_2"

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entries:** Cả hai entries được tạo thành công
- **State:** Cả hai có thể được query độc lập

### Tại sao chạy theo cách này?

1. **Batch Creation:** Một proposal có thể tạo nhiều vest entries cùng lúc.
2. **Independent Entries:** Mỗi entry độc lập, không ảnh hưởng đến nhau.
3. **Efficiency:** Cho phép thiết lập nhiều vesting schedules trong một proposal.

---

### Test Case 4: Thành công - Tạo và Sau đó Cập nhật Vest Entry

### Đầu vào
- **Genesis State:** Không có vest entries
- **Proposed Messages:**
  1. `MsgSetVestEntry` tạo entry cho "random_vester"
  2. `MsgSetVestEntry` cập nhật entry cho "random_vester" (cùng proposal)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entry:** Chỉ entry cuối (từ message 2) được lưu
- **State:** Entry được cập nhật với thông tin từ message 2

### Tại sao chạy theo cách này?

1. **Sequential Execution:** Messages được thực thi theo thứ tự trong proposal.
2. **Last Write Wins:** Message 2 ghi đè message 1 vì cùng VesterAccount.
3. **Use Case:** Cho phép điều chỉnh thông tin trong cùng proposal.

---

### Test Case 5: Thành công - Cập nhật Vest Entry Hai Lần

### Đầu vào
- **Genesis State:** Có vest entry với VesterAccount = "random_vester"
- **Proposed Messages:**
  1. `MsgSetVestEntry` cập nhật entry lần đầu
  2. `MsgSetVestEntry` cập nhật entry lần hai (cùng proposal)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entry:** Entry được cập nhật hai lần, kết quả cuối từ message 2

### Tại sao chạy theo cách này?

1. **Multiple Updates:** Có thể cập nhật cùng entry nhiều lần trong một proposal.
2. **Final State:** Chỉ trạng thái cuối được giữ lại.
3. **Flexibility:** Cho phép tinh chỉnh vesting parameters.

---

## Test Function: TestSetVestEntry_Failure

### Test Case 1: Thất bại - Vester Account Rỗng

### Đầu vào
- **Proposed Message:**
  - `MsgSetVestEntry` với VesterAccount = "" (chuỗi rỗng)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit
- **State:** Không thay đổi

### Tại sao chạy theo cách này?

1. **Validation:** VesterAccount không thể rỗng vì đây là key để xác định entry.
2. **Early Rejection:** Validation xảy ra tại CheckTx, không cần chờ proposal execution.
3. **Data Integrity:** Đảm bảo tất cả entries có identifier hợp lệ.

---

### Test Case 2: Thất bại - Treasury Account Rỗng

### Đầu vào
- **Proposed Message:**
  - `MsgSetVestEntry` với TreasuryAccount = "" (chuỗi rỗng)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **Required Field:** TreasuryAccount là bắt buộc vì đây là nguồn tokens để vest.
2. **Validation:** Treasury account rỗng không hợp lệ cho vesting operation.

---

### Test Case 3: Thất bại - Start Time Sau End Time

### Đầu vào
- **Proposed Message:**
  - `MsgSetVestEntry` với:
    - StartTime: 2024-10-01 (sau)
    - EndTime: 2023-10-02 (trước)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **Logical Validation:** StartTime phải < EndTime cho vesting period hợp lệ.
2. **Time Logic:** Không thể có vesting period với start time sau end time.

---

### Test Case 4: Thất bại - Invalid Authority

### Đầu vào
- **Proposed Message:**
  - `MsgSetVestEntry` với Authority = địa chỉ của Bob (không phải gov module)

### Đầu ra
- **Proposal Submission:** FAIL
- **Proposals:** Không có proposals được tạo

### Tại sao chạy theo cách này?

1. **Authority Check:** Chỉ governance module có quyền set vest entries.
2. **Bảo mật:** Đảm bảo chỉ governance có thể quản lý vesting schedules.
3. **Early Rejection:** Validation tại thời điểm proposal submission.

---

### Test Case 5: Thất bại - Một Message Thất bại Gây Rollback

### Đầu vào
- **Proposed Messages:**
  1. `MsgSetVestEntry` hợp lệ cho "random_vester"
  2. `MsgSetVestEntry` với invalid authority cho "random_vester"

### Đầu ra
- **Proposal Submission:** FAIL (do message 2)
- **State:** Không có entries được tạo (bao gồm message 1)

### Tại sao chạy theo cách này?

1. **Atomic Execution:** Nếu một message thất bại, toàn bộ proposal thất bại.
2. **Rollback:** State được rollback về trước khi proposal execution.
3. **Consistency:** Đảm bảo không có partial state.

---

## Test Function: TestDeleteVestEntry_Success

### Test Case 1: Thành công - Xóa Một Vest Entry

### Đầu vào
- **Genesis State:** Có vest entry với VesterAccount = "random_vester"
- **Proposed Message:**
  - `MsgDeleteVestEntry` với VesterAccount = "random_vester"

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entry:** Được xóa khỏi state
- **Query:** Không thể query entry (trả về lỗi)

### Tại sao chạy theo cách này?

1. **Delete Operation:** `MsgDeleteVestEntry` xóa entry khỏi state.
2. **Cleanup:** Sau khi xóa, entry không còn tồn tại trong state.
3. **Governance Control:** Chỉ governance có quyền xóa entries.

---

### Test Case 2: Thành công - Xóa Hai Vest Entries

### Đầu vào
- **Genesis State:** Có 2 vest entries
- **Proposed Messages:**
  1. `MsgDeleteVestEntry` cho entry 1
  2. `MsgDeleteVestEntry` cho entry 2

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Vest Entries:** Cả hai đều được xóa

### Tại sao chạy theo cách này?

1. **Batch Deletion:** Một proposal có thể xóa nhiều entries.
2. **Independent Operations:** Mỗi deletion độc lập.

---

## Test Function: TestDeleteVestEntry_Failure

### Test Case 1: Thất bại - Vest Entry Không Tồn tại

### Đầu vào
- **Genesis State:** Không có vest entries
- **Proposed Message:**
  - `MsgDeleteVestEntry` cho "random_vester" (không tồn tại)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Không thay đổi

### Tại sao chạy theo cách này?

1. **Existence Check:** Không thể xóa entry không tồn tại.
2. **Error Handling:** Proposal thất bại nhưng không làm crash system.
3. **State Protection:** State không bị ảnh hưởng khi xóa entry không tồn tại.

---

### Test Case 2: Thất bại - Xóa Cùng Vest Entry Hai Lần

### Đầu vào
- **Genesis State:** Có 1 vest entry
- **Proposed Messages:**
  1. `MsgDeleteVestEntry` cho entry (thành công)
  2. `MsgDeleteVestEntry` cho cùng entry (thất bại vì đã xóa)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Entry vẫn bị xóa (từ message 1), nhưng proposal thất bại do message 2

### Tại sao chạy theo cách này?

1. **Sequential Execution:** Message 1 thành công, message 2 thất bại.
2. **Atomic Failure:** Khi message 2 thất bại, toàn bộ proposal thất bại và state rollback.
3. **State Consistency:** Entry không bị xóa vì proposal thất bại.

---

### Test Case 3: Thất bại - Entry Thứ Hai để Xóa Không Tồn tại

### Đầu vào
- **Genesis State:** Có 1 vest entry ("random_vester")
- **Proposed Messages:**
  1. `MsgDeleteVestEntry` cho "random_vester" (tồn tại)
  2. `MsgDeleteVestEntry` cho "random_vester_2" (không tồn tại)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Không thay đổi (rollback)

### Tại sao chạy theo cách này?

1. **Partial Failure:** Message 1 thành công nhưng message 2 thất bại.
2. **Rollback:** Toàn bộ proposal thất bại và state rollback.
3. **Consistency:** Đảm bảo không có partial deletion.

---

### Test Case 4: Thất bại - Invalid Authority

### Đầu vào
- **Genesis State:** Có vest entry
- **Proposed Message:**
  - `MsgDeleteVestEntry` với Authority = địa chỉ của Bob (không phải gov)

### Đầu ra
- **Proposal Submission:** FAIL
- **State:** Không thay đổi

### Tại sao chạy theo cách này?

1. **Authority Validation:** Chỉ governance module có quyền xóa entries.
2. **Bảo mật:** Đảm bảo chỉ governance có thể quản lý vesting.

---

## Tóm tắt Flow

### Set Vest Entry Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - VesterAccount không rỗng                                │
│    - TreasuryAccount không rỗng                             │
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
│    - Kiểm tra vest entry tồn tại                              │
│    - Nếu tồn tại: UPDATE                                     │
│    - Nếu không tồn tại: CREATE                               │
│    - Lưu trong state                                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. VERIFY STATE                                              │
│    - Query vest entry                                        │
│    - Xác minh thông tin đúng                                 │
└─────────────────────────────────────────────────────────────┘
```

### Delete Vest Entry Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - VesterAccount không rỗng                                │
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
│    - Kiểm tra vest entry tồn tại                              │
│    - Nếu tồn tại: DELETE                                     │
│    - Nếu không tồn tại: FAIL                                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. VERIFY STATE                                              │
│    - Query vest entry → Nên trả về lỗi                        │
│    - Xác minh entry không còn trong state                     │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Vest Entry State:**
   ```
   Không tồn tại → Tồn tại (CREATE)
   Tồn tại → Được cập nhật (UPDATE)
   Tồn tại → Không tồn tại (DELETE)
   ```

2. **Proposal Status:**
   ```
   SUBMITTED → PASSED (nếu tất cả messages thành công)
   SUBMITTED → FAILED (nếu bất kỳ message nào thất bại)
   ```

### Điểm quan trọng

1. **Idempotency:** 
   - `SetVestEntry` có thể được gọi nhiều lần với cùng VesterAccount
   - `DeleteVestEntry` chỉ thành công nếu entry tồn tại

2. **Atomic Execution:**
   - Nếu một message thất bại, toàn bộ proposal thất bại
   - State được rollback về trước khi proposal execution

3. **Authority:**
   - Chỉ governance module có quyền set/delete vest entries
   - Validation xảy ra tại cả submission và execution time

4. **Validation:**
   - VesterAccount và TreasuryAccount không thể rỗng
   - StartTime phải < EndTime
   - Entry phải tồn tại để xóa

5. **Batch Operations:**
   - Một proposal có thể chứa nhiều messages
   - Tất cả messages phải thành công để proposal pass

### Lý do thiết kế

1. **Governance Control:** Chỉ governance có quyền quản lý vesting schedules để đảm bảo decentralization và bảo mật.

2. **Flexibility:** Cho phép cập nhật vest entries để điều chỉnh vesting schedules khi cần.

3. **Safety:** Validation đảm bảo tính toàn vẹn dữ liệu và ngăn chặn invalid states.

4. **Atomic Operations:** Đảm bảo tính nhất quán state - không có partial updates.

5. **Error Handling:** Thông báo lỗi rõ ràng và rollback đúng cách khi có lỗi xảy ra.
