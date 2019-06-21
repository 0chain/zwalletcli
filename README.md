# ZWallet Command-line Interface for 0Chain Blockchain
ZWallet Command-line utility is useful to quickly demonstrate and understand the capabilities of 0Chain Blockchain. The utility is built using 0Chain's ClientSDK library written in Go V1.12
## Features
ZWallet Command-line utility supports following features:
1. Create Wallets
2. Get test tokens from 0Chain Faucet
3. Send tokens between Wallets
4. Lock and Unlock tokens to earn interest
5. Recover wallet using passphrases

ZWallet Command-line utility provides a self-explaining "help" option that lists out the commands it supports and the parameters each command needs to perform the intended action
## How to get it?
You can clone ZWallet Command-line utility from github repo [Here](https://github.com/0chain/zwalletcli)
## Pre-requisites
* ZWallet Command-line utility needs Go V1.12 or higher. 
* [gosdk](https://github.com/0chain/gosdk)
## How to Build the code?
1. Make sure you've Go SDK 1.12 or higher and Go configurations are set and working on your system.
2.  Clone [gosdk](https://github.com/0chain/gosdk) and follow steps in Build and Installation section.
3. Clone [zwalletcli](https://github.com/0chain/zwalletcli)
4. Go to the root directory of the local repo
5. Run the following command:

        go build -tags bn256 -o zwallet

6. zwallet application is built in the local folder. 
## Getting started with ZWallet
### Before you start
Before you start playing with ZWallet, you need to know where the blockchain is running and what encryption scheme it is using. Both of that information is stored in a configuration files under clusters folder under repo. Choose the suitable one based on your needs.

### Setup
ZWallet Command-line Utility needs to know the configuration at runtime. By default, configuration files are assumed to be under $Home/.zcn folder. So, create $Home/.zcn folder and store the chosen yml files from clusters folder as nodes.yaml file there.

ZWallet can read configuration file from any customized folder and finlename too. You need to specify the customized folder with *--configDir* flag and/or the customized filename with *--config* flag in all the commands.

### Commands
To run the commands, cd to the folder where zwallet is located.

Let's go over all the available commands and play with it. It is assumed you are using default set up. If you're using customized set up,see ***customized set up*** example below for more details.

#### command with no arguments
When you run zwallet with no arguments, it will list all the supported commands.

Command

    ./zwallet
Response

   Use Zwallet to store, send and execute smart contract on 0Chain platform.
			Complete documentation is available at https://0chain.net

    Usage:
        zwallet [command]

    Available Commands:
        faucet          Faucet smart contract
        getbalance      Get balance from sharders
        getlockedtokens Get locked tokens
        help            Help about any command
        lock            Lock tokens
        lockconfig      Get lock configuration
        recoverwallet   Recover wallet
        send            Send ZCN token to another wallet
        unlock          Unlock tokens

    Flags:
        --config string      config file (default is nodes.yaml)
        --configDir string   configuration directory (default is $HOME/.zcn)
    -h, --help               help for zwallet
        --wallet string      wallet file (default is wallet.txt)

    Use "zwallet [command] --help" for more information about a command.

### help
To get the list of required arguments for a command use help flag
Command

    ./zwallet faucet --help

Response

    Faucet smart contract.
                <methodName> <input>

    Usage:
        zwallet faucet [flags]

    Flags:
    -h, --help                help for faucet
        --input string        input
        --methodName string   methodName

    Global Flags:
        --config string      config file (default is nodes.yaml)
        --configDir string   configuration directory (default is $HOME/.zcn)
        --wallet string      wallet file (default is wallet.txt)


#### getbalance
getbalance helps in two ways 

1. to get balance of an existing wallet
2. to create a wallet if there is none

Command

    ./zwallet getbalance
Response

    No wallet in path  $Home/.zcn/wallet.txt found. Creating wallet...
    ZCN wallet created!!

    Get balance failed. 
If you open the wallet.txt file, you will see the wallet details.

    {"client_id":"44347b5640ef3f5313e5efe3c6ab0e0c83efd625ed2bf00e912479aa8813cb1d","client_key":"1e400854be8bc1a787f4528da60984f22aee9e1fa47d3aa3aef27c40e8b087077283596084a9329e82d9f4f1eaf4319415648dd47795c5ed1156c2363dbe1280","keys":[{"public_key":"1e400854be8bc1a787f4528da60984f22aee9e1fa47d3aa3aef27c40e8b087077283596084a9329e82d9f4f1eaf4319415648dd47795c5ed1156c2363dbe1280","private_key":"b4ec0f105417a833d38213f2c246bd5c37a242e251009088e1e8f7204f112f0a"}],"mnemonics":"portion hockey where day drama flame stadium daughter mad salute easily exact wood peanut actual draw ethics dwarf poverty flag ladder hockey quote awesome","version":"1.0","date_created":"2019-06-16 16:22:15.406946 -0700 PDT m=+0.007561539"}
Out of these, the client_id and menmonics fields will be useful later.

#### customized set up
Let's have a customized set up with 
1. the configuration folder - create a folder "playground" under the root folder of repo.
2. the configuration file - copy the nodes.yaml file as devi.yaml and place it under the playground folder.

With this set up lets run "getbalance" again

Command
    ./zwallet --configDir ./playground --config devi getbalance

Response

    No wallet in path  ./playground/wallet.txt found. Creating wallet...

    ZCN wallet created!!

    Get balance failed. 
#### faucet
faucet command is useful to get test tokens into your wallet for transactional purposes.

Command

     ./zwallet faucet --methodName pour --input "{Pay day}"

Response

    Execute faucet smart contract success

#### getbalance
Let's use getbalance again to check the balance.
Command

    ./zwallet getbalance

Response

    Balance: 1 

There is 1 token deposited in the wallet as specified in wallet.txt. Same way you can use faucet any number of time whenever you need additional tokens.

#### getbalance
You can also use getbalance command to create a new wallet with a desired file name. In order to use that in any of the commands, you need to use the flag --wallet [ wallet_file_name ]

Command

     ./zwallet getbalance --wallet from

Response

    No wallet in path  /Users/jay_at_0chain/.zcn/from found. Creating wallet...
    ZCN wallet created!!

    Get balance failed. 

Check the new wallet file "from" created under your $Home/.zcn

#### send
Use send command to send a transaction from one wallet to the other. Send commands take four parameters. 
1. From wallet -- default is the account in wallet.txt
2. to_client_id -- address of the wallet receiving the funds
3. desc -- description for the transaction
4. tokens -- tokens in decimals to be transferred.

Command

     ./zwallet send --wallet from --desc "testing" --toclientID "7fe5e58d94684e3ec0b7fe076c4bc2aa56c455bfc7a476155c142d42eaf0d416" --token 0.5

Response

    Send token success

When you run a getbalance on both the wallets you see the difference

#### lockconfig
0Chain has a great way of earning additional tokens by locking tokens. When you lock tokens for a period of time, you will earn interest. The terms of lock can be obtained by lockconfig command.
Command

    ./zwallet lockconfig

Response

    Configuration:
    {"ID":"6dba10422e368813802877a85039d3985d96760ed844092319743fb3a76712d9","max_lock_period":31536000000000000,"min_lock_period":60000000000,"simple_global_node":{"interest_rate":0.5,"min_lock":10}}


#### lock
Command

./zwallet lock --wallet from --durationHr 0 --durationMin 5 --token 10.0

Response

    Tokens (10.000000) locked successfully

If you run the getbalance, you see that interest would have been already paid. Those additional tokens are yours to use. How cool is that!
#### getlockedtokens
Use getlockedtokens command to get informatiion about locked tokens

Command
    ./zwallet getlockedtokens --wallet from

Response

    Locked tokens:
    {"stats":[{"pool_id":"41fd52bbc848553365ae7b1319a3732764ea699964c3c97f1d85fb45fb46572e","start_time":"2019-06-17 05:48:54 +0000 UTC","duration":"5m0s","time_left":"3m57.17069839s","locked":true,"interest_rate":0.000004756468797564688,"interest_earned":475646,"balance":100000000000}]}

In the above response, make a note of pool_id. You need this when you want to unlock. Rest of the fields are self-explanatory.
#### unlock

Use this command to unlock the locked tokens. Unless you unlock, the tokens are not released.

Command

    ./zwallet unlock --poolid 41fd52bbc848553365ae7b1319a3732764ea699964c3c97f1d85fb45fb46572e
Response

    Unlock token success
#### recoverWallet

use this command to recover wallet from a different computer. You need to provide mnemonics mentioned in the wallet as an argument to prove that you own the wallet.

Command

    ./zwallet recoverwallet --mnemonic  "portion hockey where day drama flame stadium daughter mad salute easily exact wood peanut actual draw ethics dwarf poverty flag ladder hockey quote awesome"

### Tips
1. Sometimes when a transaction is sent, it may fail with a message "verify transaction failed". In such cases you need to resend the transactions
2. Use cmdlog.log to check possible reasons for failure of transactions.
3. zwalletcmd also comes with a Makefile which simplifies a lot of these zwallet commands.  