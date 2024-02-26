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
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	event, err := controller.executionService.InstantiateModel(cmd)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusForbidden)
		return
	}
	jsonEvent := InstatiatedModelEventToJson(event)
	c.IndentedJSON(http.StatusOK, jsonEvent)
}

func (controller *ExecutionController) ExecuteTransition(c *gin.Context) {
	var jsonCmd executeTransitionCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	event, err := controller.executionService.ExecuteTransition(cmd)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusForbidden)
		return
	}
	jsonEvent := ExecutedTransitionEventToJson(event)
	c.IndentedJSON(http.StatusOK, jsonEvent)
}

func (controller *ExecutionController) ProveTermination(c *gin.Context) {
	var jsonCmd proveTerminationCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	event, err := controller.executionService.ProveTermination(cmd)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusForbidden)
		return
	}
	jsonEvent := TerminatedInstanceEventToJson(event)
	c.IndentedJSON(http.StatusOK, jsonEvent)
}

func (controller *ExecutionController) CreateInitiatingMessage(c *gin.Context) {
	var jsonCmd createInitiatingMessageCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	event, err := controller.executionService.CreateInitiatingMessage(cmd)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusForbidden)
		return
	}
	jsonEvent := CreatedInitiatingMessageEventToJson(event)
	c.IndentedJSON(http.StatusOK, jsonEvent)
}

func (controller *ExecutionController) ReceiveInitiatingMessage(c *gin.Context) {
	var jsonCmd receiveInitiatingMessageCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	event, err := controller.executionService.ReceiveInitiatingMessage(cmd)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusForbidden)
		return
	}
	jsonEvent := ReceivedInitiatingMessageEventToJson(event)
	c.IndentedJSON(http.StatusOK, jsonEvent)
}

func (controller *ExecutionController) ProveMessageExchange(c *gin.Context) {
	var jsonCmd proveMessageExchangeCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToExecutionCommand()
	if err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	event, err := controller.executionService.ProveMessageExchange(cmd)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusForbidden)
		return
	}
	jsonEvent := ProvedMessageExchangeEventToJson(event)
	c.IndentedJSON(http.StatusOK, jsonEvent)
}
