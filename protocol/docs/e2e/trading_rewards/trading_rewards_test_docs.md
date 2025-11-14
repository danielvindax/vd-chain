# Tài liệu Test: Trading Rewards E2E Tests

## Tổng quan

File test này xác minh cơ chế phân phối **Trading Rewards**. Trading rewards được phân phối cho traders dựa trên hoạt động trading của họ (phí đã trả). Test đảm bảo rằng:
1. Rewards được tính toán dựa trên trading fees
2. Rewards được phân phối từ treasury account
3. Chỉ một tài khoản taker nhận rewards mỗi block
4. Nhiều tài khoản có thể nhận rewards trong cùng block
5. Rewards multiplier ảnh hưởng đến số tiền phân phối
6. Vesting tokens được chuyển từ vester đến treasury

---

## Test Function: TestTradingRewards

### Test Case 1: Mỗi Block, Chỉ Một Tài khoản Taker Nhận Rewards

### Đầu vào
- **Vest Entry:**
  - VesterAccount: rewards_vester
  - TreasuryAccount: rewards_treasury
  - StartTime: Oct 01 2023 04:00:00
  - EndTime: Oct 01 2028 04:00:00
- **Rewards Params:**
  - FeeMultiplierPpm: 990_000 (99%)
  - MarketId: 30 (rewards token)
- **Orders:**
  - Block 2: Bob (maker) bán, Alice (taker) mua 1 BTC ở 28,003
  - Block 13: Alice (maker) mua, Bob (taker) bán 1 BTC ở 28,003
- **Oracle Prices:**
  - Rewards token: $1.95
  - BTC: $28,003

### Đầu ra
- **Block 0 (Vest Start):**
  - Vester: 200 triệu tokens
  - Treasury: 0 tokens
- **Block 1:**
  - Vester: ~199,999,997.47 tokens (vesting đã bắt đầu)
  - Treasury: ~2.53 tokens (vested)
- **Block 2:**
  - Vester: ~199,999,994.93 tokens
  - Treasury: ~3.07 tokens (sau khi phân phối ~1.99 cho Alice)
  - Alice: Số dư ban đầu + ~1.99 tokens (rewards)
  - Bob: Số dư ban đầu (không có rewards - là maker)
- **Block 13:**
  - Vester: ~199,999,967.05 tokens
  - Treasury: ~28.96 tokens
  - Bob: Số dư ban đầu + ~1.99 tokens (rewards - là taker)
  - Alice: Số dư ban đầu + ~1.99 tokens (từ block 2)

### Tại sao chạy theo cách này?

1. **Vesting:** Tokens vest từ vester account đến treasury account theo thời gian.
2. **Rewards Calculation:** Rewards = (TakerFee - MakerRebate - TakerFeeRevShare) × FeeMultiplierPpm
   - Cho 1 BTC ở $28,003: ($14.0015 - $3.08033 - $7.00075) × 0.99 = $3.8812158
   - Reward tokens = $3.8812158 / $1.95 = ~1.99 tokens
3. **Taker Only:** Chỉ tài khoản taker nhận rewards, không phải maker.
4. **One Per Block:** Chỉ một tài khoản taker nhận rewards mỗi block (taker đầu tiên trong block đó).

---

### Test Case 2: Nhiều Tài khoản Nhận Rewards

### Đầu vào
- **Vest Entry:** Giống Test Case 1
- **Rewards Params:** Giống Test Case 1
- **Orders (Block 10):**
  - BTC: Bob (maker) bán 2 BTC, Alice (taker) mua 2 BTC
  - BTC: Alice (maker) mua 2 BTC, Bob (taker) bán 2 BTC
  - ETH: Carl (maker) mua 20 ETH, Dave (taker) bán 20 ETH
  - ETH: Dave (maker) bán 20 ETH, Carl (taker) mua 20 ETH
- **Oracle Prices:**
  - Rewards token: $1.95
  - BTC: $28,003
  - ETH: $1,605

### Đầu ra
- **Block 0:**
  - Vester: 200 triệu tokens
  - Treasury: 0 tokens
- **Block 10:**
  - Vester: ~199,999,974.66 tokens
  - Treasury: ~12.82 tokens (sau khi phân phối rewards)
  - Alice: Số dư ban đầu + ~3.98 tokens (rewards từ BTC trading)
  - Bob: Số dư ban đầu + ~3.98 tokens (rewards từ BTC trading)
  - Carl: Số dư ban đầu + ~2.28 tokens (rewards từ ETH trading)
  - Dave: Số dư ban đầu + ~2.28 tokens (rewards từ ETH trading)

### Tại sao chạy theo cách này?

1. **Multiple Takers:** Nhiều tài khoản có thể nhận rewards trong cùng block.
2. **Rewards Per Trade:** Mỗi taker nhận rewards dựa trên trading fees của họ.
3. **Different Markets:** Rewards được tính riêng cho mỗi market (BTC, ETH).
4. **Total Distribution:** Tổng rewards phân phối = tổng rewards cá nhân.

---

### Test Case 3: Rewards Fee Multiplier = 0, Không Phân phối Rewards

### Đầu vào
- **Vest Entry:** Giống Test Case 1
- **Rewards Params:**
  - FeeMultiplierPpm: 0 (0% - không có rewards)
