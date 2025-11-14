# Tài liệu Test: Batch Cancel E2E Tests

## Tổng quan

File test này xác minh chức năng **Batch Cancel Order** trong CLOB module. Batch cancel cho phép users hủy nhiều orders trong một transaction, được nhóm theo CLOB pair. Test đảm bảo rằng:
1. Nhiều orders có thể được hủy trong một batch cancel transaction
2. Batch cancel hoạt động cho unfilled, partially filled, và fully filled orders
3. Batch cancel tôn trọng GoodTilBlock constraints
4. Batch cancel hoạt động qua nhiều blocks

---

## Test Function: TestBatchCancelSingleCancelFunctionality

### Test Case 1: Thành công - Hủy Unfilled Short-Term Order

### Đầu vào
- **Block 1:**
  - Đặt order: Alice mua 5 ở giá 10, GTB 5
- **Block 1:**
  - Batch cancel: Hủy order của Alice với client ID 0 trên CLOB pair 0, GTB 5

### Đầu ra
- **Order:** Được xóa khỏi memclob
- **Cancel Expiration:** Được set ở block 5
- **Fill Amount:** 0 (unfilled)

### Tại sao chạy theo cách này?

1. **Batch Cancel:** Single order được hủy qua batch cancel message.
2. **Unfilled Order:** Order chưa bao giờ match, vì vậy fill amount là 0.
3. **Removal:** Order được xóa khỏi memclob ngay lập tức.

---

### Test Case 2: Thành công - Batch Cancel Partially Filled Order (Cùng Block)

### Đầu vào
- **Block 1:**
  - Đặt order: Alice mua 5 ở giá 10, GTB 5
  - Đặt order: Bob bán 4 ở giá 10, GTB 20 (match với Alice)
  - Batch cancel: Hủy order của Alice

### Đầu ra
- **Order:** Được xóa khỏi memclob
- **Fill Amount:** 4 (40% filled)
- **Cancel Expiration:** Được set ở block 5

### Tại sao chạy theo cách này?

1. **Partial Fill:** Order được partially fill (4 trong 5) trước khi cancellation.
2. **Same Block:** Cancellation xảy ra trong cùng block với partial fill.
3. **Remaining Cancelled:** Phần unfilled còn lại (1) được hủy.

---

### Test Case 3: Thành công - Hủy Partially Filled Order (Block Tiếp theo)

### Đầu vào
- **Block 1:**
  - Đặt order: Alice mua 5 ở giá 10, GTB 5
  - Đặt order: Bob bán 4 ở giá 10, GTB 20 (match)
- **Block 2:**
  - Batch cancel: Hủy order của Alice

### Đầu ra
- **Order:** Được xóa khỏi memclob
- **Fill Amount:** 4 (40% filled)
- **Cancel Expiration:** Được set ở block 5

### Tại sao chạy theo cách này?

1. **Cross-Block:** Cancellation xảy ra trong block sau partial fill.
2. **Same Behavior:** Hoạt động giống như same-block cancellation.
3. **Fill Preserved:** Fill amount từ block trước được giữ lại.

---

### Test Case 4: Thành công - Hủy Fully Filled Order

### Đầu vào
- **Block 1:**
  - Đặt order: Alice mua 5 ở giá 10, GTB 5
  - Đặt order: Bob bán 5 ở giá 10, GTB 20 (fully match)
- **Block 2:**
  - Batch cancel: Hủy order của Alice

### Đầu ra
- **Order:** Được xóa khỏi memclob (đã được xóa sau fill)
- **Fill Amount:** 5 (100% filled)
- **Cancel Expiration:** Được set ở block 5

### Tại sao chạy theo cách này?

1. **Fully Filled:** Order được fill đầy trong block 1.
2. **Already Removed:** Order đã được xóa khỏi memclob sau fill.
3. **Cancel Succeeds:** Cancellation thành công mặc dù order đã fill.

---

### Test Case 5: Thất bại - Hủy với GTB < Order GTB Không Xóa Order

### Đầu vào
- **Block 1:**
  - Đặt order: Alice mua 5 ở giá 10, GTB 20
- **Block 2:**
  - Batch cancel: Hủy với GTB 5 (ít hơn order's GTB 20)

### Đầu ra
- **Order:** Vẫn trong memclob (không được xóa)
- **Cancel Expiration:** Được set ở block 5
- **Order Expiration:** Vẫn block 20

### Tại sao chạy theo cách này?

1. **GTB Constraint:** Cancel GTB (5) < Order GTB (20).
2. **Order Not Removed:** Order vẫn tồn tại vì cancel expire trước order.
3. **Cancel Expires First:** Cancel expire ở block 5, order ở block 20.

---

## Tóm tắt Flow

### Batch Cancel Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. TẠO BATCH CANCEL                                         │
│    - Chỉ định subaccount ID                                  │
│    - Nhóm orders theo CLOB pair                              │
│    - Liệt kê client IDs để hủy                               │
│    - Set GoodTilBlock                                        │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CHECKTX VALIDATION                                       │
│    - Validate message format                                 │
│    - Kiểm tra subaccount tồn tại                             │
│    - Xác minh CLOB pairs tồn tại                              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. DELIVERTX EXECUTION                                      │
│    - Cho mỗi CLOB pair:                                       │
│      * Cho mỗi client ID:                                    │
│        - Tìm order theo subaccount + client ID + CLOB pair  │
│        - Hủy order (nếu tồn tại và chưa expire)             │
│    - Set cancel expiration                                   │
│    - Xóa orders khỏi memclob                                  │
└─────────────────────────────────────────────────────────────┘
```

### Điểm quan trọng

1. **Batch Structure:**
   - Orders được nhóm theo CLOB pair
   - Nhiều client IDs cho mỗi CLOB pair
   - Single transaction hủy nhiều orders

2. **Order Matching:**
   - Orders được match bởi: SubaccountId + ClientId + ClobPairId
   - Phải match chính xác để hủy
   - Non-existent orders bị bỏ qua

3. **GTB Constraints:**
   - Cancel GTB phải >= order GTB để xóa order
   - Nếu cancel GTB < order GTB, order vẫn tồn tại
   - Cancel expire ở GTB của nó

4. **Fill Handling:**
   - Unfilled orders: Hủy hoàn toàn
   - Partially filled: Phần còn lại được hủy
   - Fully filled: Cancel thành công nhưng order đã được xóa

5. **Cross-Block:**
   - Có thể hủy orders từ blocks trước
   - Fill amounts được giữ lại
   - Hoạt động giống như same-block cancellation

### Lý do thiết kế

1. **Efficiency:** Batch cancel cho phép hủy nhiều orders trong một transaction.

2. **Flexibility:** Nhóm theo CLOB pair cho phép selective cancellation.

3. **Safety:** GTB constraints ngăn chặn premature order removal.

4. **Consistency:** Hoạt động nhất quán qua blocks.
