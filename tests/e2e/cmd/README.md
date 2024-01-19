# Titan chain E2E tests

## Prerequisites

- Install ignite CLI v0.27.2+

```
curl https://get.ignite.com/cli@v0.27.2! | bash
```

## Start testing

- From the project root directory run:

```
go test github.com/tokenize-titan/titan/tests/e2e/cmd -v
```

- To view blockchain logs during the test:

```
tail -f tests/e2e/cmd/titand.log
```
