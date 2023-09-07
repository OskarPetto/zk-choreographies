import { Injectable } from "@nestjs/common";
import { Element, ElementType, FlowId, Model, createModelId } from "../domain/model";

const { XMLParser } = require("fast-xml-parser");

@Injectable()
export class BpmnService {
    private options = {
        attributeNamePrefix: "",
        ignoreAttributes: false
    };
    private parser = new XMLParser(this.options);

    parseModel(bpmnString: string): Model {
        const definitions = this.parser.parse(bpmnString)['bpmn:definitions'];

        const process = definitions['bpmn:process'];
        const startEvent: Element = this.parseStartEvent(process);
        const endEvents: Element[] = this.parseEndEvents(process);
        const tasks: Element[] = this.parseTasks(process);
        const exclusiveGateways: Element[] = this.parseExclusiveGateways(process);
        const parallelGateways: Element[] = this.parseParallelGateways(process);
        const flows: FlowId[] = this.parseSequenceFlows(process);
        const elements = [startEvent, ...endEvents, ...tasks, ...exclusiveGateways, ...parallelGateways];
        return {
            id: createModelId(),
            flows,
            elements: new Map(elements.map(e => [e.id, e]))
        };
    }

    private parseStartEvent(process: any): Element {
        const startEvent = process['bpmn:startEvent'];
        return {
            id: startEvent.id,
            type: ElementType.START,
            name: startEvent.name,
            incomingFlows: [],
            outgoingFlows: [startEvent['bpmn:outgoing']]
        };
    }

    private parseEndEvents(process: any): Element[] {
        const endEvents: any[] = [].concat(process['bpmn:endEvent']);
        return endEvents.map((endEvent: any) => ({
            id: endEvent.id,
            type: ElementType.END,
            name: endEvent.name,
            incomingFlows: [endEvent['bpmn:incoming']],
            outgoingFlows: []
        }));
    }

    private parseTasks(process: any): Element[] {
        const tasks: any[] = [].concat(process['bpmn:task']);
        return tasks.map((task: any) => ({
            id: task.id,
            type: ElementType.TASK,
            name: task.name,
            incomingFlows: [task['bpmn:incoming']],
            outgoingFlows: [task['bpmn:outgoing']]
        }));
    }

    private parseExclusiveGateways(process: any): Element[] {
        const exclusiveGateways = process['bpmn:exclusiveGateway'];
        return exclusiveGateways.flatMap((exclusiveGateway: any) => this.parseExclusiveGateway(exclusiveGateway));
    }

    private parseExclusiveGateway(exclusiveGateway: any): Element[] {
        const incomingFlows = exclusiveGateway['bpmn:incoming'];
        const outgoingFlows = exclusiveGateway['bpmn:outgoing'];

        if (Array.isArray(incomingFlows)) {
            return incomingFlows.map(incomingFlow => ({
                id: `${exclusiveGateway.id}_${incomingFlow}`,
                type: ElementType.XOR_JOIN,
                name: incomingFlow.name,
                incomingFlows: [incomingFlow],
                outgoingFlows: [outgoingFlows]
            }));
        } else if (Array.isArray(outgoingFlows)) {
            return outgoingFlows.map(outgoingFlow => ({
                id: `${exclusiveGateway.id}_${outgoingFlow}`,
                type: ElementType.XOR_SPLIT,
                name: outgoingFlow.name,
                incomingFlows: [incomingFlows],
                outgoingFlows: [outgoingFlow]
            }));
        } else {
            return [];
        }
    }

    private parseParallelGateways(process: any): Element[] {
        const parallelGateways = process['bpmn:parallelGateway'];
        return parallelGateways.map((parallelGateway: any) => {
            const incomingFlows = parallelGateway['bpmn:incoming'];
            const outgoingFlows = parallelGateway['bpmn:outgoing'];
            return {
                id: parallelGateway.id,
                type: Array.isArray(incomingFlows) ? ElementType.AND_JOIN : ElementType.AND_SPLIT,
                incomingFlows: [].concat(incomingFlows),
                outgoingFlows: [].concat(outgoingFlows)
            }
        });
    }

    private parseSequenceFlows(process: any): FlowId[] {
        const sequenceFlows = [].concat(process['bpmn:sequenceFlow']);
        return sequenceFlows.map((sequenceFlow: any) => sequenceFlow.id as FlowId);
    }
}