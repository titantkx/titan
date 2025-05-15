import { loadFixture } from '@nomicfoundation/hardhat-network-helpers';
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers';
import { expect } from 'chai';
import hre, { ethers } from 'hardhat';
import utils from '../utils';

describe('Sample test', function () {
  var owner: SignerWithAddress, wallets: SignerWithAddress[];

  async function deployLockFixture() {
    [owner, ...wallets] = await ethers.getSigners();

    const configs = utils.processSmartContractConfig(hre.userConfig.smartContractConfig, hre, {});
    // !!! NOTE: this override config for test example. You must remove if want to copy this code for other test
    configs.SMART_CONTRACT_NAME = 'Greeter';
    configs.CONSTRUCTOR_ARGUMENTS = ['Hello, world!'];

    // deploy contract
    const Contract = await ethers.getContractFactory(configs.SMART_CONTRACT_NAME);
    const contract = await (<any>Contract.deploy)(...configs.CONSTRUCTOR_ARGUMENTS);
    await contract.deployed();

    return { contract, configs };
  }

  before(async function () {
    await hre.network.provider.send('hardhat_reset');
    const { contract, configs } = await loadFixture(deployLockFixture);
  });

  describe('Sample', function () {
    it("Should return the new greeting once it's changed", async function () {
      const { contract, configs } = await loadFixture(deployLockFixture);

      expect(await contract.greet()).to.equal('Hello, world!');

      const setGreetingTx = await contract.setGreeting('Hola, mundo!');

      // wait until the transaction is mined
      await setGreetingTx.wait();

      expect(await contract.greet()).to.equal('Hola, mundo!');
    });
  });
});
