import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers';
import { expect } from 'chai';
import { ethers } from 'hardhat';

describe('Sample test', function () {
  var owner: SignerWithAddress, wallets: SignerWithAddress[];

  async function deploy() {
    [owner, ...wallets] = await ethers.getSigners();

    // deploy contract
    const Contract = await ethers.getContractFactory('NativeTokensERC20');
    const contract = await Contract.deploy('ibc/abcdef', 'a', 'A', 6);
    await contract.deployed();

    return { contract };
  }

  before(async function () {});

  describe('Sample', function () {
    it("Should return the new greeting once it's changed", async function () {
      const { contract } = await deploy();

      expect(await contract.name()).to.equal('a');
    });
  });
});
