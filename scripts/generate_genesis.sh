DIR=`dirname "$0"`

rm -rf ~/.odin

# initial new node
odind init validator --chain-id odinchain
echo "lock nasty suffer dirt dream fine fall deal curtain plate husodin sound tower mom crew crawl guard rack snake before fragile course bacon range" \
    | odind keys add validator --recover --keyring-backend test
echo "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic" \
    | odind keys add requester --recover --keyring-backend test


# add accounts to genesis
odind genesis add-genesis-account validator 10000000000000loki --keyring-backend test
odind genesis add-genesis-account requester 10000000000000loki --keyring-backend test


# register initial validators
odind genesis gentx validator 100000000loki \
    --chain-id odinchain \
    --keyring-backend test

# collect genesis transactions
odind genesis collect-gentxs

sed -i -e \
    "s/^minimum-gas-prices *=.*/minimum-gas-prices = \"0.0025loki\"/" \
    ~/.odin/config/app.toml

sed -i -e \
  '/\[api\]/,+10 s/enable = .*/enable = true/' \
  ~/.odin/config/app.toml

sed -i -e \
  '/\[mempool\]/,+10 s/version = .*/version = \"v1\"/' \
  ~/.odin/config/config.toml
