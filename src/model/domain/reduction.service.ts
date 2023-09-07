import { Injectable } from '@nestjs/common';
import { Model, FlowId, Transition, TransitionType, copyModel } from './model';

@Injectable()
export class ReductionService {
  reduceModel(model: Model): Model {
    const newModel = copyModel(model);
    for (const transition of newModel.transitions.values()) {
      switch (transition.type) {
        case TransitionType.START:
        case TransitionType.END:
        case TransitionType.TASK:
          break;
        case TransitionType.XOR_SPLIT:
        case TransitionType.AND_JOIN:
          this.removeTransitionAndOutgoingFlows(newModel, transition);
          break;
        case TransitionType.XOR_JOIN:
        case TransitionType.AND_SPLIT:
          this.removeTransitionAndIncomingFlows(newModel, transition);
          break;
      }
    }
    return newModel;
  }

  private removeTransitionAndOutgoingFlows(
    model: Model,
    transitionToRemove: Transition,
  ) {
    for (const transition of model.transitions.values()) {
      const intersect = this.intersect(
        transition.incomingFlows,
        transitionToRemove.outgoingFlows,
      );
      if (intersect.length > 0) {
        transition.incomingFlows = this.setMinus(
          transition.incomingFlows,
          transitionToRemove.outgoingFlows,
        );
        transition.incomingFlows = this.union(
          transition.incomingFlows,
          transitionToRemove.incomingFlows,
        );
      }
    }
    model.flows = this.setMinus(model.flows, transitionToRemove.outgoingFlows);
    model.transitions.delete(transitionToRemove.id);
  }

  private removeTransitionAndIncomingFlows(
    model: Model,
    transitionToRemove: Transition,
  ) {
    for (const transition of model.transitions.values()) {
      const intersect = this.intersect(
        transition.outgoingFlows,
        transitionToRemove.incomingFlows,
      );
      if (intersect.length > 0) {
        transition.outgoingFlows = this.setMinus(
          transition.outgoingFlows,
          transitionToRemove.incomingFlows,
        );
        transition.outgoingFlows = this.union(
          transition.outgoingFlows,
          transitionToRemove.outgoingFlows,
        );
      }
    }
    model.flows = this.setMinus(model.flows, transitionToRemove.incomingFlows);
    model.transitions.delete(transitionToRemove.id);
  }

  private setMinus(flows1: FlowId[], flows2: FlowId[]): FlowId[] {
    return flows1.filter((flow1) => !flows2.includes(flow1));
  }

  private intersect(flows1: FlowId[], flows2: FlowId[]): FlowId[] {
    return flows1.filter((flow1) => flows2.includes(flow1));
  }

  private union(flows1: FlowId[], flows2: FlowId[]): FlowId[] {
    return [...new Set([...flows1, ...flows2])];
  }
}
