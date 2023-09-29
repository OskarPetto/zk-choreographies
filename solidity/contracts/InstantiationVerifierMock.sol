// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract InstantiationVerifierMock {
    function verifyProof(
        uint256[8] calldata proof,
        uint256[2] calldata input
    ) public view {}
}
