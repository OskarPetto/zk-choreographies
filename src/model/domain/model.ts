
export type TransitionId = string;

export enum TransitionType {
  START,
  END,
  TASK,
  XOR_SPLIT,
  XOR_JOIN,
  AND_SPLIT,
  AND_JOIN
}

export type FlowId = number;

export interface Transition {
  id: TransitionId;
  type: TransitionType;
  fromFlows: FlowId[];
  toFlows: FlowId[];
}

export type ModelId = string;

export interface Model {
  id: ModelId;
  flowCount: number;
  transitions: Map<TransitionId, Transition>;
}

export function copyModel(model: Model): Model {
  const newTransitions: Transition[] = [...model.transitions.values()].map(transition => ({
    id: transition.id,
    type: transition.type,
    fromFlows: [...transition.fromFlows],
    toFlows: [...transition.toFlows]
  }))

  return {
    id: model.id,
    flowCount: model.flowCount,
    transitions: new Map(newTransitions.map(t => [t.id, t]))
  }
}

export function modelEquals(model1: Model, model2: Model): boolean {
  return findFlowMapping(model1, model2) !== undefined;
}

export function findFlowMapping(model1: Model, model2: Model): Map<FlowId, FlowId> | undefined {
  if (model1.flowCount != model2.flowCount) {
    return undefined;
  }
  if (model1.transitions.size != model2.transitions.size) {
    return undefined;
  }
  const flowMapping = new Map<FlowId, FlowId>();
  for (const transition1 of model1.transitions.values()) {
    const transition2 = model2.transitions.get(transition1.id);
    if (!transition2) {
      return undefined;
    }
    const allFlows1 = [...transition1.fromFlows, ...transition1.toFlows];
    const allFlows2 = [...transition2.fromFlows, ...transition2.toFlows];
    if (allFlows1.length != allFlows2.length) {
      return undefined;
    }
    for (const [index, flowId1] of allFlows1.entries()) {
      const flowId2 = allFlows2[index];
      if (flowMapping.has(flowId1) && flowMapping.get(flowId1) != flowId2) {
        return undefined;
      }
      flowMapping.set(flowId1, flowId2);
    }
  }
  return flowMapping;
}