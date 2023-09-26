import { Injectable } from '@nestjs/common';
import { XMLParser } from 'fast-xml-parser';
import {
  Definitions,
  SequenceFlow,
  Message,
  Choreography,
  Participant,
  StartEvent,
  EndEvent,
  GatewayType,
  ChoreographyTask,
  ExclusiveGateway,
  ParallelGateway,
  LoopType,
} from './bpmn';

@Injectable()
export class BpmnParser {
  private readonly definitionsTag = 'bpmn2:definitions';
  private readonly messageTag = 'bpmn2:message';
  private readonly choreographyTag = 'bpmn2:choreography';
  private readonly participantTag = 'bpmn2:participant';
  private readonly participantMultiplicityTag = 'bpmn2:participantMultiplicity';
  private readonly sequenceFlowTag = 'bpmn2:sequenceFlow';
  private readonly messageFlowTag = 'bpmn2:messageFlow';
  private readonly startEventTag = 'bpmn2:startEvent';
  private readonly endEventTag = 'bpmn2:endEvent';
  private readonly exclusiveGatewayTag = 'bpmn2:exclusiveGateway';
  private readonly parallelGatewayTag = 'bpmn2:parallelGateway';
  private readonly choreographyTaskTag = 'bpmn2:choreographyTask';
  private readonly incomingTag = 'bpmn2:incoming';
  private readonly outgoingTag = 'bpmn2:outgoing';
  private readonly participantRefTag = 'bpmn2:participantRef';
  private readonly messageFlowRefTag = 'bpmn2:messageFlowRef';

  private options = {
    attributeNamePrefix: '',
    ignoreAttributes: false,
    isArray: (tagName: string) =>
      [
        this.messageTag,
        this.choreographyTag,
        this.participantTag,
        this.startEventTag,
        this.endEventTag,
        this.choreographyTaskTag,
        this.exclusiveGatewayTag,
        this.parallelGatewayTag,
        this.incomingTag,
        this.outgoingTag,
        this.messageFlowRefTag,
        this.sequenceFlowTag,
        this.messageFlowTag,
      ].includes(tagName),
  };
  private parser = new XMLParser(this.options);

  parseBpmn(bpmnString: string): Definitions {
    const definitions = this.parser.parse(bpmnString)[this.definitionsTag];
    return {
      choreographies: this.parseChoreographies(definitions),
    };
  }

  parseChoreographies(definitions: any): Choreography[] {
    const choreographies = definitions[this.choreographyTag];
    const messages = this.parseMessages(definitions);
    return choreographies.map((choreography: any) =>
      this.parseChoreography(choreography, messages),
    );
  }

  parseMessages(definitions: any): Message[] {
    const messages = definitions[this.messageTag];
    return messages.map((message: any) => ({
      id: message.id,
      name: message.name,
    }));
  }

  parseChoreography(choreography: any, messages: Message[]): Choreography {
    const sequenceFlows = this.parseSequenceFlows(choreography);
    const participants = this.parseParticipants(choreography);
    const startEvents = this.parseStartEvents(choreography);
    const endEvents = this.parseEndEvents(choreography);
    const exclusiveGateways = this.parseExclusiveGateways(choreography);
    const parallelGateways = this.parseParallelGateways(choreography);
    const choreographyTasks = this.parseChoreographyTasks(choreography);
    const relevantMessages = messages.filter((message: Message) =>
      choreographyTasks.find(
        (choreographyTask: ChoreographyTask) =>
          choreographyTask.initialMessage === message.id ||
          choreographyTask.responseMessage === message.id,
      ),
    );
    const relevantParticipants = participants.filter(
      (participant: Participant) =>
        choreographyTasks.some(
          (choreographyTask: ChoreographyTask) =>
            choreographyTask.initiatingParticipant === participant.id ||
            choreographyTask.respondingParticipant === participant.id,
        ),
    );
    return {
      id: choreography.id,
      sequenceFlows,
      participants: relevantParticipants,
      startEvents,
      endEvents,
      exclusiveGateways,
      parallelGateways,
      choreographyTasks,
      messages: relevantMessages,
    };
  }

