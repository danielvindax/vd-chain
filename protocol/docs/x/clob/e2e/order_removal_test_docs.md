# Tài liệu Test: Order Removal E2E Tests

## Tổng quan

File test này xác minh chức năng **Order Removal** trong CLOB module. Orders có thể được xóa khỏi order book vì nhiều lý do. Test đảm bảo rằng:
1. Conditional orders được xóa khi chúng cross maker orders (PostOnly violation)
2. Conditional IOC orders được xóa nếu không fill đầy
3. Self-trading xóa maker orders
4. Fully filled orders được xóa
5. Under-collateralized orders được xóa

---

## Test Function: TestConditionalOrderRemoval

### Test Case 1: Conditional PostOnly Order Crosses Maker - Removed

### Đầu vào
- **Subaccounts:**
  - Alice: 10,000 USD
  - Bob: 10,000 USD
- **Orders:**
  - Long-term order: Alice mua 5 ở giá 10 (maker)
  - Conditional order: Bob bán 10 ở giá 10, PostOnly, StopLoss trigger ở 15
- **Price Update:** Giá tăng lên 14.9 (trigger conditional order)

### Đầu ra
- **Alice Order:** Không bị xóa (maker order)
- **Bob Order:** Bị xóa (PostOnly violation - crosses maker)

### Tại sao chạy theo cách này?

1. **PostOnly Violation:** Khi conditional order trigger, nó sẽ cross existing maker order.
2. **PostOnly Rule:** PostOnly orders không thể cross existing orders, phải là maker.
3. **Removal:** Conditional order được xóa thay vì cross.

---

### Test Case 2: Conditional IOC Order Không Fill Đầy - Removed

### Đầu vào
- **Subaccounts:**
  - Carl: 10,000 USD
  - Dave: 10,000 USD
- **Orders:**
  - Long-term order: Dave bán 0.25 BTC ở 50,000
  - Conditional order: Carl mua 0.5 BTC ở 50,000, IOC, StopLoss trigger ở 50,003
- **Price Update:** Giá tăng lên 50,004 (trigger conditional order)

### Đầu ra
- **Dave Order:** Bị xóa (fully filled)
- **Carl Order:** Bị xóa (IOC không fill đầy)

### Tại sao chạy theo cách này?

1. **IOC Rule:** Immediate-Or-Cancel orders phải được fill đầy ngay lập tức hoặc bị hủy.
2. **Partial Fill:** Chỉ có 0.25 BTC available, nhưng order muốn 0.5 BTC.
3. **Removal:** IOC order được xóa vì không thể fill đầy.

---

### Test Case 3: Conditional Self Trade - Xóa Maker Order

### Đầu vào
- **Subaccount:** Alice với 10,000 USD
- **Orders:**
  - Long-term order: Alice mua 5 ở giá 10
  - Conditional order: Alice bán 20 ở giá 10, StopLoss trigger ở 15
- **Price Update:** Giá tăng lên 14.9 (trigger conditional order)

### Đầu ra
- **Long-term Order:** Bị xóa (self-trade xóa maker)
- **Conditional Order:** Không bị xóa (taker trong self-trade)

### Tại sao chạy theo cách này?

1. **Self-Trade:** Cùng subaccount có cả maker và taker orders.
2. **Maker Removal:** Self-trading xóa maker order để ngăn chặn abuse.
3. **Taker Kept:** Taker order (conditional) được giữ lại.

---

### Test Case 4: Fully Filled Maker Orders - Removed

### Đầu vào
- **Subaccounts:**
  - Alice: 10,000 USD
  - Bob: 10,000 USD
- **Orders:**
  - Long-term order: Alice mua 5 ở giá 10
  - Conditional order: Bob bán 50 ở giá 10, StopLoss trigger ở 15
- **Price Update:** Giá tăng lên 14.9 (trigger conditional order)

### Đầu ra
- **Alice Order:** Bị xóa (fully filled bởi conditional order)
- **Bob Order:** Không bị xóa (partially filled, 45 còn lại)

### Tại sao chạy theo cách này?

1. **Full Fill:** Conditional order fill đầy maker order (5 units).
2. **Maker Removal:** Fully filled maker order được xóa.
3. **Taker Partial:** Conditional order partially filled, vẫn trên book.

---

### Test Case 5: Under-Collateralized Conditional Taker - Removed

### Đầu vào
- **Subaccounts:**
  - Carl: 100,000 USD
  - Dave: 10,000 USD
- **Orders:**
  - Long-term order: Carl mua 1 BTC ở 50,000
  - Conditional order: Dave bán 1 BTC ở 50,000, StopLoss trigger ở 50,003
- **Withdrawal:** Dave rút 10,000 USD (trở nên under-collateralized)
- **Price Update:** Giá tăng lên 50,002.5 (trigger conditional order)

### Đầu ra
- **Carl Order:** Không bị xóa
- **Dave Order:** Bị xóa (thất bại collateralization check trong quá trình matching)

### Tại sao chạy theo cách này?

1. **Collateralization Check:** Khi conditional order trigger và cố gắng match, hệ thống kiểm tra collateral.
2. **Insufficient Collateral:** Dave không có đủ collateral sau withdrawal.
3. **Removal:** Order được xóa thay vì thực thi trade.

---

## Tóm tắt Flow

### Order Removal Reasons

1. **PostOnly Violation:**
   - Order sẽ cross existing maker
   - PostOnly orders phải là maker
   - Order được xóa thay vì cross

2. **IOC Not Fully Filled:**
   - IOC order không thể fill đầy ngay lập tức
   - IOC orders phải fill hoàn toàn hoặc bị hủy
   - Order được xóa

3. **Self-Trade:**
   - Cùng subaccount có maker và taker orders
   - Maker order được xóa để ngăn chặn abuse
   - Taker order được giữ lại

4. **Fully Filled:**
   - Order được fill hoàn toàn bởi matching
   - Fully filled orders được xóa khỏi book
   - State được cập nhật để reflect fill

5. **Under-Collateralized:**
   - Order thất bại collateralization check
   - Insufficient margin cho position
   - Order được xóa trước khi execution

### Điểm quan trọng

1. **Removal Timing:**
   - Orders được xóa trong DeliverTx
   - Removal xảy ra trước state update
   - Events được emit cho removed orders

2. **Removal Reasons:**
   - Được track trong order removal events
   - Các lý do khác nhau cho các scenarios khác nhau
   - Được sử dụng cho off-chain tracking

3. **State Consistency:**
   - Removed orders không có trong state
   - Không thể query removed orders
   - Fill amounts được track trước khi removal

4. **Event Emission:**
   - Order removal events được emit
   - Bao gồm removal reason
   - Được sử dụng bởi indexer cho off-chain sync

5. **Collateralization:**
   - Được kiểm tra khi order cố gắng match
   - Phải có sufficient margin
   - Under-collateralized orders được xóa

### Lý do thiết kế

1. **Order Book Integrity:** Removal ngăn chặn invalid orders ở lại trên book.

2. **Risk Management:** Under-collateralized orders được xóa để ngăn chặn bad trades.

3. **Fairness:** Self-trade removal ngăn chặn manipulation.

4. **Efficiency:** Removal của invalid orders giữ book sạch.
