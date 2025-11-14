# Tài liệu Test: Bridge Module Governance Proposals

## Tổng quan

File test này xác minh các governance proposals liên quan đến **Bridge Module**, bao gồm:
1. **Update Event Params:** Cập nhật event parameters cho bridge operations
2. **Update Propose Params:** Cập nhật propose parameters cho bridge proposals
3. **Update Safety Params:** Cập nhật safety parameters cho bridge safety controls

---

## Test Function: TestUpdateEventParams

### Test Case 1: Thành công - Cập nhật Event Params

### Đầu vào
- **Genesis State:**
  - Event params:
    - Denom: "testdenom"
    - EthChainId: 123
    - EthAddress: "0x0123"
- **Proposed Message:**
  - `MsgUpdateEventParams`:
    - Denom: "advtnt" (thay đổi)
    - EthChainId: 1 (thay đổi)
    - EthAddress: "0xabcd" (thay đổi)
    - Authority: gov module

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Event Params:** Được cập nhật với giá trị mới

### Tại sao chạy theo cách này?

1. **Event Configuration:** Event params điều khiển cách bridge events được xử lý và validate.
2. **Denom:** Denomination của tokens được sử dụng trong bridge operations.
3. **EthChainId:** Ethereum chain ID cho cross-chain bridge operations.
4. **EthAddress:** Ethereum address cho bridge contract hoặc operations.
5. **Governance Control:** Chỉ governance có quyền cập nhật event params.

---

### Test Case 2: Thất bại - ETH Address Rỗng

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateEventParams` với EthAddress = "" (chuỗi rỗng)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **Required Field:** EthAddress là bắt buộc cho bridge operations, không thể rỗng.
2. **Early Validation:** Validation tại CheckTx để từ chối sớm.
3. **Data Integrity:** Đảm bảo bridge luôn có Ethereum address hợp lệ.

---

### Test Case 3: Thất bại - Invalid Authority

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateEventParams` với Authority = địa chỉ của Bob (không phải gov module)

### Đầu ra
- **Proposal Submission:** FAIL
- **Proposals:** Không có proposals được tạo

### Tại sao chạy theo cách này?

1. **Authority Check:** Chỉ governance module có quyền cập nhật event params.
2. **Bảo mật:** Đảm bảo chỉ governance có thể thay đổi bridge event configuration.
3. **Early Rejection:** Validation tại thời điểm proposal submission.

---

## Test Function: TestUpdateProposeParams

### Test Case 1: Thành công - Cập nhật Propose Params

### Đầu vào
- **Genesis State:**
  - Propose params:
    - MaxBridgesPerBlock: 10
    - ProposeDelayDuration: 1 phút
    - SkipRatePpm: 800_000
    - SkipIfBlockDelayedByDuration: 1 phút
- **Proposed Message:**
  - `MsgUpdateProposeParams`:
    - MaxBridgesPerBlock: 7 (thay đổi)
    - ProposeDelayDuration: 1 giây (thay đổi)
    - SkipRatePpm: 700_007 (thay đổi)
    - SkipIfBlockDelayedByDuration: 1 giây (thay đổi)
    - Authority: gov module

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Propose Params:** Được cập nhật với giá trị mới

### Tại sao chạy theo cách này?

1. **Propose Configuration:** Propose params điều khiển cách bridge proposals được submit và xử lý.
2. **MaxBridgesPerBlock:** Số lượng bridges tối đa có thể được propose mỗi block.
3. **ProposeDelayDuration:** Thời gian delay trước khi proposal có thể được submit.
4. **SkipRatePpm:** Tỷ lệ (parts per million) để skip proposals trong một số điều kiện.
5. **SkipIfBlockDelayedByDuration:** Ngưỡng thời gian để skip proposals nếu block bị delay.

---

### Test Case 2: Thất bại - Propose Delay Duration Âm

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateProposeParams` với ProposeDelayDuration = -1 giây (âm)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **Non-Negative Validation:** Duration không thể âm.
2. **Early Rejection:** Validation tại CheckTx để từ chối sớm.
3. **Logical Constraint:** Duration âm không có ý nghĩa cho delay.

---

### Test Case 3: Thất bại - Skip If Block Delayed By Duration Âm

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateProposeParams` với SkipIfBlockDelayedByDuration = -1 giây (âm)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **Non-Negative Validation:** Duration không thể âm.
2. **Early Rejection:** Validation tại CheckTx.
3. **Logical Constraint:** Duration âm là không hợp lệ.

---

