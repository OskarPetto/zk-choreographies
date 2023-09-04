import { ExecutionStatus, Instance } from "src/choreography";
import { Model, Transition, TransitionType } from "src/model";

export class TestdataProvider {
    static getModel1(): Model {
        const transitions: Transition[] = [
            { id: 'As', type: TransitionType.START, fromFlows: [], toFlows: [0] },
            { id: 'Da1', type: TransitionType.TASK, fromFlows: [0], toFlows: [8] },
            { id: 'Aa', type: TransitionType.TASK, fromFlows: [0], toFlows: [1, 3] },
            { id: 'Fa', type: TransitionType.TASK, fromFlows: [1], toFlows: [2] },
            { id: 'Sso', type: TransitionType.TASK, fromFlows: [3], toFlows: [4] },
            { id: 'Ro', type: TransitionType.TASK, fromFlows: [4], toFlows: [5] },
            { id: 'Co', type: TransitionType.TASK, fromFlows: [5], toFlows: [3] },
            { id: 'Ao', type: TransitionType.TASK, fromFlows: [2, 5], toFlows: [6] },
            { id: 'Aaa', type: TransitionType.TASK, fromFlows: [6], toFlows: [8] },
            { id: 'Do', type: TransitionType.TASK, fromFlows: [2, 5], toFlows: [7] },
            { id: 'Da2', type: TransitionType.TASK, fromFlows: [7], toFlows: [8] },
            { id: 'Af', type: TransitionType.END, fromFlows: [8], toFlows: [] },
        ];
        return {
            id: 'model1',
            flowCount: 9,
            transitions: new Map(transitions.map(t => [t.id, t])),
        }
    }

    static getInstance1(): Instance {
        return {
            id: 'instance1',
            model: this.getModel1().id,
            executionStatuses: Array(this.getModel1().flowCount).fill(ExecutionStatus.NOT_ACTIVE),
            finished: false
        };
    }

    static getModel2(): Model {
        const transitions: Transition[] = [
            { id: 'As', type: TransitionType.START, fromFlows: [], toFlows: [0] },
            { id: 'Da1', type: TransitionType.TASK, fromFlows: [1], toFlows: [2] },
            { id: 'Aa', type: TransitionType.TASK, fromFlows: [3], toFlows: [4] },
            { id: 'Fa', type: TransitionType.TASK, fromFlows: [5], toFlows: [6] },
            { id: 'Sso', type: TransitionType.TASK, fromFlows: [9], toFlows: [10] },
            { id: 'Ro', type: TransitionType.TASK, fromFlows: [10], toFlows: [11] },
            { id: 'Co', type: TransitionType.TASK, fromFlows: [12], toFlows: [8] },
            { id: 'Ao', type: TransitionType.TASK, fromFlows: [15], toFlows: [16] },
            { id: 'Aaa', type: TransitionType.TASK, fromFlows: [16], toFlows: [17] },
            { id: 'Do', type: TransitionType.TASK, fromFlows: [18], toFlows: [19] },
            { id: 'Da2', type: TransitionType.TASK, fromFlows: [19], toFlows: [20] },
            { id: 'Af', type: TransitionType.END, fromFlows: [21], toFlows: [] },
            { id: 't0', type: TransitionType.XOR_SPLIT, fromFlows: [0], toFlows: [1] },
            { id: 't1', type: TransitionType.XOR_JOIN, fromFlows: [2], toFlows: [21] },
            { id: 't2', type: TransitionType.XOR_SPLIT, fromFlows: [0], toFlows: [3] },
            { id: 't3', type: TransitionType.AND_SPLIT, fromFlows: [4], toFlows: [5, 7] },
            { id: 't4', type: TransitionType.XOR_JOIN, fromFlows: [7], toFlows: [9] },
            { id: 't5', type: TransitionType.XOR_SPLIT, fromFlows: [11], toFlows: [13] },
            { id: 't6', type: TransitionType.XOR_SPLIT, fromFlows: [11], toFlows: [12] },
            { id: 't7', type: TransitionType.XOR_JOIN, fromFlows: [8], toFlows: [9] },
            { id: 't8', type: TransitionType.AND_JOIN, fromFlows: [6, 13], toFlows: [14] },
            { id: 't9', type: TransitionType.XOR_SPLIT, fromFlows: [14], toFlows: [15] },
            { id: 't10', type: TransitionType.XOR_JOIN, fromFlows: [17], toFlows: [21] },
            { id: 't11', type: TransitionType.XOR_SPLIT, fromFlows: [14], toFlows: [18] },
            { id: 't12', type: TransitionType.XOR_JOIN, fromFlows: [20], toFlows: [21] },
        ];
        return {
            id: 'model1',
            flowCount: 22,
            transitions: new Map(transitions.map(t => [t.id, t])),
        }
    }
}
