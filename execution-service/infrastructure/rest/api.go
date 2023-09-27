package rest

import (
	"execution-service/domain"
	"execution-service/execution"
	"execution-service/infrastructure/json"
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
	var jsonCmd json.InstantiateModelCommand
	if err := c.BindJSON(&jsonCmd); err != nil {
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		return
	}
	result, err := executionService.InstantiateModel(cmd)
	if err != nil {
		return
	}
	jsonResult := json.FromExecutionResult(result)
	c.IndentedJSON(http.StatusOK, jsonResult)
}

func ExecuteTransition(c *gin.Context) {
	var jsonCmd json.ExecuteTransitionCommand
	if err := c.BindJSON(&jsonCmd); err != nil {
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		return
	}
	result, err := executionService.ExecuteTransition(cmd)
	if err != nil {
		return
	}
	jsonResult := json.FromExecutionResult(result)
	c.IndentedJSON(http.StatusOK, jsonResult)
}

func ProveTermination(c *gin.Context) {
	var jsonCmd json.ProveTerminationCommand
	if err := c.BindJSON(&jsonCmd); err != nil {
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		return
	}
	result, err := executionService.ProveTermination(cmd)
	if err != nil {
		return
	}
	jsonResult := json.FromExecutionResult(result)
	c.IndentedJSON(http.StatusOK, jsonResult)
}
