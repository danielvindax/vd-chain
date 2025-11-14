# Tài liệu Test: CLOB App E2E Tests

## Tổng quan

File test này xác minh chức năng **CLOB application** cốt lõi bao gồm order hydration, concurrent operations, transaction validation, và statistics tracking. Các tests đảm bảo rằng:
1. Long-term orders được hydrate đúng từ state sang memclob khi khởi động
2. Orders có thể được match trong quá trình hydration
3. Concurrent matches và cancels hoạt động đúng
4. Transaction signature validation hoạt động đúng
5. Statistics được track đúng cho trading activity

---

## Test Function: TestHydrationInPreBlocker

### Test Case: Thành công - Hydrate Long-Term Order từ State

### Đầu vào
- **Genesis State:**
  - Long-term order tồn tại trong state (không có trong memclob)
  - Order: Carl_Num0, Mua 1 BTC ở 50,000 USDC, GoodTilTime = 10
  - Order được đặt ở block 1
  - Order expiration: time.Unix(50, 0)
- **Block:** Tiến đến block 2

### Đầu ra
- **Order in State:** Order tồn tại trong state storage
- **Order in MemClob:** Order được hydrate và tồn tại trong memclob
- **Order on Orderbook:** Order hiển thị trên orderbook

### Tại sao chạy theo cách này?

1. **State Hydration:** Long-term orders được lưu trong state nhưng cần được load vào memclob khi khởi động.
2. **PreBlocker:** PreBlocker được gọi trước mỗi block để hydrate orders từ state.
3. **Order Visibility:** Orders phải ở trong memclob để hiển thị trên orderbook.
4. **Persistence:** Orders tồn tại qua restarts, vì vậy hydration là quan trọng.

---

## Test Function: TestHydrationWithMatchPreBlocker

### Test Case: Thành công - Hydrate Orders Match Trong Quá trình Hydration

### Đầu vào
- **Genesis State:**
  - Carl: Long-term buy order (1 BTC ở 50,000 USDC)
  - Dave: Long-term sell order (1 BTC ở 50,000 USDC)
  - Cả hai orders được đặt ở block 1
  - Cả hai orders expire ở time.Unix(10, 0)
- **Block:** Tiến đến block 2

### Đầu ra
- **PreBlocker:** Orders được hydrate và match trong PreBlocker
- **State Changes Discarded:** State changes từ PreBlocker bị loại bỏ (IsCheckTx = true)
- **Operations Queue:** Match operation được thêm vào operations queue
- **Block 2:** Orders được fill đầy và xóa khỏi state
- **Final State:**
  - Carl: Long 1 BTC, USDC balance giảm 50,000 USDC
  - Dave: Short 1 BTC, USDC balance tăng 50,000 USDC

### Tại sao chạy theo cách này?

1. **Hydration Matching:** Orders match trong quá trình hydration nên tạo match operations.
2. **State Isolation:** PreBlocker chạy trong CheckTx context, vì vậy state changes bị loại bỏ.
3. **Operations Queue:** Matches được queue và xử lý trong block tiếp theo.
4. **Full Fill:** Matching orders được fill đầy và xóa khỏi state.

---

## Test Function: TestConcurrentMatchesAndCancels

### Test Case: Thành công - Concurrent Matches và Cancels

### Đầu vào
- **Accounts:** 1000 tài khoản ngẫu nhiên
- **Orders:**
  - 300 orders match (150 buys, 150 sells)
    - 50 orders mỗi size 5, 10, 15 cho cả hai phía
    - Tổng matched volume: 1,500 quantums
  - 700 orders bị hủy
    - Orders được đặt rồi ngay lập tức hủy
- **Execution:** Tất cả CheckTx calls được thực thi đồng thời
- **Block:** Tiến đến block 3

### Đầu ra
- **Matched Orders:** Tất cả 300 orders được fill đầy
- **Cancelled Orders:** Tất cả 700 orders bị hủy (không fill)
- **No Data Races:** Test pass với flag `-race` được bật

### Tại sao chạy theo cách này?

1. **Concurrency Testing:** Test rằng hệ thống xử lý concurrent operations đúng cách.
2. **Race Detection:** Sử dụng Go's race detector để tìm data races.
3. **Mixed Operations:** Test cả matches và cancels xảy ra đồng thời.
4. **Stress Test:** 1000 accounts với concurrent operations stress test hệ thống.

---

## Test Function: TestFailsDeliverTxWithIncorrectlySignedPlaceOrderTx

### Test Case: Thất bại - Incorrectly Signed Order Placement

### Đầu vào
- **Order:** Order của Alice (từ Alice_Num0)
- **Signer:** Private key của Bob (signer sai)
- **Transaction:** Order placement transaction được ký bởi Bob

### Đầu ra
- **DeliverTx:** FAIL
- **Error:** "invalid pubkey: MsgProposedOperations is invalid"
- **Transaction:** Bị từ chối

### Tại sao chạy theo cách này?