- **Orders (Block 10):**
  - BTC: Bob (maker) bán 2 BTC, Alice (taker) mua 2 BTC
  - ETH: Carl (maker) mua 20 ETH, Dave (taker) bán 20 ETH
- **Oracle Prices:** Giống Test Case 2

### Đầu ra
- **Block 0:**
  - Vester: 200 triệu tokens
  - Treasury: 0 tokens
- **Block 10:**
  - Vester: ~199,999,974.66 tokens (vesting tiếp tục)
  - Treasury: ~25.34 tokens (vested, nhưng không phân phối rewards)
  - Alice: Số dư ban đầu (không có rewards)
  - Bob: Số dư ban đầu (không có rewards)
  - Carl: Số dư ban đầu (không có rewards)
  - Dave: Số dư ban đầu (không có rewards)

### Tại sao chạy theo cách này?

1. **Zero Multiplier:** Khi FeeMultiplierPpm = 0, không có rewards được phân phối.
2. **Vesting Continues:** Tokens vẫn vest từ vester đến treasury.
3. **No Distribution:** Treasury tích lũy tokens nhưng không phân phối chúng.
4. **Traders Get Nothing:** Mặc dù trading xảy ra, không có rewards được trao.

---

## Tóm tắt Flow

### Trading Rewards Distribution Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. KHỞI TẠO VESTING                                         │
│    - Vester account có số dư ban đầu                         │
│    - Treasury account bắt đầu ở 0                            │
│    - Vest entry định nghĩa vesting schedule                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. VESTING XẢY RA                                           │
│    - Tokens vest từ vester đến treasury mỗi block           │
│    - Vesting rate = total_vest / (end_time - start_time)     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. HOẠT ĐỘNG TRADING                                        │
│    - Users đặt và match orders                              │
│    - Trading fees được thu thập                              │
│    - Taker và maker roles được xác định                     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. TÍNH TOÁN REWARDS                                        │
│    - Net fees = TakerFee - MakerRebate - TakerFeeRevShare    │
│    - Rewards = Net fees × FeeMultiplierPpm                   │
│    - Reward tokens = Rewards (USD) / Rewards token price     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. PHÂN PHỐI REWARDS                                        │
│    - Chỉ tài khoản taker nhận rewards                        │
│    - Rewards được phân phối từ treasury                      │
│    - Indexer events được emit                                │
└─────────────────────────────────────────────────────────────┘
```

### Trạng thái quan trọng

1. **Vesting State:**
   ```
   Vester: 200M → Giảm (tokens vesting ra)
   Treasury: 0 → Tăng (tokens vesting vào)
   ```

2. **Rewards Distribution:**
   ```
   Treasury: Tích lũy vested tokens
   Traders: Nhận rewards dựa trên hoạt động trading
   ```

3. **Rewards Calculation:**
   ```
   Net Fees = TakerFee - MakerRebate - TakerFeeRevShare
   Rewards = Net Fees × FeeMultiplierPpm
   Reward Tokens = Rewards (USD) / Token Price
   ```

### Điểm quan trọng

1. **Vesting:**
   - Tokens vest từ vester account đến treasury account
   - Vesting xảy ra liên tục trong vesting period
   - Vesting rate = total_vest / (end_time - start_time)

2. **Rewards Calculation:**
   - Dựa trên trading fees (taker fees)
   - Net fees = taker fee - maker rebate - taker fee revenue share
   - Rewards = net fees × fee multiplier (PPM)
   - Reward tokens = rewards (USD) / rewards token price

3. **Distribution:**
   - Chỉ tài khoản taker nhận rewards (không phải makers)
   - Rewards được phân phối từ treasury account
   - Nhiều tài khoản có thể nhận rewards trong cùng block
   - Chỉ một taker mỗi block nhận rewards (taker đầu tiên)

4. **Fee Multiplier:**
   - Điều khiển phần trăm net fees trở thành rewards
   - 990,000 PPM = 99% net fees trở thành rewards
   - 0 PPM = không phân phối rewards

5. **Indexer Events:**
   - Trading rewards events được emit cho mỗi phân phối
   - Events bao gồm account address và reward amount
   - Được sử dụng bởi off-chain systems để theo dõi rewards

6. **Oracle Prices:**
   - Rewards token price được sử dụng để chuyển đổi USD rewards sang tokens
   - Trading asset prices được sử dụng để tính trading fees
   - Prices phải có sẵn cho rewards calculation

### Lý do thiết kế

1. **Khuyến khích Trading:** Rewards khuyến khích users trade và cung cấp liquidity.

2. **Phân phối Công bằng:** Rewards dựa trên trading fees đảm bảo active traders nhận nhiều rewards hơn.

3. **Tập trung Taker:** Chỉ takers nhận rewards để khuyến khích market taking và liquidity consumption.

4. **Kiểm soát Vesting:** Cơ chế vesting kiểm soát tốc độ phân phối token theo thời gian.

5. **Linh hoạt:** Fee multiplier cho phép điều chỉnh phần trăm rewards mà không thay đổi vesting schedule.

6. **Minh bạch:** Indexer events cung cấp minh bạch vào phân phối rewards.
