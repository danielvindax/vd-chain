# Tài liệu: Khớp Lệnh và Tính Toán Funding Rate

## Tổng quan

Tài liệu này mô tả chi tiết về:
1. **Cơ chế khớp lệnh (Order Matching)** trong orderbook
2. **Tính toán Funding Rate** và các công thức liên quan
3. **Premium Calculation** dựa trên impact prices
4. **Settlement Calculation** cho funding payments

---

## 1. Khớp Lệnh (Order Matching)

### 1.1. Cơ chế khớp lệnh

Khớp lệnh được thực hiện theo nguyên tắc **Price-Time Priority**:
- **Price Priority**: Lệnh có giá tốt hơn được ưu tiên
- **Time Priority**: Trong cùng mức giá, lệnh đặt trước được ưu tiên

### 1.2. Quy trình khớp lệnh

```
┌─────────────────────────────────────────────────────────────┐
│ 1. TAKER ORDER ĐẾN                                          │
│    - Taker order là lệnh mới cần khớp                        │
│    - Có thể là market order, limit order, hoặc liquidation   │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. TÌM MAKER ORDER TỐT NHẤT                                 │
│    - Tìm best order trên phía đối diện của orderbook        │
│    - Best bid cho taker sell order                           │
│    - Best ask cho taker buy order                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. KIỂM TRA CROSSING                                         │
│    - Taker buy: giá taker >= giá maker                       │
│    - Taker sell: giá taker <= giá maker                      │
│    - Nếu không cross, dừng khớp                              │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. KIỂM TRA ORDER ID VÀ SUBACCOUNT                          │
│    - Nếu taker đang replace maker (cùng OrderId): skip      │
│    - Nếu cùng subaccount: cancel maker, continue            │
│    - Nếu khác subaccount: tiếp tục                           │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. KHỚP LỆNH                                                │
│    - Fill amount = min(taker_remaining, maker_remaining)    │
│    - Giá khớp = giá của maker order (price-time priority)   │
│    - Cập nhật remaining size cho cả hai lệnh                 │
│    - Track OrderHash vào matchedOrderHashToOrder            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 6. KIỂM TRA COLLATERALIZATION                              │
│    - Kiểm tra margin requirements sau khi khớp              │
│    - Nếu maker order fail: remove và tiếp tục               │
│    - Nếu taker order fail: dừng khớp                         │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 7. LẶP LẠI HOẶC DỪNG                                         │
│    - Nếu taker order còn remaining size: lặp lại từ bước 2   │
│    - Nếu taker order đã fill hết: dừng                       │
└─────────────────────────────────────────────────────────────┘
```

### 1.3. Công thức khớp lệnh

#### Fill Amount
```
FillAmount = min(TakerRemainingSize, MakerRemainingSize)
```

#### Match Price
```
MatchPrice = MakerOrderPrice
```
*Lưu ý: Giá khớp luôn là giá của maker order (lệnh đã có trong orderbook)*

#### Quote Quantums (Giá trị khớp)
```
QuoteQuantums = FillAmount × MatchPrice × 10^(-QuantumConversionExponent)
```

### 1.4. Điều kiện khớp lệnh

1. **Opposing Sides**: Taker và maker phải ở hai phía đối diện
   - Taker buy ↔ Maker sell
   - Taker sell ↔ Maker buy

2. **Price Crossing**: Giá phải cross
   - Taker buy: `TakerPrice >= MakerPrice`
   - Taker sell: `TakerPrice <= MakerPrice`

3. **Different Subaccounts**: Không được self-trade
   - `TakerSubaccountId != MakerSubaccountId`

4. **Sufficient Size**: Cả hai lệnh phải có remaining size > 0

### 1.5. OrderHash và Tracking Matched Orders

#### OrderHash là gì?

**OrderHash** là SHA256 hash của order proto bytes, được sử dụng để:
- Định danh duy nhất cho mỗi order
- Track các orders đã được match trong một matching cycle
- Đảm bảo không có duplicate orders trong cùng một block

#### Công thức tính OrderHash

```
OrderHash = SHA256(OrderProtoBytes)
```

