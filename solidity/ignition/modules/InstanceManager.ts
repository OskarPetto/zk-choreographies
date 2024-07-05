import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";


const InstanceManagerModule = buildModule("InstanceManagerModule", (m) => {
    const instantiationVerifier = m.contract("InstantiationVerifier");
    const transitionVerifier = m.contract("TransitionVerifier");
    const terminationVerifier = m.contract("TerminationVerifier");

    const instanceManager = m.contract("InstanceManager", [instantiationVerifier, transitionVerifier, terminationVerifier]);

    return { instanceManager };
});

export default InstanceManagerModule;
