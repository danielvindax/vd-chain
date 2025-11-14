# Tài liệu Test: Container Tests

## Tổng quan

File test này xác minh chức năng **Container Test Framework**. Container tests chạy một mạng lưới đầy đủ các node trong Docker containers, cho phép test end-to-end của blockchain. Framework test cung cấp các interface để tương tác với chain thông qua CometBFT và gRPC clients. Framework cũng chạy một HTTP server cho exchange price feeds.

---

## Test Function: TestPlaceOrder

### Test Case: Thành công - Đặt lệnh trên Network

### Đầu vào
- **Network:** Testnet đầy đủ với nhiều node chạy trong Docker containers
- **Node:** Alice node
- **Order:**
  - SubaccountId: Alice_Num0
  - ClobPairId: 0
  - Side: BUY
  - Quantums: 10,000,000
  - Subticks: 1,000,000
  - GoodTilBlock: 20

### Đầu ra
- **Giao dịch:** Broadcast thành công và được bao gồm trong block
- **Order:** Lệnh được đặt trên order book

### Tại sao chạy theo cách này?

1. **Test mạng đầy đủ:** Test đặt lệnh trong môi trường mạng thực với nhiều node.
2. **Docker Containers:** Mỗi node chạy trong một Docker container riêng, mô phỏng mạng thực.
3. **End-to-End:** Xác minh flow hoàn chỉnh từ submit giao dịch đến đặt lệnh.
4. **Tương tác mạng:** Test tương tác với mạng thông qua gRPC và CometBFT clients.

---

## Test Function: TestBankSend

### Test Case: Thành công - Gửi Token giữa các Tài khoản

### Đầu vào
- **Network:** Testnet đầy đủ với nhiều node
- **Node:** Alice node
- **Trạng thái ban đầu:**
  - Alice có số dư ban đầu
  - Bob có số dư ban đầu
- **Giao dịch:**
  - From: Bob
  - To: Alice
  - Amount: 1 USDC

### Đầu ra
- **Số dư ban đầu:** Xác minh với giá trị mong đợi
  - Số dư ban đầu của Alice khớp với mong đợi
  - Số dư ban đầu của Bob khớp với mong đợi
- **Số dư cuối:** Xác minh sau giao dịch
  - Số dư Alice tăng 1 USDC
  - Số dư Bob giảm 1 USDC

### Tại sao chạy theo cách này?

1. **Xác minh số dư:** Test rằng số dư được theo dõi và cập nhật đúng cách.
2. **Đầu ra mong đợi:** Sử dụng expect files để xác minh giá trị số dư chính xác.
3. **Consensus mạng:** Xác minh rằng giao dịch được propagate và bao gồm trong block đúng cách.
4. **Nhất quán State:** Đảm bảo tất cả node có state nhất quán sau giao dịch.

---

## Test Function: TestMarketPrices

### Test Case: Thành công - Giá thị trường cập nhật từ Exchange

### Đầu vào
- **Network:** Testnet đầy đủ với nhiều node
- **Giá Exchange được đặt trước khi bắt đầu:**
  - BTC-USD: 50,001
  - ETH-USD: 55,002
  - LINK-USD: 55,003
- **Node:** Alice node
- **Timeout:** 30 giây

### Đầu ra
- **Giá thị trường:** Giá được cập nhật để khớp với giá exchange
  - Giá BTC-USD khớp với mong đợi
  - Giá ETH-USD khớp với mong đợi
  - Giá LINK-USD khớp với mong đợi

### Tại sao chạy theo cách này?

1. **Tích hợp Price Feed:** Test tích hợp với external exchange price feeds.
2. **HTTP Server:** Framework chạy HTTP server cung cấp giá exchange.
3. **Cập nhật Oracle:** Xác minh rằng giá oracle được cập nhật từ exchange feeds.
4. **Moving Window:** Giá sử dụng moving window, vì vậy giá nên được đặt trước khi mạng bắt đầu.
5. **Polling:** Sử dụng polling với timeout để chờ giá cập nhật.

---

