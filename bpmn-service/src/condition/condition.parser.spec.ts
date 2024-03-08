import { ConditionParser as ConditionParser } from './condition.parser';

describe('ConditionParser', () => {
  let conditionParser: ConditionParser;
  const condition1 = undefined;
  const condition2 = 'product not in stock';
  const condition3 = 'stock < ordered';
  const messages3 = new Map([
    ['stock', 0],
    ['ordered', 1],
  ]);
  const condition4 = '5*a > 2*b - 2';
  const messages4 = new Map([
    ['a', 0],
    ['b', 1],
  ]);
  const condition5 = '1+3*y == -1*x';
  const messages5 = new Map([
    ['x', 0],
    ['y', 1],
  ]);
  const condition6 = 'a > 3';
  const messages6 = new Map([['a', 0]]);

  const result1 = undefined;
  const result2 = undefined;
  const result3 = {
    coefficients: [1, -1],
    messageIds: [0, 1],
    offset: 0,
    comparisonOperator: 2,
  };
  const result4 = {
    coefficients: [5, -2],
    messageIds: [0, 1],
    offset: 2,
    comparisonOperator: 1,
  };
  const result5 = {
    coefficients: [3, 1],
    messageIds: [1, 0],
    offset: 1,
    comparisonOperator: 0,
  };
  const result6 = {
    coefficients: [1],
    messageIds: [0],
    offset: -3,
    comparisonOperator: 1,
  };

  beforeAll(() => {
    conditionParser = new ConditionParser();
  });

  describe('parseConditions', () => {
    it('should parse undefined', () => {
      const result = conditionParser.parseCondition(condition1!, new Map());
      expect(result).toEqual(result1);
    });
    it('should parse text', () => {
      const result = conditionParser.parseCondition(condition2, new Map());
      expect(result).toEqual(result2);
    });
    it('should parse condition without messageIds', () => {
      const result = conditionParser.parseCondition(condition3, new Map());
      expect(result).toEqual(result2);
    });
    it('should parse condition3', () => {
      const result = conditionParser.parseCondition(condition3, messages3);
      expect(result).toEqual(result3);
    });
    it('should parse condition4', () => {
      const result = conditionParser.parseCondition(condition4, messages4);
      expect(result).toEqual(result4);
    });
    it('should parse condition5', () => {
      const result = conditionParser.parseCondition(condition5, messages5);
      expect(result).toEqual(result5);
    });
    it('should parse condition6', () => {
      const result = conditionParser.parseCondition(condition6, messages6);
      expect(result).toEqual(result6);
    });
  });
});
