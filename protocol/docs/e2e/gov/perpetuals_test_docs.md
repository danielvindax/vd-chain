# Tài liệu Test: Perpetuals Module Governance Proposals

## Tổng quan

File test này xác minh các governance proposals liên quan đến **Perpetuals Module**, bao gồm:
1. **Update Module Params:** Cập nhật parameters của perpetuals module
2. **Update Perpetual Params:** Cập nhật parameters của một perpetual cụ thể
3. **Set Liquidity Tier:** Tạo hoặc cập nhật liquidity tier

---

## Test Function: TestUpdatePerpetualsModuleParams

### Test Case 1: Thành công - Cập nhật Module Params

### Đầu vào
- **Genesis State:** Perpetuals module có params khác với proposal
- **Proposed Message:**
  - `MsgUpdateParams`:
    - FundingRateClampFactorPpm: 123_456
    - PremiumVoteClampFactorPpm: 123_456_789
    - MinNumVotesPerSample: 15
    - Authority: gov module

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Module Params:** Được cập nhật với giá trị mới

### Tại sao chạy theo cách này?

1. **Module-Level Params:** Đây là các parameters áp dụng cho toàn bộ perpetuals module, không phải một perpetual cụ thể.
2. **Governance Control:** Chỉ governance có quyền cập nhật module params.
3. **Validation:** Tất cả giá trị phải > 0 để đảm bảo tính hợp lệ.

---

### Test Case 2-4: Thất bại - Giá trị Zero

### Đầu vào
- **Proposed Message:** Một trong các params = 0:
  - FundingRateClampFactorPpm = 0
  - PremiumVoteClampFactorPpm = 0
  - MinNumVotesPerSample = 0

### Đầu ra
- **CheckTx:** FAIL
- **State:** Không thay đổi

### Tại sao chạy theo cách này?

1. **Non-Zero Validation:** Tất cả params phải > 0 vì chúng được sử dụng trong calculations.
2. **Early Rejection:** Validation tại CheckTx để từ chối sớm, không cần chờ proposal execution.

---

### Test Case 5: Thất bại - Invalid Authority

### Đầu vào
- **Proposed Message:** Authority = perpetuals module (thay vì gov module)

### Đầu ra
- **Proposal Submission:** FAIL

### Tại sao chạy theo cách này?

1. **Authority Check:** Chỉ governance module có quyền cập nhật module params.
2. **Bảo mật:** Đảm bảo chỉ governance có thể thay đổi module-level settings.

---

## Test Function: TestUpdatePerpetualsParams

### Test Case 1: Thành công - Cập nhật Perpetual Params

### Đầu vào
- **Genesis State:** 
  - Có perpetual với ID = 0
  - Có liquidity tier với ID = 123
  - Có market với ID = 4
- **Proposed Message:**
  - `MsgUpdatePerpetualParams`:
    - Id: 0
    - Ticker: "BTC-VDTN" (thay đổi)
    - MarketId: 4
    - DefaultFundingPpm: 500 (thay đổi)
    - LiquidityTier: 123
    - MarketType: không thay đổi (immutable)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Perpetual Params:** Được cập nhật (trừ MarketType)

### Tại sao chạy theo cách này?

1. **Perpetual-Level Update:** Cập nhật params của một perpetual cụ thể.
2. **MarketType Immutable:** MarketType không thể thay đổi sau khi perpetual được tạo.
3. **Dependencies:** Phải có liquidity tier và market tồn tại trước.

---

### Test Case 2: Thất bại - Ticker Rỗng

### Đầu vào
- **Proposed Message:** Ticker = "" (rỗng)

### Đầu ra
- **CheckTx:** FAIL

### Tại sao chạy theo cách này?

1. **Required Field:** Ticker là identifier của perpetual, không thể rỗng.

---

### Test Case 3: Thất bại - Default Funding PPM Vượt Quá Tối đa

### Đầu vào
- **Proposed Message:** DefaultFundingPpm = 1_000_001 (> 1 triệu)

### Đầu ra
- **CheckTx:** FAIL

### Tại sao chạy theo cách này?

1. **Boundary Check:** DefaultFundingPpm phải <= 1_000_000 (1 triệu ppm = 100%).
2. **PPM Format:** PPM (parts per million) có tối đa là 1 triệu.

---

### Test Case 4: Thất bại - Invalid Authority

### Đầu vào
- **Proposed Message:** Authority = perpetuals module (thay vì gov)

### Đầu ra
- **Proposal Submission:** FAIL

### Tại sao chạy theo cách này?

1. **Governance Control:** Chỉ governance có quyền cập nhật perpetual params.

---

### Test Case 5: Thất bại - Liquidity Tier Không Tồn tại

### Đầu vào
- **Genesis State:** Chỉ có liquidity tier ID = 123
- **Proposed Message:** LiquidityTier = 124 (không tồn tại)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`

### Tại sao chạy theo cách này?

1. **Dependency Check:** Perpetual phải tham chiếu đến một liquidity tier tồn tại.
2. **Execution-Time Validation:** Validation xảy ra khi proposal thực thi, không phải khi submit.

---

### Test Case 6: Thất bại - Market ID Không Tồn tại

### Đầu vào
- **Genesis State:** Chỉ có market ID = 4
- **Proposed Message:** MarketId = 5 (không tồn tại)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`

