import { when } from 'jest-when';
import { Test } from '@nestjs/testing';
import { ChoreographyService } from './choreography.service';
import { ChoreographyParser } from './choreography.parser';
import { ChoreographyMapper } from './choreography.mapper';
import { TestdataProvider } from '../../testdata/testdata.provider';
import { ModelReducer } from '../model/model.reducer';
import { HttpModule } from '@nestjs/axios';
import { ConfigModule } from '@nestjs/config';
import { ExecutionGateway } from 'src/execution/execution.gateway';
import { ConditionModule } from 'src/condition/condition.module';

describe('ChoreographyService', () => {
  let choreographyService: ChoreographyService;
  let choreographyParser: ChoreographyParser;
  let choreographyMapper: ChoreographyMapper;
  let modelReducer: ModelReducer;
  let executionGateway: ExecutionGateway;
  const xmlString = TestdataProvider.readBpmn('example_choreography');
  const definitions = TestdataProvider.getDefinitions2();
  const model2Reduced = TestdataProvider.getModel2Reduced();
  const model2 = TestdataProvider.getModel2();
  const saltedHash = TestdataProvider.getSaltedHash();
  beforeEach(async () => {
    const moduleRef = await Test.createTestingModule({
      controllers: [],
      providers: [
        ChoreographyService,
        ChoreographyParser,
        ChoreographyMapper,
        ModelReducer,
        ExecutionGateway,
      ],
      imports: [HttpModule, ConfigModule, ConditionModule],
    }).compile();

    choreographyService =
      moduleRef.get<ChoreographyService>(ChoreographyService);
    choreographyParser = moduleRef.get<ChoreographyParser>(ChoreographyParser);
    choreographyMapper = moduleRef.get<ChoreographyMapper>(ChoreographyMapper);
    modelReducer = moduleRef.get<ModelReducer>(ModelReducer);
    executionGateway = moduleRef.get<ExecutionGateway>(ExecutionGateway);
  });

  describe('importBpmnProcess', () => {
    it('should call parser, mapper and reducer correctly', async () => {
      when(jest.spyOn(choreographyParser, 'parseBpmn'))
        .calledWith(xmlString)
        .mockReturnValue(definitions);
      when(jest.spyOn(choreographyMapper, 'toModel'))
        .calledWith(definitions.choreographies[0])
        .mockReturnValue(model2);
      when(jest.spyOn(modelReducer, 'reduceModel'))
        .calledWith(model2)
        .mockReturnValue(model2Reduced);
      when(jest.spyOn(executionGateway, 'createModel'))
        .calledWith(model2Reduced)
        .mockResolvedValue(Promise.resolve(saltedHash));

      const result = await choreographyService.transformChoreography(xmlString);
      expect(result).toEqual({
        id: saltedHash.hash,
        xmlString: xmlString,
      });
    });
  });
});
