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

var instanceService = instance.NewInstanceService()
var instanceController = instance.NewInstanceController(instanceService)

var modelService = model.NewModelService()
var modelController = model.NewModelController(modelService)

var messageService = message.NewMessageService()
var messageController = message.NewMessageController(messageService)

var proofParameters = parameters.NewProverParameters()
var proverService = prover.NewProverService(proofParameters)

var executionService = execution.NewExecutionService(modelService, instanceService, messageService, proverService, signatureParameters)
var executionController = execution.NewExecutionController(executionService)

var signatureController = signature.NewSignatureController(signatureParameters)

func main() {

	router := gin.Default()

	router.GET("/publicKeys", signatureController.GetPublicKeys)

	router.POST("/models", modelController.CreateModel)
	router.GET("/models", modelController.FindAllModels)
	router.GET("/models/:modelId", modelController.FindModelById)
	router.GET("/models/:modelId/instances", instanceController.FindInstancesByModel)

	router.GET("/instances/:instanceId", instanceController.FindInstanceById)
	router.GET("/instances/:instanceId/messages", messageController.FindMessagesByInstance)

	router.GET("/messages/:messageId", messageController.FindMessageById)

	router.POST("/execution/instantiateModel", executionController.InstantiateModel)
	router.POST("/execution/executeTransition", executionController.ExecuteTransition)
	router.POST("/execution/proveTermination", executionController.ProveTermination)
	router.POST("/execution/createInitiatingMessage", executionController.CreateInitiatingMessage)
	router.POST("/execution/receiveInitiatingMessage", executionController.ReceiveInitiatingMessage)
	router.POST("/execution/proveMessageExchange", executionController.ProveMessageExchange)

	router.Run("localhost:8080")
}
