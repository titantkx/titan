#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=x/staking/types/expected_keepers.go -package testutil -destination x/staking/testutil/expected_keepers_mocks.go
