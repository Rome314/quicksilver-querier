# Quicksilver Querier

Quicksilver Querier is a simple CLI app that allows you to retrieve data from a QuickSilver node and write it into a CSV file. It is written in Go and uses the Cosmos SDK for interacting with the QuickSilver node.

## Features

The app can query the following data from a QuickSilver node:

- Validator to delegator mapping
- Delegator to validator mapping
- Vesting accounts details categorized by type (DelayedVestingAccount, PeriodicVestingAccount, PermanentLockedAccount, and PeriodicVestingAccount)
- IBC channels between two specified chains for their STATUS
- Client state for the channels
- All "pending" receipts in the x/interchainstaking module

The output data is structured in a CSV format that is easy to read and analyze.

## Installation

1. Clone the repository to your local machine:

```bash
git clone https://github.com/Rome314/quicksilver-querier.git
cd quickdump
make build
````

## Usage

The app accepts the following parameters:

- `node`: The URL of the node to connect to (default: quicksilver.grpc.kjnodes.com:11190)
- `format`: The output format (csv)
- `output`: The path where to store the response

The app also accepts the following task names:

- `pending-staking-receipts`
- `channels-statuses`
- `vesting-accounts`
- `validators-delegators`

## Example

### Pending Staking Receipts
To get all pending staking receipts, run:

```bash
./quickdump pending-staking-receipts --node <node_url> --format <output_format> --output <output_file>
```
### Channels Statuses
To get the status of all IBC channels, run:

```bash
./quickdump channels-statuses --node <node_url> --format <output_format> --output <output_file>
```
### Vesting Accounts
To get details of all vesting accounts, run:

```bash
./quickdump vesting-accounts --node <node_url> --format <output_format> --output <output_file>
```
### Validators Delegators
To get the mapping of validators to delegators and vice versa, run:

```bash
./quickdump validators-delegators --node <node_url> --format <output_format> --output <output_file>
```

## TODO: 
- [ ] Add endpoint checking by chains-registry
- [ ] Write tests
- [ ] Write benchmarks for query pagination to find optimal batch sizes

## Contributing
Contributions are welcome. Please submit a pull request or create an issue for any enhancements, bugs, or feature requests.
