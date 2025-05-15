import { ContractFactory } from 'ethers';
import { ethers } from 'hardhat';

const deployContract = async function (contractName: string, constructorArgs: any[]) {
  const factory = await ethers.getContractFactory(contractName);
  const contract = await factory.deploy(...(constructorArgs || []));
  await contract.deployed();
  return contract;
};

const getContractFactory = async function (contractName: string): Promise<ContractFactory> {
  let factory: ContractFactory;
  try {
    factory = await ethers.getContractFactory(contractName);
  } catch (e) {
    factory = await ethers.getContractFactory(contractName + 'Upgradeable');
  }
  return factory;
};

export { deployContract, getContractFactory };
