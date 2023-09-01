import { ExecutionStatus, Instance } from "src/choreography";
import { Model, TransitionType, indexTransitions } from "src/model";

export class TestdataProvider {
    static getModel1(): Model {
        const transitions = [
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
        ];
        return {
            id: 'model1',
            placeCount: 9,
            transitions,
            transitionIndex: indexTransitions(transitions)
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

    static getModel2(): Model {
        const transitions = [
            { id: 'As', type: TransitionType.START, fromPlaces: [], toPlaces: [0] },
            { id: 't0', type: TransitionType.XOR_SPLIT, fromPlaces: [0], toPlaces: [1] },
            { id: 'Da1', type: TransitionType.TASK, fromPlaces: [1], toPlaces: [2] },
            { id: 't1', type: TransitionType.XOR_JOIN, fromPlaces: [2], toPlaces: [21] },
            { id: 'Af', type: TransitionType.END, fromPlaces: [21], toPlaces: [] },
            { id: 't2', type: TransitionType.XOR_SPLIT, fromPlaces: [0], toPlaces: [3] },
            { id: 'Aa', type: TransitionType.TASK, fromPlaces: [3], toPlaces: [4] },
            { id: 't3', type: TransitionType.AND_SPLIT, fromPlaces: [4], toPlaces: [5, 7] },
            { id: 'Fa', type: TransitionType.TASK, fromPlaces: [5], toPlaces: [6] },
            { id: 't4', type: TransitionType.XOR_JOIN, fromPlaces: [7], toPlaces: [9] },
            { id: 'Sso', type: TransitionType.TASK, fromPlaces: [9], toPlaces: [10] },
            { id: 'Ro', type: TransitionType.TASK, fromPlaces: [10], toPlaces: [11] },
            { id: 't5', type: TransitionType.XOR_SPLIT, fromPlaces: [11], toPlaces: [13] },
            { id: 't6', type: TransitionType.XOR_SPLIT, fromPlaces: [11], toPlaces: [12] },
            { id: 'Co', type: TransitionType.TASK, fromPlaces: [12], toPlaces: [8] },
            { id: 't7', type: TransitionType.XOR_JOIN, fromPlaces: [8], toPlaces: [9] },
            { id: 't8', type: TransitionType.AND_JOIN, fromPlaces: [6, 13], toPlaces: [14] },
            { id: 't9', type: TransitionType.XOR_SPLIT, fromPlaces: [14], toPlaces: [15] },
            { id: 'Ao', type: TransitionType.TASK, fromPlaces: [15], toPlaces: [16] },
            { id: 'Aaa', type: TransitionType.TASK, fromPlaces: [16], toPlaces: [17] },
            { id: 't10', type: TransitionType.XOR_JOIN, fromPlaces: [17], toPlaces: [21] },
            { id: 't11', type: TransitionType.XOR_SPLIT, fromPlaces: [17], toPlaces: [21] },
            { id: 'Do', type: TransitionType.TASK, fromPlaces: [18], toPlaces: [19] },
            { id: 'Da2', type: TransitionType.TASK, fromPlaces: [19], toPlaces: [20] },
            { id: 't12', type: TransitionType.XOR_JOIN, fromPlaces: [20], toPlaces: [21] },
        ];
        return {
            id: 'model1',
            placeCount: 22,
            transitions,
            transitionIndex: indexTransitions(transitions)
        }
    }
}
