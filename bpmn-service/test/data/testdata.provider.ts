import * as fs from 'fs';
import * as path from 'path';
import { Model, TransitionType, defaultHash } from '../../src/model/model';
import {
  Definitions,
  GatewayType,
  LoopType,
} from 'src/choreography/choreography';

function readTextFile(filename: string) {
  const filePath = path.join(process.cwd(), filename);
  return fs.readFileSync(filePath, 'utf-8').toString();
}

function writeTextFile(filename: string, data: string) {
  const filePath = path.join(process.cwd(), filename);
  return fs.writeFileSync(filePath, data);
}

export class TestdataProvider {
  static getFloorChoreography(): string {
    return readTextFile('test/data/floor_choreography.bpmn');
  }

  static getExampleChoreography(): string {
    return readTextFile('test/data/example_choreography.bpmn');
  }

  static writeExampleChoreography() {
    writeTextFile(
      'test/data/example_choreography.json',
      JSON.stringify(this.getModel2Reduced()),
    );
  }

  static getDefinitions2(): Definitions {
    return {
      choreographies: [
        {
          id: 'Choreography_07n6r3q',
          participants: [
            {
              id: 'Participant_0x6v44d',
              name: 'Customer',
            },
            {
              id: 'Participant_0n7kwiu',
              name: 'Seller',
            },
            {
              id: 'Participant_0xxs9a7',
              name: 'Supplier',
            },
          ],
          startEvents: [
            {
              id: 'Event_1525yky',
              outgoing: 'Flow_1f6eaf2',
            },
          ],
          endEvents: [
            {
              id: 'Event_08d32d7',
              name: 'Order fullfilled',
              incoming: 'Flow_0vbgnyw',
            },
          ],
          exclusiveGateways: [
            {
              id: 'Gateway_10fv7g5',
              type: GatewayType.SPLIT,
              incoming: ['Flow_0kyymfy'],
              outgoing: ['Flow_1k5cri3', 'Flow_1stf9mf'],
            },
            {
              id: 'Gateway_1ofchbb',
              type: GatewayType.JOIN,
              incoming: ['Flow_1iwbmcz', 'Flow_0ec3gh0'],
              outgoing: ['Flow_1snas91'],
            },
          ],
          parallelGateways: [
            {
              id: 'Gateway_0g1gxjc',
              type: GatewayType.SPLIT,
              incoming: ['Flow_1snas91'],
              outgoing: ['Flow_1jplxjf', 'Flow_1h69rwb'],
            },
            {
              id: 'Gateway_0nm7ind',
              type: GatewayType.JOIN,
              incoming: ['Flow_0glm8k3', 'Flow_0433489'],
              outgoing: ['Flow_0vbgnyw'],
            },
          ],
          choreographyTasks: [
            {
              id: 'ChoreographyTask_0kp4flv',
              name: 'Submit purchase order',
              incoming: 'Flow_1f6eaf2',
              outgoing: 'Flow_0kyymfy',
              initiatingParticipant: 'Participant_0x6v44d',
              respondingParticipant: 'Participant_0n7kwiu',
              initialMessage: 'Message_1376fyb',
              responseMessage: 'Message_0zzvac5',
            },
            {
              id: 'ChoreographyTask_0nl2rhr',
              name: 'Purchase raw materials',
              incoming: 'Flow_1k5cri3',
              outgoing: 'Flow_0c5yqsz',
              initiatingParticipant: 'Participant_0n7kwiu',
              respondingParticipant: 'Participant_0xxs9a7',
              initialMessage: 'Message_0xe03aa',
              responseMessage: 'Message_1p1ke3y',
              loopType: LoopType.MULTI_INSTANCE_SEQUENTIAL,
            },
            {
              id: 'ChoreographyTask_1uie9z3',
              name: 'Confirm order',
              incoming: 'Flow_0c5yqsz',
              outgoing: 'Flow_1iwbmcz',
              initiatingParticipant: 'Participant_0n7kwiu',
              respondingParticipant: 'Participant_0x6v44d',
              initialMessage: 'Message_1dd2uoz',
            },
            {
              id: 'ChoreographyTask_1dsovf5',
              name: 'Ship product',
              incoming: 'Flow_1jplxjf',
              outgoing: 'Flow_0glm8k3',
              initiatingParticipant: 'Participant_0x6v44d',
              respondingParticipant: 'Participant_0n7kwiu',
              initialMessage: 'Message_1858jqq',
              responseMessage: 'Message_01cp9ki',
            },
            {
              id: 'ChoreographyTask_1htg6wy',
              name: 'Invoice customer',
              incoming: 'Flow_1h69rwb',
              outgoing: 'Flow_0433489',
              initiatingParticipant: 'Participant_0n7kwiu',
              respondingParticipant: 'Participant_0x6v44d',
              initialMessage: 'Message_0annsni',
              responseMessage: 'Message_05ebl37',
            },
            {
              id: 'ChoreographyTask_1e51o0k',
              name: 'Confirm order',
              incoming: 'Flow_1stf9mf',
              outgoing: 'Flow_0ec3gh0',
              initiatingParticipant: 'Participant_0n7kwiu',
              respondingParticipant: 'Participant_0x6v44d',
              initialMessage: 'Message_0gdct3r',
            },
          ],
          sequenceFlows: [
            {
              id: 'Flow_1f6eaf2',
            },
            {
              id: 'Flow_0kyymfy',
            },
            {
              id: 'Flow_1k5cri3',
              name: 'order > stock',
            },
            {
              id: 'Flow_1iwbmcz',
            },
            {
              id: 'Flow_1jplxjf',
            },
            {
              id: 'Flow_1h69rwb',
            },
            {
              id: 'Flow_0glm8k3',
            },
            {
              id: 'Flow_0433489',
            },
            {
              id: 'Flow_0vbgnyw',
            },
            {
              id: 'Flow_1stf9mf',
              name: 'order <= stock',
            },
            {
              id: 'Flow_0c5yqsz',
            },
            {
              id: 'Flow_1snas91',
            },
            {
              id: 'Flow_0ec3gh0',
            },
          ],
          messages: [
            {
              id: 'Message_0zzvac5',
              name: 'stock',
            },
            {
              id: 'Message_0gdct3r',
              name: 'confirm',
            },
            {
              id: 'Message_05ebl37',
              name: 'payment',
            },
            {
              id: 'Message_01cp9ki',
              name: 'product',
            },
            {
              id: 'Message_1p1ke3y',
              name: 'raw_materials',
            },
            {
              id: 'Message_0annsni',
              name: 'invoice',
            },
            {
              id: 'Message_1858jqq',
              name: 'shipping_address',
            },
            {
              id: 'Message_1dd2uoz',
              name: 'confirm',
            },
            {
              id: 'Message_0xe03aa',
              name: 'rm_order',
            },
            {
              id: 'Message_1376fyb',
              name: 'order',
            },
          ],
        },
      ],
    };
  }

