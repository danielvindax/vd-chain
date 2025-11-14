# Tài liệu Test: Equity Tier Limit E2E Tests

## Tổng quan

File test này xác minh chức năng **Equity Tier Limit** trong CLOB module. Equity tier limits giới hạn số lượng stateful orders (long-term và conditional orders) một subaccount có thể có mở dựa trên Total Net Collateral (TNC) của họ. Test đảm bảo rằng:
1. Subaccounts với TNC thấp hơn có ít allowed stateful orders hơn
2. Subaccounts với TNC cao hơn có nhiều allowed stateful orders hơn
3. Order cancellation có thể giải phóng slots cho orders mới
4. Cả long-term và conditional orders đều tính vào limits

---

## Test Function: TestPlaceOrder_EquityTierLimit

### Test Case 1: Thất bại - Long-Term Order Vượt Quá Max Open Stateful Orders

### Đầu vào
- **Subaccount:** Alice với TNC < $5,000
- **Existing Orders:**
  - 1 conditional order (StopLoss)
- **New Order:** Long-term order
- **Equity Tier Config:**
  - Tier 0: $0 TNC → 0 orders
  - Tier 1: $5,000 TNC → 1 order
  - Tier 2: $70,000 TNC → 100 orders

### Đầu ra
- **CheckTx:** FAIL (sau khi tiến block)
- **Error:** Sẽ vượt quá max open stateful orders

### Tại sao chạy theo cách này?

1. **Tier Limit:** Alice ở tier 1 (TNC < $5,000), limit = 1 order.
2. **Already Has 1:** Đã có 1 conditional order.
3. **Exceeds Limit:** Long-term order mới sẽ vượt quá limit của 1.
4. **Rejection:** Order bị từ chối để ngăn chặn vượt quá limit.

---

### Test Case 2: Thất bại - Conditional Order Vượt Quá Max Open Stateful Orders

### Đầu vào
- **Subaccount:** Alice với TNC < $5,000
- **Existing Orders:**
  - 1 long-term order
- **New Order:** Conditional order (StopLoss)
- **Equity Tier Config:** Giống Test Case 1

### Đầu ra
- **CheckTx:** FAIL (sau khi tiến block)
- **Error:** Sẽ vượt quá max open stateful orders

### Tại sao chạy theo cách này?

1. **Same Limit:** Conditional orders tính vào cùng limit như long-term orders.
2. **Already Has 1:** Đã có 1 long-term order.
3. **Exceeds Limit:** Conditional order mới sẽ vượt quá limit của 1.

---

### Test Case 3: Thành công - Order Cancellation Giải phóng Slot

### Đầu vào
- **Subaccount:** Alice với TNC < $5,000
- **Existing Orders:**
  - 1 conditional order (StopLoss)
- **Cancellation:** Hủy conditional order
- **New Order:** Long-term order (cùng block)
- **Equity Tier Config:** Giống Test Case 1

### Đầu ra
- **Cancellation:** SUCCESS
- **New Order:** SUCCESS
- **Final State:** 1 long-term order (conditional đã hủy)

### Tại sao chạy theo cách này?

1. **Cancellation First:** Conditional order được hủy, giải phóng slot.
2. **Slot Available:** Sau cancellation, slot có sẵn cho order mới.
3. **Same Block:** Cancellation và placement trong cùng block hoạt động.
4. **Limit Respected:** Final state có 1 order, trong limit.

---

### Test Case 4: Thất bại - Conditional Order Sẽ Vượt Quá Limit (Untriggered)

### Đầu vào
- **Subaccount:** Alice với TNC < $5,000
- **Existing Orders:**
  - 1 long-term order
- **New Order:** Conditional order (TakeProfit, untriggered)
- **Equity Tier Config:** Giống Test Case 1

### Đầu ra
- **CheckTx:** FAIL (sau khi tiến block)
- **Error:** Sẽ vượt quá max open stateful orders

### Tại sao chạy theo cách này?

1. **Untriggered Counts:** Untriggered conditional orders tính vào limit.
2. **Same Limit:** Cả triggered và untriggered conditional orders đều tính.
3. **Exceeds Limit:** Conditional order mới sẽ vượt quá limit.

---

## Tóm tắt Flow

### Equity Tier Limit Check

```
┌─────────────────────────────────────────────────────────────┐
│ 1. TÍNH TOÁN SUBACCOUNT TNC                                 │
│    - Lấy Total Net Collateral của subaccount                │
│    - Bao gồm tất cả positions và assets                      │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. XÁC ĐỊNH EQUITY TIER                                     │
│    - Tìm tier dựa trên TNC amount                           │
│    - Tier 0: $0 TNC → 0 orders                              │
│    - Tier 1: $5,000 TNC → 1 order                           │
│    - Tier 2: $70,000 TNC → 100 orders                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. ĐẾM EXISTING STATEFUL ORDERS                             │
│    - Đếm long-term orders                                   │
│    - Đếm conditional orders (triggered và untriggered)      │
│    - Đếm orders trong cùng block (uncommitted)              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. VALIDATE LIMIT                                            │
│    - Kiểm tra nếu order mới sẽ vượt quá limit                │
│    - Xem xét cancellations trong cùng block                  │
│    - Từ chối nếu sẽ vượt quá limit                           │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Equity Tiers:**
   ```
   Tier 0: $0 TNC → 0 orders
   Tier 1: $5,000 TNC → 1 order
   Tier 2: $70,000 TNC → 100 orders
   ```

2. **Order Counting:**
   ```
   Long-term orders: Tính vào limit
   Conditional orders (triggered): Tính vào limit
   Conditional orders (untriggered): Tính vào limit
   ```

### Điểm quan trọng

1. **TNC-Based Limits:**
   - Limits dựa trên Total Net Collateral
   - TNC cao hơn = nhiều allowed orders hơn
   - Bảo vệ hệ thống khỏi order book spam

2. **Stateful Orders:**
   - Long-term orders tính vào limit
   - Conditional orders tính vào limit
   - Short-term orders không tính (expire cùng block)

3. **Same Block Logic:**
   - Cancellations giải phóng slots ngay lập tức
   - Có thể hủy và đặt trong cùng block
   - Uncommitted orders tính vào limit

4. **Untriggered Conditionals:**
   - Untriggered conditional orders tính vào limit
   - Phải có slot available khi đặt
   - Triggering không thay đổi count (cùng order)

5. **Validation Timing:**
   - Được kiểm tra khi đặt order
   - Sau block advancement (cho committed orders)
   - Xem xét same-block cancellations

### Lý do thiết kế

1. **Resource Management:** Limits ngăn chặn order book spam từ low-collateral accounts.

2. **Fairness:** Higher collateral accounts có nhiều order slots hơn.

3. **Flexibility:** Cancellations cho phép users quản lý order slots của họ.

4. **Safety:** Ngăn chặn system overload từ excessive stateful orders.
