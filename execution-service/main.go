package main

import (
	"execution-service/infrastructure/rest"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/models", rest.GetModels)

	router.Run("localhost:8080")
}
