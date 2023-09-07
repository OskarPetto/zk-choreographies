import * as fs from 'fs';
import * as path from 'path';
import { BpmnService } from './bpmn.service';
import { Test } from '@nestjs/testing';
import { TestdataProvider } from 'test/data/provider';
import { findFlowMapping } from 'test/testutils';

function readTextFile(filename: string) {
  const filePath = path.join(process.cwd(), filename);
  return fs.readFileSync(filePath, 'utf-8').toString();
}

describe('BpmnService', () => {
  let bpmnService: BpmnService;
  const model3 = TestdataProvider.getModel3();

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      providers: [BpmnService],
    }).compile();

    bpmnService = app.get<BpmnService>(BpmnService);
  });

  describe('parseBpmn', () => {
    it('should parse model correctly', () => {
      const bpmnString = readTextFile('test/data/conformance_example.bpmn');
      const result = bpmnService.parseProcess(bpmnString);
      const flowMapping = findFlowMapping(model3, result);
      expect(flowMapping).toBeTruthy();
    });
  });
});
