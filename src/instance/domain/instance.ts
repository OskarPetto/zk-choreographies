import { ModelId } from 'src/model';
import { v4 as uuid } from 'uuid';

export type InstanceId = string;

export interface Instance {
  id: InstanceId;
  model: ModelId;
  tokenCounts: number[];
}

export function copyInstance(instance: Instance): Instance {
  return {
    id: instance.id,
    model: instance.model,
    tokenCounts: [...instance.tokenCounts],
  };
}

export function createInstanceId(): InstanceId {
  return uuid();
}
