import { TestdataProvider } from 'test/data/testdata.provider';
import { BpmnParser } from './bpmn.parser';

describe('BpmnParser', () => {
  let bpmnParser: BpmnParser;
  const definitions2 = TestdataProvider.getDefinitions2();
  const bpmnString = TestdataProvider.getExampleChoreography();

  beforeAll(async () => {
    bpmnParser = new BpmnParser();
  });

  describe('parseBpmn', () => {
    it('should parse bpmn process correctly', () => {
      const result = bpmnParser.parseBpmn(bpmnString);
      expect(result).toEqual(definitions2);
    });
  });
});
