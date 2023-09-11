import { Injectable } from '@nestjs/common';
import {
  Model,
  PlaceId,
  Transition,
  TransitionType,
  copyModel,
} from '../domain/model';

@Injectable()
export class ModelReducer {
  reduceModel(model: Model): Model {
    const newModel = copyModel(model);
    for (const transition of newModel.transitions) {
      switch (transition.type) {
        case TransitionType.START:
        case TransitionType.END:
        case TransitionType.TASK:
          break;
        case TransitionType.XOR_SPLIT:
        case TransitionType.AND_JOIN:
          this.removeTransitionAndOutgoingPlaces(newModel, transition);
          break;
        case TransitionType.XOR_JOIN:
        case TransitionType.AND_SPLIT:
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
    model.startPlace = placeMap.get(model.startPlace)!;
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
      const intersect = this.intersect(
        transition.incomingPlaces,
        transitionToRemove.outgoingPlaces,
      );
      if (intersect.length > 0) {
        transition.incomingPlaces = this.setMinus(
          transition.incomingPlaces,
          transitionToRemove.outgoingPlaces,
        );
        transition.incomingPlaces = this.union(
          transition.incomingPlaces,
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
    for (const transition of model.transitions.values()) {
      const intersect = this.intersect(
        transition.outgoingPlaces,
        transitionToRemove.incomingPlaces,
      );
      if (intersect.length > 0) {
        transition.outgoingPlaces = this.setMinus(
          transition.outgoingPlaces,
          transitionToRemove.incomingPlaces,
        );
        transition.outgoingPlaces = this.union(
          transition.outgoingPlaces,
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
}
