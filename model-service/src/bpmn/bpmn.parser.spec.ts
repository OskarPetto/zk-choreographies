import { Test } from '@nestjs/testing';
import { TestdataProvider } from 'test/data/testdata.provider';
import { BpmnParser } from './bpmn.parser';

describe('BpmnParser', () => {
  let bpmnParser: BpmnParser;
  const process1 = TestdataProvider.getProcess1();
  const bpmnString = TestdataProvider.getConformanceExample();

  beforeAll(async () => {
    bpmnParser = new BpmnParser();
  });

  describe('parseBpmn', () => {
    it('should parse bpmn process correctly', () => {
      const result = bpmnParser.parseBpmn(bpmnString);
      expect(result.process).toEqual(process1);
    });
  });
});