Trong đó:
- `OrderProtoBytes`: Serialized bytes của order proto message
- `SHA256`: Hàm băm SHA-256

#### OrderHash trong Matching Process

Trong quá trình khớp lệnh, hệ thống duy trì một map `matchedOrderHashToOrder` để track các orders đã được match:

```go
matchedOrderHashToOrder map[OrderHash]MatchableOrder
```

**Quy trình tracking:**

1. **Khi khớp thành công:**
   ```
   makerOrderHash = makerOrder.GetOrderHash()
   matchedOrderHashToOrder[makerOrderHash] = makerOrder
   
   if !takerOrderHashWasSet {
       takerOrderHash = takerOrder.GetOrderHash()
       matchedOrderHashToOrder[takerOrderHash] = takerOrder
       takerOrderHashWasSet = true
   }
   ```

2. **Mục đích:**
   - Đảm bảo mỗi order chỉ được match một lần trong cùng một matching cycle
   - Track tất cả orders đã tham gia matching để cập nhật state
   - Tránh duplicate processing

#### Điều kiện OrderHash

1. **Unique Identification**: Mỗi order có một OrderHash duy nhất
   - Cùng một order (cùng proto bytes) → cùng OrderHash
   - Khác order → khác OrderHash

2. **Replacement Order Check**: 
   - Nếu taker order đang replace maker order (cùng OrderId), skip maker order
   ```
   if makerOrderId == takerOrderId {
       continue  // Skip, sẽ remove sau khi matching xong
   }
   ```

3. **Self-Trade Prevention**:
   - Kiểm tra subaccount ID trước khi check OrderHash
   - Nếu cùng subaccount → cancel maker order (không match)

#### Lưu ý về OrderHash

- **Liquidation Orders**: Có OrderHash riêng dựa trên `PerpetualLiquidationInfo`
- **Stateful Orders**: OrderHash được tính từ order proto
- **Short-Term Orders**: OrderHash được tính từ order proto

#### Ví dụ

**Scenario**: Taker order match với 3 maker orders

```
Matching Cycle:
1. Match với Maker1 → 
   matchedOrderHashToOrder[Maker1Hash] = Maker1
   matchedOrderHashToOrder[TakerHash] = Taker

2. Match với Maker2 → 
   matchedOrderHashToOrder[Maker2Hash] = Maker2
   // TakerHash đã có, không cần set lại

3. Match với Maker3 → 
   matchedOrderHashToOrder[Maker3Hash] = Maker3
   // TakerHash đã có, không cần set lại

Result: matchedOrderHashToOrder chứa 4 entries:
- TakerHash → Taker Order
- Maker1Hash → Maker1 Order  
- Maker2Hash → Maker2 Order
- Maker3Hash → Maker3 Order
```

---

## 2. Tính Toán Premium

### 2.1. Impact Prices

**Impact Bid Price**: Giá trung bình khi thực hiện một lệnh bán với kích thước `ImpactNotionalQuoteQuantums`.

**Impact Ask Price**: Giá trung bình khi thực hiện một lệnh mua với kích thước `ImpactNotionalQuoteQuantums`.

#### Công thức tính Impact Price

```
ImpactPrice = AverageExecutionPrice(ImpactOrder)
```

Trong đó:
- Impact Order có size = `ImpactNotionalQuoteQuantums` (trong quote quantums)
- Impact Price được tính bằng cách:
  1. Giả lập một lệnh market với size = `ImpactNotionalQuoteQuantums`
  2. Tính giá trung bình khi lệnh này được fill hoàn toàn
  3. Nếu không đủ liquidity, Impact Price = 0 (bid) hoặc ∞ (ask)

### 2.2. Premium Calculation

Premium được tính dựa trên công thức:

```
P = (Max(0, Impact Bid - Index Price) - Max(0, Index Price - Impact Ask)) / Index Price
```

#### Piece-wise Function

Premium có thể được biểu diễn dưới dạng hàm từng phần:

**Case 1: Index < Impact Bid**
```
P = (Impact Bid / Index Price) - 1
```
- Premium > 0 (dương)
- Longs trả shorts

