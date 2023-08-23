import { ModelId, PlaceId } from "src/model/domain/model";

export enum ExecutionStatus {
  NOT_ACTIVE,
  ACTIVE
}

export type ChoreographyId = string;

export interface Choreography {
  id: ChoreographyId;
  model: ModelId;
  executionStatuses: Map<PlaceId, ExecutionStatus>;
}