# Delele old volumns
rm -rf tmp/val1/.titand/*
rm -rf tmp/val2/.titand/*

# Init val1
docker compose -f docker-compose-genesis.yml run --rm -i val1 init val1 --chain-id=titan_18887-1 --overwrite
docker compose -f docker-compose-genesis.yml run --rm -i val1 config keyring-backend test
sed -i '' 's/^indexer = ".*"/indexer = "kv"/' tmp/val1/.titand/config/config.toml
sed -i '' 's/^timeout_commit = ".*"/timeout_commit = "1s"/' tmp/val1/.titand/config/config.toml

# Init val2
docker compose -f docker-compose-genesis.yml run --rm -i val2 init val2 --chain-id=titan_18887-1 --overwrite
docker compose -f docker-compose-genesis.yml run --rm -i val2 config keyring-backend test
sed -i '' 's/^indexer = ".*"/indexer = "kv"/' tmp/val2/.titand/config/config.toml
sed -i '' 's/^timeout_commit = ".*"/timeout_commit = "1s"/' tmp/val2/.titand/config/config.toml

### On val1 machine

# Copy from existing genesis file
cp genesis.json tmp/val1/.titand/config/genesis.json

# Config genesis file
config="
.chain_id = \"titan_18887-1\" |
.validators = [] |
.app_state.bank.supply = [] |
del(.app_state.bank.balances[] | select(.address == \"titan1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ljrm26\")) |
del(.app_state.ibc.connection_genesis.connections[] | select(.client_id == \"09-localhost\")) |
.app_state.staking.params.global_min_self_delegation = \"5000000000000000000\" |
.app_state.staking.last_total_power = \"0\" |
.app_state.staking.last_validator_powers = [] |
.app_state.staking.validators = [] |
.app_state.staking.delegations = [] |
.app_state.staking.unbonding_delegations = [] |
.app_state.staking.redelegations = [] |
.app_state.distribution.delegator_withdraw_infos = [] |
.app_state.distribution.previous_proposer = \"\" |
.app_state.distribution.outstanding_rewards = [] |
.app_state.distribution.validator_accumulated_commissions = [] |
.app_state.distribution.validator_historical_rewards = [] |
.app_state.distribution.validator_current_rewards = [] |
.app_state.distribution.delegator_starting_infos = [] |
.app_state.distribution.validator_slash_events = [] |
.app_state.gov.params.min_deposit[0].denom = \"atkx\" |
.app_state.gov.params.min_deposit[0].amount = \"250000000000000000000\" |
.app_state.gov.params.max_deposit_period = \"15s\" |
.app_state.gov.params.voting_period = \"30s\"
"
echo $(jq "$config" tmp/val1/.titand/config/genesis.json) > tmp/val1/.titand/config/genesis.json

# Add faucet account
docker compose -f docker-compose-genesis.yml run --rm -i val1 keys add faucet
faucet=$(docker compose -f docker-compose-genesis.yml run --rm -i val1 keys show faucet --address)
# Add balance to faucet
docker compose -f docker-compose-genesis.yml run --rm -i val1 add-genesis-account $faucet 100000000tkx

# Add val1 account
docker compose -f docker-compose-genesis.yml run --rm -i val1 keys add val1
val1=$(docker compose -f docker-compose-genesis.yml run --rm -i val1 keys show val1 --address)
# Add balance to val1
docker compose -f docker-compose-genesis.yml run --rm -i val1 add-genesis-account $val1 1000000tkx
# val1 stakes tkx
docker compose -f docker-compose-genesis.yml run --rm -i val1 gentx val1 100000tkx --min-self-delegation 5000000000000000000

# Copy genesis file from val1 machine to val2 machine
cp tmp/val1/.titand/config/genesis.json tmp/val2/.titand/config/genesis.json

### On val2 machine

# Add val2 account
docker compose -f docker-compose-genesis.yml run --rm -i val2 keys add val2
val2=$(docker compose -f docker-compose-genesis.yml run --rm -i val2 keys show val2 --address)
# Add balance to val2
docker compose -f docker-compose-genesis.yml run --rm -i val2 add-genesis-account $val2 1000000tkx
# val2 stakes tkx
docker compose -f docker-compose-genesis.yml run --rm -i val2 gentx val2 100000tkx --min-self-delegation 5000000000000000000

# Copy val2 key to val1 machine
cp tmp/val2/.titand/keyring-test/* tmp/val1/.titand/keyring-test

# Copy generated txs and genesis file from val2 machine to val1 machine
cp tmp/val2/.titand/config/gentx/gentx-* tmp/val1/.titand/config/gentx
cp tmp/val2/.titand/config/genesis.json tmp/val1/.titand/config/genesis.json

### On val1 machine

# Collect all generated transactions into genesis file
docker compose -f docker-compose-genesis.yml run --rm -i val1 collect-gentxs
# Validate the genesis file
docker compose -f docker-compose-genesis.yml run --rm -i val1 validate-genesis

# Copy final genesis file from val1 machine to val2 machine
cp tmp/val1/.titand/config/genesis.json tmp/val2/.titand/config/genesis.json

# Expose rpc endpoint
sed -i '' 's/^laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/' tmp/val1/.titand/config/config.toml
