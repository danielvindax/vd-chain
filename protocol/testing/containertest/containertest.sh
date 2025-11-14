#!/bin/bash
set -eo pipefail

# This file initializes muliple validators for local and CI testing purposes.
# This file should be run as part of `docker-compose.yml`.

source "./genesis.sh"
source "./version.sh"

CHAIN_ID="localvindax"

# Define mnemonics for all validators.
MNEMONICS=(
	# alice
	# Consensus Address: vindaxvalcons1zf9csp5ygq95cqyxh48w3qkuckmpealrw2ug4d
	"merge panther lobster crazy road hollow amused security before critic about cliff exhibit cause coyote talent happy where lion river tobacco option coconut small"

	# bob
	# Consensus Address: vindaxvalcons1s7wykslt83kayxuaktep9fw8qxe5n73ucftkh4
	"color habit donor nurse dinosaur stable wonder process post perfect raven gold census inside worth inquiry mammal panic olive toss shadow strong name drum"

	# carl
	# Consensus Address: vindaxvalcons1vy0nrh7l4rtezrsakaadz4mngwlpdmhy64h0ls
	"school artefact ghost shop exchange slender letter debris dose window alarm hurt whale tiger find found island what engine ketchup globe obtain glory manage"

	# dave
	# Consensus Address: vindaxvalcons1stjspktkshgcsv8sneqk2vs2ws0nw2wr272vtt
	"switch boring kiss cash lizard coconut romance hurry sniff bus accident zone chest height merit elevator furnace eagle fetch quit toward steak mystery nest"
)

# Define node keys for all validators.
NODE_KEYS=(
	# Node ID: 17e5e45691f0d01449c84fd4ae87279578cdd7ec
	"8EGQBxfGMcRfH0C45UTedEG5Xi3XAcukuInLUqFPpskjp1Ny0c5XvwlKevAwtVvkwoeYYQSe0geQG/cF3GAcUA=="

	# Node ID: b69182310be02559483e42c77b7b104352713166
	"3OZf5HenMmeTncJY40VJrNYKIKcXoILU5bkYTLzTJvewowU2/iV2+8wSlGOs9LoKdl0ODfj8UutpMhLn5cORlw=="

	# Node ID: 47539956aaa8e624e0f1d926040e54908ad0eb44
	"tWV4uEya9Xvmm/kwcPTnEQIV1ZHqiqUTN/jLPHhIBq7+g/5AEXInokWUGM0shK9+BPaTPTNlzv7vgE8smsFg4w=="

	# Node ID: 5882428984d83b03d0c907c1f0af343534987052
	"++C3kWgFAs7rUfwAHB7Ffrv43muPg0wTD2/UtSPFFkhtobooIqc78UiotmrT8onuT1jg8/wFPbSjhnKRThTRZg=="
)

# Define monikers for each validator. These are made up strings and can be anything.
# This also controls in which directory the validator's home will be located. i.e. `/vindax/chain/.alice`
MONIKERS=(
	"alice"
	"bob"
	"carl"
	"dave"
)

# Define all test accounts for the chain.
TEST_ACCOUNTS=(
	"vindax199tqg4wdlnu4qjlxchpd7seg454937hjrx2642" # alice
	"vindax10fx7sy6ywd5senxae9dwytf8jxek3t2gcf2z90" # bob
	"vindax1fjg6zp6vv8t9wvy4lps03r5l4g7tkjw9wuzlhs" # carl
	"vindax1wau5mja7j7zdavtfq9lu7ejef05hm6fferxsev" # dave
)

FAUCET_ACCOUNTS=(
	"vindax10ehz9v9ncpnj8hfwlsxhcg97zv5ag5w3sgac4k" # main faucet
)

# Define dependencies for this script.
# `jq` and `dasel` are used to manipulate json and yaml files respectively.
install_prerequisites() {
	apk add curl dasel jq
}


log() { printf "[init:%(%F %T)T] %s\n" -1 "$*"; }

# Get HRP part (prefix before '1') of a bech32 address
addr_hrp() {
  local addr="$1"
  echo "${addr%%1*}"
}

# Create a temporary key to read address -> infer HRP that binary is using
# Return HRP via stdout
detect_binary_hrp() {
  local home_dir="$1"
  local tmpkey="__probe__$$"
  # Create temporary key and output JSON to get .address
  local addr
  addr="$(vindaxd keys add "$tmpkey" --keyring-backend=test --home "$home_dir" -y --output json 2>/dev/null | jq -r '.address' || true)"
  # Delete temporary key immediately (to avoid dirty keyring)
  vindaxd keys delete "$tmpkey" --keyring-backend=test --home "$home_dir" -y >/dev/null 2>&1 || true

  if [[ -z "$addr" || "$addr" == "null" ]]; then
    echo ""
    return 1
  fi
  addr_hrp "$addr"
}

