import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers';
import { ethers } from 'hardhat';

import pointerAbi from 'precompiles/pointer/abi.json';

describe('Sample test', function () {
  var owner: SignerWithAddress, wallets: SignerWithAddress[];

  async function setup() {
    [owner, ...wallets] = await ethers.getSigners();
  }

  before(async function () {
    await setup();
  });

  describe('Sample', function () {
    it('Test', async function () {
      console.log('owner', owner.address);
      console.log(
        'wallets',
        wallets.map((w) => w.address)
      );

      // const nativeERC20Addr = '0x0000000000000000000000000000000000001001';
      // const nativeERC20 = NativeTokensERC20__factory.connect(nativeERC20Addr, owner);
      // nativeERC20.name().then((name) => {
      //   console.log('name', name);
      // });

      const pointerAddr = '0x000000000000000000000000000000000000100b';
      const pointer = new ethers.Contract(pointerAddr, pointerAbi, owner);
      const tx = await pointer.addNativePointer('atkx');
      console.log('Transaction:', tx);

      const receipt = await tx.wait();
      console.log('Transaction receipt:', receipt);
      console.log('Transaction hash:', tx.hash);
    });
  });
});
