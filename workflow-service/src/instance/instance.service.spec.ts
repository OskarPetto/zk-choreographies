import { TestdataProvider } from 'test/data/testdata.provider';
import { InstanceService } from './instance.service';
import { Test } from '@nestjs/testing';

describe('InstanceService', () => {
  let instanceService: InstanceService;
  const instance1 = TestdataProvider.getInstance1();

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      providers: [InstanceService],
    }).compile();

    instanceService = app.get<InstanceService>(InstanceService);
  });

  describe('findInstance', () => {
    it('should find instance', () => {
      instanceService.saveInstance(instance1);
      const result = instanceService.findInstance(instance1.id!);
      expect(result).toEqual(instance1);
    });
  });
});
