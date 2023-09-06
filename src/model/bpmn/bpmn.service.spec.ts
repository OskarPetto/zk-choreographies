import * as fs from 'fs';
import * as path from 'path';
import { BpmnService } from './bpmn.service';
import { Test } from '@nestjs/testing';

function readTextFile(filename: string) {
    const filePath = path.join(process.cwd(), filename);
    return fs.readFileSync(filePath, 'utf-8').toString();
}

describe('BpmnService', () => {
    let bpmnService: BpmnService;

    beforeAll(async () => {
        const app = await Test.createTestingModule({
            providers: [BpmnService],
        }).compile();

        bpmnService = app.get<BpmnService>(BpmnService);
    });

    describe('parseBpmn', () => {
        it('should parse model correctly', () => {
            const bpmnString = readTextFile('test/data/conformance_example.bpmn');
            const model = bpmnService.parseModel(bpmnString);
        });
    });
});
