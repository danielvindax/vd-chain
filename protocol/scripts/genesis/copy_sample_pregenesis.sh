#!/bin/bash

./scripts/genesis/prod_pregenesis.sh vindaxd
cp /tmp/prod-chain/.vindax/config/sorted_genesis.json ./scripts/genesis/sample_pregenesis.json
