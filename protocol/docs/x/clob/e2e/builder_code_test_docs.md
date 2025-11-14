# Tài liệu Test: Builder Code Orders E2E Tests

## Tổng quan

File test này xác minh chức năng **Builder Code** trong CLOB module. Builder codes cho phép order builders nhận fees khi orders của họ được match. Test đảm bảo rằng:
1. Orders với builder codes có thể được đặt thành công
2. Builder fees được tính toán và trả đúng khi orders match
3. Orders không có builder codes hoạt động bình thường
4. Fee calculation dựa trên fill amount và fee percentage

---

## Test Function: TestBuilderCodeOrders

### Test Case: Order với Builder Code Fill và Fees được Trả

### Đầu vào
- **Orders:**
  - Alice: Buy order với builder code
    - SubaccountId: Alice_Num0
    - ClobPairId: 0
    - Side: BUY
    - Quantums: 10,000,000,000 (1 BTC)
    - Subticks: 500,000,000 (50,000 USDC/BTC)
    - GoodTilBlock: 20
    - BuilderCodeParameters:
      - BuilderAddress: Địa chỉ tài khoản Carl_Num0
      - FeePpm: 1000 (0.1%)
  - Bob: Sell order không có builder code
    - SubaccountId: Bob_Num0
    - ClobPairId: 0
    - Side: SELL
    - Quantums: 10,000,000,000 (1 BTC)
    - Subticks: 500,000,000 (50,000 USDC/BTC)
    - GoodTilBlock: 20
- **Match Operation:** Orders match ở block 2

### Đầu ra
- **CheckTx:** Cả hai orders pass CheckTx validation
- **Order Fill:** Cả hai orders được fill đầy (10,000,000,000 quantums)
- **Builder Fee:** Carl nhận 50,000,000 quantums (0.1% của 50,000 USDC)
- **Builder Balance:** Số dư của Carl tăng bằng builder fee amount

### Tại sao chạy theo cách này?

1. **Builder Code Mechanism:** Test cơ chế fee-sharing nơi order builders nhận fees.
2. **Fee Calculation:** Fee = (Fill Amount × Price × FeePpm) / 1,000,000
   - Fill Amount: 10,000,000,000 quantums (1 BTC)
   - Price: 500,000,000 subticks (50,000 USDC/BTC)
   - FeePpm: 1000 (0.1%)
   - Fee = (10,000,000,000 × 500,000,000 × 1000) / (1,000,000 × 10^8) = 50,000,000 quantums
3. **Fee Payment:** Builder fee được trả từ matched order proceeds đến builder address.
4. **Order Matching:** Cả hai orders match hoàn toàn, trigger fee payment.
5. **Balance Verification:** Builder's balance được kiểm tra trước và sau match để xác minh fee payment.

---

## Tóm tắt Flow

### Builder Code Order Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. ĐẶT ORDERS                                              │
│    - Alice đặt buy order với builder code                  │
│    - Bob đặt sell order không có builder code              │
│    - Cả hai orders pass CheckTx                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. TIẾN BLOCK                                              │
│    - Orders được match trong block 2                         │
│    - Match operation được xử lý                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. TÍNH TOÁN BUILDER FEE                                   │
│    - Fee = (Fill Amount × Price × FeePpm) / 1,000,000     │
│    - Fee được trừ từ order proceeds                        │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. TRẢ BUILDER FEE                                         │
│    - Fee được chuyển đến builder address                   │
│    - Builder balance tăng                                   │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. XÁC MINH KẾT QUẢ                                         │
│    - Orders được fill đầy                                   │
│    - Builder fee được trả đúng                             │
│    - Builder balance khớp với expected                     │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Order States:**
   ```
   Đặt → CheckTx Passed → Matched → Filled → Fee Paid
   ```

2. **Builder Fee Calculation:**
   ```
   Fill Amount × Price × FeePpm / 1,000,000 = Fee
   ```

3. **Balance Updates:**
   ```
   Pre-Match Balance → Match → Fee Payment → Post-Match Balance
   ```

### Điểm quan trọng

1. **Builder Code Parameters:**
   - BuilderAddress: Địa chỉ nhận fee
   - FeePpm: Fee percentage trong parts per million (1000 = 0.1%)

2. **Fee Payment:**
   - Fee được trả từ matched order proceeds
   - Chỉ orders với builder codes tạo fees
   - Fee được tính dựa trên fill amount và price

3. **Order Matching:**
   - Orders phải match để fees được trả
   - Partial fills dẫn đến proportional fees
   - Full fills dẫn đến full fee calculation

4. **Balance Verification:**
   - Builder balance được kiểm tra trước match
   - Balance delta được tính sau match
   - Delta nên bằng builder fee

### Lý do thiết kế

1. **Fee Sharing:** Builder codes khuyến khích order builders cung cấp liquidity.

2. **Fee Calculation:** Fee tỷ lệ với trade size và price, đảm bảo compensation công bằng.

3. **Flexibility:** Orders có thể có builder codes hoặc không, cho phép các cấu trúc fee khác nhau.

4. **Verification:** Balance checks đảm bảo fees được tính toán và trả đúng.
