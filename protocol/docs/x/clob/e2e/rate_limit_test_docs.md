# Tài liệu Test: Rate Limiting E2E Tests

## Tổng quan

File test này xác minh chức năng **Rate Limiting** trong CLOB module. Rate limits giới hạn số lượng operations (orders, cancellations, leverage updates) một subaccount có thể thực hiện trong một số blocks được chỉ định. Test đảm bảo rằng:
1. Short-term orders được rate limit
2. Stateful orders được rate limit
3. Order cancellations được rate limit
4. Batch cancellations được rate limit
5. Leverage updates được rate limit
6. Rate limits áp dụng per subaccount

---

## Test Function: TestRateLimitingOrders_RateLimitsAreEnforced

### Test Case 1: Thất bại - Short-Term Orders với Cùng Subaccount Vượt Quá Limit

### Đầu vào
- **Rate Limit Config:**
  - MaxShortTermOrdersAndCancelsPerNBlocks: 1 order per 2 blocks
- **Block 2:**
  - Đặt order: Alice mua 5 ở giá 10, CLOB 0
- **Block 2:**
  - Cố gắng đặt order: Alice mua 5 ở giá 10, CLOB 1

### Đầu ra
- **First Order CheckTx:** SUCCESS
- **Second Order CheckTx:** FAIL với lỗi `ErrBlockRateLimitExceeded`
- **Error Message:** "exceeds configured block rate limit"

### Tại sao chạy theo cách này?

1. **Rate Limit:** 1 order per 2 blocks.
2. **First Order:** Tiêu thụ limit cho blocks 2-3.
3. **Second Order:** Cố gắng đặt trong cùng block, vượt quá limit.
4. **Rejection:** CheckTx từ chối order thứ hai ngay lập tức.

---

### Test Case 2: Thất bại - Short-Term Orders với Subaccounts Khác nhau Vượt Quá Limit

### Đầu vào
- **Rate Limit Config:** Giống Test Case 1
- **Block 2:**
  - Đặt order: Alice_Num0 mua 5 ở giá 10
- **Block 2:**
  - Cố gắng đặt order: Alice_Num1 mua 5 ở giá 10

### Đầu ra
- **First Order CheckTx:** SUCCESS
- **Second Order CheckTx:** FAIL với lỗi `ErrBlockRateLimitExceeded`

### Tại sao chạy theo cách này?

1. **Per Subaccount:** Rate limits áp dụng per subaccount.
2. **Different Subaccounts:** Alice_Num0 và Alice_Num1 là các subaccounts khác nhau.
3. **Still Limited:** Ngay cả các subaccounts khác nhau của cùng owner cũng bị rate limit.
4. **Owner-Based:** Rate limits có thể dựa trên owner address, không chỉ subaccount.

---

### Test Case 3: Thất bại - Stateful Orders Vượt Quá Limit

### Đầu vào
- **Rate Limit Config:**
  - MaxStatefulOrdersPerNBlocks: 1 order per 2 blocks
- **Block 2:**
  - Đặt long-term order: Alice mua 5 ở giá 10, CLOB 0
- **Block 2:**
  - Cố gắng đặt long-term order: Alice mua 5 ở giá 10, CLOB 1

### Đầu ra
- **First Order CheckTx:** SUCCESS
- **Second Order CheckTx:** FAIL với lỗi `ErrBlockRateLimitExceeded`

### Tại sao chạy theo cách này?

1. **Stateful Limit:** Limit riêng cho stateful orders.
2. **Same Limit Logic:** Hoạt động giống như short-term order limits.
3. **Per Subaccount:** Limits áp dụng per subaccount.

---

### Test Case 4: Thất bại - Order Cancellations Vượt Quá Limit

### Đầu vào
- **Rate Limit Config:**
  - MaxShortTermOrdersAndCancelsPerNBlocks: 1 operation per 2 blocks
- **Block 2:**
  - Cancel order: Alice hủy order trên CLOB 1
- **Block 2:**
  - Cố gắng cancel order: Alice hủy order trên CLOB 0

### Đầu ra
- **First Cancel CheckTx:** SUCCESS
- **Second Cancel CheckTx:** FAIL với lỗi `ErrBlockRateLimitExceeded`

### Tại sao chạy theo cách này?

1. **Cancellation Counts:** Cancellations tính vào cùng limit như orders.
2. **Combined Limit:** Orders và cancellations chia sẻ cùng rate limit.
3. **Same Logic:** Hoạt động giống như order placement limits.

