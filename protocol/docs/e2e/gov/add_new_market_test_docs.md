# Tài liệu Test: Add New Market Proposal

## Tổng quan

File test này xác minh chức năng **add new market** thông qua governance proposals. Test đảm bảo rằng khi tạo một market mới, hệ thống sẽ:
1. Tạo Oracle Market (price feed)
2. Tạo Perpetual contract
3. Tạo CLOB Pair với trạng thái INITIALIZING
4. Sử dụng DelayMessage để chuyển CLOB Pair sang ACTIVE sau một số block
5. Enable market trong market map

---

## Test Case 1: Thành công với 4 Standard Messages (Delay Blocks = 10)

### Đầu vào
- **Proposed Messages:**
  1. `MsgCreateOracleMarket`: Tạo market param với ID = 1001
  2. `MsgCreatePerpetual`: Tạo perpetual với ID = 1001
  3. `MsgCreateClobPair`: Tạo CLOB pair với ID = 1001, status = INITIALIZING
  4. `MsgDelayMessage`: Delay message để cập nhật CLOB pair sang ACTIVE sau 10 blocks
- **Market Map:** Market ban đầu bị disabled trong market map

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **Market Param:** Được tạo với ID = 1001
- **Market Price:** Được khởi tạo với price = 0
- **Perpetual:** Được tạo với ID = 1001
- **ClobPair:** 
  - Ban đầu: status = INITIALIZING
  - Sau 10 blocks: status = ACTIVE
- **Market Map:** Market được enable sau khi CLOB pair chuyển sang ACTIVE

### Tại sao chạy theo cách này?

1. **Thứ tự Message là Quan trọng:** Messages phải được thực thi theo thứ tự:
   - Oracle Market trước (cần cho price feed)
   - Perpetual tiếp theo (cần cho CLOB pair)
   - CLOB Pair sau (phụ thuộc vào perpetual)
   - DelayMessage cuối (để kích hoạt CLOB pair)

2. **Delay Blocks = 10:** CLOB pair không được kích hoạt ngay lập tức mà phải chờ 10 blocks để đảm bảo:
   - Tất cả dependencies được thiết lập đầy đủ
   - Oracle có thời gian cập nhật price
   - Hệ thống có thời gian validate state

3. **Market Map Integration:** Market phải được enable trong market map để cho phép trading, điều này chỉ xảy ra sau khi CLOB pair chuyển sang ACTIVE.

---

## Test Case 2: Thành công với Delay Blocks = 1

### Đầu vào
- **Proposed Messages:** Giống Test Case 1
- **Delay Blocks:** 1 (thay vì 10)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **ClobPair:** Chuyển sang ACTIVE sau 1 block

### Tại sao chạy theo cách này?

1. **Minimum Delay:** Test này đảm bảo delay blocks có thể là 1 (tối thiểu), không nhất thiết phải là 10.
2. **Fast Activation:** Một số trường hợp có thể cần kích hoạt market nhanh hơn, delay = 1 cho phép điều này.

---

## Test Case 3: Thành công với Delay Blocks = 0

### Đầu vào
- **Proposed Messages:** Giống Test Case 1
- **Delay Blocks:** 0 (không delay)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED`
- **ClobPair:** Không delay, nhưng vẫn cần được kích hoạt thông qua cơ chế delay message

### Tại sao chạy theo cách này?

1. **Zero Delay:** Test này đảm bảo delay = 0 cũng được hỗ trợ, có nghĩa message có thể được thực thi ngay lập tức trong block tiếp theo.
2. **Edge Case:** Đây là edge case để đảm bảo hệ thống xử lý delay = 0 đúng cách.

---

## Test Case 4: Thành công với Delayed UpdateClobPair Message Thất bại

### Đầu vào
- **Proposed Messages:**
  1. `MsgCreateOracleMarket`: ID = 1001
  2. `MsgCreatePerpetual`: ID = 1001
  3. `MsgCreateClobPair`: ID = 1001
  4. `MsgDelayMessage`: Chứa `MsgUpdateClobPair` với ClobPairId = 9999 (không tồn tại)
- **Delay Blocks:** 10

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_PASSED` (proposal vẫn pass vì các messages khác thành công)
- **ClobPair:** Vẫn ở trạng thái INITIALIZING (không được cập nhật vì delayed message thất bại)
- **Market Map:** Market vẫn disabled

