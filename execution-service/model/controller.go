package model

import (
	"execution-service/hash"
	"fmt"
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

func (controller *ModelController) ImportModel(c *gin.Context) {
	var cmdJson ImportModelCommandJson
	if err := c.BindJSON(&cmdJson); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := cmdJson.ToModelCommand()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = controller.modelService.ImportModel(cmd)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

func (controller *ModelController) CreateModel(c *gin.Context) {
	var modelJson ModelJson
	if err := c.BindJSON(&modelJson); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	model, err := modelJson.ToModel()
	if err != nil {
		fmt.Printf("Error creating model: %+v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	result := controller.modelService.CreateModel(model)
	c.IndentedJSON(http.StatusOK, hash.ToJson(result))
}
