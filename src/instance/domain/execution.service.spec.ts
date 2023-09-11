import { Test } from '@nestjs/testing';
import { ExecutionService } from './execution.service';
import { TestdataProvider } from 'test/data/provider';
import { ExecutionStatus } from './instance';

describe('ExecutionService', () => {
  let executionService: ExecutionService;
  const instance1 = TestdataProvider.getInstance1();

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      providers: [ExecutionService],
    }).compile();

    executionService = app.get<ExecutionService>(ExecutionService);
  });

  describe('executeTransition', () => {
    it('should execute start transition', () => {
      const result = executionService.executeTransitions(instance1, ['As']);
      expect(result.executionStatuses[0]).toEqual(ExecutionStatus.ACTIVE);
    });

    it('should not alter original instance', () => {
      executionService.executeTransitions(instance1, ['As']);
      expect(instance1.executionStatuses[0]).toEqual(
        ExecutionStatus.NOT_ACTIVE,
      );
    });

    it('should execute full trace 1', () => {
      const result = executionService.executeTransitions(instance1, [
        'As',
        'Aa',
        'Fa',
        'Sso',
        'Ro',
        'Ao',
        'Aaa',
        'Af',
      ]);
      expect(result.finished).toBeTruthy();
    });

    it('should execute full trace 2', () => {
      const result = executionService.executeTransitions(instance1, [
        'As',
        'Da1',
        'Af',
      ]);
      expect(result.finished).toBeTruthy();
    });

    it('should execute full trace 3', () => {
      const result = executionService.executeTransitions(instance1, [
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
      ]);
      expect(result.finished).toBeTruthy();
    });

    it('should not be finished on incomplete trace', () => {
      const result = executionService.executeTransitions(instance1, [
        'As',
        'Aa',
        'Sso',
        'Ro',
      ]);
      expect(result.finished).toBeFalsy();
    });

    it('should throw on invalid trace 1', () => {
      expect(() =>
        executionService.executeTransitions(instance1, [
          'As',
          'Da1',
          'Da1',
          'Af',
        ]),
      ).toThrowError();
    });

    it('should throw on invalid trace 2', () => {
      expect(() =>
        executionService.executeTransitions(instance1, [
          'As',
          'Aa',
          'Fa',
          'Sso',
          'Ro',
          'Ao',
          'Da2',
          'Af',
        ]),
      ).toThrowError();
    });
  });
});
