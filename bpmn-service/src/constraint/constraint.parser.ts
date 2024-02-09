import { Injectable } from '@nestjs/common';
import { Constraint } from '../domain/constraint';
import { MessageId } from 'src/domain/model';
import { BpmnMessageId } from 'src/domain/choreography';
/* eslint @typescript-eslint/no-var-requires: "off" */
const esprima = require('esprima');

interface ConstraintPart {
  constant: number;
  messageId?: number;
}

@Injectable()
export class ConstraintParser {
  parseConstraint(
    constraint: string | undefined,
    messageIds: Map<BpmnMessageId, MessageId>,
  ): Constraint | undefined {
    if (constraint === undefined) {
      return undefined;
    }
    try {
      const script = esprima.parseScript(constraint);
      if (script.body.length != 1) {
        throw Error('not implemented');
      }
      const expression = script.body[0].expression;
      return this.parseBooleanExpression(expression, messageIds);
    } catch (e) {
      console.log(
        `parsing constraint '${constraint}' resulted in error: ${e.message}`,
      );
    }
    return undefined;
  }

  private parseBooleanExpression(
    expression: any,
    messageIds: Map<BpmnMessageId, MessageId>,
  ): Constraint {
    if (expression.type !== 'BinaryExpression') {
      throw Error('not implemented');
    }
    const comparisonOperator = this.parseComparisonOperator(
      expression.operator,
    );
    const leftParts: ConstraintPart[] = this.parseExpression(
      expression.left,
      messageIds,
      1,
    );
    const rightParts: ConstraintPart[] = this.parseExpression(
      expression.right,
      messageIds,
      -1,
    );
    const combinedParts = leftParts.concat(rightParts);
    const resultingParts: ConstraintPart[] = [];
    let c = 0;
    for (const part of combinedParts) {
      if (part.messageId !== undefined) {
        resultingParts.push(part);
      } else {
        c += part.constant;
      }
    }
    resultingParts.push({ constant: c });

    const offset = resultingParts[resultingParts.length - 1].constant;
    const coefficients = [];
    const messages = [];

    for (let i = 0; i < resultingParts.length - 1; i++) {
      coefficients.push(resultingParts[i].constant);
      messages.push(resultingParts[i].messageId ?? 0);
    }

    return {
      coefficients,
      messageIds: messages,
      offset,
      comparisonOperator: comparisonOperator,
    };
  }

  private parseExpression(
    expression: any,
    messageIds: Map<BpmnMessageId, MessageId>,
    factor: number,
  ): ConstraintPart[] {
    if (expression.type == 'Literal') {
      return [
        {
          constant: expression.value * factor,
        },
      ];
    } else if (expression.type == 'Identifier') {
      let messageId;
      if (!messageIds.has(expression.name)) {
        throw Error(`message identifier '${expression.name}' not known`);
      } else {
        messageId = messageIds.get(expression.name);
      }
      return [
        {
          constant: factor,
          messageId: messageId,
        },
      ];
    } else if (expression.type == 'UnaryExpression') {
      if (expression.operator == '-') {
        return this.parseExpression(
          expression.argument,
          messageIds,
          factor * -1,
        );
      } else {
        throw Error('not implemented');
      }
    } else if (expression.type == 'BinaryExpression') {
      if (expression.operator == '*') {
        const leftParts = this.parseExpression(
          expression.left,
          messageIds,
          factor,
        );
        const rightParts = this.parseExpression(
          expression.right,
          messageIds,
          factor,
        );
        if (leftParts.length > 1 || rightParts.length > 1) {
          throw Error('not implemented');
        }
        if (
          leftParts[0].messageId !== undefined &&
          rightParts[0].messageId !== undefined
        ) {
          throw Error('not implemented');
        }
        if (
          leftParts[0].messageId === undefined &&
          rightParts[0].messageId === undefined
        ) {
          throw Error('not implemented');
        }
        const literal =
          leftParts[0].messageId === undefined
            ? leftParts[0].constant
            : rightParts[0].constant;
        const variable =
          leftParts[0].messageId !== undefined
            ? leftParts[0].messageId
            : rightParts[0].messageId;
        return [
          {
            constant: literal,
            messageId: variable,
          },
        ];
      } else if (expression.operator == '+') {
        const leftParts = this.parseExpression(
          expression.left,
          messageIds,
          factor,
        );
        const rightParts = this.parseExpression(
          expression.right,
          messageIds,
          factor,
        );
        return leftParts.concat(rightParts);
      } else if (expression.operator == '-') {
        const leftParts = this.parseExpression(
          expression.left,
          messageIds,
          factor,
        );
        const rightParts = this.parseExpression(
          expression.right,
          messageIds,
          -1 * factor,
        );
        return leftParts.concat(rightParts);
      } else {
        throw Error('not implemented');
      }
    } else {
      throw Error('not implemented');
    }
  }

  private parseComparisonOperator(operator: string): number {
    switch (operator) {
      case '==':
        return 0;
      case '>':
        return 1;
      case '<':
        return 2;
      case '>=':
        return 3;
      case '<=':
        return 4;
    }
    throw Error('not implemented');
  }
}
