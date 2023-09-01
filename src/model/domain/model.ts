
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

export type PlaceId = number;

export interface Transition {
  id: TransitionId;
  type: TransitionType;
  fromPlaces: Set<PlaceId>;
  toPlaces: Set<PlaceId>;
}

export type ModelId = string;

export interface Model {
  id: ModelId;
  places: Set<PlaceId>;
  transitions: Map<TransitionId, Transition>;
}

export function findTransition(model: Model, transitionId: TransitionId): Transition {
  const transition = model.transitions.get(transitionId);
  if (!transition) {
    throw Error(`Transition ${transitionId} in model ${model.id} not found`);
  }
  return transition;
}

export function copyModel(model: Model): Model {
  return JSON.parse(JSON.stringify(model));
}