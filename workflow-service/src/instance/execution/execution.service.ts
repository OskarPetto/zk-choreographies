import { Injectable } from '@nestjs/common';
import { Instance } from 'src/instance/instance';
import { PetriNet, Transition } from 'src/model/petri-net/petri-net';

@Injectable()
export class ExecutionService {
  instantiatePetriNet(petriNet: PetriNet): Instance {
    const tokenCounts = Array(petriNet.placeCount).fill(0);
    tokenCounts[petriNet.startPlace] = 1;
    return {
      petriNet: petriNet.id,
      tokenCounts,
    };
  }

  executeTransition(instance: Instance, transition: Transition): Instance {
    const newInstance = this.copyInstance(instance);
    this.executeTransitionInline(newInstance, transition);
    return newInstance;
  }

  private executeTransitionInline(instance: Instance, transition: Transition) {
    if (!this.isTransitionExecutable(instance, transition)) {
      throw Error(`Transition ${transition.id} is not executable`);
    }
    for (const incomingPlaceId of transition.incomingPlaces) {
      instance.tokenCounts[incomingPlaceId] -= 1;
    }
    for (const outgoingPlaceId of transition.outgoingPlaces) {
      instance.tokenCounts[outgoingPlaceId] += 1;
    }
  }

  private isTransitionExecutable(instance: Instance, transition: Transition) {
    return [...transition.incomingPlaces]
      .map((placeId) => instance.tokenCounts[placeId])
      .every((tokenCount) => tokenCount > 0);
  }

  private copyInstance(instance: Instance): Instance {
    return {
      id: instance.id,
      petriNet: instance.petriNet,
      tokenCounts: [...instance.tokenCounts],
    };
  }
}
