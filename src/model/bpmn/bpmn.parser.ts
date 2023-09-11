import { Injectable } from '@nestjs/common';
import { XMLParser } from 'fast-xml-parser';
import { Definitions, Process, Element, SequenceFlow } from './bpmn';
import { TransitionType } from '../domain/model';

@Injectable()
export class BpmnParser {
  private readonly definitionsTag = 'bpmn:definitions';
  private readonly processTag = 'bpmn:process';
  private readonly startEventTag = 'bpmn:startEvent';
  private readonly endEventTag = 'bpmn:endEvent';
  private readonly taskTag = 'bpmn:task';
  private readonly exclusiveGatewayTag = 'bpmn:exclusiveGateway';
  private readonly parallelGatewayTag = 'bpmn:parallelGateway';
  private readonly incomingTag = 'bpmn:incoming';
  private readonly outgoingTag = 'bpmn:outgoing';
  private readonly sequenceFlowTag = 'bpmn:sequenceFlow';

  private options = {
    attributeNamePrefix: '',
    ignoreAttributes: false,
    isArray: (tagName: string) =>
      [
        this.endEventTag,
        this.taskTag,
        this.exclusiveGatewayTag,
        this.parallelGatewayTag,
        this.incomingTag,
        this.outgoingTag,
      ].includes(tagName),
  };
  private parser = new XMLParser(this.options);

  parseBpmn(bpmnString: string): Definitions {
    const definitions = this.parser.parse(bpmnString)[this.definitionsTag];
    const process = definitions[this.processTag];
    return {
      process: this.parseProcess(process),
    };
  }

  parseProcess(process: any): Process {
    const startEvent = this.parseStartEvent(process);
    const endEvents = this.parseEndEvents(process);
    const tasks = this.parseTasks(process);
    const exclusiveGateways = this.parseExclusiveGateways(process);
    const parallelGateways = this.parseParallelGateways(process);
    const sequenceFlows = this.parseSequenceFlows(process);
    return {
      id: process.id,
      startEvent,
      endEvents,
      tasks,
      exclusiveGateways,
      parallelGateways,
      sequenceFlows
    };
  }

  private parseStartEvent(process: any): Element {
    const startEvent = process[this.startEventTag];
    return {
      id: startEvent.id,
      type: TransitionType.START,
      name: startEvent.name,
      incomingSequenceFlows: [],
      outgoingSequenceFlows: startEvent[this.outgoingTag],
    };
  }

  private parseEndEvents(process: any): Element[] {
    const endEvents = process[this.endEventTag];
    return endEvents.map((endEvent: any) => ({
      id: endEvent.id,
      type: TransitionType.END,
      name: endEvent.name,
      incomingSequenceFlows: endEvent[this.incomingTag],
      outgoingSequenceFlows: [],
    }));
  }

  private parseTasks(process: any): Element[] {
    const tasks = process[this.taskTag];
    return tasks.map((task: any) => ({
      id: task.id,
      type: TransitionType.TASK,
      name: task.name,
      incomingSequenceFlows: task[this.incomingTag],
      outgoingSequenceFlows: task[this.outgoingTag],
    }));
  }

  private parseExclusiveGateways(process: any): Element[] {
    const exclusiveGateways = process[this.exclusiveGatewayTag];
    return exclusiveGateways.map((exclusiveGateway: any) => {
      const incomingPlaceIds: any[] = exclusiveGateway[this.incomingTag];
      const outgoingPlaceIds: any[] = exclusiveGateway[this.outgoingTag];
      return {
        id: exclusiveGateway.id,
        type:
          incomingPlaceIds.length > 1
            ? TransitionType.XOR_JOIN
            : TransitionType.XOR_SPLIT,
        incomingSequenceFlows: incomingPlaceIds,
        outgoingSequenceFlows: outgoingPlaceIds,
      };
    });
  }
  private parseParallelGateways(process: any): Element[] {
    const parallelGateways = process[this.parallelGatewayTag];
    return parallelGateways.map((parallelGateway: any) => {
      const incomingPlaceIds: any[] = parallelGateway[this.incomingTag];
      const outgoingPlaceIds: any[] = parallelGateway[this.outgoingTag];
      return {
        id: parallelGateway.id,
        type:
          incomingPlaceIds.length > 1
            ? TransitionType.AND_JOIN
            : TransitionType.AND_SPLIT,
        incomingSequenceFlows: incomingPlaceIds,
        outgoingSequenceFlows: outgoingPlaceIds,
      };
    });
  }

  parseSequenceFlows(process: any): SequenceFlow[] {
    const sequenceFlows = process[this.sequenceFlowTag];
    return sequenceFlows.map((sequenceFlow: any) => ({
      id: sequenceFlow.id
    }))
  }
}
