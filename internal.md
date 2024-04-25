## Documnenation for the internal commands

#### Get Authorizer Configuration - `bridge-auth-config`
`./zwallet bridge-auth-config `command can be used to view authorizer configuration. Here are the parameters for the command.

| Parameter | Required | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| --id      | Yes      | Provide Authorizer ID to view its configuration . |
| --help    |          | Syntax Help for the command                       |

Sample Command:

```
./zwallet bridge-auth-config --id $AUTHORIZER_ID
```

Sample Response:

```
{
  "id": "2f945f7310689f17afd8c8cb291e1e3ba21677243aa1d404a2293064e7983d60",
  "url": "https://demo.zus.network/authorizer01/",
  "fee": 0,
  "latitude": 0,
  "longitude": 0,
  "last_health_check": 0,
  "delegate_wallet": "",
  "min_stake": 0,
  "max_stake": 0,
  "num_delegates": 0,
  "service_charge": 0
}
```

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

#### Updating staking config of a node - `mn-update-settings`

Staking config can only be updated by the node's delegate wallet.

| Parameter         | Required | Description                                   | Default | Valid Values |
| ----------------- | -------- | --------------------------------------------- | ------- | ------------ |
| `--id`            | Yes      | Node ID of a miner or sharder                 |         |              |
| `--max_stake`     | No       | Minimum amount of tokens allowed when staking |         | valid number |
| `--min_stake`     | No       | Maximum amount of tokens allowed when staking |         | valid number |
| `--num_delegates` | No       | Maximum number of staking pools               |         | valid number |
| `--service_charge`     | No       | Service Charge |         | valid number |
| `--sharder` | No       | Whether node is sharder or not               |   False      | set true for sharder node else <empty> or false |


![Update node settings for staking](docs/mn-update-settings.png "Update node settings for staking")

Sample command

```sh
./zwallet mn-update-settings --id dc8c6c93fb42e7f6d1c0f93baf66cc77e52725f79c3428a37da28e294aa2319a --max_stake 1000000000000 --min_stake 10000000 --num_delegates 25
```