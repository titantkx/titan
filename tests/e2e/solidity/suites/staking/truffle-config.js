module.exports = {
  networks: {
    // Development network is just left as truffle's default settings
    titan: {
      host: "127.0.0.1",     // Localhost (default: none)
      port: 8545,            // Standard Ethereum port (default: none)
      // websockets: true,
      network_id: "*",       // Any network (default: none)
      gas: 7000000,          // Gas sent with each transaction
      gasPrice: 100000000000,  // 100 gwei (in wei)                
      skipDryRun: true,
      timeoutBlocks: 500,
      networkCheckTimeout: 10000000,
    },
  },
  mocha: {
    enableTimeouts: false,
    timeout: 1000 * 60 * 90,
  },
  compilers: {
    solc: {
      version: "0.5.17", // A version or constraint - Ex. "^0.5.0".
      settings: {
        optimizer: {
          enabled: true,
          runs: 10000,
        },
      },
    },
  },
}
