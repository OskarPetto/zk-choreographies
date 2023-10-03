import { Injectable } from '@nestjs/common';
import { ChoreographyParser } from './choreography.parser';
import { ChoreographyMapper } from './choreography.mapper';
import { ModelReducer } from '../model/model.reducer';
import { ModelGateway } from '../model/model.gateway';
import { Choreography } from './choreography';

@Injectable()
export class ChoreographyService {
  private choreographies: Map<string, Choreography> = new Map();

  constructor(
    private bpmnParser: ChoreographyParser,
    private bpmnMapper: ChoreographyMapper,
    private modelReducer: ModelReducer,
    private modelGateway: ModelGateway,
  ) {}

  async importChoreography(bpmnString: string): Promise<string> {
    const definitions = this.bpmnParser.parseBpmn(bpmnString);
    const choreography = definitions.choreographies[0];
    const model = this.bpmnMapper.toModel(choreography);
    const reducedModel = this.modelReducer.reduceModel(model);
    const createdModel = await this.modelGateway.createModel(reducedModel);
    this.choreographies.set(choreography.id, choreography);
    return createdModel.hash.value;
  }
}
