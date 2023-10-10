package main

import (
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/proof"

	"github.com/gin-gonic/gin"
)

var modelService = model.NewModelService()
var modelController = model.NewModelController(modelService)

var messageService = message.NewMessageService()
var instanceService = instance.NewInstanceService(modelService, messageService)
var instanceController = instance.NewInstanceController(instanceService)

var signatureParameters = parameters.NewSignatureParameters()
var proofParameters = parameters.NewProofParameters()
var proofService = proof.NewProofService(proofParameters, signatureParameters, instanceService)
var proofController = proof.NewProofController(proofService)

func main() {

	router := gin.Default()
	router.POST("/models", modelController.CreateModel)

	router.GET("/models/:modelId", modelController.FindModelById)
	router.GET("/models/choreography/:choreographyId", modelController.FindModelsByChoreography)
	router.PUT("/models", modelController.ImportModel)

	router.GET("/instances/:instanceId", instanceController.FindInstanceById)
	router.GET("/instances/model/:modelId/", instanceController.FindInstancesByModel)
	router.PUT("/instances", instanceController.ImportInstance)

	router.POST("/instances/instantiation", instanceController.InstantiateModel)
	router.POST("/instances/transition", instanceController.ExecuteTransition)

	router.PUT("/proof/instantiation", proofController.ProveInstantiation)
	router.PUT("/proof/transition", proofController.ProveTransition)
	router.PUT("/proof/termination", proofController.ProveTermination)

	router.Run("localhost:8080")
}
