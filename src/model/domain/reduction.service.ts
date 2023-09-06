import { Injectable } from '@nestjs/common';
import { Model, FlowId, Element, ElementType, copyModel } from './model';

@Injectable()
export class ReductionService {
    reduceModel(model: Model): Model {
        const newModel = copyModel(model);
        for (const element of newModel.elements.values()) {
            switch (element.type) {
                case ElementType.START:
                case ElementType.END:
                case ElementType.TASK:
                    break;
                case ElementType.XOR_SPLIT:
                case ElementType.AND_JOIN:
                    this.removeElementAndOutgoingFlows(newModel, element);
                    break;
                case ElementType.XOR_JOIN:
                case ElementType.AND_SPLIT:
                    this.removeElementAndIncomingFlows(newModel, element);
                    break;
            }
        }
        return newModel;
    }

    private removeElementAndOutgoingFlows(model: Model, elementToRemove: Element) {
        for (const element of model.elements.values()) {
            const intersect = this.intersect(element.incomingFlows, elementToRemove.outgoingFlows);
            if (intersect.length > 0) {
                element.incomingFlows = this.setMinus(element.incomingFlows, elementToRemove.outgoingFlows);
                element.incomingFlows = this.union(element.incomingFlows, elementToRemove.incomingFlows);
            }
        }
        model.flows = this.setMinus(model.flows, elementToRemove.outgoingFlows);
        model.elements.delete(elementToRemove.id);
    }

    private removeElementAndIncomingFlows(model: Model, elementToRemove: Element) {
        for (const element of model.elements.values()) {
            const intersect = this.intersect(element.outgoingFlows, elementToRemove.incomingFlows);
            if (intersect.length > 0) {
                element.outgoingFlows = this.setMinus(element.outgoingFlows, elementToRemove.incomingFlows);
                element.outgoingFlows = this.union(element.outgoingFlows, elementToRemove.outgoingFlows);
            }
        }
        model.flows = this.setMinus(model.flows, elementToRemove.incomingFlows);
        model.elements.delete(elementToRemove.id);
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

    findFlowMapping(model1: Model, model2: Model): Map<FlowId, FlowId> | undefined {
        if (model1.flows.length !== model2.flows.length) {
            return undefined;
        }
        if (model1.elements.size !== model2.elements.size) {
            return undefined;
        }
        const flowMapping = new Map<FlowId, FlowId>();
        for (const element1 of model1.elements.values()) {
            const element2 = model2.elements.get(element1.id);
            if (!element2) {
                return undefined;
            }
            const allFlows1 = [...element1.incomingFlows, ...element1.outgoingFlows];
            const allFlows2 = [...element2.incomingFlows, ...element2.outgoingFlows];
            if (allFlows1.length !== allFlows2.length) {
                return undefined;
            }
            for (const [index, flowId1] of allFlows1.entries()) {
                const flowId2 = allFlows2[index];
                if (flowMapping.has(flowId1) && flowMapping.get(flowId1) !== flowId2) {
                    return undefined;
                }
                flowMapping.set(flowId1, flowId2);
            }
        }
        return flowMapping;
    }
}
