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
            const instance = choreographyService.executeTransition(instance1, startTransition);
            expect(instance.executionStatuses[0]).toEqual(ExecutionStatus.ACTIVE);
        });

        it('should execute full trace', () => {
            const trace = ['As', 'Aa', 'Fa', 'Sso', 'Ro', 'Ao', 'Aaa', 'Af'];
            let instance = instance1;
            for (const transitionId of trace) {
                const transition = Testdata.getModel1Transition(transitionId);
                instance = choreographyService.executeTransition(instance, transition);
            }
            expect(instance.finished).toBeTruthy();
        });
    });
});
