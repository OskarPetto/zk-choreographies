import { Injectable } from '@nestjs/common';
import {
  Element,
  ElementType,
  FlowId,
  Model,
  createModelId,
} from '../domain/model';
import { XMLParser } from 'fast-xml-parser';

@Injectable()
export class BpmnService {
  private options = {
    attributeNamePrefix: '',
    ignoreAttributes: false,
  };
  private parser = new XMLParser(this.options);

  parseProcess(bpmnString: string): Model {
    const definitions = this.parser.parse(bpmnString)['bpmn:definitions'];

    const process = definitions['bpmn:process'];
    const flows: FlowId[] = this.parseSequenceFlows(process);
    const startEvent: Element = this.parseStartEvent(process);
    const endEvents: Element[] = this.parseEndEvents(process);
    const tasks: Element[] = this.parseTasks(process);
    const exclusiveGateways: Element[] = this.parseExclusiveGateways(process);
    const parallelGateways: Element[] = this.parseParallelGateways(process);
    const elements = [
      startEvent,
      ...endEvents,
      ...tasks,
      ...exclusiveGateways,
      ...parallelGateways,
    ];
    return {
      id: createModelId(),
      flows,
      elements: new Map(elements.map((e) => [e.id, e])),
    };
  }

  private parseStartEvent(process: any): Element {
    const startEvent = process['bpmn:startEvent'];
    return {
      id: startEvent.id,
      type: ElementType.START,
      name: startEvent.name,
      incomingFlows: [],
      outgoingFlows: [startEvent['bpmn:outgoing']],
    };
  }

  private parseEndEvents(process: any): Element[] {
    const endEvents: any[] = [].concat(process['bpmn:endEvent']);
    return endEvents.map((endEvent: any) => ({
      id: endEvent.id,
      type: ElementType.END,
      name: endEvent.name,
      incomingFlows: [endEvent['bpmn:incoming']],
      outgoingFlows: [],
    }));
  }

  private parseTasks(process: any): Element[] {
    const tasks: any[] = [].concat(process['bpmn:task']);
    return tasks.map((task: any) => ({
      id: task.id,
      type: ElementType.TASK,
      name: task.name,
      incomingFlows: [task['bpmn:incoming']],
      outgoingFlows: [task['bpmn:outgoing']],
    }));
  }

  private parseExclusiveGateways(process: any): Element[] {
    const exclusiveGateways = process['bpmn:exclusiveGateway'];
    return exclusiveGateways.flatMap((exclusiveGateway: any) =>
      this.parseExclusiveGateway(exclusiveGateway),
    );
  }

  private parseExclusiveGateway(exclusiveGateway: any): Element[] {
    const incomingFlowIds = exclusiveGateway['bpmn:incoming'];
    const outgoingFlowIds = exclusiveGateway['bpmn:outgoing'];

    if (Array.isArray(incomingFlowIds)) {
      return incomingFlowIds.map((incomingFlowId) => ({
        id: `${exclusiveGateway.id}_${incomingFlowId}`,
        type: ElementType.XOR_JOIN,
        incomingFlows: [incomingFlowId],
        outgoingFlows: [outgoingFlowIds],
      }));
    } else if (Array.isArray(outgoingFlowIds)) {
      return outgoingFlowIds.map((outgoingFlowId) => ({
        id: `${exclusiveGateway.id}_${outgoingFlowId}`,
        type: ElementType.XOR_SPLIT,
        incomingFlows: [incomingFlowIds],
        outgoingFlows: [outgoingFlowId],
      }));
    } else {
      return [];
    }
  }

  private parseParallelGateways(process: any): Element[] {
    const parallelGateways = process['bpmn:parallelGateway'];
    return parallelGateways.map((parallelGateway: any) => {
      const incomingFlowIds = parallelGateway['bpmn:incoming'];
      const outgoingFlowIds = parallelGateway['bpmn:outgoing'];
      return {
        id: parallelGateway.id,
        type: Array.isArray(incomingFlowIds)
          ? ElementType.AND_JOIN
          : ElementType.AND_SPLIT,
        incomingFlows: [].concat(incomingFlowIds),
        outgoingFlows: [].concat(outgoingFlowIds),
      };
    });
  }

  private parseSequenceFlows(process: any): FlowId[] {
    const sequenceFlows = [].concat(process['bpmn:sequenceFlow']);
    return sequenceFlows.map((sequenceFlow: any) => sequenceFlow.id as FlowId);
  }
}
