#!/bin/bash
set -eo pipefail

# Below is the shared genesis configuration for local development, as well as our testnets.
# Any changes to genesis state in those environments belong here.
# If you are making a change to genesis which is _required_ for the chain to function,
# then that change probably belongs in `DefaultGenesis` for the module, and not here.

NINE_ZEROS="000000000"
EIGHTEEN_ZEROS="$NINE_ZEROS$NINE_ZEROS"

# Address of the `subaccounts` module account.
# Obtained from `authtypes.NewModuleAddress(subaccounttypes.ModuleName)`.
SUBACCOUNTS_MODACC_ADDR="vindax1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5vcx39"
REWARDS_VESTER_ACCOUNT_ADDR="vindax1ltyc6y4skclzafvpznpt2qjwmfwgsndp4y7tj7"
# Address of the `bridge` module account.
# Obtained from `authtypes.NewModuleAddress(bridgetypes.ModuleName)`.
BRIDGE_MODACC_ADDR="vindax1zlefkpe3g0vvm9a4h0jf9000lmqutlh9j7tmen"
USDC_DENOM="ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5"
REWARD_TOKEN="avdtn"
NATIVE_TOKEN="avdtn" # public testnet token
DEFAULT_SUBACCOUNT_QUOTE_BALANCE=100000000000000000
DEFAULT_SUBACCOUNT_QUOTE_BALANCE_FAUCET=900000000000000000
DEFAULT_SUBACCOUNT_QUOTE_BALANCE_VAULT=1000000000
MEGAVAULT_MAIN_VAULT_ACCOUNT_ADDR="vindax18tkxrnrkqc2t0lr3zxr5g6a4hdvqksylxst3uu"
DEFAULT_MEGAVAULT_MAIN_VAULT_QUOTE_BALANCE=0 # 0 USDC
NATIVE_TOKEN_WHOLE_COIN="vdtn"
COIN_NAME="VdX Testnet Token"
# Each testnet validator has 1 million whole coins of native token.
TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE=1000000$EIGHTEEN_ZEROS # 1e24 or 1 million native tokens.
# Each testnet validator self-delegates 500k whole coins of native token.
TESTNET_VALIDATOR_SELF_DELEGATE_AMOUNT=500000$EIGHTEEN_ZEROS # 5e23 or 500k native tokens.
FAUCET_NATIVE_TOKEN_BALANCE=50000000$EIGHTEEN_ZEROS # 5e25 or 50 million native tokens. 
ETH_CHAIN_ID=11155111 # sepolia
# https://sepolia.etherscan.io/address/0xf75012c350e4ad55be2048bd67ce6e03b20de82d
ETH_BRIDGE_ADDRESS="0xf75012c350e4ad55be2048bd67ce6e03b20de82d"
TOTAL_NATIVE_TOKEN_SUPPLY=1000000000$EIGHTEEN_ZEROS # 1e27
BRIDGE_GENESIS_ACKNOWLEDGED_NEXT_ID=5
BRIDGE_GENESIS_ACKNOWLEDGED_ETH_BLOCK_HEIGHT=4322136

# Use a fix genesis time from the past.
GENESIS_TIME="2023-01-01T00:00:00Z"

