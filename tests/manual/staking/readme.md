# This test validates the functionality of the staking command within the Titand command-line interface.

## Preconditions:
* Ensure you have a titand node running.
* You have installed the titand command client.

## Environments:
```
# Export gas price to send transaction
$ export fees=800000utkx

# Export home folder
$ export home=$(pwd)/local_test_data/.titan_val1

# Export the faucet's address
$ export faucet=$(titand keys show faucet --address --home=$home --keyring-backend=test)
```

## Verify command `staking`
### Steps:
1. In the command line, enter:
```
$ titand tx staking -h
```
### Expected results:
* You should see the output as follows:
```
Staking transaction subcommands

Usage:
  titand tx staking [flags]
  titand tx staking [command]

Available Commands:
  cancel-unbond    Cancel unbonding delegation and delegate back to the validator
  create-validator create new validator initialized with a self-delegation to it
  delegate         Delegate liquid tokens to a validator
  edit-validator   edit an existing validator account
  redelegate       Redelegate illiquid tokens from one validator to another
  unbond           Unbond shares from a validator
```

## Can create validator
### Steps:
1. Create and export the validator account:
```
$ titand keys add bob --home=$home --keyring-backend=test
$ export bob=$(titand keys show bob --address --home=$home --keyring-backend=test)
```
2. Ask for `1tkx` from the faucet:
```
$ titand tx bank send $faucet $bob 1tkx \
--gas=auto --fees=$fees \
--home=$home --keyring-backend=test -y
```
3. Stake `0.5tkx`:
```
$ titand tx staking create-validator \
--pubkey=$(titand tendermint show-validator --home=$home) \
--amount=0.5tkx \
--commission-rate="0.10" \
--commission-max-rate="0.20" \
--commission-max-change-rate="0.01" \
--min-self-delegation="1" \
--from=$bob --gas=auto --fees=$fees \
--home=$home --keyring-backend=test -y
```