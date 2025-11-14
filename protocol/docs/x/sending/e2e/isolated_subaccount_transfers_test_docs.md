# Tài liệu Test: Isolated Subaccount Transfers E2E Tests

## Tổng quan

File test này xác minh chức năng **Isolated Subaccount Transfer** trong Sending module. Isolated subaccounts là các subaccounts được cô lập với các perpetual markets cụ thể. Test đảm bảo rằng:
1. Transfers giữa isolated và non-isolated subaccounts hoạt động đúng
2. Collateral pools được cập nhật đúng khi transfer giữa các loại market khác nhau
3. Transfers giữa isolated subaccounts trong các markets khác nhau hoạt động đúng
4. Transfers thất bại khi collateral pools có insufficient funds
5. Transfers trong cùng isolated market không di chuyển collateral

---

## Test Function: TestTransfer_Isolated_Non_Isolated_Subaccounts

### Test Case 1: Thành công - Transfer từ Isolated đến Non-Isolated Subaccount

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: Isolated subaccount với 1 ISO long position, 10,000 USDC
  - Bob_Num0: Non-isolated subaccount với 10,000 USDC
- **Collateral Pools:**
  - Cross collateral pool: 10,000 USDC
  - Isolated market collateral pool: 10,000 USDC
- **Transfer:**
  - From: Alice_Num0 (isolated)
  - To: Bob_Num0 (non-isolated)
  - Amount: 100 USDC

### Đầu ra
- **Subaccounts:**
  - Alice_Num0: Số dư giảm 100 USDC
  - Bob_Num0: Số dư tăng 100 USDC
- **Collateral Pools:**
  - Cross collateral pool: Tăng 100 USDC (10,100 USDC)
  - Isolated market collateral pool: Giảm 100 USDC (9,900 USDC)

### Tại sao chạy theo cách này?

1. **Collateral Pool Movement:** Khi transfer từ isolated đến non-isolated, collateral di chuyển từ isolated pool đến cross pool.
2. **Isolation:** Isolated subaccounts có separate collateral pools cho mỗi market.
3. **Pool Updates:** Collateral pools phải được cập nhật để duy trì balance.

---

### Test Case 2: Thành công - Transfer từ Non-Isolated đến Isolated Subaccount

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: Isolated subaccount với 1 ISO long position, 10,000 USDC
  - Bob_Num0: Non-isolated subaccount với 10,000 USDC
- **Collateral Pools:**
  - Cross collateral pool: 10,000 USDC
  - Isolated market collateral pool: 10,000 USDC
- **Transfer:**
  - From: Bob_Num0 (non-isolated)
  - To: Alice_Num0 (isolated)
  - Amount: 100 USDC

### Đầu ra
- **Subaccounts:**
  - Alice_Num0: Số dư tăng 100 USDC
  - Bob_Num0: Số dư giảm 100 USDC
- **Collateral Pools:**
  - Cross collateral pool: Giảm 100 USDC (9,900 USDC)
  - Isolated market collateral pool: Tăng 100 USDC (10,100 USDC)

### Tại sao chạy theo cách này?

1. **Collateral Pool Movement:** Khi transfer từ non-isolated đến isolated, collateral di chuyển từ cross pool đến isolated pool.
2. **Reverse Flow:** Hướng ngược lại của Test Case 1.
3. **Pool Updates:** Collateral pools phải được cập nhật để duy trì balance.

---

### Test Case 3: Thành công - Transfer giữa Isolated Subaccounts trong Markets Khác nhau

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: Isolated subaccount trong ISO market, 1 ISO long, 10,000 USDC
  - Bob_Num0: Isolated subaccount trong ISO2 market, 1 ISO2 long, 10,000 USDC
- **Collateral Pools:**
  - ISO market collateral pool: 10,000 USDC
  - ISO2 market collateral pool: 10,000 USDC
- **Transfer:**
  - From: Alice_Num0 (ISO market)
  - To: Bob_Num0 (ISO2 market)
  - Amount: 100 USDC

### Đầu ra
- **Subaccounts:**
  - Alice_Num0: Số dư giảm 100 USDC
  - Bob_Num0: Số dư tăng 100 USDC
- **Collateral Pools:**
  - ISO market collateral pool: Giảm 100 USDC (9,900 USDC)
  - ISO2 market collateral pool: Tăng 100 USDC (10,100 USDC)

### Tại sao chạy theo cách này?

1. **Different Markets:** Isolated subaccounts trong các markets khác nhau có separate collateral pools.
2. **Pool Movement:** Collateral di chuyển từ một isolated pool đến pool khác.
3. **Isolation:** Mỗi isolated market duy trì collateral pool riêng của nó.

---

### Test Case 4: Thất bại - Insufficient Funds trong Isolated Collateral Pool

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: Isolated subaccount với 1 ISO long position, 10,000 USDC
  - Bob_Num0: Non-isolated subaccount với 10,000 USDC
- **Collateral Pools:**
  - Cross collateral pool: 10,000 USDC
  - Isolated market collateral pool: 0 USDC (rỗng)
- **Transfer:**
  - From: Alice_Num0 (isolated)
  - To: Bob_Num0 (non-isolated)
  - Amount: 100 USDC

