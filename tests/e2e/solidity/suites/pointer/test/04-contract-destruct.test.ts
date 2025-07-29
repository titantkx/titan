import { NativeTokensERC20__factory } from '@contracts/*';
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers';
import { expect } from 'chai';
import { BigNumber } from 'ethers';
import { ethers } from 'hardhat';

describe('SelfDestruct Contract', () => {
  var faucet: SignerWithAddress, wallets: SignerWithAddress[];

  async function setup() {
    [faucet, ...wallets] = await ethers.getSigners();
  }

  beforeEach(async () => {
    await setup();
  });

  describe('Contract contain native token', () => {
    it('After destruct token send back to sender', async () => {
      const selfDestructContract = await ethers.getContractFactory('SelfDestruct');
      const selfDestruct = await selfDestructContract.deploy();
      await selfDestruct.deployed();

      // transfer some native tokens to the contract
      const tx = await faucet.sendTransaction({
        to: selfDestruct.address,
        value: ethers.utils.parseEther('1.0'), // 1 TKX
      });
      await tx.wait();

      // check balance
      const oldOwnerBalance = await faucet.getBalance();
      const oldContractBalance = await ethers.provider.getBalance(selfDestruct.address);

      expect(oldContractBalance).to.equal(ethers.utils.parseEther('1.0'));

      // now self-destruct contract
      const destructTx = await selfDestruct.selfDestruct();
      await destructTx.wait();
      // get fee paid for self-destruct
      const destructReceipt = await ethers.provider.getTransactionReceipt(destructTx.hash);
      const gasFee = destructReceipt.gasUsed.mul(destructReceipt.effectiveGasPrice || 0);

      // check balance after self-destruct
      const newOwnerBalance = await faucet.getBalance();
      const newContractBalance = await ethers.provider.getBalance(selfDestruct.address);

      expect(newContractBalance).to.equal(0);
      expect(newOwnerBalance).to.be.gt(oldOwnerBalance); // owner should have more balance
      // new owner balance should be old owner balance + contract balance - gas fee
      expect(newOwnerBalance).to.equal(oldOwnerBalance.add(oldContractBalance).sub(gasFee));
    });

    it('After destruct and sent token to zero address', async () => {
      const selfDestructContract = await ethers.getContractFactory('SelfDestruct');
      const selfDestruct = await selfDestructContract.deploy();
      await selfDestruct.deployed();

      // transfer some native tokens to the contract
      const tx = await faucet.sendTransaction({
        to: selfDestruct.address,
        value: ethers.utils.parseEther('1.0'), // 1 TKX
      });
      await tx.wait();

      // check balance
      const oldZeroAddressBalance = await ethers.provider.getBalance('0x0000000000000000000000000000000000000000');
      const oldContractBalance = await ethers.provider.getBalance(selfDestruct.address);

      expect(oldContractBalance).to.equal(ethers.utils.parseEther('1.0'));

      // now self-destruct contract
      const destructTx = await selfDestruct.selfDestructWithoutEther();
      await destructTx.wait();
      // get fee paid for self-destruct
      const destructReceipt = await ethers.provider.getTransactionReceipt(destructTx.hash);
      const gasFee = destructReceipt.gasUsed.mul(destructReceipt.effectiveGasPrice || 0);

      // check balance after self-destruct
      const newZeroAddressBalance = await ethers.provider.getBalance('0x0000000000000000000000000000000000000000');
      const newContractBalance = await ethers.provider.getBalance(selfDestruct.address);

      expect(newContractBalance).to.equal(0);
      expect(newZeroAddressBalance).to.be.gt(oldZeroAddressBalance); // zero address should have more balance
      // new zero address balance should be old zero address balance + contract balance - gas fee
      expect(newZeroAddressBalance).to.equal(oldZeroAddressBalance.add(oldContractBalance));
    });
  });

  describe('Contract contain other token than native token', () => {
    const FAUCET_TITAN_ADDR = 'titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr';
    const TOKEN_DENOM = `factory/${FAUCET_TITAN_ADDR}/test`;
    const TOKEN_DECIMALS = 3;
    const TEST_ERC20_ADDR = '0x88C14eC9481F259009c9e5a7bb7f621F4c590e21';

    const TOKEN2_DENOM = `factory/${FAUCET_TITAN_ADDR}/test2`;
    const TOKEN2_DECIMALS = 3;
    const TEST2_ERC20_ADDR = '0x1C3a428274c8f06DfecE6dc03437214A78285898';

    it('After destruct token send back to sender', async () => {
      const selfDestructContract = await ethers.getContractFactory('SelfDestruct');
      const selfDestruct = await selfDestructContract.deploy();
      await selfDestruct.deployed();

      // transfer some native tokens to the contract
      const tx = await faucet.sendTransaction({
        to: selfDestruct.address,
        value: ethers.utils.parseEther('1.0'), // 1 TKX
      });
      await tx.wait();
      // transfer some test token to the contract
      const testERC20 = NativeTokensERC20__factory.connect(TEST_ERC20_ADDR, faucet);
      const testTokenTx = await testERC20.transfer(
        selfDestruct.address,
        ethers.utils.parseUnits('1000', TOKEN_DECIMALS)
      );
      await testTokenTx.wait();
      // transfer some test2 token to the contract
      const test2ERC20 = NativeTokensERC20__factory.connect(TEST2_ERC20_ADDR, faucet);
      const test2TokenTx = await test2ERC20.transfer(
        selfDestruct.address,
        ethers.utils.parseUnits('1000', TOKEN2_DECIMALS)
      );
      await test2TokenTx.wait();

      // check balance
      const oldOwnerBalance = await faucet.getBalance();
      const oldOwnerTokenBalance = await testERC20.balanceOf(faucet.address);
      const oldOwnerToken2Balance = await test2ERC20.balanceOf(faucet.address);

      const oldContractBalance = await ethers.provider.getBalance(selfDestruct.address);
      const oldContractTokenBalance = await testERC20.balanceOf(selfDestruct.address);
      const oldContractToken2Balance = await test2ERC20.balanceOf(selfDestruct.address);

      expect(oldContractBalance).to.equal(ethers.utils.parseEther('1.0'));
      expect(oldContractTokenBalance).to.equal(ethers.utils.parseUnits('1000', TOKEN_DECIMALS));
      expect(oldContractToken2Balance).to.equal(ethers.utils.parseUnits('1000', TOKEN2_DECIMALS));

      // now self-destruct contract
      const destructTx = await selfDestruct.selfDestruct();
      await destructTx.wait();
      // get fee paid for self-destruct
      const destructReceipt = await ethers.provider.getTransactionReceipt(destructTx.hash);
      const gasFee = destructReceipt.gasUsed.mul(destructReceipt.effectiveGasPrice || 0);

      // check balance after self-destruct
      const newOwnerBalance = await faucet.getBalance();
      const newOwnerTokenBalance = await testERC20.balanceOf(faucet.address);
      const newOwnerToken2Balance = await test2ERC20.balanceOf(faucet.address);

      const newContractBalance = await ethers.provider.getBalance(selfDestruct.address);
      const newContractTokenBalance = await testERC20.balanceOf(selfDestruct.address);
      const newContractToken2Balance = await test2ERC20.balanceOf(selfDestruct.address);

      // check native token
      expect(newContractBalance).to.equal(0);
      expect(newOwnerBalance).to.be.gt(oldOwnerBalance); // owner should have more balance
      // new owner balance should be old owner balance + contract balance - gas fee
      expect(newOwnerBalance).to.equal(oldOwnerBalance.add(oldContractBalance).sub(gasFee));
      // check test token
      expect(newOwnerTokenBalance).to.equal(oldOwnerTokenBalance);
      expect(newOwnerToken2Balance).to.equal(oldOwnerToken2Balance);
      expect(newContractTokenBalance).to.equal(oldContractTokenBalance);
      expect(newContractToken2Balance).to.equal(oldContractToken2Balance);
    });

    it('After destruct and sent token to zero address', async () => {
      const selfDestructContract = await ethers.getContractFactory('SelfDestruct');
      const selfDestruct = await selfDestructContract.deploy();
      await selfDestruct.deployed();

      // transfer some native tokens to the contract
      const tx = await faucet.sendTransaction({
        to: selfDestruct.address,
        value: ethers.utils.parseEther('1.0'), // 1 TKX
      });
      await tx.wait();
      // transfer some test token to the contract
      const testERC20 = NativeTokensERC20__factory.connect(TEST_ERC20_ADDR, faucet);
      const testTokenTx = await testERC20.transfer(
        selfDestruct.address,
        ethers.utils.parseUnits('1000', TOKEN_DECIMALS)
      );
      await testTokenTx.wait();
      // transfer some test2 token to the contract
      const test2ERC20 = NativeTokensERC20__factory.connect(TEST2_ERC20_ADDR, faucet);
      const test2TokenTx = await test2ERC20.transfer(
        selfDestruct.address,
        ethers.utils.parseUnits('1000', TOKEN2_DECIMALS)
      );
      await test2TokenTx.wait();

      // check balance
      const oldContractBalance = await ethers.provider.getBalance(selfDestruct.address);
      const oldContractTokenBalance = await testERC20.balanceOf(selfDestruct.address);
      const oldContractToken2Balance = await test2ERC20.balanceOf(selfDestruct.address);

      const oldZeroAddressBalance = await ethers.provider.getBalance('0x0000000000000000000000000000000000000000');
      let oldZeroAddressTokenBalance = BigNumber.from(0);
      let oldZeroAddressToken2Balance = BigNumber.from(0);
      try {
        oldZeroAddressTokenBalance = await testERC20.balanceOf('0x0000000000000000000000000000000000000000');
        oldZeroAddressToken2Balance = await test2ERC20.balanceOf('0x0000000000000000000000000000000000000000');
      } catch (e) {}

      const destructTx = await selfDestruct.selfDestructWithoutEther();
      await destructTx.wait();

      const newContractBalance = await ethers.provider.getBalance(selfDestruct.address);
      const newContractTokenBalance = await testERC20.balanceOf(selfDestruct.address);
      const newContractToken2Balance = await test2ERC20.balanceOf(selfDestruct.address);

      const newZeroAddressBalance = await ethers.provider.getBalance('0x0000000000000000000000000000000000000000');
      let newZeroAddressTokenBalance = BigNumber.from(0);
      let newZeroAddressToken2Balance = BigNumber.from(0);
      try {
        newZeroAddressTokenBalance = await testERC20.balanceOf('0x0000000000000000000000000000000000000000');
        newZeroAddressToken2Balance = await test2ERC20.balanceOf('0x0000000000000000000000000000000000000000');
      } catch (e) {}

      expect(newContractBalance).to.equal(0);
      expect(newContractTokenBalance).to.equal(oldContractTokenBalance);
      expect(newContractToken2Balance).to.equal(oldContractToken2Balance);
      expect(newZeroAddressBalance).to.equal(oldZeroAddressBalance.add(oldContractBalance));
      expect(newZeroAddressTokenBalance).to.equal(oldZeroAddressTokenBalance);
      expect(newZeroAddressToken2Balance).to.equal(oldZeroAddressToken2Balance);
    });
  });
});
