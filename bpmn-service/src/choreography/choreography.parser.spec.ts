import { TestdataProvider } from '../../testdata/testdata.provider';
import { ChoreographyParser } from './choreography.parser';

describe('ChoreographyParser', () => {
  let choreographyParser: ChoreographyParser;
  const definitions2 = TestdataProvider.getDefinitions2();
  const xmlString = TestdataProvider.readBpmn('example_choreography');

  beforeAll(() => {
    choreographyParser = new ChoreographyParser();
  });

  describe('parseBpmn', () => {
    it('should parse bpmn choreography correctly', () => {
      const result = choreographyParser.parseBpmn(xmlString);
      expect(result).toEqual(definitions2);
    });
  });
});
