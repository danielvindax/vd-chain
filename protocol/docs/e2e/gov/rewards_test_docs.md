# Tài liệu Test: Rewards Module Governance Proposals

## Tổng quan

File test này xác minh governance proposal để **cập nhật Rewards Module Params**. Rewards module quản lý phân phối rewards cho users dựa trên hoạt động trading.

---

## Test Function: TestUpdateRewardsModuleParams

### Test Case 1: Thành công - Cập nhật Rewards Module Params

### Đầu vào
- **Genesis State:**
  - Rewards module có params:
    - TreasuryAccount: "test_treasury"
    - Denom: "avdtn"
    - DenomExponent: -18
    - MarketId: 1234
    - FeeMultiplierPpm: 700_000
- **Proposed Message:**
  - `MsgUpdateParams`:
    - TreasuryAccount: "test_treasury" (không thay đổi)
    - Denom: "avdtn" (không thay đổi)
    - DenomExponent: -5 (thay đổi từ -18)
    - MarketId: 0 (thay đổi từ 1234)
    - FeeMultiplierPpm: 700_001 (thay đổi từ 700_000)
    - Authority: gov module

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Module Params:** Được cập nhật với giá trị mới

### Tại sao chạy theo cách này?

1. **Rewards Configuration:** Các params này điều khiển cách rewards được tính toán và phân phối.
2. **FeeMultiplierPpm:** Đây là multiplier (parts per million) để tính rewards dựa trên trading fees.
3. **DenomExponent:** Exponent của denom để chuyển đổi giữa các đơn vị khác nhau.
4. **MarketId:** Market ID để theo dõi rewards cho market cụ thể.

---

### Test Case 2: Thất bại - Treasury Account Rỗng

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateParams` với TreasuryAccount = "" (chuỗi rỗng)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **Required Field:** TreasuryAccount là nguồn funds để phân phối rewards, không thể rỗng.
2. **Early Validation:** Validation tại CheckTx để từ chối sớm.
3. **Data Integrity:** Đảm bảo rewards module luôn có treasury account hợp lệ.

---

### Test Case 3: Thất bại - Denom Không Hợp lệ

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateParams` với Denom = "7avdtn" (không hợp lệ - bắt đầu bằng số)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **Denom Format:** Denom phải tuân theo format chuẩn (không thể bắt đầu bằng số).
2. **Validation:** Cosmos SDK có quy tắc về denom format.
3. **Early Rejection:** Validation tại CheckTx để ngăn chặn proposals không hợp lệ.

---

### Test Case 4: Thất bại - Fee Multiplier PPM Lớn Hơn 1 Triệu

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateParams` với FeeMultiplierPpm = 1_000_001 (> 1 triệu)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **Boundary Check:** FeeMultiplierPpm phải <= 1_000_000 (1 triệu ppm = 100%).
2. **PPM Format:** PPM (parts per million) có tối đa là 1 triệu.
3. **Logical Constraint:** Multiplier không thể > 100% trong context này.

---

### Test Case 5: Thất bại - Invalid Authority

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateParams` với Authority = địa chỉ rewards module (thay vì gov module)

### Đầu ra
- **Proposal Submission:** FAIL
- **Proposals:** Không có proposals được tạo

### Tại sao chạy theo cách này?

1. **Authority Check:** Chỉ governance module có quyền cập nhật rewards params.
2. **Bảo mật:** Đảm bảo chỉ governance có thể thay đổi rewards configuration.
3. **Early Rejection:** Validation tại thời điểm proposal submission.

---

## Tóm tắt Flow

### Update Rewards Module Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - TreasuryAccount không rỗng                             │
│    - Denom hợp lệ (format chuẩn)                            │
│    - FeeMultiplierPpm <= 1_000_000                          │
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
│    - Cập nhật rewards module params                          │
│    - Áp dụng configuration mới                                │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Params Update:**
   ```
   Params Cũ → Params Mới (UPDATE)
   ```

2. **Validation Points:**
   ```
   CheckTx: TreasuryAccount, Denom, FeeMultiplierPpm
   Submission: Authority
   ```

### Điểm quan trọng

1. **TreasuryAccount:**
   - Không được rỗng
   - Đây là nguồn funds để phân phối rewards
   - Validation tại CheckTx

2. **Denom Format:**
   - Phải tuân theo Cosmos SDK denom format
   - Không thể bắt đầu bằng số
   - Validation tại CheckTx

3. **FeeMultiplierPpm:**
   - Phải <= 1_000_000 (100%)
   - Đây là multiplier để tính rewards
   - Validation tại CheckTx

4. **Authority:**
   - Chỉ governance module có quyền
   - Validation tại thời điểm proposal submission

5. **Atomic Execution:**
   - Nếu validation thất bại, toàn bộ proposal thất bại
   - State được rollback về trước khi execution

### Lý do thiết kế

1. **Governance Control:** Chỉ governance có quyền thay đổi rewards configuration để đảm bảo decentralization.

2. **An toàn:** Validation đảm bảo không có invalid states (empty treasury, invalid denom, out of bounds multiplier).

3. **Linh hoạt:** Cho phép điều chỉnh rewards parameters khi cần để tối ưu phân phối rewards.

4. **Nhất quán:** Đảm bảo rewards module luôn có configuration hợp lệ.
