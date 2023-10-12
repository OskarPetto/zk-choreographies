import { TestdataProvider } from 'test/data/testdata.provider';
import { BpmnService } from './bpmn.service';

describe('BpmnService', () => {
  let bpmnService: BpmnService;
  const bpmnModel = TestdataProvider.getExampleChoreography();

  beforeAll(() => {
    bpmnService = new BpmnService();
  });

  describe('parseBpmn', () => {
    it('should parse bpmn choreography correctly', () => {
      const id = bpmnService.createBpmnModel(bpmnModel.xmlString);
      const result = bpmnService.findBpmnModelById(id);
      expect(result.xmlString).toEqual(bpmnModel.xmlString);
    });
  });
});
