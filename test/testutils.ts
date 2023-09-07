import { FlowId, Model } from 'src/model';

export function findFlowMapping(
  model1: Model,
  model2: Model,
): Map<FlowId, FlowId> | undefined {
  if (model1.flows.length !== model2.flows.length) {
    return undefined;
  }
  if (model1.transitions.size !== model2.transitions.size) {
    return undefined;
  }
  const flowMapping = new Map<FlowId, FlowId>();
  for (const transition1 of model1.transitions.values()) {
    const transition2 = model2.transitions.get(transition1.id);
    if (!transition2 || transition1.name !== transition2.name) {
      return undefined;
    }
    if (transition1.incomingFlows.length !== transition2.incomingFlows.length) {
      return undefined;
    }
    for (const [index, flowId1] of transition1.incomingFlows.entries()) {
      const flowId2 = flowMapping.get(flowId1);
      if (flowId2 && !transition2.incomingFlows.includes(flowId2)) {
        return undefined;
      }
      flowMapping.set(flowId1, transition2.incomingFlows[index]);
    }
    if (transition1.outgoingFlows.length !== transition2.outgoingFlows.length) {
      return undefined;
    }
    for (const [index, flowId1] of transition1.outgoingFlows.entries()) {
      const flowId2 = flowMapping.get(flowId1);
      if (flowId2 && !transition2.outgoingFlows.includes(flowId2)) {
        return undefined;
      }
      flowMapping.set(flowId1, transition2.outgoingFlows[index]);
    }
  }
  return flowMapping;
}
