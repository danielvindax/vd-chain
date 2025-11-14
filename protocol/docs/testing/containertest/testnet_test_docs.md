# Test Documentation: Container Tests

## Overview

This test file verifies the **Container Test Framework** functionality. Container tests run a full network of nodes in Docker containers, allowing end-to-end testing of the blockchain. The test framework provides interfaces to interact with the chain through CometBFT and gRPC clients. The framework also runs an HTTP server for exchange price feeds.

---

## Test Function: TestPlaceOrder

### Test Case: Success - Place Order on Network

### Input
- **Network:** Full testnet with multiple nodes running in Docker containers
- **Node:** Alice node
- **Order:**
  - SubaccountId: Alice_Num0
  - ClobPairId: 0
  - Side: BUY
  - Quantums: 10,000,000
  - Subticks: 1,000,000
  - GoodTilBlock: 20

### Output
- **Transaction:** Successfully broadcast and included in block
- **Order:** Order placed on order book

### Why It Runs This Way?

1. **Full Network Test:** Tests order placement in a real network environment with multiple nodes.
2. **Docker Containers:** Each node runs in a separate Docker container, simulating a real network.
3. **End-to-End:** Verifies the complete flow from transaction submission to order placement.
4. **Network Interaction:** Tests interaction with the network through gRPC and CometBFT clients.

---

## Test Function: TestBankSend

### Test Case: Success - Send Tokens Between Accounts

### Input
- **Network:** Full testnet with multiple nodes
- **Node:** Alice node
- **Initial State:**
  - Alice has initial balance
  - Bob has initial balance
- **Transaction:**
  - From: Bob
  - To: Alice
  - Amount: 1 USDC

### Output
- **Initial Balances:** Verified against expected values
  - Alice initial balance matches expected
  - Bob initial balance matches expected
- **Final Balances:** Verified after transaction
  - Alice balance increased by 1 USDC
  - Bob balance decreased by 1 USDC

### Why It Runs This Way?

1. **Balance Verification:** Tests that balances are correctly tracked and updated.
2. **Expected Output:** Uses expect files to verify exact balance values.
3. **Network Consensus:** Verifies that transactions are properly propagated and included in blocks.
4. **State Consistency:** Ensures all nodes have consistent state after transaction.

---

## Test Function: TestMarketPrices

### Test Case: Success - Market Prices Update from Exchange

### Input
- **Network:** Full testnet with multiple nodes
- **Exchange Prices Set Before Start:**
  - BTC-USD: 50,001
  - ETH-USD: 55,002
  - LINK-USD: 55,003
- **Node:** Alice node
- **Timeout:** 30 seconds

### Output
- **Market Prices:** Prices updated to match exchange prices
  - BTC-USD price matches expected
  - ETH-USD price matches expected
  - LINK-USD price matches expected

### Why It Runs This Way?

1. **Price Feed Integration:** Tests integration with external exchange price feeds.
2. **HTTP Server:** Framework runs HTTP server that provides exchange prices.
3. **Oracle Updates:** Verifies that oracle prices are updated from exchange feeds.
4. **Moving Window:** Prices use a moving window, so prices should be set before network start.
5. **Polling:** Uses polling with timeout to wait for prices to update.

---

## Test Function: TestUpgrade

### Test Case: Success - Upgrade Network to New Version

### Input
- **Network:** Testnet with pre-upgrade genesis
- **Node:** Alice node
- **Upgrade:** Upgrade to current version
- **Upgrader:** Alice account

### Output
- **Upgrade:** Successfully executed
- **Network:** Running on new version

### Why It Runs This Way?

1. **Upgrade Testing:** Tests the network upgrade mechanism.
2. **Pre-upgrade Genesis:** Uses special genesis state for pre-upgrade testing.
3. **Version Management:** Verifies that network can upgrade to new version.
4. **State Migration:** Ensures state is correctly migrated during upgrade.

---

## Flow Summary

### Container Test Framework Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. CREATE TESTNET                                            │
│    - Initialize testnet with nodes                           │
│    - Configure Docker containers                             │
│    - Set up HTTP server for price feeds                      │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CONFIGURE BEFORE START                                    │
│    - Set exchange prices (if needed)                        │
│    - Configure genesis state                                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. START NETWORK                                             │
│    - Start Docker containers                                 │
│    - Wait for nodes to sync                                  │
│    - Verify network is ready                                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. INTERACT WITH NETWORK                                     │
│    - Query state (balances, prices, etc.)                    │
│    - Broadcast transactions                                 │
│    - Wait for blocks                                         │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. VERIFY RESULTS                                            │
│    - Compare against expected output                         │
│    - Verify state changes                                    │
│    - Clean up containers                                     │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Network State:**
   ```
   Not Started → Starting → Running → Cleaned Up
   ```

2. **Price Updates:**
   ```
   Exchange Price → HTTP Server → Oracle → Market Price
   ```

3. **Transaction Flow:**
   ```
   Broadcast → Mempool → Block → State Update
   ```

### Key Points

1. **Docker Containers:**
   - Each node runs in a separate Docker container
   - Simulates a real network environment
   - Allows testing network consensus and propagation

2. **Price Feed Integration:**
   - HTTP server provides exchange prices
   - Prices should be set before network start
   - Oracle uses moving window for price updates

3. **Expected Output:**
   - Tests use expect files to verify exact output
   - Use `-accept` flag to update expect files
   - Ensures deterministic test results

4. **Network Interaction:**
   - Query: Read state from nodes
   - BroadcastTx: Submit transactions to network
   - Wait: Wait for blocks to be produced

5. **Cleanup:**
   - Always clean up containers after tests
   - Use `defer testnet.MustCleanUp()`
   - Prevents resource leaks

6. **Upgrade Testing:**
   - Tests network upgrade mechanism
   - Uses pre-upgrade genesis state
   - Verifies state migration

### Design Rationale

1. **End-to-End Testing:** Container tests provide full network testing, not just unit tests.

2. **Real Environment:** Docker containers simulate real network conditions.

3. **Integration Testing:** Tests integration between components (nodes, price feeds, etc.).

4. **Deterministic:** Expect files ensure tests are deterministic and reproducible.

5. **Flexibility:** Framework allows testing various scenarios (upgrades, price updates, etc.).