**Case 2: Impact Bid ≤ Index ≤ Impact Ask**
```
P = 0
```
- Premium = 0
- Không có funding

**Case 3: Impact Ask < Index**
```
P = (Impact Ask / Index Price) - 1
```
- Premium < 0 (âm)
- Shorts trả longs

### 2.3. Premium trong Parts-Per-Million (PPM)

Premium được biểu diễn trong PPM (parts-per-million):

```
PremiumPpm = Premium × 1,000,000
```

#### Công thức tính PremiumPpm

```
PremiumPpm = ((ImpactPrice / IndexPrice) - 1) × 1,000,000
```

Sau đó được clamp trong khoảng `[-MaxAbsPremiumVotePpm, MaxAbsPremiumVotePpm]`.

---

## 3. Tính Toán Funding Rate

### 3.1. Funding Rate Components

Funding Rate được tính từ hai thành phần:

```
Funding Rate = Premium + Default Funding
```

Trong đó:
- **Premium**: Được tính từ premium samples (xem phần 2)
- **Default Funding**: Giá trị mặc định được cấu hình cho perpetual

### 3.2. Premium Samples Processing

#### Thu thập Premium Samples

Trong mỗi **Funding Sample Epoch**:
- Premium được tính và lưu lại như một sample
- Số lượng samples tối thiểu = `FundingTickDuration / FundingSampleDuration`

#### Xử lý Premium Samples

1. **Remove Tail Samples**: Loại bỏ một phần trăm samples ở đầu và cuối sau khi sort
   ```
   TailRemovalRate = RemovedTailSampleRatioPpm / 1,000,000
   ```

2. **Average Remaining Samples**: Tính trung bình các samples còn lại
   ```
   PremiumPpm = Average(RemainingSamples)
   ```

### 3.3. Funding Rate Clamping

Funding rate bị giới hạn để ngăn chặn thanh toán quá mức:

```
|Funding Rate| ≤ Clamp Factor × (Initial Margin - Maintenance Margin)
```

#### Công thức chi tiết

```
MaintenanceMarginPpm = InitialMarginPpm × MaintenanceFractionPpm / 1,000,000
MaxAbsFundingClampPpm = FundingRateClampFactorPpm × (InitialMarginPpm - MaintenanceMarginPpm) / 1,000,000
```

Trong đó:
- `FundingRateClampFactorPpm`: Hệ số clamp (thường là 600% = 6,000,000 ppm)
- `InitialMarginPpm`: Initial margin requirement (ví dụ: 5% = 50,000 ppm)
- `MaintenanceFractionPpm`: Maintenance fraction (ví dụ: 60% = 600,000 ppm của initial margin)
- `MaintenanceMarginPpm`: Maintenance margin = InitialMargin × MaintenanceFraction / 1,000,000

#### Ví dụ

**Input:**
- Initial Margin: 5% (50,000 ppm)
- Maintenance Fraction: 60% (600,000 ppm)
- Clamp Factor: 600% (6,000,000 ppm)

**Calculation:**
```
MaintenanceMarginPpm = 50,000 × 600,000 / 1,000,000 = 30,000 ppm = 3%
MaxClamp = 6,000,000 × (50,000 - 30,000) / 1,000,000
         = 6,000,000 × 20,000 / 1,000,000
         = 120,000 ppm = 12%
```

Nếu funding rate tính được là -17.6% (-176,000 ppm), nó sẽ bị clamp về -12% (-120,000 ppm).

---

## 4. Funding Index Calculation

### 4.1. Funding Index Delta

Funding Index Delta được tính từ funding rate, thời gian, và giá oracle:

```
IndexDelta = FundingRatePpm × (Time / RealizationPeriod) × QuoteQuantumsPerBaseQuantum
```

#### Công thức chi tiết

```
IndexDelta = FundingRatePpm × TimeSinceLastFunding × MarketPrice × 10^(QuoteAtomicResolution) / (8 hours × 10^6 × 10^BaseAtomicResolution)
```

