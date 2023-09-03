import { ModelId, PlaceId } from "src/model";

export enum ExecutionStatus {
  NOT_ACTIVE,
  ACTIVE
}

export type InstanceId = string;

export interface Instance {
  id: InstanceId;
  model: ModelId;
  executionStatuses: ExecutionStatus[];
  finished: boolean;
}

export function copyInstance(instance: Instance): Instance {
  return JSON.parse(JSON.stringify(instance));
}