function edit_genesis() {
	GENESIS=$1/genesis.json

	# IFS stands for "Internal Field Separator" and it's a special var that determines how bash splits strings.
	# So IFS=' ' specifies that we want to split on spaces.
	# "read" is a built in bash command that reads from stdin.
	# The -a flag specifies that the input should be treated as an array and assign it to the var specified after.
	# The -r flag tells the command to not treat a Backslash as an escape character.
	IFS=' ' read -ra INPUT_TEST_ACCOUNTS <<<"${2}"
	IFS=' ' read -ra INPUT_FAUCET_ACCOUNTS <<<"${3}"
	IFS=' ' read -ra INPUT_VAULT_ACCOUNTS <<<"${4}"
	IFS=' ' read -ra INPUT_VAULT_NUMBERS <<<"${5}"

	EXCHANGE_CONFIG_JSON_DIR="$6"
	if [ -z "$EXCHANGE_CONFIG_JSON_DIR" ]; then
		# Default to using exchange_config folder within the current directory.
		EXCHANGE_CONFIG_JSON_DIR="exchange_config"
	fi

	DELAY_MSG_JSON_DIR="$7"
	if [ -z "$DELAY_MSG_JSON_DIR" ]; then
		# Default to using exchange_config folder within the current directory.
		DELAY_MSG_JSON_DIR="delaymsg_config"
	fi

	INITIAL_CLOB_PAIR_STATUS="$8"
		if [ -z "$INITIAL_CLOB_PAIR_STATUS" ]; then
		# Default to initialie clob pairs as active.
		INITIAL_CLOB_PAIR_STATUS='STATUS_ACTIVE'
	fi

	REWARDS_VESTER_ACCOUNT_BALANCE="$9"
	if [ -z "$REWARDS_VESTER_ACCOUNT_BALANCE" ]; then
		# Default to 200 million full coins.
		REWARDS_VESTER_ACCOUNT_BALANCE="200000000$EIGHTEEN_ZEROS"
	fi
	
	# Genesis time
	dasel put -t string -f "$GENESIS" '.genesis_time' -v "$GENESIS_TIME"

	# Consensus params
	dasel put -t string -f "$GENESIS" '.consensus.params.block.max_bytes' -v '4194304'
	dasel put -t string -f "$GENESIS" '.consensus.params.block.max_gas' -v '-1'
	dasel put -t string -f "$GENESIS" '.consensus.params.abci.vote_extensions_enable_height' -v '1'

	# Update crisis module.
	dasel put -t string -f "$GENESIS" '.app_state.crisis.constant_fee.denom' -v "$NATIVE_TOKEN"

	# Update gov module.
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.min_deposit.[0].denom' -v "$NATIVE_TOKEN"
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.expedited_min_deposit.[0].denom' -v "$NATIVE_TOKEN"
	# reduced deposit period
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.max_deposit_period' -v '120s'
	# reduced voting period
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.expedited_voting_period' -v '60s'
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.voting_period' -v '120s'
	# set initial deposit ratio to prevent spamming
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.min_initial_deposit_ratio' -v '0.20000' # 20%
	# setting to 1 disables cancelling proposals
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.proposal_cancel_ratio' -v '1'
	# require 75% of votes for an expedited proposal to pass
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.expedited_threshold' -v '0.75000' # 75%

	# Update staking module.
	dasel put -t string -f "$GENESIS" '.app_state.staking.params.unbonding_time' -v '1814400s' # 21 days
	dasel put -t string -f "$GENESIS" '.app_state.staking.params.bond_denom' -v "$NATIVE_TOKEN"

	# Update assets module.
	dasel put -t int -f "$GENESIS" '.app_state.assets.assets.[0].id' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.assets.assets.[0].symbol' -v 'USDC'
	dasel put -t string -f "$GENESIS" '.app_state.assets.assets.[0].denom' -v "$USDC_DENOM"
	dasel put -t string -f "$GENESIS" '.app_state.assets.assets.[0].denom_exponent' -v '-6'
	dasel put -t bool -f "$GENESIS" '.app_state.assets.assets.[0].has_market' -v 'false'
	dasel put -t int -f "$GENESIS" '.app_state.assets.assets.[0].market_id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.assets.assets.[0].atomic_resolution' -v '-6'

	# Update bridge module.
	dasel put -t string -f "$GENESIS" '.app_state.bridge.event_params.denom' -v "$NATIVE_TOKEN"
	dasel put -t int -f "$GENESIS" '.app_state.bridge.event_params.eth_chain_id' -v "$ETH_CHAIN_ID"
	dasel put -t string -f "$GENESIS" '.app_state.bridge.event_params.eth_address' -v "$ETH_BRIDGE_ADDRESS"
	dasel put -t int -f "$GENESIS" '.app_state.bridge.acknowledged_event_info.next_id' -v "$BRIDGE_GENESIS_ACKNOWLEDGED_NEXT_ID"
	dasel put -t int -f "$GENESIS" '.app_state.bridge.acknowledged_event_info.eth_block_height' -v "$BRIDGE_GENESIS_ACKNOWLEDGED_ETH_BLOCK_HEIGHT"

	# Update ibc module.
	dasel put -t bool -f "$GENESIS" '.app_state.ibc.client_genesis.create_localhost' -v "false"
	dasel put -t json -f "$GENESIS" '.app_state.ibc.client_genesis.params.allowed_clients' -v "[]"
	dasel put -t string -f "$GENESIS" '.app_state.ibc.client_genesis.params.allowed_clients.[]' -v "07-tendermint"

	# Update perpetuals module.
	# Liquidity Tiers.
	# TODO(OTE-208): Finalize default values for open interest caps.
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers' -v "[]"
	# Liquidity Tier: Large-Cap
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].id' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].name' -v 'Large-Cap'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].initial_margin_ppm' -v '20000' # 2%
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].maintenance_fraction_ppm' -v '600000' # 60% of IM
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].impact_notional' -v '10000000000' # 10_000 USDC (500 USDC / 5%)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].open_interest_lower_cap' -v '0' # OIMF doesn't apply to Large-Cap
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].open_interest_upper_cap' -v '0' # OIMF doesn't apply to Large-Cap

	# Liquidity Tier: Small-Cap
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].id' -v '1'
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].name' -v 'Small-Cap'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].initial_margin_ppm' -v '100000' # 10%
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].maintenance_fraction_ppm' -v '500000' # 50% of IM
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].impact_notional' -v '5000000000' # 5_000 USDC (500 USDC / 10%)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].open_interest_lower_cap' -v '20000000000000' # 20 million USDC
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].open_interest_upper_cap' -v '50000000000000' # 50 million USDC

	# Liquidity Tier: Long-Tail
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].id' -v '2'
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].name' -v 'Long-Tail'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].initial_margin_ppm' -v '200000' # 20%
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].maintenance_fraction_ppm' -v '500000' # 50% of IM
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].impact_notional' -v '2500000000' # 2_500 USDC (500 USDC / 20%)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].open_interest_lower_cap' -v '5000000000000' # 5 million USDC
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].open_interest_upper_cap' -v '10000000000000' # 10 million USDC

	# Liquidity Tier: Safety
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].id' -v '3'
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].name' -v 'Safety'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].initial_margin_ppm' -v '1000000' # 100%
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].maintenance_fraction_ppm' -v '200000' # 20% of IM
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].impact_notional' -v '2500000000' # 2_500 USDC (2_500 USDC / 100%)
	# For `Safety` IMF is already at 100%; still we set OIMF for completeness.
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].open_interest_lower_cap' -v '2000000000000' # 2 million USDC
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].open_interest_upper_cap' -v '5000000000000' # 5 million USDC

	# Liquidity Tier: Isolated
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[4].id' -v '4'
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[4].name' -v 'Isolated'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[4].initial_margin_ppm' -v '50000' # 5%
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[4].maintenance_fraction_ppm' -v '600000' # 60% of IM
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[4].impact_notional' -v '2500000000' # 2_500 USDC (125 USDC / 5%)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[4].open_interest_lower_cap' -v '500000000000' # 500k USDC
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[4].open_interest_upper_cap' -v '1000000000000' # 1 million USDC

	# Liquidity Tier: Mid-Cap
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[5].id' -v '5'
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[5].name' -v 'Mid-Cap'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[5].initial_margin_ppm' -v '50000' # 5%
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[5].maintenance_fraction_ppm' -v '600000' # 60% of IM
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[5].impact_notional' -v '5000000000' # 5_000 USDC (250 USDC / 5%)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[5].open_interest_lower_cap' -v '40000000000000' # 40 million USDC
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[5].open_interest_upper_cap' -v '100000000000000' # 100 million USDC

	# Liquidity Tier: FX
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[6].id' -v '6'
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[6].name' -v 'FX'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[6].initial_margin_ppm' -v '10000' # 1%
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[6].maintenance_fraction_ppm' -v '500000' # 50% of IM
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[6].impact_notional' -v '2500000000' # 2_500 USDC (25 USDC / 1%)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[6].open_interest_lower_cap' -v '500000000000' # 500k USDC
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[6].open_interest_upper_cap' -v '1000000000000' # 1 million USDC

	# Params.
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.params.funding_rate_clamp_factor_ppm' -v '6000000' # 600 % (same as 75% on hourly rate)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.params.premium_vote_clamp_factor_ppm' -v '60000000' # 6000 % (some multiples of funding rate clamp factor)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.params.min_num_votes_per_sample' -v '15' # half of expected number of votes

	# Perpetuals.
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals' -v "[]"

	# Perpetual: BTC-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.ticker' -v 'BTC-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.market_id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.atomic_resolution' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.liquidity_tier' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.market_type' -v '1'

	# Perpetual: ETH-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.ticker' -v 'ETH-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.id' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.market_id' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.atomic_resolution' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.liquidity_tier' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.market_type' -v '1'

	# Perpetual: LINK-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.ticker' -v 'LINK-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.id' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.market_id' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.market_type' -v '1'

	# Perpetual: POL-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.ticker' -v 'POL-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.id' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.market_id' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.atomic_resolution' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.market_type' -v '1'

	# Perpetual: CRV-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.ticker' -v 'CRV-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.id' -v '4'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.market_id' -v '4'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.atomic_resolution' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.market_type' -v '1'

	# Perpetual: SOL-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.ticker' -v 'SOL-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.id' -v '5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.market_id' -v '5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.atomic_resolution' -v '-7'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.market_type' -v '1'

	# Perpetual: ADA-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.ticker' -v 'ADA-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.id' -v '6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.market_id' -v '6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.atomic_resolution' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.market_type' -v '1'

	# Perpetual: AVAX-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.ticker' -v 'AVAX-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.id' -v '7'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.market_id' -v '7'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.atomic_resolution' -v '-7'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.market_type' -v '1'

	# Perpetual: FIL-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.ticker' -v 'FIL-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.id' -v '8'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.market_id' -v '8'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.market_type' -v '1'

	# Perpetual: LTC-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.ticker' -v 'LTC-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.id' -v '9'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.market_id' -v '9'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.atomic_resolution' -v '-7'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.market_type' -v '1'

	# Perpetual: DOGE-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.ticker' -v 'DOGE-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.id' -v '10'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.market_id' -v '10'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.atomic_resolution' -v '-4'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.market_type' -v '1'

	# Perpetual: ATOM-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.ticker' -v 'ATOM-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.id' -v '11'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.market_id' -v '11'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.market_type' -v '1'

	# Perpetual: DOT-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.ticker' -v 'DOT-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.id' -v '12'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.market_id' -v '12'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.market_type' -v '1'

	# Perpetual: UNI-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.ticker' -v 'UNI-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.id' -v '13'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.market_id' -v '13'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.market_type' -v '1'

	# Perpetual: BCH-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.ticker' -v 'BCH-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.id' -v '14'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.market_id' -v '14'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.atomic_resolution' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.market_type' -v '1'

	# Perpetual: TRX-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.ticker' -v 'TRX-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.id' -v '15'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.market_id' -v '15'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.atomic_resolution' -v '-4'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.market_type' -v '1'

	# Perpetual: NEAR-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.ticker' -v 'NEAR-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.id' -v '16'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.market_id' -v '16'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.liquidity_tier' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.market_type' -v '1'

	# Perpetual: MKR-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.ticker' -v 'MKR-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.id' -v '17'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.market_id' -v '17'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.atomic_resolution' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.liquidity_tier' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.market_type' -v '1'

	# Initialize bridge module account balance as total native token supply.
	bridge_module_account_balance=$TOTAL_NATIVE_TOKEN_SUPPLY
	total_accounts_quote_balance=0
	acct_idx=0
	# Update subaccounts module for load testing accounts and update bridge module account balance.
	for acct in "${INPUT_TEST_ACCOUNTS[@]}"; do
		add_subaccount "$GENESIS" "$acct_idx" "$acct" "$DEFAULT_SUBACCOUNT_QUOTE_BALANCE"
		total_accounts_quote_balance=$(($total_accounts_quote_balance + $DEFAULT_SUBACCOUNT_QUOTE_BALANCE))
		bridge_module_account_balance=$(echo "$bridge_module_account_balance - $TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE" | bc)
		acct_idx=$(($acct_idx + 1))
	done
	# Update subaccounts module for faucet accounts and update bridge module account balance.
	for acct in "${INPUT_FAUCET_ACCOUNTS[@]}"; do
		add_subaccount "$GENESIS" "$acct_idx" "$acct" "$DEFAULT_SUBACCOUNT_QUOTE_BALANCE_FAUCET"
		total_accounts_quote_balance=$(($total_accounts_quote_balance + $DEFAULT_SUBACCOUNT_QUOTE_BALANCE_FAUCET))
		bridge_module_account_balance=$(echo "$bridge_module_account_balance - $TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE" | bc)
		acct_idx=$(($acct_idx + 1))
	done
	# Update subaccounts module for vault accounts.
	for acct in "${INPUT_VAULT_ACCOUNTS[@]}"; do
		add_subaccount "$GENESIS" "$acct_idx" "$acct" "$DEFAULT_SUBACCOUNT_QUOTE_BALANCE_VAULT"
		total_accounts_quote_balance=$(($total_accounts_quote_balance + $DEFAULT_SUBACCOUNT_QUOTE_BALANCE_VAULT))
		acct_idx=$(($acct_idx + 1))
	done
	# Update subaccounts module for megavault main vault account.
	if [ "$DEFAULT_MEGAVAULT_MAIN_VAULT_QUOTE_BALANCE" -gt 0 ]; then
		add_subaccount "$GENESIS" "$acct_idx" "$MEGAVAULT_MAIN_VAULT_ACCOUNT_ADDR" "$DEFAULT_MEGAVAULT_MAIN_VAULT_QUOTE_BALANCE"
		total_accounts_quote_balance=$(($total_accounts_quote_balance + $DEFAULT_MEGAVAULT_MAIN_VAULT_QUOTE_BALANCE))
		acct_idx=$(($acct_idx + 1))
	fi

	next_bank_idx=0
	if (( total_accounts_quote_balance > 0 )); then
		# Initialize subaccounts module account balance within bank module.
		dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[]" -v "{}"
		dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[0].address" -v "${SUBACCOUNTS_MODACC_ADDR}"
		dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[0].coins.[]" -v "{}"
		dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[0].coins.[0].denom" -v "$USDC_DENOM"
		# TODO(DEC-969): For testnet, ensure subaccounts module balance >= sum of subaccount quote balances.
		dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[0].coins.[0].amount" -v "${total_accounts_quote_balance}"
		next_bank_idx=$(($next_bank_idx+1))
	fi

	if [ $(echo "$REWARDS_VESTER_ACCOUNT_BALANCE > 0" | bc -l) -eq 1 ]; then
		# Initialize bank balance of reward vester account.
		dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[]" -v "{}"
		dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].address" -v "${REWARDS_VESTER_ACCOUNT_ADDR}"
		dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].coins.[]" -v "{}"
		dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].coins.[0].denom" -v "${REWARD_TOKEN}"
		dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].coins.[0].amount" -v "$REWARDS_VESTER_ACCOUNT_BALANCE"
		next_bank_idx=$(($next_bank_idx+1))

		bridge_module_account_balance=$(echo "$bridge_module_account_balance - $REWARDS_VESTER_ACCOUNT_BALANCE" | bc)
	fi

	# Initialize bank balance of bridge module account.
	dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[]" -v "{}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].address" -v "${BRIDGE_MODACC_ADDR}"
	dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].coins.[]" -v "{}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].coins.[0].denom" -v "${NATIVE_TOKEN}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].coins.[0].amount" -v "${bridge_module_account_balance}"

	# Set denom metadata
	set_denom_metadata "$NATIVE_TOKEN" "$NATIVE_TOKEN_WHOLE_COIN" "$COIN_NAME"

	# Use ATOM-USD as test oracle price of the reward token.
	dasel put -t int -f "$GENESIS" '.app_state.rewards.params.market_id' -v '11'

	# Update clob module.
	# Clob: BTC-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].id' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[0].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].perpetual_clob_metadata.perpetual_id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].subticks_per_tick' -v '100000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].quantum_conversion_exponent' -v '-9'

	# Clob: ETH-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].id' -v '1'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[1].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].perpetual_clob_metadata.perpetual_id' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].subticks_per_tick' -v '100000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].quantum_conversion_exponent' -v '-9'

	# Clob: LINK-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[2].id' -v '2'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[2].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[2].perpetual_clob_metadata.perpetual_id' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[2].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[2].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[2].quantum_conversion_exponent' -v '-9'

	# Clob: POL-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[3].id' -v '3'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[3].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[3].perpetual_clob_metadata.perpetual_id' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[3].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[3].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[3].quantum_conversion_exponent' -v '-9'

	# Clob: CRV-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[4].id' -v '4'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[4].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[4].perpetual_clob_metadata.perpetual_id' -v '4'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[4].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[4].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[4].quantum_conversion_exponent' -v '-9'

	# Clob: SOL-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[5].id' -v '5'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[5].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[5].perpetual_clob_metadata.perpetual_id' -v '5'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[5].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[5].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[5].quantum_conversion_exponent' -v '-9'

	# Clob: ADA-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[6].id' -v '6'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[6].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[6].perpetual_clob_metadata.perpetual_id' -v '6'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[6].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[6].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[6].quantum_conversion_exponent' -v '-9'

	# Clob: AVAX-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[7].id' -v '7'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[7].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[7].perpetual_clob_metadata.perpetual_id' -v '7'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[7].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[7].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[7].quantum_conversion_exponent' -v '-9'

	# Clob: FIL-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[8].id' -v '8'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[8].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[8].perpetual_clob_metadata.perpetual_id' -v '8'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[8].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[8].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[8].quantum_conversion_exponent' -v '-9'

	# Clob: LTC-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[9].id' -v '9'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[9].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[9].perpetual_clob_metadata.perpetual_id' -v '9'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[9].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[9].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[9].quantum_conversion_exponent' -v '-9'

	# Clob: DOGE-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[10].id' -v '10'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[10].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[10].perpetual_clob_metadata.perpetual_id' -v '10'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[10].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[10].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[10].quantum_conversion_exponent' -v '-9'

	# Clob: ATOM-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[11].id' -v '11'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[11].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[11].perpetual_clob_metadata.perpetual_id' -v '11'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[11].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[11].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[11].quantum_conversion_exponent' -v '-9'

	# Clob: DOT-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[12].id' -v '12'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[12].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[12].perpetual_clob_metadata.perpetual_id' -v '12'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[12].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[12].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[12].quantum_conversion_exponent' -v '-9'

	# Clob: UNI-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[13].id' -v '13'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[13].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[13].perpetual_clob_metadata.perpetual_id' -v '13'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[13].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[13].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[13].quantum_conversion_exponent' -v '-9'

	# Clob: BCH-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[14].id' -v '14'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[14].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[14].perpetual_clob_metadata.perpetual_id' -v '14'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[14].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[14].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[14].quantum_conversion_exponent' -v '-9'

	# Clob: TRX-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[15].id' -v '15'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[15].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[15].perpetual_clob_metadata.perpetual_id' -v '15'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[15].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[15].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[15].quantum_conversion_exponent' -v '-9'

	# Clob: NEAR-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[16].id' -v '16'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[16].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[16].perpetual_clob_metadata.perpetual_id' -v '16'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[16].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[16].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[16].quantum_conversion_exponent' -v '-9'

	# Clob: MKR-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[17].id' -v '17'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[17].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[17].perpetual_clob_metadata.perpetual_id' -v '17'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[17].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[17].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[17].quantum_conversion_exponent' -v '-9'

	# Liquidations
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.max_liquidation_fee_ppm' -v '15000'  # 1.5%
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.position_block_limits.min_position_notional_liquidated' -v '1000000000' # 1_000 USDC
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.position_block_limits.max_position_portion_liquidated_ppm' -v '100000'  # 10%
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.subaccount_block_limits.max_notional_liquidated' -v '100000000000'  # 100_000 USDC
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.subaccount_block_limits.max_quantums_insurance_lost' -v '1000000000000' # 1_000_000 USDC
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.fillable_price_config.bankruptcy_adjustment_ppm' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.fillable_price_config.spread_to_maintenance_margin_ratio_ppm' -v '1500000'  # 150%

	# Block Rate Limit
	# Max 400 short term orders/cancels per block
	dasel put -t json -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_short_term_orders_and_cancels_per_n_blocks.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_short_term_orders_and_cancels_per_n_blocks.[0].limit' -v '400'
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_short_term_orders_and_cancels_per_n_blocks.[0].num_blocks' -v '1'
	# Max 2 stateful orders per block
	dasel put -t json -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_stateful_orders_per_n_blocks.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_stateful_orders_per_n_blocks.[0].limit' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_stateful_orders_per_n_blocks.[0].num_blocks' -v '1'
	# Max 20 stateful orders per 100 blocks
	dasel put -t json -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_stateful_orders_per_n_blocks.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_stateful_orders_per_n_blocks.[1].limit' -v '20'
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_stateful_orders_per_n_blocks.[1].num_blocks' -v '100'

	# Equity Tier Limit
	# Max 0 open short term orders for $0 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[0].limit' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[0].usd_tnc_required' -v '0'
	# Max 1 open short term orders for $20 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[1].limit' -v '1'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[1].usd_tnc_required' -v '20000000'
	# Max 5 open short term orders for $100 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[2].limit' -v '5'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[2].usd_tnc_required' -v '100000000'
	# Max 10 open short term orders for $1000 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[3].limit' -v '10'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[3].usd_tnc_required' -v '1000000000'
	# Max 100 open short term orders for $10,000 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[4].limit' -v '100'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[4].usd_tnc_required' -v '10000000000'
	# Max 200 open short term orders for $100,000 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[5].limit' -v '1000'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[5].usd_tnc_required' -v '100000000000'
	# Max 0 open stateful orders for $0 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[0].limit' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[0].usd_tnc_required' -v '0'
	# Max 1 open stateful orders for $20 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[1].limit' -v '1'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[1].usd_tnc_required' -v '20000000'
	# Max 5 open stateful orders for $100 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[2].limit' -v '5'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[2].usd_tnc_required' -v '100000000'
	# Max 10 open stateful orders for $1000 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[3].limit' -v '10'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[3].usd_tnc_required' -v '1000000000'
	# Max 100 open stateful orders for $10,000 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[4].limit' -v '100'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[4].usd_tnc_required' -v '10000000000'
	# Max 200 open stateful orders for $100,000 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[5].limit' -v '200'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[5].usd_tnc_required' -v '100000000000'


  # Fee Tiers
  # Schedule a delayed message to swap fee tiers to the standard schedule after ~120 days of blocks.
	dasel put -t int -f "$GENESIS" '.app_state.delaymsg.next_delayed_message_id' -v '1'
	dasel put -t json -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[0].id' -v '0'

	delaymsg=$(cat "$DELAY_MSG_JSON_DIR/perpetual_fee_params_msg.json" | jq -c '.')
	dasel put -t json -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[0].msg' -v "$delaymsg"
	# Schedule the message to execute in ~7 days (at 1.6s per block.)
	dasel put -t int -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[0].block_height' -v '378000'
	# Uncomment the following to schedule the message to execute in ~120 days (at 1.6s per block.)
	# dasel put -t int -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[0].block_height' -v '6480000'

	# ICA Host Params
	update_ica_host_params
	# ICA Controller Params
	update_ica_controller_params

	# Listing
	# Set hard cap for markets
	dasel put -t int -f "$GENESIS" ".app_state.listing.hard_cap_for_markets" -v '500'
	# Set default listing vault deposit params
	dasel put -t string -f "$GENESIS" ".app_state.listing.listing_vault_deposit_params.new_vault_deposit_amount" -v "10000000000" # 10_000 USDC
	dasel put -t string -f "$GENESIS" ".app_state.listing.listing_vault_deposit_params.main_vault_deposit_amount" -v "0" # 0 USDC
	dasel put -t int -f "$GENESIS" ".app_state.listing.listing_vault_deposit_params.num_blocks_to_lock_shares" -v '2592000' # 30 days

	# Vaults
	# Set default quoting params.
	dasel put -t int -f "$GENESIS" ".app_state.vault.default_quoting_params.spread_min_ppm" -v '3000'
	# Set operator params.
	dasel put -t string -f "$GENESIS" ".app_state.vault.operator_params.operator" -v 'vindax1vlthgax23ca9syk7xgaz347xmf4nunefwppm9c'
	dasel put -t string -f "$GENESIS" ".app_state.vault.operator_params.metadata.name" -v 'Governance'
	dasel put -t string -f "$GENESIS" ".app_state.vault.operator_params.metadata.description" -v 'Governance Module Account'
	# Set total shares and owner shares.
	if [ -z "${INPUT_TEST_ACCOUNTS[0]}" ]; then
		vault_owner_address='vindax199tqg4wdlnu4qjlxchpd7seg454937hjrx2642' # alice as default vault owner
	else
		vault_owner_address="${INPUT_TEST_ACCOUNTS[0]}"
	fi
	total_deposit=$((DEFAULT_SUBACCOUNT_QUOTE_BALANCE_VAULT * ${#INPUT_VAULT_NUMBERS[@]})) 
	dasel put -t string -f "$GENESIS" ".app_state.vault.total_shares.num_shares" -v "${total_deposit}"
	dasel put -t json -f "$GENESIS" ".app_state.vault.owner_shares.[]" -v '{}'
	dasel put -t string -f "$GENESIS" ".app_state.vault.owner_shares.[0].owner" -v "${vault_owner_address}"
	dasel put -t string -f "$GENESIS" ".app_state.vault.owner_shares.[0].shares.num_shares" -v "${total_deposit}"
	# Set vaults.
	vault_idx=0
	for number in "${INPUT_VAULT_NUMBERS[@]}"; do
		dasel put -t json -f "$GENESIS" '.app_state.vault.vaults.[]' -v '{}'
		dasel put -t string -f "$GENESIS" ".app_state.vault.vaults.[${vault_idx}].vault_id.type" -v 'VAULT_TYPE_CLOB'
		dasel put -t int -f "$GENESIS" ".app_state.vault.vaults.[${vault_idx}].vault_id.number" -v "${number}"
		dasel put -t string -f "$GENESIS" ".app_state.vault.vaults.[${vault_idx}].vault_params.status" -v 'VAULT_STATUS_QUOTING'
		vault_idx=$(($vault_idx + 1))
	done

	# Update accountplus module.
	dasel put -t bool -f "$GENESIS" '.app_state.vindaxaccountplus.params.is_smart_account_active' -v 'true'
}

function add_subaccount() {
	GEN_FILE=$1
	IDX=$2
	ACCOUNT=$3
	QUOTE_BALANCE=$4

	dasel put -t json -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[]" -v "{}"
	dasel put -t json -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].id" -v "{}"
	dasel put -t string -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].id.owner" -v "$ACCOUNT"
	dasel put -t int -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].id.number" -v '0'
	dasel put -t bool -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].margin_enabled" -v 'true'

	dasel put -t json -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].asset_positions.[]" -v '{}'
	dasel put -t int -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].asset_positions.[0].asset_id" -v '0'
	dasel put -t string -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].asset_positions.[0].quantums" -v "${QUOTE_BALANCE}"
	dasel put -t int -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].asset_positions.[0].index" -v '0'
}

