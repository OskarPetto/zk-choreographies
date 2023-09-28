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

func (controller *InstanceController) GetInstances(c *gin.Context) {
	modelId := c.Param("modelId")
	instances := controller.instanceService.FindInstancesByModel(modelId)
	c.IndentedJSON(http.StatusOK, instances)
}
