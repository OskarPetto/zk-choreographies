import { when } from 'jest-when';
import { Test } from '@nestjs/testing';
import { BpmnService } from './bpmn.service';
import { BpmnParser } from './bpmn.parser';
import { BpmnMapper } from './bpmn.mapper';
import { TestdataProvider } from 'test/data/testdata.provider';
import { ModelReducer } from '../model/model.reducer';
import { ModelService } from '../model/model.service';

describe('CatsController', () => {
  // let bpmnService: BpmnService;
  // let bpmnParser: BpmnParser;
  // let bpmnMapper: BpmnMapper;
  // let modelReducer: ModelReducer;
  // let modelService: ModelService;
  // const bpmnString = TestdataProvider.getExampleChoreography();
  // const definitions = {
  //   process: TestdataProvider.getChoreography1(),
  // };
  // const model1 = TestdataProvider.getModel1();
  const model2 = TestdataProvider.getModel2();

  beforeEach(async () => {
    // const moduleRef = await Test.createTestingModule({
    //   controllers: [],
    //   providers: [
    //     BpmnService,
    //     BpmnParser,
    //     BpmnMapper,
    //     ModelReducer,
    //     ModelService,
    //   ],
    // }).compile();

    // bpmnService = await moduleRef.resolve(BpmnService);
    // bpmnParser = await moduleRef.resolve(BpmnParser);
    // bpmnMapper = await moduleRef.resolve(BpmnMapper);
    // modelReducer = await moduleRef.resolve(ModelReducer);
    // modelService = await moduleRef.resolve(ModelService);
  });

  describe('importBpmnProcess', () => {
    it('should call parser, mapper, reducer and service correctly', async () => {
      // when(jest.spyOn(bpmnParser, 'parseBpmn'))
      //   .calledWith(bpmnString)
      //   .mockReturnValue(definitions);
      // when(jest.spyOn(bpmnMapper, 'toModel'))
      //   .calledWith(definitions.process)
      //   .mockReturnValue(model1);
      // when(jest.spyOn(modelReducer, 'reduceModel'))
      //   .calledWith(model1)
      //   .mockReturnValue(model2);

      // jest.spyOn(modelService, 'saveModel');

      // bpmnService.importBpmn(bpmnString);

      // expect(modelService.saveModel).toBeCalledWith(model2);
    });
  });
});
