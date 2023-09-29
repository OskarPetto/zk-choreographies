const Pairing = artifacts.require("Pairing");
const InstantiationVerifier = artifacts.require("InstantiationVerifier");
const TransitionVerifier = artifacts.require("TransitionVerifier");
const TerminationVerifier = artifacts.require("TerminationVerifier");
const InstanceManager = artifacts.require("InstanceManager");

module.exports = function (deployer) {
  deployer.then(async () => {
    await deployer.deploy(Pairing);
    deployer.link(Pairing, InstantiationVerifier);
    deployer.link(Pairing, TransitionVerifier);
    deployer.link(Pairing, TerminationVerifier);
    await deployer.deploy(InstantiationVerifier);
    await deployer.deploy(TransitionVerifier);
    await deployer.deploy(TerminationVerifier);
    await deployer.deploy(InstanceManager, InstantiationVerifier.address, TransitionVerifier.address, TerminationVerifier.address);
  });
};
