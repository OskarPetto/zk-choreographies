pragma solidity ^0.8.0;

import "./InstantiationVerifier.sol";
import "./TransitionVerifier.sol";
import "./TerminationVerifier.sol";

struct Proof {
    uint256[2] a;
    uint256[2][2] b;
    uint256[2] c;
}

contract InstanceManager {
    mapping(uint => uint) public instancesPerModel;

    InstantiationVerifier instantiationVerifier;
    TransitionVerifier transitionVerifier;
    TerminationVerifier terminationVerifier;

    function instantiate(uint model, uint instance, Proof memory proof) public {
        require(instancesPerModel[model] == 0);
        require(
            instantiationVerifier.verifyProof(
                proof.a,
                proof.b,
                proof.c,
                [model, instance]
            )
        );

        instancesPerModel[model] = instance;
    }

    function transition(
        uint model,
        uint currentInstance,
        uint nextInstance,
        Proof memory proof
    ) public {
        require(instancesPerModel[model] == currentInstance);
        require(
            transitionVerifier.verifyProof(
                proof.a,
                proof.b,
                proof.c,
                [model, currentInstance, nextInstance, 0] // fourth value should not exist
            )
        );
        instancesPerModel[model] = nextInstance;
    }

    function termination(uint model, uint instance, Proof memory proof) public {
        require(instancesPerModel[model] == instance);
        require(
            terminationVerifier.verifyProof(
                proof.a,
                proof.b,
                proof.c,
                [model, instance]
            )
        );
        delete instancesPerModel[model];
    }
}
