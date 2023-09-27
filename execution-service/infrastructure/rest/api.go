package rest

import (
	"execution-service/domain"
	"execution-service/execution"
	"net/http"

	"github.com/gin-gonic/gin"
)

var instanceService = domain.NewInstanceService()
var executionService = execution.NewExecutionService()

func GetInstances(c *gin.Context) {
	id := c.Param("id")
	instances := instanceService.FindInstancesByModel(id)
	c.IndentedJSON(http.StatusOK, instances)
}

func InstantiateModel(c *gin.Context) {
	var cmd execution.InstantiateModelCommand
	if err := c.BindJSON(&cmd); err != nil {
		return
	}
	result, err := executionService.InstantiateModel(cmd)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}

func ExecuteTransition(c *gin.Context) {
	var cmd execution.ExecuteTransitionCommand
	if err := c.BindJSON(&cmd); err != nil {
		return
	}
	result, err := executionService.ExecuteTransition(cmd)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}

func ProveTermination(c *gin.Context) {
	var cmd execution.ProveTerminationCommand
	if err := c.BindJSON(&cmd); err != nil {
		return
	}
	result, err := executionService.ProveTermination(cmd)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}
