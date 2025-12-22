package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketHistoryController struct {
	Service *service.TicketHistoryService
}

func NewTicketHistoryController(s *service.TicketHistoryService) *TicketHistoryController {
	return &TicketHistoryController{Service: s}
}

// CreateTicketHistory
func (tc *TicketHistoryController) CreateTicketHistory(c *gin.Context) {
	var ticketHistory model.TicketHistory
	if err := c.ShouldBindJSON(&ticketHistory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketHistory(&ticketHistory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketHistory,
	})
}

// GetAllTicketHistories
func (tc *TicketHistoryController) GetAllTicketHistories(c *gin.Context) {
	ticketHistories, err := tc.Service.GetAllTicketHistories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketHistories,
	})
}

// GetTicketHistoryByID
func (tc *TicketHistoryController) GetTicketHistoryByID(c *gin.Context) {
	id := c.Param("id")
	ticketHistory, err := tc.Service.GetTicketHistoryByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketHistory == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket history not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketHistory,
	})
}

// UpdateTicketHistory
func (tc *TicketHistoryController) UpdateTicketHistory(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketHistory
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTicketHistory(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketHistory, err := tc.Service.GetTicketHistoryByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketHistory,
	})
}

// DeleteTicketHistory
func (tc *TicketHistoryController) DeleteTicketHistory(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketHistory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket history deleted",
	})
}