#!/bin/sh

bpmn_choreography=`cat ../bpmn-service/test/data/BikeRental_example.bpmn`
bpmn_choreography="${bpmn_choreography//\"/\\\"}"
bpmn_choreography=$(echo -e "$bpmn_choreography" | tr -d '\n')
transform_choreography_command="{\"xmlString\":\"${bpmn_choreography}\"}"

model=$(curl -d "$transform_choreography_command" -X POST http://localhost:3000/choreographies  -H 'Content-Type: application/json')
echo $model