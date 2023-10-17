package message

import (
	"execution-service/signature"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	messageService MessageService
}

func NewMessageController(messageService MessageService) MessageController {
	return MessageController{
		messageService: messageService,
	}
}

func (controller *MessageController) FindMessagesByInstance(c *gin.Context) {
	instanceId := c.Param("instanceId")
	instances := controller.messageService.FindMessagesByInstance(instanceId)
	jsonMessages := make([]MessageJson, len(instances))
	for i, instance := range instances {
		jsonMessages[i] = MessageToJson(instance)
	}
	c.IndentedJSON(http.StatusOK, jsonMessages)
}

func (controller *MessageController) FindMessageById(c *gin.Context) {
	messageId := c.Param("messageId")
	message, err := controller.messageService.FindMessageById(messageId)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, MessageToJson(message))
}

func (controller *MessageController) ImportMessage(c *gin.Context) {
	var cmdJson ImportMessageCommandJson
	if err := c.BindJSON(&cmdJson); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := cmdJson.ToMessageCommand()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, err := controller.messageService.ImportMessage(cmd)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.IndentedJSON(http.StatusOK, signature.ToJson(result))
}

func (controller *MessageController) CreateMessage(c *gin.Context) {
	var cmdJson CreateMessageCommandJson
	if err := c.BindJSON(&cmdJson); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	cmd, err := cmdJson.ToMessageCommand()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, err := controller.messageService.CreateMessage(cmd)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.IndentedJSON(http.StatusOK, CreateMessageResultToJson(result))
}
