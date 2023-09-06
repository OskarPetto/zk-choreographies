import { v4 as uuid } from 'uuid';

export type FlowId = string;

export type ElementId = string;

export enum ElementType {
  START,
  END,
  TASK,
  XOR_SPLIT,
  XOR_JOIN,
  AND_SPLIT,
  AND_JOIN
}

export interface Element {
  id: ElementId;
  type: ElementType;
  name?: string;
  incomingFlows: FlowId[];
  outgoingFlows: FlowId[];
}

export type ModelId = string;

export interface Model {
  id: ModelId;
  flows: FlowId[];
  elements: Map<ElementId, Element>;
}

export function copyModel(model: Model): Model {
  const elements: Element[] = [...model.elements.values()].map(element => ({
    id: element.id,
    type: element.type,
    name: element.name,
    incomingFlows: [...element.incomingFlows],
    outgoingFlows: [...element.outgoingFlows],
  }));

  return {
    id: model.id,
    flows: [...model.flows],
    elements: new Map(elements.map(t => [t.id, t])),
  }
}

export function createModelId(): ModelId {
  return uuid();
}