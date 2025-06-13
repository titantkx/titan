import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers';
import { ethers } from 'hardhat';

import { NativeTokensERC20__factory } from '@contracts/factories/contracts/NativeTokensERC20__factory.ts';
import { expect } from 'chai';
import { BigNumber } from 'ethers';
import utils from 'utils';

describe('ERC20 test', function () {
  var faucet: SignerWithAddress, wallets: SignerWithAddress[];

  const FAUCET_TITAN_ADDR = 'titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr';
  const TOKEN_DENOM = `factory/${FAUCET_TITAN_ADDR}/test`;
  const TOKEN_DECIMALS = 3;
  const TEST_ERC20_ADDR = '0x88C14eC9481F259009c9e5a7bb7f621F4c590e21';

  async function setup() {
    [faucet, ...wallets] = await ethers.getSigners();
  }

  before(async function () {
    await setup();
  });

  it('Can transfer test token from faucet to another', async function () {
    const testERC20 = NativeTokensERC20__factory.connect(TEST_ERC20_ADDR, faucet);

    const recipient = wallets[0];

    // check balance of owner and recipient via BANK
    const oldOwnerBalanceBank = await utils.getAddressBalance(FAUCET_TITAN_ADDR, TOKEN_DENOM);
    const oldRecipientBalanceBank = await utils.getAddressBalance(
      await utils.evmAddressToTitanAddress(recipient.address),
      TOKEN_DENOM
    );

    // check balance of owner and recipient via EVM
    const oldOwnerBalance = await testERC20.balanceOf(faucet.address);
    const oldRecipientBalance = await testERC20.balanceOf(recipient.address);

    const tx = await testERC20.transfer(recipient.address, BigNumber.from('1000'));

    await tx.wait();

    // check balance of owner and recipient via BANK
    const newOwnerBalanceBank = await utils.getAddressBalance(FAUCET_TITAN_ADDR, TOKEN_DENOM);
    const newRecipientBalanceBank = await utils.getAddressBalance(
      await utils.evmAddressToTitanAddress(recipient.address),
      TOKEN_DENOM
    );

    // check balance of owner and recipient via EVM
    const newOwnerBalance = await testERC20.balanceOf(faucet.address);
    const newRecipientBalance = await testERC20.balanceOf(recipient.address);

    expect(newOwnerBalanceBank).to.equal(oldOwnerBalanceBank.sub(BigNumber.from('1000')));
    expect(newRecipientBalanceBank).to.equal(oldRecipientBalanceBank.add(BigNumber.from('1000')));
    expect(newOwnerBalance).to.equal(oldOwnerBalance.sub(BigNumber.from('1000')));
    expect(newRecipientBalance).to.equal(oldRecipientBalance.add(BigNumber.from('1000')));
  });
});
