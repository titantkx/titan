import hre, { ethers } from 'hardhat';
import utils from '../utils';

async function main() {
  const { SMART_CONTRACT_NAME, CONSTRUCTOR_ARGUMENTS } = utils.processSmartContractConfig(
    hre.userConfig.smartContractConfig,
    hre
  );
  console.log('!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! DEPLOY NORMAL !!!!!!!!!!!!!!!!!!!!!!!!!!');
  console.log(`SMART_CONTRACT_NAME=${SMART_CONTRACT_NAME}`);
  console.log(`Current NETWORK=${hre.network.name}`);
  const constructorArguments: Array<any> = CONSTRUCTOR_ARGUMENTS || [];

  // Hardhat always runs the compile task when running scripts with its command
  // line interface.
  //
  const [deployer] = await ethers.getSigners();
  console.log('Deploying contracts with the account:', deployer.address);
  console.log('Account balance:', ethers.utils.formatEther(await deployer.getBalance()));

  // // If this script is run directly using `node` you may want to call compile
  // // manually to make sure everything is compiled
  // await hre.run('compile');

  // We get the contract to deploy
  const Contract = await hre.ethers.getContractFactory(SMART_CONTRACT_NAME);

  //

  await utils.confirm('Do you want to deploy NORMAL ?');
  console.log('constructorArguments', constructorArguments);
  await utils.confirm('Are you sure want to deploy NORMAL ?');
  //

  const contract = await (<any>Contract.deploy)(...constructorArguments);
  console.log('Smartcontract deploy transaction:', contract.deployTransaction.hash);
  console.log('Smartcontract deploy to:', contract.address);

  utils.writeDeployInfo({
    SMART_CONTRACT_NAME,
    NODE_ENV: process.env.NODE_ENV,
    networkName: hre.network.name,
    transaction: contract.deployTransaction.hash,
    type: 'DEPLOY',
    contractAddress: contract.address,
  });

  await contract.deployed();

  console.log('contract deployed to:', contract.address);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
