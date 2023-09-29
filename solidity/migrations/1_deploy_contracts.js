const Pairing = artifacts.require("Pairing");
const InstantiationVerifier = artifacts.require("InstantiationVerifier");
const TransitionVerifier = artifacts.require("TransitionVerifier");
const TerminationVerifier = artifacts.require("TerminationVerifier");
const InstanceManager = artifacts.require("InstanceManager");

module.exports = function (deployer) {
  deployer.deploy(Pairing);
  deployer.link(Pairing, InstantiationVerifier);
  deployer.link(Pairing, TransitionVerifier);
  deployer.link(Pairing, TerminationVerifier);
  deployer.deploy(InstantiationVerifier);
  deployer.deploy(TransitionVerifier);
  deployer.deploy(TerminationVerifier);
  deployer.deploy(InstanceManager);
};
