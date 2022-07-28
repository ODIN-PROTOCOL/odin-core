<p>&nbsp;</p>
<p align="center">

<img src="./odinprotocol_logo.png" width="500px" alt="odin logo">

</p>

<p align="center">OdinChain - Decentralized Data Delivery Network.<br/><br/>

<a href="https://pkg.go.dev/badge/github.com/ODIN-PROTOCOL/odin-core">
    <img src="https://pkg.go.dev/badge/github.com/ODIN-PROTOCOL/odin-core">
</a>
<a href="https://goreportcard.com/badge/github.com/ODIN-PROTOCOL/odin-core">
    <img src="https://goreportcard.com/badge/github.com/ODIN-PROTOCOL/odin-core">
</a>
<a href="https://github.com/ODIN-PROTOCOL/odin-core/workflows/Tests/badge.svg">
    <img src="https://github.com/ODIN-PROTOCOL/odin-core/workflows/Tests/badge.svg">
</a>

<p align="center">
  <a href="https://app.gitbook.com/@geodb/s/odin-protocol/"><strong>Documentation »</strong></a>
  <br />
  <br/>
  <a href="https://odinprotocol.io/docs/odin-whitepaper.pdf">Whitepaper</a> | 
  <a href="https://odinprotocol.io/docs/odin-tokenomics.pdf">Tokenomics paper</a>
</p>

<br/>

_Current TestNet name is "**havi** - one of the names of supreme god Odin."_ <br>
_Name:_ **odin-testnet-havi-1**

## Installation

### Binaries

