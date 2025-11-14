# Test Documentation: Trading Rewards E2E Tests

## Overview

This test file verifies the **Trading Rewards** distribution mechanism. Trading rewards are distributed to traders based on their trading activity (fees paid). The test ensures that:
1. Rewards are calculated based on trading fees
2. Rewards are distributed from treasury account
3. Only one taker account gets rewards per block
4. Multiple accounts can receive rewards in the same block
5. Rewards multiplier affects distribution amount
6. Vesting tokens are transferred from vester to treasury

---

## Test Function: TestTradingRewards

### Test Case 1: Every Block, Only One Taker Account Gets Rewards

### Input
- **Vest Entry:**
  - VesterAccount: rewards_vester
  - TreasuryAccount: rewards_treasury
  - StartTime: Oct 01 2023 04:00:00
  - EndTime: Oct 01 2028 04:00:00
- **Rewards Params:**
  - FeeMultiplierPpm: 990_000 (99%)
  - MarketId: 30 (rewards token)
- **Orders:**
  - Block 2: Bob (maker) sells, Alice (taker) buys 1 BTC at 28,003
  - Block 13: Alice (maker) buys, Bob (taker) sells 1 BTC at 28,003
- **Oracle Prices:**
  - Rewards token: $1.95
  - BTC: $28,003

### Output
- **Block 0 (Vest Start):**
  - Vester: 200 million tokens
  - Treasury: 0 tokens
- **Block 1:**
  - Vester: ~199,999,997.47 tokens (vesting started)
  - Treasury: ~2.53 tokens (vested)
- **Block 2:**
  - Vester: ~199,999,994.93 tokens
  - Treasury: ~3.07 tokens (after distributing ~1.99 to Alice)
  - Alice: Starting balance + ~1.99 tokens (rewards)
  - Bob: Starting balance (no rewards - was maker)
- **Block 13:**
  - Vester: ~199,999,967.05 tokens
  - Treasury: ~28.96 tokens
  - Bob: Starting balance + ~1.99 tokens (rewards - was taker)
  - Alice: Starting balance + ~1.99 tokens (from block 2)

### Why It Runs This Way?

1. **Vesting:** Tokens vest from vester account to treasury account over time.
2. **Rewards Calculation:** Rewards = (TakerFee - MakerRebate - TakerFeeRevShare) × FeeMultiplierPpm
   - For 1 BTC at $28,003: ($14.0015 - $3.08033 - $7.00075) × 0.99 = $3.8812158
   - Reward tokens = $3.8812158 / $1.95 = ~1.99 tokens
3. **Taker Only:** Only the taker account receives rewards, not the maker.
4. **One Per Block:** Only one taker account gets rewards per block (the first taker in that block).

---

### Test Case 2: Multiple Accounts Get Rewards

### Input
- **Vest Entry:** Same as Test Case 1
- **Rewards Params:** Same as Test Case 1
- **Orders (Block 10):**
  - BTC: Bob (maker) sells 2 BTC, Alice (taker) buys 2 BTC
  - BTC: Alice (maker) buys 2 BTC, Bob (taker) sells 2 BTC
  - ETH: Carl (maker) buys 20 ETH, Dave (taker) sells 20 ETH
  - ETH: Dave (maker) sells 20 ETH, Carl (taker) buys 20 ETH
- **Oracle Prices:**
  - Rewards token: $1.95
  - BTC: $28,003
  - ETH: $1,605

### Output
- **Block 0:**
  - Vester: 200 million tokens
  - Treasury: 0 tokens
- **Block 10:**
  - Vester: ~199,999,974.66 tokens
  - Treasury: ~12.82 tokens (after distributing rewards)
  - Alice: Starting balance + ~3.98 tokens (rewards from BTC trading)
  - Bob: Starting balance + ~3.98 tokens (rewards from BTC trading)
  - Carl: Starting balance + ~2.28 tokens (rewards from ETH trading)
  - Dave: Starting balance + ~2.28 tokens (rewards from ETH trading)

### Why It Runs This Way?

1. **Multiple Takers:** Multiple accounts can receive rewards in the same block.
2. **Rewards Per Trade:** Each taker receives rewards based on their trading fees.
3. **Different Markets:** Rewards are calculated separately for each market (BTC, ETH).
4. **Total Distribution:** Total rewards distributed = sum of individual rewards.

---

### Test Case 3: Rewards Fee Multiplier = 0, No Rewards Distributed

