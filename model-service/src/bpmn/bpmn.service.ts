import { Injectable } from '@nestjs/common';
import { BpmnParser } from './bpmn.parser';
import { BpmnMapper } from './bpmn.mapper';
import { PetriNetReducer } from '../petri-net/perti-net.reducer';
import { PetriNetService } from '../petri-net/petri-net.service';

@Injectable()
export class BpmnService {
  constructor(
    private bpmnParser: BpmnParser,
    private bpmnMapper: BpmnMapper,
    private petriNetReducer: PetriNetReducer,
    private petriNetService: PetriNetService,
  ) { }

  importBpmn(bpmnString: string) {
    const definitions = this.bpmnParser.parseBpmn(bpmnString);
    const petriNet = this.bpmnMapper.toPetriNet(definitions.process);
    const reducedPetriNet = this.petriNetReducer.reducePetriNet(petriNet);
    this.petriNetService.savePetriNet(reducedPetriNet);
  }
}