### Tại sao chạy theo cách này?

1. **Delayed Message Failure:** Khi delayed message thất bại, nó không làm thất bại toàn bộ proposal vì proposal đã được thực thi thành công.
2. **Partial Success:** Các messages khác (tạo market, perpetual, clob pair) vẫn thành công, chỉ delayed update message thất bại.
3. **State Consistency:** ClobPair vẫn ở INITIALIZING, không bị ảnh hưởng bởi delayed message thất bại.

---

## Test Case 5: Thất bại - Messages được Sắp xếp Sai Thứ tự

### Đầu vào
- **Proposed Messages (Thứ tự Sai):**
  1. `MsgCreateOracleMarket`: ID = 1001
  2. `MsgCreateClobPair`: ID = 1001 (trước khi tạo perpetual - SAI!)
  3. `MsgCreatePerpetual`: ID = 1001
  4. `MsgDelayMessage`: Cập nhật CLOB pair

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Không có gì được tạo (full rollback)

### Tại sao chạy theo cách này?

1. **Dependency Order:** CLOB Pair phụ thuộc vào Perpetual, vì vậy Perpetual phải được tạo trước.
2. **Atomic Execution:** Thực thi proposal là atomic - nếu một message thất bại, toàn bộ proposal thất bại và state được rollback.
3. **Validation:** Hệ thống validate dependencies và từ chối nếu thứ tự sai.

---

## Test Case 6: Thất bại - Objects Đã Tồn tại

### Đầu vào
- **Proposed Messages:**
  1. `MsgCreateOracleMarket`: ID = 5 (đã tồn tại trong genesis)
  2. `MsgCreatePerpetual`: ID = 5 (đã tồn tại)
  3. `MsgCreateClobPair`: ID = 5 (đã tồn tại)
  4. `MsgDelayMessage`: Cập nhật CLOB pair

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Không có gì mới được tạo

### Tại sao chạy theo cách này?

1. **Idempotency:** Không thể tạo objects với IDs đã tồn tại.
2. **Error Handling:** Hệ thống phát hiện conflict và từ chối proposal.
3. **State Protection:** Đảm bảo không có duplicate IDs trong hệ thống.

---

## Test Case 7: Thất bại - Invalid Signer (Proposal Submission)

### Đầu vào
- **Proposed Messages:**
  1. `MsgCreateOracleMarket`: Authority = CLOB module address (SAI! Phải là gov module)
  2. `MsgCreatePerpetual`: Authority = gov module (đúng)
  3. `MsgCreateClobPair`: Authority = gov module (đúng)
  4. `MsgDelayMessage`: Authority = gov module (đúng)

### Đầu ra
- **Proposal Submission:** FAIL (không thể submit)
- **Proposals:** Không có proposals được tạo

### Tại sao chạy theo cách này?

1. **Authority Validation:** Mỗi message phải có authority đúng:
   - `MsgCreateOracleMarket` phải có authority = gov module
   - Validation xảy ra tại thời điểm proposal submission
2. **Early Rejection:** Proposal bị từ chối ngay khi submit, không cần chờ execution.
3. **Bảo mật:** Đảm bảo chỉ authority đúng mới có thể tạo objects.

---

## Test Case 8: Thất bại - Invalid Signer trên MsgDelayMessage

### Đầu vào
- **Proposed Messages:**
  1. `MsgCreateOracleMarket`: Authority = gov module (đúng)
  2. `MsgCreatePerpetual`: Authority = gov module (đúng)
  3. `MsgCreateClobPair`: Authority = gov module (đúng)
  4. `MsgDelayMessage`: 
     - Authority = gov module (đúng)
     - Nhưng wrapped message (`MsgUpdateClobPair`) có authority = gov module (SAI! Phải là delaymsg module)

