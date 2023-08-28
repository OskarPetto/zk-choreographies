import { Injectable } from '@nestjs/common';
import { InstanceId, Instance, ExecutionStatus } from './instance';
import { InstanceService } from './instance.service';
import { Model, Transition, TransitionId, TransitionType, ModelService, } from 'src/model';

@Injectable()
export class ChoreographyService {
    constructor(
        private instanceService: InstanceService,
        private modelService: ModelService
    ) { }

    executeTransition(instanceId: InstanceId, transitionId: TransitionId) {
        const instance = this.instanceService.findInstance(instanceId);
        const model = this.modelService.findModel(instance.model);
        const transition = this.modelService.findTransition(model, transitionId);
        if (!this.canTransitionFire(instance, transition)) {
            throw Error(`Transition ${transitionId} cannot fire`);
        }
        this.deactivateFromPlaces(instance, transition);
        this.activateToPlaces(instance, transition);
    }

    private canTransitionFire(Instance: Instance, transition: Transition) {
        for (const fromPlace of transition.fromPlaces) {
            const exectionStatus = Instance.executionStatuses.get(fromPlace);
            if (exectionStatus === ExecutionStatus.NOT_ACTIVE) {
                return false;
            }
        }
        return true;
    }


    private deactivateFromPlaces(instance: Instance, transition: Transition) {
        for (const fromPlace of transition.fromPlaces) {
            instance.executionStatuses.set(fromPlace, ExecutionStatus.NOT_ACTIVE);
        }
    }

    private activateToPlaces(instance: Instance, transition: Transition) {
        for (const toPlace of transition.toPlaces) {
            instance.executionStatuses.set(toPlace, ExecutionStatus.ACTIVE);
        }
    }
}