# Print important environment information
log_env() {
  log "Binary: $(command -v vindaxd || echo 'not found')"
  log "Version: $(vindaxd version 2>/dev/null || echo 'unknown')"
  log "CHAIN_ID=$CHAIN_ID"
  log "USDC_DENOM=${USDC_DENOM:-<unset>}, NATIVE_TOKEN=${NATIVE_TOKEN:-<unset>}"
  log "TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE=${TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE:-<unset>}"
  log "TESTNET_VALIDATOR_SELF_DELEGATE_AMOUNT=${TESTNET_VALIDATOR_SELF_DELEGATE_AMOUNT:-<unset>}"

  # Get expected HRP from TEST_ACCOUNTS array (first element)
  local first_acct="${TEST_ACCOUNTS[0]}"
  local expected_hrp="$(addr_hrp "$first_acct")"
  log "Expected HRP from TEST_ACCOUNTS: $expected_hrp"
}

# Check HRP: if binary HRP != configured address HRP -> clear warning
preflight_hrp_check() {
  local home_dir="$1"
  local expected_hrp="$2"

  local binary_hrp
  binary_hrp="$(detect_binary_hrp "$home_dir" || true)"
  if [[ -z "$binary_hrp" ]]; then
    log "WARN: Unable to detect binary HRP (keygen probe failed)."
  else
    log "Detected binary HRP: $binary_hrp (home: $home_dir)"
    if [[ "$binary_hrp" != "$expected_hrp" ]]; then
      log "ERROR: HRP mismatch → binary uses '$binary_hrp' but TEST/FAUCET accounts use '$expected_hrp'."
      log "       This mismatch causes error 'failed to get address from Keybase ... key not found'."
      log "       Instructions: rebuild binary with SetBech32Prefix('$(addr_hrp "${NATIVE_TOKEN:-vindax}")'...) or change all addresses to HRP '$binary_hrp'."
      return 2
    fi
  fi
  return 0
}

# Quick check parse address with current binary (for easier log understanding)
probe_parse_address() {
  local addr="$1"
  local tag="$2"
  # Try to sign 1 fake tx (dry-run not available), so only log simple parse with regex and compare HRP
  local hrp="$(addr_hrp "$addr")"
  if [[ -z "$hrp" ]]; then
    log "[$tag] INVALID: '$addr' does not have valid bech32 format (no '1' found)."
    return 1
  fi
  log "[$tag] Address='$addr' | HRP='$hrp'"
}

