import { Test } from '@nestjs/testing';
import { ChoreographyService } from './choreography.service';
import { TestdataProvider } from 'test/data/provider';
import { ExecutionStatus } from './instance';
import { findTransition } from 'src/model';

describe('ChoreographyService', () => {
    let choreographyService: ChoreographyService;
    const model1 = TestdataProvider.getModel1();
    const instance1 = TestdataProvider.getInstance1();

    beforeAll(async () => {
        const app = await Test.createTestingModule({
            providers: [ChoreographyService],
        }).compile();

        choreographyService = app.get<ChoreographyService>(ChoreographyService);
    });

    describe('executeTransition', () => {
        it('should execute start transition', () => {
            const startTransition = findTransition(model1, 'As');
            const result = choreographyService.executeTransition(instance1, startTransition);
            expect(result.executionStatuses.get(0)).toEqual(ExecutionStatus.ACTIVE);
        });

        it('should not alter original instance', () => {
            const startTransition = findTransition(model1, 'As');
            const result = choreographyService.executeTransition(instance1, startTransition);
            expect(instance1.executionStatuses.get(0)).toEqual(ExecutionStatus.NOT_ACTIVE);
        });

        it('should execute full trace 1', () => {
            const trace = ['As', 'Aa', 'Fa', 'Sso', 'Ro', 'Ao', 'Aaa', 'Af'];
            const transitions = trace.map(transitionId => findTransition(model1, transitionId));
            const result = choreographyService.executeTransitions(instance1, transitions);
            expect(result.finished).toBeTruthy();
        });

        it('should execute full trace 2', () => {
            const trace = ['As', 'Da1', 'Af'];
            const transitions = trace.map(transitionId => findTransition(model1, transitionId));
            const result = choreographyService.executeTransitions(instance1, transitions);
            expect(result.finished).toBeTruthy();
        });

        it('should execute full trace 3', () => {
            const trace = ['As', 'Aa', 'Sso', 'Ro', 'Co', 'Fa', 'Sso', 'Ro', 'Do', 'Da2', 'Af'];
            const transitions = trace.map(transitionId => findTransition(model1, transitionId));
            const result = choreographyService.executeTransitions(instance1, transitions);
            expect(result.finished).toBeTruthy();
        });

        it('should not be finished on incomplete trace', () => {
            const trace = ['As', 'Aa', 'Sso', 'Ro'];
            const transitions = trace.map(transitionId => findTransition(model1, transitionId));
            const result = choreographyService.executeTransitions(instance1, transitions);
            expect(result.finished).toBeFalsy();
        });

        it('should throw on invalid trace 1', () => {
            const trace = ['As', 'Da1', 'Da1', 'Af'];
            const transitions = trace.map(transitionId => findTransition(model1, transitionId));
            expect(() => choreographyService.executeTransitions(instance1, transitions)).toThrowError();
        });

        it('should throw on invalid trace 2', () => {
            const trace = ['As', 'Aa', 'Fa', 'Sso', 'Ro', 'Ao', 'Da2', 'Af'];
            const transitions = trace.map(transitionId => findTransition(model1, transitionId));
            expect(() => choreographyService.executeTransitions(instance1, transitions)).toThrowError();
        });
    });
});
