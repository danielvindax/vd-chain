# Tài liệu Test: App Query E2E Tests

## Tổng quan

File test này xác minh chức năng **Parallel Query** trong ứng dụng. Test đảm bảo rằng ứng dụng có thể xử lý các query đồng thời an toàn mà không có data race. Test sử dụng Go's race detector để xác minh thread safety.

---

## Test Function: TestParallelQuery

### Test Case: Thành công - Query song song không có Data Race

### Đầu vào
- **Thao tác đồng thời:**
  - Thread 1: Tiến block (block 2-49)
  - Thread 2: Query app/version lặp lại
  - Thread 3: Query store/blocktime/key trực tiếp
  - Thread 4: Query gRPC PreviousBlockInfo
- **Đồng bộ:** Atomic boolean để điều phối thread
- **Thực thi:** Tất cả thao tác chạy đồng thời cho đến khi đạt giới hạn block

### Đầu ra
- **Không có Data Race:** Test pass với flag `-race` được bật
- **Kết quả Query:** Tất cả query trả về kết quả hợp lệ
- **Tính đơn điệu Height:** Query heights tăng đơn điệu
- **Nhất quán:** Store queries và gRPC queries trả về dữ liệu nhất quán

### Tại sao chạy theo cách này?

1. **Test đồng thời:** Test rằng ứng dụng xử lý các query đồng thời đúng cách.
2. **Phát hiện Race:** Sử dụng Go's race detector để tìm data race.
3. **Nhiều loại Query:** Test các đường dẫn query khác nhau (app, store, gRPC).
4. **Stress Test:** Query đồng thời trong khi block tiến stress test hệ thống.
5. **Điều phối Atomic:** Sử dụng atomic boolean để tối đa hóa khả năng data race.

---

## Tóm tắt Flow

### Parallel Query Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. BẮT ĐẦU CÁC THREAD ĐỒNG THỜI                            │
│    - Thread tiến block                                      │
│    - Thread query app/version                                │
│    - Thread query store                                      │
│    - Thread query gRPC                                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. THỰC THI ĐỒNG THỜI                                        │
│    - Block tiến trong khi query thực thi                    │
│    - Không có đồng bộ giữa các thread                       │
│    - Atomic boolean điều phối hoàn thành                    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. THỰC THI QUERY                                            │
│    - App/version: Query phiên bản ứng dụng                  │
│    - Store: Query blocktime store trực tiếp                  │
│    - gRPC: Query PreviousBlockInfo qua gRPC                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. XÁC MINH                                                 │
│    - Tất cả query trả về kết quả hợp lệ                     │
│    - Heights tăng đơn điệu                                   │
│    - Không phát hiện data race                              │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Trạng thái Query:**
   ```
   Yêu cầu Query → Thực thi Query → Trả về Kết quả → Xác minh Kết quả
   ```

2. **Thực thi đồng thời:**
   ```
   Thread 1: Tiến Block
   Thread 2: Query App/Version
   Thread 3: Query Store
   Thread 4: Query gRPC
   ```

### Điểm quan trọng

1. **Loại Query:**
   - App/Version: Query phiên bản ứng dụng
   - Store: Query store trực tiếp (blocktime key)
   - gRPC: Query dịch vụ gRPC (PreviousBlockInfo)

2. **Đồng thời:**
   - Nhiều thread thực thi query đồng thời
   - Block tiến trong khi query thực thi
   - Không có đồng bộ giữa các thread query

3. **Phát hiện Race:**
   - Sử dụng flag `-race` của Go để phát hiện data race
   - Atomic boolean tối đa hóa khả năng race
   - Wait group điều phối hoàn thành thread

4. **Xác minh:**
   - Heights phải tăng đơn điệu
   - Store và gRPC queries phải trả về dữ liệu nhất quán
   - Tất cả query phải trả về kết quả hợp lệ

### Lý do thiết kế

1. **Thread Safety:** Ứng dụng phải thread-safe cho các query đồng thời.

2. **Phát hiện Race:** Go's race detector giúp tìm data race trong quá trình test.

3. **Stress Testing:** Query đồng thời trong khi block tiến stress test hệ thống.

4. **Nhiều đường dẫn:** Test các đường dẫn query khác nhau để đảm bảo tất cả đều thread-safe.
