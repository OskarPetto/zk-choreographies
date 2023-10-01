#!/bin/bash

bpmn=$(< floor_choreography.bpmn)
bpmn=$(printf '%s' "$bpmn" | jq -sR .)

json="{\"bpmnString\":$bpmn}"

#modelId=$(curl -H "Content-Type:application/json" -X POST --data "$json" http://localhost:3000/choreographies)
#echo $modelId

modelId=GKFmv2yRkhAf9TcFyU9QpNaM1oP4Q1jweMWy0p9ZG+I=

curl http://localhost:8080/models/$modelId