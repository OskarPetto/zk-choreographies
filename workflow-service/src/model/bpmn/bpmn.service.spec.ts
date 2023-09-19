import { when } from 'jest-when';
import { Test } from '@nestjs/testing';
import { BpmnService } from './bpmn.service';
import { BpmnParser } from './bpmn.parser';
import { BpmnMapper } from './bpmn.mapper';
import { TestdataProvider } from 'test/data/testdata.provider';
import { PetriNetReducer } from '../petri-net/perti-net.reducer';
import { PetriNetService } from '../petri-net/petri-net.service';

describe('CatsController', () => {
  let bpmnService: BpmnService;
  let bpmnParser: BpmnParser;
  let bpmnMapper: BpmnMapper;
  let petriNetReducer: PetriNetReducer;
  let petriNetService: PetriNetService;
  const bpmnString = TestdataProvider.getConformanceExample();
  const definitions = {
    process: TestdataProvider.getProcess1(),
  };
  const petriNet1 = TestdataProvider.getPetriNet1();
  const petriNet2 = TestdataProvider.getPetriNet2();

  beforeEach(async () => {
    const moduleRef = await Test.createTestingModule({
      controllers: [],
      providers: [
        BpmnService,
        BpmnParser,
        BpmnMapper,
        PetriNetReducer,
        PetriNetService,
      ],
    }).compile();

    bpmnService = await moduleRef.resolve(BpmnService);
    bpmnParser = await moduleRef.resolve(BpmnParser);
    bpmnMapper = await moduleRef.resolve(BpmnMapper);
    petriNetReducer = await moduleRef.resolve(PetriNetReducer);
    petriNetService = await moduleRef.resolve(PetriNetService);
  });

  describe('importBpmnProcess', () => {
    it('should call parser, mapper, reducer and service correctly', async () => {
      when(jest.spyOn(bpmnParser, 'parseBpmn'))
        .calledWith(bpmnString)
        .mockReturnValue(definitions);
      when(jest.spyOn(bpmnMapper, 'toPetriNet'))
        .calledWith(definitions.process)
        .mockReturnValue(petriNet1);
      when(jest.spyOn(petriNetReducer, 'reducePetriNet'))
        .calledWith(petriNet1)
        .mockReturnValue(petriNet2);

      jest.spyOn(petriNetService, 'savePetriNet');

      bpmnService.importBpmn(bpmnString);

      expect(petriNetService.savePetriNet).toBeCalledWith(petriNet2);
    });
  });
});
