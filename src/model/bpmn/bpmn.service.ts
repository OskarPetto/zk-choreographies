import { Injectable } from "@nestjs/common";
import { ElementType, Model, createModelId } from "../domain/model";
import { logObject } from "src/utils/utils";

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
        logObject(process);
        const model = {
            id: createModelId(),
            flows: [],
            elements: new Map()
        };
        this.parseStartEvent(process, model);
        this.parseEndEvents(process, model);
        this.parseTasks(process, model);
        this.parseExclusiveGateways(process, model);
        this.parseParallelGateways(process, model);
        return model;
    }

    private parseStartEvent(process: any, model: Model) {
    }

    private parseEndEvents(process: any, model: Model) {

    }

    private parseTasks(process: any, model: Model) {

    }

    private parseExclusiveGateways(process: any, model: Model) {

    }

    private parseParallelGateways(process: any, model: Model) {

    }
}