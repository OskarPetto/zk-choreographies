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
import { ConstraintModule } from 'src/constraint/constraint.module';
import { BpmnService } from 'src/bpmn/bpmn.service';

describe('ChoreographyService', () => {
  let choreographyService: ChoreographyService;
  let bpmnService: BpmnService;
  let choreographyParser: ChoreographyParser;
  let choreographyMapper: ChoreographyMapper;
  let modelReducer: ModelReducer;
  let modelGateway: ModelGateway;
  const bpmnModel = TestdataProvider.getExampleChoreography();
  const definitions = TestdataProvider.getDefinitions2();
  const model2Reduced = TestdataProvider.getModel2Reduced();
  const model2 = TestdataProvider.getModel2();

  beforeEach(async () => {
    const moduleRef = await Test.createTestingModule({
      controllers: [],
      providers: [
        BpmnService,
        ChoreographyService,
        ChoreographyParser,
        ChoreographyMapper,
        ModelReducer,
        ModelGateway,
      ],
      imports: [HttpModule, ConfigModule, ConstraintModule],
    }).compile();

    choreographyService =
      moduleRef.get<ChoreographyService>(ChoreographyService);
    bpmnService = moduleRef.get<BpmnService>(BpmnService);
    choreographyParser = moduleRef.get<ChoreographyParser>(ChoreographyParser);
    choreographyMapper = moduleRef.get<ChoreographyMapper>(ChoreographyMapper);
    modelReducer = moduleRef.get<ModelReducer>(ModelReducer);
    modelGateway = moduleRef.get<ModelGateway>(ModelGateway);
  });

  describe('importBpmnProcess', () => {
    it('should call parser, mapper, reducer and service correctly', async () => {
      when(jest.spyOn(bpmnService, 'findBpmnModelById'))
        .calledWith(bpmnModel.id)
        .mockReturnValue(bpmnModel);
      when(jest.spyOn(choreographyParser, 'parseBpmn'))
        .calledWith(bpmnModel.xmlString)
        .mockReturnValue(definitions);
      when(jest.spyOn(choreographyMapper, 'toModel'))
        .calledWith(bpmnModel.id, definitions.choreographies[0])
        .mockReturnValue(model2);
      when(jest.spyOn(modelReducer, 'reduceModel'))
        .calledWith(model2)
        .mockReturnValue(model2Reduced);

      jest
        .spyOn(modelGateway, 'createModel')
        .mockImplementation((model) => Promise.resolve(model));

      choreographyService.transformChoreography(bpmnModel.id);

      expect(modelGateway.createModel).toBeCalledWith(model2Reduced);
    });
  });
});
