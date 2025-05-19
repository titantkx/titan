// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

address constant POINTER_PRECOMPILE_ADDRESS = 0x000000000000000000000000000000000000100b;

IPointer constant POINTER_CONTRACT = IPointer(POINTER_PRECOMPILE_ADDRESS);

enum PointerType {
    Native
}

interface IPointer {
    event PointerRegistered(
        PointerType pointerType,
        string token,
        address pointer
    );

    function addNativePointer(
        string memory token
    ) external returns (address ret);
}
