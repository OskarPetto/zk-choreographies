import { TestdataProvider } from 'test/data/testdata.provider';
import { PetriNetReducer } from './perti-net.reducer';

describe('ReductionService', () => {
  let petriNetReducer: PetriNetReducer;
  const petriNet1 = TestdataProvider.getPetriNet1();
  const petriNet2 = TestdataProvider.getPetriNet2();

  beforeAll(async () => {
    petriNetReducer = new PetriNetReducer();
  });

  describe('reducePetriNet', () => {
    it('should reduce petriNet correctly', () => {
      const result = petriNetReducer.reducePetriNet(petriNet1);
      expect(result).toEqual(petriNet2);
    });
  });
});
