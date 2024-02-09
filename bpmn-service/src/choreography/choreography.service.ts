import { Injectable } from '@nestjs/common';
import { ChoreographyParser } from './choreography.parser';
import { ChoreographyMapper } from './choreography.mapper';
import { ModelReducer } from '../model/model.reducer';
import { ExecutionGateway } from 'src/execution/execution.gateway';
import { Choreography } from 'src/domain/choreography';

@Injectable()
export class ChoreographyService {
  choreographies: Map<string, Choreography>;

  constructor(
    private choreographyParser: ChoreographyParser,
    private choreographyMapper: ChoreographyMapper,
    private modelReducer: ModelReducer,
    private executionGateway: ExecutionGateway,
  ) {
    this.choreographies = new Map();
  }

  async transformChoreography(xmlString: string): Promise<Choreography> {
    const definitions = this.choreographyParser.parseBpmn(xmlString);
    const parsedChoreography = definitions.choreographies[0];
    const model = this.choreographyMapper.toModel(parsedChoreography);
    const reducedModel = this.modelReducer.reduceModel(model);
    const saltedHash = await this.executionGateway.createModel(reducedModel);
    const choreography = {
      id: saltedHash.hash,
      xmlString: xmlString,
    };
    this.choreographies.set(choreography.id, choreography);
    return choreography;
  }

  findAllChoreographies(): Choreography[] {
    return [...this.choreographies.values()];
  }
}