# Create all validators for the chain including a full-node.
# Initialize their genesis files and home directories.
create_validators() {
	# Create temporary directory for all gentx files.
	mkdir /tmp/gentx

	# Iterate over all validators and set up their home directories, as well as generate `gentx` transaction for each.
	for i in "${!MONIKERS[@]}"; do
		VAL_HOME_DIR="$HOME/chain/.${MONIKERS[$i]}"
		VAL_CONFIG_DIR="$VAL_HOME_DIR/config"

		# Initialize the chain and validator files.
		vindaxd init "${MONIKERS[$i]}" -o --chain-id=$CHAIN_ID --home "$VAL_HOME_DIR"

		# Overwrite the randomly generated `priv_validator_key.json` with a key generated deterministically from the mnemonic.
		vindaxd tendermint gen-priv-key --home "$VAL_HOME_DIR" --mnemonic "${MNEMONICS[$i]}"

		# Note: `vindaxd init` non-deterministically creates `node_id.json` for each validator.
		# This is inconvenient for persistent peering during testing in Terraform configuration as the `node_id`
		# would change with every build of this container.
		#
		# For that reason we overwrite the non-deterministically generated one with a deterministic key defined in this file here.
		new_file=$(jq ".priv_key.value = \"${NODE_KEYS[$i]}\"" "$VAL_CONFIG_DIR"/node_key.json)
		cat <<<"$new_file" >"$VAL_CONFIG_DIR"/node_key.json

		edit_config "$VAL_CONFIG_DIR"

		# Using "*" as a subscript results in a single arg: "vindax1... vindax1... vindax1..."
		# Using "@" as a subscript results in separate args: "vindax1..." "vindax1..." "vindax1..."
		# Note: `edit_genesis` must be called before `add-genesis-account` or `update_genesis_use_test_exchange`.
		edit_genesis "$VAL_CONFIG_DIR" "${TEST_ACCOUNTS[*]}" "${FAUCET_ACCOUNTS[*]}" "" "" "" "" "" ""
		# Configure the genesis file to only use the test exchange to compute index prices.
		update_genesis_use_test_exchange "$VAL_CONFIG_DIR"

		echo "${MNEMONICS[$i]}" | vindaxd keys add "${MONIKERS[$i]}" --recover --keyring-backend=test --home "$VAL_HOME_DIR"
		
		EXPECTED_HRP="$(addr_hrp "${TEST_ACCOUNTS[0]}")"
		preflight_hrp_check "$VAL_HOME_DIR" "$EXPECTED_HRP" || {
		log "Abort at moniker='${MONIKERS[$i]}' due to HRP mismatch. Please fix HRP and run again."
		exit 1
		}

		# Log try parsing addresses to be added for easier debugging
		for acct in "${TEST_ACCOUNTS[@]}"; do
		probe_parse_address "$acct" "TEST_ACCOUNT"
		done
		for acct in "${FAUCET_ACCOUNTS[@]}"; do
		probe_parse_address "$acct" "FAUCET_ACCOUNT"
		done

		for acct in "${TEST_ACCOUNTS[@]}"; do
			vindaxd add-genesis-account "$acct" 100000000000000000$USDC_DENOM,$TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE$NATIVE_TOKEN --home "$VAL_HOME_DIR"
		done
		for acct in "${FAUCET_ACCOUNTS[@]}"; do
			vindaxd add-genesis-account "$acct" 900000000000000000$USDC_DENOM,$TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE$NATIVE_TOKEN --home "$VAL_HOME_DIR"
		done

		vindaxd gentx "${MONIKERS[$i]}" $TESTNET_VALIDATOR_SELF_DELEGATE_AMOUNT$NATIVE_TOKEN --moniker="${MONIKERS[$i]}" --keyring-backend=test --chain-id=$CHAIN_ID --home "$VAL_HOME_DIR"

		# Copy the gentx to a shared directory.
		cp -a "$VAL_CONFIG_DIR/gentx/." /tmp/gentx
	done

	# Copy gentxs to the first validator's home directory to build the genesis json file
	FIRST_VAL_HOME_DIR="$HOME/chain/.${MONIKERS[0]}"
	FIRST_VAL_CONFIG_DIR="$FIRST_VAL_HOME_DIR/config"

	rm -rf "$FIRST_VAL_CONFIG_DIR/gentx"
	mkdir "$FIRST_VAL_CONFIG_DIR/gentx"
	cp -r /tmp/gentx "$FIRST_VAL_CONFIG_DIR"

	# Build the final genesis.json file that all validators and the full-nodes will use.
	vindaxd collect-gentxs --home "$FIRST_VAL_HOME_DIR"

	# Copy this genesis file to each of the other validators
	for i in "${!MONIKERS[@]}"; do
		if [[ "$i" == 0 ]]; then
			# Skip first moniker as it already has the correct genesis file.
			continue
		fi

		VAL_HOME_DIR="$HOME/chain/.${MONIKERS[$i]}"
		VAL_CONFIG_DIR="$VAL_HOME_DIR/config"
		rm -rf "$VAL_CONFIG_DIR/genesis.json"
		cp "$FIRST_VAL_CONFIG_DIR/genesis.json" "$VAL_CONFIG_DIR/genesis.json"
	done
}

setup_cosmovisor() {
	for i in "${!MONIKERS[@]}"; do
		VAL_HOME_DIR="$HOME/chain/.${MONIKERS[$i]}"
		export DAEMON_NAME=vindaxd
		export DAEMON_HOME="$HOME/chain/.${MONIKERS[$i]}"

		cosmovisor init /bin/vindaxd
	done
}

download_preupgrade_binary() {
	arch="$(apk --print-arch)"
	url_arch=""
	case "$arch" in
		'x86_64')
			url_arch='amd64'
			;;
		'aarch64')
			url_arch='arm64'
			;;
		*)
			echo >&2 "unexpected architecture '$arch'"
			exit 1
			;;
	esac
	tar_url="https://github.com/danielvindax/vd-chain/releases/download/protocol%2F$PREUPGRADE_VERSION_FULL_NAME/vindaxd-$PREUPGRADE_VERSION_FULL_NAME-linux-$url_arch.tar.gz"
	tar_path='/tmp/vindaxd/vindaxd.tar.gz'
	mkdir -p /tmp/vindaxd
	curl -vL $tar_url -o $tar_path
	vindaxd_path=$(tar -xvf $tar_path --directory /tmp/vindaxd)
	cp /tmp/vindaxd/$vindaxd_path /bin/vindaxd_preupgrade
}

# TODO(DEC-1894): remove this function once we migrate off of persistent peers.
# Note: DO NOT add more config modifications in this method. Use `cmd/config.go` to configure
# the default config values.
edit_config() {
	CONFIG_FOLDER=$1

	# Disable pex
	dasel put -t bool -f "$CONFIG_FOLDER"/config.toml '.p2p.pex' -v 'false'

	# Default `timeout_commit` is 999ms. For local testnet, use a larger value to make 
	# block time longer for easier troubleshooting.
	dasel put -t string -f "$CONFIG_FOLDER"/config.toml '.consensus.timeout_commit' -v '5s'
}
log_env
install_prerequisites
setup_cosmovisor
download_preupgrade_binary
create_validators
