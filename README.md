# zwallet - a CLI for 0chain Blockchain

zwallet is a command line interface (CLI) to quickly demonstrate and understand the capabilities of 0Chain Blockchain. The utility is built using 0Chain's goSDK library written in Go. There are some cool videos on how to lock tokens. Check out this [video](https://youtu.be/Eiz9mqdFtZo) on how to use the sample makefile for sending and receiving tokens to the 0Chain testnet, and another [video](https://youtu.be/g44VczBzmXo) on how to lock tokens and earn interest.

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

## How to get it?

    git clone https://github.com/0chain/zboxcli.git

## Pre-requisites

    Go V1.12 or higher.

### [How to build on Linux](https://github.com/0chain/zwalletcli/wiki/Build-Linux)

### [How to build on Windows](https://github.com/0chain/zwalletcli/wiki/Build-Windows)

## Getting started with zwallet

### Before you start

Before you start playing with zwalet, you need to access the blockchain. Go to network folder in the repo, and choose a network. Copy it to your ~/.zcn folder and then rename it as config.yaml file.

    mkdir ~/.zcn
    cp network/one.yaml ~/.zcn/config.yaml

Sample config.yaml

      ---
      block_worker: http://localhost:9091
      signature_scheme: bls0chain
      min_submit: 50 # in percentage
      min_confirmation: 50 # in percentage
      confirmation_chain_length: 3

We use blockWorker to connect to the network instead of giving network details directly, It will fetch the network details automatically from the blockWorker's network API.

#### Override Network Details

By default it will use the miner/sharder values which it will get using the `block_worker_url/network`. In case you want to override those values and give custom miner/sharder to use, You can create a `network.yaml` in your ~/.zcn (config) folder and paste the miner/sharder values in below format.

    miners:
      - http://localhost:7071
      - http://localhost:7072
      - http://localhost:7073
    sharders:
      - http://localhost:7171

Note: This is also useful for the Mac OS users running local cluster and having trouble with docker internal IPs (block_worker return docker IPs in local)

### Setup

The zwallet command line uses the ~/.zcn/config.yaml file at runtime to point to the network specified in that file.

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

### Tips 

1. Sometimes when a transaction is sent, it may fail with a message "verify transaction failed". In such cases you need to resend the transactions
2. Use cmdlog.log to check possible reasons for failure of transactions.
3. zwallet also comes with a Makefile which simplifies a lot of these zwalletcli commands.
4. Use `--verbose` flag on any of the command to see more information.