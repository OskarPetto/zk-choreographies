import { Model } from 'src/model';
import { v4 as uuid } from 'uuid';

export type InstanceId = string;

export interface Instance {
  id: InstanceId;
  model: Model;
  tokenCounts: number[];
  finished: boolean;
}

export function copyInstance(instance: Instance): Instance {
  return {
    id: instance.id,
    model: instance.model,
    finished: instance.finished,
    tokenCounts: [...instance.tokenCounts],
  };
}

export function createInstanceId(): InstanceId {
  return uuid();
}
