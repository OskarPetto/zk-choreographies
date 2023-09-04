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
                    this.removeTransitionAndToFlows(newModel, transition);
                    break;
                case TransitionType.XOR_JOIN:
                case TransitionType.AND_SPLIT:
                    this.removeTransitionAndFromFlows(newModel, transition);
                    break;
            }
        }
        this.updateFlowCount(newModel)
        return newModel;
    }

    private updateFlowCount(model: Model) {
        let flows: FlowId[] = [];
        for (const transition of model.transitions.values()) {
            flows = [...flows, ...transition.fromFlows, ...transition.toFlows];
        }
        model.flowCount = new Set(flows).size;
    }

    private removeTransitionAndToFlows(model: Model, transitionToRemove: Transition) {
        for (const transition of model.transitions.values()) {
            const intersect = this.intersect(transition.fromFlows, transitionToRemove.toFlows);
            if (intersect.length > 0) {
                transition.fromFlows = this.setMinus(transition.fromFlows, transitionToRemove.toFlows);
                transition.fromFlows = this.union(transition.fromFlows, transitionToRemove.fromFlows);
            }
        }
        model.transitions.delete(transitionToRemove.id);
    }

    private removeTransitionAndFromFlows(model: Model, transitionToRemove: Transition) {
        for (const transition of model.transitions.values()) {
            const intersect = this.intersect(transition.toFlows, transitionToRemove.fromFlows);
            if (intersect.length > 0) {
                transition.toFlows = this.setMinus(transition.toFlows, transitionToRemove.fromFlows);
                transition.toFlows = this.union(transition.toFlows, transitionToRemove.toFlows);
            }
        }
        model.transitions.delete(transitionToRemove.id);
    }

    private setMinus(flows1: FlowId[], flows2: FlowId[]): FlowId[] {
        return flows1.filter(flow1 => !flows2.includes(flow1));
    }

    private intersect(flows1: FlowId[], flows2: FlowId[]): FlowId[] {
        return flows1.filter(flow1 => flows2.includes(flow1));
    }

    private union(flows1: FlowId[], flows2: FlowId[]): FlowId[] {
        return [...new Set([...flows1, ...flows2])];
    }
}
