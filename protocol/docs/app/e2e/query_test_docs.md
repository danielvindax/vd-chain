# Test Documentation: App Query E2E Tests

## Overview

This test file verifies **Parallel Query** functionality in the application. The test ensures that the application can handle concurrent queries safely without data races. The test uses Go's race detector to verify thread safety.

---

## Test Function: TestParallelQuery

### Test Case: Success - Parallel Queries Without Data Races

### Input
- **Concurrent Operations:**
  - Thread 1: Advance blocks (blocks 2-49)
  - Thread 2: Query app/version repeatedly
  - Thread 3: Query store/blocktime/key directly
  - Thread 4: Query gRPC PreviousBlockInfo
- **Synchronization:** Atomic boolean to coordinate threads
- **Execution:** All operations run concurrently until block limit reached

### Output
- **No Data Races:** Test passes with `-race` flag enabled
- **Query Results:** All queries return valid results
- **Height Monotonicity:** Query heights are monotonically increasing
- **Consistency:** Store queries and gRPC queries return consistent data

### Why It Runs This Way?

1. **Concurrency Testing:** Tests that the application handles concurrent queries correctly.
2. **Race Detection:** Uses Go's race detector to find data races.
3. **Multiple Query Types:** Tests different query paths (app, store, gRPC).
4. **Stress Test:** Concurrent queries while blocks advance stress tests the system.
5. **Atomic Coordination:** Uses atomic boolean to maximize potential for data races.

---

## Flow Summary

### Parallel Query Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. START CONCURRENT THREADS                                  │
│    - Block advancement thread                                │
│    - App/version query thread                                │
│    - Store query thread                                       │
│    - gRPC query thread                                       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CONCURRENT EXECUTION                                      │
│    - Blocks advance while queries execute                    │
│    - No synchronization between threads                      │
│    - Atomic boolean coordinates completion                   │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. QUERY EXECUTION                                           │
│    - App/version: Query application version                  │
│    - Store: Query blocktime store directly                    │
│    - gRPC: Query PreviousBlockInfo via gRPC                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. VERIFICATION                                              │
│    - All queries return valid results                        │
│    - Heights are monotonically increasing                    │
│    - No data races detected                                  │
└─────────────────────────────────────────────────────────────┘
```

### Important States

1. **Query States:**
   ```
   Query Request → Execute Query → Return Result → Verify Result
   ```

2. **Concurrent Execution:**
   ```
   Thread 1: Block Advancement
   Thread 2: App/Version Queries
   Thread 3: Store Queries
   Thread 4: gRPC Queries
   ```

### Key Points

1. **Query Types:**
   - App/Version: Application version query
   - Store: Direct store query (blocktime key)
   - gRPC: gRPC service query (PreviousBlockInfo)

2. **Concurrency:**
   - Multiple threads execute queries concurrently
   - Blocks advance while queries execute
   - No synchronization between query threads

3. **Race Detection:**
   - Uses Go's `-race` flag to detect data races
   - Atomic boolean maximizes potential for races
   - Wait group coordinates thread completion

4. **Verification:**
   - Heights must be monotonically increasing
   - Store and gRPC queries must return consistent data
   - All queries must return valid results

### Design Rationale

1. **Thread Safety:** Application must be thread-safe for concurrent queries.

2. **Race Detection:** Go's race detector helps find data races during testing.

3. **Stress Testing:** Concurrent queries while blocks advance stress tests the system.

4. **Multiple Paths:** Tests different query paths to ensure all are thread-safe.

