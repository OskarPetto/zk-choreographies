import * as fs from 'fs';
import * as path from 'path';
import { ExecutionStatus, Instance } from 'src/instance';
import { Model, Transition, TransitionType } from 'src/model';
import { Process } from 'src/model/bpmn/bpmn';

function readTextFile(filename: string) {
  const filePath = path.join(process.cwd(), filename);
  return fs.readFileSync(filePath, 'utf-8').toString();
}

export class TestdataProvider {
  static getConformanceExample(): string {
    return readTextFile('test/data/conformance_example.bpmn');
  }

  static getProcess1(): Process {
    return {
      id: 'conformance_example',
      startEvent: {
        id: 'As',
        type: TransitionType.START,
        name: 'Application submitted (As)',
        incomingSequenceFlows: [],
        outgoingSequenceFlows: ['Flow_0xqg6as'],
      },
      endEvents: [
        {
          id: 'Af',
          type: TransitionType.END,
          name: 'Application finished (Af)',
          incomingSequenceFlows: ['Flow_1qh0ya1'],
          outgoingSequenceFlows: [],
        },
      ],
      tasks: [
        {
          id: 'Da1',
          type: TransitionType.TASK,
          name: 'Decline application (Da1)',
          incomingSequenceFlows: ['Flow_07cyta2'],
          outgoingSequenceFlows: ['Flow_0ka5y4w'],
        },
        {
          id: 'Aa',
          type: TransitionType.TASK,
          name: 'Accept application (Aa)',
          incomingSequenceFlows: ['Flow_1n5x723'],
          outgoingSequenceFlows: ['Flow_0w362at'],
        },
        {
          id: 'Fa',
          type: TransitionType.TASK,
          name: 'Finalize application (Fa)',
          incomingSequenceFlows: ['Flow_0i668iz'],
          outgoingSequenceFlows: ['Flow_1yee6jg'],
        },
        {
          id: 'Sso',
          type: TransitionType.TASK,
          name: 'Select and send offer (Sso)',
          incomingSequenceFlows: ['Flow_1g76mn1'],
          outgoingSequenceFlows: ['Flow_19w6dwb'],
        },
        {
          id: 'Ro',
          type: TransitionType.TASK,
          name: 'Receive offer (Ro)',
          incomingSequenceFlows: ['Flow_19w6dwb'],
          outgoingSequenceFlows: ['Flow_0damgea'],
        },
        {
          id: 'Co',
          type: TransitionType.TASK,
          name: 'Cancel offer (Co)',
          incomingSequenceFlows: ['Flow_0x7rbwv'],
          outgoingSequenceFlows: ['Flow_1dlqvig'],
        },
        {
          id: 'Ao',
          type: TransitionType.TASK,
          name: 'Accept offer (Ao)',
          incomingSequenceFlows: ['Flow_0oshibp'],
          outgoingSequenceFlows: ['Flow_0u5u7d3'],
        },
        {
          id: 'Do',
          type: TransitionType.TASK,
          name: 'Decline offer (Do)',
          incomingSequenceFlows: ['Flow_1b9hzto'],
          outgoingSequenceFlows: ['Flow_18msx9i'],
        },
        {
          id: 'Aaa',
          type: TransitionType.TASK,
          name: 'Approve and activate application (Aaa)',
          incomingSequenceFlows: ['Flow_0u5u7d3'],
          outgoingSequenceFlows: ['Flow_03xp19s'],
        },
        {
          id: 'Da2',
          type: TransitionType.TASK,
          name: 'Decline application (Da2)',
          incomingSequenceFlows: ['Flow_18msx9i'],
          outgoingSequenceFlows: ['Flow_12yxkzr'],
        },
      ],
      exclusiveGateways: [
        {
          id: 'Gateway_1o9s8fw',
          type: TransitionType.XOR_SPLIT,
          incomingSequenceFlows: ['Flow_0xqg6as'],
          outgoingSequenceFlows: ['Flow_07cyta2', 'Flow_1n5x723'],
        },
        {
          id: 'Gateway_1fwxfgu',
          type: TransitionType.XOR_JOIN,
          incomingSequenceFlows: ['Flow_0vs3ms0', 'Flow_1dlqvig'],
          outgoingSequenceFlows: ['Flow_1g76mn1'],
        },
        {
          id: 'Gateway_1way65i',
          type: TransitionType.XOR_SPLIT,
          incomingSequenceFlows: ['Flow_0damgea'],
          outgoingSequenceFlows: ['Flow_0q0zlpw', 'Flow_0x7rbwv'],
        },
        {
          id: 'Gateway_1vl4hvy',
          type: TransitionType.XOR_SPLIT,
          incomingSequenceFlows: ['Flow_0jf7vbw'],
          outgoingSequenceFlows: ['Flow_0oshibp', 'Flow_1b9hzto'],
        },
        {
          id: 'Gateway_1vzsa13',
          type: TransitionType.XOR_JOIN,
          incomingSequenceFlows: [
            'Flow_03xp19s',
            'Flow_12yxkzr',
            'Flow_0ka5y4w',
          ],
          outgoingSequenceFlows: ['Flow_1qh0ya1'],
        },
      ],
      parallelGateways: [
        {
          id: 'Gateway_1rgq5gy',
          type: TransitionType.AND_SPLIT,
          incomingSequenceFlows: ['Flow_0w362at'],
          outgoingSequenceFlows: ['Flow_0i668iz', 'Flow_0vs3ms0'],
        },
        {
          id: 'Gateway_1nglicj',
          type: TransitionType.AND_JOIN,
          incomingSequenceFlows: ['Flow_0q0zlpw', 'Flow_1yee6jg'],
          outgoingSequenceFlows: ['Flow_0jf7vbw'],
        },
      ],
      sequenceFlows: [
        {
          id: 'Flow_0xqg6as'
        },
        {
          id: 'Flow_07cyta2'
        },
        {
          id: 'Flow_1n5x723'
        },
        {
          id: 'Flow_0w362at'
        },
        {
          id: 'Flow_0i668iz'
        },
        {
          id: 'Flow_0vs3ms0'
        },
        {
          id: 'Flow_1g76mn1'
        },
        {
          id: 'Flow_19w6dwb'
        },
        {
          id: 'Flow_0damgea'
        },
        {
          id: 'Flow_0q0zlpw'
        },
        {
          id: 'Flow_1yee6jg'
        },
        {
          id: 'Flow_0x7rbwv'
        },
        {
          id: 'Flow_1dlqvig'
        },
        {
          id: 'Flow_0jf7vbw'
        },
        {
          id: 'Flow_0oshibp'
        },
        {
          id: 'Flow_1b9hzto'
        },
        {
          id: 'Flow_0u5u7d3'
        },
        {
          id: 'Flow_18msx9i'
        },
        {
          id: 'Flow_03xp19s'
        },
        {
          id: 'Flow_12yxkzr'
        },
        {
          id: 'Flow_1qh0ya1'
        },
        {
          id: 'Flow_0ka5y4w'
        }
      ]
    };
  }