# Modify the genesis file to use only the test exchange for computing index prices. The test exchange is configured
# to serve prices for BTC, ETH and LINK. This must be called after edit_genesis to ensure all markets exist.
function update_genesis_use_test_exchange() {
	GENESIS=$1/genesis.json

	# For BTC, ETH and LINK, remove all exchanges except the test exchange.
	btc_exchange_config_json=$(cat <<-EOF
	{
		"exchanges": [
			{
				"exchangeName": "TestExchange",
				"ticker": "BTC-USD"
			}
		]
	}
	EOF
	)
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[0].exchange_config_json' -v "$btc_exchange_config_json"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[0].min_exchanges' -v '1'

	eth_exchange_config_json=$(cat <<-EOF
	{
		"exchanges": [
			{
				"exchangeName": "TestExchange",
				"ticker": "ETH-USD"
			}
		]
	}
	EOF
	)
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[1].exchange_config_json' -v "$eth_exchange_config_json"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[1].min_exchanges' -v '1'

	link_exchange_config_json=$(cat <<-EOF
	{
		"exchanges": [
			{
				"exchangeName": "TestExchange",
				"ticker": "LINK-USD"
			}
		]
	}
	EOF
	)
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[2].exchange_config_json' -v "$link_exchange_config_json"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[2].min_exchanges' -v '1'

	# All remaining markets can just use the LINK ticker so the daemon will start. All markets must have at least 1
	# exchange. With only one exchange configured, there should not be enough prices to meet the minimum exchange
	# count, and these markets will not have index prices.
	for market_idx in {3..34}
	do
		dasel put -t string -f "$GENESIS" ".app_state.prices.market_params.[$market_idx].exchange_config_json" -v "$link_exchange_config_json"
	done
}

