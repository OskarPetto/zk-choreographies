import { Injectable } from '@nestjs/common';
import { InstanceId, Instance, ExecutionStatus } from './instance';
import { Model } from 'src/model';
import { v4 as uuid } from 'uuid';

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

    instantiateModel(model: Model): Instance {
        return {
            id: this.createInstanceId(),
            model: model.id,
            executionStatuses: new Map(Array.from(model.places).map(placeId => [placeId, ExecutionStatus.NOT_ACTIVE])),
            finished: false
        };
    }

    saveInstance(instance: Instance) {
        this.instances.set(instance.id, instance);
    }

    private createInstanceId(): InstanceId {
        return uuid();
    }
}
