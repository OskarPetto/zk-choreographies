import { Injectable } from '@nestjs/common';
import { Instance, ExecutionStatus, copyInstance } from './instance';
import { Transition, TransitionType, } from 'src/model';

@Injectable()
export class ChoreographyService {
    executeTransitions(instance: Instance, transitions: Transition[]): Instance {
        let newInstance = copyInstance(instance);
        for (const transition of transitions) {
            this.updateInstance(newInstance, transition);
        }
        return newInstance;
    }

    executeTransition(instance: Instance, transition: Transition): Instance {
        const newInstance = copyInstance(instance);
        this.updateInstance(newInstance, transition);
        return newInstance;
    }

    private updateInstance(instance: Instance, transition: Transition) {
        if (!this.isTransitionExecutable(instance, transition)) {
            throw Error(`Transition ${transition.id} cannot fire`);
        }
        for (const fromPlace of transition.fromPlaces) {
            instance.executionStatuses.set(fromPlace, ExecutionStatus.NOT_ACTIVE);
        }
        for (const toPlace of transition.toPlaces) {
            instance.executionStatuses.set(toPlace, ExecutionStatus.ACTIVE);
        }
        if (transition.type == TransitionType.END) {
            instance.finished = true;
        }
    }

    private isTransitionExecutable(instance: Instance, transition: Transition) {
        if (instance.finished) {
            return false;
        }
        return [...transition.fromPlaces]
            .map(placeId => instance.executionStatuses.get(placeId))
            .every(executionStatus => executionStatus === ExecutionStatus.ACTIVE);
    }
}
