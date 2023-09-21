import { TestdataProvider } from 'test/data/testdata.provider';
import { BpmnMapper } from './bpmn.mapper';
import { Transition } from '../petri-net/petri-net';

describe('BpmnMapper', () => {
  let bpmnMapper: BpmnMapper;
  const process1 = TestdataProvider.getProcess1();
  const petriNet1 = TestdataProvider.getPetriNet1();

  beforeAll(async () => {
    bpmnMapper = new BpmnMapper();
  });

  describe('toPetriNet', () => {
    it('should map bpmn process correctly', () => {
      const result = bpmnMapper.toPetriNet(process1);
      expect(result.id).toEqual(petriNet1.id);
      petriNet1.transitions.forEach((transition: Transition) =>
        expect(result.transitions).toContainEqual(transition),
      );
      result.transitions.forEach((transition: Transition) =>
        expect(petriNet1.transitions).toContainEqual(transition),
      );
    });
  });
});
