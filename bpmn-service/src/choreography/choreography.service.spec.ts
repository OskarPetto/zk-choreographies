import { when } from 'jest-when';
import { Test } from '@nestjs/testing';
import { ChoreographyService } from './choreography.service';
import { ChoreographyParser } from './choreography.parser';
import { ChoreographyMapper } from './choreography.mapper';
import { TestdataProvider } from '../../test/data/testdata.provider';
import { ModelReducer } from '../model/model.reducer';
import { HttpModule } from '@nestjs/axios';
import { ConfigModule } from '@nestjs/config';
import { ConstraintModule } from 'src/constraint/constraint.module';

describe('ChoreographyService', () => {
  let choreographyService: ChoreographyService;
  let choreographyParser: ChoreographyParser;
  let choreographyMapper: ChoreographyMapper;
  let modelReducer: ModelReducer;
  const xmlString = TestdataProvider.getExampleChoreography();
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
      ],
      imports: [HttpModule, ConfigModule, ConstraintModule],
    }).compile();

    choreographyService =
      moduleRef.get<ChoreographyService>(ChoreographyService);
    choreographyParser = moduleRef.get<ChoreographyParser>(ChoreographyParser);
    choreographyMapper = moduleRef.get<ChoreographyMapper>(ChoreographyMapper);
    modelReducer = moduleRef.get<ModelReducer>(ModelReducer);
  });

  describe('importBpmnProcess', () => {
    it('should call parser, mapper and reducer correctly', async () => {
      when(jest.spyOn(choreographyParser, 'parseBpmn'))
        .calledWith(xmlString)
        .mockReturnValue(definitions);
      when(jest.spyOn(choreographyMapper, 'toModel'))
        .calledWith(xmlString, definitions.choreographies[0])
        .mockReturnValue(model2);
      when(jest.spyOn(modelReducer, 'reduceModel'))
        .calledWith(model2)
        .mockReturnValue(model2Reduced);

      const result = choreographyService.transformChoreography(xmlString);
      expect(result).toEqual(model2Reduced);
    });
  });
});
