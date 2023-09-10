import { Test } from '@nestjs/testing';

import { TestdataProvider } from 'test/data/provider';
import { ModelReducer } from './model.reducer';
import { findPlaceMapping } from 'test/testutils';

describe('ReductionService', () => {
  let modelReducer: ModelReducer;
  const model1 = TestdataProvider.getModel1();
  const model2 = TestdataProvider.getModel2();

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      providers: [ModelReducer],
    }).compile();

    modelReducer = app.get<ModelReducer>(ModelReducer);
  });

  describe('reduceModel', () => {
    it('should reduce model correctly', () => {
      const result = modelReducer.reduceModel(model1);
      const placeMapping = findPlaceMapping(model2, result);
      expect(placeMapping).toBeTruthy();
    });
  });
});
