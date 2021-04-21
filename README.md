# zwallet - a CLI for 0chain wallet

`zwallet` is a command line interface (CLI) to demonstrate the functionalities of 0Chain. 

The CLI uses the [0chain Go SDK](https://github.com/0chain/gosdk) to do most of its functions.

## Architecture

`zwallet` is configurable to work with any 0chain network. It uses a config file and a wallet file stored locally on the filesystem. 

For most of its transactional commands, `zwallet` uses the `0dns` to discover the network nodes, then submits the transaction(s) requested to the miners, and finally waits the confirmation on the sharders.  

![alt text](docs/architecture.png "Architecture")

## Getting started

### 1. Installation

**Prerequisites**
- go 1.13

**Procedures**
1. Clone the `zwalletcli` repo and 

```sh
git clone https://github.com/0chain/zwalletcli.git
cd zwalletcli
```

2. Execute install
```sh
make install
```

3. Add config yaml at `~/.zcn/config.yaml`

The following will set `https://one.devnet-0chain.net` as your blockchain network.
```sh
cat > ~/.zcn/config.yaml << EOF
block_worker: http://one.devnet-0chain.net/dns
signature_scheme: bls0chain
min_submit: 50 # in percentage
min_confirmation: 50 # in percentage
confirmation_chain_length: 3
EOF
```

To know more about the config fields, details are found [here](#zcnconfigyaml).

4. Run `zwallet`
```sh
./zwallet
```

A list of `zwallet` commands should be displayed.

----
For detailed steps on the installation, follow any of the following:
- [How to build on Linux](https://github.com/0chain/zwalletcli/wiki/Build-Linux)
- [How to build on Windows](https://github.com/0chain/zwalletcli/wiki/Build-Windows)

### 2. Run `zwallet` commands

The following steps assume the current directory is inside the `zwalletcli` repo.

1. Register a new wallet

The wallet information is stored on `/.zcn/wallet.json`. Initially, there is no wallet yet.

When you execute any `zwallet` command, it will automatically create a wallet if it cannot find any, and save that locally.

Run the `ls-miners` command and see that it creates a wallet first before completing the command requested which is listing the miner nodes.
```sh
./zwallet ls-miners
```

Output
```
No wallet in path  <home dir>/.zcn/wallet.json found. Creating wallet...
ZCN wallet created!!
Creating related read pool for storage smart-contract...
Read pool created successfully
- ID:         cdb9b5a29cb5f48b350481694c4645c2db24500e3af210e22e2d10477a68bad2
- Host:       one.devnet-0chain.net
- Port:       31203
- ID:         3d9a10dac6fb3903d4a5283a42ae07b29d8e5d228afcce9bfc14e3e9dbc82748
- Host:       one.devnet-0chain.net
- Port:       31201
- ID:         aaa721d5fbf4ca83e20c8c40874ebcb144b86f57173633ff1702968677c2fa98
- Host:       one.devnet-0chain.net
- Port:       31202
```

2. Get some tokens

Faucet smart contract is available on dev networks and can be used to get tokens.

Run the `faucet` command to get 1 token.
```sh
./zwallet faucet --methodName pour --input "need token"
```
Output
```
Execute faucet smart contract success with txn :  915cfc6fa81eb3622c7082436a8ff752420e89dee16069a625d5206dc93ac3ca
```

3. Check wallet balance

Run the `getblance` command
```sh
./zwallet getbalance
```
Output
```
Balance: 1 (1.76 USD)
```

4. Lock tokens to get interest

Tokens can be locked into a pool to gain interest. It uses the Interest Pool smart contract.
   
Run the `lock` command and provide the amount of tokens and how long to lock them. 
```sh
./zwallet lock --tokens 0.5 --durationMin 5 
```
Output
```
Tokens (0.500000) locked successfully
```

Check balance right after and see that the locked tokens is deducted but has already gained interest.
```sh
./zwallet getbalance

Balance: 0.5000004743 (0.8800008347680001 USD)
```


That's it! You are now ready to use `zwallet`.

## Flags

`zwallet` accept global flags to override default configuration. The flags can be used in any command.

| Flag | Description | Default |
| ----- | ----------- | ---------- |
| `--config` | [Config file](#zcnconfigyaml) | `config.yaml` |
| `configDir` | Config directory | `~/.zcn` |
| `network` | [Network file](#zcnnetworkyaml) | `network.yaml` |
| `verbose` | Enable verbose logging | `false` |
| `wallet` | Wallet file | `wallet.json` |

## Features

1. Creating and restoring wallets
2. Explore network
3. Getting tokens from faucet
4. Sending tokens
5. Staking for interest
6. Staking on miners and sharders 
7. Vesting pool


### Creating and restoring wallets

#### Creating wallet

Simply run any `zwallet` command and it will create a wallet if none exist yet.

Here is a `faucet` command to create wallet at default location`~/.zcn/wallet.json`.

```sh
./zwallet faucet --methodName pour --input "new wallet"
```

Command to create a second wallet at `~/.zcn/new_wallet.json`
```sh
./zwallet faucet --methodName pour --input "new wallet" --wallet new_wallet.json
```

Verify second wallet
```sh
cat ~/.zcn/new_wallet.json
```

#### Recovering wallet

It is critical that the wallet is backed up particularly the mnemonics in order to recover it later on.

To recover wallet with mnemonic, use `recoverwallet` command. 

Example

```sh
zwallet recoverwallet --mnemonic "pull floor crop best weasel suit solid gown filter kitten loan absent noodle nation potato planet demise online ten affair rich panel rent sell" --wallet recovered_wallet.json
```

Verify recovered wallet
```sh
cat ~/.zcn/recovered_wallet.json
```

#### Creating multisig wallet

TODO

```
   createmswallet     create multisig wallet
```

### Exploring network

```
   getblobbers        Get registered blobbers from sharders
   getid              Get Miner or Sharder ID from its URL
   ls-miners          Get list of all active miners fro Miner SC
   ls-sharders        Get list of all active sharders fro Miner SC
```

### Getting and sending tokens

TODO 
```
   faucet             Faucet smart contract
   send               Send ZCN tokens to another wallet
   getbalance         Get balance from sharders
   verify             verify transaction
```

### Staking for interest

TODO

```
   getlockedtokens    Get locked tokens
   lock               Lock tokens
   lockconfig         Get lock configuration
   unlock             Unlock tokens
```

### Staking on miners and sharders

TODO

```
mn-config          Get miner SC global info.
mn-info            Get miner/sharder info from Miner SC.
mn-lock            Add miner/sharder stake.
mn-pool-info       Get miner/sharder pool info from Miner SC.
mn-unlock          Unlock miner/sharder stake.
mn-update-settings Change miner/sharder settings in Miner SC.
mn-user-info       Get miner/sharder user pools info from Miner SC.
```

### Vesting pool

TODO 

```
   vp-add             Add a vesting pool
   vp-config          Check out vesting pool configurations.
   vp-delete          Delete a vesting pool
   vp-info            Check out vesting pool information.
   vp-list            Check out vesting pools list.
   vp-stop            Stop vesting for one of destinations and unlock tokens not vested
   vp-trigger         Trigger a vesting pool work.
   vp-unlock          Unlock tokens of a vesting pool
```


## Config

### ~/.zcn/config.yaml

`~/.zcn/config.yaml` is a required `zwallet` config.

| Field | Description | Value type |
| ----- | ----------- | ---------- |
| `block_worker` | The URL to chain network DNS that provides the lists of miners and sharders | string |
| `signature_scheme` | The signature scheme used in the network. This would be `bls0chain` for most networks | string |
| `min_submit` | The desired minimum success ratio (in percent) to meet when submitting transactions to miners | integer |
| `min_confirmation` | The desired minimum success ratio (in percent) to meet when verifying transactions on sharders | integer |
| `confirmation_chain_length` | The desired chain length to meet when verifying transactions | integer |

### (Optional) ~/.zcn/network.yaml 

Network nodes are automatically discovered using the `block_worker` provided on `~/.zcn/config.yaml`.

To override/limit the nodes used on `zwallet`, create `~/.zcn/network.yaml` as shown below. 

```sh
cat > ~/.zcn/network.yaml << EOF
miners:
  - http://one.devnet-0chain.net:31201
  - http://one.devnet-0chain.net:31202
  - http://one.devnet-0chain.net:31203
sharders:
  - http://one.devnet-0chain.net:31101
EOF
```

Overriding the nodes can be useful in local chain setup. In some cases, the block worker might return URLs with IP/alias only accessible within the docker network.

## Video resources

TODO verify video content
- [Send and receive token](https://youtu.be/Eiz9mqdFtZo)
- [Lock tokens and earn interest](https://youtu.be/g44VczBzmXo)


## Troubleshooting

1. For more logging, add `--verbose` when running commands.
2. `zwallet getbalance` says it failed

This happens when the wallet has no token.

```sh
zwallet getbalance

Get balance failed.
```


---
# OLD README

## Features

[zwallet](#Command-with-no-arguments) supports following features:

1. [Create Wallet](#Command-with-no-arguments)
2. [Create multisig wallet](#Create-multisig-wallet)
3. [Get test tokens from 0Chain Faucet](#Faucet)
4. [Send tokens between Wallets](#Send)
5. [Lock tokens to earn interest](#Lock)
6. [Unlock locked tokens](#Unlock)
7. [Recover wallet using passphrase](#Recover)
8. [Get balance](#Get-balance)
9. [Get blobbers list](#Get-blobbers)
10. [Get miner or sharder ID](#Get-id)
11. [Get locked tokens](#Get-locked-tokens)
12. [Get lock configuration](#Get-lock-config)
13. [Verify transaction](#Verify)
14. [Vesting pool](#Vesting)
15. [Miner SC](#Miner-SC)

zwallet CLI provides a self-explaining "help" option that lists commands and parameters they need to perform the intended action

## Commands

Note in this document, we will show only the commands, response will vary depending on your usage, so may not be provided in all places.

### Command with no arguments

When you run zwallet with no arguments, it will list all the supported commands.
If you don't have a wallet yet, It will create one.

Command

    ./zwallet

Response

    Use zwallet to store, send and execute smart contract on 0Chain platform.
    Complete documentation is available at https://0chain.net

    Usage:
      zwallet [command]

    Available Commands:

      createmswallet     create multisig wallet
      faucet             Faucet smart contract
      getbalance         Get balance from sharders
      getblobbers        Get registered blobbers from sharders
      getid              Get Miner or Sharder ID from its URL
      getlockedtokens    Get locked tokens
      help               Help about any command
      lock               Lock tokens
      lockconfig         Get lock configuration
      ls-miners          Get list of all active miners fro Miner SC
      ls-sharders        Get list of all active sharders fro Miner SC
      mn-config          Get miner SC global info.
      mn-info            Get miner/sharder info from Miner SC.
      mn-lock            Add miner/sharder stake.
      mn-pool-info       Get miner/sharder pool info from Miner SC.
      mn-unlock          Unlock miner/sharder stake.
      mn-update-settings Change miner/sharder settings in Miner SC.
      mn-user-info       Get miner/sharder user pools info from Miner SC.
      recoverwallet      Recover wallet
      send               Send ZCN tokens to another wallet
      unlock             Unlock tokens
      verify             verify transaction
      version            Prints version information
      vp-add             Add a vesting pool
      vp-config          Check out vesting pool configurations.
      vp-delete          Delete a vesting pool
      vp-info            Check out vesting pool information.
      vp-list            Check out vesting pools list.
      vp-stop            Stop vesting for one of destinations and unlock tokens not vested
      vp-trigger         Trigger a vesting pool work.
      vp-unlock          Unlock tokens of a vesting pool



    Flags:
          --config string      config file (default is config.yaml)
          --configDir string   configuration directory (default is $HOME/.zcn)
      -h, --help               help for zwallet
          --verbose            prints sdk log in stderr (default false)
          --wallet string      wallet file (default is wallet.json)


    Use "zwallet [command] --help" for more information about a command.

### Create multisig wallet

Before jumping on to command description a quick introduction to Multisignature Wallet.

A Multisignature Wallet is a wallet for which any transaction from this wallet needs to be voted by T(N) associated signer wallets. To create a Multisignature Wallet, you need to specify the number of signers (N) you want on that wallet and minimum number of votes (T) it needs for a transaction to be approved.

APIs

- CreateMSWallet API will create the group wallet (MultiSignature Wallet) and corresponding number of Signer Wallets. All of these wallets have to be registered first on the Blockchain.

- RegisterMultiSig API registers the group wallet with MultiSig smartcontract.

- CreateVote API creates a vote for a proposal.

- RegisterVote API will vote for a proposal. Initially, if the proposal does not exist, MultiSig smartcontract will automatically create one as identified by the ProposalID parameter. Any vote bearing the same ProposalID, will be counted as a vote for the transaction. When the threshold number of votes are registered, transaction will be automatically processed. Any extra votes will be ignored.

Note 1: All Proposals will have an expiry of 7 days from the time of creation. At this point, it cannot be changed. Any vote coming after the expiry may create a new proposal.

Note 2: Before a transaction or voting can happen, the group wallet and the signer wallets have to be activated with one or more tokens.

Back to the command _createmswallet_. This command demonstrates how to create a multi-signature wallet, create a proposal for a transaction, and vote for the transaction. Note that this command works only on bls0chain encryption enabled 0chain Blockchain instance. The encryption scheme is specified by the "signature_scheme" field in the nodes.yml file under the configDir option.

Command

    ./zwallet createmswallet --numsigners 3 --threshold 2

where

1. numsigners is the number of accounts that can sign the vote.
2. threshold is the minimum number of votes required for the transaction to pass.
3. testn is an optional argument. set it to true to test sending votes from all signer accounts. By default votes from only threshold number of signer accounts is used.

Response

    Creating and testing a multisig wallet is successful!

### Faucet

faucet command is useful to get test tokens into your wallet for transactional purposes.

Command

     ./zwallet faucet --methodName pour --input "{Pay day}"

Response

    Execute faucet smart contract success

### Get balance

getbalance helps in two ways

- to get balance of an existing wallet
- to create a wallet if there is none

Command

    ./zwallet getbalance

Response

    No wallet in path  $Home/.zcn/wallet.txt found. Creating wallet...
    ZCN wallet created!!

If you do not have any balance, you get the following.

    Get balance failed.

Use [faucet](#Faucet) and try again

    Balance: 1

If you open the wallet.json file, you will see the wallet details.

    {"client_id":"7ea92e2e104489092f067b2644b22c5b1da001c0730f1fd7e990fd6b6bacaedd","client_key":"221fe6a29dbd09496556778aff010ff12dadc2fe0aec9c4d9e8a48c18cb00e13138a974d7728cc67597a6f92ba9355513eda4726e4fa1124b0ed5a8c9cc4490b","keys":[{"public_key":"221fe6a29dbd09496556778aff010ff12dadc2fe0aec9c4d9e8a48c18cb00e13138a974d7728cc67597a6f92ba9355513eda4726e4fa1124b0ed5a8c9cc4490b","private_key":"94b1e9b2adf1dd1c141797aaf586b60b144011983723a266e19d6d6aaf859e1b"}],"mnemonics":"night acid purpose slim junk wrist clown lyrics engine faint select capable swallow direct armor buzz student degree omit fiction favorite air volume learn","version":"1.0","date_created":"2020-03-06 13:52:24.763873623 +0530 IST m=+0.004795132"}

Out of these, the client_id and menmonics fields will be useful later.

### Send

Use send command to send a transaction from one wallet to the other. Send commands take four parameters.

- --from wallet -- default is the account in `~/.zcn/wallet.json`
- --to_client_id -- address of the wallet receiving the funds (this must be a registered wallet)
- --desc -- description for the transaction
- --tokens -- tokens in decimals to be transferred.

Command

     ./zwallet send --desc "testing" --to_client_id "7fe5e58d94684e3ec0b7fe076c4bc2aa56c455bfc7a476155c142d42eaf0d416" --tokens 0.5

Response

    Send tokens success

When you run a getbalance on both the wallets you see the difference

NOTE

   This command may return a success response in the case of multiple transactions from the same wallet falling within the same round. So to be sure that send transaction has committed, wait a few rounds (several seconds) and check balances and/or inspect finalized blockchain transactions. This is not regarded as an error, rather that the send request was successful but not actually transacted. This is an issue with the cli tools rather than the gosdk itself.

### Lock

Command

    ./zwallet lock --durationHr 0 --durationMin 5 --tokens 0.1

Response

    Tokens (0.100000) locked successfully

If you run the getbalance, you see that interest would have been already paid. Those additional tokens are yours to use. How cool is that!

### Get locked tokens

Use getlockedtokens command to get information about locked tokens

Command

    ./zwallet getlockedtokens

Response

    Locked tokens:
    {"stats":[{"pool_id":"41fd52bbc848553365ae7b1319a3732764ea699964c3c97f1d85fb45fb46572e","start_time":"2019-06-17 05:48:54 +0000 UTC","duration":"5m0s","time_left":"3m57.17069839s","locked":true,"interest_rate":0.000004756468797564688,"interest_earned":475646,"balance":100000000000}]}

In the above response, make a note of pool_id. You need this when you want to unlock. Rest of the fields are self-explanatory.

### Unlock

Use this command to unlock the locked tokens. Unless you unlock, the tokens are not released.

You can only unlock tokens once the lock duration has passed. The time left and lock status is available when running `getlockedtokens`.

Command

    ./zwallet unlock --pool_id 41fd52bbc848553365ae7b1319a3732764ea699964c3c97f1d85fb45fb46572e

Response

    Unlock tokens success

### Recover

Use this command to recover wallet from a different computer. You need to provide mnemonics mentioned in the wallet as an argument to prove that you own the wallet.

Command

    ./zwallet recoverwallet --mnemonic  "portion hockey where day drama flame stadium daughter mad salute easily exact wood peanut actual draw ethics dwarf poverty flag ladder hockey quote awesome"

### Get lock config

0Chain has a great way of earning additional tokens by locking tokens. When you lock tokens for a period of time, you will earn interest. The terms of lock can be obtained by lockconfig command.
Command

    ./zwallet lockconfig

Response

    Configuration:
    {"ID":"6dba10422e368813802877a85039d3985d96760ed844092319743fb3a76712d9","max_lock_period":31536000000000000,"min_lock_period":60000000000,"simple_global_node":{"interest_rate":0.5,"min_lock":10}}


### Get id

Use this command to get ID of a miner or sharder.

Command

    ./zwallet getid --url http://localhost:7071

Response

    URL: http://localhost:7071
    ID: 31810bd1258ae95955fb40c7ef72498a556d3587121376d9059119d280f34929

### Get blobbers

Use this command to get list of blobbers.
Command

    ./zwallet getblobbers

Response

    Blobbers:
           URL          |                                ID
    +-----------------------+------------------------------------------------------------------+
      http://localhost:5054 | 2a4d5a5c6c0976873f426128d2ff23a060ee715bccf0fd3ca5e987d57f25b78e
      http://localhost:5053 | 2f051ca6447d8712a020213672bece683dbd0d23a81fdf93ff273043a0764d18
      http://localhost:5052 | 7a90e6790bcd3d78422d7a230390edc102870fe58c15472073922024985b1c7d
      http://localhost:5051 | f65af5d64000c7cd2883f4910eb69086f9d6e6635c744e62afcfab58b938ee25


### Verify

Use this command to verify a transaction.
Command

    ./zwallet verify --hash f65af5d64000c7cd2883f4910eb69086f9d6e6635c744e62afcfab58b938ee25

Response

    Transaction status.
    Creating and testing a multisig wallet is successful!

### Vesting

#### Add a vesting pool

Create a vesting pool.

Flags

    - description, description for vesting pool, limited by SC configurations
      'max_description_length'
    - duration, vesting duration in form of [Golang duration](https://pkg.go.dev/time?tab=doc#ParseDuration),
      the value limited by SC configurations 'min_lock_period' and
      'max_lock_period'
    - lock, amount of tokens to lock in the pool, the provided amount should fit
      the amount of destinations; also, the amount limited by SC configurations
      'min_lock'
    - d, colon separated values consist of D:V, where D is vesting destination
      id, and V is value to be vested for the destination,; the flag can be
      repeated many times

Example

```
./zwallet vp-add                                                              \
    --description "for testing"                                               \
    --duration 5m                                                             \
    --lock 5                                                                  \
    --d 9fe14ab61ad7172f3cb9629fa34ca449229579ddf2d2a0fe3a58086344522d8e:1    \
    --d e7f451fdfe12a385045fedcb2e26d5ceb50f460c19e9a58e105dde17fc624588:2
```

Successful output contains pool ID for further requests.

#### Check out vesting pool configurations.

Get Vesting SC configurations.

Example

```
./zwallet vp-config
```

#### Delete a vesting pool

Delete a vesting pool. Stop vesting all destinations, unlock all the rest and
delete the pool.

Example

```
./zwallet vp-delete --pool_ID <pool_id>
```

#### Check out vesting pool information.

Information about a vesting pool for current moment.

Example

```
./zwallet vp-info --pool_id <pool_id>
```

#### Check out vesting pools list.

Get list of all vesting pools (IDs) for current client.

```
./zwallet vp-list
```

#### Stop vesting for one of destinations and unlock tokens not vested

Stop vesting for a destination and unlock the rest.

Example

```
./zwallet vp-stop --pool_id <pool_id> --d <destination_id>
```

#### Trigger a vesting pool work.

Trigger vesting for a vesting pool for current time. It moves all vested
tokens of the pool (for current time) to the destinations of the pool. Only
pool owner can trigger.

```
./zwallet vp-trigger --pool_id <pool_id>
```

#### Unlock tokens of a vesting pool

1. By pool owner. Unlock all tokens over required if any.
2. By a destination. Unlock tokens vested for the destination.

```
./zwallet vp-unlock --pool_id <pool_id>
```

### Miner SC

#### Get Miners list from Miner SC

    ./zwallet ls-miners

#### Get Sharders list from Miner SC

    ./zwallet ls-sharders

#### Get SC configurations and state

    ./zwallet mn-config

Response

    view_change:           7301
    max_n:                 8
    min_n:                 2
    max_s:                 3
    min_s:                 1
    t_percent:             0.51
    k_percent:             0.75
    last_round:            7369
    max_stake:             100
    min_stake:             0
    interest_rate:         0.001
    reward_rate:           1
    share_ratio:           0.1
    block_reward:          0.7
    max_charge:            0.5
    epoch:                 15000000
    reward_decline_rate:   0.1
    interest_decline_rate: 0.1
    max_mint:              4000000
    minted:                147.3653
    max_delegates:         200


#### Node information

Get miner/sharder information from Miner SC. 

    ./zwallet mn-info --id NODE_ID

Response

    {"simple_miner":{"id":"31810bd1258ae95955fb40c7ef72498a556d3587121376d9059119d280f34929","n2n_host":"198.18.0.71","host":"localhost","port":7071,"public_key":"255452b9f49ebb8c8b8fcec9f0bd8a4284e540be1286bd562578e7e59765e41a7aada04c9e2ad3e28f79aacb0f1be66715535a87983843fea81f23d8011e728b","short_name":"localhost.m0","build_tag":"2a366103470715432bcac43405ab722823a00c23","total_stake":1000000000,"delegate_wallet":"31810bd1258ae95955fb40c7ef72498a556d3587121376d9059119d280f34929","service_charge":0.1,"number_of_delegates":10,"min_stake":0,"max_stake":1000000000000,"stat":{"generator_rewards":3343620000000,"generator_fees":48},"node_type":"miner","last_health_check":1612439517},"active":{"215befba83bd2d4aaeddc89ae07f4205f322d2f0cce2829f9a9ff5a5fc5ece61":{"stats":{"delegate_id":"bf325bb5b978c32ab38226d1c26857cb78171f837e98f33f3b6ffccc5a6bb8c2","high":1330000000,"low":1000000,"interest_paid":10000000,"reward_paid":1208970000000,"number_rounds":0,"status":"ACTIVE"},"pool":{"pool":{"id":"215befba83bd2d4aaeddc89ae07f4205f322d2f0cce2829f9a9ff5a5fc5ece61","balance":1000000000},"lock":{"delete_view_change_set":false,"delete_after_view_change":0,"owner":"bf325bb5b978c32ab38226d1c26857cb78171f837e98f33f3b6ffccc5a6bb8c2"}}}}}


#### Lock stake for a node.

Lock stake for miner or sharder. 

    ./zwallet mn-lock --id NODE_ID

Response

    locked with: e5f87e4a82be6297c4a39caebff87c6258f3be861b8698b01f6fbf38d227fa6f

#### Check out stake pool info.

Get miner/sharder stake pool info from Miner SC.

    ./zwallet mn-pool-info --id NODE_ID --pool_id POOL_ID

Response

    {"stats":{"delegate_id":"bf325bb5b978c32ab38226d1c26857cb78171f837e98f33f3b6ffccc5a6bb8c2","high":1330000000,"low":1330000000,"interest_paid":0,"reward_paid":42560000000,"number_rounds":0,"status":"ACTIVE"},"pool":{"pool":{"id":"215befba83bd2d4aaeddc89ae07f4205f322d2f0cce2829f9a9ff5a5fc5ece61","balance":1000000000},"lock":{"delete_view_change_set":false,"delete_after_view_change":0,"owner":"bf325bb5b978c32ab38226d1c26857cb78171f837e98f33f3b6ffccc5a6bb8c2"}}}

#### Unlock a stake

Unlock miner/sharder stake pool. Tokens will be released next VC OR at next reward round.

    ./zwallet mn-unlock --id NODE_ID --pool_id POOL_ID

Response

    tokens will be unlocked next VC

#### Update node settings.

Change miner/sharder settings in Miner SC by delegate_wallet owner.

    ./zwallet mn-update-settings --id NODE_ID  [flags]

Flags are:

    --max_stake float     max stake allowed
    --min_stake float     min stake allowed
    --num_delegates int   max number of delegate pools

### User pools of Miner SC.

Get list of pools of Miner SC of a user.

    ./zwallet mn-user-info

Response

    - node: 31810bd1258ae95955fb40c7ef72498a556d3587121376d9059119d280f34929 (miner)
        - pool_id:        e5f87e4a82be6297c4a39caebff87c6258f3be861b8698b01f6fbf38d227fa6f
          balance:        0.1
          interests paid: 0
          rewards paid:   0.399
          status:         active
          stake %:        100 %

Optional flag `--client_id` can be used to get pools information for given
user. Current user used by default.

There is `--json` flag to print result as JSON.
