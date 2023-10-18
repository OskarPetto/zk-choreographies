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

var signatureParameters = parameters.NewSignatureParameters()

var modelService = model.NewModelService()
var modelController = model.NewModelController(modelService)

var instanceService = instance.NewInstanceService(modelService)
var instanceController = instance.NewInstanceController(instanceService)

var messageService = message.NewMessageService(instanceService)
var messageController = message.NewMessageController(messageService)

var proofParameters = parameters.NewProverParameters()
var proverService = prover.NewProverService(proofParameters)

var executionService = execution.NewExecutionService(modelService, instanceService, messageService, proverService, signatureParameters)
var executionController = execution.NewExecutionController(executionService)

func main() {

	router := gin.Default()

	router.POST("/models", modelController.CreateModel)
	router.GET("/models", modelController.FindAllModels)
	router.GET("/models/:modelId", modelController.FindModelById)
	router.PUT("/models", modelController.ImportModel)

	router.GET("/models/:modelId/instances", instanceController.FindInstancesByModel)
	router.GET("/instances/:instanceId", instanceController.FindInstanceById)
	router.PUT("/instances", instanceController.ImportInstance)

	router.GET("/instances/:instanceId/messages", messageController.FindMessagesByInstance)
	router.GET("/messages/:messageId", messageController.FindMessageById)

	router.POST("/execution/instantiation", executionController.InstantiateModel)
	router.POST("/execution/transition", executionController.ExecuteTransition)
	router.POST("/execution/termination", executionController.TerminateInstance)
	router.POST("/execution/send", executionController.SendMessage)
	router.POST("/execution/receive", executionController.ReceiveMessage)

	router.Run("localhost:8080")
}
