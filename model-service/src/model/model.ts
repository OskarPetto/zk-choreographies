export type PlaceId = number;
export type ParticipantId = number;
export type MessageId = number;
export type TransitionId = string;
export type ModelId = string;

export enum TransitionType {
  REQUIRED,
  OPTIONAL_INCOMING,
  OPTIONAL_OUTGOING,
}

export interface Transition {
  id: TransitionId;
  type: TransitionType;
  name?: string;
  incomingPlaces: PlaceId[];
  outgoingPlaces: PlaceId[];
  participant?: ParticipantId;
  message?: MessageId;
}

export interface Model {
  id: ModelId;
  placeCount: number;
  participantCount: number;
  messageCount: number;
  startPlace: PlaceId;
  endPlaces: PlaceId[];
  transitions: Transition[];
}
