package main

import (
	"execution-service/execution"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/prover"

	"github.com/gin-gonic/gin"
)

var modelService = model.NewModelService()
var modelController = model.NewModelController(modelService)

var instanceService = instance.NewInstanceService()
var instanceController = instance.NewInstanceController(instanceService)

var messageService = message.NewMessageService()
var signatureParameters = parameters.NewSignatureParameters()
var proofParameters = parameters.NewProverParameters()
var proverService = prover.NewProverService(proofParameters)
var executionService = execution.NewExecutionService(modelService, instanceService, messageService, proverService, signatureParameters)
var executionController = execution.NewExecutionController(executionService)

func main() {

	router := gin.Default()
	router.POST("/models", modelController.CreateModel)

	router.GET("/models/:modelId", modelController.FindModelById)
	router.GET("/models/choreography/:choreographyId", modelController.FindModelsByChoreography)

	router.GET("/instances/:instanceId", instanceController.FindInstanceById)
	router.GET("/instances/model/:modelId/", instanceController.FindInstancesByModel)

	router.POST("/execution/instantiation", executionController.InstantiateModel)
	router.POST("/execution/transition", executionController.ExecuteTransition)
	router.POST("/execution/termination", executionController.TerminateInstance)

	router.Run("localhost:8080")
}
