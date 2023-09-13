import { Injectable } from '@nestjs/common';
import { Element, Process, SequenceFlow, SequenceFlowId } from './bpmn';
import {
  PetriNet,
  PlaceId,
  Transition,
  TransitionType,
} from '../petri-net/petri-net';

@Injectable()
export class BpmnMapper {
  toPetriNet(process: Process): PetriNet {
    const elements = [
      process.startEvent,
      ...process.endEvents,
      ...process.tasks,
      ...process.exclusiveGateways,
      ...process.parallelGateways,
    ];
    const placeIds = this.createSequenceFlowMapping(process.sequenceFlows);
    const transitions = elements.flatMap((element) =>
      this.toTransitions(placeIds, element),
    );
    return {
      id: process.id,
      startPlace: placeIds.size,
      placeCount: placeIds.size + 1,
      transitions,
    };
  }

  private toTransitions(
    placeIds: Map<SequenceFlowId, PlaceId>,
    element: Element,
  ): Transition[] {
    switch (element.type) {
      case TransitionType.START:
        return this.startToTransitions(element, placeIds);
      case TransitionType.XOR_SPLIT:
        return this.xorSplitToTransitions(element, placeIds);
      case TransitionType.XOR_JOIN:
        return this.xorJoinToTransitions(element, placeIds);
      default:
        return this.elementToTransitions(element, placeIds);
    }
  }

  private startToTransitions(
    element: Element,
    placeIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const outgoingPlaces = element.outgoingSequenceFlows.map(
      (sequenceFlowId) => placeIds.get(sequenceFlowId)!,
    );
    return [
      {
        id: element.id,
        type: element.type,
        name: element.name,
        incomingPlaces: [placeIds.size],
        outgoingPlaces,
      },
    ];
  }

  private elementToTransitions(
    element: Element,
    placeIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const incomingPlaces = element.incomingSequenceFlows.map(
      (sequenceFlowId) => placeIds.get(sequenceFlowId)!,
    );
    const outgoingPlaces = element.outgoingSequenceFlows.map(
      (sequenceFlowId) => placeIds.get(sequenceFlowId)!,
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
    placeIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const outgoingPlaces = element.outgoingSequenceFlows.map(
      (sequenceFlowId) => placeIds.get(sequenceFlowId)!,
    );
    return element.incomingSequenceFlows.map((incomingSequenceFlowId) => ({
      id: `${element.id}_${incomingSequenceFlowId}`,
      type: element.type,
      name: element.name,
      incomingPlaces: [placeIds.get(incomingSequenceFlowId)!],
      outgoingPlaces: outgoingPlaces,
    }));
  }

  private xorSplitToTransitions(
    element: Element,
    placeIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const incomingPlaces = element.incomingSequenceFlows.map(
      (sequenceFlowId) => placeIds.get(sequenceFlowId)!,
    );
    return element.outgoingSequenceFlows.map((outgoingSequenceFlowId) => ({
      id: `${element.id}_${outgoingSequenceFlowId}`,
      type: element.type,
      name: element.name,
      incomingPlaces: incomingPlaces,
      outgoingPlaces: [placeIds.get(outgoingSequenceFlowId)!],
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
