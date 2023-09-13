import { Injectable } from '@nestjs/common';
import { PetriNetId, PetriNet } from './petri-net';
import { v4 as uuid } from 'uuid';

@Injectable()
export class PetriNetService {
  petriNets: Map<PetriNetId, PetriNet> = new Map();

  savePetriNet(petriNet: PetriNet) {
    if (!petriNet.id) {
      petriNet.id = this.createPetriNetId();
    }
    this.petriNets.set(petriNet.id, petriNet);
  }

  findPetriNet(petriNetId: PetriNetId): PetriNet {
    const petriNet = this.petriNets.get(petriNetId);
    if (!petriNet) {
      throw Error(`PetriNet ${petriNetId} not found`);
    }
    return petriNet;
  }

  private createPetriNetId(): PetriNetId {
    return uuid();
  }
}
