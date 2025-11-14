# Tài liệu Test: Perpetuals Authorities E2E Tests

## Tổng quan

File test này xác minh chức năng **Authority Management** trong Perpetuals module. Authorities là các địa chỉ có quyền thực hiện các thao tác đặc quyền trong perpetuals module. Test đảm bảo rằng:
1. Governance module được công nhận là authority
2. DelayMsg module được công nhận là authority
3. Các địa chỉ không hợp lệ không được công nhận là authorities
4. Authority checks hoạt động đúng

---

## Test Function: TestHasAuthority

### Test Case 1: Thành công - Governance Module là Authority

### Đầu vào
- **Authority Address:** Địa chỉ Governance module
  - Address: `authtypes.NewModuleAddress(govtypes.ModuleName)`
- **Check:** `HasAuthority(authorityAddress)`

### Đầu ra
- **Result:** `true`
- **Authority:** Governance module được công nhận là authority

### Tại sao chạy theo cách này?

1. **Governance Authority:** Governance module cần authority để cập nhật perpetual parameters qua proposals.
2. **Module Address:** Module addresses được derive từ module names.
3. **Permission Check:** Hệ thống kiểm tra nếu address có authority trước khi cho phép thao tác đặc quyền.

---

### Test Case 2: Thành công - DelayMsg Module là Authority

### Đầu vào
- **Authority Address:** Địa chỉ DelayMsg module
  - Address: `authtypes.NewModuleAddress(delaymsgtypes.ModuleName)`
- **Check:** `HasAuthority(authorityAddress)`

### Đầu ra
- **Result:** `true`
- **Authority:** DelayMsg module được công nhận là authority

### Tại sao chạy theo cách này?

1. **DelayMsg Authority:** DelayMsg module cần authority để thực thi delayed messages cập nhật perpetual parameters.
2. **Delayed Updates:** Cho phép cập nhật parameters theo lịch trình qua delayed messages.
3. **Module Integration:** DelayMsg module tích hợp với perpetuals cho parameter updates.

---

### Test Case 3: Thất bại - Random Invalid Address Không phải Authority

### Đầu vào
- **Authority Address:** Random invalid address
  - Address: `"random"`
- **Check:** `HasAuthority(authorityAddress)`

### Đầu ra
- **Result:** `false`
- **Authority:** Random address không được công nhận là authority

### Tại sao chạy theo cách này?

1. **Bảo mật:** Chỉ các địa chỉ được ủy quyền mới có thể thực hiện thao tác đặc quyền.
2. **Validation:** Hệ thống validate authority trước khi cho phép thao tác.
3. **Access Control:** Ngăn chặn truy cập trái phép vào perpetual parameters.

---

## Tóm tắt Flow

### Authority Check Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. AUTHORITY REQUEST                                         │
│    - Address được cung cấp cho authority check               │
│    - Hệ thống kiểm tra nếu address được ủy quyền            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. AUTHORITY VALIDATION                                     │
│    - Kiểm tra nếu address khớp với governance module        │
│    - Kiểm tra nếu address khớp với delaymsg module           │
│    - Kiểm tra nếu address trong danh sách authorized         │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. RESULT                                                    │
│    - Trả về true nếu được ủy quyền                          │
│    - Trả về false nếu không được ủy quyền                    │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Authority States:**
   ```
   Address → Authority Check → Authorized / Not Authorized
   ```

2. **Authorized Addresses:**
   ```
   Governance Module Address → Authorized
   DelayMsg Module Address → Authorized
   Other Addresses → Not Authorized
   ```

### Điểm quan trọng

1. **Module Authorities:**
   - Governance module: Có thể cập nhật parameters qua proposals
   - DelayMsg module: Có thể thực thi delayed parameter updates
   - Other modules: Không được ủy quyền theo mặc định

2. **Authority Check:**
   - `HasAuthority(address)` kiểm tra nếu address được ủy quyền
   - Trả về kết quả boolean
   - Được sử dụng trước khi cho phép thao tác đặc quyền

3. **Bảo mật:**
   - Chỉ các địa chỉ được ủy quyền mới có thể thực hiện thao tác đặc quyền
   - Ngăn chặn cập nhật parameters trái phép
   - Đảm bảo tính toàn vẹn hệ thống

### Lý do thiết kế

1. **Access Control:** Authority system cung cấp fine-grained access control cho perpetual parameters.

2. **Module Integration:** Cho phép các modules khác (governance, delaymsg) tương tác với perpetuals.

3. **Bảo mật:** Ngăn chặn truy cập trái phép vào các system parameters quan trọng.

4. **Flexibility:** Có thể được mở rộng để thêm nhiều authorized addresses nếu cần.
