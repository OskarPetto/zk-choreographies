export type SequenceFlowId = string;

export interface SequenceFlow {
  id: SequenceFlowId;
  name?: string;
}

export interface StartEvent {
  id: string;
  name?: string;
  outgoing: SequenceFlowId;
}

export interface EndEvent {
  id: string;
  name?: string;
  incoming: SequenceFlowId;
}

export enum GatewayType {
  SPLIT,
  JOIN,
}

export interface ExclusiveGateway {
  id: string;
  type: GatewayType;
  incoming: SequenceFlowId[];
  outgoing: SequenceFlowId[];
}

export interface ParallelGateway {
  id: string;
  type: GatewayType;
  incoming: SequenceFlowId[];
  outgoing: SequenceFlowId[];
}

export type MessageId = string;

export interface Message {
  id: MessageId;
  name: string;
}

export type ParticipantId = string;

export interface Participant {
  id: ParticipantId;
  name: string;
}

export enum LoopType {
  STANDARD,
  MULTI_INSTANCE_SEQUENTIAL,
  MULTI_INSTANCE_PARALLEL,
}

export interface ChoreographyTask {
  id: string;
  name?: string;
  incoming: SequenceFlowId[];
  outgoing: SequenceFlowId;
  initiatingParticipant: ParticipantId;
  respondingParticipant: ParticipantId;
  initialMessage?: MessageId;
  responseMessage?: MessageId;
  loopType?: LoopType;
}

export interface Choreography {
  id: string;
  startEvents: StartEvent[];
  endEvents: EndEvent[];
  participants: Participant[];
  choreographyTasks: ChoreographyTask[];
  exclusiveGateways: ExclusiveGateway[];
  parallelGateways: ParallelGateway[];
  sequenceFlows: SequenceFlow[];
  messages: Message[];
}

export interface Definitions {
  choreographies: Choreography[];
}
