# Building and running the application

## Building the `odin` application

If you want to build the `odin` application in this repo to see the functionalities, **Go 1.18.0+** is required .

Add some parameters to environment is necessary if you have never used the `go mod` before.

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=\$PATH:\$GOBIN" >> ~/.bash_profile
echo "export GO111MODULE=on" >> ~/.bash_profile
source ~/.bash_profile
```

Now, you can install and run the application.

```bash
# Clone the source of the tutorial repository
git clone https://github.com/ODIN-PROTOCOL/odin-core
cd odin-core

# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands:
odind version
odind help
```

## Running test application locally

You can use the following script to generate a test environment to run BandChain locally. This will create the default genesis file with one validator, as well as some test accounts.

```bash
./scripts/generate_genesis.sh
cp ./docker-config/single-validator/priv_validator_key.json ~/.odin/config/priv_validator_key.json
cp ./docker-config/single-validator/node_key.json ~/.odin/config/node_key.json
```

Once done, you can optionally add data sources or oracle scripts to the genesis file using `odind`.

```bash
odind add-data-source ...
odind add-oracle-script ...
```

You can now start the chain with `odind`.

```bash
odind start
```

On a separate tab, you should run the oracle daemon script to ensure your validator responds to oracle requests.

```bash
./scripts/start_yoda.sh validator
```

To send an oracle request to the chain, use `odind`.

```bash
odind tx oracle request [ORACLE_SCRIPT_ID] [ASK_COUNT] [MIN_COUNT] -c [CALLDATA] --from requester --gas auto --keyring-backend test --from requester
```
