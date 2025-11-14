# Tài liệu Test: Wind Down Market Proposal

## Tổng quan

File test này xác minh chức năng **Wind Down Market** (đóng market) thông qua governance proposals. Test đảm bảo rằng khi một CLOB pair chuyển sang `STATUS_FINAL_SETTLEMENT`, hệ thống sẽ:
1. Hủy tất cả stateful orders đang mở
2. Thực hiện final settlement deleveraging
3. Chặn đặt lệnh mới

---

## Test Case 1: Final Settlement Deleveraging - Tài khoản TNC Không Âm

### Đầu vào
- **Initial Subaccounts:**
  - `Carl_Num0`: 1 BTC Short position với 100,000 USD collateral
  - `Dave_Num0`: 1 BTC Long position với 50,000 USD collateral
- **ClobPair:** BTC-USD với trạng thái ban đầu `STATUS_ACTIVE`
- **Proposal:** Chuyển ClobPair sang `STATUS_FINAL_SETTLEMENT`

### Đầu ra
- **Subaccounts sau final settlement:**
  - `Carl_Num0`: Chỉ còn 50,000 USD (short position đóng ở oracle price)
  - `Dave_Num0`: Có 100,000 USD (long position đóng ở oracle price)
- **ClobPair status:** `STATUS_FINAL_SETTLEMENT`
- **Events:** Indexer events được emit cho ClobPair update

### Tại sao chạy theo cách này?

1. **Non-negative TNC (Total Net Collateral):** Cả hai subaccounts có TNC dương, có nghĩa là có đủ collateral cho settlement.
2. **Deleveraging at Oracle Price:** Vì cả hai accounts có TNC dương, deleveraging được thực hiện ở oracle price (reference price).
3. **Result:** 
   - Carl (short 1 BTC) phải trả Dave (long 1 BTC) một số tiền bằng oracle price
   - Nếu oracle price là 50,000 USD/BTC:
     - Carl ban đầu: 100,000 USD - phải trả 50,000 USD = còn lại 50,000 USD
     - Dave ban đầu: 50,000 USD + nhận 50,000 USD = 100,000 USD

---

## Test Case 2: Final Settlement Deleveraging - Tài khoản TNC Âm

### Đầu vào
- **Initial Subaccounts:**
  - `Carl_Num0`: 1 BTC Short position với 49,999 USD collateral (TNC âm)
  - `Dave_Num0`: 1 BTC Long position với 50,001 USD collateral
- **ClobPair:** BTC-USD với trạng thái ban đầu `STATUS_ACTIVE`
- **Proposal:** Chuyển ClobPair sang `STATUS_FINAL_SETTLEMENT`

### Đầu ra
- **Subaccounts sau final settlement:**
  - `Carl_Num0`: Rỗng (không còn gì do TNC âm)
  - `Dave_Num0`: Có 100,000 USD
- **ClobPair status:** `STATUS_FINAL_SETTLEMENT`

### Tại sao chạy theo cách này?

1. **Negative TNC:** Carl có TNC âm (49,999 USD < 50,000 USD cần để cover short position), có nghĩa là tài khoản này undercollateralized.
2. **Deleveraging at Bankruptcy Price:** Khi một account có TNC âm, deleveraging được thực hiện ở "bankruptcy price" - giá mà account không còn gì sau settlement.
3. **Result:**
   - Carl không có đủ funds để settle đầy đủ, vì vậy tất cả collateral của Carl (49,999 USD) được chuyển cho Dave
   - Dave nhận 49,999 USD từ Carl + 50,001 USD ban đầu = 100,000 USD
   - Carl mất tất cả và tài khoản trở thành rỗng

---

## Test Case 3: Hủy Stateful Orders Đang Mở

### Đầu vào
- **Subaccounts:**
  - `Alice_Num0`: 10,000 USD
  - `Bob_Num0`: 10,000 USD
- **Preexisting Stateful Orders:**
  - Long-term order từ Alice: Mua 5 units ở giá 5, GoodTilBlockTime = 5
  - Long-term order từ Bob: Bán 10 units ở giá 10, GoodTilBlockTime = 10, PostOnly
  - Conditional order từ Alice: Mua 1 BTC ở giá 50,000, GoodTilBlockTime = 10, StopLoss trigger = 50,001
- **Proposal:** Chuyển ClobPair sang `STATUS_FINAL_SETTLEMENT`

### Đầu ra
- **Stateful Orders:** Tất cả 3 orders được xóa khỏi state
- **Indexer Events:** Events được emit cho mỗi order bị xóa với lý do `ORDER_REMOVAL_REASON_FINAL_SETTLEMENT`
- **ClobPair status:** `STATUS_FINAL_SETTLEMENT`

### Tại sao chạy theo cách này?