  static getModel1(): Model {
    const transitions: Transition[] = [
      {
        id: 'As',
        name: 'Application submitted (As)',
        type: TransitionType.START,
        incomingPlaces: [],
        outgoingPlaces: [0],
      },
      {
        id: 'Gateway_1o9s8fw_Flow_07cyta2',
        type: TransitionType.XOR_SPLIT,
        incomingPlaces: [0],
        outgoingPlaces: [1],
      },
      {
        id: 'Gateway_1o9s8fw_Flow_1n5x723',
        type: TransitionType.XOR_SPLIT,
        incomingPlaces: [0],
        outgoingPlaces: [2],
      },
      {
        id: 'Da1',
        name: 'Decline application (Da1)',
        type: TransitionType.TASK,
        incomingPlaces: [1],
        outgoingPlaces: [21],
      },
      {
        id: 'Aa',
        name: 'Accept application (Aa)',
        type: TransitionType.TASK,
        incomingPlaces: [2],
        outgoingPlaces: [3],
      },
      {
        id: 'Gateway_1rgq5gy',
        type: TransitionType.AND_SPLIT,
        incomingPlaces: [3],
        outgoingPlaces: [4, 5],
      },
      {
        id: 'Fa',
        name: 'Finalize application (Fa)',
        type: TransitionType.TASK,
        incomingPlaces: [4],
        outgoingPlaces: [10],
      },
      {
        id: 'Gateway_1fwxfgu_Flow_0vs3ms0',
        type: TransitionType.XOR_JOIN,
        incomingPlaces: [5],
        outgoingPlaces: [6],
      },
      {
        id: 'Gateway_1fwxfgu_Flow_1dlqvig',
        type: TransitionType.XOR_JOIN,
        incomingPlaces: [12],
        outgoingPlaces: [6],
      },
      {
        id: 'Sso',
        name: 'Select and send offer (Sso)',
        type: TransitionType.TASK,
        incomingPlaces: [6],
        outgoingPlaces: [7],
      },
      {
        id: 'Ro',
        name: 'Receive offer (Ro)',
        type: TransitionType.TASK,
        incomingPlaces: [7],
        outgoingPlaces: [8],
      },
      {
        id: 'Gateway_1way65i_Flow_0q0zlpw',
        type: TransitionType.XOR_SPLIT,
        incomingPlaces: [8],
        outgoingPlaces: [9],
      },
      {
        id: 'Gateway_1way65i_Flow_0x7rbwv',
        type: TransitionType.XOR_SPLIT,
        incomingPlaces: [8],
        outgoingPlaces: [11],
      },
      {
        id: 'Gateway_1nglicj',
        type: TransitionType.AND_JOIN,
        incomingPlaces: [9, 10],
        outgoingPlaces: [13],
      },
      {
        id: 'Co',
        name: 'Cancel offer (Co)',
        type: TransitionType.TASK,
        incomingPlaces: [11],
        outgoingPlaces: [12],
      },
      {
        id: 'Gateway_1vl4hvy_Flow_0oshibp',
        type: TransitionType.XOR_SPLIT,
        incomingPlaces: [13],
        outgoingPlaces: [14],
      },
      {
        id: 'Gateway_1vl4hvy_Flow_1b9hzto',
        type: TransitionType.XOR_SPLIT,
        incomingPlaces: [13],
        outgoingPlaces: [15],
      },
      {
        id: 'Ao',
        name: 'Accept offer (Ao)',
        type: TransitionType.TASK,
        incomingPlaces: [14],
        outgoingPlaces: [16],
      },
      {
        id: 'Do',
        name: 'Decline offer (Do)',
        type: TransitionType.TASK,
        incomingPlaces: [15],
        outgoingPlaces: [17],
      },
      {
        id: 'Aaa',
        name: 'Approve and activate application (Aaa)',
        type: TransitionType.TASK,
        incomingPlaces: [16],
        outgoingPlaces: [18],
      },
      {
        id: 'Da2',
        name: 'Decline application (Da2)',
        type: TransitionType.TASK,
        incomingPlaces: [17],
        outgoingPlaces: [19],
      },
      {
        id: 'Gateway_1vzsa13_Flow_03xp19s',
        type: TransitionType.XOR_JOIN,
        incomingPlaces: [18],
        outgoingPlaces: [20],
      },
      {
        id: 'Gateway_1vzsa13_Flow_12yxkzr',
        type: TransitionType.XOR_JOIN,
        incomingPlaces: [19],
        outgoingPlaces: [20],
      },
      {
        id: 'Gateway_1vzsa13_Flow_0ka5y4w',
        type: TransitionType.XOR_JOIN,
        incomingPlaces: [21],
        outgoingPlaces: [20],
      },
      {
        id: 'Af',
        name: 'Application finished (Af)',
        type: TransitionType.END,
        incomingPlaces: [20],
        outgoingPlaces: [],
      },
    ];
    return {
      id: 'conformance_example',
      placeCount: 22,
      transitions,
    };
  }

