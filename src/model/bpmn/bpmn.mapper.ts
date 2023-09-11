import { Injectable } from '@nestjs/common';
import { Element, Process } from './bpmn';
import { Model, Transition, TransitionType } from '../domain/model';

@Injectable()
export class BpmnMapper {
  toModel(process: Process): Model {
    const elements = [
      process.startEvent,
      ...process.endEvents,
      ...process.tasks,
      ...process.exclusiveGateways,
      ...process.parallelGateways,
    ];
    const transitions = elements.flatMap((element) =>
      this.toTransitions(element),
    );
    return {
      id: process.id,
      transitions,
    };
  }

  private toTransitions(element: Element): Transition[] {
    if (element.type === TransitionType.XOR_SPLIT) {
      return element.outgoingSequenceFlows.map((outgoingSequenceFlowId) => ({
        id: `${element.id}_${outgoingSequenceFlowId}`,
        type: element.type,
        name: element.name,
        incomingPlaces: element.incomingSequenceFlows,
        outgoingPlaces: [outgoingSequenceFlowId],
      }));
    } else if (element.type === TransitionType.XOR_JOIN) {
      return element.incomingSequenceFlows.map((incomingSequenceFlowId) => ({
        id: `${element.id}_${incomingSequenceFlowId}`,
        type: element.type,
        name: element.name,
        incomingPlaces: [incomingSequenceFlowId],
        outgoingPlaces: element.outgoingSequenceFlows,
      }));
    } else {
      return [
        {
          id: element.id,
          type: element.type,
          name: element.name,
          incomingPlaces: element.incomingSequenceFlows,
          outgoingPlaces: element.outgoingSequenceFlows,
        },
      ];
    }
  }
}
