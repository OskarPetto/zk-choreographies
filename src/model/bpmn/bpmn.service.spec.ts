import * as fs from 'fs';
import * as path from 'path';
import { BpmnService } from './bpmn.service';
import { Test } from '@nestjs/testing';
import { TestdataProvider } from 'test/data/provider';
import { ReductionService } from '../domain/reduction.service';
import { logObject } from 'src/utils/utils';
import { findFlowMapping } from 'test/testutils';

function readTextFile(filename: string) {
    const filePath = path.join(process.cwd(), filename);
    return fs.readFileSync(filePath, 'utf-8').toString();
}

describe('BpmnService', () => {
    let bpmnService: BpmnService;
    let reductionService: ReductionService;
    const model3 = TestdataProvider.getModel3();

    beforeAll(async () => {
        const app = await Test.createTestingModule({
            providers: [BpmnService, ReductionService],
        }).compile();

        bpmnService = app.get<BpmnService>(BpmnService);
        reductionService = app.get<ReductionService>(ReductionService);
    });

    describe('parseBpmn', () => {
        it('should parse model correctly', () => {
            const bpmnString = readTextFile('test/data/conformance_example.bpmn');
            const result = bpmnService.parseModel(bpmnString);
            const flowMapping = findFlowMapping(model3, result);
            expect(flowMapping).toBeTruthy();
        });
    });
});
