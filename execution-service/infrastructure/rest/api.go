package rest

import (
	"execution-service/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

var modelService = domain.NewModelService()

func GetModels(c *gin.Context) {
	models := modelService.FindAllModels()
	c.IndentedJSON(http.StatusOK, models)
}
