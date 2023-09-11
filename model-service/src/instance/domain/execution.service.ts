import { Injectable } from '@nestjs/common';
import { Instance, copyInstance } from 'src/instance';
import { Transition } from 'src/model';

@Injectable()
export class ExecutionService {
  executeTransition(instance: Instance, transition: Transition): Instance {
    const newInstance = copyInstance(instance);
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
}
