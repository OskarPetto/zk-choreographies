package signature

import (
	"execution-service/domain"
	"execution-service/parameters"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignatureController struct {
	signatureParameters parameters.SignatureParameters
}

func NewSignatureController(signatureParameters parameters.SignatureParameters) SignatureController {
	return SignatureController{
		signatureParameters: signatureParameters,
	}
}

func (controller *SignatureController) GetPublicKeys(c *gin.Context) {
	publicKeys := controller.signatureParameters.GetPublicKeys(domain.IdentityCount)
	publicKeyStrings := toStringArray(publicKeys)
	c.IndentedJSON(http.StatusOK, publicKeyStrings)
}
