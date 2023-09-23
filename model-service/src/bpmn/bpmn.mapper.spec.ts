import { TestdataProvider } from 'test/data/testdata.provider';
import { BpmnMapper } from './bpmn.mapper';
import { Transition } from '../model/model';

describe('BpmnMapper', () => {
  let bpmnMapper: BpmnMapper;
  const process1 = TestdataProvider.getChoreography1();
  const model1 = TestdataProvider.getModel1();

  beforeAll(async () => {
    bpmnMapper = new BpmnMapper();
  });

  describe('toModel', () => {
    it('should map bpmn process correctly', () => {
      const result = bpmnMapper.toModel(process1);
      expect(result.id).toEqual(model1.id);
      model1.transitions.forEach((transition: Transition) =>
        expect(result.transitions).toContainEqual(transition),
      );
      result.transitions.forEach((transition: Transition) =>
        expect(model1.transitions).toContainEqual(transition),
      );
    });
  });
});
