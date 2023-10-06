import { SequenceFlow } from 'src/choreography/choreography';
import { ConstraintParser } from './constraint.parser';
import { defaultConstraint } from './constraint';
import { TestdataProvider } from 'test/data/testdata.provider';

describe('ConstraintParser', () => {
  let constraintParser: ConstraintParser;
  const constraint1 = undefined
  const constraint2 = 'product not in stock'
  const constraint3 = 'stock < ordered'
  const messages3 = new Map([['stock', 0], ['ordered', 1]])
  const constraint4 = '5*a > 2*b - 2'
  const messages4 = new Map([['a', 0], ['b', 1]])
  const constraint5 = '1+3*y == -1*x'
  const messages5 = new Map([['x', 0], ['y', 1]])
  const constraint6 = 'a > 3'
  const messages6 = new Map([['a', 0]])

  const result1 = undefined
  const result2 = undefined
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
  }
  const result5 = {
    coefficients: [3, 1],
    messageIds: [1, 0],
    offset: 1,
    comparisonOperator: 0,
  }
  const result6 = {
    coefficients: [1],
    messageIds: [0],
    offset: -3,
    comparisonOperator: 1,
  }

  beforeAll(() => {
    constraintParser = new ConstraintParser();
  });

  describe('parseConstraints', () => {
    it('should parse undefined', () => {
      const result = constraintParser.parseConstraint(constraint1!, new Map());
      expect(result).toEqual(result1);
    });
    it('should parse text', () => {
      const result = constraintParser.parseConstraint(constraint2, new Map());
      expect(result).toEqual(result2);
    });
    it('should parse constraint without messageIds', () => {
      const result = constraintParser.parseConstraint(constraint3, new Map());
      expect(result).toEqual(result2);
    });
    it('should parse constraint3', () => {
      const result = constraintParser.parseConstraint(constraint3, messages3);
      expect(result).toEqual(result3);
    });
    it('should parse constraint4', () => {
      const result = constraintParser.parseConstraint(constraint4, messages4);
      expect(result).toEqual(result4);
    });
    it('should parse constraint5', () => {
      const result = constraintParser.parseConstraint(constraint5, messages5);
      expect(result).toEqual(result5);
    });
    it('should parse constraint6', () => {
      const result = constraintParser.parseConstraint(constraint6, messages6);
      expect(result).toEqual(result6);
    });
  });
});
