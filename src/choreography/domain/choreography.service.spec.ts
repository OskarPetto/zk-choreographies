import { Test } from '@nestjs/testing';
import { ChoreographyService } from './choreography.service';
import { Testdata } from 'test/testdata';
import { ExecutionStatus } from './instance';

describe('ChoreographyService', () => {
    let choreographyService: ChoreographyService;
    const model1 = Testdata.getModel1();
    const instance1 = Testdata.getInstance1();

    beforeAll(async () => {
        const app = await Test.createTestingModule({
            providers: [ChoreographyService],
        }).compile();

        choreographyService = app.get<ChoreographyService>(ChoreographyService);
    });

    describe('executeTransition', () => {
        it('should execute start transition', () => {
            const startTransition = Testdata.getModel1Transition('As');
            const result = choreographyService.executeTransition(instance1, startTransition);
            expect(result.executionStatuses[0]).toEqual(ExecutionStatus.ACTIVE);
        });

        it('should execute full trace 1', () => {
            const trace = ['As', 'Aa', 'Fa', 'Sso', 'Ro', 'Ao', 'Aaa', 'Af'];
            const transitions = trace.map(transitionId => Testdata.getModel1Transition(transitionId));
            const result = choreographyService.executeTransitions(instance1, transitions);
            expect(result.finished).toBeTruthy();
        });

        it('should execute full trace 2', () => {
            const trace = ['As', 'Da1', 'Af'];
            const transitions = trace.map(transitionId => Testdata.getModel1Transition(transitionId));
            const result = choreographyService.executeTransitions(instance1, transitions);
            expect(result.finished).toBeTruthy();
        });

        it('should execute full trace 3', () => {
            const trace = ['As', 'Aa', 'Sso', 'Ro', 'Co', 'Fa', 'Sso', 'Ro', 'Do', 'Da2', 'Af'];
            const transitions = trace.map(transitionId => Testdata.getModel1Transition(transitionId));
            const result = choreographyService.executeTransitions(instance1, transitions);
            expect(result.finished).toBeTruthy();
        });

        it('should not be finished on incomplete trace', () => {
            const trace = ['As', 'Aa', 'Sso', 'Ro'];
            const transitions = trace.map(transitionId => Testdata.getModel1Transition(transitionId));
            const result = choreographyService.executeTransitions(instance1, transitions);
            expect(result.finished).toBeFalsy();
        });

        it('should throw on invalid trace 1', () => {
            const trace = ['As', 'Da1', 'Da1', 'Af'];
            const transitions = trace.map(transitionId => Testdata.getModel1Transition(transitionId));
            expect(() => choreographyService.executeTransitions(instance1, transitions)).toThrowError();
        });

        it('should throw on invalid trace 2', () => {
            const trace = ['As', 'Aa', 'Fa', 'Sso', 'Ro', 'Ao', 'Da2', 'Af'];
            const transitions = trace.map(transitionId => Testdata.getModel1Transition(transitionId));
            expect(() => choreographyService.executeTransitions(instance1, transitions)).toThrowError();
        });
    });
});
