import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers';
import { NativeTokensERC20__factory } from '@pre-contracts/factory/contracts/NativeTokensERC20__factory';
import { ethers } from 'hardhat';

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
      // create factory token
      // titand --home local_test_data/.titan_val1 --keyring-backend test tx tokenfactory create-denom test
      // --node http://localhost:26657 --chain-id titan_18887-1 --from faucet --gas=auto --gas-prices=100000000000atkx --gas-adjustment=2 -y

      // const nativeERC20Addr = '0x0000000000000000000000000000000000001001';
      // const nativeERC20 = NativeTokensERC20__factory.connect(nativeERC20Addr, owner);
      // nativeERC20.name().then((name) => {
      //   console.log('name', name);
      // });

      // set denom metadata

      // {
      //   "description": "",
      //   "denom_units": [
      //     {
      //       "denom": "factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/test",
      //       "exponent": 0,
      //       "aliases": []
      //     },
      //     {
      //       "denom": "test",
      //       "exponent": 3,
      //       "aliases": []
      //     }
      //   ],
      //   "base": "factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/test",
      //   "display": "test",
      //   "name": "test",
      //   "symbol": "TEST",
      //   "uri": "",
      //   "uri_hash": ""
      // }

      // titand --home local_test_data/.titan_val1 --keyring-backend test tx tokenfactory set-denom-metadata "{\"description\":\"\",\"denom_units\":[{\"denom\":\"factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/test\",\"exponent\":0,\"aliases\":[]},{\"denom\":\"test\",\"exponent\":3,\"aliases\":[]}],\"base\":\"factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/test\",\"display\":\"test\",\"name\":\"test\",\"symbol\":\"TEST\",\"uri\":\"\",\"uri_hash\":\"\"}"
      // --node http://localhost:26657 --chain-id titan_18887-1 --from faucet --gas=auto --gas-prices=100000000000atkx --gas-adjustment=2 -y

      //{"description":"","denom_units":[{"denom":"factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/test","exponent":0,"aliases":[]},{"denom":"test","exponent":3,"aliases":[]}],"base":"factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/test","display":"test","name":"test","symbol":"TEST","uri":"","uri_hash":""}

      // const pointerAddr = '0x000000000000000000000000000000000000100b';
      // const pointer = new ethers.Contract(pointerAddr, pointerAbi, owner);
      // const tx = await pointer.addNativePointer('factory/titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr/test');
      // console.log('Transaction:', tx);

      // const receipt = await tx.wait();
      // console.log('Transaction receipt:', receipt);
      // console.log('Transaction hash:', tx.hash);

      // 0x88C14eC9481F259009c9e5a7bb7f621F4c590e21

      const nativeERC20Addr = '0x88C14eC9481F259009c9e5a7bb7f621F4c590e21';
      const nativeERC20 = NativeTokensERC20__factory.connect(nativeERC20Addr, owner);
      const name = await nativeERC20.name();
      console.log('name', name);
      const symbol = await nativeERC20.symbol();
      console.log('symbol', symbol);
      const decimals = await nativeERC20.decimals();
      console.log('decimals', decimals);
      const denom = await nativeERC20.denom();
      console.log('denom', denom);
      const totalSupply = await nativeERC20.totalSupply();
      console.log('totalSupply', totalSupply.toString());
      const balance = await nativeERC20.balanceOf(owner.address);
      console.log('balance', balance.toString());
    });
  });
});
