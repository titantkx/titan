#!/bin/sh

set -e

# Detect platform
platform=$(uname)
if [ "$platform" = "Darwin" ]; then
    SED_INPLACE="sed -i ''"
else
    SED_INPLACE="sed -i"
fi

# Delele old volumns
rm -rf tmp/val1/.titand/*
rm -rf tmp/val2/.titand/*

# Init val1
docker compose run --rm -i val1 init val1 --chain-id=titan_18887-1 --overwrite
docker compose run --rm -i val1 config keyring-backend test
$SED_INPLACE 's/^indexer = ".*"/indexer = "kv"/' tmp/val1/.titand/config/config.toml
$SED_INPLACE 's/^timeout_commit = ".*"/timeout_commit = "0.5s"/' tmp/val1/.titand/config/config.toml
# Init val2
docker compose run --rm -i val2 init val2 --chain-id=titan_18887-1 --overwrite
docker compose run --rm -i val2 config keyring-backend test
$SED_INPLACE 's/^indexer = ".*"/indexer = "kv"/' tmp/val2/.titand/config/config.toml
$SED_INPLACE 's/^timeout_commit = ".*"/timeout_commit = "0.5s"/' tmp/val2/.titand/config/config.toml

### On val1 machine

# Add faucet account
docker compose run --rm -i val1 keys add faucet
faucet=$(docker compose run --rm -i val1 keys show faucet --address)
# Add balance to faucet
docker compose run --rm -i val1 add-genesis-account $faucet 97999000tkx

# Add val1 account
docker compose run --rm -i val1 keys add val1
val1=$(docker compose run --rm -i val1 keys show val1 --address)
# Add balance to val1
docker compose run --rm -i val1 add-genesis-account $val1 1000000tkx
# val1 stakes tkx
docker compose run --rm -i val1 gentx val1 100000tkx --min-self-delegation 5000000000000000000

# Add reward-pool-admin account
docker compose run --rm -i val1 keys add reward-pool-admin
reward_pool_admin=$(docker compose run --rm -i val1 keys show reward-pool-admin --address)
# Add balance to reward-pool-admin
docker compose run --rm -i val1 add-genesis-account $reward_pool_admin 1000tkx

# Config genesis file
config="
.app_state.validatorreward.params.authority = \"$reward_pool_admin\" |
.app_state.validatorreward.params.rate = \"0.190000000000000000\" |
.app_state.staking.params.bond_denom = \"atkx\" |
.app_state.staking.params.unbonding_time = \"120s\" |
.app_state.staking.params.max_validators = 8 |
.app_state.staking.params.global_min_self_delegation = \"5000000000000000000\" |
.app_state.crisis.constant_fee.denom = \"atkx\" |
.app_state.crisis.constant_fee.amount = \"10000000000000000000\" |
.app_state.feemarket.params.base_fee = \"100000000000\" |
.app_state.feemarket.params.min_gas_price = \"100000000000.000000000000000000\" |
.app_state.slashing.params.signed_blocks_window = \"100\" |
.app_state.slashing.params.min_signed_per_window = \"0.500000000000000000\" |
.app_state.slashing.params.downtime_jail_duration = \"60s\" |
.app_state.slashing.params.slash_fraction_double_sign = \"0.050000000000000000\" |
.app_state.slashing.params.slash_fraction_downtime = \"0.000100000000000000\" |
.app_state.gov.params.min_deposit[0].denom = \"atkx\" |
.app_state.gov.params.min_deposit[0].amount = \"250000000000000000000\" |
.app_state.gov.params.max_deposit_period = \"15s\" |
.app_state.gov.params.voting_period = \"30s\" |
.app_state.evm.params.evm_denom = \"atkx\" |
.app_state.bank.denom_metadata[0].name = \"titan tkx\" |
.app_state.bank.denom_metadata[0].symbol = \"TKX\" |
.app_state.bank.denom_metadata[0].display = \"tkx\" |
.app_state.bank.denom_metadata[0].base = \"atkx\" |
.app_state.bank.denom_metadata[0].description = \"The native token of the Titan network.\" |
.app_state.bank.denom_metadata[0].denom_units[0].denom = \"atkx\" |
.app_state.bank.denom_metadata[0].denom_units[0].exponent = 0 |
.app_state.bank.denom_metadata[0].denom_units[0].aliases[0] = \"attotkx\" |
.app_state.bank.denom_metadata[0].denom_units[1].denom = \"utkx\" |
.app_state.bank.denom_metadata[0].denom_units[1].exponent = 12 |
.app_state.bank.denom_metadata[0].denom_units[1].aliases[0] = \"microtkx\" |
.app_state.bank.denom_metadata[0].denom_units[2].denom = \"mtkx\" |
.app_state.bank.denom_metadata[0].denom_units[2].exponent = 15 |
.app_state.bank.denom_metadata[0].denom_units[2].aliases[0] = \"millitkx\" |
.app_state.bank.denom_metadata[0].denom_units[3].denom = \"tkx\" |
.app_state.bank.denom_metadata[0].denom_units[3].exponent = 18
"
echo "$(jq "$config" tmp/val1/.titand/config/genesis.json)" >tmp/val1/.titand/config/genesis.json

# Copy genesis file from val1 machine to val2 machine
cp tmp/val1/.titand/config/genesis.json tmp/val2/.titand/config/genesis.json

### On val2 machine

# Add val2 account
docker compose run --rm -i val2 keys add val2
val2=$(docker compose run --rm -i val2 keys show val2 --address)
# Add balance to val2
docker compose run --rm -i val2 add-genesis-account $val2 1000000tkx
# val2 stakes tkx
docker compose run --rm -i val2 gentx val2 100000tkx --min-self-delegation 5000000000000000000

# Copy val2 key to val1 machine
cp tmp/val2/.titand/keyring-test/* tmp/val1/.titand/keyring-test

# Copy generated txs and genesis file from val2 machine to val1 machine
cp tmp/val2/.titand/config/gentx/gentx-* tmp/val1/.titand/config/gentx
cp tmp/val2/.titand/config/genesis.json tmp/val1/.titand/config/genesis.json

### On val1 machine

# Collect all generated transactions into genesis file
docker compose run --rm -i val1 collect-gentxs
# Validate the genesis file
docker compose run --rm -i val1 validate-genesis

# Copy final genesis file from val1 machine to val2 machine
cp tmp/val1/.titand/config/genesis.json tmp/val2/.titand/config/genesis.json

# Add val2 node to seed peers
val2id=$(docker compose run --rm -i val2 tendermint show-node-id)
$SED_INPLACE "s/^seeds = \"\"/seeds = \"$val2id@val2:26656\"/" tmp/val1/.titand/config/config.toml

# Expose rpc endpoint
$SED_INPLACE 's/^laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/' tmp/val1/.titand/config/config.toml

### On val2 machine

# Add val1 node to seed peers
val1id=$(docker compose run --rm -i val1 tendermint show-node-id)
$SED_INPLACE "s/^seeds = \"\"/seeds = \"$val1id@val1:26656\"/" tmp/val2/.titand/config/config.toml
