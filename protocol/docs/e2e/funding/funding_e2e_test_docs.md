# Tài liệu Test: Funding E2E Tests

## Tổng quan

File test này xác minh **Funding Mechanism** cho perpetual markets. Funding là khoản thanh toán định kỳ giữa long và short positions dựa trên premium (chênh lệch giữa mark price và index price). Test đảm bảo rằng:
1. Funding premiums được tính toán đúng dựa trên order book impact prices
2. Funding index được cập nhật đúng dựa trên premiums
3. Funding settlements được tính toán và áp dụng đúng cho subaccounts
4. Funding rate clamping hoạt động khi premiums vượt quá giới hạn

---

## Test Function: TestFunding

### Test Case 1: Index Price dưới Impact Bid, Positive Funding, Longs trả Shorts

### Đầu vào
- **Orders:**
  - Unmatched orders để tạo funding premiums:
    - Bob: Bán 2 BTC ở 28,005 (impact ask)
    - Alice: Mua 2 BTC ở 28,000 (impact bid)
  - Matched orders để thiết lập positions:
    - Bob: Bán 1 BTC ở 28,003 (matched)
    - Alice: Mua 0.8 BTC ở 28,003 (matched)
    - Carl: Mua 0.2 BTC ở 28,003 (matched)
- **Initial Index Price:** 28,002
- **Index Price for Premium:** 27,960 (dưới impact bid)
- **Oracle Price for Funding Index:** 27,000

### Đầu ra
- **Funding Premiums:** ~1,430 ppm (0.143%)
- **Funding Index:** 482
- **Settlements:**
  - Alice (long 0.8 BTC): Trả $3.856
  - Bob (short 1 BTC): Nhận $4.82
  - Carl (long 0.2 BTC): Trả $0.964

### Tại sao chạy theo cách này?

1. **Premium Calculation:** Khi index price (27,960) dưới impact bid (28,000), premium là dương.
   - Premium = (28,000 / 27,960) - 1 ≈ 0.00143 (0.143%)
2. **Longs Pay Shorts:** Premium dương có nghĩa longs trả shorts.
3. **Funding Index:** Được tính từ premium samples trong funding epoch.
4. **Settlement:** Được áp dụng khi subaccount nhận transfer, dựa trên chênh lệch funding index.

---

### Test Case 2: Index Price trên Impact Ask, Negative Funding, Final Funding Rate Clamped

### Đầu vào
- **Orders:** Giống Test Case 1
- **Initial Index Price:** 28,002
- **Index Price for Premium:** 34,000 (trên impact ask)
- **Oracle Price for Funding Index:** 33,500

### Đầu ra
- **Funding Premiums:** -176,323 ppm (-17.6%, nhưng bị clamped)
- **Funding Index:** -50,250 (clamped đến -12% dựa trên margin requirements)
- **Settlements:**
  - Alice (long 0.8 BTC): Nhận $402 (shorts trả longs)
  - Bob (short 1 BTC): Trả $502.5
  - Carl (long 0.2 BTC): Nhận $100.5

### Tại sao chạy theo cách này?

1. **Premium Calculation:** Khi index price (34,000) trên impact ask (28,005), premium là âm.
   - Premium = (28,005 / 34,000) - 1 ≈ -0.176 (-17.6%)
2. **Funding Rate Clamp:** Funding rate bị clamped để ngăn chặn thanh toán quá mức.
   - Clamp = premium_rate_clamp_factor × (initial_margin - maintenance_margin)
   - Clamp = 600% × (5% - 3%) = 12% = 120,000 ppm
3. **Shorts Pay Longs:** Premium âm có nghĩa shorts trả longs (sau khi clamped).
4. **Settlement:** Được áp dụng với funding rate bị clamped.

---

### Test Case 3: Index Price giữa Impact Bid và Ask, Zero Funding

### Đầu vào
- **Orders:** Giống Test Case 1
- **Initial Index Price:** 28,002
- **Index Price for Premium:** 28,003 (giữa impact bid 28,000 và ask 28,005)
- **Oracle Price for Funding Index:** 27,500

### Đầu ra
- **Funding Premiums:** Không có (zero)
- **Funding Index:** 0
- **Settlements:**
  - Alice: $0
  - Bob: $0
  - Carl: $0

