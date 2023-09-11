import { TransitionType } from '../domain/model';

export type SequenceFlowId = string;

export interface SequenceFlow {
  id: SequenceFlowId;
}

export interface Element {
  id: string;
  type: TransitionType;
  name?: string;
  incomingSequenceFlows: SequenceFlowId[];
  outgoingSequenceFlows: SequenceFlowId[];
}

export interface Process {
  id: string;
  startEvent: Element;
  endEvents: Element[];
  tasks: Element[];
  exclusiveGateways: Element[];
  parallelGateways: Element[];
  sequenceFlows: SequenceFlow[];
}

export interface Definitions {
  process: Process;
}
