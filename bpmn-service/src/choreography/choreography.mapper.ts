import { Injectable } from '@nestjs/common';
import {
  ChoreographyTask,
  EndEvent,
  ExclusiveGateway,
  GatewayType,
  ParallelGateway,
  Participant,
  SequenceFlow,
  SequenceFlowId,
  StartEvent,
  ParticipantId as BpmnParticipantId,
  MessageId as BpmnMessageId,
  Message,
  Choreography,
  LoopType,
} from './choreography';
import {
  Model,
  ParticipantId,
  MessageId,
  PlaceId,
  Transition,
  TransitionType,
} from '../model/model';
import { ConstraintParser } from 'src/constraint/constraint.parser';

interface ConstraintMapping {
  sequenceFlowNames: Map<PlaceId, string | undefined>;
  messageIdPerMessageName: Map<string, MessageId>;
}

@Injectable()
export class ChoreographyMapper {
  constructor(private constraintParser: ConstraintParser) {}
  toModel(xmlString: string, choreography: Choreography): Model {
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
          additionalPlaceIds,
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
            choreographyTask.initiatingParticipant === participantId,
        ),
    );

    const constraintMapping = this.createConstraintMapping(
      choreography.sequenceFlows,
      sequenceFlowPlaceIds,
      choreography.messages,
      messageIds,
    );
    this.addConstraints(transitions, constraintMapping);

    return {
      source: xmlString,
      placeCount: sequenceFlowPlaceIds.size + additionalPlaceIds.length + 2,
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

  private createConstraintMapping(
    sequenceFlows: SequenceFlow[],
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
    messages: Message[],
    messageIds: Map<BpmnMessageId, MessageId>,
  ): ConstraintMapping {
    const constraintMapping: ConstraintMapping = {
      sequenceFlowNames: new Map(),
      messageIdPerMessageName: new Map(),
    };
    for (const sequenceFlow of sequenceFlows) {
      const placeId = sequenceFlowPlaceIds.get(sequenceFlow.id)!;
      constraintMapping.sequenceFlowNames.set(placeId, sequenceFlow.name);
    }
    for (const message of messages) {
      const messageId = messageIds.get(message.id)!;
      constraintMapping.messageIdPerMessageName.set(message.name, messageId);
    }
    return constraintMapping;
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
    additionalPlaceIds: PlaceId[],
  ): Transition {
    const outgoingPlaceId = sequenceFlowPlaceIds.get(startEvent.outgoing)!;
    return {
      id: startEvent.id,
      type: TransitionType.REQUIRED,
      name: startEvent.name,
      incomingPlaces: [sequenceFlowPlaceIds.size + additionalPlaceIds.length],
      outgoingPlaces: [outgoingPlaceId],
    };
  }

  private endEventToTransition(
    endEvent: EndEvent,
    sequenceFlowPlaceIds: Map<SequenceFlowId, PlaceId>,
    additionalPlaceIds: PlaceId[],
  ): Transition[] {
    const incomingPlaceId = sequenceFlowPlaceIds.get(endEvent.incoming)!;
    return [
      {
        id: endEvent.id,
        type: TransitionType.REQUIRED,
        name: endEvent.name,
        incomingPlaces: [incomingPlaceId],
        outgoingPlaces: [
          sequenceFlowPlaceIds.size + additionalPlaceIds.length + 1,
        ],
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
    additionalPlaceIds: PlaceId[],
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

    const initialMessage = choreographyTask.initialMessage
      ? messageIds.get(choreographyTask.initialMessage)
      : undefined;
    const responseMessage = choreographyTask.responseMessage
      ? messageIds.get(choreographyTask.responseMessage)
      : undefined;

    if (choreographyTask.loopType === undefined) {
      const additionalPlaceId =
        sequenceFlowPlaceIds.size + additionalPlaceIds.length;
      additionalPlaceIds.push(additionalPlaceId);

      return [
        {
          id: `${choreographyTask.id}_${choreographyTask.initiatingParticipant}`,
          type: TransitionType.REQUIRED,
          name: choreographyTask.name,
          incomingPlaces: [incomingPlaceId],
          outgoingPlaces: [additionalPlaceId],
          initiatingParticipant: initiatingParticipantId,
          respondingParticipant: respondingParticipantId,
          message: initialMessage,
        },
        {
          id: `${choreographyTask.id}_${choreographyTask.respondingParticipant}`,
          type: TransitionType.REQUIRED,
          name: choreographyTask.name,
          incomingPlaces: [additionalPlaceId],
          outgoingPlaces: [outgoingPlaceId],
          initiatingParticipant: respondingParticipantId,
          respondingParticipant: initiatingParticipantId,
          message: responseMessage,
        },
      ];
    } else if (
      choreographyTask.loopType === LoopType.MULTI_INSTANCE_SEQUENTIAL
    ) {
      return [
        {
          id: `${choreographyTask.id}_0`,
          type: TransitionType.REQUIRED,
          name: choreographyTask.name,
          incomingPlaces: [incomingPlaceId],
          outgoingPlaces: [outgoingPlaceId],
          initiatingParticipant: initiatingParticipantId,
          respondingParticipant: respondingParticipantId,
          message: initialMessage,
        },
        {
          id: `${choreographyTask.id}_loop`,
          type: TransitionType.REQUIRED,
          name: choreographyTask.name,
          incomingPlaces: [outgoingPlaceId],
          outgoingPlaces: [outgoingPlaceId],
          initiatingParticipant: initiatingParticipantId,
          respondingParticipant: respondingParticipantId,
          message: initialMessage,
        },
      ];
    } else {
      throw Error('not supported');
    }
  }

  addConstraints(
    transitions: Transition[],
    constraintMapping: ConstraintMapping,
  ) {
    transitions.forEach((transition) => {
      for (const incomingPlace of transition.incomingPlaces) {
        const constraintString =
          constraintMapping.sequenceFlowNames.get(incomingPlace);
        if (constraintString !== undefined) {
          const constraint = this.constraintParser.parseConstraint(
            constraintString,
            constraintMapping.messageIdPerMessageName,
          );
          transition.messageConstraint = constraint;
          break;
        }
      }
    });
  }
}
