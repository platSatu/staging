package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketBlastController struct {
	Service *service.TicketBlastService
}

func NewTicketBlastController(s *service.TicketBlastService) *TicketBlastController {
	return &TicketBlastController{Service: s}
}

// CreateTicketBlast
func (tc *TicketBlastController) CreateTicketBlast(c *gin.Context) {
	var ticketBlast model.TicketBlast
	if err := c.ShouldBindJSON(&ticketBlast); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketBlast(&ticketBlast); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketBlast,
	})
}

// GetAllTicketBlasts
func (tc *TicketBlastController) GetAllTicketBlasts(c *gin.Context) {
	ticketBlasts, err := tc.Service.GetAllTicketBlasts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketBlasts,
	})
}

// GetTicketBlastByID
func (tc *TicketBlastController) GetTicketBlastByID(c *gin.Context) {
	id := c.Param("id")
	ticketBlast, err := tc.Service.GetTicketBlastByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketBlast == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket blast not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketBlast,
	})
}

// UpdateTicketBlast
func (tc *TicketBlastController) UpdateTicketBlast(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketBlast
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTicketBlast(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketBlast, err := tc.Service.GetTicketBlastByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket blast",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketBlast,
	})
}

// DeleteTicketBlast
func (tc *TicketBlastController) DeleteTicketBlast(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketBlast(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket blast deleted",
	})
}
