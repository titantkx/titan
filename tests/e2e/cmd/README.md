# Titan chain E2E tests

## Prerequisites

Required tools:

- jq
- docker
- docker compose

## Start testing

From the project root directory run:

```shell
# Basic test
export TEST_TYPE=basic
go test github.com/titantkx/titan/tests/e2e/cmd -v

# Upgrade test
export TEST_TYPE=upgrade
go test github.com/titantkx/titan/tests/e2e/cmd -v

# Upgrade test from an exported genesis file
export TEST_TYPE=upgrade-from-genesis
go test github.com/titantkx/titan/tests/e2e/cmd -v
```

### Upgrade test

1. Update `UpgradeName` with new upgrade name in `tests/e2e/cmd/setup/upgrade/setup.go`.

2. Update `image` in `tests/e2e/cmd/setup/upgrade/docker-compose-genesis.yml` is old titand version.

3. Run test

    ```shell
    make test-e2e-upgrade
    ```

### Upgrade test from an exported genesis file

1. Update `UpgradeName` with new upgrade name in `tests/e2e/cmd/setup/upgrade-from-genesis/setup.go`.

2. Update `image` in `tests/e2e/cmd/setup/upgrade-from-genesis/docker-compose-genesis.yml` is old titand version.

3. To run the upgrade test from an exported genesis file, you need to export the genesis file from the current chain state first.

    ```shell
    titand export --for-zero-height > genesis.json
    ```

4. Run test

    ```shell
    make test-e2e-upgrade-from-genesis
    ```

## To view blockchain logs during the test

```shell
tail -f tests/e2e/cmd/titand.log
```
