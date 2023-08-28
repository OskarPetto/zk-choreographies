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
}