### Đầu ra
- **DeliverTx:** FAIL
- **Error:** "insufficient funds"
- **Error Code:** `ErrInsufficientFunds`
- **Subaccounts:** Không thay đổi (transfer thất bại)
- **Collateral Pools:** Không thay đổi (transfer thất bại)

### Tại sao chạy theo cách này?

1. **Pool Balance:** Collateral pool phải có đủ funds cho transfer.
2. **Validation:** Hệ thống validate pool balance trước khi cho phép transfer.
3. **Failure Handling:** Transfer thất bại nếu pool có insufficient funds.

---

### Test Case 5: Thất bại - Insufficient Funds giữa Isolated Markets

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: Isolated subaccount trong ISO market, 1 ISO long, 10,000 USDC
  - Bob_Num0: Isolated subaccount trong ISO2 market, 1 ISO2 long, 10,000 USDC
- **Collateral Pools:**
  - ISO market collateral pool: 0 USDC (rỗng)
  - ISO2 market collateral pool: 10,000 USDC
- **Transfer:**
  - From: Alice_Num0 (ISO market)
  - To: Bob_Num0 (ISO2 market)
  - Amount: 100 USDC

### Đầu ra
- **DeliverTx:** FAIL
- **Error:** "insufficient funds"
- **Error Code:** `ErrInsufficientFunds`
- **Subaccounts:** Không thay đổi (transfer thất bại)
- **Collateral Pools:** Không thay đổi (transfer thất bại)

### Tại sao chạy theo cách này?

1. **Pool Balance:** Source collateral pool phải có đủ funds.
2. **Validation:** Hệ thống validate pool balance trước khi cho phép transfer.
3. **Failure Handling:** Transfer thất bại nếu source pool có insufficient funds.

---

### Test Case 6: Thành công - Transfer trong Cùng Isolated Market

### Đầu vào
- **Subaccounts:**
  - Alice_Num0: Isolated subaccount trong ISO market, 1 ISO long, 10,000 USDC
  - Bob_Num0: Isolated subaccount trong ISO market, 1 ISO long, 10,000 USDC
- **Collateral Pools:**
  - ISO market collateral pool: 10,000 USDC
- **Transfer:**
  - From: Alice_Num0 (ISO market)
  - To: Bob_Num0 (ISO market)
  - Amount: 100 USDC

### Đầu ra
- **Subaccounts:**
  - Alice_Num0: Số dư giảm 100 USDC
  - Bob_Num0: Số dư tăng 100 USDC
- **Collateral Pools:**
  - ISO market collateral pool: Không thay đổi (10,000 USDC)

### Tại sao chạy theo cách này?

1. **Same Market:** Cả hai subaccounts đều trong cùng isolated market.
2. **No Pool Movement:** Collateral ở trong cùng pool.
3. **Efficiency:** Không cần di chuyển collateral giữa các pools.

---

## Tóm tắt Flow

### Isolated Subaccount Transfer Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. TẠO TRANSFER MESSAGE                                     │
│    - Sender subaccount ID                                    │
│    - Receiver subaccount ID                                  │
│    - Asset ID và amount                                      │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. XÁC ĐỊNH LOẠI MARKET                                     │
│    - Kiểm tra nếu sender là isolated                        │
│    - Kiểm tra nếu receiver là isolated                       │
│    - Xác định loại market                                   │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. VALIDATE COLLATERAL POOLS                                 │
│    - Kiểm tra source pool balance                            │
│    - Xác minh sufficient funds                               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. CẬP NHẬT COLLATERAL POOLS                                │
│    - Nếu markets khác nhau: Di chuyển collateral giữa pools  │
│    - Nếu cùng market: Không di chuyển pool                   │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. CẬP NHẬT SUBACCOUNTS                                      │
│    - Giảm sender balance                                     │
│    - Tăng receiver balance                                   │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Transfer States:**
   ```
   Tạo Transfer → Validate Pools → Cập nhật Pools → Cập nhật Subaccounts → Hoàn thành
   ```

2. **Collateral Pool Updates:**
   ```
   Isolated → Non-Isolated: Isolated Pool ↓, Cross Pool ↑
   Non-Isolated → Isolated: Cross Pool ↓, Isolated Pool ↑
   Isolated → Isolated (Khác nhau): Source Pool ↓, Dest Pool ↑
   Isolated → Isolated (Cùng): Không thay đổi
   ```

### Điểm quan trọng

1. **Collateral Pools:**
   - Cross collateral pool: Cho non-isolated subaccounts
   - Isolated market pools: Một pool cho mỗi isolated market
   - Pools phải duy trì balance

2. **Transfer Rules:**
   - Markets khác nhau: Collateral di chuyển giữa pools
   - Cùng market: Không di chuyển pool
   - Insufficient funds: Transfer thất bại

3. **Validation:**
   - Pool balance phải đủ
   - Transfer amount phải hợp lệ
   - Subaccounts phải tồn tại

### Lý do thiết kế

1. **Isolation:** Isolated markets duy trì separate collateral pools cho risk management.

2. **Pool Management:** Collateral pools đảm bảo sufficient funds cho positions.

3. **Efficiency:** Transfers cùng market không di chuyển collateral để hiệu quả.

4. **Safety:** Validation ngăn chặn transfers khi pools có insufficient funds.