### Tại sao chạy theo cách này?

1. **Market Dependency:** Perpetual phải tham chiếu đến một market tồn tại.
2. **Execution Validation:** Kiểm tra xảy ra khi thực thi proposal.

---

## Test Function: TestSetLiquidityTier

### Test Case 1: Thành công - Tạo Liquidity Tier Mới

### Đầu vào
- **Genesis State:** Không có liquidity tier ID = 5678
- **Proposed Message:**
  - `MsgSetLiquidityTier`:
    - Id: 5678
    - Name: "Test Tier"
    - InitialMarginPpm: 765_432
    - MaintenanceFractionPpm: 345_678
    - ImpactNotional: 654_321

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Liquidity Tier:** Được tạo mới

### Tại sao chạy theo cách này?

1. **Create Operation:** Khi liquidity tier không tồn tại, sẽ tạo mới.
2. **Idempotency:** Có thể cập nhật sau khi tạo.

---

### Test Case 2: Thành công - Cập nhật Liquidity Tier Hiện có

### Đầu vào
- **Genesis State:** Có liquidity tier ID = 5678
- **Proposed Message:** `MsgSetLiquidityTier` với cùng ID

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Liquidity Tier:** Được cập nhật

### Tại sao chạy theo cách này?

1. **Update Operation:** Khi liquidity tier đã tồn tại, sẽ cập nhật thay vì tạo mới.

---

### Test Case 3-5: Thất bại - Giá trị Không Hợp lệ

### Đầu vào
- **Proposed Message:** Một trong các trường hợp:
  - InitialMarginPpm = 1_000_001 (> tối đa)
  - MaintenanceFractionPpm = 1_000_001 (> tối đa)
  - ImpactNotional = 0

### Đầu ra
- **CheckTx:** FAIL

### Tại sao chạy theo cách này?

1. **Boundary Validation:** 
   - Giá trị PPM phải <= 1_000_000
   - ImpactNotional phải > 0 (không thể = 0)

---

### Test Case 6: Thất bại - Invalid Authority

### Đầu vào
- **Proposed Message:** Authority = perpetuals module (thay vì gov)

### Đầu ra
- **Proposal Submission:** FAIL

### Tại sao chạy theo cách này?

1. **Governance Control:** Chỉ governance có quyền set liquidity tiers.

---

## Tóm tắt Flow

### Update Module Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - Tất cả params > 0                                       │
│    - Authority = gov module                                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                                │
│    - Validate authority                                    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTION                                        │
│    - Cập nhật module params                                   │
│    - Áp dụng cho tất cả perpetuals                          │
└─────────────────────────────────────────────────────────────┘
```

### Update Perpetual Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - Ticker không rỗng                                       │
│    - DefaultFundingPpm <= 1_000_000                         │
│    - Authority = gov module                                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTION                                        │
│    - Kiểm tra perpetual tồn tại                               │
│    - Kiểm tra liquidity tier tồn tại                          │
│    - Kiểm tra market tồn tại                                  │
│    - Cập nhật params (trừ MarketType)                        │
└─────────────────────────────────────────────────────────────┘
```

### Set Liquidity Tier Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - InitialMarginPpm <= 1_000_000                          │
│    - MaintenanceFractionPpm <= 1_000_000                     │
│    - ImpactNotional > 0                                      │
│    - Authority = gov module                                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTION                                        │
│    - Kiểm tra liquidity tier tồn tại                          │
│    - Nếu tồn tại: UPDATE                                     │
│    - Nếu không tồn tại: CREATE                               │
└─────────────────────────────────────────────────────────────┘
```

### Điểm quan trọng

1. **Module vs Perpetual Params:**
   - Module params: Áp dụng cho toàn bộ module
   - Perpetual params: Áp dụng cho perpetual cụ thể

2. **Immutable Fields:**
   - MarketType của perpetual không thể thay đổi sau khi tạo

3. **Dependencies:**
   - Perpetual phải tham chiếu đến liquidity tier và market tồn tại
   - Validation xảy ra tại execution time

4. **PPM Format:**
   - PPM (parts per million) có tối đa = 1_000_000 (100%)
   - Tất cả giá trị PPM phải <= 1_000_000

5. **Authority:**
   - Tất cả updates phải có authority = gov module
   - Validation tại cả submission và execution time

### Lý do thiết kế

1. **Governance Control:** Chỉ governance có quyền thay đổi perpetuals configuration để đảm bảo decentralization.

2. **An toàn:** Validation đảm bảo không có invalid states (zero values, out of bounds).

3. **Dependency Management:** Đảm bảo perpetuals chỉ tham chiếu đến objects tồn tại.

4. **Linh hoạt:** Cho phép cập nhật params khi cần để điều chỉnh market conditions.
