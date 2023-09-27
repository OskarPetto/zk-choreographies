package main

import (
	"execution-service/infrastructure/rest"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/models/:modelId/instances", rest.GetInstances)
	router.POST("/models/:modelId/instances", rest.InstantiateModel)
	router.POST("/models/:modelId/instances/:instanceId", rest.ExecuteTransition)
	router.PUT("/instantiation-proof", rest.ProveInstantiation)
	router.PUT("/transition-proof", rest.ProveTransition)
	router.PUT("/termination-proof", rest.ProveTermination)

	router.Run("localhost:8080")
}