1. **Stateful Orders:** Đây là các orders tồn tại qua nhiều blocks (long-term và conditional orders), khác với short-term orders chỉ tồn tại trong 1 block.
2. **Must Cancel When Wind Down:** Khi market đóng, tất cả pending orders phải được hủy vì chúng không thể được thực thi nữa.
3. **Events:** Indexer cần được thông báo về order cancellations để cập nhật off-chain state.
4. **Both Sides:** Test này đảm bảo orders từ cả hai phía (buy và sell) đều được hủy, bao gồm conditional orders.

---

## Test Case 4: Chặn Đặt Lệnh Mới

### Đầu vào
- **Subaccounts:**
  - `Alice_Num0`: 10,000 USD
  - `Bob_Num0`: 10,000 USD
  - `Carl_Num0`: 10,000 USD
- **Orders cố gắng đặt (sau wind down):**
  - Short-term orders: Mua 10, Bán 15, IOC order
  - Long-term orders: Mua 5, Bán 10
  - Conditional orders: Mua 1 BTC với stop loss
- **Proposal:** ClobPair đã được chuyển sang `STATUS_FINAL_SETTLEMENT`

### Đầu ra
- **Tất cả CheckTx responses:** FAIL với log chứa "trading is disabled for clob pair"
- **Không có orders được đặt:** Tất cả orders bị từ chối

### Tại sao chạy theo cách này?

1. **Protect Final Settlement State:** Khi market ở trạng thái final settlement, lệnh mới không được phép vì:
   - Market đang trong quá trình đóng
   - Chỉ final settlement deleveraging được phép
   - Không có hoạt động trading mới có thể xảy ra

2. **Validation at CheckTx:** Validation được thực hiện tại `CheckTx` để từ chối orders sớm, trước khi chúng vào mempool.

3. **All Order Types:** Test này đảm bảo short-term, long-term, và conditional orders đều bị chặn, bất kể loại order.

---

## Tóm tắt Flow

### Wind Down Market Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. KHỞI TẠO GENESIS STATE                                   │
│    - Tạo ClobPair với trạng thái ACTIVE                     │
│    - Tạo subaccounts với positions                          │
│    - Đặt stateful orders (nếu có)                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. SUBMIT GOVERNANCE PROPOSAL                                │
│    - Tạo MsgUpdateClobPair với trạng thái FINAL_SETTLEMENT │
│    - Submit proposal qua governance module                   │
│    - Validators vote (trong test: tất cả vote YES)          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. PROPOSAL EXECUTED                                         │
│    - Proposal status: PROPOSAL_STATUS_PASSED                │
│    - ClobPair status chuyển sang FINAL_SETTLEMENT          │
│    - Indexer events được emit                                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. HỦY STATEFUL ORDERS                                      │
│    - Tất cả long-term orders được xóa                       │
│    - Tất cả conditional orders được xóa                     │
│    - Indexer events được emit cho mỗi order bị xóa         │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. FINAL SETTLEMENT DELEVERAGING                            │
│    - Liquidations daemon cung cấp SubaccountOpenPositionInfo │
│    - Hệ thống xác định accounts cần deleverage              │
│    - Thực hiện deleveraging:                                 │
│      * TNC không âm → deleverage ở oracle price              │
│      * TNC âm → deleverage ở bankruptcy price               │
│    - Cập nhật subaccount balances                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 6. CHẶN LỆNH MỚI                                            │
│    - CheckTx validation từ chối tất cả lệnh mới             │
│    - Log: "trading is disabled for clob pair"                │
│    - Áp dụng cho tất cả loại orders                          │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **ClobPair Status Transition:**
   ```
   STATUS_ACTIVE → STATUS_FINAL_SETTLEMENT
   ```

2. **Subaccount State Changes:**
   - Positions được đóng (deleveraged)
   - Số dư được cập nhật dựa trên oracle/bankruptcy price
   - Tài khoản TNC âm có thể bị xóa sạch

3. **Order State:**
   - Stateful orders: Được xóa khỏi state
   - Lệnh mới: Bị từ chối tại CheckTx

### Điểm quan trọng

1. **Timing:** Final settlement deleveraging chỉ xảy ra sau khi proposal execution và liquidations daemon đã cung cấp thông tin position.

2. **Price Determination:**
   - **Oracle Price:** Được sử dụng cho accounts có TNC dương
   - **Bankruptcy Price:** Được sử dụng cho accounts có TNC âm (giá mà account không còn gì)

3. **Event Emission:** Indexer cần được thông báo về:
   - ClobPair status update
   - Stateful order removals
   - Để đồng bộ state off-chain

4. **Validation:** CheckTx validation đảm bảo không có lệnh mới nào có thể được đặt sau khi market đã được wind down.

### Lý do thiết kế

1. **An toàn:** Wind down market là quy trình quan trọng phải đảm bảo:
   - Tất cả positions được settle đúng cách
   - Không có hoạt động trading mới
   - State được cập nhật nhất quán

2. **Công bằng:** 
   - Tài khoản TNC không âm được settle ở oracle price (giá trị công bằng)
   - Tài khoản TNC âm được settle ở bankruptcy price (mất tất cả)

3. **Minh bạch:** Indexer events đảm bảo off-chain systems có thể theo dõi tất cả thay đổi.
