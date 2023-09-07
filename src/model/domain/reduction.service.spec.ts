import { Test } from '@nestjs/testing';

import { TestdataProvider } from 'test/data/provider';
import { ReductionService } from './reduction.service';
import { findFlowMapping } from 'test/testutils';

describe('ReductionService', () => {
    let reductionService: ReductionService;
    const model1 = TestdataProvider.getModel1();
    const model2 = TestdataProvider.getModel2();

    beforeAll(async () => {
        const app = await Test.createTestingModule({
            providers: [ReductionService],
        }).compile();

        reductionService = app.get<ReductionService>(ReductionService);
    });

    describe('reduceModel', () => {
        it('should reduce model correctly', () => {
            const result = reductionService.reduceModel(model2);
            const flowMapping = findFlowMapping(model1, result);
            expect(flowMapping).toBeTruthy();
        });
    });
});
