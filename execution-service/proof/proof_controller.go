package proof

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProofController struct {
	proofService ProofService
}

func NewProofController(proofService ProofService) ProofController {
	return ProofController{
		proofService: proofService,
	}
}

func (controller *ProofController) ProveInstantiation(c *gin.Context) {
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

func (controller *ProofController) ProveTransition(c *gin.Context) {
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

func (controller *ProofController) ProveTermination(c *gin.Context) {
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
