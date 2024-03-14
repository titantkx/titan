# Titan chain E2E tests

## Prerequisites

Required tools:

- jq
- docker
- docker compose

## Start testing

From the project root directory run:

```
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

To view blockchain logs during the test:

```
tail -f tests/e2e/cmd/titand.log
```