Trong đó:
- `FundingRatePpm`: Funding rate trong parts-per-million (8-hour rate)
- `TimeSinceLastFunding`: Thời gian tính bằng giây kể từ funding tick cuối
- `RealizationPeriod`: Chu kỳ funding (8 hours = 28,800 giây)
- `MarketPrice`: Giá oracle của perpetual
- `QuoteAtomicResolution`: Atomic resolution của quote currency
- `BaseAtomicResolution`: Atomic resolution của base currency

#### Implementation

```go
// 1. Nhân với time delta
result = TimeSinceLastFunding × FundingRatePpm

// 2. Nhân với giá (chuyển từ base quantums sang quote quantums)
result = BaseToQuoteQuantums(result, BaseAtomicResolution, MarketPrice, MarketExponent)

// 3. Chia cho realization period (8 hours)
result = result / (60 × 60 × 8)
```

### 4.2. Funding Index Update

Funding Index được cập nhật tại mỗi funding tick epoch:

```
NewFundingIndex = CurrentFundingIndex + IndexDelta
```

Nếu funding rate = 0, funding index không thay đổi.

---

## 5. Settlement Calculation

### 5.1. Settlement Formula

Settlement được tính khi subaccount nhận transfer hoặc khi position thay đổi:

```
Settlement = -(FundingIndex - PositionFundingIndex) × PositionSize
```

Trong đó:
- `FundingIndex`: Funding index hiện tại của perpetual
- `PositionFundingIndex`: Funding index khi position được mở/cập nhật lần cuối
- `PositionSize`: Kích thước position (dương cho long, âm cho short)

#### Lưu ý về dấu

- **Long Position** (PositionSize > 0):
  - Nếu FundingIndex tăng → Settlement < 0 → Long trả tiền
  - Nếu FundingIndex giảm → Settlement > 0 → Long nhận tiền

- **Short Position** (PositionSize < 0):
  - Nếu FundingIndex tăng → Settlement > 0 → Short nhận tiền
  - Nếu FundingIndex giảm → Settlement < 0 → Short trả tiền

### 5.2. Settlement trong Quote Quantums

Settlement được tính trong quote quantums (PPM):

```
SettlementPpm = -(FundingIndex - PositionFundingIndex) × PositionSize
```

Sau đó chuyển sang quote quantums:

```
SettlementQuoteQuantums = SettlementPpm / 1,000,000
```

### 5.3. Funding Payment

Funding Payment là giá trị ngược lại của Settlement:

```
FundingPayment = -Settlement
```

- Positive settlement → Negative funding payment (position nhận tiền)
- Negative settlement → Positive funding payment (position trả tiền)

---

## 6. Ví dụ Tính Toán

### 6.1. Ví dụ: Premium Calculation

**Input:**
- Impact Bid: $28,000
- Impact Ask: $28,005
- Index Price: $27,960

**Calculation:**
```
Index < Impact Bid → Premium = (28,000 / 27,960) - 1 = 0.00143 = 0.143%
PremiumPpm = 0.00143 × 1,000,000 = 1,430 ppm
```

**Result:** Premium = 1,430 ppm (dương) → Longs trả shorts

---

### 6.2. Ví dụ: Funding Rate Calculation

**Input:**
- Premium: 1,430 ppm
- Default Funding: 0 ppm
- Initial Margin: 5% (50,000 ppm)
- Maintenance Margin: 3% (30,000 ppm)
- Clamp Factor: 600% (6,000,000 ppm)

**Calculation:**
```
Funding Rate = 1,430 + 0 = 1,430 ppm

Max Clamp = 6,000,000 × (50,000 - 30,000) / 1,000,000
          = 6,000,000 × 20,000 / 1,000,000
          = 120,000 ppm = 12%

|1,430| < 120,000 → Không cần clamp
```

**Result:** Funding Rate = 1,430 ppm (0.143%)

---

### 6.3. Ví dụ: Funding Index Delta

**Input:**
- Funding Rate: 1,430 ppm (8-hour rate)
- Time Since Last Funding: 28,800 giây (8 hours)
- Market Price: $28,000
- Base Atomic Resolution: -8 (BTC)
- Quote Atomic Resolution: -6 (USDC)

