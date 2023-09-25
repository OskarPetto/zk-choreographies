import { TestdataProvider } from 'test/data/testdata.provider';
import { ModelReducer } from './model.reducer';
import { logObject } from 'test/testutils';
import { BpmnMapper } from 'src/bpmn/bpmn.mapper';

describe('ReductionService', () => {
  let modelReducer: ModelReducer;
  const model2 = TestdataProvider.getModel2();
  const model2Reduced = TestdataProvider.getModel2Reduced();

  beforeAll(async () => {
    modelReducer = new ModelReducer();
  });

  describe('reduceModel', () => {
    it('should reduce model correctly', () => {
      const result = modelReducer.reduceModel(model2);
      expect(result).toEqual(model2Reduced)
    });
  });
});
