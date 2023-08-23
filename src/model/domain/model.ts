
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
  transitions: Map<TransitionId, Transition>;
}