**Calculation:**
```
IndexDelta = 1,430 × 28,800 × 28,000 × 10^(-6) / (28,800 × 10^8)
          = 1,430 × 28,000 × 10^(-6) / 10^8
          = 1,430 × 28,000 / 10^14
          = 40,040,000 / 10^14
          = 0.0000004004 (trong quote quantums per base quantum)
```

**Result:** IndexDelta ≈ 400 ppm (trong quote quantums per base quantum)

---

### 6.4. Ví dụ: Settlement Calculation

**Input:**
- Current Funding Index: 1,000 ppm
- Position Funding Index: 0 ppm (mở position khi index = 0)
- Position Size: 0.8 BTC (long position)

**Calculation:**
```
SettlementPpm = -(1,000 - 0) × 0.8 × 10^8
              = -1,000 × 80,000,000
              = -80,000,000,000 ppm

SettlementQuoteQuantums = -80,000,000,000 / 1,000,000
                        = -80,000 quote quantums
                        = -$0.08 (vì USDC có resolution -6)
```

**Result:** Settlement = -$0.08 → Long trả $0.08 cho funding

---

## 7. Tóm tắt Công thức

### 7.1. Order Matching

```
FillAmount = min(TakerRemaining, MakerRemaining)
MatchPrice = MakerOrderPrice
QuoteQuantums = FillAmount × MatchPrice × 10^(-QuantumConversionExponent)
```

### 7.2. Premium

```
If Index < Impact Bid:  P = (Impact Bid / Index) - 1
If Impact Bid ≤ Index ≤ Impact Ask:  P = 0
If Impact Ask < Index:  P = (Impact Ask / Index) - 1

PremiumPpm = P × 1,000,000
```

### 7.3. Funding Rate

```
FundingRate = Premium + DefaultFunding
MaxClamp = ClampFactor × (InitialMargin - MaintenanceMargin)
FinalFundingRate = Clamp(FundingRate, -MaxClamp, MaxClamp)
```

### 7.4. Funding Index

```
IndexDelta = FundingRatePpm × Time × MarketPrice × 10^QuoteRes / (8h × 10^6 × 10^BaseRes)
NewIndex = CurrentIndex + IndexDelta
```

### 7.5. Settlement

```
SettlementPpm = -(FundingIndex - PositionFundingIndex) × PositionSize
SettlementQuoteQuantums = SettlementPpm / 1,000,000
FundingPayment = -Settlement
```

---

## 8. Điểm quan trọng

### 8.1. Units và Resolutions

- **PPM (Parts-Per-Million)**: 1,000,000 ppm = 100%
- **Atomic Resolution**: Số chữ số thập phân
  - BTC: -8 (1 BTC = 10^8 base quantums)
  - USDC: -6 (1 USDC = 10^6 quote quantums)

### 8.2. Funding Rate Period

- Funding rate được tính cho **8-hour period**
- Funding tick epoch thường là 8 hours
- Funding sample epoch thường là vài phút (để thu thập samples)

### 8.3. Precision

- Tất cả phép nhân được thực hiện trước phép chia để tránh mất precision
- Sử dụng `big.Int` và `big.Rat` cho tính toán chính xác
- Rounding hướng về zero (truncated division)

### 8.4. Edge Cases

1. **Empty Orderbook**: Premium = 0
2. **Insufficient Liquidity**: Impact Price = 0 hoặc ∞ → Premium = 0
3. **Zero Funding Rate**: Funding Index không thay đổi
4. **No Position Change**: Settlement = 0 nếu FundingIndex không đổi

---

## 9. References

- Code Implementation:
  - Order Matching: `protocol/x/clob/memclob/memclob.go`
  - Premium Calculation: `protocol/x/clob/memclob/memclob.go::GetPricePremium()`
  - Funding Rate: `protocol/x/perpetuals/keeper/perpetual.go::MaybeProcessNewFundingTickEpoch()`
  - Funding Index: `protocol/x/perpetuals/funding/funding.go::GetFundingIndexDelta()`
  - Settlement: `protocol/x/perpetuals/lib/lib.go::GetSettlementPpmWithPerpetual()`

- Test Documentation:
  - `protocol/docs/e2e/funding/funding_e2e_test_docs.md`

