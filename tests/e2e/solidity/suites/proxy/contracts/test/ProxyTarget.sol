pragma solidity 0.4.26;

contract ProxyTargetWithoutFallback {
    event Pong();

    function ping() external {
      emit Pong();
    }
}

contract ProxyTargetWithFallback is ProxyTargetWithoutFallback {
    event ReceivedEth();

    function () external payable {
      emit ReceivedEth();
    }
}
