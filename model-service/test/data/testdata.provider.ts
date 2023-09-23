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
  static getBarbaraReChoreography(): string {
    return readTextFile('test/data/barbara_re_choreography.bpmn');
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
}