# Modify the genesis file to add test volatile market. Market TEST-USD will be added as market 33.
function update_genesis_use_test_volatile_market() {
	GENESIS=$1/genesis.json
	TEST_USD_MARKET_ID=33

	# Market: TEST-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.last().pair' -v 'TEST-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.last().id' -v "${TEST_USD_MARKET_ID}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.last().exponent' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.last().min_price_change_ppm' -v '250' # 0.025%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.last().id' -v "${TEST_USD_MARKET_ID}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.last().exponent' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.last().price' -v '10000000'          # $100 = 1 TEST.

	# Liquidity Tier: For TEST-USD. 1% leverage and regular 1m nonlinear margin thresholds.
	NUM_LIQUIDITY_TIERS=$(jq -c '.app_state.perpetuals.liquidity_tiers | length' < ${GENESIS})
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().id' -v "${NUM_LIQUIDITY_TIERS}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().name' -v 'test-usd-100x-liq-tier-linear'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().initial_margin_ppm' -v '10007' # 1% + a little prime (100x leverage)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().maintenance_fraction_ppm' -v '500009' # 50% of IM + a little prime
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().impact_notional' -v '50000000000' # 50_000 USDC (500 USDC / 1%)

	# Liquidity Tier: For TEST-USD. 1% leverage and 100 nonlinear margin thresholds.
	NUM_LIQUIDITY_TIERS_2=$(jq -c '.app_state.perpetuals.liquidity_tiers | length' < ${GENESIS})
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().id' -v "${NUM_LIQUIDITY_TIERS_2}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().name' -v 'test-usd-100x-liq-tier-nonlinear'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().initial_margin_ppm' -v '10007' # 1% + a little prime (100x leverage)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().maintenance_fraction_ppm' -v '500009' # 50% of IM + a little prime
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().impact_notional' -v '50000000000' # 50_000 USDC (500 USDC / 1%)

	# Perpetual: TEST-USD
	NUM_PERPETUALS=$(jq -c '.app_state.perpetuals.perpetuals | length' < ${GENESIS})
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.ticker' -v 'TEST-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.id' -v "${NUM_PERPETUALS}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.market_id' -v "${TEST_USD_MARKET_ID}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.atomic_resolution' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.default_funding_ppm' -v '100'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.liquidity_tier' -v "${NUM_LIQUIDITY_TIERS}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.market_type' -v '1'


	# Clob: TEST-USD
	NUM_CLOB_PAIRS=$(jq -c '.app_state.clob.clob_pairs | length' < ${GENESIS})
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.last().id' -v "${NUM_CLOB_PAIRS}"
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.last().status' -v 'STATUS_ACTIVE'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.last().perpetual_clob_metadata.perpetual_id' -v "${NUM_PERPETUALS}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.last().step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.last().subticks_per_tick' -v '100' # $0.01 ticks
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.last().quantum_conversion_exponent' -v '-8'
}

