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
		c.Error(err)
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, ToJson(instance))
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
