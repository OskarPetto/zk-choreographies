import { ignition } from 'hardhat';
import InstanceManagerModule from '../ignition/modules/InstanceManager';
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";

import * as fs from 'fs';
import { expect } from 'chai';
import { buildModule } from '@nomicfoundation/ignition-core';

const InstanceManagerWithMockModule = buildModule("InstanceManagerWithMockModule", (m) => {
  const instantiationVerifier = m.contract("InstantiationVerifierMock");
  const transitionVerifier = m.contract("TransitionVerifier");
  const terminationVerifier = m.contract("TerminationVerifier");

  const instanceManager = m.contract("InstanceManager", [instantiationVerifier, transitionVerifier, terminationVerifier]);

  return { instanceManager };
});

describe('InstanceManager', function () {

  let proofs: any;
  this.beforeAll(function() {
    const proofsBuffer = fs.readFileSync('test/proofs.json');
    proofs = JSON.parse(proofsBuffer.toString());
  })

  it('instantiation', async function () {
    const { instanceManager } = await ignition.deploy(InstanceManagerModule);

    const firstProof = proofs[0];
    const instantiationProof = firstProof.value;
    const instance = firstProof.input[0]

    await instanceManager.instantiate(instantiationProof, instance);

    const isStored = await instanceManager.instances(instance);
    expect(isStored).to.be.true;
  });

  it('termination', async function () {
    const { instanceManager } = await ignition.deploy(InstanceManagerWithMockModule);

    const lastProof = proofs[proofs.length - 1];
    const terminationProof = lastProof.value;
    const instance = lastProof.input[0]

    await instanceManager.instantiate(terminationProof, instance); // uses mock
    await instanceManager.terminate(terminationProof, instance);

    const isStored = await instanceManager.instances(instance);
    expect(isStored).to.be.false;
  });

  it('transition', async function () {
    const { instanceManager } = await ignition.deploy(InstanceManagerWithMockModule);

    const secondProof = proofs[1];
    const transitionProof = secondProof.value;
    const currentInstance = secondProof.input[0]
    const nextInstance = secondProof.input[1]

    await instanceManager.instantiate(transitionProof, currentInstance); // uses mock
    await instanceManager.transition(transitionProof, currentInstance, nextInstance);
    
    expect(await instanceManager.instances(currentInstance)).to.be.false;
    expect(await instanceManager.instances(nextInstance)).to.be.true;
  });
});
