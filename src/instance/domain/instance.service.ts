import { Injectable } from '@nestjs/common';
import {
  InstanceId,
  Instance,
  createInstanceId,
  ExecutionStatus,
} from './instance';
import { Model, PlaceId } from 'src/model';

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
    const places = this.collectPlaces(model);
    return {
      id: createInstanceId(),
      model: model.id,
      executionStatuses: new Map(
        places.map((placeId) => [placeId, ExecutionStatus.NOT_ACTIVE]),
      ),
      finished: false,
    };
  }

  private collectPlaces(model: Model): PlaceId[] {
    let places: PlaceId[] = [];
    for (const transition of model.transitions.values()) {
      places = [
        ...places,
        ...transition.incomingPlaces,
        ...transition.outgoingPlaces,
      ];
    }
    return [...new Set(places)];
  }
}
