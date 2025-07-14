require('dotenv').config({ path: process.env.NODE_ENV ? `.env.${process.env.NODE_ENV}` : '.env' });
console.log(`Current NODE_ENV=${process.env.NODE_ENV ? process.env.NODE_ENV : '!!!!!!MAIN!!!!!'}`);

import { HardhatUserConfig, task } from 'hardhat/config';

import '@nomicfoundation/hardhat-toolbox';
import '@openzeppelin/hardhat-upgrades';
import 'hardhat-gas-trackooor';
import 'hardhat-storage-layout';
import 'hardhat-tracer';
import 'solidity-coverage';
import 'tsconfig-paths/register';

const MILLISECOND_PER_SECOND = 1000;

if (process.env.NODE_ENV === 'test') {
}

// This is a sample Hardhat task. To learn how to create your own go to
// https://hardhat.org/guides/create-task.html
task('accounts', 'Prints the list of accounts', async (taskArgs, hre) => {
  const accounts = await hre.ethers.getSigners();

  for (const account of accounts) {
    console.log(account.address);
  }
});

const config: HardhatUserConfig = {
  solidity: {
    version: '0.8.16',
    settings: {
      optimizer: {
        enabled: true,
        runs: 1000,
      },
      outputSelection: {
        '*': {
          '*': ['storageLayout'],
        },
      },
    },
  },
  networks: {
    titan: {
      url: 'http://127.0.0.1:8545',
      gasMultiplier: 1.1,
      timeout: 600000 * MILLISECOND_PER_SECOND,
    },
  },
  gasReporter: {
    enabled: process.env.REPORT_GAS === 'TRUE',
    currency: 'USD',
    coinmarketcap: process.env.COIN_MARKET_CAP_KEY || undefined,
  },
};

export default config;
