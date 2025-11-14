# Tài liệu Test: Long-Term Orders E2E Tests

## Tổng quan

File test này xác minh chức năng **Long-Term Order** trong CLOB module. Long-term orders là các orders tồn tại qua nhiều blocks cho đến khi được fill, hủy, hoặc expire. Test đảm bảo rằng:
1. Long-term orders có thể được đặt và tồn tại qua blocks
2. Long-term orders có thể được hủy
3. Order cancellation và placement trong cùng block được xử lý đúng
4. Fully filled orders không thể được hủy

---

## Test Function: TestPlaceOrder_StatefulCancelFollowedByPlaceInSameBlockErrorsInCheckTx

### Test Case: Thất bại - Hủy và Đặt Cùng Order trong Cùng Block

### Đầu vào
- **Block 2:**
  - Đặt long-term order: Alice mua 5 ở giá 10
- **Block 3:**
  - Hủy long-term order
  - Cố gắng đặt cùng order lại

### Đầu ra
- **Cancel CheckTx:** SUCCESS
- **Place CheckTx:** FAIL với lỗi "An uncommitted stateful order cancellation with this OrderId already exists"
- **Final State:** Order bị hủy, order mới không được đặt

### Tại sao chạy theo cách này?

1. **Uncommitted Cancellation:** Khi một order bị hủy trong cùng block, cancellation là uncommitted.
2. **Conflict Detection:** Hệ thống phát hiện conflict giữa cancellation và placement của cùng order.
3. **Early Rejection:** CheckTx từ chối placement để ngăn chặn invalid state.

---

## Test Function: TestCancelFullyFilledStatefulOrderInSameBlockItIsFilled

### Test Case: Thất bại - Hủy Fully Filled Order

### Đầu vào
- **Block 2:**
  - Đặt long-term order: Alice mua 5 ở giá 10
- **Block 3:**
  - Đặt matching order: Bob bán 5 ở giá 10 (fill đầy order của Alice)
  - Cố gắng hủy order của Alice

### Đầu ra
- **Match CheckTx:** SUCCESS
- **Cancel CheckTx:** SUCCESS
- **DeliverTx:** Cancel transaction THẤT BẠI với lỗi `ErrStatefulOrderCancellationFailedForAlreadyRemovedOrder`
- **Final State:** Order được fill đầy, cancellation thất bại

### Tại sao chạy theo cách này?

1. **Order Filled First:** Matching order fill long-term order trước khi cancellation thực thi.
2. **Cancellation Fails:** Không thể hủy một order đã được xóa (filled).
3. **Transaction Ordering:** DeliverTx xử lý transactions theo thứ tự, vì vậy fill xảy ra trước cancellation.

---

## Test Function: TestCancelStatefulOrder

### Test Case 1: Thành công - Hủy Order trong Cùng Block

### Đầu vào
- **Block 2:**
  - Đặt long-term order
  - Hủy cùng order

### Đầu ra
- **Both CheckTx:** SUCCESS
- **Final State:** Order không tồn tại trong state

### Tại sao chạy theo cách này?

1. **Same Block Cancellation:** Order có thể được hủy trong cùng block nó được đặt.
2. **State Cleanup:** Order được xóa khỏi state ngay lập tức.

---

### Test Case 2: Thành công - Hủy Order trong Block Tương lai

### Đầu vào
- **Block 2:**
  - Đặt long-term order
- **Block 3:**
  - Hủy order

### Đầu ra
- **Place CheckTx:** SUCCESS
- **Cancel CheckTx:** SUCCESS
- **Final State:** Order được xóa khỏi state

### Tại sao chạy theo cách này?

1. **Persistent Orders:** Long-term orders tồn tại qua blocks.
2. **Future Cancellation:** Orders có thể được hủy trong bất kỳ block tương lai nào trước khi expire.

---

## Tóm tắt Flow

### Long-Term Order Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│ 1. ĐẶT ORDER                                                │
│    - Order được đặt trên order book                         │
│    - Order tồn tại trong state                               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. ORDER TỒN TẠI                                            │
│    - Order vẫn trên book qua blocks                          │
│    - Có thể được match bất cứ lúc nào                        │
│    - Có thể được hủy bất cứ lúc nào                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. ORDER TERMINATION                                         │
│    - Option A: Fully filled → Removed                       │
│    - Option B: Cancelled → Removed                           │
│    - Option C: Expired → Removed                             │
└─────────────────────────────────────────────────────────────┘
```

### Cancel Order Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. SUBMIT CANCELLATION                                       │
│    - Tạo MsgCancelOrderStateful                              │
│    - Chỉ định order ID và good til block time               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CHECKTX VALIDATION                                        │
│    - Xác minh order tồn tại trong state                      │
│    - Kiểm tra cancellation parameters                        │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. DELIVERTX EXECUTION                                       │
│    - Xác minh order vẫn tồn tại (chưa fill)                  │
│    - Xóa order khỏi state                                    │
│    - Emit order removal events                               │
└─────────────────────────────────────────────────────────────┘
```

### Điểm quan trọng

1. **Long-Term Orders:**
   - Tồn tại qua nhiều blocks
   - GoodTilBlockTime chỉ định expiration time
   - Có thể được match, hủy, hoặc expire

2. **Cancellation:**
   - Có thể hủy trong cùng block hoặc blocks tương lai
   - Không thể hủy nếu order đã fill
   - Không thể hủy nếu order không tồn tại

3. **Conflict Detection:**
   - Hệ thống phát hiện conflicts giữa cancellation và placement
   - CheckTx từ chối conflicting operations sớm
   - Ngăn chặn invalid state

4. **State Management:**
   - Orders được track trong keeper state
   - Cancellation xóa order khỏi state
   - Fill xóa order khỏi state

5. **Transaction Ordering:**
   - DeliverTx xử lý transactions theo thứ tự
   - Transaction đầu tiên modify order thắng
   - Các transactions conflicting sau đó thất bại

### Lý do thiết kế

1. **Flexibility:** Long-term orders cho phép users đặt orders tồn tại cho đến khi fill hoặc hủy.

2. **Safety:** Conflict detection ngăn chặn invalid state từ cancellation/placement conflicts.

3. **Efficiency:** Early rejection tại CheckTx ngăn chặn wasted computation.

4. **Consistency:** Transaction ordering đảm bảo deterministic state updates.
