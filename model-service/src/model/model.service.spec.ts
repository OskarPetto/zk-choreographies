import { TestdataProvider } from 'test/data/testdata.provider';
import { ModelService } from './model.service';

describe('ModelService', () => {
  let modelService: ModelService;
  const model2 = TestdataProvider.getModel2();

  beforeAll(async () => {
    modelService = new ModelService();
    modelService.saveModel(model2);
  });

  describe('findModel', () => {
    it('should find model1', () => {
      const result = modelService.findModel(model2.id);
      expect(result).toEqual(model2);
    });
  });
});
