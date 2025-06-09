import { HardhatEthersHelpers } from '@nomiclabs/hardhat-ethers/types';
import { bech32 } from 'bech32';
import { execSync } from 'child_process';
import { BigNumber } from 'ethers';
import prompts from 'prompts';

const ROOT_DIR = `${__dirname}` + '/' + '../../../../../../';
const MAX_BLOCK_INTERVAL = 5000; // 5 second

async function confirm(question: string) {
  const answer = await prompts({
    type: 'confirm',
    name: 'confirmed',
    message: question,
    initial: false,
  });

  if (!answer.confirmed) {
    console.info('\r\n');
    console.info(`----ABORT SCRIPT----`);
    throw new Error(`Do not accept question: ${question}`);
  }
}

/**
 * Run a command synchronously and return the output.
 * @param command The command to run.
 * @param cwd The working directory to run the command in. RELATIVE TO ROOT_DIR
 * @returns The output of the command.
 */
async function runCommandSync(command: string, cwd?: string): Promise<string> {
  const workingDir = cwd ? `${ROOT_DIR}` + cwd : `${ROOT_DIR}`;
  // console.info(`Running command: ${command} in ${workingDir}`);
  return execSync(command, { encoding: 'utf8', cwd: workingDir, stdio: ['ignore', 'pipe', 'pipe'] });
}

async function waitForNumBlocks(ethers: HardhatEthersHelpers, num: number): Promise<void> {
  let currentBlock = await ethers.provider.getBlockNumber();
  const startBlock = currentBlock;
  const targetBlock = currentBlock + num;
  const startTime = Date.now();
  while (currentBlock < targetBlock) {
    await new Promise((resolve) => setTimeout(resolve, 1000)); // wait for 1 second
    currentBlock = await ethers.provider.getBlockNumber();
    // console.info(`Waiting: ${startBlock}==${currentBlock}==>${targetBlock}`);

    // Check if the time limit has been reached
    const elapsedTime = Date.now() - startTime;
    if (elapsedTime > MAX_BLOCK_INTERVAL * num) {
      throw new Error(`Timeout: waited too long for block number to reach ${num}`);
    }
  }
}

async function runTitandTx(ethers: HardhatEthersHelpers, command: string, signerName = 'faucet'): Promise<any> {
  let params = `--home local_test_data/.titan_val1 --keyring-backend test --output json`;
  params += ` --node http://localhost:26657 --chain-id titan_18887-1 --from ${signerName} --gas=auto --gas-prices=100000000000atkx --gas-adjustment=2 -y`;

  const cmd = command + ` ${params}`;
  const output = await runCommandSync(cmd);
  await waitForNumBlocks(ethers, 1);

  return output;
}

async function runTitandQuery(command: string): Promise<any> {
  let params = `--home local_test_data/.titan_val1 --output json`;
  params += ` --node http://localhost:26657 --chain-id titan_18887-1`;

  const cmd = command + ` ${params}`;
  const output = await runCommandSync(cmd);
  return JSON.parse(output);
}

async function getAddressBalance(address: string, tokenDenom: string): Promise<BigNumber> {
  const balance = await runTitandQuery(`titand q bank balances ${address}`);
  const tokenBalance = balance.balances.find((b: { denom: string }) => b.denom === tokenDenom);
  if (!tokenBalance) {
    return BigNumber.from(0);
  }
  return BigNumber.from(tokenBalance.amount);
}

async function evmAddressToTitanAddress(evmAddress: string): Promise<string> {
  // Convert EVM address to Titan address format use bech32 encoding
  const titanPrefix = 'titan';
  // convert evmAddress to bytes
  const bytes = Buffer.from(evmAddress.slice(2), 'hex');
  // convert bytes to bech32
  const titanAddress = bech32.encode(titanPrefix, bech32.toWords(bytes));
  return titanAddress;
}

export default {
  confirm,
  runCommandSync,
  waitForNumBlocks,
  runTitandTx,
  runTitandQuery,
  getAddressBalance,
  evmAddressToTitanAddress,
};
