import { Injectable } from '@nestjs/common';
import { Instance } from './instance';
import { Model, PlaceId, Transition } from 'src/model';

@Injectable()
export class ConformanceService {
  isExecutionValid(
    instanceBefore: Instance,
    instanceAfter: Instance,
    transition: Transition,
    model: Model,
  ): boolean {
    if (instanceBefore.id !== instanceAfter.id) {
      return false;
    }
    if (!this.isInstanceValid(instanceBefore, model)) {
      return false;
    }
    if (!this.isInstanceValid(instanceAfter, model)) {
      return false;
    }
    if (!this.isTransitionInModel(transition, model)) {
      return false;
    }
    return this.areTokenChangesValid(instanceBefore, instanceAfter, transition);
  }

  private areTokenChangesValid(
    instanceBefore: Instance,
    instanceAfter: Instance,
    transition: Transition,
  ): boolean {
    for (let placeId = 0; placeId < instanceBefore.tokenCounts.length; placeId++) {
      const tokenCountBefore = instanceBefore.tokenCounts[placeId];
      const tokenCountAfter = instanceAfter.tokenCounts[placeId];
      if (transition.incomingPlaces.includes(placeId)) {
        if (tokenCountAfter !== tokenCountBefore - 1) {
          return false;
        }
      } else if (transition.outgoingPlaces.includes(placeId)) {
        if (tokenCountAfter !== tokenCountBefore + 1) {
          return false;
        }
      } else {
        if (tokenCountAfter !== tokenCountBefore) {
          return false;
        }
      }
    }
    return true;
  }

  private isTransitionInModel(transition: Transition, model: Model): boolean {
    const transitionInModel = model.transitions.find(
      (t) => t.id === transition.id,
    );
    if (!transitionInModel) {
      return false;
    }
    return (
      transition.type === transitionInModel.type &&
      transition.name === transitionInModel.name &&
      this.placesEqual(
        transition.incomingPlaces,
        transitionInModel.incomingPlaces,
      ) &&
      this.placesEqual(
        transition.outgoingPlaces,
        transitionInModel.outgoingPlaces,
      )
    );
  }

  private placesEqual(places1: PlaceId[], places2: PlaceId[]): boolean {
    return (
      places1.length === places2.length &&
      places1.every((place1, index) => place1 === places2[index])
    );
  }
  private isInstanceValid(instance: Instance, model: Model): boolean {
    if (!this.isInstanceOfModel(instance, model)) {
      return false;
    }
    return instance.tokenCounts.every((tokenCount) => tokenCount >= 0);
  }

  private isInstanceOfModel(instance: Instance, model: Model): boolean {
    return (
      instance.model === model.id &&
      instance.tokenCounts.length === model.placeCount
    );
  }
}
