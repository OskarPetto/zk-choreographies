import { Injectable } from '@nestjs/common';
import { PetriNet, PlaceId, Transition, TransitionType } from './petri-net';

@Injectable()
export class PetriNetReducer {
  reducePetriNet(petriNet: PetriNet): PetriNet {
    const newPetriNet = this.copyPetriNet(petriNet);
    for (const transition of newPetriNet.transitions) {
      switch (transition.type) {
        case TransitionType.START:
        case TransitionType.END:
        case TransitionType.TASK:
          break;
        case TransitionType.XOR_SPLIT:
        case TransitionType.AND_JOIN:
          this.removeTransitionAndOutgoingPlaces(newPetriNet, transition);
          break;
        case TransitionType.XOR_JOIN:
        case TransitionType.AND_SPLIT:
          this.removeTransitionAndIncomingPlaces(newPetriNet, transition);
          break;
      }
    }
    this.repairPlaceIds(newPetriNet);
    return newPetriNet;
  }

  private repairPlaceIds(petriNet: PetriNet) {
    const places = this.collectPlaces(petriNet);
    const placeMap: Map<PlaceId, PlaceId> = new Map();
    let index = 0;
    for (const place of places) {
      placeMap.set(place, index++);
    }
    for (const transition of petriNet.transitions) {
      transition.incomingPlaces = transition.incomingPlaces.map(
        (place) => placeMap.get(place)!,
      );
      transition.outgoingPlaces = transition.outgoingPlaces.map(
        (place) => placeMap.get(place)!,
      );
    }
    petriNet.placeCount = places.length;
    petriNet.startPlace = placeMap.get(petriNet.startPlace)!;
  }

  private collectPlaces(petriNet: PetriNet): PlaceId[] {
    let places: PlaceId[] = [];
    for (const transition of petriNet.transitions) {
      places = [
        ...places,
        ...transition.incomingPlaces,
        ...transition.outgoingPlaces,
      ];
    }
    return [...new Set(places)].sort((a, b) => a - b);
  }

  private removeTransitionAndOutgoingPlaces(
    petriNet: PetriNet,
    transitionToRemove: Transition,
  ) {
    for (const transition of petriNet.transitions) {
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
    petriNet.transitions = petriNet.transitions.filter(
      (t) => t.id !== transitionToRemove.id,
    );
  }

  private removeTransitionAndIncomingPlaces(
    petriNet: PetriNet,
    transitionToRemove: Transition,
  ) {
    for (const transition of petriNet.transitions.values()) {
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
    petriNet.transitions = petriNet.transitions.filter(
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

  private copyPetriNet(petriNet: PetriNet): PetriNet {
    const transitions: Transition[] = [...petriNet.transitions.values()].map(
      (transition) => ({
        id: transition.id,
        type: transition.type,
        name: transition.name,
        incomingPlaces: [...transition.incomingPlaces],
        outgoingPlaces: [...transition.outgoingPlaces],
      }),
    );

    return {
      id: petriNet.id,
      placeCount: petriNet.placeCount,
      startPlace: petriNet.startPlace,
      transitions,
    };
  }
}