# Modify the genesis file with ICA Host params that are consistent with
# v3.0.0 upgrade.
function update_ica_host_params() {
	dasel put -t json -f "$GENESIS" '.app_state.interchainaccounts.host_genesis_state.params.allow_messages' -v '[]'
	dasel put -t string -f "$GENESIS" '.app_state.interchainaccounts.host_genesis_state.params.allow_messages.[]' -v "/ibc.applications.transfer.v1.MsgTransfer"
	dasel put -t string -f "$GENESIS" '.app_state.interchainaccounts.host_genesis_state.params.allow_messages.[]' -v "/cosmos.bank.v1beta1.MsgSend"
	dasel put -t string -f "$GENESIS" '.app_state.interchainaccounts.host_genesis_state.params.allow_messages.[]' -v "/cosmos.staking.v1beta1.MsgDelegate"
	dasel put -t string -f "$GENESIS" '.app_state.interchainaccounts.host_genesis_state.params.allow_messages.[]' -v "/cosmos.staking.v1beta1.MsgBeginRedelegate"
	dasel put -t string -f "$GENESIS" '.app_state.interchainaccounts.host_genesis_state.params.allow_messages.[]' -v "/cosmos.staking.v1beta1.MsgUndelegate"
	dasel put -t string -f "$GENESIS" '.app_state.interchainaccounts.host_genesis_state.params.allow_messages.[]' -v "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation"
	dasel put -t string -f "$GENESIS" '.app_state.interchainaccounts.host_genesis_state.params.allow_messages.[]' -v "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress"
	dasel put -t string -f "$GENESIS" '.app_state.interchainaccounts.host_genesis_state.params.allow_messages.[]' -v "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"
	dasel put -t string -f "$GENESIS" '.app_state.interchainaccounts.host_genesis_state.params.allow_messages.[]' -v "/cosmos.distribution.v1beta1.MsgFundCommunityPool"
	dasel put -t string -f "$GENESIS" '.app_state.interchainaccounts.host_genesis_state.params.allow_messages.[]' -v "/cosmos.gov.v1.MsgVote"
}

