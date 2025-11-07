#!/usr/bin/env bash
set -euo pipefail

# =====================[ Configuration ]=====================
BIN="${BIN:-dydxprotocold}"
KEYRING="${KEYRING:-test}"             # test | file | os
HOME_DIR="${HOME_DIR:-$HOME/.dydxprotocold}" # change if different
RESET="${RESET:-true}"

# =====================[ Account list & mnemonic ]=====================
NAMES=(
  "alice" "bob" "carl" "dave" "emily"
  "fiona" "greg" "henry" "ian" "jeff"
)

MNEMONICS=(
"merge panther lobster crazy road hollow amused security before critic about cliff exhibit cause coyote talent happy where lion river tobacco option coconut small"
"color habit donor nurse dinosaur stable wonder process post perfect raven gold census inside worth inquiry mammal panic olive toss shadow strong name drum"
"school artefact ghost shop exchange slender letter debris dose window alarm hurt whale tiger find found island what engine ketchup globe obtain glory manage"
"switch boring kiss cash lizard coconut romance hurry sniff bus accident zone chest height merit elevator furnace eagle fetch quit toward steak mystery nest"
"brave way sting spin fog process matrix glimpse volcano recall day lab raccoon hand path pig rent mixture just way blouse alone upon prefer"
"suffer claw truly wife simple mean still mammal bind cake truly runway attack burden lazy peanut unusual such shock twice appear gloom priority kind"
"step vital slight present group gallery flower gap copy sweet travel bitter arena reject evidence deal ankle motion dismiss trim armed slab life future"
"piece choice region bike tragic error drive defense air venture bean solve income upset physical sun link actor task runway match gauge brand march"
"burst section toss rotate law thumb shoe wire only decide meadow aunt flight humble story mammal radar scene wrist essay taxi leisure excess milk"
"fashion charge estate devote jaguar fun swift always road lend scrap panic matter core defense high gas athlete permit crane assume pact fitness matrix"
)

# =====================[ Environment setup ]=====================
command -v "$BIN" >/dev/null 2>&1 || { echo "❌ Binary not found: $BIN"; exit 1; }

if [ "$RESET" = "true" ] && [ -d "$HOME_DIR" ]; then
  echo "🧹 Removing old test keyring data: $HOME_DIR"
  rm -rf "$HOME_DIR"
fi

mkdir -p "$HOME_DIR"

echo "🔑 Creating keyring and generating addresses..."
echo

# =====================[ Generate addresses ]=====================
for i in "${!NAMES[@]}"; do
  NAME="${NAMES[$i]}"
  MN="${MNEMONICS[$i]}"

  # Import key from mnemonic
  printf "%s\n" "$MN" | $BIN keys add "$NAME" --recover \
    --keyring-backend "$KEYRING" \
    --home "$HOME_DIR" >/dev/null

  ADDR=$($BIN keys show "$NAME" -a --keyring-backend "$KEYRING" --home "$HOME_DIR")
  VALOPER=$($BIN keys show "$NAME" --bech val --address --keyring-backend "$KEYRING" --home "$HOME_DIR")
  VALCONS=$($BIN tendermint show-validator --home "$HOME_DIR" 2>/dev/null || echo "—")

  echo "👤 $NAME"
  echo "  Account:  $ADDR"
  echo "  Valoper:  ${VALOPER:-—}"
  echo "  Valcons:  ${VALCONS:-—}"
  echo
done

echo "✅ Complete! All addresses have been generated."
