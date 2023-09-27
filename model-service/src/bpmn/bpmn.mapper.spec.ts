import { TestdataProvider } from 'test/data/testdata.provider';
import { BpmnMapper } from './bpmn.mapper';

describe('BpmnMapper', () => {
  let bpmnMapper: BpmnMapper;
  const definitions2 = TestdataProvider.getDefinitions2();
  const model2 = TestdataProvider.getModel2();

  beforeAll(() => {
    bpmnMapper = new BpmnMapper();
    jest.useFakeTimers();
    jest.setSystemTime(TestdataProvider.getDate());
  });

  describe('toModel', () => {
    it('should map bpmn process correctly', () => {
      const result = bpmnMapper.toModel(definitions2.choreographies[0]);
      expect(result).toEqual(model2);
    });
  });

  afterAll(() => {
    jest.useRealTimers();
  });
});