### Test Case 4: Thất bại - Skip Rate PPM Vượt Quá Giới Hạn

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateProposeParams` với SkipRatePpm = 1_000_001 (> 1 triệu)

### Đầu ra
- **CheckTx:** FAIL
- **Proposal:** Không được submit

### Tại sao chạy theo cách này?

1. **Boundary Check:** SkipRatePpm phải <= 1_000_000 (1 triệu ppm = 100%).
2. **PPM Format:** PPM (parts per million) có tối đa là 1 triệu.
3. **Logical Constraint:** Rate không thể vượt quá 100%.

---

### Test Case 5: Thất bại - Invalid Authority

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateProposeParams` với Authority = địa chỉ của Alice (không phải gov module)

### Đầu ra
- **Proposal Submission:** FAIL
- **Proposals:** Không có proposals được tạo

### Tại sao chạy theo cách này?

1. **Authority Check:** Chỉ governance module có quyền cập nhật propose params.
2. **Bảo mật:** Đảm bảo chỉ governance có thể thay đổi bridge proposal configuration.

---

## Test Function: TestUpdateSafetyParams

### Test Case 1: Thành công - Cập nhật Safety Params

### Đầu vào
- **Genesis State:**
  - Safety params:
    - IsDisabled: false
    - DelayBlocks: 10
- **Proposed Message:**
  - `MsgUpdateSafetyParams`:
    - IsDisabled: true (thay đổi)
    - DelayBlocks: 5 (thay đổi)
    - Authority: gov module

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Safety Params:** Được cập nhật với giá trị mới

### Tại sao chạy theo cách này?

1. **Safety Configuration:** Safety params điều khiển cơ chế an toàn cho bridge operations.
2. **IsDisabled:** Flag để enable/disable bridge safety checks.
3. **DelayBlocks:** Số lượng blocks để delay trước khi thực thi bridge operations.
4. **Governance Control:** Chỉ governance có quyền cập nhật safety params.

---

### Test Case 2: Thất bại - Invalid Authority

### Đầu vào
- **Proposed Message:**
  - `MsgUpdateSafetyParams` với Authority = địa chỉ của Alice (không phải gov module)

### Đầu ra
- **Proposal Submission:** FAIL
- **Proposals:** Không có proposals được tạo

### Tại sao chạy theo cách này?

1. **Authority Check:** Chỉ governance module có quyền cập nhật safety params.
2. **Bảo mật:** Đảm bảo chỉ governance có thể thay đổi bridge safety configuration.
3. **Early Rejection:** Validation tại thời điểm proposal submission.

---

## Tóm tắt Flow

### Update Event Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - EthAddress không rỗng                                   │
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
│    - Cập nhật event params                                    │
│    - Áp dụng configuration mới                                │
└─────────────────────────────────────────────────────────────┘
```

### Update Propose Params Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. VALIDATE INPUT                                            │
│    - ProposeDelayDuration >= 0                               │
│    - SkipIfBlockDelayedByDuration >= 0                      │
│    - SkipRatePpm <= 1_000_000                                │
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
│    - Cập nhật propose params                                  │
│    - Áp dụng configuration mới                                │
└─────────────────────────────────────────────────────────────┘
```

### Update Safety Params Process

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
│    - Cập nhật safety params                                   │
│    - Áp dụng configuration mới                                │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Event Params Update:**
   ```
   Params Cũ → Params Mới (UPDATE)
   ```

2. **Propose Params Update:**
   ```
   Params Cũ → Params Mới (UPDATE)
   ```

3. **Safety Params Update:**
   ```
   Params Cũ → Params Mới (UPDATE)
   ```

### Điểm quan trọng

1. **Event Params:**
   - Denom: Token denomination cho bridge operations
   - EthChainId: Ethereum chain ID cho cross-chain operations
   - EthAddress: Không được rỗng
   - Validation tại CheckTx

2. **Propose Params:**
   - MaxBridgesPerBlock: Số lượng bridges tối đa mỗi block
   - ProposeDelayDuration: Phải >= 0
   - SkipRatePpm: Phải <= 1_000_000
   - SkipIfBlockDelayedByDuration: Phải >= 0
   - Validation tại CheckTx

3. **Safety Params:**
   - IsDisabled: Enable/disable safety checks
   - DelayBlocks: Số lượng blocks để delay
   - Không có CheckTx validation (chỉ authority check)

4. **Authority:**
   - Chỉ governance module có quyền cập nhật tất cả params
   - Validation tại thời điểm proposal submission

5. **Atomic Execution:**
   - Nếu validation thất bại, toàn bộ proposal thất bại
   - State được rollback về trước khi execution

### Lý do thiết kế

1. **Governance Control:** Chỉ governance có quyền thay đổi bridge configuration để đảm bảo decentralization và bảo mật.

2. **An toàn:** Validation đảm bảo không có invalid states (empty address, negative duration, out of bounds rate).

3. **Linh hoạt:** Cho phép điều chỉnh bridge parameters khi cần để tối ưu bridge operations.

4. **Nhất quán:** Đảm bảo bridge module luôn có configuration hợp lệ.
