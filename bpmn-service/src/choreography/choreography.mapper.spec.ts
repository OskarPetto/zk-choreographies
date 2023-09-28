import { TestdataProvider } from '../../test/data/testdata.provider';
import { ChoreographyMapper } from './choreography.mapper';

describe('ChoreographyMapper', () => {
  let choreographyMapper: ChoreographyMapper;
  const definitions2 = TestdataProvider.getDefinitions2();
  const model2 = TestdataProvider.getModel2();

  beforeAll(() => {
    choreographyMapper = new ChoreographyMapper();
    jest.useFakeTimers();
    jest.setSystemTime(TestdataProvider.getDate());
  });

  describe('toModel', () => {
    it('should map bpmn choreography correctly', () => {
      const result = choreographyMapper.toModel(definitions2.choreographies[0]);
      expect(result).toEqual(model2);
    });
  });

  afterAll(() => {
    jest.useRealTimers();
  });
});
