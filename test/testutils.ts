import { FlowId, Model } from "src/model";

export function findFlowMapping(model1: Model, model2: Model): Map<FlowId, FlowId> | undefined {
    if (model1.flows.length !== model2.flows.length) {
        return undefined;
    }
    if (model1.elements.size !== model2.elements.size) {
        return undefined;
    }
    const flowMapping = new Map<FlowId, FlowId>();
    for (const element1 of model1.elements.values()) {
        const element2 = model2.elements.get(element1.id);
        if (!element2 || element1.name !== element2.name) {
            return undefined;
        }
        if (element1.incomingFlows.length !== element2.incomingFlows.length) {
            return undefined;
        }
        for (const [index, flowId1] of element1.incomingFlows.entries()) {
            const flowId2 = flowMapping.get(flowId1);
            if (flowId2 && !element2.incomingFlows.includes(flowId2)) {
                return undefined;
            }
            flowMapping.set(flowId1, element2.incomingFlows[index]);
        }
        if (element1.outgoingFlows.length !== element2.outgoingFlows.length) {
            return undefined;
        }
        for (const [index, flowId1] of element1.outgoingFlows.entries()) {
            const flowId2 = flowMapping.get(flowId1);
            if (flowId2 && !element2.outgoingFlows.includes(flowId2)) {
                return undefined;
            }
            flowMapping.set(flowId1, element2.outgoingFlows[index]);
        }
    }
    return flowMapping;
}