package execution

import (
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
	var jsonCmd instantiateModelCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, err := controller.executionService.InstantiateModel(cmd)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	jsonResult := ToJson(result)
	c.IndentedJSON(http.StatusOK, jsonResult)
}

func (controller *ExecutionController) ExecuteTransition(c *gin.Context) {
	var jsonCmd executeTransitionCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, err := controller.executionService.ExecuteTransition(cmd)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	jsonResult := ToJson(result)
	c.IndentedJSON(http.StatusOK, jsonResult)
}

func (controller *ExecutionController) TerminateInstance(c *gin.Context) {
	var jsonCmd terminateInstanceCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, err := controller.executionService.TerminateInstance(cmd)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	jsonResult := ToJson(result)
	c.IndentedJSON(http.StatusOK, jsonResult)
}