  static getModel2(): Model {
    return {
      hash: defaultHash(),
      source: TestdataProvider.getExampleChoreography(),
      placeCount: 20,
      participantCount: 2,
      messageCount: 10,
      startPlaces: [18],
      endPlaces: [19],
      transitions: [
        {
          id: 'Event_1525yky',
          type: TransitionType.REQUIRED,
          incomingPlaces: [18],
          outgoingPlaces: [0],
        },
        {
          id: 'Event_08d32d7',
          type: TransitionType.REQUIRED,
          name: 'Order fullfilled',
          incomingPlaces: [8],
          outgoingPlaces: [19],
        },
        {
          id: 'Gateway_10fv7g5_Flow_1k5cri3',
          type: TransitionType.OPTIONAL_OUTGOING,
          incomingPlaces: [1],
          outgoingPlaces: [2],
        },
        {
          id: 'Gateway_10fv7g5_Flow_1stf9mf',
          type: TransitionType.OPTIONAL_OUTGOING,
          incomingPlaces: [1],
          outgoingPlaces: [9],
        },
        {
          id: 'Gateway_1ofchbb_Flow_1iwbmcz',
          type: TransitionType.OPTIONAL_INCOMING,
          incomingPlaces: [3],
          outgoingPlaces: [11],
        },
        {
          id: 'Gateway_1ofchbb_Flow_0ec3gh0',
          type: TransitionType.OPTIONAL_INCOMING,
          incomingPlaces: [12],
          outgoingPlaces: [11],
        },
        {
          id: 'Gateway_0g1gxjc',
          type: TransitionType.OPTIONAL_INCOMING,
          incomingPlaces: [11],
          outgoingPlaces: [4, 5],
        },
        {
          id: 'Gateway_0nm7ind',
          type: TransitionType.OPTIONAL_OUTGOING,
          incomingPlaces: [6, 7],
          outgoingPlaces: [8],
        },
        {
          id: 'ChoreographyTask_0kp4flv_Participant_0x6v44d',
          type: TransitionType.REQUIRED,
          name: 'Submit purchase order',
          incomingPlaces: [0],
          outgoingPlaces: [13],
          initiatingParticipant: 0,
          respondingParticipant: 1,
          message: 9,
        },
        {
          id: 'ChoreographyTask_0kp4flv_Participant_0n7kwiu',
          type: TransitionType.REQUIRED,
          name: 'Submit purchase order',
          incomingPlaces: [13],
          outgoingPlaces: [1],
          initiatingParticipant: 1,
          respondingParticipant: 0,
          message: 0,
        },
        {
          id: 'ChoreographyTask_0nl2rhr_0',
          type: TransitionType.REQUIRED,
          name: 'Purchase raw materials',
          incomingPlaces: [2],
          outgoingPlaces: [10],
          initiatingParticipant: 1,
          respondingParticipant: 2,
          message: 8,
          messageConstraint: {
            coefficients: [1, -1],
            messageIds: [9, 0],
            offset: 0,
            comparisonOperator: 1,
          },
        },
        {
          id: 'ChoreographyTask_0nl2rhr_loop',
          type: TransitionType.REQUIRED,
          name: 'Purchase raw materials',
          incomingPlaces: [10],
          outgoingPlaces: [10],
          initiatingParticipant: 1,
          respondingParticipant: 2,
          message: 8,
        },
        {
          id: 'ChoreographyTask_1uie9z3_Participant_0n7kwiu',
          type: TransitionType.REQUIRED,
          name: 'Confirm order',
          incomingPlaces: [10],
          outgoingPlaces: [14],
          initiatingParticipant: 1,
          respondingParticipant: 0,
          message: 7,
        },
        {
          id: 'ChoreographyTask_1uie9z3_Participant_0x6v44d',
          type: TransitionType.REQUIRED,
          name: 'Confirm order',
          incomingPlaces: [14],
          outgoingPlaces: [3],
          initiatingParticipant: 0,
          respondingParticipant: 1,
        },
        {
          id: 'ChoreographyTask_1dsovf5_Participant_0x6v44d',
          type: TransitionType.REQUIRED,
          name: 'Ship product',
          incomingPlaces: [4],
          outgoingPlaces: [15],
          initiatingParticipant: 0,
          respondingParticipant: 1,
          message: 6,
        },
        {
          id: 'ChoreographyTask_1dsovf5_Participant_0n7kwiu',
          type: TransitionType.REQUIRED,
          name: 'Ship product',
          incomingPlaces: [15],
          outgoingPlaces: [6],
          initiatingParticipant: 1,
          respondingParticipant: 0,
          message: 3,
        },
        {
          id: 'ChoreographyTask_1htg6wy_Participant_0n7kwiu',
          type: TransitionType.REQUIRED,
          name: 'Invoice customer',
          incomingPlaces: [5],
          outgoingPlaces: [16],
          initiatingParticipant: 1,
          respondingParticipant: 0,
          message: 5,
        },
        {
          id: 'ChoreographyTask_1htg6wy_Participant_0x6v44d',
          type: TransitionType.REQUIRED,
          name: 'Invoice customer',
          incomingPlaces: [16],
          outgoingPlaces: [7],
          initiatingParticipant: 0,
          respondingParticipant: 1,
          message: 2,
        },
        {
          id: 'ChoreographyTask_1e51o0k_Participant_0n7kwiu',
          type: TransitionType.REQUIRED,
          name: 'Confirm order',
          incomingPlaces: [9],
          outgoingPlaces: [17],
          initiatingParticipant: 1,
          respondingParticipant: 0,
          message: 1,
          messageConstraint: {
            coefficients: [1, -1],
            messageIds: [9, 0],
            offset: 0,
            comparisonOperator: 4,
          },
        },
        {
          id: 'ChoreographyTask_1e51o0k_Participant_0x6v44d',
          type: TransitionType.REQUIRED,
          name: 'Confirm order',
          incomingPlaces: [17],
          outgoingPlaces: [12],
          initiatingParticipant: 0,
          respondingParticipant: 1,
        },
      ],
    };
  }