  static getModel2(): Model {
    const transitions: Transition[] = [
      {
        id: 'As',
        name: 'Application submitted (As)',
        type: TransitionType.START,
        incomingPlaces: [],
        outgoingPlaces: [0],
      },
      {
        id: 'Da1',
        name: 'Decline application (Da1)',
        type: TransitionType.TASK,
        incomingPlaces: [0],
        outgoingPlaces: [8],
      },
      {
        id: 'Aa',
        name: 'Accept application (Aa)',
        type: TransitionType.TASK,
        incomingPlaces: [0],
        outgoingPlaces: [1, 2],
      },
      {
        id: 'Fa',
        name: 'Finalize application (Fa)',
        type: TransitionType.TASK,
        incomingPlaces: [1],
        outgoingPlaces: [5],
      },
      {
        id: 'Sso',
        name: 'Select and send offer (Sso)',
        type: TransitionType.TASK,
        incomingPlaces: [2],
        outgoingPlaces: [3],
      },
      {
        id: 'Ro',
        name: 'Receive offer (Ro)',
        type: TransitionType.TASK,
        incomingPlaces: [3],
        outgoingPlaces: [4],
      },
      {
        id: 'Co',
        name: 'Cancel offer (Co)',
        type: TransitionType.TASK,
        incomingPlaces: [4],
        outgoingPlaces: [2],
      },
      {
        id: 'Ao',
        name: 'Accept offer (Ao)',
        type: TransitionType.TASK,
        incomingPlaces: [5, 4],
        outgoingPlaces: [6],
      },
      {
        id: 'Do',
        name: 'Decline offer (Do)',
        type: TransitionType.TASK,
        incomingPlaces: [5, 4],
        outgoingPlaces: [7],
      },
      {
        id: 'Aaa',
        name: 'Approve and activate application (Aaa)',
        type: TransitionType.TASK,
        incomingPlaces: [6],
        outgoingPlaces: [8],
      },
      {
        id: 'Da2',
        name: 'Decline application (Da2)',
        type: TransitionType.TASK,
        incomingPlaces: [7],
        outgoingPlaces: [8],
      },
      {
        id: 'Af',
        name: 'Application finished (Af)',
        type: TransitionType.END,
        incomingPlaces: [8],
        outgoingPlaces: [],
      },
    ];
    return {
      id: 'conformance_example',
      placeCount: 9,
      transitions,
    };
  }

  static getInstance1(): Instance {
    return {
      id: 'instance1',
      model: this.getModel2(),
      executionStatuses: [
        ExecutionStatus.NOT_ACTIVE,
        ExecutionStatus.NOT_ACTIVE,
        ExecutionStatus.NOT_ACTIVE,
        ExecutionStatus.NOT_ACTIVE,
        ExecutionStatus.NOT_ACTIVE,
        ExecutionStatus.NOT_ACTIVE,
        ExecutionStatus.NOT_ACTIVE,
        ExecutionStatus.NOT_ACTIVE,
        ExecutionStatus.NOT_ACTIVE,
      ],
      finished: false,
    };
  }
}
