import { Injectable } from '@nestjs/common';
import { ChoreographyParser } from './choreography.parser';
import { ChoreographyMapper } from './choreography.mapper';
import { ModelReducer } from '../model/model.reducer';
import { ModelGateway } from '../model/model.gateway';
import { BpmnService } from 'src/bpmn/bpmn.service';

@Injectable()
export class ChoreographyService {
  constructor(
    private bpmnService: BpmnService,
    private choreographyParser: ChoreographyParser,
    private choreographyMapper: ChoreographyMapper,
    private modelReducer: ModelReducer,
    private modelGateway: ModelGateway,
  ) {}

  async transformChoreography(id: string): Promise<string> {
    const bpmnModel = this.bpmnService.findBpmnModelById(id);
    const definitions = this.choreographyParser.parseBpmn(bpmnModel.xmlString);
    const choreography = definitions.choreographies[0];
    const model = this.choreographyMapper.toModel(id, choreography);
    const reducedModel = this.modelReducer.reduceModel(model);
    const createdModel = await this.modelGateway.createModel(reducedModel);
    return createdModel.hash.value;
  }
}
