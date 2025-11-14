# Tài liệu Test: Affiliates Register E2E Tests

## Tổng quan

File test này xác minh chức năng **Affiliate Registration** trong Affiliates module. Affiliates cho phép users đăng ký referral relationships nơi một referee (user) có thể được liên kết với một affiliate (referrer). Test đảm bảo rằng:
1. Chỉ referee có thể đăng ký affiliate relationship
2. Affiliate không thể đăng ký chính họ
3. Các địa chỉ không liên quan không thể đăng ký relationship
4. Signature validation hoạt động đúng

---

## Test Function: TestRegisterAffiliateInvalidSigner

### Test Case 1: Thành công - Valid Signer (Referee)

### Đầu vào
- **Referee:** Địa chỉ của Bob
- **Affiliate:** Địa chỉ của Alice
- **Signer:** Private key của Bob (referee)
- **Message:** `MsgRegisterAffiliate`
  - Referee: Địa chỉ của Bob
  - Affiliate: Địa chỉ của Alice

### Đầu ra
- **CheckTx:** SUCCESS
- **Transaction:** Được chấp nhận
- **Relationship:** Affiliate relationship được đăng ký

### Tại sao chạy theo cách này?

1. **Referee Authorization:** Chỉ referee có thể đăng ký affiliate relationship.
2. **Self-Registration:** Referee phải ký transaction bằng chính họ.
3. **Relationship Creation:** Tạo referral relationship giữa referee và affiliate.

---

### Test Case 2: Thất bại - Invalid Signer (Affiliate)

### Đầu vào
- **Referee:** Địa chỉ của Bob
- **Affiliate:** Địa chỉ của Alice
- **Signer:** Private key của Alice (affiliate, signer sai)
- **Message:** `MsgRegisterAffiliate`
  - Referee: Địa chỉ của Bob
  - Affiliate: Địa chỉ của Alice

### Đầu ra
- **CheckTx:** FAIL
- **Error:** "pubKey does not match signer address"
- **Transaction:** Bị từ chối

### Tại sao chạy theo cách này?

1. **Authorization:** Chỉ referee có thể đăng ký relationship.
2. **Bảo mật:** Ngăn chặn affiliates đăng ký chính họ.
3. **Signature Validation:** Hệ thống validate rằng signer khớp với referee address.

---

### Test Case 3: Thất bại - Invalid Signer (Non-Related Address)

### Đầu vào
- **Referee:** Địa chỉ của Bob
- **Affiliate:** Địa chỉ của Alice
- **Signer:** Private key của Carl (địa chỉ không liên quan, signer sai)
- **Message:** `MsgRegisterAffiliate`
  - Referee: Địa chỉ của Bob
  - Affiliate: Địa chỉ của Alice

### Đầu ra
- **CheckTx:** FAIL
- **Error:** "pubKey does not match signer address"
- **Transaction:** Bị từ chối

### Tại sao chạy theo cách này?

1. **Authorization:** Chỉ referee có thể đăng ký relationship.
2. **Bảo mật:** Ngăn chặn bên thứ ba đăng ký relationships.
3. **Signature Validation:** Hệ thống validate rằng signer khớp với referee address.

---

## Tóm tắt Flow

### Affiliate Registration Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. TẠO MESSAGE                                               │
│    - Referee address                                         │
│    - Affiliate address                                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. KÝ TRANSACTION                                            │
│    - Referee ký với private key của họ                       │
│    - Signature được bao gồm trong transaction                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. CHECKTX VALIDATION                                        │
│    - Validate signature khớp với referee address            │
│    - Kiểm tra nếu signer là referee                          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. ĐĂNG KÝ                                                   │
│    - Nếu hợp lệ: Đăng ký affiliate relationship             │
│    - Nếu không hợp lệ: Từ chối transaction                  │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Registration States:**
   ```
   Chưa Đăng ký → Yêu cầu Đăng ký → Validated → Đã Đăng ký / Bị Từ chối
   ```

2. **Signer Validation:**
   ```
   Transaction → Trích xuất Signer → So sánh với Referee → Authorized / Unauthorized
   ```

### Điểm quan trọng

1. **Referee Authorization:**
   - Chỉ referee có thể đăng ký affiliate relationship
   - Referee phải ký transaction
   - Signature phải khớp với referee address

2. **Bảo mật:**
   - Affiliate không thể đăng ký chính họ
   - Bên thứ ba không thể đăng ký relationships
   - Signature validation ngăn chặn đăng ký trái phép

3. **Relationship Creation:**
   - Tạo referral relationship giữa referee và affiliate
   - Relationship có thể được sử dụng cho rewards và tracking
   - One-to-one relationship (một referee, một affiliate)

### Lý do thiết kế

1. **User Control:** Referee kiểm soát đăng ký affiliate của chính họ.

2. **Bảo mật:** Signature validation ngăn chặn đăng ký trái phép.

3. **Simplicity:** Mô hình relationship one-to-one đơn giản.

4. **Flexibility:** Cho phép users chọn affiliate của họ.