  static getModel2Reduced(): Model {
    return {
      hash: defaultHash(),
      source: TestdataProvider.getExampleChoreography(),
      placeCount: 14,
      participantCount: 2,
      messageCount: 10,
      startPlaces: [12],
      endPlaces: [13],
      transitions: [
        {
          id: 'Event_1525yky',
          type: TransitionType.REQUIRED,
          incomingPlaces: [12],
          outgoingPlaces: [0],
        },
        {
          id: 'Event_08d32d7',
          type: TransitionType.REQUIRED,
          name: 'Order fullfilled',
          incomingPlaces: [4, 5],
          outgoingPlaces: [13],
        },
        {
          id: 'ChoreographyTask_0kp4flv_Participant_0x6v44d',
          type: TransitionType.REQUIRED,
          name: 'Submit purchase order',
          incomingPlaces: [0],
          outgoingPlaces: [7],
          message: 9,
          initiatingParticipant: 0,
          respondingParticipant: 1,
        },
        {
          id: 'ChoreographyTask_0kp4flv_Participant_0n7kwiu',
          type: TransitionType.REQUIRED,
          name: 'Submit purchase order',
          incomingPlaces: [7],
          outgoingPlaces: [1],
          message: 0,
          initiatingParticipant: 1,
          respondingParticipant: 0,
        },
        {
          id: 'ChoreographyTask_0nl2rhr_0',
          type: TransitionType.REQUIRED,
          name: 'Purchase raw materials',
          incomingPlaces: [1],
          outgoingPlaces: [6],
          message: 8,
          initiatingParticipant: 1,
          respondingParticipant: 2,
          messageConstraint: {
            coefficients: [1, -1],
            messageIds: [9, 0],
            offset: 0,
            comparisonOperator: 1,
          },
        },
        {
          id: 'ChoreographyTask_0nl2rhr_loop',
          type: TransitionType.REQUIRED,
          name: 'Purchase raw materials',
          incomingPlaces: [6],
          outgoingPlaces: [6],
          message: 8,
          initiatingParticipant: 1,
          respondingParticipant: 2,
        },
        {
          id: 'ChoreographyTask_1uie9z3_Participant_0n7kwiu',
          type: TransitionType.REQUIRED,
          name: 'Confirm order',
          incomingPlaces: [6],
          outgoingPlaces: [8],
          message: 7,
          initiatingParticipant: 1,
          respondingParticipant: 0,
        },
        {
          id: 'ChoreographyTask_1uie9z3_Participant_0x6v44d',
          type: TransitionType.REQUIRED,
          name: 'Confirm order',
          incomingPlaces: [8],
          outgoingPlaces: [2, 3],
          initiatingParticipant: 0,
          respondingParticipant: 1,
        },
        {
          id: 'ChoreographyTask_1dsovf5_Participant_0x6v44d',
          type: TransitionType.REQUIRED,
          name: 'Ship product',
          incomingPlaces: [2],
          outgoingPlaces: [9],
          message: 6,
          initiatingParticipant: 0,
          respondingParticipant: 1,
        },
        {
          id: 'ChoreographyTask_1dsovf5_Participant_0n7kwiu',
          type: TransitionType.REQUIRED,
          name: 'Ship product',
          incomingPlaces: [9],
          outgoingPlaces: [4],
          message: 3,
          initiatingParticipant: 1,
          respondingParticipant: 0,
        },
        {
          id: 'ChoreographyTask_1htg6wy_Participant_0n7kwiu',
          type: TransitionType.REQUIRED,
          name: 'Invoice customer',
          incomingPlaces: [3],
          outgoingPlaces: [10],
          message: 5,
          initiatingParticipant: 1,
          respondingParticipant: 0,
        },
        {
          id: 'ChoreographyTask_1htg6wy_Participant_0x6v44d',
          type: TransitionType.REQUIRED,
          name: 'Invoice customer',
          incomingPlaces: [10],
          outgoingPlaces: [5],
          message: 2,
          initiatingParticipant: 0,
          respondingParticipant: 1,
        },
        {
          id: 'ChoreographyTask_1e51o0k_Participant_0n7kwiu',
          type: TransitionType.REQUIRED,
          name: 'Confirm order',
          incomingPlaces: [1],
          outgoingPlaces: [11],
          message: 1,
          initiatingParticipant: 1,
          respondingParticipant: 0,
          messageConstraint: {
            coefficients: [1, -1],
            messageIds: [9, 0],
            offset: 0,
            comparisonOperator: 4,
          },
        },
        {
          id: 'ChoreographyTask_1e51o0k_Participant_0x6v44d',
          type: TransitionType.REQUIRED,
          name: 'Confirm order',
          incomingPlaces: [11],
          outgoingPlaces: [2, 3],
          initiatingParticipant: 0,
          respondingParticipant: 1,
        },
      ],
    };
  }

  static getDate(): Date {
    return new Date(Date.parse('2023-09-27T22:57:44.261Z'));
  }
}
