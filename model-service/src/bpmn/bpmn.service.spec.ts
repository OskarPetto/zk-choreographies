import { when } from 'jest-when';
import { Test } from '@nestjs/testing';
import { BpmnService } from './bpmn.service';
import { BpmnParser } from './bpmn.parser';
import { BpmnMapper } from './bpmn.mapper';
import { ModelReducer } from 'src/model/model.reducer';
import { ModelService } from 'src/model/model.service';
import { TestdataProvider } from 'test/data/testdata.provider';

describe('CatsController', () => {
  let bpmnService: BpmnService;
  let bpmnParser: BpmnParser;
  let bpmnMapper: BpmnMapper;
  let modelReducer: ModelReducer;
  let modelService: ModelService;
  const bpmnString = TestdataProvider.getConformanceExample();
  const definitions = {
    process: TestdataProvider.getProcess1()
  };
  const model1 = TestdataProvider.getModel1();
  const model2 = TestdataProvider.getModel2();

  beforeEach(async () => {
    const moduleRef = await Test.createTestingModule({
      controllers: [],
      providers: [BpmnService, BpmnParser, BpmnMapper, ModelReducer, ModelService],
    }).compile();

    bpmnService = await moduleRef.resolve(BpmnService);
    bpmnParser = await moduleRef.resolve(BpmnParser)
    bpmnMapper = await moduleRef.resolve(BpmnMapper)
    modelReducer = await moduleRef.resolve(ModelReducer)
    modelService = await moduleRef.resolve(ModelService)
  });

  describe('importBpmnProcess', () => {
    it('should call parser, mapper, reducer and service correctly', async () => {
      when(jest.spyOn(bpmnParser, 'parseBpmn'))
        .calledWith(bpmnString)
        .mockReturnValue(definitions);
      when(jest.spyOn(bpmnMapper, 'toModel'))
        .calledWith(definitions.process)
        .mockReturnValue(model1);
      when(jest.spyOn(modelReducer, 'reduceModel'))
        .calledWith(model1)
        .mockReturnValue(model2);

      jest.spyOn(modelService, 'saveModel');

      bpmnService.importBpmn(bpmnString);

      expect(modelService.saveModel).toBeCalledWith(model2);
    });
  });
});