## Test Function: TestUpgrade

### Test Case: Thành công - Nâng cấp Network lên Phiên bản mới

### Đầu vào
- **Network:** Testnet với genesis trước nâng cấp
- **Node:** Alice node
- **Upgrade:** Nâng cấp lên phiên bản hiện tại
- **Upgrader:** Tài khoản Alice

### Đầu ra
- **Upgrade:** Thực thi thành công
- **Network:** Chạy trên phiên bản mới

### Tại sao chạy theo cách này?

1. **Test nâng cấp:** Test cơ chế nâng cấp mạng.
2. **Genesis trước nâng cấp:** Sử dụng genesis state đặc biệt cho test trước nâng cấp.
3. **Quản lý phiên bản:** Xác minh rằng mạng có thể nâng cấp lên phiên bản mới.
4. **Migration State:** Đảm bảo state được migrate đúng cách trong quá trình nâng cấp.

---

## Tóm tắt Flow

### Container Test Framework Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. TẠO TESTNET                                              │
│    - Khởi tạo testnet với các node                          │
│    - Cấu hình Docker containers                             │
│    - Thiết lập HTTP server cho price feeds                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CẤU HÌNH TRƯỚC KHI BẮT ĐẦU                              │
│    - Đặt giá exchange (nếu cần)                            │
│    - Cấu hình genesis state                                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. BẮT ĐẦU NETWORK                                          │
│    - Khởi động Docker containers                            │
│    - Chờ các node sync                                      │
│    - Xác minh mạng đã sẵn sàng                             │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. TƯƠNG TÁC VỚI NETWORK                                   │
│    - Query state (số dư, giá, v.v.)                        │
│    - Broadcast giao dịch                                   │
│    - Chờ block                                              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. XÁC MINH KẾT QUẢ                                        │
│    - So sánh với đầu ra mong đợi                           │
│    - Xác minh thay đổi state                                │
│    - Dọn dẹp containers                                     │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Trạng thái Network:**
   ```
   Chưa bắt đầu → Đang bắt đầu → Đang chạy → Đã dọn dẹp
   ```

2. **Cập nhật giá:**
   ```
   Giá Exchange → HTTP Server → Oracle → Giá thị trường
   ```

3. **Flow giao dịch:**
   ```
   Broadcast → Mempool → Block → Cập nhật State
   ```

### Điểm quan trọng

1. **Docker Containers:**
   - Mỗi node chạy trong một Docker container riêng
   - Mô phỏng môi trường mạng thực
   - Cho phép test consensus mạng và propagation

2. **Tích hợp Price Feed:**
   - HTTP server cung cấp giá exchange
   - Giá nên được đặt trước khi mạng bắt đầu
   - Oracle sử dụng moving window cho cập nhật giá

3. **Đầu ra mong đợi:**
   - Test sử dụng expect files để xác minh đầu ra chính xác
   - Sử dụng flag `-accept` để cập nhật expect files
   - Đảm bảo kết quả test xác định

4. **Tương tác mạng:**
   - Query: Đọc state từ các node
   - BroadcastTx: Submit giao dịch đến mạng
   - Wait: Chờ block được tạo

5. **Dọn dẹp:**
   - Luôn dọn dẹp containers sau test
   - Sử dụng `defer testnet.MustCleanUp()`
   - Ngăn chặn rò rỉ tài nguyên

6. **Test nâng cấp:**
   - Test cơ chế nâng cấp mạng
   - Sử dụng genesis state trước nâng cấp
   - Xác minh migration state

### Lý do thiết kế

1. **End-to-End Testing:** Container tests cung cấp test mạng đầy đủ, không chỉ unit tests.

2. **Môi trường thực:** Docker containers mô phỏng điều kiện mạng thực.

3. **Integration Testing:** Test tích hợp giữa các component (node, price feeds, v.v.).

4. **Xác định:** Expect files đảm bảo test xác định và có thể tái tạo.

5. **Linh hoạt:** Framework cho phép test các scenario khác nhau (nâng cấp, cập nhật giá, v.v.).
