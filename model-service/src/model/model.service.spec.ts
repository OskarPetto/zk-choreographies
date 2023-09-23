import { TestdataProvider } from 'test/data/testdata.provider';
import { ModelService } from './model.service';

describe('ModelService', () => {
  let modelService: ModelService;
  const model1 = TestdataProvider.getModel1();
  const model2 = TestdataProvider.getModel2();

  beforeAll(async () => {
    modelService = new ModelService();
    model1.id = 'other';
    modelService.saveModel(model1);
    modelService.saveModel(model2);
  });

  describe('findModel', () => {
    it('should find model1', () => {
      const result = modelService.findModel(model1.id);
      expect(result).toEqual(model1);
    });
  });
});
