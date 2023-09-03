import { Model, ModelId } from "src/model";
import { v4 as uuid } from 'uuid';

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

export function instantiateModel(model: Model): Instance {
  return {
    id: creadeId(),
    model: model.id,
    executionStatuses: Array(model.placeCount).fill(ExecutionStatus.NOT_ACTIVE),
    finished: false
  };
}

function creadeId(): InstanceId {
  return uuid();
}
