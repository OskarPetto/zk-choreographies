import { Test } from '@nestjs/testing';

import { InstanceService } from './instance.service';
import { ExecutionStatus } from './instance';
import { TestdataProvider } from 'test/data/provider';

describe('InstanceService', () => {
    let instanceService: InstanceService;
    const model1 = TestdataProvider.getModel1();

    beforeAll(async () => {
        const app = await Test.createTestingModule({
            providers: [InstanceService],
        }).compile();

        instanceService = app.get<InstanceService>(InstanceService);
    });

    describe('instantiateModel', () => {
        it('should instantiate model correctly', () => {
            const instance = instanceService.instantiateModel(model1);
            expect(instance.model).toEqual(model1.id);
            expect(Array.from(instance.executionStatuses.values())).toEqual(Array(model1.places.size).fill(ExecutionStatus.NOT_ACTIVE));
            expect(instance.finished).toBeFalsy();
        });
    });
});
