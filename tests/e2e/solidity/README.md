# Solidity tests

Increasingly difficult tests are provided:

- [Basic](./suites/basic): simple Counter example, for basic calls, transactions, and events
- [Initialize](./suites/initialize): initialization contract and tests from [aragonOS](https://github.com/aragon/aragonOS)
- [Initialize (Buidler)](./suites/initialize-buidler): initialization contract and tests from [aragonOS](https://github.com/aragon/aragonOS), using [buidler](https://buidler.dev/)
- [Proxy](./suites/proxy): depositable delegate proxy contract and tests from [aragonOS](https://github.com/aragon/aragonOS)
- [Staking](./suites/staking): Staking contracts and full test suite from [aragon/staking](http://github.com/aragon/staking)

### Quick start

**Prerequisite**: in the repo's root, run `make install` to install the `titand` and `titand` binaries. When done, come back to this directory.

**Prerequisite**: install the individual solidity packages. They're set up as individual reops in a yarn monorepo workspace. Install them all via `yarn install`.

To run the tests, you can use the `test-helper.js` utility to test all suites under `ganache` or `titan` network. The `test-helper.js` will help you spawn an `titand` process before running the tests.

You can simply run `yarn test --network titan` to run all tests with ethermint network, or you can run `yarn test --network ganache` to use ganache shipped with truffle. In most cases, there two networks should produce identical test results.

If you only want to run a few test cases, append the name of tests following by the command line. For example, use `yarn test --network titan --allowTests=basic` to run the `basic` test under `titan` network.

If you need to take more control, you can also run `titand` using:

```sh
./init-test-node.sh
```

Keep the terminal window open, go into any of the tests and run `yarn test-titan`. You should see `titand` accepting transactions and producing blocks. You should be able to query for any transaction via:

- `titand query tx <cosmos-sdk tx>`
- `curl localhost:8545 -H "Content-Type:application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["<titan tx>"],"id":1}'`

From here, in your other available terminal,
And obviously more, via the Ethereum JSON-RPC API).

When in doubt, you can also run the tests against a Ganache instance via `yarn test-ganache`, to make sure they are behaving correctly.
