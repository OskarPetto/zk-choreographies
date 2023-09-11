import { Injectable } from '@nestjs/common';
import { Instance, copyInstance } from 'src/instance';
import { Model, Transition, TransitionId, TransitionType } from 'src/model';

@Injectable()
export class ExecutionService {
  executeTransitions(
    instance: Instance,
    transitionIds: TransitionId[],
  ): Instance {
    const transitions = transitionIds.map((transitionId) =>
      this.findTransition(instance.model, transitionId),
    );
    const newInstance = copyInstance(instance);
    for (const transition of transitions) {
      this.executeTransition(newInstance, transition);
    }
    return newInstance;
  }

  private executeTransition(instance: Instance, transition: Transition) {
    if (!this.isTransitionExecutable(instance, transition)) {
      throw Error(`Transition ${transition.id} is not executable`);
    }
    for (const incomingPlaceId of transition.incomingPlaces) {
      instance.tokenCounts[incomingPlaceId] -= 1;
    }
    for (const outgoingPlaceId of transition.outgoingPlaces) {
      instance.tokenCounts[outgoingPlaceId] += 1;
    }
    if (transition.type === TransitionType.END) {
      instance.finished = true;
    }
  }

  private isTransitionExecutable(instance: Instance, transition: Transition) {
    if (instance.finished) {
      return false;
    }
    return [...transition.incomingPlaces]
      .map((placeId) => instance.tokenCounts[placeId])
      .every((tokenCount) => tokenCount > 0);
  }

  private findTransition(model: Model, transitionId: TransitionId): Transition {
    const transition = model.transitions.find((t) => t.id === transitionId);
    if (!transition) {
      throw Error(`Transition ${transitionId} in model ${model.id} not found`);
    }
    return transition;
  }
}
