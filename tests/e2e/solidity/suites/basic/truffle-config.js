module.exports = {
  networks: {
    // Development network is just left as truffle's default settings
    titan: {
      host: "127.0.0.1",     // Localhost (default: none)
      port: 8545,            // Standard Ethereum port (default: none)
      network_id: "*",       // Any network (default: none)
      gas: 5000000,          // Gas sent with each transaction
      gasPrice: 100000000000,  // 100 gwei (in wei)                
    },
  },
  compilers: {
    solc: {
      version: "0.5.17",
    },
  },
}
