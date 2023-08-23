import { Injectable } from '@nestjs/common';
import { ChoreographyId, Choreography, ExecutionStatus } from './choreography';
import { Model, ModelService, Transition, TransitionId, TransitionType } from 'model';

@Injectable()
export class ChoreographyService {
    constructor(private modelService: ModelService) { }

    choreographies: Map<ChoreographyId, Choreography>;

    fire(choreographyId: ChoreographyId, transitionId: TransitionId) {
        const choreography = this.choreographies.get(choreographyId);
        const model = this.modelService.find(choreography.model);
        const transition = model.transitions.get(transitionId);
        this.deactivateFromPlaces(choreography, transition);
        this.activateToPlaces(choreography, model, transition);
    }

    private deactivateFromPlaces(choreography: Choreography, transition: Transition) {
        for (const fromPlace of transition.fromPlaces) {
            const exectionStatus = choreography.executionStatuses.get(fromPlace);
            if (exectionStatus === ExecutionStatus.NOT_ACTIVE) {
                throw Error(`Place ${fromPlace} is not active`);
            }
            choreography.executionStatuses.set(fromPlace, ExecutionStatus.NOT_ACTIVE);
        }
    }

    private activateToPlaces(choreography: Choreography, model: Model, transition: Transition) {
        for (const toPlace of transition.toPlaces) {
            choreography.executionStatuses.set(toPlace, ExecutionStatus.ACTIVE);
            for (const nextTransition of model.transitions.values()) {
                if (nextTransition.fromPlaces.has(toPlace) && nextTransition.type === TransitionType.AND_JOIN) {
                    break;
                }
            }
        }
    }
}
