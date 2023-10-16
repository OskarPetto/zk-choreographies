import { Injectable } from '@nestjs/common';
import { Model, PlaceId, Transition, TransitionType } from './model';

@Injectable()
export class ModelReducer {
  reduceModel(model: Model): Model {
    const newModel = this.copyModel(model);
    for (const transition of newModel.transitions) {
      switch (transition.type) {
        case TransitionType.REQUIRED:
          break;
        case TransitionType.OPTIONAL_OUTGOING:
          this.removeTransitionAndOutgoingPlaces(newModel, transition);
          break;
        case TransitionType.OPTIONAL_INCOMING:
          this.removeTransitionAndIncomingPlaces(newModel, transition);
          break;
      }
    }
    this.repairPlaceIds(newModel);
    return newModel;
  }

  private repairPlaceIds(model: Model) {
    const places = this.collectPlaces(model);
    const placeMap: Map<PlaceId, PlaceId> = new Map();
    let index = 0;
    for (const place of places) {
      placeMap.set(place, index++);
    }
    for (const transition of model.transitions) {
      transition.incomingPlaces = transition.incomingPlaces.map(
        (place) => placeMap.get(place)!,
      );
      transition.outgoingPlaces = transition.outgoingPlaces.map(
        (place) => placeMap.get(place)!,
      );
    }
    model.placeCount = places.length;
    model.startPlaces = model.startPlaces.map(
      (startPlace) => placeMap.get(startPlace)!,
    );
    model.endPlaces = model.endPlaces.map(
      (endPlace) => placeMap.get(endPlace)!,
    );
  }

  private collectPlaces(model: Model): PlaceId[] {
    let places: PlaceId[] = [];
    for (const transition of model.transitions) {
      places = [
        ...places,
        ...transition.incomingPlaces,
        ...transition.outgoingPlaces,
      ];
    }
    return [...new Set(places)].sort((a, b) => a - b);
  }

  private removeTransitionAndOutgoingPlaces(
    model: Model,
    transitionToRemove: Transition,
  ) {
    for (const transition of model.transitions) {
      if (transition.id === transitionToRemove.id) {
        continue;
      }
      const incomingIntersect = this.intersect(
        transition.incomingPlaces,
        transitionToRemove.outgoingPlaces,
      );
      if (incomingIntersect.length > 0) {
        transition.incomingPlaces = this.setMinus(
          transition.incomingPlaces,
          transitionToRemove.outgoingPlaces,
        );
        transition.incomingPlaces = this.union(
          transition.incomingPlaces,
          transitionToRemove.incomingPlaces,
        );
        if (transitionToRemove.constraint) {
          if (transition.constraint) {
            throw Error(
              `Cannot reduce model because transition ${transition.id} already has constraint`,
            );
          } else {
            transition.constraint = transitionToRemove.constraint;
          }
        }
      }
      const outgoingIntersect = this.intersect(
        transition.outgoingPlaces,
        transitionToRemove.outgoingPlaces,
      );
      if (outgoingIntersect.length > 0) {
        transition.outgoingPlaces = this.setMinus(
          transition.outgoingPlaces,
          transitionToRemove.outgoingPlaces,
        );
        transition.outgoingPlaces = this.union(
          transition.outgoingPlaces,
          transitionToRemove.incomingPlaces,
        );
      }
    }
    model.transitions = model.transitions.filter(
      (t) => t.id !== transitionToRemove.id,
    );
  }

  private removeTransitionAndIncomingPlaces(
    model: Model,
    transitionToRemove: Transition,
  ) {
    for (const transition of model.transitions) {
      if (transition.id === transitionToRemove.id) {
        continue;
      }
      const outgoingIntersect = this.intersect(
        transition.outgoingPlaces,
        transitionToRemove.incomingPlaces,
      );
      if (outgoingIntersect.length > 0) {
        transition.outgoingPlaces = this.setMinus(
          transition.outgoingPlaces,
          transitionToRemove.incomingPlaces,
        );
        transition.outgoingPlaces = this.union(
          transition.outgoingPlaces,
          transitionToRemove.outgoingPlaces,
        );
      }
      const incomingIntersect = this.intersect(
        transition.incomingPlaces,
        transitionToRemove.incomingPlaces,
      );
      if (incomingIntersect.length > 0) {
        transition.incomingPlaces = this.setMinus(
          transition.incomingPlaces,
          transitionToRemove.incomingPlaces,
        );
        transition.incomingPlaces = this.union(
          transition.incomingPlaces,
          transitionToRemove.outgoingPlaces,
        );
      }
    }

    model.transitions = model.transitions.filter(
      (t) => t.id !== transitionToRemove.id,
    );
  }

  private setMinus(places1: PlaceId[], places2: PlaceId[]): PlaceId[] {
    return places1.filter((place1) => !places2.includes(place1));
  }

  private intersect(places1: PlaceId[], places2: PlaceId[]): PlaceId[] {
    return places1.filter((place1) => places2.includes(place1));
  }

  private union(places1: PlaceId[], places2: PlaceId[]): PlaceId[] {
    return [...new Set([...places1, ...places2])];
  }

  private copyModel(model: Model): Model {
    const transitions: Transition[] = [...model.transitions.values()].map(
      (transition) => ({
        id: transition.id,
        type: transition.type,
        name: transition.name,
        incomingPlaces: [...transition.incomingPlaces],
        outgoingPlaces: [...transition.outgoingPlaces],
        sender: transition.sender,
        recipient: transition.recipient,
        message: transition.message,
        constraint: transition.constraint
          ? { ...transition.constraint }
          : undefined,
      }),
    );

    return {
      source: model.source,
      placeCount: model.placeCount,
      participantCount: model.participantCount,
      messageCount: model.messageCount,
      startPlaces: [...model.startPlaces],
      endPlaces: [...model.endPlaces],
      transitions,
    };
  }
}
