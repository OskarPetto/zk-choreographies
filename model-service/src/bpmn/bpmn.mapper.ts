import { Injectable } from '@nestjs/common';
import { ChoreographyTask, EndEvent, ExclusiveGateway, GatewayType, ParallelGateway, Participant, SequenceFlow, SequenceFlowId, StartEvent, ParticipantId as BpmnParticipantId, MessageId as BpmnMessageId, Message, Choreography } from './bpmn';
import {
  Model,
  ParticipantId,
  MessageId,
  PlaceId,
  Transition,
  TransitionType,
} from '../model/model';
import { v4 as uuid } from 'uuid';

@Injectable()
export class BpmnMapper {
  toModel(choreography: Choreography): Model {
    const sequenceFlowPlaceIds = this.createSequenceFlowMapping(choreography.sequenceFlows);
    const participantIds = this.createParticipantMapping(choreography.participants);
    const messageIds = this.createMessageMapping(choreography.messages);
    const additionalPlaceIds: PlaceId[] = [];
    const choreographyTaskTransitions = choreography.choreographyTasks.flatMap(choreographyTask => this.choreographyTaskToTransitions(choreographyTask, sequenceFlowPlaceIds, participantIds, messageIds, additionalPlaceIds))

    const startTransition = this.startEventToTransition(choreography.startEvent, sequenceFlowPlaceIds, additionalPlaceIds);
    const endTransitions = choreography.endEvents.flatMap(endEvent => this.endEventToTransitions(endEvent, sequenceFlowPlaceIds, additionalPlaceIds));
    const exclusiveGatewayTransitions = choreography.exclusiveGateways.flatMap(exclusiveGateway => this.exclusiveGatewayToTransitions(exclusiveGateway, sequenceFlowPlaceIds));
    const parallelGatewayTransitions = choreography.parallelGateways.flatMap(parallelGateway => this.parallelGatewayToTransitions(parallelGateway, sequenceFlowPlaceIds));

    const transitions = [
      startTransition,
      ...endTransitions,
      ...exclusiveGatewayTransitions,
      ...parallelGatewayTransitions,
      ...choreographyTaskTransitions
    ]
    return {
      id: choreography.id,
      placeCount: sequenceFlowPlaceIds.size + additionalPlaceIds.length + 2,
      participantCount: participantIds.size,
      messageCount: messageIds.size,
      startPlace: startTransition.incomingPlaces[0],
      endPlaces: endTransitions.flatMap(endTransition => endTransition.outgoingPlaces),
      transitions,
    };
  }

  private createSequenceFlowMapping(
    sequenceFlows: SequenceFlow[],
  ): Map<SequenceFlowId, PlaceId> {
    const map: Map<SequenceFlowId, PlaceId> = new Map();
    let index = 0;
    for (const sequenceFlow of sequenceFlows.sort()) {
      map.set(sequenceFlow.id, index++);
    }
    return map;
  }

  private createParticipantMapping(
    participants: Participant[],
  ): Map<BpmnParticipantId, ParticipantId> {
    const map: Map<BpmnParticipantId, ParticipantId> = new Map();
    let index = 0;
    for (const participant of participants.sort()) {
      map.set(participant.id, index++);
    }
    return map;
  }

  private createMessageMapping(
    messages: Message[],
  ): Map<BpmnMessageId, MessageId> {
    const map: Map<BpmnMessageId, MessageId> = new Map();
    let index = 0;
    for (const message of messages.sort()) {
      map.set(message.id, index++);
    }
    return map;
  }


  private startEventToTransition(
    startEvent: StartEvent,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
    additionalPlaceIds: PlaceId[]
  ): Transition {
    const outgoingPlaceId = sequenceFlowPlaceIds.get(startEvent.outgoing)!;
    return {
      id: startEvent.id,
      type: TransitionType.START,
      name: startEvent.name,
      incomingPlaces: [sequenceFlowPlaceIds.size + additionalPlaceIds.length],
      outgoingPlaces: [outgoingPlaceId],
    };
  }

  private endEventToTransitions(
    endEvent: EndEvent,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
    additionalPlaceIds: PlaceId[]
  ): Transition[] {
    const incomingPlaceId = sequenceFlowPlaceIds.get(endEvent.incoming)!;
    return [
      {
        id: endEvent.id,
        type: TransitionType.END,
        name: endEvent.name,
        incomingPlaces: [incomingPlaceId],
        outgoingPlaces: [sequenceFlowPlaceIds.size + additionalPlaceIds.length + 1],
      },
    ];
  }

