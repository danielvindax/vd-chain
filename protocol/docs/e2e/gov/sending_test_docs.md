# Tài liệu Test: Sending Module Governance Proposals

## Tổng quan

File test này xác minh governance proposal để **gửi tokens từ module account đến user account hoặc module account khác**. Đây là cơ chế để governance phân phối tokens từ treasury hoặc module accounts.

---

## Test Function: TestSendFromModuleToAccount

### Test Case 1: Thành công - Gửi từ Module đến User Account

### Đầu vào
- **Genesis State:**
  - Community Treasury module có số dư: 200 avdtn
- **Proposed Message:**
  - `MsgSendFromModuleToAccount`:
    - SenderModuleName: "community_treasury"
    - Recipient: Địa chỉ của Alice
    - Coin: 123 avdtn
    - Authority: gov module

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Module Balance:** 200 - 123 = 77 avdtn
- **Alice Balance:** Tăng 123 avdtn

### Tại sao chạy theo cách này?

1. **Module-to-Account Transfer:** Cho phép governance phân phối tokens từ module accounts (như treasury) đến user accounts.
2. **Use Cases:** 
   - Airdrops
   - Phân phối rewards
   - Treasury disbursements
3. **Governance Control:** Chỉ governance có quyền thực hiện transfers từ module accounts.

---

### Test Case 2: Thành công - Gửi từ Module đến Module Account

### Đầu vào
- **Genesis State:**
  - Community Treasury module có số dư: 123 avdtn
- **Proposed Message:**
  - `MsgSendFromModuleToAccount`:
    - SenderModuleName: "community_treasury"
    - Recipient: Địa chỉ Community Vester module
    - Coin: 123 avdtn
    - Authority: gov module

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Treasury Balance:** 0 avdtn (tất cả đã chuyển)
- **Vester Balance:** Tăng 123 avdtn

### Tại sao chạy theo cách này?

1. **Module-to-Module Transfer:** Cho phép chuyển tokens giữa các module accounts.
2. **Use Cases:**
   - Funding vesting contracts
   - Rebalancing module accounts
   - Treasury management
3. **Flexibility:** Hỗ trợ cả user accounts và module accounts làm recipients.

---

### Test Case 3: Thất bại - Số dư Không Đủ

### Đầu vào
- **Genesis State:**
  - Community Treasury module có số dư: 123 avdtn
- **Proposed Message:**
  - `MsgSendFromModuleToAccount`:
    - Coin: 124 avdtn (nhiều hơn số dư)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Không thay đổi (số dư vẫn giữ nguyên)

### Tại sao chạy theo cách này?

1. **Balance Check:** Module account phải có số dư đủ để chuyển.
2. **Execution-Time Validation:** Validation xảy ra khi proposal thực thi.
3. **State Protection:** Đảm bảo không có số dư âm.

---

### Test Case 4: Thất bại - Invalid Authority

### Đầu vào
- **Proposed Message:**
  - `MsgSendFromModuleToAccount` với Authority = địa chỉ sending module (thay vì gov module)

### Đầu ra
- **Proposal Submission:** FAIL
- **Proposals:** Không có proposals được tạo

### Tại sao chạy theo cách này?

1. **Authority Check:** Chỉ governance module có quyền gửi từ module accounts.
2. **Bảo mật:** Đảm bảo chỉ governance có thể kiểm soát module account transfers.
3. **Early Rejection:** Validation tại thời điểm proposal submission.

---

## Tóm tắt Flow

### Send From Module To Account Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - Authority = gov module                                  │
│    - SenderModuleName hợp lệ                                │
│    - Recipient address hợp lệ                               │
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
│    - Kiểm tra sender module balance >= coin amount         │
│    - Chuyển coins từ sender module                           │
│    - Chuyển coins đến recipient                             │
│    - Cập nhật số dư                                          │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Balance Changes:**
   ```
   Sender Module: Số dư -= Số tiền
   Recipient: Số dư += Số tiền
   ```

2. **Validation Points:**
   ```
   Submission: Authority check
   Execution: Balance check
   ```

### Điểm quan trọng

1. **Balance Validation:**
   - Sender module phải có số dư đủ
   - Validation xảy ra tại execution time
   - Nếu không đủ, proposal thất bại và state rollback

2. **Authority:**
   - Chỉ governance module có quyền
   - Validation tại thời điểm proposal submission

3. **Recipient Types:**
   - Hỗ trợ cả user accounts và module accounts
   - Recipient address phải hợp lệ

4. **Atomic Execution:**
   - Nếu số dư không đủ, toàn bộ proposal thất bại
   - State được rollback về trước khi execution

5. **Use Cases:**
   - Treasury disbursements
   - Phân phối rewards
   - Rebalancing module accounts
   - Funding vesting contracts

### Lý do thiết kế

1. **Governance Control:** Chỉ governance có quyền chuyển từ module accounts để đảm bảo decentralization và bảo mật.

2. **An toàn:** Balance validation đảm bảo không có số dư âm hoặc không đủ funds.

3. **Linh hoạt:** Hỗ trợ cả user và module accounts làm recipients để hỗ trợ nhiều use cases.

4. **Minh bạch:** Tất cả transfers từ module accounts đi qua governance proposals, đảm bảo tính minh bạch.
