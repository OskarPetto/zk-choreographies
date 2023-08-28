import { Test } from '@nestjs/testing';

import { InstanceService } from './instance.service';
import { ModelService, ModelProvider } from 'src/model';

describe('InstanceService', () => {
    let instanceService: InstanceService;
    let modelService: ModelService;
    const model1 = ModelProvider.getModel1();

    beforeAll(async () => {
        const app = await Test.createTestingModule({
            providers: [InstanceService, ModelService],
        }).compile();

        instanceService = app.get<InstanceService>(InstanceService);
        modelService = app.get<ModelService>(ModelService);
    });

    describe('instantiateModel', () => {
        it('should reference model correctly', () => {
            jest.spyOn(modelService, 'findModel').mockImplementation(() => model1);
            const instance = instanceService.instantiateModel(model1.id);
            expect(instance.model).toEqual(model1.id);
        });
    });
});
