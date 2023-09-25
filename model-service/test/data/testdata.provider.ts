import * as fs from 'fs';
import * as path from 'path';
import {
  Model,
  Transition,
  TransitionType,
} from 'src/model/model';
import { Definitions, GatewayType } from 'src/bpmn/bpmn';

function readTextFile(filename: string) {
  const filePath = path.join(process.cwd(), filename);
  return fs.readFileSync(filePath, 'utf-8').toString();
}

export class TestdataProvider {
  static getFloorChoreography(): string {
    return readTextFile('test/data/floor_choreography.bpmn');
  }

  static getExampleChoreography(): string {
    return readTextFile('test/data/example_choreography.bpmn');
  }

  static getDefinitions2(): Definitions {
    return {
      choreographies: [
        {
          id: "Choreography_07n6r3q",
          participants: [
            {
              id: "Participant_0x6v44d",
              name: "Customer",
            },
            {
              id: "Participant_0n7kwiu",
              name: "Seller",
            },
            {
              id: "Participant_0xxs9a7",
              name: "Supplier",
              maxMultiplicity: 2
            }
          ],
          startEvent: {
            id: "Event_1525yky",
            outgoing: "Flow_1f6eaf2",
          },
          endEvents: [
            {
              id: "Event_08d32d7",
              name: "Order fullfilled",
              incoming: "Flow_0vbgnyw"
            }
          ],
          exclusiveGateways: [
            {
              id: "Gateway_10fv7g5",
              default: "Flow_0572p2k",
              type: GatewayType.SPLIT,
              incoming: [
                "Flow_0kyymfy"
              ],
              outgoing: [
                "Flow_1k5cri3",
                "Flow_0572p2k"
              ]
            },
            {
              id: "Gateway_035usvb",
              type: GatewayType.JOIN,
              incoming: [
                "Flow_0572p2k",
                "Flow_1oa3jbg"
              ],
              outgoing: [
                "Flow_0rtcg7d"
              ]
            }
          ],
          parallelGateways: [
            {
              id: "Gateway_0g1gxjc",
              type: GatewayType.SPLIT,
              incoming: [
                "Flow_1iwbmcz"
              ],
              outgoing: [
                "Flow_1jplxjf",
                "Flow_1h69rwb"
              ]
            },
            {
              id: "Gateway_0nm7ind",
              type: GatewayType.JOIN,
              incoming: [
                "Flow_0glm8k3",
                "Flow_0433489"
              ],
              outgoing: [
                "Flow_0vbgnyw"
              ]
            },
          ],
          choreographyTasks: [
            {
              id: "ChoreographyTask_0kp4flv",
              name: "Submit purchase order",
              incoming: "Flow_1f6eaf2",
              outgoing: "Flow_0kyymfy",
              initiatingParticipant: "Participant_0x6v44d",
              respondingParticipant: "Participant_0n7kwiu",
              initialMessage: "Message_1376fyb"
            },
            {
              id: "ChoreographyTask_0nl2rhr",
              name: "Purchase raw materials",
              incoming: "Flow_1k5cri3",
              outgoing: "Flow_1oa3jbg",
              initiatingParticipant: "Participant_0n7kwiu",
              respondingParticipant: "Participant_0xxs9a7",
              initialMessage: "Message_0xe03aa",
              responseMessage: "Message_1p1ke3y",
            },
            {
              id: "ChoreographyTask_1uie9z3",
              name: "Confirm order",
              incoming: "Flow_0rtcg7d",
              outgoing: "Flow_1iwbmcz",
              initiatingParticipant: "Participant_0n7kwiu",
              respondingParticipant: "Participant_0x6v44d",
              initialMessage: "Message_1dd2uoz",
            },
            {
              id: "ChoreographyTask_1dsovf5",
              name: "Ship product",
              incoming: "Flow_1jplxjf",
              outgoing: "Flow_0glm8k3",
              initiatingParticipant: "Participant_0x6v44d",
              respondingParticipant: "Participant_0n7kwiu",
              initialMessage: "Message_1858jqq",
              responseMessage: "Message_01cp9ki"
            },
            {
              id: "ChoreographyTask_1htg6wy",
              name: "Invoice customer",
              incoming: "Flow_1h69rwb",
              outgoing: "Flow_0433489",
              initiatingParticipant: "Participant_0n7kwiu",
              respondingParticipant: "Participant_0x6v44d",
              initialMessage: "Message_0annsni",
              responseMessage: "Message_05ebl37"
            }
          ],
          sequenceFlows: [
            {
              id: "Flow_1f6eaf2"
            },
            {
              id: "Flow_0kyymfy"
            },
            {
              id: "Flow_1k5cri3",
              name: "product not in stock"
            },
            {
              id: "Flow_1oa3jbg"
            },
            {
              id: "Flow_0572p2k",
              name: "product in stock"
            },
            {
              id: "Flow_1iwbmcz"
            },
            {
              id: "Flow_1jplxjf"
            },
            {
              id: "Flow_1h69rwb"
            },
            {
              id: "Flow_0glm8k3"
            },
            {
              id: "Flow_0433489"
            },
            {
              id: "Flow_0vbgnyw"
            },
            {
              id: "Flow_0rtcg7d"
            }
          ],
          messages: [
            {
              id: "Message_05ebl37",
              name: "Payment"
            },
            {
              id: "Message_01cp9ki",
              name: "Product"
            },
            {
              id: "Message_1p1ke3y",
              name: "Raw materials"
            },
            {
              id: "Message_0annsni",
              name: "Invoice"
            },
            {
              id: "Message_1858jqq",
              name: "Shipping address"
            },
            {
              id: "Message_1dd2uoz",
              name: "Confirm"
            },
            {
              id: "Message_0xe03aa",
              name: "Raw materials request"
            },
            {
              id: "Message_1376fyb",
              name: "Purchase order"
            }
          ]
        }
      ]
    };
  }

