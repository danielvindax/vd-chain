# Vindax Chain – Rebrand Notes

This repository is a fork of the dYdX v4 chain.  
It is being rebranded into **Vindax Chain (vd-chain)** with a new identity, chain ID, denom, and module path.

## ✅ New naming and network identity
| Item | Value |
|------|-------|
| Project name | **Vindax Chain** |
| Binary | `vindaxd` |
| Chain ID (mainnet) | `vindax-1` |
| Chain ID (dev/local) | `localvindax` | `vindax-testnet` |
| Bech32 prefix | `vindax` |
| Base denom | `avdtn` |
| Display denom | `vdtn` |

## ✅ Go module rename
`go.mod` updated to: module github.com/danielvindax/vd-chain
All internal imports will progressively be updated to use the new module path instead of the previous dYdX paths.

## ✅ Status
- [x] Fork from dYdX v4