package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketEventController struct {
	Service *service.TicketEventService
}

func NewTicketEventController(s *service.TicketEventService) *TicketEventController {
	return &TicketEventController{Service: s}
}

func (tc *TicketEventController) GetAllTicketEventsPublic(c *gin.Context) {
	ticketEvents, err := tc.Service.GetAllTicketEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketEvents,
	})
}

// CreateTicketEvent
func (tc *TicketEventController) CreateTicketEvent(c *gin.Context) {
	var ticketEvent model.TicketEvent
	if err := c.ShouldBindJSON(&ticketEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketEvent(&ticketEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketEvent,
	})
}

// GetAllTicketEvents
func (tc *TicketEventController) GetAllTicketEvents(c *gin.Context) {
	ticketEvents, err := tc.Service.GetAllTicketEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketEvents,
	})
}

// GetTicketEventByID
func (tc *TicketEventController) GetTicketEventByID(c *gin.Context) {
	id := c.Param("id")
	ticketEvent, err := tc.Service.GetTicketEventByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket event not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketEvent,
	})
}

// UpdateTicketEvent - Perbaiki untuk match dengan service yang butuh ID terpisah
func (tc *TicketEventController) UpdateTicketEvent(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketEvent

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Panggil service dengan ID terpisah
	if err := tc.Service.UpdateTicketEvent(id, &updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketEvent, err := tc.Service.GetTicketEventByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket event",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketEvent,
	})
}

// DeleteTicketEvent
func (tc *TicketEventController) DeleteTicketEvent(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketEvent(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket event deleted",
	})
}
