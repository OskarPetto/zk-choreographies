package signature

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignatureController struct {
	signatureService SignatureService
}

func NewSignatureController(signatureService SignatureService) SignatureController {
	return SignatureController{
		signatureService: signatureService,
	}
}

func (controller *SignatureController) FindSignatureByInstance(c *gin.Context) {
	instanceId := c.Param("instanceId")
	signature, err := controller.signatureService.FindSignatureByInstance(instanceId)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, ToJson(signature))
}

func (controller *SignatureController) ImportSignature(c *gin.Context) {
	var signatureJson SignatureJson
	if err := c.BindJSON(&signatureJson); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	signature, err := signatureJson.ToSignature()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = controller.signatureService.ImportSignature(signature)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}
