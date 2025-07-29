// SPDX-License-Identifier: MIT
pragma solidity ^0.8.16;

contract SelfDestruct {
  // accept receive tokens
  receive() external payable {}

  // function to self-destruct the contract and send all ether to the caller
  function selfDestruct() external {
    selfdestruct(payable(msg.sender));
  }

  // function to self-destruct the contract and not send any ether
  function selfDestructWithoutEther() external {
    selfdestruct(payable(address(0)));
  }
}
