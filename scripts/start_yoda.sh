#!/bin/bash

rm -rf ~/.yoda

# config chain id
yoda config chain-id odinchain

# add validator to yoda config
yoda config validator $(odind keys show validator -a --bech val --keyring-backend test)

# setup execution endpoint
yoda config executor "rest:$EXECUTOR_URL?timeout=10s"

# setup broadcast-timeout to yoda config
yoda config broadcast-timeout "5m"

# setup rpc-poll-interval to yoda config
yoda config rpc-poll-interval "1s"

# setup max-try to yoda config
yoda config max-try 5

echo "y" | odind tx oracle activate --from validator --gas-prices 0.0025loki --keyring-backend test --chain-id odinchain

# wait for activation transaction success
sleep 2

for i in $(eval echo {1..1})
do
  # add reporter key
  yoda keys add reporter$i
done

# send odin tokens to reporters
echo "y" | odind tx bank send validator $(yoda keys list -a) 1000000loki --gas-prices 0.0025loki --keyring-backend test --chain-id odinchain

# wait for sending odin tokens transaction success
sleep 2

# add reporter to odinchain
echo "y" | odind tx oracle add-reporters $(yoda keys list -a) --from validator --gas-prices 0.0025loki --keyring-backend test --chain-id odinchain

# wait for addding reporter transaction success
sleep 2

# run yoda
yoda run
