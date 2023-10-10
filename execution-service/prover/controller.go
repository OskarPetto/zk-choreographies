package prover

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProverController struct {
	proofService ProverService
}

func NewProverController(proofService ProverService) ProverController {
	return ProverController{
		proofService: proofService,
	}
}

func (controller *ProverController) ProveInstantiation(c *gin.Context) {
	var jsonCmd ProveInstantiationCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToProofCommand()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, err := controller.proofService.ProveInstantiation(cmd)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	jsonResult := result.ToJson()
	c.IndentedJSON(http.StatusOK, jsonResult)
}

func (controller *ProverController) ProveTransition(c *gin.Context) {
	var jsonCmd ProveTransitionCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToProofCommand()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, err := controller.proofService.ProveTransition(cmd)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	jsonResult := result.ToJson()
	c.IndentedJSON(http.StatusOK, jsonResult)
}

func (controller *ProverController) ProveTermination(c *gin.Context) {
	var jsonCmd ProveTerminationCommandJson
	if err := c.BindJSON(&jsonCmd); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := jsonCmd.ToProofCommand()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, err := controller.proofService.ProveTermination(cmd)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	jsonResult := result.ToJson()
	c.IndentedJSON(http.StatusOK, jsonResult)
}
