import { Injectable } from '@nestjs/common';
import { Element, Process, SequenceFlow, SequenceFlowId } from './bpmn';
import { Model, PlaceId, Transition, TransitionType } from '../domain/model';

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
    const placeIds = this.createPlaceIdMapping(process.sequenceFlows);
    const transitions = elements.flatMap((element) =>
      this.toTransitions(placeIds, element),
    );
    return {
      id: process.id,
      placeCount: placeIds.size,
      transitions,
    };
  }

  private toTransitions(
    placeIds: Map<SequenceFlowId, PlaceId>,
    element: Element,
  ): Transition[] {
    const incomingPlaces = element.incomingSequenceFlows.map(
      (sequenceFlowId) => placeIds.get(sequenceFlowId)!,
    );
    const outgoingPlaces = element.outgoingSequenceFlows.map(
      (sequenceFlowId) => placeIds.get(sequenceFlowId)!,
    );

    if (element.type === TransitionType.XOR_SPLIT) {
      return element.outgoingSequenceFlows.map((outgoingSequenceFlowId) => ({
        id: `${element.id}_${outgoingSequenceFlowId}`,
        type: element.type,
        name: element.name,
        incomingPlaces: incomingPlaces,
        outgoingPlaces: [placeIds.get(outgoingSequenceFlowId)!],
      }));
    } else if (element.type === TransitionType.XOR_JOIN) {
      return element.incomingSequenceFlows.map((incomingSequenceFlowId) => ({
        id: `${element.id}_${incomingSequenceFlowId}`,
        type: element.type,
        name: element.name,
        incomingPlaces: [placeIds.get(incomingSequenceFlowId)!],
        outgoingPlaces: outgoingPlaces,
      }));
    } else {
      return [
        {
          id: element.id,
          type: element.type,
          name: element.name,
          incomingPlaces,
          outgoingPlaces,
        },
      ];
    }
  }

  private createPlaceIdMapping(
    sequenceFlows: SequenceFlow[],
  ): Map<SequenceFlowId, PlaceId> {
    const map: Map<SequenceFlowId, PlaceId> = new Map();
    let index = 0;
    for (const sequenceFlow of sequenceFlows) {
      map.set(sequenceFlow.id, index++);
    }
    return map;
  }
}
