import { TestdataProvider } from 'test/data/testdata.provider';
import { ModelReducer } from './model.reducer';

describe('ReductionService', () => {
  let modelReducer: ModelReducer;
  const model1 = TestdataProvider.getModel1();
  const model2 = TestdataProvider.getModel2();

  beforeAll(async () => {
    modelReducer = new ModelReducer();
  });

  describe('reduceModel', () => {
    it('should reduce model correctly', () => {
      const result = modelReducer.reduceModel(model1);
      expect(result).toEqual(model2);
    });
  });
});