### Đầu ra
- **Proposal Status:** `PROPOSAL_STATUS_FAILED`
- **State:** Market, Perpetual, CLOB Pair được tạo nhưng proposal thất bại khi thực thi delayed message

### Tại sao chạy theo cách này?

1. **Nested Authority:** `MsgDelayMessage` chứa một message khác (`MsgUpdateClobPair`), và message đó cũng có authority riêng của nó.
2. **Delayed Validation:** Authority của wrapped message chỉ được validate khi delayed message được thực thi, không phải khi proposal được submit.
3. **Partial State:** Các messages trước đó thực thi thành công, nhưng delayed message thất bại làm proposal thất bại.

---

## Tóm tắt Flow

### Add New Market Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. KHỞI TẠO GENESIS STATE                                   │
│    - Market map với market disabled                          │
│    - Không có market/perpetual/clob pair với ID mới         │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                               │
│    - Validate authority của tất cả messages                 │
│    - Validate thứ tự message                                │
│    - Validate không có duplicate IDs                        │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTION (Nếu Submit Thành công)                │
│    a. Thực thi MsgCreateOracleMarket                        │
│       → Tạo MarketParam và MarketPrice (price = 0)          │
│    b. Thực thi MsgCreatePerpetual                            │
│       → Tạo Perpetual contract                             │
│    c. Thực thi MsgCreateClobPair                             │
│       → Tạo CLOB Pair với trạng thái INITIALIZING            │
│    d. Thực thi MsgDelayMessage                              │
│       → Lên lịch MsgUpdateClobPair thực thi sau N blocks      │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. DELAYED MESSAGE EXECUTION                                │
│    - Sau N blocks (0, 1, hoặc 10), delayed message thực thi │
│    - Validate authority của wrapped message                 │
│    - Cập nhật CLOB Pair status: INITIALIZING → ACTIVE       │
│    - Enable market trong market map                         │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. TRADING ENABLED                                          │
│    - Sau khi CLOB Pair = ACTIVE và market enabled           │
│    - Users có thể đặt orders (nhưng cần oracle price > 0)   │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **CLOB Pair Status Transition:**
   ```
   Không tồn tại → INITIALIZING → ACTIVE
   ```

2. **Market Map State:**
   ```
   Disabled → Enabled (khi CLOB Pair = ACTIVE)
   ```

3. **Oracle Price:**
   ```
   Không tồn tại → 0 (khởi tạo) → Giá thực tế (từ oracle)
   ```

### Điểm quan trọng

1. **Thứ tự Message:** Thứ tự message là QUAN TRỌNG:
   - Oracle Market → Perpetual → CLOB Pair → DelayMessage
   - Thứ tự sai sẽ làm proposal thất bại

2. **Authority:**
   - Proposal messages: Authority = gov module
   - Delayed UpdateClobPair: Authority = delaymsg module
   - Validation xảy ra tại cả submission và execution time

3. **Delay Blocks:**
   - Có thể là 0, 1, hoặc bất kỳ số nào
   - Cho phép hệ thống có thời gian setup trước khi kích hoạt

4. **Atomic Execution:**
   - Nếu một message thất bại, toàn bộ proposal thất bại
   - State được rollback về trước khi proposal execution

5. **Market Map Integration:**
   - Market phải được enable trong market map để trading
   - Chỉ được enable sau khi CLOB Pair = ACTIVE

### Lý do thiết kế

1. **An toàn:** Cơ chế delay đảm bảo market không được kích hoạt ngay lập tức, cho phép validation và setup.

2. **Quản lý Dependencies:** Thứ tự message đảm bảo dependencies được tạo đúng cách.

3. **Linh hoạt:** Delay blocks có thể được điều chỉnh khi cần.

4. **Xử lý Lỗi:** Atomic execution đảm bảo tính nhất quán state - không có partial state.

5. **Tích hợp:** Tích hợp market map đảm bảo market có thể được khám phá và trade.