  static getModel2(): Model {
    return {
      id: "Choreography_07n6r3q",
      placeCount: 18,
      participantCount: 2,
      messageCount: 8,
      startPlace: 16,
      endPlaces: [17],
      transitions: [
        {
          id: "Event_1525yky",
          type: TransitionType.REQUIRED,
          incomingPlaces: [16],
          outgoingPlaces: [0]
        },
        {
          id: "Event_08d32d7",
          type: TransitionType.REQUIRED,
          name: "Order fullfilled",
          incomingPlaces: [10],
          outgoingPlaces: [17]
        },
        {
          id: "Gateway_10fv7g5_Flow_1k5cri3",
          type: TransitionType.OPTIONAL_OUTGOING,
          incomingPlaces: [1],
          outgoingPlaces: [2]
        },
        {
          id: "Gateway_10fv7g5_Flow_0572p2k",
          type: TransitionType.OPTIONAL_OUTGOING,
          incomingPlaces: [1],
          outgoingPlaces: [4]
        },
        {
          id: "Gateway_035usvb_Flow_0572p2k",
          type: TransitionType.OPTIONAL_INCOMING,
          incomingPlaces: [4],
          outgoingPlaces: [11]
        },
        {
          id: "Gateway_035usvb_Flow_1oa3jbg",
          type: TransitionType.OPTIONAL_INCOMING,
          incomingPlaces: [3],
          outgoingPlaces: [11]
        },
        {
          id: "Gateway_0g1gxjc",
          type: TransitionType.OPTIONAL_INCOMING,
          incomingPlaces: [5],
          outgoingPlaces: [6, 7]
        },
        {
          id: "Gateway_0nm7ind",
          type: TransitionType.OPTIONAL_OUTGOING,
          incomingPlaces: [8, 9],
          outgoingPlaces: [10]
        },
        {
          id: "ChoreographyTask_0kp4flv_Participant_0x6v44d",
          type: TransitionType.REQUIRED,
          name: "Submit purchase order",
          incomingPlaces: [0],
          outgoingPlaces: [12],
          participant: 0,
          message: 7
        },
        {
          id: "ChoreographyTask_0kp4flv_Participant_0n7kwiu",
          type: TransitionType.REQUIRED,
          name: "Submit purchase order",
          incomingPlaces: [12],
          outgoingPlaces: [1],
          participant: 1
        },
        {
          id: "ChoreographyTask_0nl2rhr_loop",
          type: TransitionType.REQUIRED,
          name: "Purchase raw materials",
          incomingPlaces: [2],
          outgoingPlaces: [2],
          participant: 1,
          message: 6,
        },
        {
          id: "ChoreographyTask_0nl2rhr_end",
          type: TransitionType.REQUIRED,
          name: "Purchase raw materials",
          incomingPlaces: [2],
          outgoingPlaces: [3],
          participant: 1,
          message: 6,
        },
        {
          id: "ChoreographyTask_1uie9z3_Participant_0n7kwiu",
          type: TransitionType.REQUIRED,
          name: "Confirm order",
          incomingPlaces: [11],
          outgoingPlaces: [13],
          participant: 1,
          message: 5,
        },
        {
          id: "ChoreographyTask_1uie9z3_Participant_0x6v44d",
          type: TransitionType.REQUIRED,
          name: "Confirm order",
          incomingPlaces: [13],
          outgoingPlaces: [5],
          participant: 0,
        },
        {
          id: "ChoreographyTask_1dsovf5_Participant_0x6v44d",
          type: TransitionType.REQUIRED,
          name: "Ship product",
          incomingPlaces: [6],
          outgoingPlaces: [14],
          participant: 0,
          message: 4,
        },
        {
          id: "ChoreographyTask_1dsovf5_Participant_0n7kwiu",
          type: TransitionType.REQUIRED,
          name: "Ship product",
          incomingPlaces: [14],
          outgoingPlaces: [8],
          participant: 1,
          message: 1,
        },
        {
          id: "ChoreographyTask_1htg6wy_Participant_0n7kwiu",
          type: TransitionType.REQUIRED,
          name: "Invoice customer",
          incomingPlaces: [7],
          outgoingPlaces: [15],
          participant: 1,
          message: 3,
        },
        {
          id: "ChoreographyTask_1htg6wy_Participant_0x6v44d",
          type: TransitionType.REQUIRED,
          name: "Invoice customer",
          incomingPlaces: [15],
          outgoingPlaces: [9],
          participant: 0,
          message: 0,
        },
      ]
    };
  }

