
export type TransitionId = string;

export enum TransitionType {
  START,
  END,
  TASK,
  // XOR_SPLIT,
  // XOR_JOIN,
  // AND_SPLIT,
  // AND_JOIN
}

export type PlaceId = number;

export interface Transition {
  id: TransitionId;
  type: TransitionType;
  fromPlaces: PlaceId[];
  toPlaces: PlaceId[];
}

export type ModelId = string;

export interface Model {
  id: ModelId;
  placeCount: number;
  transitions: Transition[];
}

export function findTransition(model: Model, transitionId: TransitionId) {
  const transition = model.transitions.find(transition => transition.id === transitionId);
  if (!transition) {
    throw Error(`Transition ${transitionId} in model ${model.id} not found`);
  }
  return transition;
}