function update_ica_controller_params() {
	dasel put -t bool -f "$GENESIS" '.app_state.interchainaccounts.controller_genesis_state.params.controller_enabled' -v "false"
}

# Modify the genesis file to only use fixed price exchange.
function update_all_markets_with_fixed_price_exchange() {
    GENESIS=$1/genesis.json

    # Read the number of markets
    NUM_MARKETS=$(jq -c '.app_state.prices.market_params | length' < "${GENESIS}")

    # Loop through each market and update the parameters
    for ((j = 0; j < NUM_MARKETS; j++)); do
        # Get the current ticker
        TICKER=$(jq -r ".app_state.prices.market_params[$j].pair" < "${GENESIS}")

        # Update the exchange_config_json using the EOF syntax
        exchange_config_json=$(cat <<-EOF
{
    "exchanges": [
        {
            "exchangeName": "TestFixedPriceExchange",
            "ticker": "${TICKER}"
        }
    ]
}
EOF
        )
        dasel put -t string -f "$GENESIS" ".app_state.prices.market_params.[$j].exchange_config_json" -v "$exchange_config_json"

        # Update the min_exchanges
        dasel put -t int -f "$GENESIS" ".app_state.prices.market_params.[$j].min_exchanges" -v "1"
    done
}

# Modify the genesis file with reduced complete bridge delay (for testing in non-prod envs).
update_genesis_complete_bridge_delay() {
	GENESIS=$1/genesis.json

	# Reduce complete bridge delay to 600 blocks.
	dasel put -t int -f "$GENESIS" '.app_state.bridge.safety_params.delay_blocks' -v "$2"
}

