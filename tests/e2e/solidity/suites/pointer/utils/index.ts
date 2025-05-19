import { HardhatEthersHelpers } from '@nomiclabs/hardhat-ethers/types';
import { execSync } from 'child_process';
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
  console.info(`Running command: ${command} in ${workingDir}`);
  try {
    return execSync(command, { encoding: 'utf8', cwd: workingDir });
  } catch (error) {
    // console.error(`Command failed: ${command}`);
    // console.error(error);
    throw error;
  }
}

async function waitForNumBlocks(ethers: HardhatEthersHelpers, num: number): Promise<void> {
  let currentBlock = await ethers.provider.getBlockNumber();
  const startBlock = currentBlock;
  const targetBlock = currentBlock + num;
  const startTime = Date.now();
  while (currentBlock < targetBlock) {
    await new Promise((resolve) => setTimeout(resolve, 1000)); // wait for 1 second
    currentBlock = await ethers.provider.getBlockNumber();
    console.info(`Waiting: ${startBlock}==${currentBlock}==>${targetBlock}`);

    // Check if the time limit has been reached
    const elapsedTime = Date.now() - startTime;
    if (elapsedTime > MAX_BLOCK_INTERVAL * num) {
      throw new Error(`Timeout: waited too long for block number to reach ${num}`);
    }
  }
}

async function runTitandTx(ethers: HardhatEthersHelpers, command: string, signerName = 'faucet') {
  let params = `--home local_test_data/.titan_val1 --keyring-backend test`;
  params += ` --node http://localhost:26657 --chain-id titan_18887-1 --from ${signerName} --gas=auto --gas-prices=100000000000atkx --gas-adjustment=2 -y`;

  const cmd = command + ` ${params}`;
  await runCommandSync(cmd);
  await waitForNumBlocks(ethers, 1);
}

export default {
  confirm,
  runCommandSync,
  waitForNumBlocks,
  runTitandTx,
};