1. **Signature Validation:** Transactions phải được ký bởi account đúng.
2. **Bảo mật:** Ngăn chặn order placement trái phép.
3. **DeliverTx Validation:** Validation xảy ra trong DeliverTx, không chỉ CheckTx.

---

## Test Function: TestFailsDeliverTxWithUnsignedTransactions

### Test Case: Thất bại - Unsigned Order Placement

### Đầu vào
- **Order:** Order của Alice (từ Alice_Num0)
- **Transaction:** Order placement transaction không có signatures

### Đầu ra
- **DeliverTx:** FAIL
- **Error:** "Error: no signatures supplied: MsgProposedOperations is invalid"
- **Transaction:** Bị từ chối

### Tại sao chạy theo cách này?

1. **Signature Requirement:** Tất cả transactions phải được ký.
2. **Bảo mật:** Ngăn chặn unsigned transactions được xử lý.
3. **Validation:** Signature validation xảy ra trong DeliverTx.

---

## Test Function: TestStats

### Test Case: Thành công - Statistics Tracking

### Đầu vào
- **Epochs:** Nhiều epochs với trading activity
- **Orders:**
  - Block 2-5: Alice (maker) và Bob (taker) trade 10,000 notional
  - Block 6: Alice và Bob trade 5,000 notional (cùng epoch)
  - Block 8: Alice và Bob trade 5,000 notional (epoch mới)
- **Time:** Epochs tiến dựa trên StatsEpochDuration

### Đầu ra
- **User Stats:**
  - Alice: MakerNotional = 20,000, TakerNotional = 0
  - Bob: MakerNotional = 0, TakerNotional = 20,000
- **Global Stats:** NotionalTraded = 20,000
- **Epoch Stats:**
  - Epoch 0: 15,000 notional
  - Epoch 2: 5,000 notional
- **Window Expiration:** Stats expire sau window duration

### Tại sao chạy theo cách này?

1. **Statistics Tracking:** Theo dõi trading activity cho rewards và analytics.
2. **Maker/Taker:** Phân biệt giữa maker và taker roles.
3. **Epochs:** Statistics được track theo epoch.
4. **Window:** Statistics expire sau window duration.
5. **Aggregation:** User stats, global stats, và epoch stats đều được track.

---

## Tóm tắt Flow

### Order Hydration Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. GENESIS STATE                                             │
│    - Long-term orders được lưu trong state                  │
│    - Orders không có trong memclob                           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. PREBLOCKER                                                │
│    - Load orders từ state                                   │
│    - Hydrate orders vào memclob                              │
│    - Kiểm tra matches                                        │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. ORDERBOOK                                                 │
│    - Orders hiển thị trên orderbook                          │
│    - Orders có thể được match                                │
└─────────────────────────────────────────────────────────────┘
```

### Concurrent Operations Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. CONCURRENT CHECKTX                                        │
│    - Nhiều goroutines thực thi CheckTx                       │
│    - Orders được đặt và hủy đồng thời                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. BLOCK ADVANCEMENT                                         │
│    - Tất cả transactions được bao gồm trong block            │
│    - Matches và cancels được xử lý                           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. VERIFICATION                                              │
│    - Matched orders được fill đầy                            │
│    - Cancelled orders không fill                             │
│    - Không phát hiện data races                              │
└─────────────────────────────────────────────────────────────┘
```

### Statistics Tracking Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. ORDER MATCHING                                            │
│    - Orders match trong block                                │
│    - Maker và taker được xác định                           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. STATISTICS UPDATE                                         │
│    - User stats được cập nhật                                │
│    - Global stats được cập nhật                              │
│    - Epoch stats được cập nhật                               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. EPOCH ADVANCEMENT                                         │
│    - Epoch mới bắt đầu                                       │
│    - Previous epoch stats được giữ lại                      │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. WINDOW EXPIRATION                                         │
│    - Old epoch stats expire                                 │
│    - Stats được xóa khỏi window                             │
└─────────────────────────────────────────────────────────────┘
```

### Điểm quan trọng

1. **Order Hydration:**
   - Long-term orders phải được hydrate từ state khi khởi động
   - Hydration xảy ra trong PreBlocker
   - Orders trở nên hiển thị trên orderbook sau hydration

2. **Concurrent Operations:**
   - Hệ thống phải xử lý concurrent CheckTx calls
   - Race detector giúp tìm data races
   - Stress testing với nhiều accounts

3. **Transaction Validation:**
   - Signatures phải hợp lệ
   - Signer phải khớp với transaction sender
   - Validation xảy ra trong DeliverTx

4. **Statistics:**
   - Maker/Taker distinction là quan trọng
   - Statistics được track theo epoch
   - Statistics expire sau window duration

### Lý do thiết kế

1. **State Hydration:** Đảm bảo orders tồn tại qua restarts và có sẵn cho matching.

2. **Concurrency:** Hệ thống phải xử lý high concurrency trong production.

3. **Bảo mật:** Signature validation ngăn chặn transactions trái phép.

4. **Analytics:** Statistics cho phép rewards và analytics features.
