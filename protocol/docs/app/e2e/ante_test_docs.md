# Tài liệu Test: App Ante Handler E2E Tests

## Tổng quan

File test này xác minh chức năng **Parallel Ante Handler** trong ứng dụng. Test đảm bảo rằng ante handler có thể xử lý các giao dịch CLOB và các module khác đồng thời mà không có data race. Test sử dụng Go's race detector để xác minh thread safety.

---

## Test Function: TestParallelAnteHandler_ClobAndOther

### Test Case: Thành công - Giao dịch CLOB và Transfer song song

### Đầu vào
- **Tài khoản:** 10 tài khoản ngẫu nhiên
- **Thao tác đồng thời:**
  - Thread 1: Tiến block (block 2-49)
  - Threads 2-11: Rút tiền từ subaccounts (một thread mỗi tài khoản)
  - Threads 12-21: Đặt và hủy lệnh CLOB (một thread mỗi tài khoản)
- **Giao dịch:**
  - Withdraw: `MsgWithdrawFromSubaccount` (1 USDC mỗi giao dịch)
  - Place Order: `MsgPlaceOrder` (1 quantum, giá 10)
  - Cancel Order: `MsgCancelOrderShortTerm`
- **Thực thi:** Tất cả lời gọi CheckTx được thực thi đồng thời
- **Block:** Tiến đến block 50

### Đầu ra
- **Không có Data Race:** Test pass với flag `-race` được bật
- **Giao dịch:** Tất cả giao dịch pass CheckTx
- **Trạng thái cuối:** Số dư subaccount khớp với giá trị mong đợi
  - Balance = Initial Balance - (Transfer Count × 1 USDC)

### Tại sao chạy theo cách này?

1. **Test đồng thời:** Test rằng ante handler xử lý các giao dịch đồng thời đúng cách.
2. **Phát hiện Race:** Sử dụng Go's race detector để tìm data race.
3. **Thao tác hỗn hợp:** Test cả giao dịch CLOB và sending module đồng thời.
4. **Stress Test:** Nhiều tài khoản với thao tác đồng thời stress test hệ thống.
5. **Cô lập tài khoản:** Mỗi tài khoản có thread riêng để tối đa hóa contention.

---

## Tóm tắt Flow

### Parallel Ante Handler Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. BẮT ĐẦU CÁC THREAD ĐỒNG THỜI                            │
│    - Thread tiến block                                      │
│    - Threads withdraw (một thread mỗi tài khoản)          │
│    - Threads lệnh CLOB (một thread mỗi tài khoản)           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CHECKTX ĐỒNG THỜI                                        │
│    - Giao dịch withdraw được thực thi đồng thời            │
│    - Giao dịch lệnh CLOB được thực thi đồng thời            │
│    - Không có đồng bộ giữa các thread                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. XỬ LÝ ANTE HANDLER                                       │
│    - Ante handler xử lý giao dịch                          │
│    - Xác thực chữ ký                                        │
│    - Xác thực account sequence                              │
│    - Trừ phí                                                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. TIẾN BLOCK                                               │
│    - Tất cả giao dịch được bao gồm trong block              │
│    - State được cập nhật                                     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. XÁC MINH                                                 │
│    - Số dư subaccount khớp với mong đợi                     │
│    - Không phát hiện data race                              │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Trạng thái giao dịch:**
   ```
   Tạo Giao dịch → CheckTx → Ante Handler → DeliverTx → Cập nhật State
   ```

2. **Thực thi đồng thời:**
   ```
   Thread 1: Tiến Block
   Threads 2-11: Giao dịch Withdraw
   Threads 12-21: Giao dịch Lệnh CLOB
   ```

3. **Trạng thái tài khoản:**
   ```
   Số dư ban đầu → Giao dịch Withdraw → Số dư cuối
   ```

### Điểm quan trọng

1. **Ante Handler:**
   - Xử lý giao dịch trước DeliverTx
   - Xác thực chữ ký
   - Kiểm tra account sequences
   - Trừ phí

2. **Đồng thời:**
   - Nhiều thread thực thi CheckTx đồng thời
   - Block tiến trong khi giao dịch thực thi
   - Không có đồng bộ giữa các thread giao dịch

3. **Phát hiện Race:**
   - Sử dụng flag `-race` của Go để phát hiện data race
   - Atomic boolean tối đa hóa khả năng race
   - Wait group điều phối hoàn thành thread

4. **Xác minh:**
   - Số dư subaccount phải khớp với giá trị mong đợi
   - Balance = Initial - (Transfer Count × Transfer Amount)
   - Tất cả giao dịch phải pass CheckTx

### Lý do thiết kế

1. **Thread Safety:** Ante handler phải thread-safe cho các giao dịch đồng thời.

2. **Phát hiện Race:** Go's race detector giúp tìm data race trong quá trình test.

3. **Stress Testing:** Giao dịch đồng thời trong khi block tiến stress test hệ thống.

4. **Thao tác hỗn hợp:** Test cả giao dịch CLOB và module khác để đảm bảo tương thích.
