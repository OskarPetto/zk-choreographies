// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./InstantiationVerifier.sol";
import "./TransitionVerifier.sol";
import "./TerminationVerifier.sol";

contract InstanceManager {
    mapping(uint => uint) public instancesWithTheirModel;

    InstantiationVerifier instantiationVerifier;
    TransitionVerifier transitionVerifier;
    TerminationVerifier terminationVerifier;

    event Instantiation(uint model, uint instance);
    event Transition(uint model, uint currentInstance, uint nextInstance);
    event Termination(uint model, uint instance);

    constructor(
        address _instantiationVerifier,
        address _transitionVerifier,
        address _terminationVerifier
    ) {
        instantiationVerifier = InstantiationVerifier(_instantiationVerifier);
        transitionVerifier = TransitionVerifier(_transitionVerifier);
        terminationVerifier = TerminationVerifier(_terminationVerifier);
    }

    function instantiate(
        uint[8] memory proof,
        uint model,
        uint instance
    ) public {
        instantiationVerifier.verifyProof(proof, [model, instance]);
        instancesWithTheirModel[instance] = model;
        emit Instantiation(model, instance);
    }

    function transition(
        uint[8] memory proof,
        uint model,
        uint currentInstance,
        uint nextInstance
    ) public {
        require(instancesWithTheirModel[currentInstance] == model);
        transitionVerifier.verifyProof(
            proof,
            [model, currentInstance, nextInstance]
        );
        delete instancesWithTheirModel[currentInstance];
        instancesWithTheirModel[nextInstance] = model;
        emit Transition(model, currentInstance, nextInstance);
    }

    function terminate(uint[8] memory proof, uint model, uint instance) public {
        require(instancesWithTheirModel[instance] == model);
        terminationVerifier.verifyProof(proof, [model, instance]);

        delete instancesWithTheirModel[instance];
        emit Termination(model, instance);
    }
}
