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

func (controller *ModelController) FindModelById(c *gin.Context) {
	modelId := c.Param("modelId")
	model, err := controller.modelService.FindModelById(modelId)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, ToJson(model))
}

func (controller *ModelController) FindAllModels(c *gin.Context) {
	models := controller.modelService.FindAllModels()
	jsonModels := make([]ModelJson, len(models))
	for i, instance := range models {
		jsonModels[i] = ToJson(instance)
	}
	c.IndentedJSON(http.StatusOK, jsonModels)
}
