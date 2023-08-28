import { Test } from '@nestjs/testing';

import { InstanceService } from './instance.service';
import { ExecutionStatus } from './instance';
import { Testdata } from 'test/testdata';

describe('InstanceService', () => {
    let instanceService: InstanceService;
    const model1 = Testdata.getModel1();

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
            expect(instance.executionStatuses).toEqual(Array(model1.placeCount).fill(ExecutionStatus.NOT_ACTIVE));
            expect(instance.finished).toBeFalsy();
        });
    });
});