  private parallelGatewayToTransitions(
    parallelGateway: ParallelGateway,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const incomingPlaceIds = parallelGateway.incoming.map(
      (sequenceFlowId) => sequenceFlowPlaceIds.get(sequenceFlowId)!,
    );
    const outgoingPlaceIds = parallelGateway.outgoing.map(
      (sequenceFlowId) => sequenceFlowPlaceIds.get(sequenceFlowId)!,
    );
    return [
      {
        id: parallelGateway.id,
        type: parallelGateway.type === GatewayType.JOIN ? TransitionType.AND_JOIN : TransitionType.AND_SPLIT,
        incomingPlaces: incomingPlaceIds,
        outgoingPlaces: outgoingPlaceIds,
      },
    ];
  }

  private exclusiveGatewayToTransitions(exclusiveGateway: ExclusiveGateway,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>): Transition[] {
    if (exclusiveGateway.type === GatewayType.JOIN) {
      return this.exclusiveJoinGatewayToTransitions(exclusiveGateway, sequenceFlowPlaceIds);
    } else {
      return this.exclusiveSplitGatewayToTransitions(exclusiveGateway, sequenceFlowPlaceIds);
    }
  }

  private exclusiveJoinGatewayToTransitions(
    exclusiveGateway: ExclusiveGateway,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const outgoingPlaceId = sequenceFlowPlaceIds.get(exclusiveGateway.outgoing[0])!;

    return exclusiveGateway.incoming.map((incomingSequenceFlowId) => ({
      id: `${exclusiveGateway.id}_${incomingSequenceFlowId}`,
      type: TransitionType.XOR_JOIN,
      incomingPlaces: [sequenceFlowPlaceIds.get(incomingSequenceFlowId)!],
      outgoingPlaces: [outgoingPlaceId],
    }));
  }

  private exclusiveSplitGatewayToTransitions(
    exclusiveGateway: ExclusiveGateway,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const incomingPlaceId = sequenceFlowPlaceIds.get(exclusiveGateway.incoming[0])!;

    return exclusiveGateway.outgoing.map((outgoingSequenceFlowId) => ({
      id: `${exclusiveGateway.id}_${outgoingSequenceFlowId}`,
      type: TransitionType.XOR_SPLIT,
      incomingPlaces: [incomingPlaceId],
      outgoingPlaces: [sequenceFlowPlaceIds.get(outgoingSequenceFlowId)!],
    }));
  }


  private choreographyTaskToTransitions(
    choreographyTask: ChoreographyTask,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
    participantIds: Map<BpmnParticipantId, ParticipantId>,
    messageIds: Map<BpmnMessageId, MessageId>,
    additionalPlaceIds: PlaceId[]
  ): Transition[] {
    const incomingPlaceId = sequenceFlowPlaceIds.get(choreographyTask.incoming)!;
    const outgoingPlaceId = sequenceFlowPlaceIds.get(choreographyTask.outgoing)!;

    const additionalPlaceId = sequenceFlowPlaceIds.size + additionalPlaceIds.length;
    additionalPlaceIds.push(additionalPlaceId);

    const initialMessage = choreographyTask.initialMessage ? messageIds.get(choreographyTask.initialMessage) : undefined;
    const responseMessage = choreographyTask.responseMessage ? messageIds.get(choreographyTask.responseMessage) : undefined;

    return [
      {
        id: `${choreographyTask.id}_${choreographyTask.initiatingParticipant}`,
        type: TransitionType.TASK,
        name: choreographyTask.name,
        incomingPlaces: [incomingPlaceId],
        outgoingPlaces: [additionalPlaceId],
        participant: participantIds.get(choreographyTask.initiatingParticipant),
        message: initialMessage,
      },
      {
        id: `${choreographyTask.id}_${choreographyTask.respondingParticipant}`,
        type: TransitionType.TASK,
        name: choreographyTask.name,
        incomingPlaces: [additionalPlaceId],
        outgoingPlaces: [outgoingPlaceId],
        participant: participantIds.get(choreographyTask.respondingParticipant),
        message: responseMessage,
      },
    ];
  }

}
