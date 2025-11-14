# Tài liệu Test: Order Matches E2E Tests

## Tổng quan

File test này xác minh **Order Matching Validation** trong CLOB module. Test đảm bảo rằng order matching operations được validate đúng trong DeliverTx, bao gồm:
1. IOC orders không thể được match hai lần
2. Partially filled conditional IOC orders không thể được match lại
3. IOC orders có thể match với nhiều makers trong single operation
4. IOC orders không thể là taker trong nhiều separate matches

---

## Test Function: TestDeliverTxMatchValidation

### Test Case 1: Thất bại - Partially Filled IOC Taker Order Không thể Match Hai Lần

### Đầu vào
- **Block 1:**
  - Đặt order: Bob mua 5 ở giá 40
  - Đặt order: Alice bán 10 ở giá 15, IOC
  - Match operation: Alice (IOC taker) match với Bob (maker), fill 5
- **Block 2:**
  - Đặt order: Bob mua 5 ở giá 40
  - Đặt order: Alice bán 10 ở giá 15, IOC (cùng order)
  - Match operation: Cố gắng match Alice IOC order lại

### Đầu ra
- **Block 1 DeliverTx:** SUCCESS
- **Block 2 DeliverTx:** FAIL với lỗi "IOC order is already filled, remaining size is cancelled."

### Tại sao chạy theo cách này?

1. **IOC Rule:** IOC orders phải được fill đầy ngay lập tức hoặc bị hủy.
2. **Partial Fill:** Order được partially fill trong block 1 (5 trong 10).
3. **Remaining Cancelled:** Remaining size (5) được hủy sau partial fill.
4. **Cannot Reuse:** Không thể match cùng IOC order lại trong block sau.

---

### Test Case 2: Thất bại - Không thể Match Partially Filled Conditional IOC Order

### Đầu vào
- **Block 1:**
  - Đặt conditional IOC order: Alice mua 1 BTC ở 50,000, TakeProfit trigger ở 49,999
  - Đặt long-term order: Dave bán 0.25 BTC ở 50,000
- **Block 2:**
  - Conditional order trigger và partially match (0.25 BTC filled)
- **Block 3:**
  - Đặt order: Dave bán 1 BTC ở 50,000
  - Match operation: Cố gắng match conditional IOC order lại

### Đầu ra
- **Block 2 DeliverTx:** SUCCESS (partial fill)
- **Block 3 DeliverTx:** FAIL với lỗi `ErrStatefulOrderDoesNotExist`

### Tại sao chạy theo cách này?

1. **Conditional IOC:** Conditional IOC order được partially fill trong block 2.
2. **Order Removed:** Sau partial fill, conditional IOC order được xóa khỏi state.
3. **Cannot Match:** Không thể match order không còn tồn tại trong state.

---

### Test Case 3: Thành công - IOC Order Match với Nhiều Makers trong Single Operation

### Đầu vào
- **Orders:**
  - Bob mua 5 ở giá 40 (maker 1)
  - Bob mua 5 ở giá 40 (maker 2)
  - Alice bán 10 ở giá 15, IOC (taker)
- **Match Operation:**
  - Alice IOC order match với cả hai Bob orders
  - Fill 5 từ maker 1, fill 5 từ maker 2

### Đầu ra
- **DeliverTx:** SUCCESS
- **Result:** Alice order fill đầy (10), cả hai Bob orders fill đầy (5 mỗi)

### Tại sao chạy theo cách này?

1. **Multiple Makers:** IOC order có thể match với nhiều maker orders trong single operation.
2. **Full Fill:** Order được fill đầy bằng cách kết hợp fills từ nhiều makers.
3. **Single Operation:** Tất cả matches xảy ra trong một match operation.

---

### Test Case 4: Thất bại - IOC Order Không thể Là Taker trong Nhiều Matches

### Đầu vào
- **Orders:**
  - Bob mua 5 ở giá 40 (maker 1)
  - Alice bán 10 ở giá 15, IOC (taker)
  - Bob mua 5 ở giá 40 (maker 2)
- **Match Operations:**
  - Match 1: Alice IOC với maker 1 (fill 5)
  - Match 2: Alice IOC với maker 2 (fill 5)

### Đầu ra
- **DeliverTx:** FAIL với lỗi "IOC order is already filled, remaining size is cancelled."

### Tại sao chạy theo cách này?

1. **Multiple Matches:** IOC order không thể là taker trong nhiều separate match operations.
2. **First Match:** Sau match đầu tiên, order được coi là filled/cancelled.
3. **Second Match Fails:** Không thể sử dụng cùng IOC order trong match operation thứ hai.

---

## Tóm tắt Flow

### IOC Order Matching Rules

1. **Single Match Operation:**
   - IOC order có thể match với nhiều makers trong một operation
   - Tất cả fills xảy ra atomically
   - Order fill đầy hoặc bị hủy

2. **No Multiple Matches:**
   - IOC order không thể là taker trong nhiều separate operations
   - Sau match đầu tiên, order được fill/cancel
   - Subsequent matches thất bại

3. **Partial Fill Handling:**
   - Nếu IOC order partially fill, remaining size được hủy
   - Order được xóa khỏi state sau partial fill
   - Không thể match lại trong blocks sau

### Điểm quan trọng

1. **IOC Order Behavior:**
   - Phải fill hoàn toàn ngay lập tức hoặc bị hủy
   - Không thể tồn tại qua blocks nếu partially filled
   - Remaining size được hủy sau partial fill

2. **Match Operation:**
   - Single match operation có thể bao gồm nhiều maker fills
   - Tất cả fills xảy ra atomically
   - Order state được cập nhật sau operation

3. **State Management:**
   - Partially filled IOC orders được xóa khỏi state
   - Không thể query hoặc match removed orders
   - State consistency được duy trì

4. **Validation:**
   - DeliverTx validate match operations
   - Kiểm tra order existence và state
   - Từ chối invalid matches

### Lý do thiết kế

1. **Immediate Execution:** IOC orders đảm bảo immediate execution hoặc cancellation.

2. **State Consistency:** Ngăn chặn matching orders không còn tồn tại.

3. **Atomic Operations:** Single match operation đảm bảo atomic fills.

4. **Safety:** Validation ngăn chặn invalid match operations.
