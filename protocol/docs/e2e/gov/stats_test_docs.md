# Tài liệu Test: Stats Module Governance Proposals

## Tổng quan

File test này xác minh governance proposal để **cập nhật Stats Module Params**. Stats module quản lý theo dõi thống kê và window duration cho các metrics khác nhau.

---

## Test Function: TestUpdateParams

### Test Case 1: Thành công - Cập nhật Stats Module Params

### Đầu vào
- **Genesis State:**
  - Stats module có params:
    - WindowDuration: khác với proposal
- **Proposed Message:**
  - `MsgUpdateParams`:
    - WindowDuration: 1 giờ (thay đổi)
    - Authority: gov module

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Module Params:** Được cập nhật với WindowDuration mới

### Tại sao chạy theo cách này?

1. **Stats Configuration:** WindowDuration điều khiển time window cho statistics calculations.
2. **Time Window:** Parameter này định nghĩa thời gian thống kê được theo dõi và tổng hợp.
3. **Governance Control:** Chỉ governance có quyền cập nhật stats params.
4. **Linh hoạt:** Cho phép điều chỉnh statistics window khi cần.

---

### Test Case 2: Thất bại - Invalid Authority

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateParams` với Authority = địa chỉ stats module (thay vì gov module)

### Đầu ra
- **Proposal Submission:** FAIL
- **Proposals:** Không có proposals được tạo
- **State:** Không thay đổi

### Tại sao chạy theo cách này?

1. **Authority Check:** Chỉ governance module có quyền cập nhật stats params.
2. **Bảo mật:** Đảm bảo chỉ governance có thể thay đổi statistics configuration.
3. **Early Rejection:** Validation tại thời điểm proposal submission.

---

## Tóm tắt Flow

### Update Stats Module Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
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
│    - Cập nhật stats module params                             │
│    - Áp dụng WindowDuration mới                              │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Params Update:**
   ```
   Params Cũ → Params Mới (UPDATE)
   ```

2. **WindowDuration:**
   ```
   Duration Cũ → Duration Mới (UPDATE)
   ```

### Điểm quan trọng

1. **WindowDuration:**
   - Định nghĩa time window cho statistics tracking
   - Được sử dụng để tổng hợp và tính toán statistics
   - Có thể được điều chỉnh để tối ưu statistics collection

2. **Authority:**
   - Chỉ governance module có quyền cập nhật
   - Validation tại thời điểm proposal submission

3. **Atomic Execution:**
   - Nếu validation thất bại, toàn bộ proposal thất bại
   - State được rollback về trước khi execution

4. **Simplicity:**
   - Stats module có params tối thiểu (chỉ WindowDuration)
   - Quy trình cập nhật đơn giản

### Lý do thiết kế

1. **Governance Control:** Chỉ governance có quyền thay đổi statistics configuration để đảm bảo decentralization.

2. **Linh hoạt:** Cho phép điều chỉnh statistics window khi cần để tối ưu metrics collection.

3. **Simplicity:** Params tối thiểu giữ module đơn giản và tập trung.

4. **Nhất quán:** Đảm bảo stats module luôn có configuration hợp lệ.
