import { Injectable } from '@nestjs/common';
import { InstanceId, Instance, ExecutionStatus } from './instance';
import { ModelId, ModelService, PlaceId } from 'src/model';
import { v4 as uuid } from 'uuid';

@Injectable()
export class InstanceService {
    constructor(private modelService: ModelService) { }

    instances: Map<InstanceId, Instance>;

    findInstance(instanceId: InstanceId): Instance {
        const instance = this.instances.get(instanceId);
        if (!instance) {
            throw Error(`Instance ${instanceId} not found`);
        }
        return instance;
    }

    instantiateModel(modelId: ModelId): Instance {
        const model = this.modelService.findModel(modelId);
        const InstanceId = this.createInstanceId();
        const mapPlaceFunction: (_: PlaceId) => [PlaceId, ExecutionStatus] = placeId => [placeId, ExecutionStatus.NOT_ACTIVE];
        return {
            id: InstanceId,
            model: model.id,
            executionStatuses: new Map(model.places.map(mapPlaceFunction))
        };
    }

    saveInstance(instance: Instance) {
        this.instances.set(instance.id, instance);
    }

    private createInstanceId(): InstanceId {
        return uuid();
    }
}
