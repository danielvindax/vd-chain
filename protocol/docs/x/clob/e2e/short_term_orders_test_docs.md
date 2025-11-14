# Tài liệu Test: Short-Term Orders E2E Tests

## Tổng quan

File test này xác minh chức năng **Short-Term Order** trong CLOB (Central Limit Order Book) module. Short-term orders là các orders expire ở cuối block hiện tại. Test đảm bảo rằng:
1. Orders có thể được đặt trên order book
2. Orders có thể được match với existing orders
3. Off-chain và on-chain messages được emit đúng
4. Order state được track đúng

---

## Test Function: TestPlaceOrder

### Test Case 1: Thành công - Đặt Order trên Order Book

### Đầu vào
- **Order:**
  - SubaccountId: Alice_Num0
  - ClobPairId: 0
  - Side: BUY
  - Quantums: 5
  - Subticks: 1,000,000 (price = 10)
  - GoodTilBlock: 20

### Đầu ra
- **CheckTx:** SUCCESS
- **Off-chain Messages:**
  - Order place message
  - Order update message (với fill amount = 0)
- **On-chain Messages:**
  - Indexer block event trong block tiếp theo

### Tại sao chạy theo cách này?

1. **Order Placement:** Order được đặt thành công trên order book.
2. **Off-chain Updates:** Indexer cần được thông báo về orders mới cho off-chain systems.
3. **On-chain Events:** Block events được emit trong block tiếp theo cho on-chain tracking.

---

### Test Case 2: Thành công - Match Order Đầy đủ

### Đầu vào
- **Orders:**
  - Alice: Mua 5 ở giá 10
  - Bob: Bán 5 ở giá 10

### Đầu ra
- **CheckTx:** Cả hai orders SUCCESS
- **Orders Filled:** Cả hai orders được fill đầy
- **Off-chain Messages:**
  - Order place messages cho cả hai orders
  - Order update messages hiển thị full fill
- **On-chain Messages:**
  - Subaccount update events cho cả Alice và Bob
  - Positions được cập nhật đúng

### Tại sao chạy theo cách này?

1. **Order Matching:** Khi buy và sell orders match về giá và số lượng, chúng được fill đầy.
2. **Position Updates:** Positions của cả hai subaccounts được cập nhật sau matching.
3. **Event Emission:** Events được emit cho cả maker và taker.

---

## Tóm tắt Flow

### Place Order Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. SUBMIT ORDER                                             │
│    - Tạo MsgPlaceOrder                                       │
│    - Ký transaction                                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CHECKTX VALIDATION                                       │
│    - Validate order format                                   │
│    - Kiểm tra collateralization                             │
│    - Xác minh order parameters                               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. ORDER PLACEMENT                                          │
│    - Thêm order vào order book                               │
│    - Emit off-chain order place message                      │
│    - Emit off-chain order update message                     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. NEXT BLOCK                                               │
│    - Emit on-chain indexer block event                       │
│    - Cập nhật order state                                    │
└─────────────────────────────────────────────────────────────┘
```

### Match Order Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. PLACE FIRST ORDER                                         │
│    - Order được thêm vào order book                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. PLACE MATCHING ORDER                                      │
│    - Order match với existing order                         │
│    - Cả hai orders được fill đầy                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. EXECUTE TRADE                                            │
│    - Cập nhật subaccount positions                          │
│    - Emit subaccount update events                           │
│    - Xóa filled orders khỏi book                             │
└─────────────────────────────────────────────────────────────┘
```

### Điểm quan trọng

1. **Short-Term Orders:**
   - Expire ở cuối block hiện tại
   - GoodTilBlock chỉ định block number khi order expire
   - Phải được match trong cùng block hoặc chúng expire

2. **Order Matching:**
   - Orders match khi giá và số lượng tương thích
   - Maker-taker model: order đầu tiên là maker, matching order là taker
   - Cả hai orders được xóa khỏi book sau full match

3. **Off-chain Messages:**
   - Order place message: thông báo indexer về order mới
   - Order update message: thông báo indexer về fill amount
   - Transaction hash được bao gồm trong message headers

4. **On-chain Events:**
   - Subaccount update events: theo dõi thay đổi positions
   - Indexer block events: theo dõi block-level state changes
   - Events được emit trong block tiếp theo sau order placement

5. **State Tracking:**
   - Order state được track trong keeper
   - Fill amounts được track theo từng order
   - Positions được cập nhật sau matching

### Lý do thiết kế

1. **Efficiency:** Short-term orders cung cấp execution nhanh cho nhu cầu trading ngay lập tức.

2. **Liquidity:** Khuyến khích active trading bằng cách yêu cầu matching ngay lập tức.

3. **State Management:** Orders expire tự động, ngăn chặn stale orders làm lộn xộn book.

4. **Event Emission:** Cả off-chain và on-chain events đảm bảo complete state tracking.
