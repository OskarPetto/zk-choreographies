// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./InstantiationVerifier.sol";
import "./TransitionVerifier.sol";
import "./TerminationVerifier.sol";

contract InstanceManager {
    mapping(uint => bool) public instances;

    InstantiationVerifier instantiationVerifier;
    TransitionVerifier transitionVerifier;
    TerminationVerifier terminationVerifier;

    event Instantiation(uint instance);
    event Transition(uint currentInstance, uint nextInstance);
    event Termination(uint instance);

    constructor(
        address _instantiation,
        address _transition,
        address _termination
    ) {
        instantiationVerifier = InstantiationVerifier(_instantiation);
        transitionVerifier = TransitionVerifier(_transition);
        terminationVerifier = TerminationVerifier(_termination);
    }

    function instantiate(uint[8] memory proof, uint next) public {
        require(instances[next] == false);
        instantiationVerifier.verifyProof(proof, [next]);
        instances[next] = true;
        emit Instantiation(next);
    }

    function transition(uint[8] memory proof, uint current, uint next) public {
        require(instances[current] == true);
        transitionVerifier.verifyProof(proof, [current, next]);
        delete instances[current];
        instances[next] = true;
        emit Transition(current, next);
    }

    function terminate(uint[8] memory proof, uint current) public {
        require(instances[current] == true);
        terminationVerifier.verifyProof(proof, [current]);
        delete instances[current];
        emit Termination(current);
    }
}
