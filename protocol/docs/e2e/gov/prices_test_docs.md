# Tài liệu Test: Prices Module Governance Proposals

## Tổng quan

File test này xác minh governance proposal để **cập nhật Market Param** trong Prices module. Market param chứa thông tin về price feed configuration cho một market.

---

## Test Function: TestUpdateMarketParam

### Test Case 1: Thành công - Cập nhật Market Param

### Đầu vào
- **Genesis State:**
  - Có market param với ID = 0, Pair = "btc-avdtn", MinPriceChangePpm = 1_000
  - Market tồn tại trong market map
- **Proposed Message:**
  - `MsgUpdateMarketParam`:
    - Id: 0
    - Pair: "btc-avdtn" (không thay đổi)
    - MinPriceChangePpm: 2_002 (thay đổi từ 1_000)
    - Authority: gov module

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Market Param:** MinPriceChangePpm được cập nhật thành 2_002

### Tại sao chạy theo cách này?

1. **Market Param Update:** Cho phép điều chỉnh ngưỡng thay đổi giá tối thiểu cho một market.
2. **MinPriceChangePpm:** Đây là ngưỡng (parts per million) để xác định khi nào thay đổi giá đủ lớn để được coi là thay đổi đáng kể.
3. **Governance Control:** Chỉ governance có quyền cập nhật market params để đảm bảo tính nhất quán.

---

### Test Case 2: Thất bại - Market Param Không Tồn tại

### Đầu vào
- **Genesis State:** Chỉ có market param với ID = 0
- **Proposed Message:** 
  - `MsgUpdateMarketParam` với ID = 1 (không tồn tại)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Không thay đổi

### Tại sao chạy theo cách này?

1. **Existence Check:** Không thể cập nhật market param không tồn tại.
2. **Execution-Time Validation:** Validation xảy ra khi proposal thực thi, không phải khi submit.
3. **State Protection:** Đảm bảo không có partial updates.

---

### Test Case 3: Thất bại - Pair Name Mới Không Tồn tại trong Market Map

### Đầu vào
- **Genesis State:** 
  - Có market param với ID = 0, Pair = "btc-avdtn"
  - Market map chỉ có "btc-avdtn"
- **Proposed Message:**
  - `MsgUpdateMarketParam` với Pair = "nonexistent-pair"

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Không thay đổi

### Tại sao chạy theo cách này?

1. **Market Map Integration:** Market param phải tham chiếu đến một pair tồn tại trong market map.
2. **Consistency:** Đảm bảo market param và market map luôn đồng bộ với nhau.
3. **Dependency:** Market param phụ thuộc vào market map configuration.

---

### Test Case 4: Thất bại - Pair Rỗng

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateMarketParam` với Pair = "" (chuỗi rỗng)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **Required Field:** Pair là identifier của market, không thể rỗng.
2. **Early Validation:** Validation tại CheckTx để từ chối sớm.
3. **Data Integrity:** Đảm bảo tất cả market params có pair name hợp lệ.

---

### Test Case 5: Thất bại - Invalid Authority

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateMarketParam` với Authority = địa chỉ của Alice (không phải gov module)

### Đầu ra
- **Proposal Submission:** FAIL
- **Proposals:** Không có proposals được tạo

### Tại sao chạy theo cách này?

1. **Authority Check:** Chỉ governance module có quyền cập nhật market params.
2. **Bảo mật:** Đảm bảo chỉ governance có thể thay đổi price feed configuration.
3. **Early Rejection:** Validation tại thời điểm proposal submission.

---

## Tóm tắt Flow

### Update Market Param Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - Pair không rỗng                                         │
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
│    - Kiểm tra market param tồn tại                            │
│    - Kiểm tra pair tồn tại trong market map                  │
│    - Cập nhật market param                                    │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Market Param State:**
   ```
   Tồn tại → Được cập nhật (UPDATE)
   Không tồn tại → FAIL
   ```

2. **Market Map Integration:**
   ```
   Pair phải tồn tại trong market map
   Không thể cập nhật với pair không tồn tại
   ```

### Điểm quan trọng

1. **Existence Validation:**
   - Market param phải tồn tại để cập nhật
   - Pair phải tồn tại trong market map

2. **Authority:**
   - Chỉ governance module có quyền cập nhật
   - Validation tại cả submission và execution time

3. **MinPriceChangePpm:**
   - Đây là ngưỡng để xác định thay đổi giá đáng kể
   - Có thể được điều chỉnh để tinh chỉnh tần suất cập nhật giá

4. **Market Map Dependency:**
   - Market param phụ thuộc vào market map
   - Đảm bảo tính nhất quán giữa hai hệ thống

5. **Atomic Execution:**
   - Nếu validation thất bại, toàn bộ proposal thất bại
   - State được rollback về trước khi execution

### Lý do thiết kế

1. **Governance Control:** Chỉ governance có quyền thay đổi price feed configuration để đảm bảo decentralization và tính nhất quán.

2. **Consistency:** Market param và market map phải luôn đồng bộ với nhau để đảm bảo price feeds hoạt động đúng.

3. **An toàn:** Validation đảm bảo không có invalid states (empty pair, non-existent market).

4. **Linh hoạt:** Cho phép điều chỉnh MinPriceChangePpm để tối ưu tần suất cập nhật giá.
