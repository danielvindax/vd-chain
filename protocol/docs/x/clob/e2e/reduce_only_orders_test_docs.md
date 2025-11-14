# Tài liệu Test: Reduce-Only Orders E2E Tests

## Tổng quan

File test này xác minh chức năng **Reduce-Only Order** trong CLOB module. Reduce-only orders là các orders chỉ có thể giảm existing position, không mở position mới hoặc tăng existing position. Test đảm bảo rằng:
1. Reduce-only orders có thể partially match để giảm position
2. Reduce-only orders không thể tăng position size
3. Reduce-only orders hoạt động với IOC orders
4. Reduce-only orders hoạt động qua nhiều blocks

---

## Test Function: TestReduceOnlyOrders

### Test Case 1: Thành công - IOC Reduce-Only Order Partially Match, Maker Fully Filled (Cùng Block)

### Đầu vào
- **Subaccounts:**
  - Carl: 100,000 USD
  - Alice: 1 BTC Long, 500,000 USD
- **Orders (Block 1):**
  - Carl: Mua 10 ở giá 500,000 (maker)
  - Alice: Bán 15 ở giá 500,000, IOC, Reduce-Only (taker)
- **Match:** Order của Alice match với order của Carl

### Đầu ra
- **Carl Order:** Fully filled (10)
- **Alice Order:** Partially filled (10), remaining cancelled
- **Carl Position:** 10 (new position opened)
- **Alice Position:** 0.999999 (giảm từ 1 BTC)

### Tại sao chạy theo cách này?

1. **Reduce-Only:** Order của Alice chỉ có thể giảm 1 BTC long position của cô ấy.
2. **Partial Match:** Order match 10 units, giảm position từ 1 BTC xuống 0.999999 BTC.
3. **IOC Behavior:** Remaining 5 units bị hủy (IOC rule).
4. **Maker Filled:** Maker order của Carl được fill đầy.

---

### Test Case 2: Thành công - IOC Reduce-Only Order Partially Match (Block Thứ hai)

### Đầu vào
- **Subaccounts:**
  - Carl: 100,000 USD
  - Alice: 1 BTC Long, 500,000 USD
- **Orders:**
  - Block 1: Carl mua 10 ở giá 500,000
  - Block 2: Alice bán 15 ở giá 500,000, IOC, Reduce-Only

### Đầu ra
- **Carl Order:** Fully filled (10)
- **Alice Order:** Partially filled (10), remaining cancelled
- **Alice Position:** Giảm từ 1 BTC xuống 0.999999 BTC

### Tại sao chạy theo cách này?

1. **Cross-Block Matching:** Reduce-only order có thể match với orders từ blocks trước.
2. **Same Behavior:** Reduce-only logic hoạt động giống qua blocks.
3. **Position Reduction:** Position giảm bằng matched amount.

---

### Test Case 3: Thành công - IOC Reduce-Only Order Partially Match, Maker Partially Filled

### Đầu vào
- **Subaccounts:**
  - Carl: 100,000 USD
  - Alice: 1 BTC Long, 500,000 USD
- **Orders:**
  - Block 1: Carl mua 80 ở giá 500,000
  - Block 2: Alice bán 15 ở giá 500,000, IOC, Reduce-Only

### Đầu ra
- **Carl Order:** Partially filled (15), 65 còn lại trên book
- **Alice Order:** Partially filled (15), remaining cancelled
- **Alice Position:** Giảm từ 1 BTC xuống 0.9999985 BTC

### Tại sao chạy theo cách này?

1. **Partial Fill Both:** Cả hai orders được partially fill.
2. **Maker Remains:** Order của Carl vẫn trên book với remaining size.
3. **Position Reduced:** Position của Alice giảm bằng filled amount.

---

## Tóm tắt Flow

### Reduce-Only Order Logic

```
┌─────────────────────────────────────────────────────────────┐
│ 1. KIỂM TRA EXISTING POSITION                               │
│    - Query current position của subaccount                   │
│    - Xác định position side (long/short)                    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. VALIDATE ORDER DIRECTION                                  │
│    - Reduce-only buy: Chỉ hợp lệ nếu short position tồn tại │
│    - Reduce-only sell: Chỉ hợp lệ nếu long position tồn tại │
│    - Từ chối nếu không có position hoặc wrong direction      │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. TÍNH TOÁN MAX FILL AMOUNT                                │
│    - Max fill = min(order size, position size)              │
│    - Không thể fill nhiều hơn existing position              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. THỰC THI MATCH                                           │
│    - Fill lên đến max fill amount                           │
│    - Giảm position bằng fill amount                         │
│    - Hủy remaining nếu IOC                                  │
└─────────────────────────────────────────────────────────────┘
```

### Điểm quan trọng

1. **Reduce-Only Constraint:**
   - Chỉ có thể giảm existing position
   - Không thể mở position mới
   - Không thể tăng position size

2. **Position Direction:**
   - Reduce-only buy: Chỉ hoạt động với short position
   - Reduce-only sell: Chỉ hoạt động với long position
   - Phải khớp với position direction

3. **Fill Amount:**
   - Giới hạn bởi existing position size
   - Không thể fill nhiều hơn position
   - Partial fills được phép

4. **IOC Compatibility:**
   - Reduce-only hoạt động với IOC orders
   - Remaining size bị hủy nếu không fill đầy
   - Immediate execution hoặc cancellation

5. **Cross-Block Matching:**
   - Có thể match với orders từ blocks trước
   - Position được kiểm tra tại match time
   - Hoạt động giống như same-block matching

### Lý do thiết kế

1. **Risk Management:** Reduce-only orders giúp users đóng positions mà không vô tình mở positions mới.

2. **Position Control:** Đảm bảo users chỉ có thể giảm risk, không tăng nó.

3. **Flexibility:** Hoạt động với các loại orders khác nhau (IOC, regular, etc.).

4. **Safety:** Ngăn chặn accidental position increases.
