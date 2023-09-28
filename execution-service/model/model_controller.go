package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModelController struct {
	modelService ModelService
}

func NewModelController(modelService ModelService) ModelController {
	return ModelController{
		modelService: modelService,
	}
}

func (controller *ModelController) GetModel(c *gin.Context) {
	modelId := c.Param("modelId")
	model, err := controller.modelService.FindModelById(modelId)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, ToJson(model))
}

func (controller *ModelController) PutModel(c *gin.Context) {
	var modelJson ModelJson
	if err := c.BindJSON(&modelJson); err != nil {
		return
	}
	model, err := modelJson.ToModel()
	if err != nil {
		return
	}
	controller.modelService.SaveModel(model)
	c.Status(http.StatusOK)
}
