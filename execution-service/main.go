//go:build wireinject
// +build wireinject

package main

import (
	"execution-service/execution"
	"execution-service/instance"
	"execution-service/model"
	"execution-service/proof"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var instanceProviders = wire.NewSet(instance.NewInstanceController, instance.NewInstanceService)
var modelProviders = wire.NewSet(model.NewModelController, model.NewModelService)
var executionProviders = wire.NewSet(execution.NewExecutionController, execution.InitializeExecutionService)
var proofProviders = wire.NewSet(proof.NewProofController, proof.InitializeProofService)

func InitializeInstanceController() instance.InstanceController {
	wire.Build(instanceProviders)
	return instance.InstanceController{}
}

func InitializeModelController() model.ModelController {
	wire.Build(modelProviders)
	return model.ModelController{}
}

func InitializeExecutionController() execution.ExecutionController {
	wire.Build(executionProviders)
	return execution.ExecutionController{}
}

func InitializeProofController() proof.ProofController {
	wire.Build(proofProviders)
	return proof.ProofController{}
}

func main() {

	instanceController := InitializeInstanceController()
	modelController := InitializeModelController()
	executionController := InitializeExecutionController()
	proofController := InitializeProofController()

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
