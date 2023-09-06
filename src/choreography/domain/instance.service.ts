import { Injectable } from '@nestjs/common';
import { InstanceId, Instance, ExecutionStatus, creadeInstanceId } from './instance';
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

    saveInstance(instance: Instance) {
        this.instances.set(instance.id, instance);
    }

    instantiateModel(model: Model): Instance {
        return {
            id: creadeInstanceId(),
            model: model.id,
            executionStatuses: new Map(model.flows.map(flowId => [flowId, ExecutionStatus.NOT_ACTIVE])),
            finished: false
        };
    }
}
