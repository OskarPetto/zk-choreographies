import { ExecutionStatus, Instance } from 'src/instance';
import { Model, Element, ElementType } from 'src/model';

export class TestdataProvider {
  static getModel1(): Model {
    const elements: Element[] = [
      {
        id: 'As',
        name: 'Application submitted',
        type: ElementType.START,
        incomingFlows: [],
        outgoingFlows: ['0'],
      },
      {
        id: 'Da1',
        name: 'Decline application',
        type: ElementType.TASK,
        incomingFlows: ['0'],
        outgoingFlows: ['8'],
      },
      {
        id: 'Aa',
        name: 'Accept application',
        type: ElementType.TASK,
        incomingFlows: ['0'],
        outgoingFlows: ['1', '3'],
      },
      {
        id: 'Fa',
        name: 'Finalize application',
        type: ElementType.TASK,
        incomingFlows: ['1'],
        outgoingFlows: ['2'],
      },
      {
        id: 'Sso',
        name: 'Select and send offer',
        type: ElementType.TASK,
        incomingFlows: ['3'],
        outgoingFlows: ['4'],
      },
      {
        id: 'Ro',
        name: 'Receive offer',
        type: ElementType.TASK,
        incomingFlows: ['4'],
        outgoingFlows: ['5'],
      },
      {
        id: 'Co',
        name: 'Cancel offer',
        type: ElementType.TASK,
        incomingFlows: ['5'],
        outgoingFlows: ['3'],
      },
      {
        id: 'Ao',
        name: 'Accept offer',
        type: ElementType.TASK,
        incomingFlows: ['2', '5'],
        outgoingFlows: ['6'],
      },
      {
        id: 'Aaa',
        name: 'Approve and activate application',
        type: ElementType.TASK,
        incomingFlows: ['6'],
        outgoingFlows: ['8'],
      },
      {
        id: 'Do',
        name: 'Decline offer',
        type: ElementType.TASK,
        incomingFlows: ['2', '5'],
        outgoingFlows: ['7'],
      },
      {
        id: 'Da2',
        name: 'Decline application',
        type: ElementType.TASK,
        incomingFlows: ['7'],
        outgoingFlows: ['8'],
      },
      {
        id: 'Af',
        name: 'Application finished',
        type: ElementType.END,
        incomingFlows: ['8'],
        outgoingFlows: [],
      },
    ];
    return {
      id: 'model1',
      flows: [...Array(9).keys()].map((i) => '' + i),
      elements: new Map(elements.map((t) => [t.id, t])),
    };
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
        ['8', ExecutionStatus.NOT_ACTIVE],
      ]),
      finished: false,
    };
  }

  static getModel2(): Model {
    const elements: Element[] = [
      {
        id: 'As',
        name: 'Application submitted',
        type: ElementType.START,
        incomingFlows: [],
        outgoingFlows: ['0'],
      },
      {
        id: 'Da1',
        name: 'Decline application',
        type: ElementType.TASK,
        incomingFlows: ['1'],
        outgoingFlows: ['2'],
      },
      {
        id: 'Aa',
        name: 'Accept application',
        type: ElementType.TASK,
        incomingFlows: ['3'],
        outgoingFlows: ['4'],
      },
      {
        id: 'Fa',
        name: 'Finalize application',
        type: ElementType.TASK,
        incomingFlows: ['5'],
        outgoingFlows: ['6'],
      },
      {
        id: 'Sso',
        name: 'Select and send offer',
        type: ElementType.TASK,
        incomingFlows: ['9'],
        outgoingFlows: ['10'],
      },
      {
        id: 'Ro',
        name: 'Receive offer',
        type: ElementType.TASK,
        incomingFlows: ['10'],
        outgoingFlows: ['11'],
      },
      {
        id: 'Co',
        name: 'Cancel offer',
        type: ElementType.TASK,
        incomingFlows: ['12'],
        outgoingFlows: ['8'],
      },
      {
        id: 'Ao',
        name: 'Accept offer',
        type: ElementType.TASK,
        incomingFlows: ['15'],
        outgoingFlows: ['16'],
      },
      {
        id: 'Aaa',
        name: 'Approve and activate application',
        type: ElementType.TASK,
        incomingFlows: ['16'],
        outgoingFlows: ['17'],
      },
      {
        id: 'Do',
        name: 'Decline offer',
        type: ElementType.TASK,
        incomingFlows: ['18'],
        outgoingFlows: ['19'],
      },
      {
        id: 'Da2',
        name: 'Decline application',
        type: ElementType.TASK,
        incomingFlows: ['19'],
        outgoingFlows: ['20'],
      },
      {
        id: 'Af',
        name: 'Application finished',
        type: ElementType.END,
        incomingFlows: ['21'],
        outgoingFlows: [],
      },
      {
        id: 't0',
        type: ElementType.XOR_SPLIT,
        incomingFlows: ['0'],
        outgoingFlows: ['1'],
      },
      {
        id: 't2',
        type: ElementType.XOR_SPLIT,
        incomingFlows: ['0'],
        outgoingFlows: ['3'],
      },
      {
        id: 't3',
        type: ElementType.AND_SPLIT,
        incomingFlows: ['4'],
        outgoingFlows: ['5', '7'],
      },
      {
        id: 't4',
        type: ElementType.XOR_JOIN,
        incomingFlows: ['7'],
        outgoingFlows: ['9'],
      },
      {
        id: 't7',
        type: ElementType.XOR_JOIN,
        incomingFlows: ['8'],
        outgoingFlows: ['9'],
      },
      {
        id: 't5',
        type: ElementType.XOR_SPLIT,
        incomingFlows: ['11'],
        outgoingFlows: ['13'],
      },
      {
        id: 't6',
        type: ElementType.XOR_SPLIT,
        incomingFlows: ['11'],
        outgoingFlows: ['12'],
      },
      {
        id: 't8',
        type: ElementType.AND_JOIN,
        incomingFlows: ['6', '13'],
        outgoingFlows: ['14'],
      },
      {
        id: 't9',
        type: ElementType.XOR_SPLIT,
        incomingFlows: ['14'],
        outgoingFlows: ['15'],
      },
      {
        id: 't11',
        type: ElementType.XOR_SPLIT,
        incomingFlows: ['14'],
        outgoingFlows: ['18'],
      },
      {
        id: 't1',
        type: ElementType.XOR_JOIN,
        incomingFlows: ['2'],
        outgoingFlows: ['21'],
      },
      {
        id: 't10',
        type: ElementType.XOR_JOIN,
        incomingFlows: ['17'],
        outgoingFlows: ['21'],
      },
      {
        id: 't12',
        type: ElementType.XOR_JOIN,
        incomingFlows: ['20'],
        outgoingFlows: ['21'],
      },
    ];
    return {
      id: 'model2',
      flows: [...Array(22).keys()].map((i) => '' + i),
      elements: new Map(elements.map((t) => [t.id, t])),
    };
  }

  static getModel3(): Model {
    const elements: Element[] = [
      {
        id: 'StartEvent_1t34b5j',
        name: 'Application submitted (As)',
        type: ElementType.START,
        incomingFlows: [],
        outgoingFlows: ['0'],
      },
      {
        id: 'Activity_1ujl89c',
        name: 'Decline application (Da1)',
        type: ElementType.TASK,
        incomingFlows: ['1'],
        outgoingFlows: ['2'],
      },
      {
        id: 'Activity_1f7vwaw',
        name: 'Accept application (Aa)',
        type: ElementType.TASK,
        incomingFlows: ['3'],
        outgoingFlows: ['4'],
      },
      {
        id: 'Activity_1srex89',
        name: 'Finalize application (Fa)',
        type: ElementType.TASK,
        incomingFlows: ['5'],
        outgoingFlows: ['6'],
      },
      {
        id: 'Activity_1tgknny',
        name: 'Select and send offer (Sso)',
        type: ElementType.TASK,
        incomingFlows: ['9'],
        outgoingFlows: ['10'],
      },
      {
        id: 'Activity_0r9antn',
        name: 'Receive offer (Ro)',
        type: ElementType.TASK,
        incomingFlows: ['10'],
        outgoingFlows: ['11'],
      },
      {
        id: 'Activity_1q88w1a',
        name: 'Cancel offer (Co)',
        type: ElementType.TASK,
        incomingFlows: ['12'],
        outgoingFlows: ['8'],
      },
      {
        id: 'Activity_1g33b4i',
        name: 'Accept offer (Ao)',
        type: ElementType.TASK,
        incomingFlows: ['15'],
        outgoingFlows: ['16'],
      },
      {
        id: 'Activity_0iyazhn',
        name: 'Approve and activate application (Aaa)',
        type: ElementType.TASK,
        incomingFlows: ['16'],
        outgoingFlows: ['17'],
      },
      {
        id: 'Activity_05537al',
        name: 'Decline offer (Do)',
        type: ElementType.TASK,
        incomingFlows: ['18'],
        outgoingFlows: ['19'],
      },
      {
        id: 'Activity_19r4cu3',
        name: 'Decline application (Da2)',
        type: ElementType.TASK,
        incomingFlows: ['19'],
        outgoingFlows: ['20'],
      },
      {
        id: 'Event_1bsroo5',
        name: 'Application finished (Af)',
        type: ElementType.END,
        incomingFlows: ['21'],
        outgoingFlows: [],
      },
      {
        id: 'Gateway_1o9s8fw_Flow_07cyta2',
        type: ElementType.XOR_SPLIT,
        incomingFlows: ['0'],
        outgoingFlows: ['1'],
      },
      {
        id: 'Gateway_1o9s8fw_Flow_1n5x723',
        type: ElementType.XOR_SPLIT,
        incomingFlows: ['0'],
        outgoingFlows: ['3'],
      },
      {
        id: 'Gateway_1rgq5gy',
        type: ElementType.AND_SPLIT,
        incomingFlows: ['4'],
        outgoingFlows: ['5', '7'],
      },
      {
        id: 'Gateway_1fwxfgu_Flow_0vs3ms0',
        type: ElementType.XOR_JOIN,
        incomingFlows: ['7'],
        outgoingFlows: ['9'],
      },
      {
        id: 'Gateway_1fwxfgu_Flow_1dlqvig',
        type: ElementType.XOR_JOIN,
        incomingFlows: ['8'],
        outgoingFlows: ['9'],
      },
      {
        id: 'Gateway_1way65i_Flow_0q0zlpw',
        type: ElementType.XOR_SPLIT,
        incomingFlows: ['11'],
        outgoingFlows: ['13'],
      },
      {
        id: 'Gateway_1way65i_Flow_0x7rbwv',
        type: ElementType.XOR_SPLIT,
        incomingFlows: ['11'],
        outgoingFlows: ['12'],
      },
      {
        id: 'Gateway_1nglicj',
        type: ElementType.AND_JOIN,
        incomingFlows: ['6', '13'],
        outgoingFlows: ['14'],
      },
      {
        id: 'Gateway_1vl4hvy_Flow_0oshibp',
        type: ElementType.XOR_SPLIT,
        incomingFlows: ['14'],
        outgoingFlows: ['15'],
      },
      {
        id: 'Gateway_1vl4hvy_Flow_1b9hzto',
        type: ElementType.XOR_SPLIT,
        incomingFlows: ['14'],
        outgoingFlows: ['18'],
      },
      {
        id: 'Gateway_1vzsa13_Flow_0ka5y4w',
        type: ElementType.XOR_JOIN,
        incomingFlows: ['2'],
        outgoingFlows: ['21'],
      },
      {
        id: 'Gateway_1vzsa13_Flow_03xp19s',
        type: ElementType.XOR_JOIN,
        incomingFlows: ['17'],
        outgoingFlows: ['21'],
      },
      {
        id: 'Gateway_1vzsa13_Flow_12yxkzr',
        type: ElementType.XOR_JOIN,
        incomingFlows: ['20'],
        outgoingFlows: ['21'],
      },
    ];
    return {
      id: 'conformance_example.bpmn',
      flows: [...Array(22).keys()].map((i) => '' + i),
      elements: new Map(elements.map((t) => [t.id, t])),
    };
  }
}
