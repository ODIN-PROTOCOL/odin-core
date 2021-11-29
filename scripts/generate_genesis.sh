DIR=`dirname "$0"`

rm -rf ~/.odin

# initial new node
odind init validator --chain-id odinchain
echo "lock nasty suffer dirt dream fine fall deal curtain plate husband sound tower mom crew crawl guard rack snake before fragile course bacon range" \
    | odind keys add validator --recover --keyring-backend test
echo "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic" \
    | odind keys add requester --recover --keyring-backend test

cp ./docker-config/single-validator/priv_validator_key.json ~/.odin/config/priv_validator_key.json
cp ./docker-config/single-validator/node_key.json ~/.odin/config/node_key.json

# add accounts to genesis
odind add-genesis-account validator 10000000000000odin --keyring-backend test
odind add-genesis-account requester 10000000000000odin --keyring-backend test


# register initial validators
odind gentx validator 100000000odin \
    --chain-id odinchain \
    --keyring-backend test

# collect genesis transactions
odind collect-gentxs


