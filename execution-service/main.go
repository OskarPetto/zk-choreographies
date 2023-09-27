package main

import (
	"execution-service/infrastructure/rest"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/models/:id/instances", rest.GetInstances)
	router.POST("/execution/instantiateModel", rest.InstantiateModel)
	router.POST("/execution/executeTransition", rest.ExecuteTransition)
	router.PUT("/proof/proveInstantiation", rest.ProveInstantiation)
	router.PUT("/proof/proveTransition", rest.ProveTransition)
	router.PUT("/proof/proveTermination", rest.ProveTermination)

	router.Run("localhost:8080")
}
