# Titan chain E2E tests

## Prerequisites

Required tools:

- jq
- docker
- docker compose

## Start testing

From the project root directory run:

```
go test github.com/tokenize-titan/titan/tests/e2e/cmd -v
```

To view blockchain logs during the test:

```
tail -f tests/e2e/cmd/titand.log
```
