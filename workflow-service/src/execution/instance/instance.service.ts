import { Injectable } from '@nestjs/common';
import { InstanceId, Instance } from './instance';
import { v4 as uuid } from 'uuid';

@Injectable()
export class InstanceService {
  instances: Map<InstanceId, Instance> = new Map();

  findInstance(instanceId: InstanceId): Instance {
    const instance = this.instances.get(instanceId);
    if (!instance) {
      throw Error(`Instance ${instanceId} not found`);
    }
    return instance;
  }

  saveInstance(instance: Instance) {
    if (!instance.id) {
      instance.id = this.createInstanceId();
    }
    this.instances.set(instance.id, instance);
  }

  private createInstanceId(): InstanceId {
    return uuid();
  }
}
