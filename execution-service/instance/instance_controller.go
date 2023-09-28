package instance

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type InstanceController struct {
	instanceService InstanceService
}

func NewInstanceController(instanceService InstanceService) InstanceController {
	return InstanceController{
		instanceService: instanceService,
	}
}

func (controller *InstanceController) FindInstanceById(c *gin.Context) {
	instanceId := c.Param("instanceId")
	instance, err := controller.instanceService.FindInstanceById(instanceId)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, ToJson(instance))
}

func (controller *InstanceController) ImportInstance(c *gin.Context) {
	var instanceJson InstanceJson
	if err := c.BindJSON(&instanceJson); err != nil {
		return
	}
	instance, err := instanceJson.ToInstance()
	if err != nil {
		return
	}
	controller.instanceService.ImportInstance(instance)
	c.Status(http.StatusOK)
}

func (controller *InstanceController) FindInstancesByModel(c *gin.Context) {
	modelId := c.Param("modelId")
	instances := controller.instanceService.FindInstancesByModel(modelId)
	jsonInstances := make([]InstanceJson, len(instances))
	for i, instance := range instances {
		jsonInstances[i] = ToJson(instance)
	}
	c.IndentedJSON(http.StatusOK, jsonInstances)
}
