import { when } from 'jest-when';
import { Test } from '@nestjs/testing';
import { ChoreographyService } from './choreography.service';
import { ChoreographyParser } from './choreography.parser';
import { ChoreographyMapper } from './choreography.mapper';
import { TestdataProvider } from '../../test/data/testdata.provider';
import { ModelReducer } from '../model/model.reducer';
import { ModelGateway } from '../model/model.gateway';
import { HttpModule } from '@nestjs/axios';
import { ConfigModule } from '@nestjs/config';

describe('ChoreographyService', () => {
  let bpmnService: ChoreographyService;
  let bpmnParser: ChoreographyParser;
  let bpmnMapper: ChoreographyMapper;
  let modelReducer: ModelReducer;
  let modelGateway: ModelGateway;
  const bpmnString = TestdataProvider.getExampleChoreography();
  const definitions = TestdataProvider.getDefinitions2();
  const model2Reduced = TestdataProvider.getModel2Reduced();
  const model2 = TestdataProvider.getModel2();

  beforeEach(async () => {
    const moduleRef = await Test.createTestingModule({
      controllers: [],
      providers: [
        ChoreographyService,
        ChoreographyParser,
        ChoreographyMapper,
        ModelReducer,
        ModelGateway,
      ],
      imports: [HttpModule, ConfigModule]
    }).compile();

    bpmnService = moduleRef.get<ChoreographyService>(ChoreographyService);
    bpmnParser = moduleRef.get<ChoreographyParser>(ChoreographyParser);
    bpmnMapper = moduleRef.get<ChoreographyMapper>(ChoreographyMapper);
    modelReducer = moduleRef.get<ModelReducer>(ModelReducer);
    modelGateway = moduleRef.get<ModelGateway>(ModelGateway);
  });

  describe('importBpmnProcess', () => {
    it('should call parser, mapper, reducer and service correctly', async () => {
      when(jest.spyOn(bpmnParser, 'parseBpmn'))
        .calledWith(bpmnString)
        .mockReturnValue(definitions);
      when(jest.spyOn(bpmnMapper, 'toModel'))
        .calledWith(definitions.choreographies[0])
        .mockReturnValue(model2);
      when(jest.spyOn(modelReducer, 'reduceModel'))
        .calledWith(model2)
        .mockReturnValue(model2Reduced);

      jest.spyOn(modelGateway, 'createModel');

      bpmnService.importChoreography(bpmnString);

      expect(modelGateway.createModel).toBeCalledWith(model2Reduced);
    });
  });
});
