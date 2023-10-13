import { Injectable } from '@nestjs/common';
import { ChoreographyParser } from './choreography.parser';
import { ChoreographyMapper } from './choreography.mapper';
import { ModelReducer } from '../model/model.reducer';
import { Model } from 'src/model/model';

@Injectable()
export class ChoreographyService {
  constructor(
    private choreographyParser: ChoreographyParser,
    private choreographyMapper: ChoreographyMapper,
    private modelReducer: ModelReducer,
  ) {}

  transformChoreography(xmlString: string): Model {
    const definitions = this.choreographyParser.parseBpmn(xmlString);
    const choreography = definitions.choreographies[0];
    const model = this.choreographyMapper.toModel(xmlString, choreography);
    const reducedModel = this.modelReducer.reduceModel(model);
    return reducedModel;
  }
}
