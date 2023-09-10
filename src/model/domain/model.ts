import { v4 as uuid } from 'uuid';

export type PlaceId = string;

export type TransitionId = string;

export enum TransitionType {
  START,
  END,
  TASK,
  XOR_SPLIT,
  XOR_JOIN,
  AND_SPLIT,
  AND_JOIN,
}

export interface Transition {
  id: TransitionId;
  type: TransitionType;
  name?: string;
  incomingPlaces: PlaceId[];
  outgoingPlaces: PlaceId[];
}

export type ModelId = string;

export interface Model {
  id: ModelId;
  transitions: Map<TransitionId, Transition>;
}

export function copyModel(model: Model): Model {
  const transitions: Transition[] = [...model.transitions.values()].map(
    (transition) => ({
      id: transition.id,
      type: transition.type,
      name: transition.name,
      incomingPlaces: [...transition.incomingPlaces],
      outgoingPlaces: [...transition.outgoingPlaces],
    }),
  );

  return {
    id: model.id,
    transitions: new Map(transitions.map((t) => [t.id, t])),
  };
}

export function createModelId(): ModelId {
  return uuid();
}