You can find the latest binaries on our [releases](https://github.com/ODIN-PROTOCOL/odin-core/releases) page.

### Building from source

To install OdinChain's daemon `odind`, you need to have [Go](https://golang.org/) (version 1.18.0 or later)
and [gcc](https://gcc.gnu.org/) installed on our machine. Navigate to the Golang
project [download page](https://golang.org/dl/) and gcc [install page](https://gcc.gnu.org/install/), respectively for
install and setup instructions.

## Running a Validator Node on the OdinChain TestNet

The following steps shows how to set up a validator node on the odinchain testnet. For similar instructions on running a
validator node on our testnet, please refer
to [this article](https://medium.com/odinprotocol/odinchain-guanyu-testnet-3-successful-upgrade-how-to-join-as-a-validator-2766ca6717d4)

We recommend the following for running a odinChain Validator:

- **2 or more** CPU cores
- **8 GB **of RAM
- At least **256GB** of disk storage

## Setting Up Validator Node

### Downloading the binaries

We will be assuming that you will be running your node on a Ubuntu 18.04 LTS machine that is allowing connections to
port 26656.

To start, you’ll need to install the various utility tools and Golang on the machine.

```bash
sudo apt-get update
sudo apt-get upgrade -y
sudo apt-get install -y build-essential curl wget

wget https://go.dev/dl/go1.18.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz
rm go1.18.1.linux-amd64.tar.gz

echo "export PATH=\$PATH:/usr/local/go/bin:~/go/bin" >> $HOME/.profile
source ~/.profile
```

### Build OdinChain Daemon

Next, you will need to clone and build OdinChain. The canonical version for this GuanYu Mainnet is v1.2.6.

```bash
git clone https://github.com/ODIN-PROTOCOL/odin-core
cd odin-core
git checkout testnet-{name}
make install

# Check that the correction version of odind is installed
odind version --long
```

### Creating OdinChain Account and Setup Config

Once installed, you can use the `odind` CLI to create a new OdinChain wallet address and initialize the chain. Please make sure to keep your mnemonic safe!

```bash
# Create a new odin wallet. Do not lose your mnemonic!
odind keys add <YOUR_WALLET>

# Initialize a blockchain environment for generating genesis transaction.
odind init --chain-id odin-testnet-havi-1 <YOUR_MONIKER>
```

You can then download the official genesis file from the repository. You should also add the initial peer nodes to your
Tendermint configuration file.

```bash
# Download genesis file from the repository.
wget https://raw.githubusercontent.com/ODIN-PROTOCOL/odin-core/master/testnets/odin-testnet-{name}/genesis.json
# Check genesis hash
sudo apt-get install jq
# Move the genesis file to the proper location
mv genesis.json $HOME/.odin/config
# Add some persistent peers
sed -E -i \
  's/persistent_peers = \".*\"/persistent_peers = \"492a4e30c10194e1d8f6fa194ba3f63b1aa73484@35.195.4.110:26656,417c2df701780c7f8751bc4a298411374082ef9e@34.78.138.110:26656,ea43cac04a556d01050a09a5699c3ba272a91116@34.78.239.23:26656,4edb332575e5108b131f0a7c0d9ac237569634ad@34.77.171.169:26656' \
  $HOME/.odin/config/config.toml
```

### Starting the Blockchain Daemon

With all configurations ready, you can start your blockchain node with a single command. In this tutorial, however, we
will show you a simple way to set up `systemd` to run the node daemon with auto-restart.

- Create a config file, using the contents below, at `/etc/systemd/system/odind.service`. You will need to edit the
  default ubuntu username to reflect your machine’s username. Note that you may need to use sudo as it lives in a
  protected folder

```
[Unit]
Description=odinChain Node Daemon
After=network-online.target
[Service]
User=ubuntu
ExecStart=/home/ubuntu/go/bin/odind start
Restart=always
RestartSec=3
LimitNOFILE=4096
[Install]
WantedBy=multi-user.target
```

- Install the service and start the node

```
sudo systemctl enable odind
sudo systemctl start odind
```

While not required, it is recommended that you run your validator node behind your sentry nodes for DDOS mitigation.
See [this thread](https://forum.cosmos.network/t/sentry-node-architecture-overview/454) for some example setups. Your
node will now start connecting to other nodes and syncing the blockchain state.

### ⚠️ Wait Until Your Chain is Fully Sync

You can tail the log output with `journalctl -u odind.service -f`. If all goes well, you should see that the node daemon
has started syncing. Now you should wait until your node has caught up with the most recent block.

```bash
... odind: I[..] Executed block  ... module=state height=20000 ...
... odind: I[..] Committed state ... module=state height=20000 ...
... odind: I[..] Executed block  ... module=state height=20001 ...
... odind: I[..] Committed state ... module=state height=20001 ...
```

⚠️ **NOTE:** You should not proceed to the next step until your node caught up to the latest block.

### Send Yourself odin Token

With everything ready, you will need some odin tokens to apply as a validator. You can use `odind` keys list command to
show your address.

```bash
odind keys list
- name: ...
  type: local
  address: odin1g3fd6rslryv498tjqmmjcnq5dlr0r6udm2rxjk
  pubkey: ...
  mnemonic: ""
  threshold: 0
  pubkeys: []
```

### Apply to Become Block Validator

Once you have some odin tokens, you can apply to become a validator by sending `MsgCreateValidator` transaction.

```bash
odind tx staking create-validator \
    --amount <your-amount-to-stake>loki \
    --commission-max-change-rate 0.01 \
    --commission-max-rate 0.2 \
    --commission-rate 0.1 \
    --from <your-wallet-name> \
    --min-self-delegation 1 \
    --moniker <your-moniker> \
    --pubkey $(odind tendermint show-validator) \
    --chain-id odin-testnet-havi-1
```

Once the transaction is mined, you should see yourself on the [validator page](https://testnet.odinprotocol.io/validators).
Congratulations. You are now a working OdinChain testnet validator!

### Setting Up Yoda — The Oracle Daemon

For Phase 1, OdinChain validators are also responsible for responding to oracle data requests. Whenever someone submits
a request message to OdinChain, the chain will autonomously choose a subset of active oracle validators to perform the
data query.

The validators are chosen submit a report message to OdinChain within a given timeframe as specified by a chain
parameter. We provide a program called yoda to do this task for you.

Yoda uses an external executor to resolve requests to data sources. Currently, it
supports [AWS Lambda](https://aws.amazon.com/lambda/) (through the REST interface).

In future releases, `yoda` will support more executors and allow you to specify multiple executors to add redundancy.
Please use [this link](https://github.com/odinprotocol/odinchain/wiki/AWS-lambda-executor-setup) to setup lambda
function.

## Resources

- Peers:
    - Testnet:
        - [Odin Testnet](https://node.testnet.odinprotocol.io)

## Community

- [Official Website](https://odinprotocol.io)

## License & Contributing

...
