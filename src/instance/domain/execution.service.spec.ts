import { Test } from '@nestjs/testing';
import { ExecutionService } from './execution.service';
import { TestdataProvider } from 'test/data/provider';
import { Model, Transition, TransitionId } from 'src/model';

function findTransition(model: Model, transitionId: TransitionId): Transition {
  const transition = model.transitions.find((t) => t.id === transitionId);
  if (!transition) {
    throw Error(`Transition ${transitionId} in model ${model.id} not found`);
  }
  return transition;
}

describe('ExecutionService', () => {
  let executionService: ExecutionService;
  const model2 = TestdataProvider.getModel2();
  const instance1 = TestdataProvider.getInstance1();

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      providers: [ExecutionService],
    }).compile();

    executionService = app.get<ExecutionService>(ExecutionService);
  });

  describe('executeTransition', () => {
    it('should execute start transition', () => {
      const startTransition = findTransition(model2, 'As');
      const result = executionService.executeTransitions(instance1, [
        startTransition,
      ]);
      expect(result.tokenCounts[0]).toEqual(1);
    });

    it('should not alter original instance', () => {
      const startTransition = findTransition(model2, 'As');
      executionService.executeTransitions(instance1, [startTransition]);
      expect(instance1.tokenCounts[0]).toEqual(0);
    });

    it('should execute full trace without error 1', () => {
      const trace = ['As', 'Aa', 'Fa', 'Sso', 'Ro', 'Ao', 'Aaa', 'Af'];
      const transitions = trace.map((transitionId) =>
        findTransition(model2, transitionId),
      );
      executionService.executeTransitions(instance1, transitions);
    });

    it('should execute full trace without error 2', () => {
      const trace = ['As', 'Da1', 'Af'];
      const transitions = trace.map((transitionId) =>
        findTransition(model2, transitionId),
      );
      executionService.executeTransitions(instance1, transitions);
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
      const transitions = trace.map((transitionId) =>
        findTransition(model2, transitionId),
      );
      executionService.executeTransitions(instance1, transitions);
    });

    it('should throw on invalid trace 1', () => {
      const trace = ['As', 'Da1', 'Da1', 'Af'];
      const transitions = trace.map((transitionId) =>
        findTransition(model2, transitionId),
      );
      expect(() =>
        executionService.executeTransitions(instance1, transitions),
      ).toThrowError();
    });

    it('should throw on invalid trace 2', () => {
      const trace = ['As', 'Aa', 'Fa', 'Sso', 'Ro', 'Ao', 'Da2', 'Af'];
      const transitions = trace.map((transitionId) =>
        findTransition(model2, transitionId),
      );
      expect(() =>
        executionService.executeTransitions(instance1, transitions),
      ).toThrowError();
    });
  });
});
