import { Injectable } from '@nestjs/common';
import {
  Model,
  ParticipantId,
  MessageId,
  PlaceId,
  Transition,
  TransitionType,
} from '../domain/model';
import { ConditionParser } from 'src/condition/condition.parser';
import {
  ParsedChoreography,
  SequenceFlow,
  SequenceFlowId,
  Message,
  Participant,
  StartEvent,
  EndEvent,
  ParallelGateway,
  GatewayType,
  ExclusiveGateway,
  ChoreographyTask,
} from './choreography.parser';
import { BpmnMessageId, BpmnParticipantId } from 'src/domain/choreography';

interface ConditionMapping {
  sequenceFlowNames: Map<PlaceId, string | undefined>;
  messageIdPerMessageName: Map<string, MessageId>;
}

@Injectable()
export class ChoreographyMapper {
  constructor(private conditionParser: ConditionParser) { }
  toModel(choreography: ParsedChoreography): Model {
    const sequenceFlowPlaceIds = this.createSequenceFlowMapping(
      choreography.sequenceFlows,
    );
    const participantIds = this.createParticipantMapping(
      choreography.participants,
    );
    const messageIds = this.createMessageMapping(choreography.messages);

    const additionalPlaceIds: PlaceId[] = [];
    const choreographyTaskTransitions = choreography.choreographyTasks.flatMap(
      (choreographyTask) =>
        this.choreographyTaskToTransitions(
          choreographyTask,
          sequenceFlowPlaceIds,
          participantIds,
          messageIds,
        ),
    );

    const startTransitions = choreography.startEvents.flatMap((startEvent) =>
      this.startEventToTransition(
        startEvent,
        sequenceFlowPlaceIds,
        additionalPlaceIds,
      ),
    );
    const endTransitions = choreography.endEvents.flatMap((endEvent) =>
      this.endEventToTransition(
        endEvent,
        sequenceFlowPlaceIds,
        additionalPlaceIds,
      ),
    );
    const exclusiveGatewayTransitions = choreography.exclusiveGateways.flatMap(
      (exclusiveGateway) =>
        this.exclusiveGatewayToTransitions(
          exclusiveGateway,
          sequenceFlowPlaceIds,
        ),
    );
    const parallelGatewayTransitions = choreography.parallelGateways.flatMap(
      (parallelGateway) =>
        this.parallelGatewayToTransitions(
          parallelGateway,
          sequenceFlowPlaceIds,
        ),
    );

    this.fixTransitionsOfExclusiveGatewayPairs(
      exclusiveGatewayTransitions,
      choreography.exclusiveGatewayPairs,
    );

    const transitions = [
      ...startTransitions,
      ...endTransitions,
      ...exclusiveGatewayTransitions,
      ...parallelGatewayTransitions,
      ...choreographyTaskTransitions,
    ];
    const relevantParticipants = [...participantIds.values()].filter(
      (participantId) =>
        choreographyTaskTransitions.some(
          (choreographyTask) =>
            choreographyTask.initiatingParticipant === participantId ||
            choreographyTask.respondingParticipant === participantId,
        ),
    );

    const conditionMapping = this.createConditionMapping(
      choreography.sequenceFlows,
      sequenceFlowPlaceIds,
      choreography.messages,
      messageIds,
    );
    this.addConditions(transitions, conditionMapping);
    return {
      placeCount: sequenceFlowPlaceIds.size + additionalPlaceIds.length,
      participantCount: relevantParticipants.length,
      messageCount: messageIds.size,
      startPlaces: startTransitions.flatMap(
        (startTransition) => startTransition.incomingPlaces,
      ),
      endPlaces: endTransitions.flatMap(
        (endTransition) => endTransition.outgoingPlaces,
      ),
      transitions,
    };
  }

  fixTransitionsOfExclusiveGatewayPairs(
    exclusiveGatewayTransitions: Transition[],
    exclusiveGatewayPairs: string[][],
  ) {
    exclusiveGatewayPairs.forEach((exclusiveGatewayPair) => {
      const concatinatedId1 = `${exclusiveGatewayPair[0]}_${exclusiveGatewayPair[2]}`;
      const concatinatedId2 = `${exclusiveGatewayPair[1]}_${exclusiveGatewayPair[2]}`;
      const transition = exclusiveGatewayTransitions.find(
        (exclusiveGatewayTransition) => {
          return (
            exclusiveGatewayTransition.id === concatinatedId1 ||
            exclusiveGatewayTransition.id === concatinatedId2
          );
        },
      );
      if (transition != null) {
        transition.type = TransitionType.REQUIRED;
      }
    });
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

  private createConditionMapping(
    sequenceFlows: SequenceFlow[],
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
    messages: Message[],
    messageIds: Map<BpmnMessageId, MessageId>,
  ): ConditionMapping {
    const conditionMapping: ConditionMapping = {
      sequenceFlowNames: new Map(),
      messageIdPerMessageName: new Map(),
    };
    for (const sequenceFlow of sequenceFlows) {
      const placeId = sequenceFlowPlaceIds.get(sequenceFlow.id)!;
      conditionMapping.sequenceFlowNames.set(placeId, sequenceFlow.name);
    }
    for (const message of messages) {
      const messageId = messageIds.get(message.id)!;
      conditionMapping.messageIdPerMessageName.set(message.name, messageId);
    }
    return conditionMapping;
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
      if (message.name !== undefined) {
        map.set(message.id, index++);
      }
    }
    return map;
  }

