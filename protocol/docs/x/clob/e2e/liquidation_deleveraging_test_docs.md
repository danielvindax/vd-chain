# Tài liệu Test: Liquidation và Deleveraging E2E Tests

## Tổng quan

File test này xác minh chức năng **Liquidation và Deleveraging** trong CLOB module. Khi một subaccount trở nên undercollateralized (dưới maintenance margin), nó có thể được liquidate. Test đảm bảo rằng:
1. Liquidations tôn trọng position block limits (MinPositionNotionalLiquidated, MaxPositionPortionLiquidatedPpm)
2. Liquidations tôn trọng subaccount block limits (MaxNotionalLiquidated)
3. Liquidations hoạt động cho cả long và short positions
4. Insurance fund cover losses khi cần

---

## Test Function: TestLiquidationConfig

### Test Case 1: Liquidating Short - Tôn trọng MinPositionNotionalLiquidated

### Đầu vào
- **Subaccounts:**
  - Carl: 1 BTC Short, 50,499 USD collateral (undercollateralized)
  - Dave: 1 BTC Long, 50,000 USD collateral
- **Order:** Dave bán 1 BTC ở 50,000
- **Liquidation Config:**
  - MinPositionNotionalLiquidated: $100,000
  - MaxPositionPortionLiquidatedPpm: 1% (10,000 ppm)
  - Oracle Price: 50,000

### Đầu ra
- **Liquidation:** Toàn bộ position được liquidate (1 BTC)
- **Carl Balance:** 50,499 - 50,000 - 250 (fees) = 249 USD
- **Dave Balance:** 50,000 + 50,000 = 100,000 USD

### Tại sao chạy theo cách này?

1. **Minimum Notional:** 1% của $50,000 = $500, nhưng minimum là $100,000.
2. **Entire Position:** Vì $500 < $100,000, toàn bộ position được liquidate.
3. **Full Liquidation:** Tất cả 1 BTC được liquidate để đáp ứng minimum requirement.

---

### Test Case 2: Liquidating Long - Tôn trọng MaxPositionPortionLiquidatedPpm

### Đầu vào
- **Subaccounts:**
  - Carl: 1 BTC Short, 100,000 USD
  - Dave: 1 BTC Long, 49,501 USD (undercollateralized)
- **Order:** Carl mua 1 BTC ở 50,000
- **Liquidation Config:**
  - MinPositionNotionalLiquidated: $1,000
  - MaxPositionPortionLiquidatedPpm: 10% (100,000 ppm)
  - Oracle Price: 50,000

### Đầu ra
- **Liquidation:** 10% của position được liquidate (0.1 BTC)
- **Dave Balance:** -49,501 + 5,000 - 25 (fees) = -44,526 USD
- **Dave Position:** 0.9 BTC long còn lại

### Tại sao chạy theo cách này?

1. **Portion Limit:** 10% của $50,000 = $5,000 worth của BTC.
2. **Partial Liquidation:** Chỉ 0.1 BTC (10%) được liquidate.
3. **Remaining Position:** 0.9 BTC position vẫn còn.

---

### Test Case 3: Liquidating Short - Tôn trọng MaxNotionalLiquidated

### Đầu vào
- **Subaccounts:**
  - Carl: 1 BTC Short, 50,499 USD (undercollateralized)
  - Dave: 1 BTC Long, 50,000 USD
- **Order:** Dave bán 1 BTC ở 49,500
- **Liquidation Config:**
  - MaxNotionalLiquidated: $5,000 mỗi block
  - Oracle Price: 50,000

### Đầu ra
- **Liquidation:** Chỉ $5,000 worth được liquidate (0.1 BTC)
- **Carl Balance:** 50,499 - 5,000 - 25 (fees) = 45,474 USD
- **Carl Position:** 0.9 BTC short còn lại

### Tại sao chạy theo cách này?

1. **Subaccount Limit:** Tối đa $5,000 có thể được liquidate mỗi block.
2. **Partial Liquidation:** Chỉ 0.1 BTC ($5,000 worth) được liquidate.
3. **Remaining Position:** 0.9 BTC position vẫn còn cho future liquidation.

---

## Tóm tắt Flow

### Liquidation Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. PHÁT HIỆN UNDERCOLLATERALIZATION                        │
│    - Subaccount TNC < maintenance margin                    │
│    - Liquidations daemon xác định liquidatable accounts     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. TÍNH TOÁN LIQUIDATION AMOUNT                             │
│    - Áp dụng position block limits                           │
│      * MinPositionNotionalLiquidated                         │
│      * MaxPositionPortionLiquidatedPpm                       │
│    - Áp dụng subaccount block limits                         │
│      * MaxNotionalLiquidated                                 │
│    - Lấy minimum của tất cả limits                           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. TÌM MATCHING ORDERS                                      │
│    - Tìm kiếm order book cho matching orders                 │
│    - Sử dụng fillable price config                           │
│    - Match ở best available price                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. THỰC THI LIQUIDATION                                     │
│    - Đóng phần position                                      │
│    - Chuyển funds đến counterparty                           │
│    - Charge liquidation fee                                  │
│    - Cập nhật subaccount state                               │
└─────────────────────────────────────────────────────────────┘
```

### Liquidation Limits

1. **Position Block Limits:**
   - MinPositionNotionalLiquidated: Minimum $ amount để liquidate
   - MaxPositionPortionLiquidatedPpm: Maximum % của position để liquidate
   - Áp dụng per position per block

2. **Subaccount Block Limits:**
   - MaxNotionalLiquidated: Maximum $ amount per subaccount per block
   - MaxQuantumsInsuranceLost: Maximum insurance fund loss per block
   - Áp dụng per subaccount per block

### Điểm quan trọng

1. **Liquidation Triggers:**
   - Subaccount TNC < maintenance margin
   - Được phát hiện bởi liquidations daemon
   - Liquidatable accounts được xác định mỗi block

2. **Liquidation Amount:**
   - Được tính dựa trên nhiều limits
   - Minimum của position limits và subaccount limits
   - Đảm bảo controlled liquidation rate

3. **Price Discovery:**
   - Sử dụng fillable price config
   - Match ở best available price trên order book
   - Có thể sử dụng insurance fund nếu không có matching orders

4. **Liquidation Fees:**
   - Được charge cho liquidated account
   - MaxLiquidationFeePpm set maximum fee
   - Fees compensate liquidators

5. **Partial Liquidation:**
   - Có thể liquidate phần position
   - Remaining position vẫn mở
   - Có thể được liquidate lại trong future blocks

6. **Insurance Fund:**
   - Cover losses khi liquidation price không thuận lợi
   - MaxQuantumsInsuranceLost giới hạn fund exposure
   - Bảo vệ protocol khỏi excessive losses

### Lý do thiết kế

1. **Risk Management:** Liquidation limits ngăn chặn excessive liquidation trong single block.

2. **Market Stability:** Controlled liquidation rate ngăn chặn market disruption.

3. **Fairness:** Limits đảm bảo tất cả liquidatable accounts được đối xử công bằng.

4. **Safety:** Insurance fund bảo vệ protocol khỏi extreme market conditions.
