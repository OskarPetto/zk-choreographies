import { Injectable } from '@nestjs/common';
import { Element, Process, SequenceFlow, SequenceFlowId } from './bpmn';
import {
  Model,
  PlaceId,
  Transition,
  TransitionType,
} from '../model/model';

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
    const sequenceFlowPlaceIds = this.createSequenceFlowMapping(process.sequenceFlows);
    const transitions = elements.flatMap((element) =>
      this.toTransitions(sequenceFlowPlaceIds, element),
    );
    return {
      id: process.id,
      startPlace: sequenceFlowPlaceIds.size,
      endPlace: sequenceFlowPlaceIds.size + 1,
      placeCount: sequenceFlowPlaceIds.size + 2,
      participantCount: 1,
      transitions,
    };
  }

  private toTransitions(
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
    element: Element,
  ): Transition[] {
    switch (element.type) {
      case TransitionType.START:
        return this.startToTransitions(element, sequenceFlowPlaceIds);
      case TransitionType.END:
        return this.endToTransitions(element, sequenceFlowPlaceIds);
      case TransitionType.XOR_SPLIT:
        return this.xorSplitToTransitions(element, sequenceFlowPlaceIds);
      case TransitionType.XOR_JOIN:
        return this.xorJoinToTransitions(element, sequenceFlowPlaceIds);
      default:
        return this.elementToTransitions(element, sequenceFlowPlaceIds);
    }
  }

  private startToTransitions(
    element: Element,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const outgoingPlaces = element.outgoingSequenceFlows.map(
      (sequenceFlowId) => sequenceFlowPlaceIds.get(sequenceFlowId)!,
    );
    return [
      {
        id: element.id,
        type: element.type,
        name: element.name,
        incomingPlaces: [sequenceFlowPlaceIds.size],
        outgoingPlaces,
      },
    ];
  }

  private endToTransitions(
    element: Element,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const incomingPlaces = element.incomingSequenceFlows.map(
      (sequenceFlowId) => sequenceFlowPlaceIds.get(sequenceFlowId)!,
    );
    return [
      {
        id: element.id,
        type: element.type,
        name: element.name,
        incomingPlaces,
        outgoingPlaces: [sequenceFlowPlaceIds.size + 1],
      },
    ];
  }

  private elementToTransitions(
    element: Element,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const incomingPlaces = element.incomingSequenceFlows.map(
      (sequenceFlowId) => sequenceFlowPlaceIds.get(sequenceFlowId)!,
    );
    const outgoingPlaces = element.outgoingSequenceFlows.map(
      (sequenceFlowId) => sequenceFlowPlaceIds.get(sequenceFlowId)!,
    );
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

  private xorJoinToTransitions(
    element: Element,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const outgoingPlaces = element.outgoingSequenceFlows.map(
      (sequenceFlowId) => sequenceFlowPlaceIds.get(sequenceFlowId)!,
    );
    return element.incomingSequenceFlows.map((incomingSequenceFlowId) => ({
      id: `${element.id}_${incomingSequenceFlowId}`,
      type: element.type,
      name: element.name,
      incomingPlaces: [sequenceFlowPlaceIds.get(incomingSequenceFlowId)!],
      outgoingPlaces: outgoingPlaces,
    }));
  }

  private xorSplitToTransitions(
    element: Element,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const incomingPlaces = element.incomingSequenceFlows.map(
      (sequenceFlowId) => sequenceFlowPlaceIds.get(sequenceFlowId)!,
    );
    return element.outgoingSequenceFlows.map((outgoingSequenceFlowId) => ({
      id: `${element.id}_${outgoingSequenceFlowId}`,
      type: element.type,
      name: element.name,
      incomingPlaces: incomingPlaces,
      outgoingPlaces: [sequenceFlowPlaceIds.get(outgoingSequenceFlowId)!],
    }));
  }

  private createSequenceFlowMapping(
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
