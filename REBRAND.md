# Vindax Chain – Rebrand Notes

This repository is a fork of the dYdX v4 chain.  
It is being rebranded into **Vindax Chain (vd-chain)** with a new identity, chain ID, denom, and module path.

## ✅ New naming and network identity
| Item | Value |
|------|-------|
| Project name | **Vindax Chain** |
| Binary | `vdxd` |
| Chain ID (mainnet) | `vindax-1` |
| Chain ID (dev/local) | `vindax-local-1` |
| Bech32 prefix | `vindax` |
| Base denom | `uvdx` |
| Display denom | `VDX` |

## ✅ Go module rename
`go.mod` updated to: module github.com/danielvindax/vd-chain
All internal imports will progressively be updated to use the new module path instead of the previous dYdX paths.

## ✅ Status
- [x] Fork from dYdX v4
- [ ] Updated go.mod
- [ ] Updated imports to new module
- [ ] Rebrand Bech32 prefix → `vindax`
- [ ] Base denom → `uvdx`
- [ ] Rename binary → `vdxprotocold`    
- [ ] Update README and docs
- [ ] Generate new local genesis
- [ ] Build CI / Docker

## ✅ Local devnet quick start (after rebrand)
```bash
make build
./build/vdxprotocold init local --chain-id vindax-local-1
./build/vdxprotocold keys add alice --keyring-backend test
./build/vdxprotocold add-genesis-account $(./build/vdxprotocold keys show alice -a --keyring-backend test) 100000000uvdx
./build/vdxprotocold gentx alice 1000000uvdx --chain-id vindax-local-1 --keyring-backend test
./build/vdxprotocold collect-gentxs
./build/vdxprotocold start
