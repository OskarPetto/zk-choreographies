package main

import (
	"execution-service/execution"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/prover"
	"execution-service/state"

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

var stateService = state.NewStateService(modelService, instanceService, messageService, signatureParameters)
var stateController = state.NewStateController(stateService)

func main() {

	router := gin.Default()
	router.POST("/models", modelController.CreateModel)

	router.GET("/models/:modelId", modelController.FindModelById)
	router.GET("/models", modelController.FindModelsByChoreography)

	router.GET("/instances/:instanceId", instanceController.FindInstanceById)
	router.GET("/instances", instanceController.FindInstancesByModel)

	router.POST("/instantiation", executionController.InstantiateModel)
	router.POST("/transition", executionController.ExecuteTransition)
	router.POST("/termination", executionController.TerminateInstance)

	router.PUT("/state", stateController.ImportState)

	router.Run("localhost:8080")
}
