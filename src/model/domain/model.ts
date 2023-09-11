import { v4 as uuid } from 'uuid';

export type PlaceId = number;

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
  placeCount: number;
  transitions: Transition[];
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
    placeCount: model.placeCount,
    transitions,
  };
}

export function createModelId(): ModelId {
  return uuid();
}
