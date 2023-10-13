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
        address _instantiationVerifier,
        address _transitionVerifier,
        address _terminationVerifier
    ) {
        instantiationVerifier = InstantiationVerifier(_instantiationVerifier);
        transitionVerifier = TransitionVerifier(_transitionVerifier);
        terminationVerifier = TerminationVerifier(_terminationVerifier);
    }

    function instantiate(uint[8] memory proof, uint instance) public {
        require(instances[instance] == false);
        instantiationVerifier.verifyProof(proof, [instance]);
        instances[instance] = true;
        emit Instantiation(instance);
    }

    function transition(
        uint[8] memory proof,
        uint currentInstance,
        uint nextInstance
    ) public {
        require(instances[currentInstance] == true);
        transitionVerifier.verifyProof(proof, [currentInstance, nextInstance]);
        instances[currentInstance] = false;
        instances[nextInstance] = true;
        emit Transition(currentInstance, nextInstance);
    }

    function terminate(uint[8] memory proof, uint instance) public {
        require(instances[instance] == true);
        terminationVerifier.verifyProof(proof, [instance]);

        delete instances[instance];
        emit Termination(instance);
    }
}