### Input
- **Vest Entry:** Same as Test Case 1
- **Rewards Params:**
  - FeeMultiplierPpm: 0 (0% - no rewards)
- **Orders (Block 10):**
  - BTC: Bob (maker) sells 2 BTC, Alice (taker) buys 2 BTC
  - ETH: Carl (maker) buys 20 ETH, Dave (taker) sells 20 ETH
- **Oracle Prices:** Same as Test Case 2

### Output
- **Block 0:**
  - Vester: 200 million tokens
  - Treasury: 0 tokens
- **Block 10:**
  - Vester: ~199,999,974.66 tokens (vesting continues)
  - Treasury: ~25.34 tokens (vested, but no rewards distributed)
  - Alice: Starting balance (no rewards)
  - Bob: Starting balance (no rewards)
  - Carl: Starting balance (no rewards)
  - Dave: Starting balance (no rewards)

### Why It Runs This Way?

1. **Zero Multiplier:** When FeeMultiplierPpm = 0, no rewards are distributed.
2. **Vesting Continues:** Tokens still vest from vester to treasury.
3. **No Distribution:** Treasury accumulates tokens but doesn't distribute them.
4. **Traders Get Nothing:** Even though trading occurs, no rewards are given.

---

## Flow Summary

### Trading Rewards Distribution Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. INITIALIZE VESTING                                        │
│    - Vester account has initial balance                      │
│    - Treasury account starts at 0                           │
│    - Vest entry defines vesting schedule                     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. VESTING OCCURS                                            │
│    - Tokens vest from vester to treasury each block          │
│    - Vesting rate = total_vest / (end_time - start_time)     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. TRADING ACTIVITY                                          │
│    - Users place and match orders                            │
│    - Trading fees are collected                              │
│    - Taker and maker roles are identified                    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. CALCULATE REWARDS                                         │
│    - Net fees = TakerFee - MakerRebate - TakerFeeRevShare    │
│    - Rewards = Net fees × FeeMultiplierPpm                   │
│    - Reward tokens = Rewards (USD) / Rewards token price     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. DISTRIBUTE REWARDS                                        │
│    - Only taker accounts receive rewards                     │
│    - Rewards distributed from treasury                       │
│    - Indexer events emitted                                  │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Vesting State:**
   ```
   Vester: 200M → Decreasing (tokens vesting out)
   Treasury: 0 → Increasing (tokens vesting in)
   ```

2. **Rewards Distribution:**
   ```
   Treasury: Accumulates vested tokens
   Traders: Receive rewards based on trading activity
   ```

3. **Rewards Calculation:**
   ```
   Net Fees = TakerFee - MakerRebate - TakerFeeRevShare
   Rewards = Net Fees × FeeMultiplierPpm
   Reward Tokens = Rewards (USD) / Token Price
   ```

### Key Points

1. **Vesting:**
   - Tokens vest from vester account to treasury account
   - Vesting occurs continuously over vesting period
   - Vesting rate = total_vest / (end_time - start_time)

2. **Rewards Calculation:**
   - Based on trading fees (taker fees)
   - Net fees = taker fee - maker rebate - taker fee revenue share
   - Rewards = net fees × fee multiplier (PPM)
   - Reward tokens = rewards (USD) / rewards token price

3. **Distribution:**
   - Only taker accounts receive rewards (not makers)
   - Rewards distributed from treasury account
   - Multiple accounts can receive rewards in same block
   - Only one taker per block gets rewards (first taker)

4. **Fee Multiplier:**
   - Controls what percentage of net fees become rewards
   - 990,000 PPM = 99% of net fees become rewards
   - 0 PPM = no rewards distributed

5. **Indexer Events:**
   - Trading rewards events emitted for each distribution
   - Events include account address and reward amount
   - Used by off-chain systems to track rewards

6. **Oracle Prices:**
   - Rewards token price used to convert USD rewards to tokens
   - Trading asset prices used to calculate trading fees
   - Prices must be available for rewards calculation

### Design Rationale

1. **Incentivize Trading:** Rewards incentivize users to trade and provide liquidity.

2. **Fair Distribution:** Rewards based on trading fees ensure active traders receive more rewards.

3. **Taker Focus:** Only takers receive rewards to incentivize market taking and liquidity consumption.

4. **Vesting Control:** Vesting mechanism controls token distribution rate over time.

5. **Flexibility:** Fee multiplier allows adjusting rewards percentage without changing vesting schedule.

6. **Transparency:** Indexer events provide transparency into rewards distribution.

