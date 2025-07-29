import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers';
import { expect } from 'chai';
import { ethers } from 'hardhat';

import { NativeTokensERC20__factory } from '@contracts/factories/contracts/NativeTokensERC20__factory.ts';
import pointerAbi from '@precompiles/pointer/abi.json';
import { BigNumber } from 'ethers';
import utils from 'utils';

describe('Register test', function () {
  var owner: SignerWithAddress, wallets: SignerWithAddress[];
  const FAUCET_TITAN_ADDR = 'titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr';

  async function setup() {
    [owner, ...wallets] = await ethers.getSigners();
  }

  before(async function () {
    await setup();
  });

  it('Can not register pointer contract for base token', async function () {
    const pointerAddr = '0x000000000000000000000000000000000000100b';
    const pointer = new ethers.Contract(pointerAddr, pointerAbi, owner);
    try {
      const tx = await pointer.addNativePointer(`atkx`);
      await tx.wait();
      expect.fail('Should not be able to register pointer for base token');
    } catch (e: unknown) {
      expect(e).to.be.instanceOf(Error);
      expect((e as Error).message).to.include('cannot create pointer for base token atkx');
    }
    await utils.waitForNumBlocks(ethers, 1);
  });

  it('Can register pointer contract', async function () {
    // create factory token
    const tokenDenom = 'test';
    const tokenDecimals = 3;

    const createTokenCmd = `titand tx tokenfactory create-denom ${tokenDenom}`;
    try {
      await utils.runTitandTx(ethers, createTokenCmd);
    } catch (e: unknown) {
      if (e instanceof Error && e.message.includes('attempting to create a denom that already exists')) {
        console.log('Token already exists.');
      } else {
        throw e;
      }
    }

    const tokenMetadata = {
      description: '',
      denom_units: [
        {
          denom: `factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`,
          exponent: 0,
          aliases: [],
        },
        {
          denom: tokenDenom,
          exponent: tokenDecimals,
          aliases: [],
        },
      ],
      base: `factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`,
      display: tokenDenom,
      name: tokenDenom,
      symbol: tokenDenom.toUpperCase(),
      uri: '',
      uri_hash: '',
    };
    const setTokenDenomCmd = `titand tx tokenfactory set-denom-metadata ${JSON.stringify(
      JSON.stringify(tokenMetadata)
    )} `;

    await utils.runTitandTx(ethers, setTokenDenomCmd);

    // mint new token to faucet
    const mintCmd = `titand tx tokenfactory mint 1000000000factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom} ${FAUCET_TITAN_ADDR}`;
    await utils.runTitandTx(ethers, mintCmd);

    // get total supply
    const totalSupply = (await utils.runTitandQuery(`titand q bank total`)) as {
      supply: { denom: string; amount: string }[];
    };
    const newTokenSupply = totalSupply.supply.find(
      (supply) => supply.denom === `factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`
    );
    expect(newTokenSupply).to.not.be.undefined;

    // get faucet balance
    const faucetBalance = (await utils.runTitandQuery(`titand q bank balances ${FAUCET_TITAN_ADDR}`)) as {
      balances: { denom: string; amount: string }[];
    };
    const faucetTokenBalance = faucetBalance.balances.find(
      (balance) => balance.denom === `factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`
    );
    expect(faucetTokenBalance).to.not.be.undefined;

    /// deploy pointer contract
    ///
    ///
    const pointerAddr = '0x000000000000000000000000000000000000100b';
    const pointer = new ethers.Contract(pointerAddr, pointerAbi, owner);
    const tx = await pointer.addNativePointer(`factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`);
    const receipt = await tx.wait();
    // console.log('Transaction hash:', tx.hash);
    // check event
    const event = receipt.events?.find((event: any) => event.event === 'PointerRegistered');
    expect(event.args.pointerType).to.equal(0);
    expect(event.args.token).to.equal(`factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`);
    const pointerContract = event.args.pointer;
    expect(pointerContract).to.equal('0x88C14eC9481F259009c9e5a7bb7f621F4c590e21');

    ///  Check pointer contract
    ///
    ///

    // 0x88c14ec9481f259009c9e5a7bb7f621f4c590e21;
    // 0x1C3a428274c8f06DfecE6dc03437214A78285898
    const testERC20Addr = '0x88C14eC9481F259009c9e5a7bb7f621F4c590e21';
    // const test2ERC20Addr = '0x1C3a428274c8f06DfecE6dc03437214A78285898';
    const testERC20 = NativeTokensERC20__factory.connect(testERC20Addr, owner);
    // const test2ERC20 = NativeTokensERC20__factory.connect(test2ERC20Addr, owner);
    expect(await testERC20.name()).to.equal(tokenDenom);
    expect(await testERC20.symbol()).to.equal(tokenDenom.toUpperCase());
    expect(await testERC20.decimals()).to.equal(tokenDecimals);
    expect(await testERC20.denom()).to.equal(`factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`);
    expect(await testERC20.totalSupply()).to.equal(BigNumber.from(newTokenSupply?.amount));
    expect(await testERC20.balanceOf(owner.address)).to.equal(BigNumber.from(faucetTokenBalance?.amount));
  });

  it('Can register another pointer contract', async function () {
    // create factory token
    const tokenDenom = 'test2';
    const tokenDecimals = 3;

    const createTokenCmd = `titand tx tokenfactory create-denom ${tokenDenom}`;
    try {
      await utils.runTitandTx(ethers, createTokenCmd);
    } catch (e: unknown) {
      if (e instanceof Error && e.message.includes('attempting to create a denom that already exists')) {
        console.log('Token already exists.');
      } else {
        throw e;
      }
    }

    const tokenMetadata = {
      description: '',
      denom_units: [
        {
          denom: `factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`,
          exponent: 0,
          aliases: [],
        },
        {
          denom: tokenDenom,
          exponent: tokenDecimals,
          aliases: [],
        },
      ],
      base: `factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`,
      display: tokenDenom,
      name: tokenDenom,
      symbol: tokenDenom.toUpperCase(),
      uri: '',
      uri_hash: '',
    };
    const setTokenDenomCmd = `titand tx tokenfactory set-denom-metadata ${JSON.stringify(
      JSON.stringify(tokenMetadata)
    )} `;

    await utils.runTitandTx(ethers, setTokenDenomCmd);

    // mint new token to faucet
    const mintCmd = `titand tx tokenfactory mint 1000000000factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom} ${FAUCET_TITAN_ADDR}`;
    await utils.runTitandTx(ethers, mintCmd);

    // get total supply
    const totalSupply = (await utils.runTitandQuery(`titand q bank total`)) as {
      supply: { denom: string; amount: string }[];
    };
    const newTokenSupply = totalSupply.supply.find(
      (supply) => supply.denom === `factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`
    );
    expect(newTokenSupply).to.not.be.undefined;

    // get faucet balance
    const faucetBalance = (await utils.runTitandQuery(`titand q bank balances ${FAUCET_TITAN_ADDR}`)) as {
      balances: { denom: string; amount: string }[];
    };
    const faucetTokenBalance = faucetBalance.balances.find(
      (balance) => balance.denom === `factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`
    );
    expect(faucetTokenBalance).to.not.be.undefined;

    /// deploy pointer contract
    ///
    ///

    const pointerAddr = '0x000000000000000000000000000000000000100b';
    const pointer = new ethers.Contract(pointerAddr, pointerAbi, owner);
    const tx = await pointer.addNativePointer(`factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`);
    const receipt = await tx.wait();
    // console.log('Transaction hash:', tx.hash);
    // check event
    const event = receipt.events?.find((event: any) => event.event === 'PointerRegistered');
    expect(event.args.pointerType).to.equal(0);
    expect(event.args.token).to.equal(`factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`);
    const pointerContract = event.args.pointer;
    expect(pointerContract).to.equal('0x1C3a428274c8f06DfecE6dc03437214A78285898');

    ///  Check pointer contract
    ///
    ///
    const test2ERC20Addr = '0x1C3a428274c8f06DfecE6dc03437214A78285898';
    const test2ERC20 = NativeTokensERC20__factory.connect(test2ERC20Addr, owner);
    expect(await test2ERC20.name()).to.equal(tokenDenom);
    expect(await test2ERC20.symbol()).to.equal(tokenDenom.toUpperCase());
    expect(await test2ERC20.decimals()).to.equal(tokenDecimals);
    expect(await test2ERC20.denom()).to.equal(`factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/${tokenDenom}`);
    expect(await test2ERC20.totalSupply()).to.equal(BigNumber.from(newTokenSupply?.amount));
    expect(await test2ERC20.balanceOf(owner.address)).to.equal(BigNumber.from(faucetTokenBalance?.amount));
  });
});
