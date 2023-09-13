import { Test } from '@nestjs/testing';
import { ExecutionService } from './execution.service';
import { TestdataProvider } from 'test/data/testdata.provider';
import { findTransitions } from 'test/testutils';
import { Transition } from 'src/model/petri-net/petri-net';
import { Instance } from './instance/instance';

function executeTransitions(
  executionService: ExecutionService,
  instance: Instance,
  transitions: Transition[],
): Instance {
  for (const transition of transitions) {
    instance = executionService.executeTransition(instance, transition);
  }
  return instance;
}

describe('ExecutionService', () => {
  let executionService: ExecutionService;
  const petriNet2 = TestdataProvider.getPetriNet2();
  const instance1 = TestdataProvider.getInstance1();

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      providers: [ExecutionService],
    }).compile();

    executionService = app.get<ExecutionService>(ExecutionService);
  });

  describe('instantiatePetriNet', () => {
    it('should instantiate petriNet correctly', () => {
      const result = executionService.instantiatePetriNet(petriNet2);
      expect(result.petriNet).toEqual(instance1.petriNet);
      expect(result.tokenCounts).toEqual(instance1.tokenCounts);
    });
  });

  describe('executeTransition', () => {
    it('should execute start transition', () => {
      const transitions = findTransitions(petriNet2, ['As']);
      const result = executionService.executeTransition(
        instance1,
        transitions[0],
      );
      expect(result.tokenCounts[0]).toEqual(1);
    });

    it('should not alter original instance', () => {
      const transitions = findTransitions(petriNet2, ['As']);
      executionService.executeTransition(instance1, transitions[0]);
      expect(instance1.tokenCounts[0]).toEqual(0);
    });

    it('should execute full trace without error 1', () => {
      const trace = ['As', 'Aa', 'Fa', 'Sso', 'Ro', 'Ao', 'Aaa', 'Af'];
      const transitions = findTransitions(petriNet2, trace);
      executeTransitions(executionService, instance1, transitions);
    });

    it('should execute full trace without error 2', () => {
      const trace = ['As', 'Da1', 'Af'];
      const transitions = findTransitions(petriNet2, trace);
      executeTransitions(executionService, instance1, transitions);
    });

    it('should execute full trace without error 3', () => {
      const trace = [
        'As',
        'Aa',
        'Sso',
        'Ro',
        'Co',
        'Fa',
        'Sso',
        'Ro',
        'Do',
        'Da2',
        'Af',
      ];
      const transitions = findTransitions(petriNet2, trace);
      executeTransitions(executionService, instance1, transitions);
    });

    it('should throw on invalid trace 1', () => {
      const trace = ['As', 'Da1', 'Da1', 'Af'];
      const transitions = findTransitions(petriNet2, trace);
      expect(() => {
        executeTransitions(executionService, instance1, transitions);
      }).toThrowError();
    });

    it('should throw on invalid trace 2', () => {
      const trace = ['As', 'Aa', 'Fa', 'Sso', 'Ro', 'Ao', 'Da2', 'Af'];
      const transitions = findTransitions(petriNet2, trace);
      expect(() => {
        executeTransitions(executionService, instance1, transitions);
      }).toThrowError();
    });
  });
});
