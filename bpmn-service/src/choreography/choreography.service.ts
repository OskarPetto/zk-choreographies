import { Injectable } from '@nestjs/common';
import { ChoreographyParser } from './choreography.parser';
import { ChoreographyMapper } from './choreography.mapper';
import { ModelReducer } from '../model/model.reducer';
import { ExecutionGateway } from 'src/execution/execution.gateway';
import { firstValueFrom } from 'rxjs';
import { SaltedHash } from 'src/domain/execution';
import { Choreography } from 'src/domain/choreography';
import { AxiosResponse } from 'axios';

@Injectable()
export class ChoreographyService {
  choreographies: Map<string, Choreography>

  constructor(
    private choreographyParser: ChoreographyParser,
    private choreographyMapper: ChoreographyMapper,
    private modelReducer: ModelReducer,
    private executionGateway: ExecutionGateway
  ) {
    this.choreographies = new Map();
  }

  async transformChoreography(xmlString: string): Promise<SaltedHash> {
    const definitions = this.choreographyParser.parseBpmn(xmlString);
    const choreography = definitions.choreographies[0];
    const model = this.choreographyMapper.toModel(xmlString, choreography);
    const reducedModel = this.modelReducer.reduceModel(model);
    const saltedHash = await this.executionGateway.createModel(reducedModel);
    choreography.modelId = saltedHash.hash;
    this.choreographies.set(choreography.id, choreography);
    return saltedHash;
  }
}
