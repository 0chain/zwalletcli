# zwallet - a CLI for 0chain wallet

`zwallet` is a command line interface (CLI) to demonstrate the functionalities of 0Chain.

The CLI utilizes the [0chain Go SDK](https://github.com/0chain/gosdk).

- [zwallet - a CLI for 0chain wallet](#zwallet---a-cli-for-0chain-wallet)
  - [Architecture](#architecture)
  - [Getting started](#getting-started)
    - [1. Installation](#1-installation)
    - [2. Run `zwallet` commands](#2-run-zwallet-commands)
  - [Global parameters](#global-parameters)
  - [Commands](#commands)
    - [Creating and restoring wallets](#creating-and-restoring-wallets)
      - [Creating wallet - (any command)](#creating-wallet---any-command)
      - [Recovering wallet - `recoverwallet`](#recovering-wallet---recoverwallet)
    - [Exploring network nodes](#exploring-network-nodes)
      - [Listing all miners - `ls-miners`](#listing-all-miners---ls-miners)
      - [Listing all sharders -`ls-sharders`](#listing-all-sharders--ls-sharders)
      - [Listing all blobbers - `getblobbers`](#listing-all-blobbers---getblobbers)
      - [Getting node ID by URL - `getid`](#getting-node-id-by-url---getid)
    - [Getting and sending tokens](#getting-and-sending-tokens)
      - [Getting tokens with Faucet smart contract - `faucet`](#getting-tokens-with-faucet-smart-contract---faucet)
      - [Checking balance - `getbalance`](#checking-balance---getbalance)
      - [Sending tokens to another wallet - `send`](#sending-tokens-to-another-wallet---send)
      - [Verifying a transaction - `verify`](#verifying-a-transaction---verify)
    - [Staking on miners and sharders](#staking-on-miners-and-sharders)
      - [Getting the staking config - `mn-config`](#getting-the-staking-config---mn-config)
      - [Getting a miner or sharder info for staking - `mn-info`](#getting-a-miner-or-sharder-info-for-staking---mn-info)
      - [Locking a stake on a node - `mn-lock`](#locking-a-stake-on-a-node---mn-lock)
    - [Getting the stake pools of a wallet - `mn-user-info`](#getting-the-stake-pools-of-a-wallet---mn-user-info)
      - [Getting the stake pool info - `mn-pool-info`](#getting-the-stake-pool-info---mn-pool-info)
      - [Unlock a stake - `mn-unlock`](#unlock-a-stake---mn-unlock)
      - [Updating staking config of a node - `mn-update-settings`](#updating-staking-config-of-a-node---mn-update-settings)
  - [Config](#config)
    - [~/.zcn/config.yaml](#zcnconfigyaml)
    - [(Optional) ~/.zcn/network.yaml](#optional-zcnnetworkyaml)
 
## Architecture

`zwallet` can be configured to work with any 0chain network. It uses a config and a wallet file stored on the local filesystem.

For most transactions, `zwallet` uses the `0dns` to discover the network nodes, then creates and submits transaction(s) to the miners, and finally waits for transaction confirmation on the sharders.

![architecture](docs/architecture.png "Architecture")

## Getting started

### 1. Installation

**Prerequisites**

- Go: Installation instructions for Mac, Linux and Windows can be found [here](https://go.dev/doc/install).

**Procedures**

1. Clone the `zwalletcli` repo and install

```sh
git clone https://github.com/0chain/zwalletcli.git
cd zwalletcli
make install
```

2. Add config yaml at `~/.zcn/config.yaml`

The following script sets `https://demo.zus.network` as your network.

```sh
cat > ~/.zcn/config.yaml << EOF
block_worker: https://demo.zus.network/dns
signature_scheme: bls0chain
min_submit: 50 # in percentage
min_confirmation: 50 # in percentage
confirmation_chain_length: 3
EOF
```

To understand more about the config properties, head over [here](https://github.com/0chain/zwalletcli/blob/staging/network/config.yaml).

3. Run `zwallet` to display the list of supported commands.

```sh
./zwallet
```
----
For detailed steps on the installation, follow any of the following:

- [How to build on Linux/Mac](https://github.com/0chain/zwalletcli/wiki/Build-on-Linux-and-Mac)
- [How to build on Windows](https://github.com/0chain/zwalletcli/wiki/Build-Windows)

### 2. Run `zwallet` commands

The following steps assume that your terminal's working directory is inside the `zwalletcli` repo.

## Global parameters

`zwallet` accept global parameters to override default configuration and can be used in any command.

| Parameter     | Description                     | Default        |
| ------------- | ------------------------------- | -------------- |
| `--help`      | Show help                       |                |
| `--config`    | [Config file]()   | `config.yaml`  |
| `--configDir` | Config directory                | `~/.zcn`       |
| `--network`   | [Network file](#zcnnetworkyaml) | `network.yaml` |
| `--verbose`   | Enable verbose logging          | `false`        |
| `--wallet`    | Wallet file                     | `wallet.json`  |

## Commands

### Creating and restoring wallets

#### Creating wallet - (any command)

Simply run  `create-wallet` command and it will create a wallet if none exist yet.

![create wallet](docs/createwallet.png "Create wallet")

Here is a sample with `create-wallet` command and this creates a wallet at default location`~/.zcn/wallet.json`

```sh
./zwallet create-wallet
```
Sample Output 

```
wallet saved in /home/ubuntu/.zcn/wallet.json
{"client_id":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
"client_key":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
"keys":[{"public_key":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
"private_key":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}],
"mnemonics":"xxxx xxxx xxxx xxxxx",
"version":"1.0","date_created":"2023-05-03T12:44:46+05:30","nonce":0}

```
Here is a sample with `faucet` command and this creates a wallet at default location`~/.zcn/wallet.json`

```sh
./zwallet faucet --methodName pour --input "new wallet"
```

Another `faucet` command to create a second wallet at `~/.zcn/new_wallet.json`

```sh
./zwallet faucet --methodName pour --input "new wallet" --wallet new_wallet.json
```

Sample Output
```
No wallet in path  <home directory>/.zcn/new_wallet.json found. Creating wallet...
ZCN wallet created!!
Creating related read pool for storage smart-contract...
Read pool created successfully
```

Verify second wallet

```sh
cat ~/.zcn/new_wallet.json
```

#### Recovering wallet - `recoverwallet`

`recoverwallet` is used when restoring a lost wallet or when loading the wallet from a different computer.

Given a wallet's mnemonics, you can recover and recreate your wallet.

| Parameter    | Required | Description                            | Default | Valid Values |
| ------------ | -------- | -------------------------------------- | ------- | ------------ |
| `--mnemonic` | Yes      | The mnemonics of the wallet to recover |         |              |
| `--offline`  | No       | recover wallet without registration on blockchain|  false  |              |

![recover wallet](docs/recoverwallet.png "Recover wallet")

Sample command

```sh
./zwallet recoverwallet --wallet recovered_wallet.json --mnemonic "pull floor crop best weasel suit solid gown filter kitten loan absent noodle nation potato planet demise online ten affair rich panel rent sell"
```

Sample output

```
Wallet recovered!!
```

Verify recovered wallet

```sh
cat ~/.zcn/recovered_wallet.json
```
### Exploring network nodes

#### Listing all miners - `ls-miners`

The list of miners are retrieved using the Miner smart contract.

| Parameter | Required | Description          | Default | Valid Values      |
| --------- | -------- | -------------------- | ------- | ----------------- |
| `--json`  | No       | Print output as JSON |         | <empty to enable> |

![List miner nodes](docs/ls-miners.png "List miner nodes")

```sh
./zwallet ls-miners
```

Sample output

```
- ID:         cdb9b5a29cb5f48b350481694c4645c2db24500e3af210e22e2d10477a68bad2
- Host:       demo1.zus.network
- Port:       31203
- ID:         3d9a10dac6fb3903d4a5283a42ae07b29d8e5d228afcce9bfc14e3e9dbc82748
- Host:       demo2.zus.network
- Port:       31201
- ID:         aaa721d5fbf4ca83e20c8c40874ebcb144b86f57173633ff1702968677c2fa98
- Host:       demo3.zus.network
- Port:       31202
```

#### Listing all sharders -`ls-sharders`

The list of sharders are retrieved using the latest finalized magic block. All registered sharders can be retrieved with the `--all` parameter.

| Parameter | Required | Description                             | Default | Valid Values      |
| --------- | -------- | --------------------------------------- | ------- | ----------------- |
| `--json`  | No       | Print output as JSON                    |         | <empty to enable> |
| `--all`   | No       | Print also registered nodes on Miner SC |         | <empty to enable> |

![List sharder nodes](docs/ls-sharders.png "List sharder nodes")

```sh
./zwallet ls-sharders --all
```

Sample output

```
MagicBlock Sharders
ID: 12e317e5d7a4a0a914ec26074e28f00502c735ddf7ac7d156b34e83e39792a9d
  - N2NHost: demo1.zus.network
  - Host: demo1.zus.network
  - Port: 31102
ID: 675502b613ba1c5985636e3e92b9a857855a52155e3316bb40fe9607e14167fb
  - N2NHost: demo2.zus.network
  - Host: demo1.zus.network
  - Port: 31101

Registered Sharders
ID: 675502b613ba1c5985636e3e92b9a857855a52155e3316bb40fe9607e14167fb
  - N2NHost: demo1.zus.network
  - Host: demo1.zus.network
  - Port: 31101
ID: 12e317e5d7a4a0a914ec26074e28f00502c735ddf7ac7d156b34e83e39792a9d
  - N2NHost: demo2.zus.network
  - Host: demo2.zus.network
  - Port: 31102
ID: 43f4f011698db6f2078e6ceb1cd981ab3bd35d07b7ac6fdf7c77aec1feee09be
  - N2NHost: demo3.zus.network
  - Host: demo3.zus.network
  - Port: 31101
ID: fd02f4436692bd9f679fae809f4f140fd4daaa35769ae9c6db1ab9664f766c22
  - N2NHost: 144.76.91.241
  - Host: demo4.zus.network
  - Port: 31102

```

#### Listing all blobbers - `getblobbers`

The list of blobbers are retrieved using the Storage smart contract.

![List blobber nodes](docs/getblobbers.png "List blobber nodes")

```sh
./zwallet getblobbers
```

Sample output

```
Blobbers:
                 URL                 |                                ID                                |          CAP           |     R / W PRICE     | DEMAND
+------------------------------------+------------------------------------------------------------------+------------------------+---------------------+--------+
  http://demo1.zus.network:31305 | 011a5444c9fe53137da7c3d871f3bf4bbf3c01607b14aa95ff2de43b5537d4b6 | 271.5 GiB / 1000.0 GiB | 0.010000 / 0.010000 |    0.1
  http://demo1.zus.network:31306 | 2efc85d6a2f36380e1e77b843cd9f4fe55668271cae4925ab38a92504176e5df | 107.8 GiB / 1000.0 GiB | 0.010000 / 0.010000 |    0.1
  http://demo1.zus.network:31302 | 34934babf0781c21736023ff89bc554928d77c028a968ef7344a460611d5a8d2 | 104.3 GiB / 1000.0 GiB | 0.010000 / 0.010000 |    0.1
```

#### Getting node ID by URL - `getid`

Print the ID of a blockchain node.

| Parameter | Required | Description                               | Default | Valid Values |
| --------- | -------- | ----------------------------------------- | ------- | ------------ |
| `--url`   | Yes      | URL to the node (miner, sharder, blobber) |         |              |

![Get node ID](docs/getid.png "Get node ID")

The following command get the details of the sharder on a given URL

```sh
./zwallet getid --url http://demo1.zus.network:31101
```

Output

```
URL: http://demo1.zus.network:31101
ID: 675502b613ba1c5985636e3e92b9a857855a52155e3316bb40fe9607e14167fb
```

### Getting and sending tokens

[Video walkthrough](https://youtu.be/Eiz9mqdFtZo)

#### Getting tokens with Faucet smart contract - `faucet`

Tokens can be retrieved and added to your wallet through the Faucet smart contract.

| Parameter      | Required | Description                                                  | Default | Valid Values     |
| -------------- | -------- | ------------------------------------------------------------ | ------- | ---------------- |
| `--methodName` | Yes      | Smart Contract method to call (`pour` - get tokens, `refill` - return tokens) |         | `pour`, `refill` |
| `--input`      | Yes      | Request description                                          |         | any string       |
| `--tokens`     | No       | Amount of tokens (maximum of 1.0)                            | 1.0     | (0 - 1.0]        |

![Faucet](docs/faucet.png "Faucet")

The following command will give 1 token to the default wallet.

```sh
./zwallet faucet --methodName pour --input "need token"
```

The following command will return 0.5 token to faucet.

```sh
./zwallet faucet --methodName refill --input "not using" --tokens 0.5
```

Sample output from `faucet` prints the transaction

```
Execute faucet smart contract success with txn :  d25acd4a339f38a9ce4d1fa91b287302fab713ef4385522e16d18fd147b2ebaf
```

#### Checking balance - `getbalance`

Wallet balances are retrieved from sharders.

> Note: Balance would not show any [locked tokens](#locking-tokens-for-interest---lock).

![Get wallet balance](docs/getbalance.png "Get wallet balance")

```sh
./zwallet getbalance
```

Sample output

```
Balance: 3 (4.2299999999999995 USD)
```

To check the balance of another wallet, use `--wallet` global parameter.

```sh
./zwallet getbalance --wallet another_wallet.json
```

> Note: When there is no token on the wallet yet, output will show `Get balance failed.`

#### Sending tokens to another wallet - `send`

![Send tokens to another wallet](docs/send.png "Send tokens to another wallet")

Transferring tokens from a wallet to another is done through `send`

| Parameter        | Required | Description                    | Default | Valid Values |
| ---------------- | -------- | ------------------------------ | ------- | ------------ |
| `--to_client_id` | Yes      | Client ID of the recipient     |         |              |
| `--tokens`       | Yes      | Amount of tokens to send       |         | valid number |
| `--desc`         | Yes      | Transfer description           |         | any string   |
| `--fee`          | No       | Amount of tokens to use as fee | 0.0     | valid number |

The following sends 0.2 token from the default wallet to the specified client ID.

```sh
./zwallet send --to_client_id e7ebb698213b6bda097c0a14ccbe574356e99e9b666e4baeae540da1d9b51e7e --tokens .2 --desc "gift"
```

Output

```
Send tokens success
```

To use a different wallet as sender, use `--wallet` global parameter.

```sh
./zwallet send --to_client_id e7ebb698213b6bda097c0a14ccbe574356e99e9b666e4baeae540da1d9b51e7e --tokens .2 --desc "gift" --wallet another_wallet.json
```

#### Verifying a transaction - `verify`

`verify` checks whether a given transaction hash was confirmed on the blockchain.

| Parameter | Required | Description                   | Default | Valid Values           |
| --------- | -------- | ----------------------------- | ------- | ---------------------- |
| `--hash`  | Yes      | Hash of transaction to verify |         | valid transaction hash |

Note: Not all `zwallet` commands (eg. `send`) prints the transaction hash created. To see more details printed including the hashes, use `--verbose` global parameter.

![Verify transaction confirmation](docs/verify.png "Verify transaction confirmation")

Sample command

```sh
./zwallet verify --hash 867c240b640e3d128643330af383cb3a0a229ebce08cae667edd7766c7ccc850
```

Output

```
Transaction verification success
```

To see more details about the transaction on `verify`, use `--verbose` global parameter.

```sh
./zwallet verify --hash 867c240b640e3d128643330af383cb3a0a229ebce08cae667edd7766c7ccc850 --verbose
```

### Staking on miners and sharders

[Miner smart contract](https://github.com/0chain/0chain/blob/master/code/go/0chain.net/smartcontract/minersc/READEME.md) allows staking on the miner and sharder nodes.

The maximum number of stake pools per node is limited to the number of delegates allowed. To find out the number of delegates and the minimum and maximum tokens allowed, query the staking config.

#### Getting the staking config - `mn-config`

`mn-config` display the global info of Miner SC for staking.

![Miner SC global info](docs/mn-config.png "Miner SC global info")

```sh
./zwallet mn-config
```

Sample output

```
view_change:           0
max_n:                 8
min_n:                 2
max_s:                 3
min_s:                 1
t_percent:             0.51
k_percent:             0.75
last_round:            5527098
max_stake:             100
min_stake:             0
interest_rate:         5.55e-07
reward_rate:           1
share_ratio:           0.8
block_reward:          0.21
max_charge:            0.5
epoch:                 15000000
reward_decline_rate:   0.1
interest_decline_rate: 0.1
max_mint:              4000000
minted:                928552.464
max_delegates:         200
```

#### Getting a miner or sharder info for staking - `mn-info`

Node stats for staking are retrieved from Miner SC.

| Parameter | Required | Description                                                  | Default | Valid Values |
| --------- | -------- | ------------------------------------------------------------ | ------- | ------------ |
| `--id`    | Yes      | Node ID of a miner or sharder (get at `ls-miners` or `ls-sharders`) |         |              |

![Node stat for staking](docs/mn-info.png "Node stat for staking")

Sample command

```sh
./zwallet mn-info --id 68ed8f16e1d50e3217425b3e98fb7a39e5d7201fe4b1dccfe8477b5c54761b24
```

Sample output

```json
{"simple_miner":{"id":"68ed8f16e1d50e3217425b3e98fb7a39e5d7201fe4b1dccfe8477b5c54761b24","n2n_host":"five.devnet-0chain.net","host":"five.devnet-0chain.net","port":31102,"path":"sharder02","public_key":"458170c28496333426f9866989f7b335564d13d00f50db503275e7ec50a008031764ba5df42c20c85d76e7fe7eda43f39afdb650b8ffa9ed340e6fc50c35ae22","short_name":"localhost.s02","build_tag":"50fb047353c4c3d268c0e0ebfd8e63f1d10c6813","total_stake":0,"delegate_wallet":"68ed8f16e1d50e3217425b3e98fb7a39e5d7201fe4b1dccfe8477b5c54761b24","service_charge":0.1,"number_of_delegates":10,"min_stake":0,"max_stake":1000000000000,"stat":{"sharder_rewards":1160553450000000},"node_type":"sharder","last_health_check":1619554377},"pending":{"01978379a586de2882638345e215baaf8382093609d910da5ac1a833e2814f6f":{"stats":{"delegate_id":"133807913c66ec0b4342612f23fecd1852b456152433b8380cd2abcd411d4c07","high":0,"low":-1,"interest_paid":0,"reward_paid":0,"number_rounds":0,"status":"PENDING"},"pool":{"pool":{"id":"01978379a586de2882638345e215baaf8382093609d910da5ac1a833e2814f6f","balance":10000000000},"lock":{"delete_view_change_set":false,"delete_after_view_change":0,"owner":"ff12c78ee4a985b4fc2ac52ec8a24e9df2bd912636da15437b0eb7707b99abf4"}}}}}
```

Reformatted output

```json
{
  "simple_miner":{
    "id":"68ed8f16e1d50e3217425b3e98fb7a39e5d7201fe4b1dccfe8477b5c54761b24",
    "n2n_host":"demo1.zus.network",
    "host":"demo1.zus.network",
    "port":31102,
    "path":"sharder02",
    "public_key":"458170c28496333426f9866989f7b335564d13d00f50db503275e7ec50a008031764ba5df42c20c85d76e7fe7eda43f39afdb650b8ffa9ed340e6fc50c35ae22",
    "short_name":"localhost.s02",
    "build_tag":"50fb047353c4c3d268c0e0ebfd8e63f1d10c6813",
    "total_stake":0,
    "delegate_wallet":"68ed8f16e1d50e3217425b3e98fb7a39e5d7201fe4b1dccfe8477b5c54761b24",
    "service_charge":0.1,
    "number_of_delegates":10,
    "min_stake":0,
    "max_stake":1000000000000,
    "stat":{
      "sharder_rewards":1160553450000000
    },
    "node_type":"sharder",
    "last_health_check":1619554377
  },
  "pending":{
    "01978379a586de2882638345e215baaf8382093609d910da5ac1a833e2814f6f":{
      "stats":{
        "delegate_id":"133807913c66ec0b4342612f23fecd1852b456152433b8380cd2abcd411d4c07",
        "high":0,
        "low":-1,
        "interest_paid":0,
        "reward_paid":0,
        "number_rounds":0,
        "status":"PENDING"
      },
      "pool":{
        "pool":{
          "id":"01978379a586de2882638345e215baaf8382093609d910da5ac1a833e2814f6f",
          "balance":10000000000
        },
        "lock":{
          "delete_view_change_set":false,
          "delete_after_view_change":0,
          "owner":"ff12c78ee4a985b4fc2ac52ec8a24e9df2bd912636da15437b0eb7707b99abf4"
        }
      }
    }
  }
}
```

#### Locking a stake on a node - `mn-lock`

Staking tokens on a node gains additional tokens over time. Tokens locked for staking can be unlocked anytime although have to wait for the next view change cycle.

Note however that if a node becomes offline, all stake pools are automatically unlocked and tokens are returned to wallets.

| Parameter  | Required | Description                                                  | Default | Valid Values |
| ---------- | -------- | ------------------------------------------------------------ | ------- | ------------ |
| `--id`     | Yes      | Node ID of a miner or sharder to stake for (get at `ls-miners` or `ls-sharders`) |         |              |
| `--tokens` | Yes      | Amounts of token to stake                                    |         | valid number |

![Staking tokens on node](docs/mn-lock.png "Staking tokens on node")

Sample command

```sh
./zwallet mn-lock --id dc8c6c93fb42e7f6d1c0f93baf66cc77e52725f79c3428a37da28e294aa2319a --tokens 1
```

The output would print the stake pool id.

```sh
locked with: b488738546d84aed9d3dcb2bbe24c161bc4338638669e64e814631efd430fd85
```

If the locking of stakes is failing, verify the following.

1. Wallet has enough tokens
2. Node ID is valid
3. Node has available delegate

### Getting the stake pools of a wallet - `mn-user-info`

| Parameter     | Required | Description             | Default                        | Valid Values      |
| ------------- | -------- | ----------------------- | ------------------------------ | ----------------- |
| `--client_id` | No       | Client ID of the wallet | Wallet at `~/.zcn/wallet.json` |                   |
| `--json`      | No       | Print output as JSON    |                                | <empty to enable> |

![Wallet stake pools](docs/mn-user-info.png "Wallet stake pools")

Sample command

```sh
./zwallet mn-user-info
```

Sample output

```
- node: dc8c6c93fb42e7f6d1c0f93baf66cc77e52725f79c3428a37da28e294aa2319a (sharder)
  - pool_id:        b488738546d84aed9d3dcb2bbe24c161bc4338638669e64e814631efd430fd85
    balance:        1
    interests paid: 0
    rewards paid:   0
    status:         pending
    stake %:        100 %
```

Sample reformatted JSON output

```json
{
  "pools":{
    "sharder":{
      "dc8c6c93fb42e7f6d1c0f93baf66cc77e52725f79c3428a37da28e294aa2319a":[
        {
          "id":"b488738546d84aed9d3dcb2bbe24c161bc4338638669e64e814631efd430fd85",
          "balance":10000000000,
          "interest_paid":0,
          "reward_paid":0,
          "status":"PENDING",
          "high":0,
          "low":-1
        }
      ]
    }
  }
}
```

#### Getting the stake pool info - `mn-pool-info`

| Parameter   | Required | Description                                                  | Default | Valid Values |
| ----------- | -------- | ------------------------------------------------------------ | ------- | ------------ |
| `--id`      | Yes      | Node ID of a miner or sharder (get at `ls-miners` or `ls-sharders`) |         |              |

![Stake pool info](docs/mn-pool-info.png "Stake pool info")

Sample command

```sh
./zwallet mn-pool-info --id dc8c6c93fb42e7f6d1c0f93baf66cc77e52725f79c3428a37da28e294aa2319a
```

Sample output

```json
{"pool_id":"4db8e916a7fa73f79097adf11d47bd535ac5fe3566e26b0e71d205c0171181df","balance":1000000000,"reward":152712000000,"status":0,"round_created":632,"delegate_id":"c542815b92849bae60032761e71eab2ba0c703589252cabbc38d48b5846e38bf","staked_at":1658695083}
```

Reformatted output

```json
{
  "pool_id":"4db8e916a7fa73f79097adf11d47bd535ac5fe3566e26b0e71d205c0171181df",
  "balance":1000000000,
  "reward":152712000000,
  "status":0,
  "round_created":632,
  "delegate_id":"c542815b92849bae60032761e71eab2ba0c703589252cabbc38d48b5846e38bf",
  "staked_at":1658695083
}
```

If the `pool-id` flag is not present, a list of pools will be returned.
Sample command

```sh
./zwallet mn-pool-info --id dc8c6c93fb42e7f6d1c0f93baf66cc77e52725f79c3428a37da28e294aa2319a
```

Sample output

```json
[{"pool_id":"4db8e916a7fa73f79097adf11d47bd535ac5fe3566e26b0e71d205c0171181df","balance":1000000000,"reward":12096000000,"status":1,"round_created":632,"delegate_id":"c542815b92849bae60032761e71eab2ba0c703589252cabbc38d48b5846e38bf","staked_at":1658695083}]
```

Reformatted output

```json
[
  {
    "pool_id":"4db8e916a7fa73f79097adf11d47bd535ac5fe3566e26b0e71d205c0171181df",
    "balance":1000000000,
    "reward":12096000000,
    "status":1,
    "round_created":632,
    "delegate_id":"c542815b92849bae60032761e71eab2ba0c703589252cabbc38d48b5846e38bf",
    "staked_at":1658695083
  }
]
```

#### Unlock a stake - `mn-unlock`

| Parameter   | Required | Description                                                  | Default | Valid Values |
| ----------- | -------- | ------------------------------------------------------------ | ------- | ------------ |
| `--id`      | Yes      | Node ID of a miner or sharder to unlock stakes from (get at `mn-user-info`) |         |              |

![Unlock a stake](docs/mn-unlock.png "Unlock a stake")

```sh
./zwallet mn-unlock --id dc8c6c93fb42e7f6d1c0f93baf66cc77e52725f79c3428a37da28e294aa2319a
```

Output

```
tokens will be unlocked next VC
```

Tokens are released on the next view change cycle or at the next reward round.


#### Updating staking config of a node - `mn-update-settings`

Staking config can only be updated by the node's delegate wallet.

| Parameter         | Required | Description                                   | Default | Valid Values |
| ----------------- | -------- | --------------------------------------------- | ------- | ------------ |
| `--id`            | Yes      | Node ID of a miner or sharder                 |         |              |
| `--max_stake`     | No       | Minimum amount of tokens allowed when staking |         | valid number |
| `--min_stake`     | No       | Maximum amount of tokens allowed when staking |         | valid number |
| `--num_delegates` | No       | Maximum number of staking pools               |         | valid number |

![Update node settings for staking](docs/mn-update-settings.png "Update node settings for staking")

Sample command

```sh
./zwallet mn-update-settings --id dc8c6c93fb42e7f6d1c0f93baf66cc77e52725f79c3428a37da28e294aa2319a --max_stake 1000000000000 --min_stake 10000000 --num_delegates 25
```
## Config

### ~/.zcn/config.yaml

`~/.zcn/config.yaml` is a required `zwallet` config.

| Field                       | Description                                                  | Value type |
| --------------------------- | ------------------------------------------------------------ | ---------- |
| `block_worker`              | The URL to chain network DNS that provides the lists of miners and sharders | string     |
| `signature_scheme`          | The signature scheme used in the network. This would be `bls0chain` for most networks | string     |
| `min_submit`                | The desired minimum success ratio (in percent) to meet when submitting transactions to miners | integer    |
| `min_confirmation`          | The desired minimum success ratio (in percent) to meet when verifying transactions on sharders | integer    |
| `confirmation_chain_length` | The desired chain length to meet when verifying transactions | integer    |

### (Optional) ~/.zcn/network.yaml

Network nodes are automatically discovered using the `block_worker` provided on `~/.zcn/config.yaml`.

To override/limit the nodes used on `zwallet`, create `~/.zcn/network.yaml` as shown below.

```sh
cat > ~/.zcn/network.yaml << EOF
miners:
  - http://demo1.zus.network:31201
  - http://demo1.zus.network:31202
  - http://demo1.zus.network:31203
sharders:
  - http://demo1.zus.network:31101
EOF
```

Overriding the nodes can be useful in local chain setup. In some cases, the block worker might return URLs with IP/alias only accessible within the docker network.
