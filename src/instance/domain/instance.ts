import { FlowId, ModelId } from "src/model";
import { v4 as uuid } from 'uuid';

export enum ExecutionStatus {
  NOT_ACTIVE,
  ACTIVE
}

export type InstanceId = string;

export interface Instance {
  id: InstanceId;
  model: ModelId;
  executionStatuses: Map<FlowId, ExecutionStatus>;
  finished: boolean;
}

export function copyInstance(instance: Instance): Instance {
  return {
    id: instance.id,
    model: instance.model,
    finished: instance.finished,
    executionStatuses: new Map(instance.executionStatuses)
  }
}

export function createInstanceId(): InstanceId {
  return uuid();
}