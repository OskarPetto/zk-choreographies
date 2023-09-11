import { Injectable } from '@nestjs/common';
import {
  InstanceId,
  Instance,
  createInstanceId,
  ExecutionStatus,
} from './instance';
import { Model, copyModel } from 'src/model';

@Injectable()
export class InstanceService {
  instances: Map<InstanceId, Instance>;

  findInstance(instanceId: InstanceId): Instance {
    const instance = this.instances.get(instanceId);
    if (!instance) {
      throw Error(`Instance ${instanceId} not found`);
    }
    return instance;
  }

  saveInstance(instance: Instance) {
    this.instances.set(instance.id, instance);
  }

  instantiateModel(model: Model): Instance {
    return {
      id: createInstanceId(),
      model: copyModel(model),
      executionStatuses: Array(model.placeCount).fill(ExecutionStatus.NOT_ACTIVE),
      finished: false,
    };
  }
}
