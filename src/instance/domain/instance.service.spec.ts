import { TestdataProvider } from 'test/data/provider';
import { InstanceService } from './instance.service';
import { Test } from '@nestjs/testing';

describe('InstanceService', () => {
  let instanceService: InstanceService;
  const model2 = TestdataProvider.getModel2();
  const instance1 = TestdataProvider.getInstance1();

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      providers: [InstanceService],
    }).compile();

    instanceService = app.get<InstanceService>(InstanceService);
  });

  describe('instantiateModel', () => {
    it('should instantiate model correctly', () => {
      const result = instanceService.instantiateModel(model2);
      expect(result.model).toEqual(instance1.model);
      expect(result.tokenCounts).toEqual(instance1.tokenCounts);
    });
  });
});
