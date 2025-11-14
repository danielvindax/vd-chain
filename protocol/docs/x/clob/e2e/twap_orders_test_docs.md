# Tài liệu Test: TWAP Orders E2E Tests

## Tổng quan

File test này xác minh chức năng **TWAP (Time-Weighted Average Price) Order** trong CLOB module. TWAP orders chia một large order thành nhiều smaller suborders được thực thi theo thời gian ở regular intervals. Test đảm bảo rằng:
1. TWAP orders được chia thành suborders dựa trên duration và interval
2. Suborders được đặt ở regular intervals
3. Suborders sử dụng oracle price với price tolerance
4. TWAP orders catch up nếu suborders expire unfilled
5. Duplicate TWAP orders bị từ chối

---

## Test Function: TestTwapOrderPlacementAndCatchup

### Test Case: Thành công - TWAP Order Placement và Suborder Execution

### Đầu vào
- **TWAP Order:**
  - SubaccountId: Alice_Num0
  - Side: BUY
  - Quantums: 100,000,000,000 (10 BTC)
  - Duration: 300 giây (5 phút)
  - Interval: 60 giây (1 phút)
  - PriceTolerance: 0% (market order)
  - GoodTilBlockTime: 300 giây từ bây giờ

### Đầu ra
- **TWAP Order Placement:**
  - RemainingLegs: 4 (5 total - 1 triggered)
  - RemainingQuantums: 100,000,000,000
- **First Suborder:**
  - Quantums: 20,000,000,000 (100B / 5 = 20B per leg)
  - Subticks: 200,000,000 ($20,000 oracle price)
  - GoodTilBlockTime: 3 giây từ bây giờ
  - Side: BUY (giống parent)
- **After 30 Seconds:**
  - Suborder expired và được xóa
  - TWAP order vẫn có 4 remaining legs
- **After 60 Seconds Total:**
  - Second suborder được đặt
  - Quantums: 25,000,000,000 (100B / 4 = 25B, catching up)
  - Subticks: 200,000,000 (cùng oracle price)

### Tại sao chạy theo cách này?

1. **Suborder Calculation:** Total quantums chia cho số lượng legs.
   - 5 legs: 100B / 5 = 20B per leg
   - Sau 1 leg: 100B / 4 = 25B per leg (catchup)
2. **Interval-Based:** Suborders được đặt ở regular intervals (60 giây).
3. **Oracle Price:** Suborders sử dụng current oracle price.
4. **Catchup Logic:** Nếu suborder expire unfilled, next suborder có size lớn hơn để catch up.

---

## Test Function: TestDuplicateTWAPOrderPlacement

### Test Case: Thất bại - Duplicate TWAP Order

### Đầu vào
- **Block 1:**
  - Đặt TWAP order: Alice mua 100B quantums qua 4 legs
- **Block 2:**
  - Cố gắng đặt cùng TWAP order (cùng OrderId)

### Đầu ra
- **First Order CheckTx:** SUCCESS
- **Second Order CheckTx:** FAIL
- **Error:** "A stateful order with this OrderId already exists"

### Tại sao chạy theo cách này?

1. **Duplicate Detection:** Hệ thống phát hiện duplicate OrderId.
2. **Stateful Order:** TWAP orders là stateful orders.
3. **Rejection:** Không thể đặt duplicate stateful order.

---

## Test Function: TestTWAPOrderWithMatchingOrders

### Test Case: Thành công - TWAP Suborder Match với Existing Order

### Đầu vào
- **TWAP Order:** Alice mua 100B quantums qua 4 legs
- **Existing Order:** Bob bán matching quantity ở compatible price
- **Suborder:** First suborder được đặt và match với order của Bob

### Đầu ra
- **Suborder:** Fully filled
- **TWAP Order:** Remaining quantums giảm
- **Next Suborder:** Next suborder được đặt ở next interval

### Tại sao chạy theo cách này?

1. **Suborder Matching:** TWAP suborders có thể match như regular orders.
2. **Fill Tracking:** Fills giảm remaining quantums trong TWAP order.
3. **Continued Execution:** TWAP order tiếp tục đặt suborders cho đến khi hoàn thành.

---

## Tóm tắt Flow

### TWAP Order Execution Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. ĐẶT TWAP ORDER                                          │
│    - Chỉ định total quantums, duration, interval            │
│    - Tính số lượng legs = duration / interval                │
│    - Lưu TWAP order trong state                              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. ĐẶT FIRST SUBORDER                                       │
│    - Tính suborder size = total / num_legs                  │
│    - Lấy current oracle price                               │
│    - Áp dụng price tolerance                                │
│    - Đặt suborder trên order book                            │
│    - Lên lịch next suborder trigger                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. SUBORDER EXECUTION                                       │
│    - Suborder có thể match với existing orders               │
│    - Nếu filled: Cập nhật remaining quantums                  │
│    - Nếu expired: Xóa và catch up                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. TRIGGER NEXT SUBORDER                                    │
│    - Tại next interval time                                  │
│    - Tính catchup size nếu previous expired                 │
│    - Đặt next suborder                                       │
│    - Lặp lại cho đến khi tất cả legs được thực thi           │
└─────────────────────────────────────────────────────────────┘
```

### Catchup Logic

```
Nếu previous suborder expired unfilled:
  Remaining quantums = original remaining - 0 (không có gì filled)
  Remaining legs = original legs - 1
  Next suborder size = remaining quantums / remaining legs
  
Ví dụ:
  Original: 100B quantums, 5 legs
  First suborder: 20B (expires unfilled)
  Catchup: 100B / 4 = 25B per remaining leg
```

### Điểm quan trọng

1. **TWAP Parameters:**
   - Duration: Tổng thời gian để thực thi order
   - Interval: Thời gian giữa các suborders
   - PriceTolerance: Maximum price deviation từ oracle
   - Số lượng legs = duration / interval

2. **Suborder Calculation:**
   - Initial: total_quantums / num_legs
   - Catchup: remaining_quantums / remaining_legs
   - Đảm bảo tất cả quantums được thực thi đến cuối duration

3. **Oracle Price:**
   - Suborders sử dụng current oracle price
   - Price tolerance cho phép deviation
   - Market orders: tolerance = 0

4. **Suborder Lifecycle:**
   - Được đặt tại interval time
   - Có thể match với existing orders
   - Expire nếu không fill bởi GoodTilBlockTime
   - Được xóa và next suborder catch up

5. **State Tracking:**
   - TWAP order placement được track trong state
   - Remaining legs và quantums được track
   - Trigger times được lên lịch cho next suborders

6. **Duplicate Prevention:**
   - Không thể đặt duplicate TWAP order
   - Cùng OrderId bị từ chối
   - Ngăn chặn accidental duplicate execution

### Lý do thiết kế

1. **Price Impact Reduction:** Chia large orders giảm market impact.

2. **Time Distribution:** Thực thi order theo thời gian cho better average price.

3. **Flexibility:** Configurable duration và interval cho các strategies khác nhau.

4. **Catchup Logic:** Đảm bảo order hoàn thành ngay cả khi một số suborders expire.

5. **Oracle Integration:** Sử dụng oracle price cho fair execution.
