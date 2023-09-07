import { Injectable } from '@nestjs/common';
import { Instance, ExecutionStatus, copyInstance } from 'src/instance';
import { Transition, TransitionType } from 'src/model';

@Injectable()
export class ExecutionService {
  executeTransitions(instance: Instance, transitions: Transition[]): Instance {
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
    for (const incomingFlowId of transition.incomingFlows) {
      instance.executionStatuses.set(
        incomingFlowId,
        ExecutionStatus.NOT_ACTIVE,
      );
    }
    for (const outgoingFlowId of transition.outgoingFlows) {
      instance.executionStatuses.set(outgoingFlowId, ExecutionStatus.ACTIVE);
    }
    if (transition.type == TransitionType.END) {
      instance.finished = true;
    }
  }

  private isTransitionExecutable(instance: Instance, transition: Transition) {
    if (instance.finished) {
      return false;
    }
    return [...transition.incomingFlows]
      .map((flowId) => instance.executionStatuses.get(flowId))
      .every((executionStatus) => executionStatus === ExecutionStatus.ACTIVE);
  }
}
