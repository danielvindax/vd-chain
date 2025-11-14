# Tài liệu Test: Fee Tiers Module Governance Proposals

## Tổng quan

File test này xác minh governance proposal để **cập nhật Perpetual Fee Params** trong Fee Tiers module. Fee tiers định nghĩa cấu trúc phí cho perpetual trading dựa trên trading volume và activity.

---

## Test Function: TestUpdateFeeTiersModuleParams

### Test Case 1: Thành công - Cập nhật Perpetual Fee Params

### Đầu vào
- **Genesis State:**
  - Fee tiers module có params khác với proposal
- **Proposed Message:**
  - `MsgUpdatePerpetualFeeParams`:
    - Tiers:
      - Tier 0:
        - Name: "test_tier_0"
        - MakerFeePpm: 11_000
        - TakerFeePpm: 22_000
        - Không có volume requirements (tier đầu tiên)
      - Tier 1:
        - Name: "test_tier_1"
        - AbsoluteVolumeRequirement: 200_000
        - TotalVolumeShareRequirementPpm: 100_000
        - MakerVolumeShareRequirementPpm: 50_000
        - MakerFeePpm: 1_000
        - TakerFeePpm: 2_000
    - Authority: gov module

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Fee Params:** Được cập nhật với cấu trúc tier mới

### Tại sao chạy theo cách này?

1. **Fee Tier Structure:** Fee tiers định nghĩa các mức phí khác nhau dựa trên trading volume và activity.
2. **First Tier:** Tier đầu tiên (tier 0) phải không có volume requirements - đây là tier mặc định cho tất cả traders.
3. **Volume Requirements:** Các tier cao hơn yêu cầu traders đáp ứng volume thresholds để đủ điều kiện cho phí thấp hơn.
4. **Maker/Taker Fees:** Phí khác nhau cho makers (liquidity providers) và takers (liquidity consumers).
5. **Governance Control:** Chỉ governance có quyền cập nhật fee tiers.

---

### Test Case 2: Thất bại - Không có Tiers

### Đầu vào
- **Proposed Message:**
  - `MsgUpdatePerpetualFeeParams` với Tiers = [] (mảng rỗng)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **Required Field:** Ít nhất một tier là bắt buộc cho cấu trúc phí.
2. **Early Validation:** Validation tại CheckTx để từ chối sớm.
3. **Data Integrity:** Đảm bảo fee tiers module luôn có cấu trúc tier hợp lệ.

---

### Test Case 3: Thất bại - Tier Đầu tiên Có Volume Requirement Khác Không

### Đầu vào
- **Proposed Message:**
  - `MsgUpdatePerpetualFeeParams` với:
    - Tier 0:
      - AbsoluteVolumeRequirement: 1 (khác không - SAI!)
      - MakerFeePpm: 1_000
      - TakerFeePpm: 2_000

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **First Tier Rule:** Tier đầu tiên (tier 0) phải có volume requirements bằng không vì đây là tier mặc định cho tất cả traders.
2. **Logical Constraint:** Nếu tier đầu tiên có volume requirements, traders mới không thể đủ điều kiện cho bất kỳ tier nào.
3. **Early Rejection:** Validation tại CheckTx để ngăn chặn configuration không hợp lệ.

---

### Test Case 4: Thất bại - Tổng của Maker Fee Thấp nhất và Taker Fee Thấp nhất là Âm

### Đầu vào
- **Proposed Message:**
  - `MsgUpdatePerpetualFeeParams` với:
    - Tier 0:
      - MakerFeePpm: -1_000 (âm - maker fee thấp nhất)
      - TakerFeePpm: 2_000
    - Tier 1:
      - MakerFeePpm: -888
      - TakerFeePpm: 500 (taker fee thấp nhất)
    - Tổng của fees thấp nhất: -1_000 + 500 = -500 (âm)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **Fee Validation:** Tổng của maker fee thấp nhất và taker fee thấp nhất trên tất cả tiers phải không âm.
2. **Economic Constraint:** Đảm bảo cấu trúc phí khả thi về mặt kinh tế - ít nhất một kết hợp của maker/taker fees nên không âm.
3. **Early Rejection:** Validation tại CheckTx để ngăn chặn cấu trúc phí không hợp lệ.

---

### Test Case 5: Thất bại - Invalid Authority

### Đầu vào
- **Proposed Message:**
  - `MsgUpdatePerpetualFeeParams` với Authority = địa chỉ fee tiers module (thay vì gov module)

### Đầu ra
- **Proposal Submission:** FAIL
- **Proposals:** Không có proposals được tạo

### Tại sao chạy theo cách này?

1. **Authority Check:** Chỉ governance module có quyền cập nhật fee tiers params.
2. **Bảo mật:** Đảm bảo chỉ governance có thể thay đổi cấu trúc phí.
3. **Early Rejection:** Validation tại thời điểm proposal submission.

---

## Tóm tắt Flow

### Update Perpetual Fee Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - Tiers array không rỗng                                  │
│    - Tier đầu tiên có volume requirements bằng không        │
│    - Tổng của maker fee thấp nhất + taker fee thấp nhất >= 0 │
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
│    - Cập nhật perpetual fee params                            │
│    - Áp dụng cấu trúc tier mới                                │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Fee Params Update:**
   ```
   Tiers Cũ → Tiers Mới (UPDATE)
   ```

2. **Tier Structure:**
   ```
   Tier 0: Không có volume requirements (tier mặc định)
   Tier 1+: Volume requirements để đủ điều kiện
   ```

### Điểm quan trọng

1. **Tier Structure:**
   - Tier đầu tiên (tier 0) phải có volume requirements bằng không
   - Các tier cao hơn yêu cầu volume thresholds
   - Mỗi tier có maker và taker fees

2. **Volume Requirements:**
   - AbsoluteVolumeRequirement: Ngưỡng trading volume tuyệt đối
   - TotalVolumeShareRequirementPpm: Tổng volume share (PPM)
   - MakerVolumeShareRequirementPpm: Maker volume share (PPM)

3. **Fee Validation:**
   - Tổng của maker fee thấp nhất + taker fee thấp nhất phải >= 0
   - Đảm bảo cấu trúc phí khả thi về mặt kinh tế
   - Validation tại CheckTx

4. **Authority:**
   - Chỉ governance module có quyền cập nhật
   - Validation tại thời điểm proposal submission

5. **Atomic Execution:**
   - Nếu validation thất bại, toàn bộ proposal thất bại
   - State được rollback về trước khi execution

### Lý do thiết kế

1. **Governance Control:** Chỉ governance có quyền thay đổi cấu trúc phí để đảm bảo decentralization và công bằng.

2. **An toàn:** Validation đảm bảo không có invalid states (empty tiers, invalid first tier, negative fee sum).

3. **Linh hoạt:** Cho phép điều chỉnh fee tiers khi cần để tối ưu trading incentives.

4. **Economic Viability:** Fee validation đảm bảo cấu trúc phí bền vững về mặt kinh tế.

5. **User Experience:** Tier đầu tiên không có requirements đảm bảo tất cả traders có thể tham gia.
