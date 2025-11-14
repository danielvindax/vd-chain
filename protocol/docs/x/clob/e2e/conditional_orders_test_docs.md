# Tài liệu Test: Conditional Orders E2E Tests

## Tổng quan

File test này xác minh chức năng **Conditional Order** trong CLOB module. Conditional orders là các orders được trigger khi các điều kiện giá nhất định được đáp ứng. Test đảm bảo rằng:
1. Conditional orders được đặt nhưng không được trigger khi conditions không được đáp ứng
2. Conditional orders được trigger khi price conditions được đáp ứng
3. Các loại trigger khác nhau (TakeProfit, StopLoss) hoạt động đúng
4. Conditional orders có thể match với existing orders khi được trigger

---

## Test Function: TestConditionalOrder

### Test Case 1: TakeProfit/Buy - Không Trigger (Không có Price Update)

### Đầu vào
- **Subaccount:** Alice với 100,000 USD
- **Order:** Conditional Buy 1 BTC ở giá 50,000, TakeProfit trigger ở 49,999
- **Price Updates:** Không có

### Đầu ra
- **Order State:** Tồn tại trong state, không được trigger
- **Triggered State:** false sau tất cả blocks

### Tại sao chạy theo cách này?

1. **TakeProfit Logic:** TakeProfit/Buy trigger khi giá xuống dưới trigger price.
2. **No Price Update:** Không có price update, trigger condition không bao giờ được đáp ứng.
3. **Order Persists:** Order vẫn trong state chờ trigger condition.

---

### Test Case 2: StopLoss/Buy - Không Trigger (Không có Price Update)

### Đầu vào
- **Subaccount:** Alice với 100,000 USD
- **Order:** Conditional Buy 1 BTC ở giá 50,000, StopLoss trigger ở 50,001
- **Price Updates:** Không có

### Đầu ra
- **Order State:** Tồn tại trong state, không được trigger
- **Triggered State:** false sau tất cả blocks

### Tại sao chạy theo cách này?

1. **StopLoss Logic:** StopLoss/Buy trigger khi giá lên trên trigger price.
2. **No Price Update:** Không có price update, trigger condition không bao giờ được đáp ứng.
3. **Order Persists:** Order vẫn trong state.

---

### Test Case 3: TakeProfit/Buy - Triggered bởi Price Update

### Đầu vào
- **Subaccount:** Alice với 100,000 USD
- **Order:** Conditional Buy 1 BTC ở giá 50,000, TakeProfit trigger ở 49,999
- **Price Update:** Giá giảm xuống 49,997 (dưới trigger)

### Đầu ra
- **Order State:** Triggered và được đặt trên order book
- **Triggered State:** true sau price update
- **Order:** Có thể match với existing orders

### Tại sao chạy theo cách này?

1. **Price Condition Met:** Giá (49,997) < trigger (49,999), condition được đáp ứng.
2. **Order Triggered:** Conditional order trở thành active order.
3. **Matching:** Triggered order có thể match với existing orders.

---

### Test Case 4: StopLoss/Buy - Triggered bởi Price Update

### Đầu vào
- **Subaccount:** Alice với 100,000 USD
- **Order:** Conditional Buy 1 BTC ở giá 50,000, StopLoss trigger ở 50,001
- **Price Update:** Giá tăng lên 50,003 (trên trigger)

### Đầu ra
- **Order State:** Triggered và được đặt trên order book
- **Triggered State:** true sau price update

### Tại sao chạy theo cách này?

1. **Price Condition Met:** Giá (50,003) > trigger (50,001), condition được đáp ứng.
2. **Order Triggered:** Conditional order trở thành active order.

---

### Test Case 5: TakeProfit/Sell - Triggered và Matched

### Đầu vào
- **Subaccounts:**
  - Bob: 100,000 USD, 1 BTC long
  - Alice: 100,000 USD
- **Orders:**
  - Long-term order: Bob bán 1 BTC ở 50,000
  - Conditional order: Alice mua 1 BTC ở 50,000, TakeProfit trigger ở 49,999
- **Price Update:** Giá giảm xuống 49,997

### Đầu ra
- **Conditional Order:** Triggered
- **Orders Matched:** Cả hai orders match đầy đủ
- **Positions:** Alice có 1 BTC long, Bob có 0 BTC

### Tại sao chạy theo cách này?

1. **Trigger Condition:** Giá giảm xuống dưới trigger, order được trigger.
2. **Immediate Matching:** Triggered order match với existing order trên book.
3. **Trade Execution:** Cả hai orders được fill, positions được cập nhật.

---

## Tóm tắt Flow

### Conditional Order Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│ 1. ĐẶT CONDITIONAL ORDER                                    │
│    - Order được đặt trong conditional state                  │
│    - Trigger condition được chỉ định                          │
│    - Order không active trên order book                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CHỜ TRIGGER                                               │
│    - Theo dõi price updates                                   │
│    - Kiểm tra trigger condition mỗi block                   │
│    - Order vẫn trong conditional state                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. TRIGGER CONDITION ĐƯỢC ĐÁP ỨNG                           │
│    - Price update đáp ứng trigger condition                  │
│    - Order chuyển sang active state                           │
│    - Order được đặt trên order book                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. ORDER ACTIVE                                             │
│    - Order có thể match với existing orders                  │
│    - Order hoạt động như regular order                        │
│    - Có thể được fill, hủy, hoặc expire                     │
└─────────────────────────────────────────────────────────────┘
```

### Trigger Types

1. **TakeProfit/Buy:**
   - Trigger khi giá < trigger price
   - Được sử dụng để mua ở giá thấp hơn (profit taking)

2. **StopLoss/Buy:**
   - Trigger khi giá > trigger price
   - Được sử dụng để mua ở giá cao hơn (stop loss protection)

3. **TakeProfit/Sell:**
   - Trigger khi giá > trigger price
   - Được sử dụng để bán ở giá cao hơn (profit taking)

4. **StopLoss/Sell:**
   - Trigger khi giá < trigger price
   - Được sử dụng để bán ở giá thấp hơn (stop loss protection)

### Điểm quan trọng

1. **Conditional State:**
   - Orders bắt đầu trong conditional state
   - Không active trên order book cho đến khi triggered
   - Không thể match cho đến khi triggered

2. **Trigger Conditions:**
   - Được kiểm tra mỗi block sau price updates
   - Giá phải vượt qua trigger price để activate
   - Một khi triggered, order trở thành active

3. **Price Updates:**
   - Oracle price updates trigger condition checks
   - Price updates đến từ price feed
   - Nhiều price updates có thể xảy ra mỗi block

4. **Order Matching:**
   - Một khi triggered, order hoạt động như regular order
   - Có thể match ngay lập tức nếu compatible order tồn tại
   - Có thể vẫn trên book nếu không có match

5. **State Tracking:**
   - Triggered state được track theo từng order
   - Order state transitions: conditional → triggered → filled/cancelled

### Lý do thiết kế

1. **Risk Management:** Conditional orders cho phép users đặt automatic orders dựa trên price movements.

2. **Flexibility:** Các loại trigger khác nhau hỗ trợ nhiều trading strategies.

3. **Efficiency:** Orders chỉ trở thành active khi conditions được đáp ứng, giảm order book clutter.

4. **Safety:** Trigger conditions ngăn chặn accidental order execution.