### Tại sao chạy theo cách này?

1. **Zero Premium:** Khi index price ở giữa impact bid và ask, premium là zero.
2. **No Funding:** Không có funding payments khi premium là zero.
3. **Funding Index:** Không thay đổi (bắt đầu ở 0, vẫn ở 0).

---

## Tóm tắt Flow

### Funding Calculation Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. ĐẶT ORDERS                                               │
│    - Đặt unmatched orders để set impact prices              │
│    - Đặt matched orders để mở positions                     │
│    - Cập nhật initial index price                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. FUNDING SAMPLE EPOCHS                                    │
│    - Tiến đến funding tick epoch                            │
│    - Cập nhật index price cho premium calculation           │
│    - Thu thập premium samples (60 samples mỗi funding tick) │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. TÍNH TOÁN FUNDING PREMIUMS                               │
│    - Tính premium = (impact_price / index_price) - 1        │
│    - Nếu index < impact_bid: positive premium (longs trả)    │
│    - Nếu index > impact_ask: negative premium (shorts trả)  │
│    - Nếu index giữa bid/ask: zero premium                   │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. CẬP NHẬT FUNDING INDEX                                   │
│    - Tính funding index từ premium samples                   │
│    - Áp dụng funding rate clamp nếu cần                      │
│    - Cập nhật perpetual funding index                        │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. SETTLE FUNDING                                           │
│    - Khi subaccount nhận transfer                            │
│    - Tính settlement = (funding_index - position_index) × size │
│    - Áp dụng settlement vào subaccount balance              │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Premium Calculation:**
   ```
   Index < Impact Bid → Positive Premium → Longs Pay Shorts
   Index > Impact Ask → Negative Premium → Shorts Pay Longs
   Index Giữa → Zero Premium → No Funding
   ```

2. **Funding Index:**
   ```
   Initial: 0
   Sau Epoch: Cập nhật dựa trên premium samples
   Clamped: Nếu premium vượt quá giới hạn dựa trên margin
   ```

3. **Settlement:**
   ```
   Position Mở: Funding Index = 0
   Sau Funding: Funding Index Cập nhật
   Khi Transfer: Settlement = (New Index - Old Index) × Size
   ```

### Điểm quan trọng

1. **Impact Prices:**
   - Impact bid: Giá bid tốt nhất từ order book
   - Impact ask: Giá ask tốt nhất từ order book
   - Được sử dụng để tính premium khi index price nằm ngoài phạm vi

2. **Premium Calculation:**
   - Premium = (Impact Price / Index Price) - 1
   - Positive: Longs trả shorts
   - Negative: Shorts trả longs
   - Zero: Không có funding

3. **Funding Rate Clamp:**
   - Ngăn chặn funding payments quá mức
   - Dựa trên margin requirements
   - Công thức: clamp_factor × (initial_margin - maintenance_margin)

4. **Funding Index:**
   - Theo dõi cumulative funding theo thời gian
   - Được cập nhật tại mỗi funding tick
   - Được sử dụng để tính settlement khi position được đóng hoặc transfer

5. **Settlement:**
   - Được tính khi subaccount nhận transfer
   - Settlement = (Current Funding Index - Position Funding Index) × Position Size
   - Được áp dụng vào subaccount balance

6. **Premium Samples:**
   - 60 samples được thu thập mỗi funding tick epoch
   - Samples được sử dụng để tính average premium
   - Premium samples reset khi bắt đầu funding tick epoch mới

### Lý do thiết kế

1. **Công bằng:** Funding đảm bảo long và short positions được cân bằng bằng cách chuyển giá trị dựa trên premium.

2. **Price Discovery:** Premium phản ánh sự khác biệt giữa mark price (order book) và index price (oracle).

3. **An toàn:** Funding rate clamping ngăn chặn thanh toán quá mức có thể gây ra liquidations.

4. **Hiệu quả:** Funding index cho phép tính toán settlement hiệu quả mà không cần tính lại tất cả historical premiums.

5. **Minh bạch:** Premium samples được thu thập theo thời gian để đảm bảo tính toán funding công bằng và chính xác.
