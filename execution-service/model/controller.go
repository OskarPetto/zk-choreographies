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

func (controller *ModelController) CreateModel(c *gin.Context) {
	var modelJson ModelJson
	if err := c.BindJSON(&modelJson); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	model, err := modelJson.ToModel()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result := controller.modelService.CreateModel(model)
	c.IndentedJSON(http.StatusOK, ToJson(result))
}

func (controller *ModelController) FindModelsByChoreography(c *gin.Context) {
	choreographyId := c.Query("choreographyId")
	if choreographyId == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	models := controller.modelService.FindModelsByChoreography(choreographyId)
	jsonModels := make([]ModelJson, len(models))
	for i, instance := range models {
		jsonModels[i] = ToJson(instance)
	}
	c.IndentedJSON(http.StatusOK, jsonModels)
}