---

### Test Case 5: Thất bại - Batch Cancellations Vượt Quá Limit

### Đầu vào
- **Rate Limit Config:**
  - MaxShortTermOrdersAndCancelsPerNBlocks: 2 operations per 2 blocks
- **Block 2:**
  - Batch cancel: Alice hủy 3 orders (tính là 1 operation)
- **Block 2:**
  - Cố gắng batch cancel: Alice hủy 3 orders

### Đầu ra
- **First Batch Cancel CheckTx:** SUCCESS
- **Second Batch Cancel CheckTx:** FAIL với lỗi `ErrBlockRateLimitExceeded`

### Tại sao chạy theo cách này?

1. **Batch Counts as One:** Batch cancel tính là 1 operation, không phải per order.
2. **Limit Exceeded:** Batch cancel thứ hai vượt quá limit của 2 per 2 blocks.
3. **Efficiency:** Batch operations hiệu quả hơn cho rate limits.

---

### Test Case 6: Thất bại - Leverage Updates Vượt Quá Limit

### Đầu vào
- **Rate Limit Config:**
  - MaxLeverageUpdatesPerNBlocks: 1 update per 2 blocks
- **Block 2:**
  - Update leverage: Alice cập nhật leverage cho perpetual 0 lên 5x
- **Block 2:**
  - Cố gắng update leverage: Alice cập nhật leverage cho perpetual 1 lên 10x

### Đầu ra
- **First Update CheckTx:** SUCCESS
- **Second Update CheckTx:** FAIL với lỗi `ErrBlockRateLimitExceeded`

### Tại sao chạy theo cách này?

1. **Separate Limit:** Leverage updates có rate limit riêng.
2. **Per Subaccount:** Limits áp dụng per subaccount.
3. **Prevents Spam:** Ngăn chặn excessive leverage update operations.

---

## Tóm tắt Flow

### Rate Limit Check Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. XÁC ĐỊNH LOẠI OPERATION                                  │
│    - Short-term order/cancel                                 │
│    - Stateful order                                          │
│    - Leverage update                                         │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. LẤY RATE LIMIT CONFIG                                    │
│    - Tìm limit cho operation type                            │
│    - Lấy NumBlocks và Limit                                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. ĐẾM RECENT OPERATIONS                                    │
│    - Đếm operations trong N blocks cuối                      │
│    - Bao gồm current block                                  │
│    - Đếm per subaccount                                     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. VALIDATE LIMIT                                            │
│    - Kiểm tra nếu count >= limit                            │
│    - Từ chối nếu vượt quá limit                             │
│    - Cho phép nếu trong limit                               │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Rate Limit Config:**
   ```
   MaxShortTermOrdersAndCancelsPerNBlocks: Limit per N blocks
   MaxStatefulOrdersPerNBlocks: Limit per N blocks
   MaxLeverageUpdatesPerNBlocks: Limit per N blocks
   ```

2. **Operation Counting:**
   ```
   Đếm operations trong sliding window của N blocks
   Bao gồm operations từ current block
   Đếm per subaccount
   ```

### Điểm quan trọng

1. **Per Subaccount Limits:**
   - Limits áp dụng per subaccount
   - Các subaccounts khác nhau có limits riêng
   - Owner address cũng có thể được xem xét

2. **Sliding Window:**
   - Đếm operations trong N blocks cuối
   - Window trượt khi blocks tiến
   - Operations expire sau N blocks

3. **Operation Types:**
   - Short-term orders và cancellations chia sẻ limit
   - Stateful orders có limit riêng
   - Leverage updates có limit riêng

4. **Batch Operations:**
   - Batch cancel tính là 1 operation
   - Hiệu quả hơn individual cancels
   - Vẫn subject to rate limits

5. **CheckTx Validation:**
   - Rate limits được kiểm tra tại CheckTx
   - Early rejection ngăn chặn wasted computation
   - Error code: `ErrBlockRateLimitExceeded`

6. **Block Advancement:**
   - Limits tồn tại qua blocks
   - Operations tính vào limit cho N blocks
   - Sau N blocks, operations không còn tính

### Lý do thiết kế

1. **Spam Prevention:** Rate limits ngăn chặn order book spam từ single subaccount.

2. **Fairness:** Đảm bảo tất cả users có fair access đến order book.

3. **System Stability:** Ngăn chặn system overload từ excessive operations.

4. **Flexibility:** Các limits khác nhau cho các operation types khác nhau.
