import { Injectable } from '@nestjs/common';
import { ChoreographyParser } from './choreography.parser';
import { ChoreographyMapper } from './choreography.mapper';
import { ModelReducer } from '../model/model.reducer';
import { ModelGateway } from '../model/model.gateway';

@Injectable()
export class ChoreographyService {
  constructor(
    private choreographyParser: ChoreographyParser,
    private choreographyMapper: ChoreographyMapper,
    private modelReducer: ModelReducer,
    private modelGateway: ModelGateway,
  ) {}

  async transformChoreography(xmlString: string): Promise<string> {
    const definitions = this.choreographyParser.parseBpmn(xmlString);
    const choreography = definitions.choreographies[0];
    const model = this.choreographyMapper.toModel(xmlString, choreography);
    const reducedModel = this.modelReducer.reduceModel(model);
    const createdModel = await this.modelGateway.createModel(reducedModel);
    return createdModel.hash.value;
  }
}
