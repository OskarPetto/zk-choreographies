package main

import (
	"execution-service/execution"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/prover"
	"execution-service/signature"

	"github.com/gin-gonic/gin"
)

var signatureParameters = parameters.NewSignatureParameters()

var modelService = model.NewModelService()
var modelController = model.NewModelController(modelService)

var instanceService = instance.NewInstanceService(modelService)
var instanceController = instance.NewInstanceController(instanceService)

var messageService = message.NewMessageService(modelService, instanceService, signatureParameters)
var messageController = message.NewMessageController(messageService)

var signatureService = signature.NewSignatureService(instanceService)
var signatureController = signature.NewSignatureController(signatureService)

var proofParameters = parameters.NewProverParameters()
var proverService = prover.NewProverService(proofParameters)

var executionService = execution.NewExecutionService(modelService, instanceService, messageService, proverService, signatureParameters, signatureService)
var executionController = execution.NewExecutionController(executionService)

func main() {

	router := gin.Default()

	router.GET("/models", modelController.FindAllModels)
	router.GET("/models/:modelId", modelController.FindModelById)
	router.PUT("/models", modelController.ImportModel)
	router.POST("/models", modelController.CreateModel)

	router.GET("/models/:modelId/instances", instanceController.FindInstancesByModel)
	router.GET("/instances/:instanceId", instanceController.FindInstanceById)
	router.PUT("/instances", instanceController.ImportInstance)

	router.GET("/instances/:instanceId/messages", messageController.FindMessagesByInstance)
	router.GET("/messages/:messageId", messageController.FindMessageById)
	router.PUT("/messages", messageController.ImportMessage)
	router.POST("/messages", messageController.CreateMessage)

	router.GET("/instances/:instanceId/signatures", signatureController.FindSignatureByInstance)
	router.PUT("/signatures", signatureController.ImportSignature)

	router.POST("/instantiation", executionController.InstantiateModel)
	router.POST("/transition", executionController.ExecuteTransition)
	router.POST("/termination", executionController.TerminateInstance)

	router.Run("localhost:8080")
}