  private startEventToTransition(
    startEvent: StartEvent,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
    additionalPlaceIds: PlaceId[],
  ): Transition {
    const outgoingPlaceId = sequenceFlowPlaceIds.get(startEvent.outgoing)!;
    const newPlaceId = sequenceFlowPlaceIds.size + additionalPlaceIds.length;
    additionalPlaceIds.push(newPlaceId);
    return {
      id: startEvent.id,
      type: TransitionType.REQUIRED,
      name: startEvent.name,
      incomingPlaces: [newPlaceId],
      outgoingPlaces: [outgoingPlaceId],
    };
  }

  private endEventToTransition(
    endEvent: EndEvent,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
    additionalPlaceIds: PlaceId[],
  ): Transition[] {
    const incomingPlaceId = sequenceFlowPlaceIds.get(endEvent.incoming)!;
    const newPlaceId = sequenceFlowPlaceIds.size + additionalPlaceIds.length;
    additionalPlaceIds.push(newPlaceId);
    return [
      {
        id: endEvent.id,
        type: TransitionType.REQUIRED,
        name: endEvent.name,
        incomingPlaces: [incomingPlaceId],
        outgoingPlaces: [newPlaceId],
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
        type:
          parallelGateway.type === GatewayType.JOIN
            ? TransitionType.OPTIONAL_OUTGOING
            : TransitionType.OPTIONAL_INCOMING,
        incomingPlaces: incomingPlaceIds,
        outgoingPlaces: outgoingPlaceIds,
      },
    ];
  }

  private exclusiveGatewayToTransitions(
    exclusiveGateway: ExclusiveGateway,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    if (exclusiveGateway.type === GatewayType.JOIN) {
      return this.exclusiveJoinGatewayToTransitions(
        exclusiveGateway,
        sequenceFlowPlaceIds,
      );
    } else {
      return this.exclusiveSplitGatewayToTransitions(
        exclusiveGateway,
        sequenceFlowPlaceIds,
      );
    }
  }

  private exclusiveJoinGatewayToTransitions(
    exclusiveGateway: ExclusiveGateway,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const outgoingPlaceId = sequenceFlowPlaceIds.get(
      exclusiveGateway.outgoing[0],
    )!;

    return exclusiveGateway.incoming.map((incomingSequenceFlowId) => ({
      id: `${exclusiveGateway.id}_${incomingSequenceFlowId}`,
      type: TransitionType.OPTIONAL_INCOMING,
      incomingPlaces: [sequenceFlowPlaceIds.get(incomingSequenceFlowId)!],
      outgoingPlaces: [outgoingPlaceId],
    }));
  }

  private exclusiveSplitGatewayToTransitions(
    exclusiveGateway: ExclusiveGateway,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
  ): Transition[] {
    const incomingPlaceId = sequenceFlowPlaceIds.get(
      exclusiveGateway.incoming[0],
    )!;

    return exclusiveGateway.outgoing.map((outgoingSequenceFlowId) => {
      return {
        id: `${exclusiveGateway.id}_${outgoingSequenceFlowId}`,
        type: TransitionType.OPTIONAL_OUTGOING,
        incomingPlaces: [incomingPlaceId],
        outgoingPlaces: [sequenceFlowPlaceIds.get(outgoingSequenceFlowId)!],
      };
    });
  }

  private choreographyTaskToTransitions(
    choreographyTask: ChoreographyTask,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
    participantIds: Map<BpmnParticipantId, ParticipantId>,
    messageIds: Map<BpmnMessageId, MessageId>,
  ): Transition[] {
    const incomingPlaceId = sequenceFlowPlaceIds.get(
      choreographyTask.incoming,
    )!;
    const outgoingPlaceId = sequenceFlowPlaceIds.get(
      choreographyTask.outgoing,
    )!;

    const initiatingParticipantId = participantIds.get(
      choreographyTask.initiatingParticipant,
    );
    const respondingParticipantId = participantIds.get(
      choreographyTask.respondingParticipant,
    );

    const initiatingMessageId = choreographyTask.initialMessage
      ? messageIds.get(choreographyTask.initialMessage)
      : undefined;
    const respondingMessageId = choreographyTask.responseMessage
      ? messageIds.get(choreographyTask.responseMessage)
      : undefined;

    return [
      {
        id: choreographyTask.id,
        type: TransitionType.REQUIRED,
        name: choreographyTask.name,
        incomingPlaces: [incomingPlaceId],
        outgoingPlaces: [outgoingPlaceId],
        initiatingParticipant: initiatingParticipantId,
        respondingParticipant: respondingParticipantId,
        initiatingMessage: initiatingMessageId,
        respondingMessage: respondingMessageId,
      },
    ];
  }

  addConditions(transitions: Transition[], conditionMapping: ConditionMapping) {
    transitions.forEach((transition) => {
      for (const incomingPlace of transition.incomingPlaces) {
        const conditionString =
          conditionMapping.sequenceFlowNames.get(incomingPlace);
        if (conditionString !== undefined) {
          const condition = this.conditionParser.parseCondition(
            conditionString,
            conditionMapping.messageIdPerMessageName,
          );
          transition.condition = condition;
          break;
        }
      }
    });
  }
}
