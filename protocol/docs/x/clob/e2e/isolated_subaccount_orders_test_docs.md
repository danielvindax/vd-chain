# Tài liệu Test: Isolated Subaccount Orders E2E Tests

## Tổng quan

File test này xác minh chức năng **Isolated Subaccount Order** trong CLOB module. Isolated subaccounts là các subaccounts chỉ có thể trade trong các isolated markets cụ thể. Test đảm bảo rằng:
1. Isolated subaccounts không thể đặt orders cho cross-market (non-isolated) perpetuals
2. Isolated subaccounts có thể đặt orders cho isolated market của họ
3. Isolated subaccounts có thể match với các subaccounts khác trong isolated market
4. Isolated subaccounts sử dụng isolated collateral pool

---

## Test Function: TestIsolatedSubaccountOrders

### Test Case 1: Thất bại - Isolated Subaccount Không thể Đặt Cross-Market Order

### Đầu vào
- **Subaccounts:**
  - Alice: 1 ISO Long, 10,000 USD (isolated to ISO market)
  - Bob: 10,000 USD
- **Perpetuals:**
  - BTC-USD (market 0)
  - ETH-USD (market 1)
  - ISO-USD (market 3, isolated)
- **Orders:**
  - Alice: Cố gắng mua 5 BTC ở giá 10 (CLOB 0 - cross-market)
  - Bob: Bán 5 BTC ở giá 10 (CLOB 0)

### Đầu ra
- **Alice Order:** Bị từ chối (invalid cho isolated subaccount)
- **Bob Order:** Được chấp nhận
- **Orders Filled:** Không có (order của Alice invalid)
- **Subaccounts:** Không thay đổi

### Tại sao chạy theo cách này?

1. **Isolation Constraint:** Subaccount của Alice bị isolated chỉ cho ISO market.
2. **Cross-Market Rejection:** Không thể đặt orders cho BTC market (market 0).
3. **Validation:** Hệ thống validate subaccount có thể trade trong requested market.
4. **Protection:** Ngăn chặn isolated subaccounts trade cross-markets.

---

### Test Case 2: Thành công - Isolated Subaccount Đặt Order trong Isolated Market

### Đầu vào
- **Subaccounts:**
  - Alice: 1 ISO Long, 10,000 USD (isolated to ISO market)
  - Bob: 10,000 USD
- **Perpetuals:**
  - ISO-USD (market 3, isolated)
- **Orders:**
  - Alice: Mua 1 ISO ở giá 10 (CLOB 3 - isolated market)
  - Bob: Bán 1 ISO ở giá 10 (CLOB 3)

### Đầu ra
- **Both Orders:** Được chấp nhận
- **Orders Matched:** Cả hai orders match đầy đủ
- **Alice Position:** 2 ISO Long (1 existing + 1 từ match)
- **Bob Position:** 1 ISO Short

### Tại sao chạy theo cách này?

1. **Isolated Market:** Alice có thể trade trong ISO market (isolated market của cô ấy).
2. **Order Matching:** Orders match bình thường trong isolated market.
3. **Position Updates:** Positions được cập nhật đúng sau match.

---

### Test Case 3: Thành công - Isolated Subaccount Sử dụng Isolated Collateral Pool

### Đầu vào
- **Subaccounts:**
  - Alice: 1 ISO Long, 10,000 USD (isolated)
  - Bob: 10,000 USD
- **Collateral Pools:**
  - Main pool: 10,000 USD
  - ISO isolated pool: 10,000 USD
- **Orders:**
  - Alice: Mua 1 ISO ở giá 10
  - Bob: Bán 1 ISO ở giá 10

### Đầu ra
- **Orders Matched:** Thành công
- **Collateral Pool:** ISO isolated pool được sử dụng cho Alice
- **Main Pool:** Không bị ảnh hưởng

### Tại sao chạy theo cách này?

1. **Isolated Pool:** Isolated subaccounts sử dụng separate collateral pool.
2. **Pool Isolation:** Main pool và isolated pools là riêng biệt.
3. **Risk Isolation:** Isolated markets có isolated risk.

---

## Tóm tắt Flow

### Isolated Subaccount Order Validation

```
┌─────────────────────────────────────────────────────────────┐
│ 1. NHẬN ORDER                                              │
│    - Order chỉ định CLOB pair / perpetual                    │
│    - Order chỉ định subaccount ID                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. KIỂM TRA SUBACCOUNT ISOLATION                            │
│    - Query isolated market của subaccount (nếu có)          │
│    - Kiểm tra nếu subaccount bị isolated                    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. VALIDATE MARKET                                          │
│    - Nếu isolated: Kiểm tra nếu order market = isolated market │
│    - Nếu không isolated: Cho phép bất kỳ market nào          │
│    - Từ chối nếu isolated subaccount cố gắng cross-market   │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. XỬ LÝ ORDER                                              │
│    - Nếu hợp lệ: Xử lý bình thường                          │
│    - Sử dụng isolated collateral pool nếu applicable        │
│    - Cập nhật positions trong isolated market                │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Subaccount Isolation:**
   ```
   Non-Isolated: Có thể trade bất kỳ market nào
   Isolated: Chỉ có thể trade isolated market
   ```

2. **Collateral Pools:**
   ```
   Main Pool: Cho non-isolated markets
   Isolated Pool: Cho isolated markets (mỗi market)
   ```

### Điểm quan trọng

1. **Isolation Constraint:**
   - Isolated subaccounts chỉ có thể trade isolated market của họ
   - Không thể đặt orders cho các markets khác
   - Validation tại CheckTx

2. **Market Matching:**
   - Isolated subaccounts có thể match với bất kỳ subaccount nào trong isolated market
   - Matching hoạt động bình thường trong isolated market
   - Cross-market matching bị ngăn chặn

3. **Collateral Pools:**
   - Isolated markets có separate collateral pools
   - Isolated subaccounts sử dụng isolated pool
   - Risk isolated per market

4. **Position Management:**
   - Positions được track theo market
   - Isolated positions riêng biệt với cross-market positions
   - Collateral requirements per pool

5. **Validation:**
   - CheckTx validate market access
   - Early rejection cho invalid orders
   - Error messages rõ ràng

### Lý do thiết kế

1. **Risk Isolation:** Isolated markets ngăn chặn risk spillover đến main system.

2. **Capital Efficiency:** Isolated markets có thể có risk parameters khác nhau.

3. **Market Segregation:** Ngăn chặn isolated subaccounts ảnh hưởng đến main markets.

4. **Flexibility:** Cho phép các markets mới với risk profiles khác nhau.
