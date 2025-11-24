package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EventUserController struct {
	Service *service.EventUserService
}

func NewEventUserController(service *service.EventUserService) *EventUserController {
	return &EventUserController{Service: service}
}

func (c *EventUserController) RegisterEventUser(ctx *gin.Context) {
	var input model.EventUser
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := c.Service.CreateEventUser(input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": created})
}
