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

var providerBindings = wire.NewSet(adapter.NewModelAdapter, wire.Bind(new(model.ModelPort), new(*adapter.ModelAdapter)))
var instanceProviders = wire.NewSet(instance.NewInstanceController, instance.NewInstanceService)
var executionProviders = wire.NewSet(providerBindings, execution.NewExecutionController, execution.InitializeExecutionService)
var proofProviders = wire.NewSet(providerBindings, proof.NewProofController, proof.InitializeProofService)

func InitializeInstanceController() instance.InstanceController {
	wire.Build(instanceProviders)
	return instance.InstanceController{}
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
	executionController := InitializeExecutionController()
	proofController := InitializeProofController()

	router := gin.Default()
	router.GET("/models/:modelId/instances", instanceController.GetInstances)
	router.POST("/models/:modelId/instances", executionController.InstantiateModel)
	router.POST("/models/:modelId/instances/:instanceId", executionController.ExecuteTransition)
	router.PUT("/instantiation-proof", proofController.ProveInstantiation)
	router.PUT("/transition-proof", proofController.ProveTransition)
	router.PUT("/termination-proof", proofController.ProveTermination)

	router.Run("localhost:8080")
}
