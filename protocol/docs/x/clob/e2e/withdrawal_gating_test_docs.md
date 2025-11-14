# Tài liệu Test: Withdrawal Gating E2E Tests

## Tổng quan

File test này xác minh chức năng **Withdrawal Gating** trong CLOB module. Khi một subaccount có Total Net Collateral (TNC) âm và không thể được deleverage, withdrawals và transfers từ market đó bị chặn (gated) để bảo vệ hệ thống. Test đảm bảo rằng:
1. Withdrawals bị gated khi negative TNC subaccounts tồn tại
2. Gating áp dụng cho isolated markets riêng biệt
3. Gating chặn withdrawals cho affected markets
4. Gating unblock khi negative TNC được giải quyết

---

## Test Function: TestWithdrawalGating_NegativeTncSubaccount_BlocksThenUnblocks

### Test Case 1: Withdrawals Gated - Non-Overlapping Bankruptcy Prices

### Đầu vào
- **Subaccounts:**
  - Carl: 1 BTC Short, 49,999 USD (negative TNC, undercollateralized)
  - Dave: 1 BTC Long, 50,000 USD (short)
  - Dave_Num1: 10,000 USD
- **Oracle Price:** $50,500 / BTC
- **Liquidation Order:** Dave bán 0.25 BTC ở $50,000
- **Liquidation:** Cố gắng nhưng deleveraging thất bại (non-overlapping bankruptcy prices)
- **Withdrawal:** Dave_Num1 cố gắng withdraw từ BTC market

### Đầu ra
- **Liquidation:** Thất bại (deleveraging không thể được thực hiện)
- **Carl State:** Vẫn có negative TNC
- **Withdrawals Gated:** BTC market withdrawals bị chặn
- **Error:** "WithdrawalsAndTransfersBlocked: failed to apply subaccount updates"
- **Gated Perpetual:** BTC perpetual ID được đánh dấu là gated
- **Negative TNC Seen At Block:** Block 4

### Tại sao chạy theo cách này?

1. **Negative TNC:** Carl có negative TNC (49,999 < 50,000 cần).
2. **Deleveraging Fails:** Không thể deleverage vì bankruptcy prices không overlap.
3. **System Protection:** Withdrawals bị gated để ngăn chặn further capital outflow.
4. **Market Isolation:** Gating áp dụng cho specific perpetual/market.

---

### Test Case 2: Withdrawals Gated - Isolated Market

### Đầu vào
- **Subaccounts:**
  - Carl: 1 ISO Short, 49 USD (negative TNC)
  - Dave: 1 ISO Long, 50 USD (short)
  - Alice: 1 ISO Long, 10,000 USD (isolated subaccount)
- **Oracle Price:** $50.5 / ISO
- **Liquidation:** Cố gắng nhưng deleveraging thất bại
- **Withdrawal:** Alice cố gắng withdraw từ ISO market

### Đầu ra
- **Withdrawals Gated:** ISO market withdrawals bị chặn cho isolated subaccounts
- **Gated Perpetual:** ISO perpetual ID được đánh dấu là gated
- **Error:** "WithdrawalsAndTransfersBlocked"

### Tại sao chạy theo cách này?

1. **Isolated Market:** ISO là isolated market với separate collateral pool.
2. **Isolated Subaccount:** Alice có isolated subaccount cho ISO market.
3. **Market-Specific Gating:** Gating chỉ áp dụng cho ISO market cho isolated subaccounts.
4. **Protection:** Ngăn chặn capital outflow từ isolated market khi negative TNC tồn tại.

---

### Test Case 3: Withdrawals Không Gated - Non-Isolated Subaccount

### Đầu vào
- **Subaccounts:**
  - Carl: 1 ISO Short, 49 USD (negative TNC)
  - Dave: 1 ISO Long, 50 USD (short)
  - Alice: 10,000 USD (non-isolated subaccount)
