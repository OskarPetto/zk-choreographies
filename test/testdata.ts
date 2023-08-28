import { ExecutionStatus, Instance } from "src/choreography";
import { Model, TransitionType } from "src/model";

export class Testdata {
    static getModel1(): Model {
        return {
            id: 'model1',
            placeCount: 9,
            transitions: [
                { id: 'As', type: TransitionType.START, fromPlaces: [], toPlaces: [0] },
                { id: 'Da1', type: TransitionType.TASK, fromPlaces: [0], toPlaces: [8] },
                { id: 'Aa', type: TransitionType.TASK, fromPlaces: [0], toPlaces: [1, 3] },
                { id: 'Fa', type: TransitionType.TASK, fromPlaces: [1], toPlaces: [2] },
                { id: 'Sso', type: TransitionType.TASK, fromPlaces: [3], toPlaces: [4] },
                { id: 'Ro', type: TransitionType.TASK, fromPlaces: [4], toPlaces: [5] },
                { id: 'Co', type: TransitionType.TASK, fromPlaces: [5], toPlaces: [3] },
                { id: 'Ao', type: TransitionType.TASK, fromPlaces: [2, 5], toPlaces: [6] },
                { id: 'Aaa', type: TransitionType.TASK, fromPlaces: [6], toPlaces: [8] },
                { id: 'Do', type: TransitionType.TASK, fromPlaces: [2, 5], toPlaces: [7] },
                { id: 'Da2', type: TransitionType.TASK, fromPlaces: [7], toPlaces: [8] },
                { id: 'Af', type: TransitionType.END, fromPlaces: [8], toPlaces: [] },
            ]
        }
    }

    static getInstance1(): Instance {
        return {
            id: 'instance1',
            model: this.getModel1().id,
            executionStatuses: Array(this.getModel1().placeCount).fill(ExecutionStatus.NOT_ACTIVE),
            finished: false
        };
    }
}
