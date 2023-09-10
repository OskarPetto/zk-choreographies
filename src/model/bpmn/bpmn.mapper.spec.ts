import { Test } from '@nestjs/testing';
import { TestdataProvider } from 'test/data/provider';
import { BpmnMapper } from './bpmn.mapper';

describe('BpmnMapper', () => {
  let bpmnMapper: BpmnMapper;
  const process1 = TestdataProvider.getProcess1();
  const model1 = TestdataProvider.getModel1();

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      providers: [BpmnMapper],
    }).compile();

    bpmnMapper = app.get<BpmnMapper>(BpmnMapper);
  });

  describe('toModel', () => {
    it('should map bpmn process correctly', () => {
      const result = bpmnMapper.toModel(process1);
      expect(result).toEqual(model1);
    });
  });
});
