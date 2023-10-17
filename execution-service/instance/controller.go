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

func (controller *InstanceController) ImportInstance(c *gin.Context) {
	var instanceJson InstanceJson
	if err := c.BindJSON(&instanceJson); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	instance, err := instanceJson.ToInstance()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = controller.instanceService.ImportInstance(instance)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}
