import { Injectable } from '@nestjs/common';
import { ChoreographyId, Choreography, ExecutionStatus } from './choreography';
import { Model, Transition, TransitionId, TransitionType } from 'src/model/domain/model';
import { ModelService } from 'src/model/domain/model.service';

@Injectable()
export class ChoreographyService {
    constructor(private modelService: ModelService) { }

    choreographies: Map<ChoreographyId, Choreography>;

    fire(choreographyId: ChoreographyId, transitionId: TransitionId) {
        const choreography = this.choreographies.get(choreographyId);
        if (!choreography) {
            throw Error(`Choreography ${choreographyId} not found`);
        }
        const model = this.modelService.find(choreography.model);
        if (!model) {
            throw Error(`Model ${choreography.model} not found`);
        }
        const transition = model.transitions.get(transitionId);
        if (!transition) {
            throw Error(`Transition ${transitionId} in choreography ${choreographyId} not found`);
        }
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
