package execution

import (
	"execution-service/instance"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExecutionController struct {
	executionService ExecutionService
}

func NewExecutionController(executionService ExecutionService) ExecutionController {
	return ExecutionController{
		executionService: executionService,
	}
}

func (controller *ExecutionController) InstantiateModel(c *gin.Context) {
	modelId := c.Param("modelId")
	var jsonCmd InstantiateModelCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand(modelId)
	if err != nil {
		return
	}
	result, err := controller.executionService.InstantiateModel(cmd)
	if err != nil {
		return
	}
	jsonResult := instance.ToJson(result)
	c.IndentedJSON(http.StatusOK, jsonResult)
}

func (controller *ExecutionController) ExecuteTransition(c *gin.Context) {
	modelId := c.Param("modelId")
	instanceId := c.Param("instanceId")

	var jsonCmd ExecuteTransitionCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand(modelId, instanceId)
	if err != nil {
		return
	}
	result, err := controller.executionService.ExecuteTransition(cmd)
	if err != nil {
		return
	}
	jsonResult := instance.ToJson(result)
	c.IndentedJSON(http.StatusOK, jsonResult)
}
