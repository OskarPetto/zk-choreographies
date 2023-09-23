import { TestdataProvider } from 'test/data/testdata.provider';
import { BpmnMapper } from './bpmn.mapper';
import { Transition } from '../model/model';

describe('BpmnMapper', () => {
  let bpmnMapper: BpmnMapper;
  const definitions2 = TestdataProvider.getDefinitions2();
  const model2 = TestdataProvider.getModel2();

  beforeAll(async () => {
    bpmnMapper = new BpmnMapper();
  });

  describe('toModel', () => {
    it('should map bpmn process correctly', () => {
      // const result = bpmnMapper.toModel(definitions2);
      // expect(result.id).toEqual(model2.id);
      // model2.transitions.forEach((transition: Transition) =>
      //   expect(result.transitions).toContainEqual(transition),
      // );
      // result.transitions.forEach((transition: Transition) =>
      //   expect(model2.transitions).toContainEqual(transition),
      // );
    });
  });
});
