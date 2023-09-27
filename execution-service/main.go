package main

import (
	"execution-service/infrastructure/rest"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/models/:modelId/instances", rest.GetInstances)
	router.POST("/models/:modelId/instantiate", rest.InstantiateModel)
	router.POST("/models/:modelId/instances/:instanceId", rest.ExecuteTransition)
	router.PUT("/proof/proveInstantiation", rest.ProveInstantiation)
	router.PUT("/proof/proveTransition", rest.ProveTransition)
	router.PUT("/proof/proveTermination", rest.ProveTermination)

	router.Run("localhost:8080")
}
