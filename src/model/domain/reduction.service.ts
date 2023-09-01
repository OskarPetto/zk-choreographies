import { Injectable } from '@nestjs/common';
import { Model, PlaceId, Transition } from './model';

@Injectable()
export class ReductionService {

    reduceModel(model: Model): Model {


        return model;
    }

    private removeTransition(model: Model, transitionToRemove: Transition): Model {
        for (const transition of model.transitions.values()) {
            if (transition.id === transitionToRemove.id) {
                continue;
            }

            transition.toPlaces = this.setMinus(transition.toPlaces, transitionToRemove.fromPlaces);
            transition.toPlaces = this.union(transition.toPlaces, transitionToRemove.toPlaces);
        }
    }

    private setMinus(places1: Set<PlaceId>, places2: Set<PlaceId>): Set<PlaceId> {
        const difference: Set<PlaceId> = new Set();

        for (const place of places1) {
            if (!places2.has(place)) {
                difference.add(place);
            }
        }

        return difference;
    }

    private union(places1: Set<PlaceId>, places2: Set<PlaceId>): Set<PlaceId> {
        return new Set([...places1, ...places2]);
    }
}
