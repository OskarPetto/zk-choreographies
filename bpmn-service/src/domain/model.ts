import { Condition } from './condition';

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
  type?: TransitionType;
  name?: string;
  incomingPlaces: PlaceId[];
  outgoingPlaces: PlaceId[];
  initiatingParticipant?: ParticipantId;
  respondingParticipant?: ParticipantId;
  initiatingMessage?: MessageId;
  respondingMessage?: MessageId;
  condition?: Condition;
}

export interface Model {
  placeCount: number;
  participantCount: number;
  messageCount: number;
  startPlaces: PlaceId[];
  endPlaces: PlaceId[];
  transitions: Transition[];
}