  parseSequenceFlows(choreography: any): SequenceFlow[] {
    const sequenceFlows = choreography[this.sequenceFlowTag];
    return sequenceFlows.map((sequenceFlow: any) => ({
      id: sequenceFlow.id,
      name: sequenceFlow.name,
    }));
  }

  private parseParticipants(choreography: any): Participant[] {
    return choreography[this.participantTag].map((participant: any) => {
      return {
        id: participant.id,
        name: participant.name,
      };
    });
  }

  private parseStartEvents(choreography: any): StartEvent[] {
    const startEvents = choreography[this.startEventTag];
    return startEvents.map((startEvent: any) => ({
      id: startEvent.id,
      name: startEvent.name,
      outgoing: startEvent[this.outgoingTag][0],
    }));
  }

  private parseEndEvents(choreography: any): EndEvent[] {
    const endEvents = choreography[this.endEventTag];
    return endEvents.map((endEvent: any) => ({
      id: endEvent.id,
      name: endEvent.name,
      incoming: endEvent[this.incomingTag][0],
    }));
  }

  private parseExclusiveGateways(choreography: any): ExclusiveGateway[] {
    const exclusiveGateways = choreography[this.exclusiveGatewayTag];
    return exclusiveGateways.map((exclusiveGateway: any) => {
      const incoming: any[] = exclusiveGateway[this.incomingTag];
      const outgoing: any[] = exclusiveGateway[this.outgoingTag];
      return {
        id: exclusiveGateway.id,
        type: incoming.length > 1 ? GatewayType.JOIN : GatewayType.SPLIT,
        default: exclusiveGateway.default,
        incoming,
        outgoing,
      };
    });
  }
  private parseParallelGateways(choreography: any): ParallelGateway[] {
    const parallelGateways = choreography[this.parallelGatewayTag];
    return parallelGateways.map((parallelGateway: any) => {
      const incoming: any[] = parallelGateway[this.incomingTag];
      const outgoing: any[] = parallelGateway[this.outgoingTag];
      return {
        id: parallelGateway.id,
        type: incoming.length > 1 ? GatewayType.JOIN : GatewayType.SPLIT,
        incoming,
        outgoing,
      };
    });
  }

  private parseChoreographyTasks(choreography: any): ChoreographyTask[] {
    const allMessageFlows: Map<string, any> = new Map(
      choreography[this.messageFlowTag].map((messageFlow: any) => [
        messageFlow.id,
        messageFlow,
      ]),
    );
    const choreographyTasks = choreography[this.choreographyTaskTag];
    return choreographyTasks.map((choreographyTask: any) => {
      const participantRefs = choreographyTask[this.participantRefTag];
      const initiatingParticipantRef =
        choreographyTask.initiatingParticipantRef;
      const respondingParticipantRef = participantRefs.find(
        (participantRef: any) =>
          participantRef != choreographyTask.initiatingParticipantRef,
      )!;
      const messageFlows = choreographyTask[this.messageFlowRefTag].map(
        (messageFlowRef: any) => allMessageFlows.get(messageFlowRef)!,
      );
      const initialMessageRef = messageFlows.find(
        (messageFlow: any) =>
          messageFlow.sourceRef === initiatingParticipantRef,
      )?.messageRef;
      const responseMessageRef = messageFlows.find(
        (messageFlow: any) =>
          messageFlow.sourceRef === respondingParticipantRef,
      )?.messageRef;

      return {
        id: choreographyTask.id,
        name: choreographyTask.name,
        incoming: choreographyTask[this.incomingTag][0],
        outgoing: choreographyTask[this.outgoingTag][0],
        initiatingParticipant: initiatingParticipantRef,
        respondingParticipant: respondingParticipantRef,
        initialMessage: initialMessageRef,
        responseMessage: responseMessageRef,
        loopType: this.parseLoopType(choreographyTask),
      };
    });
  }

  private parseLoopType(choreographyTask: any): LoopType | undefined {
    switch (choreographyTask.loopType) {
      case "Standard":
        return LoopType.STANDARD;
      case "MultiInstanceSequential":
        return LoopType.MULTI_INSTANCE_SEQUENTIAL;
      case "MultiInstanceParallel":
        return LoopType.MULTI_INSTANCE_PARALLEL;
      default:
        return undefined;
    }
  }
}
