#!/bin/sh

weber_choreography=`cat ../bpmn-service/test/data/weber_choreography.bpmn`
weber_choreography="${weber_choreography//\"/\\\"}"
weber_choreography=$(echo -e "$weber_choreography" | tr -d '\n')
transform_choreography_command="{\"xmlString\":\"${weber_choreography}\"}"

model=$(curl -d "$transform_choreography_command" -X POST http://localhost:3000/choreographies  -H 'Content-Type: application/json')
echo $model