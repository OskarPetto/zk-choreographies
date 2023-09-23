export type PlaceId = number;
export type ParticipantId = number;
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

export interface Message {

}

export interface Transition {
  id: TransitionId;
  type: TransitionType;
  name?: string;
  participant?: ParticipantId;
  incomingPlaces: PlaceId[];
  outgoingPlaces: PlaceId[];
}

export type ModelId = string;

export interface Model {
  id: ModelId;
  placeCount: number;
  startPlace: PlaceId;
  endPlace: PlaceId;
  participantCount: number;
  transitions: Transition[];
}
