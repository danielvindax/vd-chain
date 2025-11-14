# Test Documentation: Short-Term Orders E2E Tests

## Overview

This test file verifies **Short-Term Order** functionality in the CLOB (Central Limit Order Book) module. Short-term orders are orders that expire at the end of the current block. The test ensures that:
1. Orders can be placed on the order book
2. Orders can be matched with existing orders
3. Off-chain and on-chain messages are correctly emitted
4. Order state is correctly tracked

---

## Test Function: TestPlaceOrder

### Test Case 1: Success - Place Order on Order Book

### Input
- **Order:**
  - SubaccountId: Alice_Num0
  - ClobPairId: 0
  - Side: BUY
  - Quantums: 5
  - Subticks: 1,000,000 (price = 10)
  - GoodTilBlock: 20

### Output
- **CheckTx:** SUCCESS
- **Off-chain Messages:**
  - Order place message
  - Order update message (with fill amount = 0)
- **On-chain Messages:**
  - Indexer block event in next block

### Why It Runs This Way?

1. **Order Placement:** Order is successfully placed on the order book.
2. **Off-chain Updates:** Indexer needs to be notified about new orders for off-chain systems.
3. **On-chain Events:** Block events are emitted in the next block for on-chain tracking.

---

### Test Case 2: Success - Match Order Fully

### Input
- **Orders:**
  - Alice: Buy 5 at price 10
  - Bob: Sell 5 at price 10

### Output
- **CheckTx:** Both orders SUCCESS
- **Orders Filled:** Both orders fully filled
- **Off-chain Messages:**
  - Order place messages for both orders
  - Order update messages showing full fill
- **On-chain Messages:**
  - Subaccount update events for both Alice and Bob
  - Positions updated correctly

### Why It Runs This Way?

1. **Order Matching:** When buy and sell orders match in price and quantity, they are fully filled.
2. **Position Updates:** Both subaccounts' positions are updated after matching.
3. **Event Emission:** Events are emitted for both maker and taker.

---

## Flow Summary

### Place Order Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. SUBMIT ORDER                                             │
│    - Create MsgPlaceOrder                                   │
│    - Sign transaction                                        │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CHECKTX VALIDATION                                       │
│    - Validate order format                                  │
│    - Check collateralization                                │
│    - Verify order parameters                                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. ORDER PLACEMENT                                          │
│    - Add order to order book                                 │
│    - Emit off-chain order place message                      │
│    - Emit off-chain order update message                     │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. NEXT BLOCK                                               │
│    - Emit on-chain indexer block event                       │
│    - Update order state                                      │
└─────────────────────────────────────────────────────────────┘
```

### Match Order Process

```
┌─────────────────────────────────────────────────────────────┐
│ 1. PLACE FIRST ORDER                                         │
│    - Order added to order book                               │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. PLACE MATCHING ORDER                                      │
│    - Order matches with existing order                       │
│    - Both orders fully filled                                │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. EXECUTE TRADE                                            │
│    - Update subaccount positions                             │
│    - Emit subaccount update events                           │
│    - Remove filled orders from book                          │
└─────────────────────────────────────────────────────────────┘
```

### Key Points

1. **Short-Term Orders:**
   - Expire at the end of the current block
   - GoodTilBlock specifies the block number when order expires
   - Must be matched within the same block or they expire

2. **Order Matching:**
   - Orders match when price and quantity are compatible
   - Maker-taker model: first order is maker, matching order is taker
   - Both orders are removed from book after full match

3. **Off-chain Messages:**
   - Order place message: notifies indexer about new order
   - Order update message: notifies indexer about fill amount
   - Transaction hash included in message headers

4. **On-chain Events:**
   - Subaccount update events: track position changes
   - Indexer block events: track block-level state changes
   - Events emitted in next block after order placement

5. **State Tracking:**
   - Order state tracked in keeper
   - Fill amounts tracked per order
   - Positions updated after matching

### Design Rationale

1. **Efficiency:** Short-term orders provide fast execution for immediate trading needs.

2. **Liquidity:** Encourages active trading by requiring immediate matching.

3. **State Management:** Orders expire automatically, preventing stale orders from cluttering the book.

4. **Event Emission:** Both off-chain and on-chain events ensure complete state tracking.

