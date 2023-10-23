const fs = require('fs');
const proofsString = fs.readFileSync('test/proofs.json');
const proofs = JSON.parse(proofsString);

const InstantiationVerifier = artifacts.require("InstantiationVerifier");
const TransitionVerifier = artifacts.require("TransitionVerifier");
const TerminationVerifier = artifacts.require("TerminationVerifier");

const InstantiationVerifierMock = artifacts.require("InstantiationVerifierMock");

const InstanceManager = artifacts.require("InstanceManager");


contract('InstanceManager', (accounts) => {
  it('instantiation', async () => {
    const instanceManager = await InstanceManager.new(InstantiationVerifier.address, accounts[0], accounts[0]);

    const firstProof = proofs[0];
    const instantiationProof = firstProof.value;
    const instance = firstProof.input[0]

    await instanceManager.instantiate(instantiationProof, instance);

    const isStored = await instanceManager.instances(instance);
    assert.equal(isStored, true, "instance was stored");
  });

  it('termination', async () => {
    const instantiationVerifierMock = await InstantiationVerifierMock.new();
    const instanceManager = await InstanceManager.new(instantiationVerifierMock.address, accounts[0], TerminationVerifier.address);

    const lastProof = proofs[proofs.length - 1];
    const terminationProof = lastProof.value;
    const instance = lastProof.input[0]


    await instanceManager.instantiate(terminationProof, instance); // uses mock
    await instanceManager.terminate(terminationProof, instance);

    const isStored = await instanceManager.instances(instance);
    assert.equal(isStored, false, "instance was deleted");
  });

  it('transition', async () => {
    const instantiationVerifierMock = await InstantiationVerifierMock.new();
    const instanceManager = await InstanceManager.new(instantiationVerifierMock.address, TransitionVerifier.address, accounts[0]);

    const secondProof = proofs[1];
    const transitionProof = secondProof.value;
    const currentInstance = secondProof.input[0]
    const nextInstance = secondProof.input[1]

    await instanceManager.instantiate(transitionProof, currentInstance); // uses mock
    await instanceManager.transition(transitionProof, currentInstance, nextInstance);

    let isCurrentInstanceStored = await instanceManager.instances(currentInstance);
    assert.equal(isCurrentInstanceStored, false, "currentInstance was deleted");
    let isNextInstanceStored = await instanceManager.instances(nextInstance);
    assert.equal(isNextInstanceStored, true, "nextInstance was stored");
  });
});