# Set the denom metadata, which is for human readability.
set_denom_metadata() {
	local BASE_DENOM=$1
	local WHOLE_COIN_DENOM=$2
	local COIN_NAME=$3
	dasel put -t json -f "$GENESIS" ".app_state.bank.denom_metadata" -v "[]"
	dasel put -t json -f "$GENESIS" ".app_state.bank.denom_metadata.[]" -v "{}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].description" -v "The native token of the network"
	dasel put -t json -f "$GENESIS" ".app_state.bank.denom_metadata.[0].denom_units" -v "[]"
	dasel put -t json -f "$GENESIS" ".app_state.bank.denom_metadata.[0].denom_units.[]" -v "{}"
	# Base denom is the minimum unit of the a token and the denom used by `x/bank`.
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].denom_units.[0].denom" -v "$BASE_DENOM"
	dasel put -t json -f "$GENESIS" ".app_state.bank.denom_metadata.[0].denom_units.[]" -v "{}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].denom_units.[1].denom" -v "$WHOLE_COIN_DENOM"
	dasel put -t int -f "$GENESIS" ".app_state.bank.denom_metadata.[0].denom_units.[1].exponent" -v 18
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].base" -v "$BASE_DENOM"
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].name" -v "$COIN_NAME"
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].display" -v "$WHOLE_COIN_DENOM"
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].symbol" -v "$WHOLE_COIN_DENOM"
}
