import { Test } from '@nestjs/testing';
import { ExecutionService } from './execution.service';
import { TestdataProvider } from 'test/data/provider';
import { ExecutionStatus } from './instance';
import { Model, Transition, TransitionId } from 'src/model';

function findTransition(model: Model, transitionId: TransitionId): Transition {
  const transition = model.transitions.get(transitionId);
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
      expect(result.executionStatuses.get('Flow_0xqg6as')).toEqual(
        ExecutionStatus.ACTIVE,
      );
    });

    it('should not alter original instance', () => {
      const startTransition = findTransition(model2, 'As');
      executionService.executeTransitions(instance1, [startTransition]);
      expect(instance1.executionStatuses.get('Flow_0xqg6as')).toEqual(
        ExecutionStatus.NOT_ACTIVE,
      );
    });

    it('should execute full trace 1', () => {
      const trace = ['As', 'Aa', 'Fa', 'Sso', 'Ro', 'Ao', 'Aaa', 'Af'];
      const transitions = trace.map((transitionId) =>
        findTransition(model2, transitionId),
      );
      const result = executionService.executeTransitions(
        instance1,
        transitions,
      );
      expect(result.finished).toBeTruthy();
    });

    it('should execute full trace 2', () => {
      const trace = ['As', 'Da1', 'Af'];
      const transitions = trace.map((transitionId) =>
        findTransition(model2, transitionId),
      );
      const result = executionService.executeTransitions(
        instance1,
        transitions,
      );
      expect(result.finished).toBeTruthy();
    });

    it('should execute full trace 3', () => {
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
      const result = executionService.executeTransitions(
        instance1,
        transitions,
      );
      expect(result.finished).toBeTruthy();
    });

    it('should not be finished on incomplete trace', () => {
      const trace = ['As', 'Aa', 'Sso', 'Ro'];
      const transitions = trace.map((transitionId) =>
        findTransition(model2, transitionId),
      );
      const result = executionService.executeTransitions(
        instance1,
        transitions,
      );
      expect(result.finished).toBeFalsy();
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
