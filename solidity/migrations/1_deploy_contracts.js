const InstantiationVerifier = artifacts.require("InstantiationVerifier");
const TransitionVerifier = artifacts.require("TransitionVerifier");
const TerminationVerifier = artifacts.require("TerminationVerifier");
const InstanceManager = artifacts.require("InstanceManager");

module.exports = function (deployer) {
  deployer.then(async () => {
    await deployer.deploy(InstantiationVerifier);
    await deployer.deploy(TransitionVerifier);
    await deployer.deploy(TerminationVerifier);
    await deployer.deploy(InstanceManager, InstantiationVerifier.address, TransitionVerifier.address, TerminationVerifier.address);
  });
};
