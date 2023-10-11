package state

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StateController struct {
	stateService StateService
}

func NewStateController(stateService StateService) StateController {
	return StateController{
		stateService: stateService,
	}
}

func (controller *StateController) ImportState(c *gin.Context) {
	var cmdJson ImportStateCommandJson
	if err := c.BindJSON(&cmdJson); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := cmdJson.ToStateCommand()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = controller.stateService.ImportState(cmd)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
