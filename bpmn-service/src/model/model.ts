import { Constraint } from 'src/constraint/constraint';

export type PlaceId = number;
export type ParticipantId = number;
export type MessageId = number;
export type TransitionId = string;

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
  constraint?: Constraint;
}

export interface Model {
  hash: { value: string; salt: string };
  choreography: string;
  placeCount: number;
  participantCount: number;
  messageCount: number;
  startPlaces: PlaceId[];
  endPlaces: PlaceId[];
  transitions: Transition[];
}
