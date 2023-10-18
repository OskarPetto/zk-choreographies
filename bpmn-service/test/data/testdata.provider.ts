import * as fs from 'fs';
import * as path from 'path';
import { Model, TransitionType } from '../../src/model/model';
import { Definitions, GatewayType } from 'src/choreography/choreography';

function readTextFile(filename: string) {
  const filePath = path.join(process.cwd(), filename);
  return fs.readFileSync(filePath, 'utf-8').toString();
}

function writeTextFile(filename: string, data: string) {
  const filePath = path.join(process.cwd(), filename);
  return fs.writeFileSync(filePath, data);
}

export class TestdataProvider {
  static readBpmn(filename: string): string {
    return readTextFile(`test/data/${filename}.bpmn`);
  }

  static writeModel(filename: string, model: Model) {
    writeTextFile(`test/data/${filename}.json`, JSON.stringify(model));
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
      source: 'test',
      placeCount: 19,
      participantCount: 3,
      messageCount: 10,
      startPlaces: [17],
      endPlaces: [18],
      transitions: [
        {
          id: 'Event_1525yky',
          type: TransitionType.REQUIRED,
          incomingPlaces: [17],
          outgoingPlaces: [0],
        },
        {
          id: 'Event_08d32d7',
          type: TransitionType.REQUIRED,
          name: 'Order fullfilled',
          incomingPlaces: [8],
          outgoingPlaces: [18],
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
          id: 'ChoreographyTask_0kp4flv_0',
          type: TransitionType.REQUIRED,
          name: 'Submit purchase order',
          incomingPlaces: [0],
          outgoingPlaces: [13],
          sender: 0,
          recipient: 1,
          message: 9,
        },
        {
          id: 'ChoreographyTask_0kp4flv_1',
          type: TransitionType.REQUIRED,
          name: 'Submit purchase order',
          incomingPlaces: [13],
          outgoingPlaces: [1],
          sender: 1,
          recipient: 0,
          message: 0,
        },
        {
          id: 'ChoreographyTask_0nl2rhr_0',
          type: TransitionType.REQUIRED,
          name: 'Purchase raw materials',
          incomingPlaces: [2],
          outgoingPlaces: [14],
          sender: 1,
          recipient: 2,
          message: 8,
          constraint: {
            coefficients: [1, -1],
            messageIds: [9, 0],
            offset: 0,
            comparisonOperator: 1,
          },
        },
        {
          id: 'ChoreographyTask_0nl2rhr_1',
          type: TransitionType.REQUIRED,
          name: 'Purchase raw materials',
          incomingPlaces: [14],
          outgoingPlaces: [10],
          sender: 2,
          recipient: 1,
          message: 4,
        },
        {
          id: 'ChoreographyTask_1uie9z3',
          type: TransitionType.REQUIRED,
          name: 'Confirm order',
          incomingPlaces: [10],
          outgoingPlaces: [3],
          sender: 1,
          recipient: 0,
          message: 7,
        },
        {
          id: 'ChoreographyTask_1dsovf5_0',
          type: TransitionType.REQUIRED,
          name: 'Ship product',
          incomingPlaces: [4],
          outgoingPlaces: [15],
          sender: 0,
          recipient: 1,
          message: 6,
        },
        {
          id: 'ChoreographyTask_1dsovf5_1',
          type: TransitionType.REQUIRED,
          name: 'Ship product',
          incomingPlaces: [15],
          outgoingPlaces: [6],
          sender: 1,
          recipient: 0,
          message: 3,
        },
        {
          id: 'ChoreographyTask_1htg6wy_0',
          type: TransitionType.REQUIRED,
          name: 'Invoice customer',
          incomingPlaces: [5],
          outgoingPlaces: [16],
          sender: 1,
          recipient: 0,
          message: 5,
        },
        {
          id: 'ChoreographyTask_1htg6wy_1',
          type: TransitionType.REQUIRED,
          name: 'Invoice customer',
          incomingPlaces: [16],
          outgoingPlaces: [7],
          sender: 0,
          recipient: 1,
          message: 2,
        },
        {
          id: 'ChoreographyTask_1e51o0k',
          type: TransitionType.REQUIRED,
          name: 'Confirm order',
          incomingPlaces: [9],
          outgoingPlaces: [12],
          sender: 1,
          recipient: 0,
          message: 1,
          constraint: {
            coefficients: [1, -1],
            messageIds: [9, 0],
            offset: 0,
            comparisonOperator: 4,
          },
        },
      ],
    };
  }

  static getModel2Reduced(): Model {
    return {
      source: 'test',
      placeCount: 13,
      participantCount: 3,
      messageCount: 10,
      startPlaces: [11],
      endPlaces: [12],
      transitions: [
        {
          id: 'Event_1525yky',
          incomingPlaces: [11],
          outgoingPlaces: [0],
        },
        {
          id: 'Event_08d32d7',
          name: 'Order fullfilled',
          incomingPlaces: [4, 5],
          outgoingPlaces: [12],
        },
        {
          id: 'ChoreographyTask_0kp4flv_0',
          name: 'Submit purchase order',
          incomingPlaces: [0],
          outgoingPlaces: [7],
          message: 9,
          sender: 0,
          recipient: 1,
        },
        {
          id: 'ChoreographyTask_0kp4flv_1',
          name: 'Submit purchase order',
          incomingPlaces: [7],
          outgoingPlaces: [1],
          message: 0,
          sender: 1,
          recipient: 0,
        },
        {
          id: 'ChoreographyTask_0nl2rhr_0',
          name: 'Purchase raw materials',
          incomingPlaces: [1],
          outgoingPlaces: [8],
          message: 8,
          sender: 1,
          recipient: 2,
          constraint: {
            coefficients: [1, -1],
            messageIds: [9, 0],
            offset: 0,
            comparisonOperator: 1,
          },
        },
        {
          id: 'ChoreographyTask_0nl2rhr_1',
          name: 'Purchase raw materials',
          incomingPlaces: [8],
          outgoingPlaces: [6],
          message: 4,
          sender: 2,
          recipient: 1,
        },
        {
          id: 'ChoreographyTask_1uie9z3',
          name: 'Confirm order',
          incomingPlaces: [6],
          outgoingPlaces: [2, 3],
          message: 7,
          sender: 1,
          recipient: 0,
        },
        {
          id: 'ChoreographyTask_1dsovf5_0',
          name: 'Ship product',
          incomingPlaces: [2],
          outgoingPlaces: [9],
          message: 6,
          sender: 0,
          recipient: 1,
        },
        {
          id: 'ChoreographyTask_1dsovf5_1',
          name: 'Ship product',
          incomingPlaces: [9],
          outgoingPlaces: [4],
          message: 3,
          sender: 1,
          recipient: 0,
        },
        {
          id: 'ChoreographyTask_1htg6wy_0',
          name: 'Invoice customer',
          incomingPlaces: [3],
          outgoingPlaces: [10],
          message: 5,
          sender: 1,
          recipient: 0,
        },
        {
          id: 'ChoreographyTask_1htg6wy_1',
          name: 'Invoice customer',
          incomingPlaces: [10],
          outgoingPlaces: [5],
          message: 2,
          sender: 0,
          recipient: 1,
        },
        {
          id: 'ChoreographyTask_1e51o0k',
          name: 'Confirm order',
          incomingPlaces: [1],
          outgoingPlaces: [2, 3],
          message: 1,
          sender: 1,
          recipient: 0,
          constraint: {
            coefficients: [1, -1],
            messageIds: [9, 0],
            offset: 0,
            comparisonOperator: 4,
          },
        },
      ],
    };
  }

  static getDate(): Date {
    return new Date(Date.parse('2023-09-27T22:57:44.261Z'));
  }
}
