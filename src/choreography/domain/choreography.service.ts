import { Injectable } from '@nestjs/common';
import { Instance, ExecutionStatus } from './instance';
import { Transition, TransitionType, } from 'src/model';

@Injectable()
export class ChoreographyService {
    executeTransition(instance: Instance, transition: Transition): Instance {
        if (!this.isTransitionExecutable(instance, transition)) {
            throw Error(`Transition ${transition.id} cannot fire`);
        }
        const newInstance: Instance = {
            ...instance,
        }
        this.setExecutionStatuses(newInstance, transition);
        if (transition.type == TransitionType.END) {
            newInstance.finished = true;
        }
        return newInstance;
    }

    private isTransitionExecutable(instance: Instance, transition: Transition) {
        if (instance.finished) {
            return false;
        }
        for (const fromPlace of transition.fromPlaces) {
            const exectionStatus = instance.executionStatuses[fromPlace];
            if (exectionStatus === ExecutionStatus.NOT_ACTIVE) {
                return false;
            }
        }
        return true;
    }

    private setExecutionStatuses(instance: Instance, transition: Transition) {
        for (const fromPlace of transition.fromPlaces) {
            instance.executionStatuses[fromPlace] = ExecutionStatus.NOT_ACTIVE;
        }
        for (const toPlace of transition.toPlaces) {
            instance.executionStatuses[toPlace] = ExecutionStatus.ACTIVE;
        }
    }
}
