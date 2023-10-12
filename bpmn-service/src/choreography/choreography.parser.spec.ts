import { TestdataProvider } from '../../test/data/testdata.provider';
import { ChoreographyParser } from './choreography.parser';

describe('ChoreographyParser', () => {
  let choreographyParser: ChoreographyParser;
  const definitions2 = TestdataProvider.getDefinitions2();
  const bpmnModel = TestdataProvider.getExampleChoreography();

  beforeAll(() => {
    choreographyParser = new ChoreographyParser();
  });

  describe('parseBpmn', () => {
    it('should parse bpmn choreography correctly', () => {
      const result = choreographyParser.parseBpmn(bpmnModel.xmlString);
      expect(result).toEqual(definitions2);
    });
  });
});
