//go:build wireinject
// +build wireinject

package main

import (
	"execution-service/adapter"
	"execution-service/execution"
	"execution-service/instance"
	"execution-service/model"
	"execution-service/proof"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var instanceProviders = wire.NewSet(instance.NewInstanceController, instance.NewInstanceService)
var providerBindings = wire.NewSet(adapter.NewModelAdapter, wire.Bind(new(model.ModelPort), new(*adapter.ModelAdapter)))
var modelProviders = wire.NewSet(proofProviders, model.NewModelController, model.NewModelService)
var executionProviders = wire.NewSet(providerBindings, execution.NewExecutionController, execution.InitializeExecutionService)
var proofProviders = wire.NewSet(providerBindings, proof.NewProofController, proof.InitializeProofService)

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
	router.PUT("/models", modelController.PutModel)
	router.GET("/models/:modelId", modelController.GetModel)

	router.GET("/models/:modelId/instances", instanceController.GetInstances)
	router.PUT("/instances", instanceController.PutInstance)
	router.GET("/instances/:instanceId", instanceController.GetInstance)

	router.POST("/models/:modelId/instances", executionController.InstantiateModel)
	router.POST("/models/:modelId/instances/:instanceId", executionController.ExecuteTransition)

	router.PUT("/proof/instantiation", proofController.ProveInstantiation)
	router.PUT("/proof/transition", proofController.ProveTransition)
	router.PUT("/proof/termination", proofController.ProveTermination)

	router.Run("localhost:8080")
}
