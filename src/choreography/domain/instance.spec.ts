import { TestdataProvider } from "test/data/provider";
import { ExecutionStatus, instantiateModel } from "./instance";

describe('Instance', () => {
    const model1 = TestdataProvider.getModel1();

    describe('fromModel', () => {
        it('should instantiate model correctly', () => {
            const instance = instantiateModel(model1);
            expect(instance.model).toEqual(model1.id);
            expect(instance.executionStatuses).toEqual(Array(model1.flowCount).fill(ExecutionStatus.NOT_ACTIVE));
            expect(instance.finished).toBeFalsy();
        });
    });
});