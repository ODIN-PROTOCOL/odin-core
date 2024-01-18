#!/bin/bash
#set -o errexit -o nounset -o pipefail

PASSWORD=${PASSWORD:-1234567890}
STAKE=${STAKE_TOKEN:-loki}
FEE=${FEE_TOKEN:-loki}
CHAIN_ID=${CHAIN_ID:-testing}
MONIKER=${MONIKER:-node001}

odind init $MONIKER --chain-id $CHAIN_ID

if ! odind keys show val1; then
  (echo $PASSWORD; echo $PASSWORD) | odind keys add val1
fi
# hardcode the val1 account for this instance
echo $PASSWORD | odind genesis add-genesis-account val1 100000000000$STAKE

# (optionally) add a few more genesis accounts
for addr in $@; do
  echo $addr
  odind genesis add-genesis-account $addr 1000000000000$STAKE
done

# submit a genesis val1 tx
## Workraround for https://github.com/cosmos/cosmos-sdk/issues/8251
(echo $PASSWORD; echo $PASSWORD; echo $PASSWORD) | odind genesis gentx val1 250000000000$STAKE --chain-id=$CHAIN_ID
odind genesis collect-gentxs
