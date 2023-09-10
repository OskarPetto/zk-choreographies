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
    };
  }

  static getModel1(): Model {
    const transitions: Transition[] = [
      {
        id: 'As',
        name: 'Application submitted (As)',
        type: TransitionType.START,
        incomingPlaces: [],
        outgoingPlaces: ['Flow_0xqg6as'],
      },
      {
        id: 'Gateway_1o9s8fw_Flow_07cyta2',
        type: TransitionType.XOR_SPLIT,
        incomingPlaces: ['Flow_0xqg6as'],
        outgoingPlaces: ['Flow_07cyta2'],
      },
      {
        id: 'Gateway_1o9s8fw_Flow_1n5x723',
        type: TransitionType.XOR_SPLIT,
        incomingPlaces: ['Flow_0xqg6as'],
        outgoingPlaces: ['Flow_1n5x723'],
      },
      {
        id: 'Da1',
        name: 'Decline application (Da1)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_07cyta2'],
        outgoingPlaces: ['Flow_0ka5y4w'],
      },
      {
        id: 'Aa',
        name: 'Accept application (Aa)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_1n5x723'],
        outgoingPlaces: ['Flow_0w362at'],
      },
      {
        id: 'Gateway_1rgq5gy',
        type: TransitionType.AND_SPLIT,
        incomingPlaces: ['Flow_0w362at'],
        outgoingPlaces: ['Flow_0i668iz', 'Flow_0vs3ms0'],
      },
      {
        id: 'Fa',
        name: 'Finalize application (Fa)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_0i668iz'],
        outgoingPlaces: ['Flow_1yee6jg'],
      },
      {
        id: 'Gateway_1fwxfgu_Flow_0vs3ms0',
        type: TransitionType.XOR_JOIN,
        incomingPlaces: ['Flow_0vs3ms0'],
        outgoingPlaces: ['Flow_1g76mn1'],
      },
      {
        id: 'Gateway_1fwxfgu_Flow_1dlqvig',
        type: TransitionType.XOR_JOIN,
        incomingPlaces: ['Flow_1dlqvig'],
        outgoingPlaces: ['Flow_1g76mn1'],
      },
      {
        id: 'Sso',
        name: 'Select and send offer (Sso)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_1g76mn1'],
        outgoingPlaces: ['Flow_19w6dwb'],
      },
      {
        id: 'Ro',
        name: 'Receive offer (Ro)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_19w6dwb'],
        outgoingPlaces: ['Flow_0damgea'],
      },
      {
        id: 'Gateway_1way65i_Flow_0q0zlpw',
        type: TransitionType.XOR_SPLIT,
        incomingPlaces: ['Flow_0damgea'],
        outgoingPlaces: ['Flow_0q0zlpw'],
      },
      {
        id: 'Gateway_1way65i_Flow_0x7rbwv',
        type: TransitionType.XOR_SPLIT,
        incomingPlaces: ['Flow_0damgea'],
        outgoingPlaces: ['Flow_0x7rbwv'],
      },
      {
        id: 'Gateway_1nglicj',
        type: TransitionType.AND_JOIN,
        incomingPlaces: ['Flow_0q0zlpw', 'Flow_1yee6jg'],
        outgoingPlaces: ['Flow_0jf7vbw'],
      },
      {
        id: 'Co',
        name: 'Cancel offer (Co)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_0x7rbwv'],
        outgoingPlaces: ['Flow_1dlqvig'],
      },
      {
        id: 'Gateway_1vl4hvy_Flow_0oshibp',
        type: TransitionType.XOR_SPLIT,
        incomingPlaces: ['Flow_0jf7vbw'],
        outgoingPlaces: ['Flow_0oshibp'],
      },
      {
        id: 'Gateway_1vl4hvy_Flow_1b9hzto',
        type: TransitionType.XOR_SPLIT,
        incomingPlaces: ['Flow_0jf7vbw'],
        outgoingPlaces: ['Flow_1b9hzto'],
      },
      {
        id: 'Ao',
        name: 'Accept offer (Ao)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_0oshibp'],
        outgoingPlaces: ['Flow_0u5u7d3'],
      },
      {
        id: 'Do',
        name: 'Decline offer (Do)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_1b9hzto'],
        outgoingPlaces: ['Flow_18msx9i'],
      },
      {
        id: 'Aaa',
        name: 'Approve and activate application (Aaa)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_0u5u7d3'],
        outgoingPlaces: ['Flow_03xp19s'],
      },
      {
        id: 'Da2',
        name: 'Decline application (Da2)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_18msx9i'],
        outgoingPlaces: ['Flow_12yxkzr'],
      },
      {
        id: 'Gateway_1vzsa13_Flow_03xp19s',
        type: TransitionType.XOR_JOIN,
        incomingPlaces: ['Flow_03xp19s'],
        outgoingPlaces: ['Flow_1qh0ya1'],
      },
      {
        id: 'Gateway_1vzsa13_Flow_12yxkzr',
        type: TransitionType.XOR_JOIN,
        incomingPlaces: ['Flow_12yxkzr'],
        outgoingPlaces: ['Flow_1qh0ya1'],
      },
      {
        id: 'Gateway_1vzsa13_Flow_0ka5y4w',
        type: TransitionType.XOR_JOIN,
        incomingPlaces: ['Flow_0ka5y4w'],
        outgoingPlaces: ['Flow_1qh0ya1'],
      },
      {
        id: 'Af',
        name: 'Application finished (Af)',
        type: TransitionType.END,
        incomingPlaces: ['Flow_1qh0ya1'],
        outgoingPlaces: [],
      },
    ];
    return {
      id: 'conformance_example',
      transitions: new Map(transitions.map((t) => [t.id, t])),
    };
  }

  static getModel2(): Model {
    const transitions: Transition[] = [
      {
        id: 'As',
        name: 'Application submitted (As)',
        type: TransitionType.START,
        incomingPlaces: [],
        outgoingPlaces: ['Flow_0xqg6as'],
      },
      {
        id: 'Da1',
        name: 'Decline application (Da1)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_0xqg6as'],
        outgoingPlaces: ['Flow_1qh0ya1'],
      },
      {
        id: 'Aa',
        name: 'Accept application (Aa)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_0xqg6as'],
        outgoingPlaces: ['Flow_0i668iz', 'Flow_1g76mn1'],
      },
      {
        id: 'Fa',
        name: 'Finalize application (Fa)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_0i668iz'],
        outgoingPlaces: ['Flow_1yee6jg'],
      },
      {
        id: 'Sso',
        name: 'Select and send offer (Sso)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_1g76mn1'],
        outgoingPlaces: ['Flow_19w6dwb'],
      },
      {
        id: 'Ro',
        name: 'Receive offer (Ro)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_19w6dwb'],
        outgoingPlaces: ['Flow_0damgea'],
      },
      {
        id: 'Co',
        name: 'Cancel offer (Co)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_0damgea'],
        outgoingPlaces: ['Flow_1g76mn1'],
      },
      {
        id: 'Ao',
        name: 'Accept offer (Ao)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_1yee6jg', 'Flow_0damgea'],
        outgoingPlaces: ['Flow_0u5u7d3'],
      },
      {
        id: 'Do',
        name: 'Decline offer (Do)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_1yee6jg', 'Flow_0damgea'],
        outgoingPlaces: ['Flow_18msx9i'],
      },
      {
        id: 'Aaa',
        name: 'Approve and activate application (Aaa)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_0u5u7d3'],
        outgoingPlaces: ['Flow_1qh0ya1'],
      },
      {
        id: 'Da2',
        name: 'Decline application (Da2)',
        type: TransitionType.TASK,
        incomingPlaces: ['Flow_18msx9i'],
        outgoingPlaces: ['Flow_1qh0ya1'],
      },
      {
        id: 'Af',
        name: 'Application finished (Af)',
        type: TransitionType.END,
        incomingPlaces: ['Flow_1qh0ya1'],
        outgoingPlaces: [],
      },
    ];
    return {
      id: 'conformance_example',
      transitions: new Map(transitions.map((t) => [t.id, t])),
    };
  }

  static getInstance1(): Instance {
    return {
      id: 'instance1',
      model: this.getModel2().id,
      executionStatuses: new Map([
        ['Flow_0xqg6as', ExecutionStatus.NOT_ACTIVE],
        ['Flow_1qh0ya1', ExecutionStatus.NOT_ACTIVE],
        ['Flow_0i668iz', ExecutionStatus.NOT_ACTIVE],
        ['Flow_1g76mn1', ExecutionStatus.NOT_ACTIVE],
        ['Flow_1yee6jg', ExecutionStatus.NOT_ACTIVE],
        ['Flow_19w6dwb', ExecutionStatus.NOT_ACTIVE],
        ['Flow_0damgea', ExecutionStatus.NOT_ACTIVE],
        ['Flow_0u5u7d3', ExecutionStatus.NOT_ACTIVE],
        ['Flow_18msx9i', ExecutionStatus.NOT_ACTIVE],
      ]),
      finished: false,
    };
  }
}
