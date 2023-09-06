import { ExecutionStatus, Instance } from "src/choreography";
import { Model, Element, ElementType, FlowId } from "src/model";

export class TestdataProvider {
    static getModel1(): Model {
        const elements: Element[] = [
            { id: 'As', name: 'Application submitted', type: ElementType.START, incomingFlows: [], outgoingFlows: ['0'] },
            { id: 'Da1', name: 'Decline application', type: ElementType.TASK, incomingFlows: ['0'], outgoingFlows: ['8'] },
            { id: 'Aa', name: 'Accept application', type: ElementType.TASK, incomingFlows: ['0'], outgoingFlows: ['1', '3'] },
            { id: 'Fa', name: 'Finalize application', type: ElementType.TASK, incomingFlows: ['1'], outgoingFlows: ['2'] },
            { id: 'Sso', name: 'Select and send offer', type: ElementType.TASK, incomingFlows: ['3'], outgoingFlows: ['4'] },
            { id: 'Ro', name: 'Receive offer', type: ElementType.TASK, incomingFlows: ['4'], outgoingFlows: ['5'] },
            { id: 'Co', name: 'Cancel offer', type: ElementType.TASK, incomingFlows: ['5'], outgoingFlows: ['3'] },
            { id: 'Ao', name: 'Accept offer', type: ElementType.TASK, incomingFlows: ['2', '5'], outgoingFlows: ['6'] },
            { id: 'Aaa', name: 'Approve and activate application', type: ElementType.TASK, incomingFlows: ['6'], outgoingFlows: ['8'] },
            { id: 'Do', name: 'Decline offer', type: ElementType.TASK, incomingFlows: ['2', '5'], outgoingFlows: ['7'] },
            { id: 'Da2', name: 'Decline application', type: ElementType.TASK, incomingFlows: ['7'], outgoingFlows: ['8'] },
            { id: 'Af', name: 'Application finished', type: ElementType.END, incomingFlows: ['8'], outgoingFlows: [] },
        ];
        return {
            id: 'model1',
            flows: [...Array(9).keys()].map(i => '' + i),
            elements: new Map(elements.map(t => [t.id, t])),
        }
    }

    static getInstance1(): Instance {
        return {
            id: 'instance1',
            model: this.getModel1().id,
            executionStatuses: new Map([
                ['0', ExecutionStatus.NOT_ACTIVE],
                ['1', ExecutionStatus.NOT_ACTIVE],
                ['2', ExecutionStatus.NOT_ACTIVE],
                ['3', ExecutionStatus.NOT_ACTIVE],
                ['4', ExecutionStatus.NOT_ACTIVE],
                ['5', ExecutionStatus.NOT_ACTIVE],
                ['6', ExecutionStatus.NOT_ACTIVE],
                ['7', ExecutionStatus.NOT_ACTIVE],
                ['8', ExecutionStatus.NOT_ACTIVE]
            ]),
            finished: false
        };
    }

    static getModel2(): Model {
        const transitions: Element[] = [
            { id: 'As', name: 'Application submitted', type: ElementType.START, incomingFlows: [], outgoingFlows: ['0'] },
            { id: 'Da1', name: 'Decline application', type: ElementType.TASK, incomingFlows: ['1'], outgoingFlows: ['2'] },
            { id: 'Aa', name: 'Accept application', type: ElementType.TASK, incomingFlows: ['3'], outgoingFlows: ['4'] },
            { id: 'Fa', name: 'Finalize application', type: ElementType.TASK, incomingFlows: ['5'], outgoingFlows: ['6'] },
            { id: 'Sso', name: 'Select and send offer', type: ElementType.TASK, incomingFlows: ['9'], outgoingFlows: ['10'] },
            { id: 'Ro', name: 'Receive offer', type: ElementType.TASK, incomingFlows: ['10'], outgoingFlows: ['11'] },
            { id: 'Co', name: 'Cancel offer', type: ElementType.TASK, incomingFlows: ['12'], outgoingFlows: ['8'] },
            { id: 'Ao', name: 'Accept offer', type: ElementType.TASK, incomingFlows: ['15'], outgoingFlows: ['16'] },
            { id: 'Aaa', name: 'Approve and activate application', type: ElementType.TASK, incomingFlows: ['16'], outgoingFlows: ['17'] },
            { id: 'Do', name: 'Decline offer', type: ElementType.TASK, incomingFlows: ['18'], outgoingFlows: ['19'] },
            { id: 'Da2', name: 'Decline application', type: ElementType.TASK, incomingFlows: ['19'], outgoingFlows: ['20'] },
            { id: 'Af', name: 'Application finished', type: ElementType.END, incomingFlows: ['21'], outgoingFlows: [] },
            { id: 't0', type: ElementType.XOR_SPLIT, incomingFlows: ['0'], outgoingFlows: ['1'] },
            { id: 't1', type: ElementType.XOR_JOIN, incomingFlows: ['2'], outgoingFlows: ['21'] },
            { id: 't2', type: ElementType.XOR_SPLIT, incomingFlows: ['0'], outgoingFlows: ['3'] },
            { id: 't3', type: ElementType.AND_SPLIT, incomingFlows: ['4'], outgoingFlows: ['5', '7'] },
            { id: 't4', type: ElementType.XOR_JOIN, incomingFlows: ['7'], outgoingFlows: ['9'] },
            { id: 't5', type: ElementType.XOR_SPLIT, incomingFlows: ['11'], outgoingFlows: ['13'] },
            { id: 't6', type: ElementType.XOR_SPLIT, incomingFlows: ['11'], outgoingFlows: ['12'] },
            { id: 't7', type: ElementType.XOR_JOIN, incomingFlows: ['8'], outgoingFlows: ['9'] },
            { id: 't8', type: ElementType.AND_JOIN, incomingFlows: ['6', '13'], outgoingFlows: ['14'] },
            { id: 't9', type: ElementType.XOR_SPLIT, incomingFlows: ['14'], outgoingFlows: ['15'] },
            { id: 't10', type: ElementType.XOR_JOIN, incomingFlows: ['17'], outgoingFlows: ['21'] },
            { id: 't11', type: ElementType.XOR_SPLIT, incomingFlows: ['14'], outgoingFlows: ['18'] },
            { id: 't12', type: ElementType.XOR_JOIN, incomingFlows: ['20'], outgoingFlows: ['21'] },
        ];
        return {
            id: 'model1',
            flows: [...Array(22).keys()].map(i => '' + i),
            elements: new Map(transitions.map(t => [t.id, t])),
        }
    }
}