- **Oracle Price:** $50.5 / ISO
- **Liquidation:** Cố gắng nhưng deleveraging thất bại
- **Withdrawal:** Alice cố gắng withdraw (không phải từ ISO market)

### Đầu ra
- **Withdrawals Not Gated:** Alice có thể withdraw (không phải từ ISO market)
- **Gated Perpetual:** ISO perpetual ID vẫn được đánh dấu là gated
- **Selective Gating:** Chỉ ISO market withdrawals bị chặn

### Tại sao chạy theo cách này?

1. **Non-Isolated Subaccount:** Alice không có ISO position.
2. **Selective Gating:** Gating chỉ ảnh hưởng withdrawals từ gated market.
3. **Other Markets:** Withdrawals từ các markets khác vẫn được phép.

---

## Tóm tắt Flow

### Withdrawal Gating Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. PHÁT HIỆN NEGATIVE TNC                                   │
│    - Subaccount có TNC < 0                                  │
│    - Liquidations daemon xác định negative TNC accounts      │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CỐ GẮNG DELEVERAGING                                    │
│    - Cố gắng deleverage negative TNC account                │
│    - Kiểm tra overlapping bankruptcy prices                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. DELEVERAGING THẤT BẠI                                    │
│    - Bankruptcy prices không overlap                        │
│    - Không thể đóng position                                │
│    - Negative TNC vẫn tồn tại                               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. GATE WITHDRAWALS                                        │
│    - Đánh dấu perpetual là gated                            │
│    - Chặn withdrawals từ gated market                       │
│    - Ghi lại block khi negative TNC được thấy               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. CHẶN WITHDRAWAL ATTEMPTS                                │
│    - Từ chối withdrawal transactions                       │
│    - Trả về lỗi: "WithdrawalsAndTransfersBlocked"          │
│    - Áp dụng cho affected markets only                       │
└─────────────────────────────────────────────────────────────┘
```

### Gating Resolution

```
┌─────────────────────────────────────────────────────────────┐
│ 1. GIẢI QUYẾT NEGATIVE TNC                                 │
│    - Position đóng qua deleveraging                         │
│    - Hoặc position đóng qua matching                        │
│    - Hoặc collateral được thêm                              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. UNGATE WITHDRAWALS                                       │
│    - Không còn negative TNC accounts                        │
│    - Xóa gating từ perpetual                                │
│    - Cho phép withdrawals lại                                │
└─────────────────────────────────────────────────────────────┘
```

### Điểm quan trọng

1. **Negative TNC Detection:**
   - Subaccount TNC < 0
   - Được phát hiện bởi liquidations daemon
   - Được ghi lại per perpetual/market

2. **Deleveraging Failure:**
   - Bankruptcy prices không overlap
   - Không thể tìm counterparty để deleverage
   - Negative TNC không thể được giải quyết

3. **Gating Mechanism:**
   - Perpetual được đánh dấu là gated
   - Withdrawals bị chặn cho gated market
   - Transfers cũng bị chặn

4. **Market Isolation:**
   - Isolated markets được gated riêng biệt
   - Non-isolated subaccounts không bị ảnh hưởng bởi isolated market gating
   - Cross-market gating có thể xảy ra

5. **Block Tracking:**
   - Block khi negative TNC được thấy lần đầu
   - Được sử dụng cho gating duration tracking
   - Giúp xác định persistent issues

6. **Error Handling:**
   - Error message rõ ràng: "WithdrawalsAndTransfersBlocked"
   - Transaction bị từ chối tại CheckTx hoặc DeliverTx
   - State được bảo vệ khỏi capital outflow

### Lý do thiết kế

1. **System Protection:** Ngăn chặn capital outflow khi system đang gặp rủi ro.

2. **Risk Containment:** Gating giới hạn exposure đến negative TNC accounts.

3. **Market Isolation:** Isolated markets được bảo vệ riêng biệt.

4. **Fairness:** Chỉ ảnh hưởng withdrawals từ affected markets.

5. **Transparency:** Error messages rõ ràng và block tracking.