  static getModel2Reduced(): Model {
    return {
      id: "Choreography_07n6r3q",
      placeCount: 12,
      participantCount: 2,
      messageCount: 8,
      startPlace: 10,
      endPlaces: [11],
      transitions: [
        {
          id: "Event_1525yky",
          type: TransitionType.REQUIRED,
          incomingPlaces: [10],
          outgoingPlaces: [0]
        },
        {
          id: "Event_08d32d7",
          type: TransitionType.REQUIRED,
          name: "Order fullfilled",
          incomingPlaces: [3, 4],
          outgoingPlaces: [11]
        },
        {
          id: "ChoreographyTask_0kp4flv_Participant_0x6v44d",
          type: TransitionType.REQUIRED,
          name: "Submit purchase order",
          incomingPlaces: [0],
          outgoingPlaces: [6],
          participant: 0,
          message: 7
        },
        {
          id: "ChoreographyTask_0kp4flv_Participant_0n7kwiu",
          type: TransitionType.REQUIRED,
          name: "Submit purchase order",
          incomingPlaces: [6],
          outgoingPlaces: [5],
          participant: 1
        },
        {
          id: "ChoreographyTask_0nl2rhr_loop",
          type: TransitionType.REQUIRED,
          name: "Purchase raw materials",
          incomingPlaces: [5],
          outgoingPlaces: [5],
          participant: 1,
          message: 6,
        },
        {
          id: "ChoreographyTask_0nl2rhr_end",
          type: TransitionType.REQUIRED,
          name: "Purchase raw materials",
          incomingPlaces: [5],
          outgoingPlaces: [5],
          participant: 1,
          message: 6,
        },
        {
          id: "ChoreographyTask_1uie9z3_Participant_0n7kwiu",
          type: TransitionType.REQUIRED,
          name: "Confirm order",
          incomingPlaces: [5],
          outgoingPlaces: [7],
          participant: 1,
          message: 5,
        },
        {
          id: "ChoreographyTask_1uie9z3_Participant_0x6v44d",
          type: TransitionType.REQUIRED,
          name: "Confirm order",
          incomingPlaces: [7],
          outgoingPlaces: [1, 2],
          participant: 0,
        },
        {
          id: "ChoreographyTask_1dsovf5_Participant_0x6v44d",
          type: TransitionType.REQUIRED,
          name: "Ship product",
          incomingPlaces: [1],
          outgoingPlaces: [8],
          participant: 0,
          message: 4,
        },
        {
          id: "ChoreographyTask_1dsovf5_Participant_0n7kwiu",
          type: TransitionType.REQUIRED,
          name: "Ship product",
          incomingPlaces: [8],
          outgoingPlaces: [3],
          participant: 1,
          message: 1,
        },
        {
          id: "ChoreographyTask_1htg6wy_Participant_0n7kwiu",
          type: TransitionType.REQUIRED,
          name: "Invoice customer",
          incomingPlaces: [2],
          outgoingPlaces: [9],
          participant: 1,
          message: 3,
        },
        {
          id: "ChoreographyTask_1htg6wy_Participant_0x6v44d",
          type: TransitionType.REQUIRED,
          name: "Invoice customer",
          incomingPlaces: [9],
          outgoingPlaces: [4],
          participant: 0,
          message: 0,
        },
      ]
    };
  }
}
