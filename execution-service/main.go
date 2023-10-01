package main

import (
	"execution-service/execution"
	"execution-service/instance"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/proof"

	"github.com/gin-gonic/gin"
)

var instanceService = instance.NewInstanceService()
var instanceController = instance.NewInstanceController(instanceService)
var modelService = model.NewModelService()
var modelController = model.NewModelController(modelService)

var executionService = execution.NewExecutionService(instanceService, modelService)
var executionController = execution.NewExecutionController(executionService)

var signatureParameters = parameters.NewSignatureParameters()
var proofParameters = parameters.NewProofParameters()
var proofService = proof.NewProofService(proofParameters, signatureParameters, instanceService, modelService)
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

	router.POST("/instances/instantiation", executionController.InstantiateModel)
	router.POST("/instances/transition", executionController.ExecuteTransition)

	router.PUT("/proof/instantiation", proofController.ProveInstantiation)
	router.PUT("/proof/transition", proofController.ProveTransition)
	router.PUT("/proof/termination", proofController.ProveTermination)

	router.Run("localhost:8080")
}
