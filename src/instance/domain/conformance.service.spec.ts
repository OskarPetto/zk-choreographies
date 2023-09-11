import { Test } from '@nestjs/testing';
import { ExecutionService } from './execution.service';
import { TestdataProvider } from 'test/data/provider';
import { ConformanceService } from './conformance.service';
import { copyInstance } from './instance';
import { findTransition } from 'test/testutils';
import { copyModel } from 'src/model';

describe('ConformanceService', () => {
  let conformanceService: ConformanceService;
  let executionService: ExecutionService;
  const model2 = TestdataProvider.getModel2();
  const instance1 = TestdataProvider.getInstance1();

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      providers: [ConformanceService, ExecutionService],
    }).compile();

    conformanceService = app.get<ConformanceService>(ConformanceService);
    executionService = app.get<ExecutionService>(ExecutionService);
  });

  describe('isExecutionValid', () => {
    it('should return true on valid execution', () => {
      const transition = findTransition(model2, 'As');
      const instance2 = executionService.executeTransition(instance1, transition);
      expect(conformanceService.isExecutionValid(instance1, instance2, transition, model2)).toBeTruthy();
    });

    it('should check instance id', () => {
      const transition = findTransition(model2, 'As');
      const instance2 = executionService.executeTransition(instance1, transition);
      instance2.id = 'other';
      expect(conformanceService.isExecutionValid(instance1, instance2, transition, model2)).toBeFalsy();
    });

    it('should check whether instance is of model', () => {
      const transition = findTransition(model2, 'As');
      const instance2 = executionService.executeTransition(instance1, transition);
      instance2.model = 'other';
      expect(conformanceService.isExecutionValid(instance1, instance2, transition, model2)).toBeFalsy();
    });

    it('should check instance tokens', () => {
      const transition = findTransition(model2, 'As');
      const instance2 = executionService.executeTransition(instance1, transition);
      instance2.tokenCounts[0] = -1;
      expect(conformanceService.isExecutionValid(instance1, instance2, transition, model2)).toBeFalsy();
    });

    it('should check transition id', () => {
      const transition = findTransition(model2, 'As');
      const instance2 = executionService.executeTransition(instance1, transition);
      const copiedTransition = findTransition(copyModel(model2), 'As');
      copiedTransition.id = 'other';
      expect(conformanceService.isExecutionValid(instance1, instance2, copiedTransition, model2)).toBeFalsy();
    });

    it('should check transition places', () => {
      const transition = findTransition(model2, 'As');
      const instance2 = executionService.executeTransition(instance1, transition);
      const copiedTransition = findTransition(copyModel(model2), 'As');
      copiedTransition.incomingPlaces = [1];
      expect(conformanceService.isExecutionValid(instance1, instance2, copiedTransition, model2)).toBeFalsy();
    });

    it('should check token changes', () => {
      const transition = findTransition(model2, 'As');
      const instance2 = executionService.executeTransition(instance1, transition);
      instance2.tokenCounts[0] = 0;
      expect(conformanceService.isExecutionValid(instance1, instance2, transition, model2)).toBeFalsy();
    });

    it('should return only trues on full example', () => {
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

      let instanceBefore = instance1;
      let instanceAfter;
      for (const transition of transitions) {
        instanceAfter = executionService.executeTransition(instanceBefore, transition);
        expect(conformanceService.isExecutionValid(instanceBefore, instanceAfter, transition, model2)).toBeTruthy();
        instanceBefore = instanceAfter;
      }
    });
  });
});
