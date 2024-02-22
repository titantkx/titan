#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=x/staking/types/expected_keepers.go -package testutil -destination x/staking/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/gov/types/expected_keepers.go -package testutil -destination x/gov/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/validatorreward/types/expected_keepers.go -package testutil -destination x/validatorreward/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/nftmint/types/expected_keepers.go -package testutil -destination x/nftmint/testutil/expected_keepers_mocks.go
