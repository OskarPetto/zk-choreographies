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
	event, err := controller.executionService.InstantiateModel(cmd)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	jsonEvent := InstatiatedModelEventToJson(event)
	c.IndentedJSON(http.StatusOK, jsonEvent)
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
	event, err := controller.executionService.ExecuteTransition(cmd)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	jsonEvent := ExecutedTransitionEventToJson(event)
	c.IndentedJSON(http.StatusOK, jsonEvent)
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
	event, err := controller.executionService.TerminateInstance(cmd)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	jsonEvent := TerminatedInstanceEventToJson(event)
	c.IndentedJSON(http.StatusOK, jsonEvent)
}

func (controller *ExecutionController) SendMessage(c *gin.Context) {
	var jsonCmd sendMessageCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	event, err := controller.executionService.SendMessage(cmd)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	jsonEvent := SentMessageEventToJson(event)
	c.IndentedJSON(http.StatusOK, jsonEvent)
}

func (controller *ExecutionController) ReceiveMessage(c *gin.Context) {
	var jsonCmd receiveMessageCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	event, err := controller.executionService.ReceiveMessage(cmd)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	jsonEvent := ReceivedMessageEventToJson(event)
	c.IndentedJSON(http.StatusOK, jsonEvent)
}
