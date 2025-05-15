import { ContractFactory, Signer } from 'ethers';
import { ethers, upgrades } from 'hardhat';

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

const deployProxy = async function (contractName: string, constructorArgs: any[]) {
  let factory = await getContractFactory(contractName);

  upgrades.silenceWarnings();
  const contract = await upgrades.deployProxy(factory, constructorArgs || [], {
    kind: 'uups',
    unsafeAllow: ['constructor', 'delegatecall'],
  });
  await contract.deployed();
  return contract;
};

const upgradeProxy = async function (proxyAddress: string, contractName: string, caller?: Signer) {
  var Contract = await ethers.getContractFactory(contractName);
  if (caller) {
    Contract = Contract.connect(caller);
  }
  // upgrades.silenceWarnings();
  const contract = await upgrades.upgradeProxy(proxyAddress, Contract, {
    kind: 'uups',
    unsafeAllow: ['constructor', 'delegatecall'],
  });

  await contract.deployed();
  return contract;
};

export { deployContract, deployProxy, getContractFactory, upgradeProxy };
