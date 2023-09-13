import { TestdataProvider } from 'test/data/testdata.provider';
import { PetriNetService } from './petri-net.service';

describe('PetriNetService', () => {
  let petriNetService: PetriNetService;
  const petriNet1 = TestdataProvider.getPetriNet1();
  const petriNet2 = TestdataProvider.getPetriNet2();

  beforeAll(async () => {
    petriNetService = new PetriNetService();
    petriNet1.id = 'other';
    petriNetService.savePetriNet(petriNet1);
    petriNetService.savePetriNet(petriNet2);
  });

  describe('findPetriNet', () => {
    it('should find petriNet1', () => {
      const result = petriNetService.findPetriNet(petriNet1.id);
      expect(result).toEqual(petriNet1);
    });
  });
});
